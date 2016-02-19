package bot

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/nlopes/slack"
)

// processMessage handles incomming messages
func (bot *Bot) processMessage(message *slack.MessageEvent) {
	user := message.User
	// Handle Violators
	_, isViolator := bot.violatorUserMap[user]
	if isViolator {
		bot.violatorMessage(message, user)
	}

	botCalled := strings.HasPrefix(
		message.Text,
		fmt.Sprintf("<@%s>", bot.Slack.GetSelfID()),
	)
	if botCalled { // Messages made directly to bot
		switch {
		case bot.isMasterUser(user):
			{
				bot.masterMessages(message)
			}
		default:
			{
				bot.niceMessage(message, user)
			}
		}
	} else {
		switch {
		// Messages that contain the word tick
		case strings.Contains(message.Text, "tick"):
			{
				bot.Slack.SendMessage(bot.Slack.NewOutgoingMessage(
					"tock", message.Channel,
				))
			}
		// Messages that references the bot will be send out 30% of the time
		case strings.Contains(message.Text, fmt.Sprintf("<@%s>", bot.Slack.GetSelfID())):
			{
				randomInt := rand.Intn(10)
				if randomInt >= 7 {
					bot.Slack.SendMessage(bot.Slack.NewOutgoingMessage(
						bot.MessageRepo.Nice.GenerateMessage(user),
						message.Channel,
					))
				}
			}
		}
	}
}

// violatorMessage has the message for a late user
func (bot *Bot) violatorMessage(message *slack.MessageEvent, user string) {
	var returnMessage string
	// Check if user is still late
	if bot.isLateUser(user) {
		returnMessage = bot.MessageRepo.Angry.GenerateMessage(user)
	} else {
		returnMessage = fmt.Sprintf(
			"<@%s>, I was about to yell at you, but then I realized you actually filled out your timesheet. Thanks! ^_^",
			user,
		)
	}
	delete(bot.violatorUserMap, user)
	bot.Slack.SendMessage(bot.Slack.NewOutgoingMessage(returnMessage, message.Channel))
}

// masterMessages contains the commands for admins
func (bot *Bot) masterMessages(message *slack.MessageEvent) {
	var returnMessage string
	switch {
	case strings.Contains(message.Text, "slap users"):
		{
			go bot.SlapLateUsers()
			returnMessage = "Slapping Users!"
		}
	case strings.Contains(message.Text, "bother users"):
		{
			bot.startviolatorUserMapUpdater()
			returnMessage = "Starting to bother users!"
		}
	case strings.Contains(message.Text, "who is late?"):
		{
			lateList, total := bot.fetchLateUsers()
			returnMessage = fmt.Sprintf("%s are late! %d people total.", lateList, total)
		}
	default:
		{
			returnMessage = fmt.Sprintf(
				"Commands:\n Message tardy users `<@%s>: slap users!`\nBother tardy users `<@%s>: bother users!`\nFind out who is late `<@%s>: who is late?`",
				bot.Slack.GetSelfID(),
				bot.Slack.GetSelfID(),
				bot.Slack.GetSelfID(),
			)
		}
	}
	bot.Slack.SendMessage(bot.Slack.NewOutgoingMessage(
		returnMessage, message.Channel,
	))
}

// niceMessage are commands for user who are not late
func (bot *Bot) niceMessage(message *slack.MessageEvent, user string) {
	var returnMessage string
	switch {
	case strings.Contains(message.Text, "hello"):
		{
			bot.Slack.SendMessage(bot.Slack.NewOutgoingMessage(
				bot.MessageRepo.Nice.GenerateMessage(user),
				message.Channel,
			))
		}
	case strings.Contains(message.Text, "status"):
		{
			go func() {
				if bot.isLateUser(user) {
					returnMessage = fmt.Sprintf("<@%s>, you're late -_-", user)
				} else {
					returnMessage = fmt.Sprintf("<@%s>, you're on time! ^_^", user)
				}
				bot.Slack.SendMessage(bot.Slack.NewOutgoingMessage(
					returnMessage, message.Channel,
				))
			}()
		}
	}
}
