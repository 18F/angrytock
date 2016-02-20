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

type replaceChannel struct {
	newState map[string]string
	response chan bool
}

type readChannel struct {
	key      string
	response chan string
}

type deleteChannel struct {
	key      string
	response chan bool
}

type writeChannel struct {
	key      string
	value    string
	response chan bool
}

type mapStateManager struct {
	state          map[string]string
	readChannel    chan *readChannel
	writeChannel   chan *writeChannel
	deleteChannel  chan *deleteChannel
	replaceChannel chan *replaceChannel
}

func newMapStateManager() *mapStateManager {
	return &mapStateManager{
		make(map[string]string),
		make(chan *readChannel),
		make(chan *writeChannel),
		make(chan *deleteChannel),
		make(chan *replaceChannel),
	}
}

func initMapStateManager() *mapStateManager {
	manager := newMapStateManager()
	go func() {
		for {
			select {
			case read := <-manager.readChannel:
				read.response <- manager.state[read.key]
			case write := <-manager.writeChannel:
				manager.state[write.key] = write.value
				write.response <- true
			case remove := <-manager.deleteChannel:
				delete(manager.state, remove.key)
				remove.response <- true

			case replace := <-manager.replaceChannel:
				manager.state = replace.newState
				replace.response <- true
			}
		}
	}()
	return manager
}

// Bot struct serves as the primary entry point for slack and tock api methods
// It stores the slack token string and a database connection for storing
// emails and usernames
type Bot struct {
	UserEmailMap    *mapStateManager
	Slack           *slackPackage.Slack
	Tock            *tockPackage.Tock
	MessageRepo     *messagesPackage.MessageRepository
	violatorUserMap *mapStateManager
	masterList      []string
}

// InitBot method initalizes a bot
func InitBot() *Bot {
	userEmailMap := initMapStateManager()
	violatorUserMap := initMapStateManager()
	masterList := strings.Split(os.Getenv("MASTER_LIST"), ",")
	slack := slackPackage.InitSlack()
	tock := tockPackage.InitTock()
	messageRepo := messagesPackage.InitMessageRepository()

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
	writeToUserEmailMap := &writeChannel{
		key:      "",
		value:    "",
		response: make(chan bool),
	}
	defer close(writeToUserEmailMap.response)
	slackUsers := bot.Slack.FetchSlackUsers()
	for _, user := range slackUsers {
		if strings.HasSuffix(user.Profile.Email, ".gov") {
			writeToUserEmailMap.key = user.Profile.Email
			writeToUserEmailMap.value = user.ID
			bot.UserEmailMap.writeChannel <- writeToUserEmailMap
			succeeded := <-writeToUserEmailMap.response
			if succeeded != true {
				log.Println("User Add Failed")
			}
			bot.updateMasterList(user.Profile.Email, user.ID)
		}
	}
}

// updateviolatorUserMap generates a new map containing the slack id and user email
// of late tock users
func (bot *Bot) updateviolatorUserMap() {
	readUserEmailMap := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(readUserEmailMap.response)
	violatorUserMap := make(map[string]string)
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			readUserEmailMap.key = user.Email
			bot.UserEmailMap.readChannel <- readUserEmailMap
			userID := <-readUserEmailMap.response
			if user.Email != "" && userID != "" {
				violatorUserMap[userID] = user.Email
			}
		},
	)
	replaceviolatorUserMap := &replaceChannel{
		newState: violatorUserMap,
		response: make(chan bool),
	}
	defer close(replaceviolatorUserMap.response)
	bot.violatorUserMap.replaceChannel <- replaceviolatorUserMap
	<-replaceviolatorUserMap.response
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
				violatorUserMap := make(map[string]string)
				replaceviolatorUserMap := &replaceChannel{
					newState: violatorUserMap,
					response: make(chan bool),
				}
				defer close(replaceviolatorUserMap.response)
				bot.violatorUserMap.replaceChannel <- replaceviolatorUserMap
				replaced := <-replaceviolatorUserMap.response
				if replaced {
					log.Println("Destroying violator list")
				}
				ticker.Stop()
				return
			}
		}
	}()
}

// SlapLateUsers collects users from tock and looks for thier slack ids in a database
func (bot *Bot) SlapLateUsers() {
	log.Println("Slapping Tock Users")
	readUserEmailMap := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(readUserEmailMap.response)
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			readUserEmailMap.key = user.Email
			bot.UserEmailMap.readChannel <- readUserEmailMap
			userID := <-readUserEmailMap.response
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
	readViolatorMap := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(readViolatorMap.response)
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
				bot.processMessage(event, readViolatorMap)
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
	readUserEmailMap := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(readUserEmailMap.response)
	found := false
	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			readUserEmailMap.key = user.Email
			bot.UserEmailMap.readChannel <- readUserEmailMap
			userID := <-readUserEmailMap.response
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
	readUserEmailMap := &readChannel{
		key:      "",
		response: make(chan string),
	}
	defer close(readUserEmailMap.response)

	bot.Tock.UserApplier(
		func(user tockPackage.User) {
			readUserEmailMap.key = user.Email
			bot.UserEmailMap.readChannel <- readUserEmailMap
			slackUserID := <-readUserEmailMap.response
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
