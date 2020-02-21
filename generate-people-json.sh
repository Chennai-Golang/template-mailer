#!/usr/bin/env bash

_invite_file_path="past-mails/2020/02-invites.json"
_reminder_file_path="past-mails/2020/02-reminders.json"
_existing_mails_file="mails/existing.json"
_current_rsvp_file="mails/current-rsvp.json"

# Combine invite and reminder list from previous meetups
jq -s '[ .[0] + .[1] | group_by(.Email | ascii_downcase)[] | add ]' $_invite_file_path $_reminder_file_path | \
  jq '[ .[] | .Email |= ascii_downcase ]' | \
  jq '[ sort_by(.Email) | .[] | .IsRSVPed = false ]' > $_existing_mails_file

# Sort files before comparing - DESTRUCTIVE!
cat $_existing_mails_file | jq '[ .[] | .Email |= ascii_downcase ]' | jq 'sort_by(.Email)' > "$_existing_mails_file-sorted"
mv "$_existing_mails_file-sorted" $_existing_mails_file

cat $_current_rsvp_file | jq '[ .[] | .Email |= ascii_downcase ]' | jq 'sort_by(.Email)' > "$_current_rsvp_file-sorted"
mv "$_current_rsvp_file-sorted" $_current_rsvp_file

# Check for duplicates
jq -r '.[] | .Email' $_existing_mails_file > /tmp/uniq-mail-check-file1.txt
jq -r '.[] | .Email' $_current_rsvp_file > /tmp/uniq-mail-check-file2.txt

# Duplicate emails
echo "------Duplicate emails!------"
comm -1 -2 /tmp/uniq-mail-check-file1.txt /tmp/uniq-mail-check-file2.txt
echo "-----------------------------"

# Remove duplicates from file1!
comm -1 -2 /tmp/uniq-mail-check-file1.txt /tmp/uniq-mail-check-file2.txt | xargs -n 1 -- ./dedup.sh $_existing_mails_file

# Combine file1 with file2 contents
jq -s '.[0] + .[1]' $_existing_mails_file $_current_rsvp_file | tee people.json
