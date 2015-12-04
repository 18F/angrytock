# angrytock
[![Go Report Card](http://goreportcard.com/badge/18F/angrytock)](http://goreportcard.com/report/18F/angrytock)

A slack bot for "reminding" [Tock](https://github.com/18F/tock) users who are late filling out their timecards.

## Bot Master User Commands
`@botname: slap users!` : Reminds users to fill in their time sheets one time.  
`@botname: bother users!` : Searches for users writing in Slack and tells them to fill in their time sheets. Will only bother user 1 time and is only active for 30 minutes.  

## Regular interactions
Will respond with spicy comment about time if user starts messages with bot name. (ie `@botname: hi`) Mentions will not trigger responses.  

## Running tests
`go test ./... -cover `

## Deployment

### Env Variables
Set the following env variables  
```
export SLACK_KEY="<<Slack Key>>"
export TOCK_URL="https://Tock-Domain.com"
export MASTER_LIST=<<EMAIL>>,<<EMAIL>>
export PORT=5000 # will be set automatically by Cloud Foundry
```

### Deploying to Cloud Foundry
`cf push`
