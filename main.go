package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

type Person struct {
	Name  string
	Email string
}

type Credential struct {
	Email    string
	Password string
}

var peopleFile string
var templateFile string
var persons []Person
var credential Credential

var subject string

func initPersonsAndCreds() {
	jsonBytes, _ := ioutil.ReadFile(peopleFile)
	json.Unmarshal(jsonBytes, &persons)

	jsonBytes, _ = ioutil.ReadFile("./secrets/credentials.json")
	err := json.Unmarshal(jsonBytes, &credential)

	if err != nil {
		log.Panic(err)
	}
}

func send(email string, body string) {
	from := credential.Email
	pass := credential.Password
	to := email

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Panic("smtp error: %s", err)
		return
	}
}

func main() {
	peopleFile = os.Args[1]
	templateFile = os.Args[2]
	subject = os.Args[3]

	initPersonsAndCreds()
	txtFile, _ := ioutil.ReadFile(templateFile)

	tmpl, err := template.New("test").Parse(string(txtFile))

	if err != nil {
		log.Panic(err)
	}

	for _, person := range persons {
		var body bytes.Buffer
		err = tmpl.Execute(&body, person)

		send(person.Email, body.String())
	}

	if err != nil {
		log.Panic(err)
	}
}
