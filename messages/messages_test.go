package messagesPackage

import (
	"strings"
	"testing"
)

var messageRepo = InitMessageRepository()

// Check that messages render properly with user name
func TestAngryMessagesMessages(t *testing.T) {
	for _, message := range messageRepo.Angry.Messages {
		if !strings.Contains(message, "<@%s>") {
			t.Errorf(message)
		}
	}
}

// Check that messages render properly with user name
func TestNiceMessagesMessages(t *testing.T) {
	for _, message := range messageRepo.Nice.Messages {
		if !strings.Contains(message, "<@%s>") {
			t.Errorf(message)
		}
	}
}

// Check that messages include a %s
func TestReminderMessagesMessages(t *testing.T) {
	for _, message := range messageRepo.Reminder.Messages {
		if !strings.Contains(message, "%s") {
			t.Errorf(message)
		}
	}
}

// Check that randomly generated messages work
func TestRandomMessages(t *testing.T) {
	message := messageRepo.Angry.fetchRandomMessage()
	if !strings.Contains(message, "<@%s>") {
		t.Errorf(message)
	}

}
