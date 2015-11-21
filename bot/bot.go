// Package bot provides an interface accessing the tock and slack apis
// The primary purpose of this packages is to collect users from tock
// who have not filled out thier time forms and use the slack api to message them.
package bot

import (
	"crypto"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"

	"github.com/18F/hmacauth"
	"github.com/boltdb/bolt"
)

// Bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	ID            string
	Connection    *websocket.Conn
	Token         string
	DB            *bolt.DB
	AuditEndpoint string
	Auth          hmacauth.HmacAuth
}

// Open url and return the body of request
func fetchData(URL string) []byte {

	res, err := http.Get(URL)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	return body

}

// InitBot method initalizes a bot
func InitBot() *Bot {

	// Collect the slack key
	slackKey := os.Getenv("SLACK_KEY")
	if slackKey == "" {
		log.Fatal("SLACK_KEY environment variable not found")
	}

	// Start a connection to the websocket
	ws, id := NewSlackConnection(slackKey)

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

	// Get Audit endpoint
	auditendpoint := os.Getenv("AUDIT_ENDPOINT")
	if auditendpoint == "" {
		log.Fatal("AUDIT_ENDPOINT environment variable not found")
	}
	// Collect HMAC secret
	HMACSecret := os.Getenv("HMAC_SECRET")
	if HMACSecret == "" {
		log.Fatal("HMAC_SECRET environment variable not found")
	}
	auth := hmacauth.NewHmacAuth(crypto.SHA1, []byte(HMACSecret), "X-Signature", nil)

	return &Bot{id, ws, slackKey, db, auditendpoint, auth}
}

// SlapLateUsers collects users from tock and looks for thier slack ids in a database
func (bot *Bot) SlapLateUsers() {
	log.Println("Slapping Tock Users")
	data := bot.FetchTockUsers()
	bot.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range data.Users {
			v := string(b.Get([]byte(user.Email)))
			if v != "" {
				bot.MessageUser(v, "Please fill out your time sheet!")
			}
		}
		return nil
	})
}

// StoreSlackUsers is a method for collecting and storing slack users in database
func (bot *Bot) StoreSlackUsers() {
	log.Println("Collecting Slack Users")
	slackUserData := bot.FetchSlackUsers()
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
	data := bot.FetchTockUsers()
	bot.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SlackUsers"))
		for _, user := range data.Users {
			v := string(b.Get([]byte(user.Email)))
			userMap[v] = user.Email
		}
		return nil
	})
	return userMap
}

// BotherSlackUsers is methods that slacks offending tock users when they type
// write in slack
func (bot *Bot) BotherSlackUsers() {
	userMap := make(map[string]string)

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
		message, err := bot.GetMessage()
		if err != nil {
			log.Fatal(err)
		}
		// Check if the user is an offending user
		_, ok := userMap[message.User]
		// If the user is an offending user message them, but remove them off the list
		if message.Type == "message" && ok == true {
			message.Text = fmt.Sprintf(
				"<@%s>! So you have time for slack, but not tock hu?!",
				message.User,
			)
			bot.PostMessage(message)
			delete(userMap, message.User)
		}

	}

}
