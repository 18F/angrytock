// The bot package provides an interface accessing the tock and slack apis
// The primary purpose of this packages is to collect users from tock
// who have not filled out thier time forms and use the slack api to message them.
package bot

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
)

// The bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	Token string
	DB    *bolt.DB
}

// Open url and return the body of request
func fetchData(Url string) []byte {

	res, err := http.Get(Url)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	return body

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
				bot.MessageUser(v, "Please fill out your time sheet!")
			}
		}
		return nil
	})
}

// Function for collecting and storing slack users in database
func (bot *Bot) StoreSlackUsers() {

	slackUserData := bot.FetchSlackUsers()

	bot.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range slackUserData.Users {
			err := b.Put([]byte(user.Profile.Email), []byte(user.ID))
			log.Print("Saved :", user.Profile.Email)
			if err != nil {
				log.Print(err)
			}
		}
		return nil
	})
}
