package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type UserList struct {
	Ok      bool `json:"ok"`
	Members []struct {
		ID       string      `json:"id"`
		TeamID   string      `json:"team_id"`
		Name     string      `json:"name"`
		Deleted  bool        `json:"deleted"`
		Status   interface{} `json:"status"`
		Color    string      `json:"color"`
		RealName string      `json:"real_name"`
		Tz       interface{} `json:"tz"`
		TzLabel  string      `json:"tz_label"`
		TzOffset int         `json:"tz_offset"`
		Profile  struct {
			BotID              string      `json:"bot_id"`
			Image24            string      `json:"image_24"`
			Image32            string      `json:"image_32"`
			Image48            string      `json:"image_48"`
			Image72            string      `json:"image_72"`
			Image192           string      `json:"image_192"`
			Image512           string      `json:"image_512"`
			Image1024          string      `json:"image_1024"`
			ImageOriginal      string      `json:"image_original"`
			RealName           string      `json:"real_name"`
			RealNameNormalized string      `json:"real_name_normalized"`
			Fields             interface{} `json:"fields"`
		} `json:"profile"`
		IsAdmin           bool `json:"is_admin"`
		IsOwner           bool `json:"is_owner"`
		IsPrimaryOwner    bool `json:"is_primary_owner"`
		IsRestricted      bool `json:"is_restricted"`
		IsUltraRestricted bool `json:"is_ultra_restricted"`
		IsBot             bool `json:"is_bot"`
	} `json:"members"`
}

func (bot *Bot) FetchSlackUsers() {

	var data UserList

	url := fmt.Sprintf("https://slack.com/api/users.list?token=%s", bot.Token)

	res, err := http.Get(url)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}

	fmt.Print(data)

}
