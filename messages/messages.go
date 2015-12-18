package messagesPackage

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// MessageArray is a torage for one type of message includes methods for choosing a
// messages from the bunch
type MessageArray struct {
	Messages []string `yaml:"responses"`
}

// fetchRandomMessage method for selecting a random message from the created messages
func (msgs MessageArray) fetchRandomMessage() string {
	return msgs.Messages[rand.Intn(len(msgs.Messages))]
}

func (msgs MessageArray) GenerateMessage(filler string) string {
	return fmt.Sprintf(
		msgs.fetchRandomMessage(),
		filler,
	)
}

// MessageRepository Contains the messages and methods for generating
// responses for the bot
type MessageRepository struct {
	Angry    *MessageArray `yaml:"AngryMessages"`
	Nice     *MessageArray `yaml:"NiceMessages"`
	Reminder *MessageArray `yaml:"ReminderMessages"`
}

// InitMessageRepository loads all of the data into a MessageRepository
// and a MessageRepository struct
func InitMessageRepository() *MessageRepository {
	var mrep MessageRepository
	messageFile := "messages.yaml"
	workingDir, _ := os.Getwd()
	if !strings.HasSuffix(workingDir, "messages") {
		messageFile = filepath.Join("messages", messageFile)
	}
	messageFile = filepath.Join(workingDir, messageFile)
	data, err := ioutil.ReadFile(messageFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if yaml.Unmarshal(data, &mrep) != nil {
		log.Fatalf("error: %v", err)
	}
	return &mrep
}
