package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/smtp"
	"text/template"
)

type Tarkan struct {
	Name  string
	Email string
}

type Credential struct {
	Email    string
	Password string
}

var tarkans []Tarkan
var credential Credential

const subject = "Initiating project newsletter"

func initTarkansAndCreds() {
	jsonBytes, _ := ioutil.ReadFile("./secrets/tarkans.json")
	json.Unmarshal(jsonBytes, &tarkans)

	jsonBytes, _ = ioutil.ReadFile("./secrets/credentials.json")
	json.Unmarshal(jsonBytes, &credential)
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
	initTarkansAndCreds()
	txtFile, _ := ioutil.ReadFile("./template.txt")

	tmpl, err := template.New("test").Parse(string(txtFile))

	if err != nil {
		panic(err)
	}

	for _, tarkan := range tarkans {
		var body bytes.Buffer
		err = tmpl.Execute(&body, tarkan)

		send(tarkan.Email, body.String())
	}

	if err != nil {
		panic(err)
	}
}
