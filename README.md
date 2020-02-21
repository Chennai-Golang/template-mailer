# Template Mailer

A simple text based template mailer using golang libs.

## Build

```bash
go build .
```

## Generate people.json

```bash
./generate-people-json.sh
```

## Run

```bash
./template-mailer {people.json} {templates/template.txt} "{subject}"
```

or

```bash
./send-email.sh
```

## TODO

* Switch to HTML based mail now that we have SendGrid
