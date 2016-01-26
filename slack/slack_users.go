package slackPackage

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// UserList is a struct representation of the slack user.list api JSON structure.
type UserList struct {
	Users []struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		Profile           struct {
			Email              string      `json:"email"`
			RealName           string      `json:"real_name"`
			RealNameNormalized string      `json:"real_name_normalized"`
		} `json:"profile"`
		RealName string      `json:"real_name"`
		TeamID   string      `json:"team_id"`
	} `json:"members"`
	Ok bool `json:"ok"`
}

// ChannelResponse is a struct representation of a response from im.open
type channelResponse struct {
	Ok          bool `json:"ok"`
	NoOp        bool `json:"no_op"`
	AlreadyOpen bool `json:"already_open"`
	Channel     struct {
		ID string
	} `json:"channel"`
}

// FetchSlackUsers fetches a list of slack users and saves thier user ids by
// emails in a database
func (slack *Slack) FetchSlackUsers() *UserList {
	// Get a list of users
	var data UserList
	body := slack.DataFetcher.FetchData("https://slack.com/api/users.list")
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}

// openChannel opens or gets a direct message channel to user
func (slack *Slack) openOrGetChannel(user string) *channelResponse {
	// Open channel
	var channelData channelResponse
	URL := fmt.Sprintf("https://slack.com/api/im.open?user=%s", user)
	body := slack.DataFetcher.FetchData(URL)
	err := json.Unmarshal(body, &channelData)
	if err != nil {
		log.Print(err)
	}
	return &channelData
}

// MessageUser is a method to send a message to users
func (slack *Slack) MessageUser(user string, message string) {
	// Open a chanel to the user
	channelData := slack.openOrGetChannel(user)
	// Send message to the user
	URL := "https://slack.com/api/chat.postMessage?"
	URL += fmt.Sprintf(
		"channel=%s&text=%s&as_user=true",
		channelData.Channel.ID,
		url.QueryEscape(message),
	)
	res, err := http.Get(slack.AddToken(URL))
	res.Body.Close()
	if err != nil {
		log.Print("Failed to make request")
	}
}
