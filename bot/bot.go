// Package bot provides an interface accessing the tock and slack apis
// The primary purpose of this packages is to collect users from tock
// who have not filled out thier time forms and use the slack api to message them.
package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/18F/angrytock/messages"
	"github.com/18F/angrytock/slack"
	"github.com/18F/angrytock/tock"
)

// Bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	UserEmailMap    map[string]string
	Slack           *slackPackage.Slack
	Tock            *tockPackage.Tock
	MessageRepo     *messagesPackage.MessageRepository
	violatorUserMap map[string]string
	masterList      []string
}

// InitBot method initalizes a bot
func InitBot() *Bot {
	userEmailMap := make(map[string]string)
	violatorUserMap := make(map[string]string)
	masterList := strings.Split(os.Getenv("MASTER_LIST"), ",")
	slack := slackPackage.InitSlack()
	tock := tockPackage.InitTock()
	messageRepo := messagesPackage.InitMessageRepository()

	return &Bot{userEmailMap, slack, tock, messageRepo, violatorUserMap, masterList}
}

// masterList checks if a user email is in the masterList and
// if it is, it will replace the users email with a slack id
func (bot *Bot) updateMasterList(userEmail string, userSlackID string) {
	for idx, masterUserEmail := range bot.masterList {
		if masterUserEmail == userEmail {
			bot.masterList[idx] = userSlackID
		}
	}
}

// StoreSlackUsers is a method for collecting and storing slack users in database
func (bot *Bot) StoreSlackUsers() {
	log.Println("Collecting Slack Users")
	slackUserData := bot.Slack.FetchSlackUsers()
	for _, user := range slackUserData.Users {
		if user.Profile.Email != "" {
			log.Println("Saved:", user.Profile.Email)
			bot.UserEmailMap[user.Profile.Email] = user.ID
			bot.updateMasterList(user.Profile.Email, user.ID)
		}
	}
}

// updateviolatorUserMap generates a new map containing the slack id and user email
// of late tock users
func (bot *Bot) updateviolatorUserMap() {
	violatorUserMap := make(map[string]string)
	data := bot.Tock.FetchTockUsers()
	for _, user := range data.Users {
		userID := bot.UserEmailMap[user.Email]
		if user.Email != "" && userID != "" {
			violatorUserMap[userID] = user.Email
		}
	}
	bot.violatorUserMap = violatorUserMap
}

// startviolatorUserMapUpdater begins a ticker that only keeps the
// violator list full for 30 minutes
func (bot *Bot) startviolatorUserMapUpdater() {
	// Collect user data
	bot.updateviolatorUserMap()
	// Create a ticker to renew the cache of tock users
	ticker := time.NewTicker(30 * time.Minute)
	// Start the go channel
	go func() {
		for {
			select {
			case <-ticker.C:
				bot.violatorUserMap = make(map[string]string)
				log.Println("Destroying violator list")
				ticker.Stop()
				return
			}
		}
	}()
}

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
			if strings.Contains(message.Text, "slap users!") {
				go bot.SlapLateUsers()
				message.Text = "Slapping Users!"
			} else if strings.Contains(message.Text, "bother users!") {
				bot.startviolatorUserMapUpdater()
				message.Text = "Starting to bother users!"
			} else {
				message.Text = fmt.Sprintf(
					"Commands:\n Message tardy users `<@%s>: slap users!`\nBother tardy users `<@%s>: bother users!`",
					bot.Slack.ID,
					bot.Slack.ID,
				)
			}
			bot.Slack.PostMessage(message)
		}
	// If the user is an offending user message them, but remove them off the list
	case isViolator:
		{
			log.Println(message.Text)
			message.Text = bot.MessageRepo.GenerateAngryMessage(user)
			bot.Slack.PostMessage(message)
			delete(bot.violatorUserMap, user)
		}
		// If this is a message directed at the bot respond in a nice way
	case strings.HasPrefix(message.Text, fmt.Sprintf("<@%s>", bot.Slack.ID)):
		{
			log.Println(message.Text)
			message.Text = bot.MessageRepo.GenerateNiceMessage(user)
			bot.Slack.PostMessage(message)
		}
	}
}

// SlapLateUsers collects users from tock and looks for thier slack ids in a database
func (bot *Bot) SlapLateUsers() {
	log.Println("Slapping Tock Users")
	data := bot.Tock.FetchTockUsers()
	for _, user := range data.Users {
		userID := bot.UserEmailMap[user.Email]
		bot.Slack.MessageUser(userID, bot.MessageRepo.GenerateReminderMessages())
	}
}

// ListenToSlackUsers starts a loop that listens to tock users
func (bot *Bot) ListenToSlackUsers() {
	log.Println("Listening to slack")
	// Creating a for loop to catch channel messages from slack
	for {
		// Get each incoming message
		message, _ := bot.Slack.GetMessage()
		// Only process messages
		if message.Type == "message" {
			bot.processMessage(message)
		}
	}
}
