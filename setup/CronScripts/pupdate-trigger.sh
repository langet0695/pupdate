#!/bin/bash
jwt=$(curl -X "POST" http://admin:<your-api-password>@localhost:8080/createToken)
echo "Collected JWT"
clean_jwt=$(echo "$jwt" | sed 's/"//g')
header="Authorization: Bearer ${clean_jwt//[$'\t\r\n ']}"
echo "Sending Messages"
mailing=$(curl -X "POST" -H "$header" http://localhost:8080/mail)
echo "$mailing"
