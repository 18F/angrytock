# angrytock
[![Go Report Card](http://goreportcard.com/badge/18F/angrytock)](http://goreportcard.com/report/18F/angrytock)

A slack bot for "reminding" [Tock](https://github.com/18F/tock) users who are late filling out their timecards.

## Usage
Set the following env variables
```
export SLACK_KEY="<<Slack Key>>"
export TOCK_URL="https://Tock-Domain.com"
export PORT=5000 # will be set automatically by cloud foundry
```

## Deployment
`cf push`
