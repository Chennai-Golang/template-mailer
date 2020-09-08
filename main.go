package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

// Person represents a meetup invitee
type Person struct {
	Name     string
	Email    string
	IsRSVPed bool
}

// Credential represents Mail credentials
type Credential struct {
	Email    string
	Password string
}

// Mailer wraps sendgrid and templating
type Mailer struct {
	*sendgrid.Client
	*template.Template
	Subject string
	Credential
	People []Person
}

var mailer Mailer

func getCredential() (*Credential, error) {
	var credential Credential

	jsonBytes, err := ioutil.ReadFile("./secrets/credentials.json")

	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	err = json.Unmarshal(jsonBytes, &credential)
	if err != nil {
		return nil, fmt.Errorf("malformed credentials: %v", err)
	}

	return &credential, nil
}

func getPeople(peopleFile string) ([]Person, error) {
	var people []Person
	jsonBytes, err := ioutil.ReadFile(peopleFile)

	if err != nil {
		return nil, fmt.Errorf("unable to read people list: %v", err)
	}

	err = json.Unmarshal(jsonBytes, &people)

	if err != nil {
		return nil, fmt.Errorf("malformed people list: %v", err)
	}

	return people, nil
}

func getTemplate(templateFile string) (*template.Template, error) {
	txtFile, err := ioutil.ReadFile(templateFile)

	if err != nil {
		return nil, fmt.Errorf("unable to read template file: %v", err)
	}

	tmpl, err := template.New("test").Parse(string(txtFile))

	if err != nil {
		return nil, fmt.Errorf("unable to parse template file: %v", err)
	}

	return tmpl, nil
}

func init() {
	peopleFile := os.Args[1]
	templateFile := os.Args[2]
	subject := os.Args[3]

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	credential, err := getCredential()
	if err != nil {
		log.Panic(err)
	}

	people, err := getPeople(peopleFile)
	if err != nil {
		log.Panic(err)
	}

	tmpl, err := getTemplate(templateFile)
	if err != nil {
		log.Panic(err)
	}

	mailer = Mailer{
		Client:     client,
		Template:   tmpl,
		Credential: *credential,
		Subject:    subject,
		People:     people,
	}
}

// Send dispatches the mail to people
func (m Mailer) Send() []error {
	var errs []error

	for index, person := range m.People {
		from := mail.NewEmail("Chennai Gophers", m.Credential.Email)
		subject := m.Subject
		to := mail.NewEmail(person.Name, person.Email)

		var body bytes.Buffer
		err := m.Template.Execute(&body, person)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		plainTextContent := body.String()

		log.Debug(plainTextContent)

		message := mail.NewSingleEmail(from, subject, to, "", "")

		message.Content = []*mail.Content{mail.NewContent("text/plain", plainTextContent)}

		response, err := m.Client.Send(message)

		log.Infof("\tStatus Code: %d", response.StatusCode)
		log.Infof("\tBody: %s", response.Body)

		if !(response.StatusCode < 400) {
			resErr := fmt.Errorf("unable to send email index: %d ; to: %v ; body: %s ; error: %v", index, person, response.Body, err)
			log.Error(resErr)
			errs = append(errs, resErr)
		} else if err != nil {
			log.Error("unable to send email:", err)
			errs = append(errs, err)
		}
	}

	return errs
}

func main() {
	log.SetLevel(log.DebugLevel)

	errs := mailer.Send()

	if len(errs) != 0 {
		log.Error("unable to send all emails!")
	}
}
