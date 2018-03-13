# AngryTock

[![Go Report Card](http://goreportcard.com/badge/18F/angrytock)](http://goreportcard.com/report/18F/angrytock)

A slack bot for "reminding" [Tock](https://github.com/18F/tock) users who are late filling out their timecards.

This bot is dependent on Tock. Please refer to the appropriate documentation for
both projects.

- [Tock documentation](https://github.com/18F/tock/tree/master/docs)
- [AngryTock documentation](https://github.com/18F/angrytock/tree/master/docs)

## Bot Master User Commands
`@botname: slap users!` : Reminds users to fill in their time sheets one time.
`@botname: bother users!` : Searches for users writing in Slack and tells them to fill in their time sheets. Will only bother user 1 time and is only active for 30 minutes.
`@botname: who is late?` : Returns a list of users who are late.


## Regular interactions
`@botname: status` : Will check the Tock API and tell the user if they have filled out their timesheet.
`@botname: say something` : Will respond to the use with a message about time.

## Running tests
`go test ./... -cover `

## Deployment

### Env Variables
Set the following env variables
```
export SLACK_KEY="<<Slack Key>>"
export TOCK_URL="https://Tock-Domain.com"
export TOCK_USER_URL="https://Tock-Domain.com/users"
export MASTER_LIST=<<EMAIL>>,<<EMAIL>>
export PORT=5000 # will be set automatically by Cloud Foundry
```

### Deploying to Cloud Foundry
`cf push`

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for additional information.

### Branch flow

- Main branch: `master`

## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.
