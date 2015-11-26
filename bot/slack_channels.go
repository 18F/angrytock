// These methods are going to be depricated because bots cannot join channels
// by themselves.

package bot

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

func (bot *Bot) fetchSlackChannels() *channelsResponse {
	var data channelsResponse
	URL := fmt.Sprintf("https://slack.com/api/channels.list?token=%s", bot.Token)
	body := FetchData(URL)
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}

func (bot *Bot) joinChannel(channelName string) {
	URL := fmt.Sprintf("https://slack.com/api/channels.join?token=%s&name=%s", bot.Token, channelName)
	body := FetchData(URL)
	log.Println(string(body))
}

// JoinAllChannels is a method to make the bot join all the channels
func (bot *Bot) JoinAllChannels() {
	channels := bot.fetchSlackChannels()
	for _, channel := range channels.Channels {
		bot.joinChannel(channel.Name)
	}
}
