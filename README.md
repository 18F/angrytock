# tock-bot

Slackbot slapbot for slapping tock slackers

## Usage
Set the following env variables
```
export SLACK_KEY="<<Slack Key>>"
export AUDIT_ENDPOINT="http://Tock-Domain/api/reporting_period_audit/"
export PORT=5000 # will be set automatically by cloud foundry
```

## Deployment
`cf push`
