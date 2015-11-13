// The bot package provides an interface accessing the tock and slack apis
// The primary purpose of this packages is to collect users from tock
// who have not filled out thier time forms and use the slack api to message them.
package bot

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// The bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	Token string
	DB *bolt.DB
}

// Method for initializing a bot
func InitBot() *Bot {

	// Collect the slack key
	slack_key := os.Getenv("SLACK_KEY")
	if slack_key == "" {
			log.Fatal("Slack key not found")
	}

	// Open connection to database
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
			log.Fatal(err)
	}

	// Create a database
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("SlackUsers"))
		if err != nil {
				log.Fatal(err)
		}
		return nil
	})

	return &Bot{slack_key, db}
}

// Collects users from tock and looks for thier slack ids in a database
func (bot *Bot) SlapLateUsers() {

	data := FetchTockUsers()
	bot.DB.View(func(tx *bolt.Tx) error {
		for _, user := range data.Users {
			b := tx.Bucket([]byte("SlackUsers"))
			v := string(b.Get([]byte(user.Email)))
			if v != "" {
				log.Print(v)
			}

		}
		return nil
	})
}

// Function for collecting and storing slack users in database
func (bot *Bot) StoreSlackUsers() {

	data := bot.FetchSlackUsers()

	bot.DB.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range data.Users {
			err := b.Put([]byte(user.Profile.Email), []byte(user.ID))
			log.Print("Saved :", user.Profile.Email)
			if err != nil {
					log.Print(err)
			}
		}
		return nil
	})
}
