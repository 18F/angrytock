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

func (bot *Bot) masterMessages(message slackPackage.Message) {

	if strings.Contains(message.Text, "slap users!") {
		go bot.SlapLateUsers()
		message.Text = "Slapping Users!"
	} else if strings.Contains(message.Text, "bother users!") {
		bot.startviolatorUserMapUpdater()
		message.Text = "Starting to bother users!"
	} else if strings.Contains(message.Text, "who is late?") {
		lateList, total := bot.fetchLateUsers()
		message.Text = fmt.Sprintf("%s are late! %d people total.", lateList, total)
	} else {
		message.Text = fmt.Sprintf(
			"Commands:\n Message tardy users `<@%s>: slap users!`\nBother tardy users `<@%s>: bother users!`\nFind out who is late `<@%s>: who is late?`",
			bot.Slack.ID,
			bot.Slack.ID,
			bot.Slack.ID,
		)
	}
	bot.Slack.PostMessage(message)

}

func (bot *Bot) violatorMessage(message slackPackage.Message, user string) {
	message.Text = bot.MessageRepo.Angry.GenerateMessage(user)
	bot.Slack.PostMessage(message)
	delete(bot.violatorUserMap, user)

}

func (bot *Bot) niceMessage(message slackPackage.Message, user string) {

	if strings.Contains(message.Text, "say something") {
		message.Text = bot.MessageRepo.Nice.GenerateMessage(user)
		bot.Slack.PostMessage(message)
	} else if strings.Contains(message.Text, "status") {
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
