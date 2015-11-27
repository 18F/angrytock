/*
Package containing methods for connecting to slack api, reading messages,
and responding.

Code "forked" from `https://github.com/rapidloop/mybot`
*/

package slackPackage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

var counter uint64

// These two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.
type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	URL   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	ID string `json:"id"`
}

// Message struct is the message structure that reads and writes to and
// from the real time messaging websocket socket
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
}

// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func start(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	var respObj responseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.URL
	id = respObj.Self.ID
	return
}

// NewSlackConnection Starts a websocket-based Real Time API session and return
//the websocket and the ID of the (bot-)user whom the token belongs to.
func NewSlackConnection(slackKey string) (*websocket.Conn, string) {
	// Check if the keys exist
	if slackKey == "" {
		fmt.Fprintf(os.Stderr, "SLACK_KEY missing from env")
		os.Exit(1)
	}
	wsurl, id, err := start(slackKey)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}
	return ws, id
}

// GetMessage is a method for getting messages from the websocket connection
func (slack *Slack) GetMessage() (m Message, err error) {
	err = websocket.JSON.Receive(slack.Connection, &m)
	return
}

// PostMessage is a method for posting messages back into the slack
func (slack *Slack) PostMessage(m Message) error {
	m.ID = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(slack.Connection, m)
}
