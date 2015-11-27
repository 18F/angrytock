// Package bot provides an interface accessing the tock and slack apis
// The primary purpose of this packages is to collect users from tock
// who have not filled out thier time forms and use the slack api to message them.
package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/geramirez/tock-bot/slack"
	"github.com/geramirez/tock-bot/tock"
)

// Bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	DB    *bolt.DB
	Slack *slackPackage.Slack
	Tock  *tockPackage.Tock
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

	return &Bot{db, slack, tock}
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

func (bot *Bot) createUserMap() map[string]string {
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
	return userMap
}

// BotherSlackUsers is methods that slacks offending tock users when they type
// write in slack
func (bot *Bot) BotherSlackUsers() {
	log.Println("Bothering Tock Users")
	// Collect user data
	userMap := make(map[string]string)
	userMap = bot.createUserMap()

	// Create a ticker to renew the cache of tock users
	ticker := time.NewTicker(20 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Fetching Tock Users")
				userMap = bot.createUserMap()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		// Get each incoming message
		message, err := bot.Slack.GetMessage()
		if err != nil {
			log.Print(err)
		}
		log.Print(message)

		// Check if the user is an offending user
		_, ok := userMap[message.User]
		// If the user is an offending user message them, but remove them off the list
		if message.Type == "message" && ok == true {
			message.Text = fmt.Sprintf(
				"<@%s>! So you have time for slack, but not tock, huh?!",
				message.User,
			)
			bot.Slack.PostMessage(message)
			delete(userMap, message.User)
		}

	}

}
