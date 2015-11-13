package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Struct representation of the slack user.list api JSON structure.
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

// Fetches a list of slack users and saves thier user ids by emails in a database
func (bot *Bot) FetchSlackUsers() *UserList {

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
	return &data

}
