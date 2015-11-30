package messagesPackage

import (
	"fmt"
	"strings"
	"testing"
)

var messageRepo = InitMessageRepository()

// Check that messages render properly with user name
func TestAngryMessagesMessages(t *testing.T) {
	for _, message := range messageRepo.AngryMessages.Messages {
		renderedMessage := fmt.Sprintf(message, "testUser")
		if strings.Contains(renderedMessage, "<@%!s(MISSING)>") {
			t.Errorf(message)
		}
	}
}

// Check that messages render properly with user name
func TestNiceMessagesMessages(t *testing.T) {
	for _, message := range messageRepo.NiceMessages.Messages {
		renderedMessage := fmt.Sprintf(message, "testUser")
		if strings.Contains(renderedMessage, "<@%!s(MISSING)>") {
			t.Errorf(message)
		}
	}
}

// Check that messages do not include username
func TestReminderMessagesMessages(t *testing.T) {
	for _, message := range messageRepo.ReminderMessages.Messages {
		if strings.Contains(message, "%s") {
			t.Errorf(message)
		}
	}
}

// Check that randomly generated messages work
func TestRandomMessages(t *testing.T) {
	message := messageRepo.AngryMessages.fetchRandomMessage()
	if !strings.Contains(message, "<@%s>") {
		t.Errorf(message)
	}

}
