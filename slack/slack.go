package slackPackage

import (
	"log"
	"net/url"
	"os"

	"github.com/18F/angrytock/helpers"
	"golang.org/x/net/websocket"
)

// Slack struct stores the data and websocket connection for slack rti
type Slack struct {
	ID          string
	Token       string
	Connection  *websocket.Conn
	DataFetcher *helpers.DataFetcher
}

// addToken add the slack token to the url
func addToken(Token string, URL string) string {

	parsedURL, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	// Add the token
	query := parsedURL.Query()
	query.Add("token", Token)
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}

// A closure for initalizing a function that add url token to url strings
// and then returns the data
func dataFetcherClosure(token string) func(URL string) []byte {
	return func(URL string) []byte {
		return helpers.FetchData(addToken(token, URL))
	}
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

	// Initalize a new data fetcher
	dataFetcher := helpers.NewDataFetcher(dataFetcherClosure(key))

	return &Slack{id, key, ws, dataFetcher}
}

// AddToken is the external facing function for adding tokens
func (slack *Slack) AddToken(URL string) string {
	return addToken(slack.Token, URL)
}
