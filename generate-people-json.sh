#!/usr/bin/env bash

_invite_file_path="past-mails/2019/1912-invite.json"
_reminder_file_path="past-mails/2019/191213-reminder.json"

# Combine invite and reminder list from previous meetups
jq -s '[ .[0] + .[1] | group_by(.Email | ascii_downcase)[] | add ]' $_invite_file_path $_reminder_file_path | \
  jq '[ .[] | .Email |= ascii_downcase ]' | \
  jq '[ sort_by(.Email) | .[] | .IsRSVPed = false ]' > mails/existing.json

_file1_path="mails/existing.json"
_file2_path="mails/current-rsvp.json"

# Sort files before comparing - DESTRUCTIVE!
cat $_file1_path | jq '[ .[] | .Email |= ascii_downcase ]' | jq 'sort_by(.Email)' > "$_file1_path-sorted"
mv "$_file1_path-sorted" $_file1_path

cat $_file2_path | jq '[ .[] | .Email |= ascii_downcase ]' | jq 'sort_by(.Email)' > "$_file2_path-sorted"
mv "$_file2_path-sorted" $_file2_path

# Check for duplicates
jq -r '.[] | .Email' $_file1_path > /tmp/uniq-mail-check-file1.txt
jq -r '.[] | .Email' $_file2_path > /tmp/uniq-mail-check-file2.txt

# Duplicate emails
echo "------Duplicate emails!------"
comm -1 -2 /tmp/uniq-mail-check-file1.txt /tmp/uniq-mail-check-file2.txt
echo "-----------------------------"

# Remove duplicates from file1!
comm -1 -2 /tmp/uniq-mail-check-file1.txt /tmp/uniq-mail-check-file2.txt | xargs -n 1 -- ./dedup.sh $_file1_path

# Combine file1 with file2 contents
jq -s '.[0] + .[1]' $_file1_path $_file2_path | tee people.json
