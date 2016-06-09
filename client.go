// Copyright 2016 Will Charczuk. Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package slack is a event driven client for the popular Slack chat application.

// A trivial example is:
//  package main
//  import (
//      "fmt"
//      "os"
//      "github.com/wcharczuk/go-slack"
//  )
//
//  func main() {
//      client := slack.NewClient(os.Getenv("SLACK_TOKEN"))
//      client.AddEventListener(slack.EventHello, func(m *slack.Message, c *slack.Client) {
//          fmt.Println("connected")
//      })
//      client.AddEventListener(slack.EventMessage, func(m *slack.Message, c *slack.Client) {
//          fmt.Prinln("message received!")
//      })
//      session, err := client.Connect() //session has the current users list and channel list
//      if err != nil {
//          fmt.Printf("%v\n", err)
//          os.Exit(1)
//      }
//  }
// The client has two phases of initialization, NewClient and Start.

package slack

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/blendlabs/go-exception"
	"github.com/gorilla/websocket"
)

// Client constants
const (
	// DefaultPingInterval is the ping interval in seconds.
	DefaultPingInterval = 30 * time.Second
)

// EventListener is a function that recieves messages from a client.
type EventListener func(client *Client, message *Message)

// --------------------------------------------------------------------------------
// Connect(token string) is the main slack entry point.
// It returns a `*Client`, which is the root struct for interacting with slack.
// Register "Listeners" to hook into incoming messages and filter out ones you
// dont' care about.
// --------------------------------------------------------------------------------

// NewClient creates a Client with a given token.
func NewClient(token string) *Client {
	c := &Client{
		Token:          token,
		EventListeners: map[Event][]EventListener{},
		ActiveChannels: []string{},
		pingInFlight:   map[int64]time.Time{},
		pingInterval:   DefaultPingInterval,
	}
	c.AddEventListener(EventChannelJoined, c.handleChannelJoined)
	c.AddEventListener(EventChannelDeleted, c.handleChannelDeleted)
	c.AddEventListener(EventChannelUnArchive, c.handleChannelUnarchive)
	c.AddEventListener(EventChannelLeft, c.handleChannelLeft)
	c.AddEventListener(EventPong, c.handlePong)
	return c
}

// Client is the mechanism with which the package consumer interacts with Slack.
type Client struct {
	Token          string
	EventListeners map[Event][]EventListener
	ActiveChannels []string

	activeLock       sync.Mutex
	socketConnection *websocket.Conn

	pingInFlight     map[int64]time.Time
	pingInFlightLock sync.Mutex
	pingInterval     time.Duration
}

// AddEventListener attaches a new Listener to the given event.
// There can be multiple listeners to an event.
// If an event is already being listened for, calling Listen will add a new listener to that event.
func (rtm *Client) AddEventListener(event Event, handler EventListener) {
	if listeners, handlesEvent := rtm.EventListeners[event]; handlesEvent {
		rtm.EventListeners[event] = append(listeners, handler)
	} else {
		rtm.EventListeners[event] = []EventListener{handler}
	}
}

// RemoveEventListener removes all listeners for an event.
func (rtm *Client) RemoveEventListener(event Event) {
	delete(rtm.EventListeners, event)
}

// Connect be4gins a session with Slack.
func (rtm *Client) Connect() (*Session, error) {
	res := Session{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/rtm.start").
		WithPostData("token", rtm.Token).
		WithPostData("no_unreads", "true").
		WithPostData("mpim_aware", "true").
		FetchJSONToObject(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	//start socket connection
	u, err := url.Parse(res.URL)
	if err != nil {
		return nil, err
	}

	rtm.socketConnection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		return nil, err
	}

	go func() { //fetch the (initial) channel list asyncronously
		rtm.activeLock.Lock()
		defer rtm.activeLock.Unlock()

		channels, chanelsErr := rtm.ChannelsList(true) //excludeArchived == true
		if chanelsErr != nil {
			return
		}

		for x := 0; x < len(channels); x++ {
			channel := channels[x]
			if channel.IsMember && !channel.IsArchived {
				rtm.ActiveChannels = append(rtm.ActiveChannels, channel.ID)
			}
		}
	}()

	go func() {
		rtm.pingLoop()
	}()

	go func() {
		rtm.listenLoop()
	}()

	return &res, nil
}

// Stop closes the connection with Slack.
func (rtm *Client) Stop() error {
	if rtm.socketConnection == nil {
		return nil
	}

	closeErr := rtm.socketConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if closeErr != nil {
		return closeErr
	}
	rtm.socketConnection.Close()
	rtm.socketConnection = nil
	return nil
}

// SendMessage sends a basic message over the open web socket connection to slack.
func (rtm *Client) SendMessage(m *Message) error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	return rtm.socketConnection.WriteJSON(m)
}

// Say sends a basic message to a given channelID.
func (rtm *Client) Say(channelID string, messageComponents ...interface{}) error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	m := &Message{Type: "message", Text: fmt.Sprint(messageComponents...), Channel: channelID}
	return rtm.SendMessage(m)
}

// Sayf is an overload that uses Printf style replacements for a basic message to a given channelID.
func (rtm *Client) Sayf(channelID, format string, messageComponents ...interface{}) error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	m := &Message{Type: "message", Text: fmt.Sprintf(format, messageComponents...), Channel: channelID}
	return rtm.SendMessage(m)
}

// Ping sends a special type of "ping" message to Slack to remind it to keep the connection open.
// Currently unused internally by Slack.
func (rtm *Client) Ping() error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	p := &Message{ID: time.Now().UTC().UnixNano(), Type: "ping"}
	rtm.pingInFlight[p.ID] = time.Now().UTC()
	return rtm.socketConnection.WriteJSON(p)
}

//--------------------------------------------------------------------------------
// INTERNAL METHODS
//--------------------------------------------------------------------------------

func (rtm *Client) pingLoop() error {
	var err error
	for rtm.socketConnection != nil {

		err = rtm.Ping()
		if err != nil {
			err = rtm.cycleConnection()
			if err != nil {

				break
			}
		}

		time.Sleep(rtm.pingInterval)
	}

	return nil
}

func (rtm *Client) handlePong(client *Client, message *Message) {

}

func (rtm *Client) cycleConnection() error {
	res := Session{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/rtm.start").
		WithPostData("token", rtm.Token).
		WithPostData("no_unreads", "true").
		WithPostData("mpim_aware", "true").
		FetchJSONToObject(&res)

	if err != nil {
		return err
	}

	u, err := url.Parse(res.URL)
	if err != nil {
		return err
	}
	rtm.socketConnection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	return err
}

func (rtm *Client) listenLoop() (err error) {
	defer func() {
		if err != nil {
			fmt.Printf("Slack :: Exiting Listen Loop, err: %#v\n", err)
		}
	}()
	var messageBytes []byte
	var mt MessageType
	var cm ChannelJoinedMessage
	var m Message

	for {
		if rtm.socketConnection == nil {
			return nil
		}
		_, messageBytes, err = rtm.socketConnection.ReadMessage()
		if err != nil {
			return err
		}

		err = json.Unmarshal(messageBytes, &mt)
		if err == nil {
			switch mt.Type {
			case EventChannelJoined:
				{
					err = json.Unmarshal(messageBytes, &cm)
					if err == nil {
						rtm.dispatch(&Message{Type: EventChannelJoined, Channel: cm.Channel.ID})
					}
				}
			default:
				{
					err = json.Unmarshal(messageBytes, &m)
					if err == nil {
						rtm.dispatch(&m)
					}
				}
			}
		}
	}
}

func (rtm *Client) dispatch(m *Message) {
	var listener EventListener
	if listeners, hasListeners := rtm.EventListeners[m.Type]; hasListeners {
		for x := 0; x < len(listeners); x++ {
			listener = listeners[x]
			go func() {
				defer func() {
					if r := recover(); r != nil {
						//log the panic somewhere.
					}
				}()

				listener(rtm, m)
			}()
		}
	}
}

func (rtm *Client) handleChannelJoined(client *Client, message *Message) {
	rtm.activeLock.Lock()
	defer rtm.activeLock.Unlock()
	rtm.ActiveChannels = append(rtm.ActiveChannels, message.Channel)
}

func (rtm *Client) handleChannelUnarchive(client *Client, message *Message) {
	rtm.activeLock.Lock()
	defer rtm.activeLock.Unlock()

	channel, err := rtm.ChannelsInfo(message.Channel)
	if err != nil {
		return
	}
	if channel.IsMember {
		rtm.ActiveChannels = append(rtm.ActiveChannels, message.Channel)
	}
}

func (rtm *Client) handleChannelLeft(client *Client, message *Message) {
	rtm.removeActiveChannel(message.Channel)
}

func (rtm *Client) handleChannelDeleted(client *Client, message *Message) {
	rtm.removeActiveChannel(message.Channel)
}

func (rtm *Client) removeActiveChannel(channelID string) {
	rtm.activeLock.Lock()
	defer rtm.activeLock.Unlock()

	currentChannels := []string{}
	for x := 0; x < len(rtm.ActiveChannels); x++ {
		currentChannelID := rtm.ActiveChannels[x]
		if channelID != currentChannelID {
			currentChannels = append(currentChannels, currentChannelID)
		}
	}
	rtm.ActiveChannels = currentChannels
}
