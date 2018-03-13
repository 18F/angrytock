package slackPackage

import (
	"fmt"
	"log"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/nlopes/slack"
)

// Slack sturct extend the slackRTM method
type Slack struct {
	*slack.RTM
}

// InitSlack initalizes the struct object
func InitSlack() *Slack {
	appEnv, _ := cfenv.Current()
	appService, _ := appEnv.Services.WithName("angrytock-credentials")
	// Collect the slack key
	key := fmt.Sprint(appService.Credentials["SLACK_KEY"])
	if key == "" {
		log.Fatal("SLACK_KEY environment variable not found")
	}
	rtm := slack.New(key).NewRTM()
	return &Slack{rtm}
}

// FetchSlackUsers fetches a list of slack users and saves thier user ids by
// this method could use the GetInfo()
func (api *Slack) FetchSlackUsers() []slack.User {
	users, err := api.GetUsers()
	if err != nil {
		log.Println(err.Error())
	}
	return users
}

// GetSelfID returns the ID of the slack bot
func (api *Slack) GetSelfID() string {
	return api.GetInfo().User.ID
}

// MessageUser opens a channel to a user if it doesn't exist and messages the user
func (api *Slack) MessageUser(user string, message string) {
	_, _, channelID, err := api.Client.OpenIMChannel(user)
	if err != nil {
		log.Println("Unable to open channel")
	}
	// Can insert images an other things here
	postParams := slack.PostMessageParameters{}
	api.Client.PostMessage(channelID, message, postParams)

}
