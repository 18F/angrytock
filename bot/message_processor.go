package bot

import (
	"fmt"
	"strings"

	"github.com/18F/angrytock/slack"
)

// processMessage handles incomming messages
func (bot *Bot) processMessage(message slackPackage.Message) {

	// Check if the user is an offending user
	user := message.User
	_, isViolator := bot.violatorUserMap[user]
	// Check if user is in masterList
	var ismasterUser bool
	for _, masterUser := range bot.masterList {
		if masterUser == user {
			ismasterUser = true
		}
	}

	switch {
	// Master user commands
	case ismasterUser && strings.HasPrefix(message.Text, fmt.Sprintf("<@%s>", bot.Slack.ID)):
		{
			bot.masterMessages(message)
		}
	// If the user is an offending user message them, but remove them off the list
	case isViolator:
		{
			bot.violatorMessage(message, user)
		}
		// If this is a message directed at the bot respond in a nice way
	case strings.HasPrefix(message.Text, fmt.Sprintf("<@%s>", bot.Slack.ID)):
		{
			bot.niceMessage(message, user)
		}
	}
}

// masterMessages contains the commands for admins
func (bot *Bot) masterMessages(message slackPackage.Message) {
	switch {
	case strings.Contains(message.Text, "slap users"):
		{
			go bot.SlapLateUsers()
			message.Text = "Slapping Users!"
		}
	case strings.Contains(message.Text, "bother users"):
		{
			bot.startviolatorUserMapUpdater()
			message.Text = "Starting to bother users!"
		}
	case strings.Contains(message.Text, "who is late?"):
		{
			lateList, total := bot.fetchLateUsers()
			message.Text = fmt.Sprintf("%s are late! %d people total.", lateList, total)
		}
	default:
		{
			message.Text = fmt.Sprintf(
				"Commands:\n Message tardy users `<@%s>: slap users!`\nBother tardy users `<@%s>: bother users!`\nFind out who is late `<@%s>: who is late?`",
				bot.Slack.ID,
				bot.Slack.ID,
				bot.Slack.ID,
			)
		}
	}
	bot.Slack.PostMessage(message)

}

// violatorMessage has the message for a late user
func (bot *Bot) violatorMessage(message slackPackage.Message, user string) {
	// Check if user is still late
	if bot.isLateUser(user) {
		message.Text = bot.MessageRepo.Angry.GenerateMessage(user)
	} else {
		message.Text = fmt.Sprintf(
			"<@%s>, I was about to yell at you, but then I realized you actually filled outyour timesheet. Thanks! ^_^",
			user,
		)
	}
	delete(bot.violatorUserMap, user)
	bot.Slack.PostMessage(message)

}

// niceMessage are commands for user who are not late
func (bot *Bot) niceMessage(message slackPackage.Message, user string) {

	switch {
	case strings.Contains(message.Text, "say"):
		{
			message.Text = bot.MessageRepo.Nice.GenerateMessage(user)
			bot.Slack.PostMessage(message)
		}
	case strings.HasSuffix(message.Text, "ping"):
		{
			message.Text = "pong!"
			bot.Slack.PostMessage(message)
		}
	case strings.Contains(message.Text, "status"):
		{
			// Start a go func because the search may take a while
			go func() {
				if bot.isLateUser(user) {
					message.Text = fmt.Sprintf("<@%s>, you're late -_-", user)
				} else {
					message.Text = fmt.Sprintf("<@%s>, you're on time! ^_^", user)
				}
				bot.Slack.PostMessage(message)
			}()
		}
	}
}
