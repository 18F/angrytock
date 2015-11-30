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
		Color             string `json:"color"`
		Deleted           bool   `json:"deleted"`
		Has2fa            bool   `json:"has_2fa"`
		ID                string `json:"id"`
		IsAdmin           bool   `json:"is_admin"`
		IsBot             bool   `json:"is_bot"`
		IsOwner           bool   `json:"is_owner"`
		IsPrimaryOwner    bool   `json:"is_primary_owner"`
		IsRestricted      bool   `json:"is_restricted"`
		IsUltraRestricted bool   `json:"is_ultra_restricted"`
		Name              string `json:"name"`
		Profile           struct {
			Email              string      `json:"email"`
			Fields             interface{} `json:"fields"`
			Image192           string      `json:"image_192"`
			Image24            string      `json:"image_24"`
			Image32            string      `json:"image_32"`
			Image48            string      `json:"image_48"`
			Image512           string      `json:"image_512"`
			Image72            string      `json:"image_72"`
			RealName           string      `json:"real_name"`
			RealNameNormalized string      `json:"real_name_normalized"`
		} `json:"profile"`
		RealName string      `json:"real_name"`
		Status   interface{} `json:"status"`
		TeamID   string      `json:"team_id"`
		Tz       string      `json:"tz"`
		TzLabel  string      `json:"tz_label"`
		TzOffset int         `json:"tz_offset"`
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
	_, err := http.Get(slack.AddToken(URL))
	if err != nil {
		log.Print("Failed to make request")
	}
}
