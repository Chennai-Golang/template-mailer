#!/usr/bin/env bash

set -e

go build .

## Invite mail test
# ./template-mailer test.json templates/invite-template.txt "Go February 2020 meetup"

## Invite mail actual
# ./template-mailer people.json templates/invite-template.txt "Go February 2020 meetup"


## Reminder mail test
./template-mailer test.json templates/reminder-template.txt "See you tomorrow at Go meetup - Feb 2020 Edition"

## Reminder mail actual
# ./template-mailer mails/current-rsvp.json templates/reminder-template.txt "See you tomorrow at Go meetup - Feb 2020 Edition"
