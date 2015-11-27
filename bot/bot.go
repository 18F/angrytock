// Package bot provides an interface accessing the tock and slack apis
// The primary purpose of this packages is to collect users from tock
// who have not filled out thier time forms and use the slack api to message them.
package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/geramirez/tock-bot/slack"
	"github.com/geramirez/tock-bot/tock"
)

// Bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	DB      *bolt.DB
	Slack   *slackPackage.Slack
	Tock    *tockPackage.Tock
	userMap map[string]string
}

// InitBot method initalizes a bot
func InitBot() *Bot {

	// Open connection to database
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Create a database bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("SlackUsers"))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	slack := slackPackage.InitSlack()
	tock := tockPackage.InitTock()
	userMap := make(map[string]string)

	return &Bot{db, slack, tock, userMap}
}

// StoreSlackUsers is a method for collecting and storing slack users in database
func (bot *Bot) StoreSlackUsers() {
	log.Println("Collecting Slack Users")
	slackUserData := bot.Slack.FetchSlackUsers()
	bot.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range slackUserData.Users {
			if user.Profile.Email != "" {
				err := b.Put([]byte(user.Profile.Email), []byte(user.ID))
				log.Print("Saved :", user.Profile.Email)
				if err != nil {
					log.Print(err)
				}
			}
		}
		return nil
	})
}

// updateUserMap generates a new map containing the slack id and user email
// of late tock users
func (bot *Bot) updateUserMap() {
	userMap := make(map[string]string)
	data := bot.Tock.FetchTockUsers()
	bot.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range data.Users {

			v := string(b.Get([]byte(user.Email)))
			if user.Email != "" && v != "" {
				userMap[v] = user.Email
			}
		}
		return nil
	})
	bot.userMap = userMap
}

// startUserMapUpdater begins a ticker to update the user map every thirty minutes
func (bot *Bot) startUserMapUpdater() {
	// Collect user data
	bot.updateUserMap()
	// Create a ticker to renew the cache of tock users
	ticker := time.NewTicker(30 * time.Minute)
	// Start the go channel
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Rebuilding userMap")
				bot.updateUserMap()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// processMessage handles incomming messages
func (bot *Bot) processMessage(message slackPackage.Message) {

	// Check if the user is an offending user
	_, isInMap := bot.userMap[message.User]

	switch {
	case isInMap:
		{
			// If the user is an offending user message them, but remove them off the list
			message.Text = fmt.Sprintf(
				"<@%s>! So you have time for slack, but not tock, huh?!",
				message.User,
			)
			bot.Slack.PostMessage(message)
			delete(bot.userMap, message.User)
		}

	case strings.Contains(message.Text, fmt.Sprintf("<@%s>", bot.Slack.ID)):
		{
			message.Text = "Random Quote"
			bot.Slack.PostMessage(message)
		}
	}

}

// SlapLateUsers collects users from tock and looks for thier slack ids in a database
func (bot *Bot) SlapLateUsers() {
	log.Println("Slapping Tock Users")
	data := bot.Tock.FetchTockUsers()
	bot.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range data.Users {
			v := string(b.Get([]byte(user.Email)))
			bot.Slack.MessageUser(v, "Please fill out your time sheet!")
		}
		return nil
	})
}

// BotherSlackUsers is methods that slacks offending tock users when they type
// write in slack
func (bot *Bot) BotherSlackUsers() {
	log.Println("Bothering Tock Users")
	// Starting the map updater channel/method
	bot.startUserMapUpdater()
	// Creating a for loop to catch channel messages from slack
	for {
		// Get each incoming message
		message, err := bot.Slack.GetMessage()
		if err != nil {
			log.Print(err, message)
		}
		// Only process messages
		if message.Type == "message" {
			bot.processMessage(message)
		}

	}
}
