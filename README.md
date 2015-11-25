# tock-bot

[![Go Report Card](http://goreportcard.com/badge/geramirez/tock-bot)](http://goreportcard.com/report/geramirez/tock-bot)

Slackbot slapbot for slapping tock slackers

## Usage
Set the following env variables
```
export SLACK_KEY="<<Slack Key>>"
export AUDIT_ENDPOINT="https://Tock-Domain/api/reporting_period_audit/"
export PORT=5000 # will be set automatically by cloud foundry
```

## Deployment
`cf push`
