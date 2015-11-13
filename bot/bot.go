package bot

import (
	"log"
	"os"
)

type Bot struct {
	Token string
}

// Function for initializing the bot
func InitBot() *Bot {
	return &Bot{os.Getenv("SLACK_KEY")}
}

func (bot *Bot) SlapLateUsers() {
	data := FetchTockUsers()
	for _, user := range data.Users {
		log.Println(user)
	}
}
