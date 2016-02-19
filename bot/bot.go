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
	"github.com/nlopes/slack"
)

type readChannel struct {
	key      string
	response chan string
}

type writeChannel struct {
	key      string
	value    string
	response chan bool
}

type mapStateManager struct {
	state        map[string]string
	readChannel  chan *readChannel
	writeChannel chan *writeChannel
}

func newMapStateManager() *mapStateManager {
	return &mapStateManager{
		make(map[string]string),
		make(chan *readChannel),
		make(chan *writeChannel),
	}
}

func (manager *mapStateManager) startMapStateManager() {
	go func() {
		for {
			select {
			case read := <-manager.readChannel:
				read.response <- manager.state[read.key]
			case write := <-manager.writeChannel:
				log.Println(write)
				manager.state[write.key] = write.value
				write.response <- true
			}
		}
	}()
}

// Bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	UserEmailMap    *mapStateManager
	Slack           *slackPackage.Slack
	Tock            *tockPackage.Tock
	MessageRepo     *messagesPackage.MessageRepository
	violatorUserMap map[string]string
	masterList      []string
}

// InitBot method initalizes a bot
func InitBot() *Bot {
	userEmailMap := newMapStateManager()
	violatorUserMap := make(map[string]string)
	masterList := strings.Split(os.Getenv("MASTER_LIST"), ",")
	slack := slackPackage.InitSlack()
	tock := tockPackage.InitTock()
	messageRepo := messagesPackage.InitMessageRepository()

	// Start the map manager
	userEmailMap.startMapStateManager()

	return &Bot{userEmailMap, slack, tock, messageRepo, violatorUserMap, masterList}
}

// Check if user is in masterList
func (bot *Bot) isMasterUser(user string) bool {
	var ismasterUser bool
	for _, masterUser := range bot.masterList {
		if masterUser == user {
			ismasterUser = true
			break
		}
	}
	return ismasterUser
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
	// Open a write channel to the bot
	writes := &writeChannel{
		key:      "",
		value:    "",
		response: make(chan bool),
	}
	defer close(writes.response)
	slackUsers := bot.Slack.FetchSlackUsers()
	for _, user := range slackUsers {
		if strings.HasSuffix(user.Profile.Email, ".gov") {
			writes.key = user.Profile.Email
			writes.value = user.ID
			bot.UserEmailMap.writeChannel <- writes
			<-writes.response
			bot.updateMasterList(user.Profile.Email, user.ID)
		}
	}
}

// updateviolatorUserMap generates a new map containing the slack id and user email
// of late tock users
func (bot *Bot) updateviolatorUserMap() {
	reads := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(reads.response)

	violatorUserMap := make(map[string]string)
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			reads.key = user.Email
			bot.UserEmailMap.readChannel <- reads
			userID := <-reads.response
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
	reads := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(reads.response)
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			reads.key = user.Email
			bot.UserEmailMap.readChannel <- reads
			userID := <-reads.response
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
	go bot.Slack.ManageConnection()
	// Creating a for loop to catch channel messages from slack
	for {
		select {
		case rtmEvent := <-bot.Slack.IncomingEvents:
			switch event := rtmEvent.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
			case *slack.ConnectedEvent:
				// Ignore PresenceChangeEvent
			case *slack.MessageEvent:
				bot.processMessage(event)
			case *slack.PresenceChangeEvent:
				// Ignore PresenceChangeEvent
			case *slack.LatencyReport:
				// Ignore LatencyReport
			case *slack.RTMError:
				// Show errors
				fmt.Printf("Error: %s\n", event.Error())
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break

			default:
				// Do nothing
			}
		}
	}
}

// isLateUser returns if the user is late.
func (bot *Bot) isLateUser(slackUserID string) bool {
	reads := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(reads.response)
	found := false
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			reads.key = user.Email
			bot.UserEmailMap.readChannel <- reads
			userID := <-reads.response
			if slackUserID == userID {
				found = true
			}
		},
	)
	return found
}

// fetchLateUsers returns a list of late users
func (bot *Bot) fetchLateUsers() (string, int) {
	var lateList string
	var counter int
	reads := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(reads.response)

	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			reads.key = user.Email
			bot.UserEmailMap.readChannel <- reads
			slackUserID := <-reads.response
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
