package slackPackage

import (
	"log"
	"net/url"
	"os"

	"github.com/geramirez/tock-bot/helpers"
	"golang.org/x/net/websocket"
)

// Slack struct stores the data and websocket connection for slack rti
type Slack struct {
	ID         string
	Token      string
	Connection *websocket.Conn
}

// InitSlack initalizes the struct object
func InitSlack() *Slack {
	// Collect the slack key
	key := os.Getenv("SLACK_KEY")
	if key == "" {
		log.Fatal("SLACK_KEY environment variable not found")
	}
	// Start a connection to the websocket
	ws, id := NewSlackConnection(key)
	// Return pointer to slack struct
	return &Slack{id, key, ws}
}

// AddToken add the slack token to the url
func (slack *Slack) AddToken(URL string) string {

	parsedURL, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	// Add the token
	query := parsedURL.Query()
	query.Add("token", slack.Token)
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()

}

// FetchData adds a token to the url before fetching data
func (slack *Slack) FetchData(URL string) []byte {
	return helpers.FetchData(slack.AddToken(URL))
}
