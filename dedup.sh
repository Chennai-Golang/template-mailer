#!/usr/bin/env bash

_FILE_PATH=$1
_EMAIL_TO_REMOVE=$2

jq --arg email $_EMAIL_TO_REMOVE 'del( .[] | select(.Email | ascii_downcase == $email ))' $_FILE_PATH > "$_FILE_PATH-removed"

mv "$_FILE_PATH-removed" $_FILE_PATH
