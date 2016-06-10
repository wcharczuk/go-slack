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
	"net/http"
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

	// DefaultPingTimeout is the timeout for in flight pings.
	DefaultPingTimeout = 30 * time.Second

	// DefaultPingMaxFails is the number of ping failures before cycling.
	DefaultPingMaxFails = 5

	// DefaultPingMaxInFlight is the maximum number of pings in flight.
	DefaultPingMaxInFlight = 5
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
		Token:           token,
		EventListeners:  map[Event][]EventListener{},
		ActiveChannels:  []string{},
		isDebug:         true,
		pingTimeout:     DefaultPingTimeout,
		pingMaxInFlight: DefaultPingMaxInFlight,
		pingMaxFails:    DefaultPingMaxFails,
		pingInFlight:    map[int64]time.Time{},
		pingInterval:    DefaultPingInterval,
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

	pingTimeout      time.Duration
	pingMaxInFlight  int
	pingMaxFails     int
	pingFails        int
	pingInFlight     map[int64]time.Time
	pingInFlightLock sync.Mutex
	pingInterval     time.Duration

	isDebug bool
}

// SetDebug turns on debug logging.
func (rtm *Client) SetDebug(value bool) {
	rtm.isDebug = value
}

// AddEventListener attaches a new Listener to the given event.
// There can be multiple listeners to an event.
// If an event is already being listened for, calling Listen will add a new listener to that event.
func (rtm *Client) AddEventListener(event Event, handler EventListener) {
	rtm.EventListeners[event] = append(rtm.EventListeners[event], handler)
}

// RemoveEventListeners removes all listeners for an event.
func (rtm *Client) RemoveEventListeners(event Event) {
	delete(rtm.EventListeners, event)
}

// Connect be4gins a session with Slack.
func (rtm *Client) Connect() (*Session, error) {
	res := Session{}
	meta, err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/rtm.start").
		WithPostData("token", rtm.Token).
		WithPostData("no_unreads", "false").
		WithPostData("mpim_aware", "true").
		FetchJSONToObjectWithMeta(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if meta.StatusCode > http.StatusOK {
		return exception.New("Non-200 Status from Slack, aborting.")
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

	// asynchronously fetch active channels.
	go rtm.fetchActiveChannels()

	// ping slack every N seconds to make sure the connection is still active.
	go rtm.pingLoop()

	// listen for messages.
	go rtm.listenLoop()

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
	rtm.dispatch(p)
	rtm.pingInFlight[p.ID] = time.Now().UTC()
	return rtm.socketConnection.WriteJSON(p)
}

//--------------------------------------------------------------------------------
// INTERNAL METHODS
//--------------------------------------------------------------------------------

func (rtm *Client) pingLoop() error {
	var err error
	for rtm.socketConnection != nil {
		err = rtm.doPing()
		if err != nil {
			break
		}
		time.Sleep(rtm.pingInterval)
	}
	return nil
}

func (rtm *Client) doPing() error {
	rtm.pingInFlightLock.Lock()
	defer rtm.pingInFlightLock.Unlock()

	var err error
	if len(rtm.pingInFlight) < rtm.pingMaxInFlight {
		err = rtm.Ping()
		if err != nil {
			rtm.logf("ping error, cycling connection: %v\n", err)
			err = rtm.cycleConnection()
			if err != nil {
				rtm.logf("error cycling connection: %v\n", err)
			}
		}
	}

	now := time.Now().UTC()
	for _, v := range rtm.pingInFlight {
		if now.Sub(v) >= rtm.pingTimeout {
			err = rtm.cycleConnection()
			if err != nil {
				rtm.logf("error cycling connection: %v\n", err)
			}
		}
	}

	for k, v := range rtm.pingInFlight {
		if now.Sub(v) >= rtm.pingTimeout {
			delete(rtm.pingInFlight, k)
		}
	}
	return nil
}

func (rtm *Client) handlePong(client *Client, message *Message) {
	rtm.pingInFlightLock.Lock()
	defer rtm.pingInFlightLock.Unlock()
	delete(rtm.pingInFlight, message.ReplyTo)
}

func (rtm *Client) cycleConnection() error {
	res := Session{}
	meta, err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/rtm.start").
		WithPostData("token", rtm.Token).
		WithPostData("no_unreads", "true").
		WithPostData("mpim_aware", "true").
		FetchJSONToObjectWithMeta(&res)

	if err != nil {
		return err
	}

	if meta.StatusCode > http.StatusOK {
		return exception.New("Non-200 Status from Slack, aborting.")
	}

	rtm.pingInFlight = map[int64]time.Time{}
	rtm.pingFails = 0

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
			rtm.logf("exiting Listen Loop, err: %#v", err)
		}
	}()
	var mt MessageType
	var messageBytes []byte

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
			m := Message{}
			err = json.Unmarshal(messageBytes, &m)
			if err == nil {
				if len(mt.Type) == 0 && m.OK != nil { //special situation where acks don't have types and we have to sniff.
					rtm.dispatch(&Message{Type: EventMessageACK, ReplyTo: m.ReplyTo, Timestamp: m.Timestamp, Text: m.Text})
				} else {
					rtm.dispatch(&m)
				}
			}
		}
	}
}

func (rtm *Client) dispatch(m *Message) {
	if listeners, hasListeners := rtm.EventListeners[m.Type]; hasListeners {
		for index := range listeners {
			go func(listener EventListener) {
				defer func() {
					if r := recover(); r != nil {
						rtm.logf("go-slack: dispatch() fatal: %#v\n", r)
					}
				}()

				listener(rtm, m)
			}(listeners[index])
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

func (rtm *Client) fetchActiveChannels() {
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

func (rtm *Client) log(args ...interface{}) {
	if rtm.isDebug {
		fmt.Println(args...)
	}
}

func (rtm *Client) logf(format string, args ...interface{}) {
	if rtm.isDebug {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("%s\n", msg)
	}
}
