package bot

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
type ChannelResponse struct {
	Ok          bool `json:"ok"`
	NoOp        bool `json:"no_op"`
	AlreadyOpen bool `json:"already_open"`
	Channel     struct {
		ID string
	} `json:"channel"`
}

// FetchSlackUsers fetches a list of slack users and saves thier user ids by
// emails in a database
func (bot *Bot) FetchSlackUsers() *UserList {

	var data UserList

	URL := fmt.Sprintf("https://slack.com/api/users.list?token=%s", bot.Token)

	body := fetchData(URL)

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data

}

// MessageUser is a method to send a message to users
func (bot *Bot) MessageUser(user string, message string) {

	var channelData ChannelResponse

	// Open channel
	URL := "https://slack.com/api/im.open?"
	URL += fmt.Sprintf("token=%s&user=%s", bot.Token, user)

	body := fetchData(URL)

	err := json.Unmarshal(body, &channelData)
	if err != nil {
		log.Print(err)
	}

	// Send message
	URL = "https://slack.com/api/chat.postMessage?"
	URL += fmt.Sprintf(
		"token=%s&channel=%s&text=%s&as_user=true",
		bot.Token,
		channelData.Channel.ID,
		url.QueryEscape(message),
	)

	_, err = http.Get(URL)
	if err != nil {
		log.Print("Failed to make request")
	}
}
