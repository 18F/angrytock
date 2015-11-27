package messagesPackage

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// MessageArray is a torage for one type of message includes methods for choosing a
// messages from the bunch
type MessageArray struct {
	Messages []string `yaml:"responses"`
}

func (msgs MessageArray) fetchRandomMessage() string {
	return msgs.Messages[0]
}

// MessageRepository Contains the messages and methods for generating
// responses for the bot
type MessageRepository struct {
	AngryMessages MessageArray `yaml:"AngryMessages"`
	NiceMessages  MessageArray `yaml:"NiceMessages"`
}

// InitMessageRepository loads all of the data into a MessageRepository
// and a MessageRepository struct
func InitMessageRepository() *MessageRepository {
	var mrep MessageRepository
	data, err := ioutil.ReadFile("messages/messages.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if yaml.Unmarshal(data, &mrep) != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println("-------------", mrep)
	return &mrep
}

// GenerateAngryMessage creates an angry message
func (mrep *MessageRepository) GenerateAngryMessage(user string) string {
	return fmt.Sprintf(
		mrep.AngryMessages.fetchRandomMessage(),
		user,
	)
}

// GenerateNiceMessage creates a nice message
func (mrep *MessageRepository) GenerateNiceMessage(user string) string {
	return fmt.Sprintf(
		mrep.NiceMessages.fetchRandomMessage(),
		user,
	)
}
