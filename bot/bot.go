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
		if strings.HasSuffix(user.Profile.Email, ".gov") {
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
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			userID := bot.UserEmailMap[user.Email]
			if user.Email != "" && userID != "" {
				violatorUserMap[userID] = user.Email
			}
		},
	)
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

// SlapLateUsers collects users from tock and looks for thier slack ids in a database
func (bot *Bot) SlapLateUsers() {
	log.Println("Slapping Tock Users")
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			userID := bot.UserEmailMap[user.Email]
			if userID != "" {
				bot.Slack.MessageUser(
					userID,
					bot.MessageRepo.Reminder.GenerateMessage(bot.Tock.UserTockURL),
				)
			}
		},
	)
}

// ListenToSlackUsers starts a loop that listens to tock users
func (bot *Bot) ListenToSlackUsers() {
	log.Println("Listening to slack")
	// Creating a for loop to catch channel messages from slack
	for {
		// Get each incoming message
		message := bot.Slack.GetMessage()
		// Only process messages
		if message.Type == "message" {
			bot.processMessage(message)
		}
	}
}

// isLateUser returns if the user is late.
func (bot *Bot) isLateUser(slackUserID string) bool {
	found := false
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			if slackUserID == bot.UserEmailMap[user.Email] {
				found = true
			}
		},
	)
	return found
}

// fetchLateUsers returns a list of late users
func (bot *Bot) fetchLateUsers() (string, int) {
	var lateList string
	var slackUserID string
	var counter int

	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			slackUserID = bot.UserEmailMap[user.Email]
			if slackUserID != "" {
				lateList += fmt.Sprintf("<@%s>, ", slackUserID)
				counter++
			}
		},
	)

	if lateList == "" {
		lateList = "No people"
	}

	return lateList, counter
}
