// These methods are going to be depricated because bots cannot join channels
// by themselves.

package slack

import (
	"encoding/json"
	"fmt"
	"log"
)

type channelsResponse struct {
	Channels []struct {
		Created    int      `json:"created"`
		Creator    string   `json:"creator"`
		ID         string   `json:"id"`
		IsArchived bool     `json:"is_archived"`
		IsChannel  bool     `json:"is_channel"`
		IsGeneral  bool     `json:"is_general"`
		IsMember   bool     `json:"is_member"`
		Members    []string `json:"members"`
		Name       string   `json:"name"`
		NumMembers int      `json:"num_members"`
		Purpose    struct {
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
			Value   string `json:"value"`
		} `json:"purpose"`
		Topic struct {
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
			Value   string `json:"value"`
		} `json:"topic"`
	} `json:"channels"`
	Ok bool `json:"ok"`
}

func (slack *Slack) fetchSlackChannels() *channelsResponse {
	var data channelsResponse
	body := slack.FetchData("https://slack.com/api/channels.list")
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}

func (slack *Slack) joinChannel(channelName string) {
	URL := fmt.Sprintf("https://slack.com/api/channels.join?&name=%s", channelName)
	//TODO: finish the joining channel method
	_ = slack.FetchData(URL)
}

// JoinAllChannels is a method to make the bot join all the channels
func (slack *Slack) JoinAllChannels() {
	channels := slack.fetchSlackChannels()
	for _, channel := range channels.Channels {
		slack.joinChannel(channel.Name)
	}
}
