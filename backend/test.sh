#!/bin/bash
go build -o climate . &&
PORT=8701 JWT_KEY=114514 DB_USER=climate_sentinel DB_NAME=climate_sentinel DB_PASSWORD=climate_sentinel SMTP_MAIL=floodguard@126.com SMTP_PASSWORD=NTbKihEAKrNLcuJm SMTP_SERVER=smtp.126.com SMTP_PORT=25 ./climate
