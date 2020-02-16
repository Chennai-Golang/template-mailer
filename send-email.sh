#!/usr/bin/env bash

set -e

go build .

./template-mailer test.json templates/invite-template.txt "Go February 2020 meetup"

# ./template-mailer people.json templates/invite-template.txt "Go February 2020 meetup"

# ./template-mailer mails/current-rsvp.json templates/invite-template.txt "See you tomorrow at Go meetup"
