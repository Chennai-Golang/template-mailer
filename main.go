package main

import (
	"bytes"
	"encoding/json"
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

func init() {
	peopleFile := os.Args[1]
	templateFile := os.Args[2]
	subject := os.Args[3]

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	var credential Credential

	jsonBytes, _ := ioutil.ReadFile("./secrets/credentials.json")
	err := json.Unmarshal(jsonBytes, &credential)

	if err != nil {
		log.Panic(err)
	}

	var people []Person
	jsonBytes, _ = ioutil.ReadFile(peopleFile)
	err = json.Unmarshal(jsonBytes, &people)

	if err != nil {
		log.Panic(err)
	}

	txtFile, _ := ioutil.ReadFile(templateFile)

	tmpl, err := template.New("test").Parse(string(txtFile))

	if err != nil {
		log.Panic(err)
	}

	mailer = Mailer{
		Client:     client,
		Template:   tmpl,
		Credential: credential,
		Subject:    subject,
		People:     people,
	}
}

// Send dispatches the mail to people
func (m Mailer) Send() []error {
	var errs []error

	for _, person := range m.People {
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

		log.Info(plainTextContent)

		message := mail.NewSingleEmail(from, subject, to, plainTextContent, "")

		response, err := m.Client.Send(message)

		log.Debug(response.StatusCode, response.Body, response.Headers)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func main() {
	errs := mailer.Send()

	if len(errs) != 0 {
		log.Error(errs)
	}
}
