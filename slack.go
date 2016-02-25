package slack

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blendlabs/go-exception"
	"github.com/blendlabs/go-request"
	"github.com/gorilla/websocket"
)

// Event is a type alias for string to differentiate Slack event types.
type Event string

const (
	// APIScheme is the protocol used to communicate with slack.
	APIScheme = "https"
	// APIEndpoint is the host used to communicate with slack.
	APIEndpoint = "slack.com"

	// EventHello is an enumerated event.
	EventHello Event = "hello"
	// EventMessage is an enumerated event.
	EventMessage Event = "message"
	// EventUserTyping is an enumerated event.
	EventUserTyping Event = "user_typing"
	// EventChannelMarked is an enumerated event.
	EventChannelMarked Event = "channel_marked"
	// EventChannelJoined is an enumerated event.
	EventChannelJoined Event = "channel_joined"
	// EventChannelLeft is an enumerated event.
	EventChannelLeft Event = "channel_left"
	// EventChannelDeleted is an enumerated event.
	EventChannelDeleted Event = "channel_deleted"
	// EventChannelRename is an enumerated event.
	EventChannelRename Event = "channel_rename"
	// EventChannelArchive is an enumerated event.
	EventChannelArchive Event = "channel_archive"
	// EventChannelUnArchive is an enumerated event.
	EventChannelUnArchive Event = "channel_unarchive"
	// EventChannelHistoryChanged is an enumerated event.
	EventChannelHistoryChanged Event = "channel_history_changed"
	// EventDNDUpdated is an enumerated event.
	EventDNDUpdated Event = "dnd_updated"
	// EventDNDUpdatedUser is an enumerated event.
	EventDNDUpdatedUser Event = "dnd_updated_user"
	// EventIMCreated is an enumerated event.
	EventIMCreated Event = "im_created"
	// EventImOpen is an enumerated event.
	EventImOpen Event = "im_open"
	// EventImClose is an enumerated event.
	EventImClose Event = "im_close"
	// EventImMarked is an enumerated event.
	EventImMarked Event = "im_marked"
	// EventImHistoryChanged is an enumerated event.
	EventImHistoryChanged Event = "im_history_changed"
	// EventGroupJoined is an enumerated event.
	EventGroupJoined Event = "group_joined"
	// EventGroupLeft is an enumerated event.
	EventGroupLeft Event = "group_left"
	// EventGroupOpen is an enumerated event.
	EventGroupOpen Event = "group_open"
	// EventGroupClose is an enumerated event.
	EventGroupClose Event = "group_close"
	// EventGroupArchive is an enumerated event.
	EventGroupArchive Event = "group_archive"
	// EventGroupUnarchive is an enumerated event.
	EventGroupUnarchive Event = "group_unarchive"
	// EventGroupRename is an enumerated event.
	EventGroupRename Event = "group_rename"
	// EventGroupMarked is an enumerated event.
	EventGroupMarked Event = "group_marked"
	// EventGroupHistoryChanged is an enumerated event.
	EventGroupHistoryChanged Event = "group_history_changed"
	// EventFileCreated is an enumerated event.
	EventFileCreated Event = "file_created"
	// EventFileShared is an enumerated event.
	EventFileShared Event = "file_shared"
	// EventFileUnshared is an enumerated event.
	EventFileUnshared Event = "file_unshared"
	// EventFilePublic is an enumerated event.
	EventFilePublic Event = "file_public"
	// EventFilePrivate is an enumerated event.
	EventFilePrivate Event = "file_private"
	// EventFileChange is an enumerated event.
	EventFileChange Event = "file_change"
	// EventFileDeleted is an enumerated event.
	EventFileDeleted Event = "file_deleted"
	// EventFileCommentAdded is an enumerated event.
	EventFileCommentAdded Event = "file_comment_added"
	// EventFileCommentEdited is an enumerated event.
	EventFileCommentEdited Event = "file_comment_edited"
	// EventFileCommentDeleted is an enumerated event.
	EventFileCommentDeleted Event = "file_comment_deleted"
	// EventPinAdded is an enumerated event.
	EventPinAdded Event = "pin_added"
	// EventPinRemoved is an enumerated event.
	EventPinRemoved Event = "pin_removed"
	// EventPresenceChange is an enumerated event.
	EventPresenceChange Event = "presence_change"
	// EventManualPresenceChange is an enumerated event.
	EventManualPresenceChange Event = "manual_presence_change"
	// EventPrefChange is an enumerated event.
	EventPrefChange Event = "pref_change"
	// EventUserChange is an enumerated event.
	EventUserChange Event = "user_change"
	// EventTeamJoin is an enumerated event.
	EventTeamJoin Event = "team_join"
	// EventStarAdded is an enumerated event.
	EventStarAdded Event = "star_added"
	// EventStarRemoved is an enumerated event.
	EventStarRemoved Event = "star_removed"
	// EventReactionAdded is an enumerated event.
	EventReactionAdded Event = "reaction_added"
	// EventReactionRemoved is an enumerated event.
	EventReactionRemoved Event = "reaction_removed"
	// EventEmojiChanged is an enumerated event.
	EventEmojiChanged Event = "emoji_changed"
	// EventCommandsChanged is an enumerated event.
	EventCommandsChanged Event = "commands_changed"
	// EventTeamPlanChanged is an enumerated event.
	EventTeamPlanChanged Event = "team_plan_changed"
	// EventTeamPrefChanged is an enumerated event.
	EventTeamPrefChanged Event = "team_pref_changed"
	// EventEmailDomainChanged is an enumerated event.
	EventEmailDomainChanged Event = "email_domain_changed"
	// EventTeamProfileChange is an enumerated event.
	EventTeamProfileChange Event = "team_profile_change"
	// EventTeamProfileDelete is an enumerated event.
	EventTeamProfileDelete Event = "team_profile_delete"
	// EventTeamProfileReorder is an enumerated event.
	EventTeamProfileReorder Event = "team_profile_reorder"
	// EventBotAdded is an enumerated event.
	EventBotAdded Event = "bot_added"
	// EventBotChanged is an enumerated event.
	EventBotChanged Event = "bot_changed"
	// EventAccountsChanged is an enumerated event.
	EventAccountsChanged Event = "accounts_changed"
	// EventTeamMigrationStarted is an enumerated event.
	EventTeamMigrationStarted Event = "team_migration_started"

	// EventSubtypeBotMessage is an enumerated sub event.
	EventSubtypeBotMessage Event = "bot_message"
	// EventSubtypeMeMessage is an enumerated sub event.
	EventSubtypeMeMessage Event = "me_message"
	// EventSubtypeMessageChanged is an enumerated sub event.
	EventSubtypeMessageChanged Event = "message_changed"
	// EventSubtypeMessageDeleted is an enumerated sub event.
	EventSubtypeMessageDeleted Event = "message_deleted"
	// EventSubtypeChannelJoin is an enumerated sub event.
	EventSubtypeChannelJoin Event = "channel_join"
	// EventSubtypeChannelLeave is an enumerated sub event.
	EventSubtypeChannelLeave Event = "channel_leave"
	// EventSubtypeChannelTopic is an enumerated sub event.
	EventSubtypeChannelTopic Event = "channel_topic"
	// EventSubtypeChannelPurpose is an enumerated sub event.
	EventSubtypeChannelPurpose Event = "channel_purpose"
	// EventSubtypeChannelName is an enumerated sub event.
	EventSubtypeChannelName Event = "channel_name"
	// EventSubtypeChannelArchive is an enumerated sub event.
	EventSubtypeChannelArchive Event = "channel_archive"
	// EventSubtypeChannelUnarchive is an enumerated sub event.
	EventSubtypeChannelUnarchive Event = "channel_unarchive"
)

// User is the struct that represents a Slack user.
type User struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	Deleted           bool         `json:"deletd"`
	Color             string       `json:"color"`
	Profile           *UserProfile `json:"profile"`
	IsBot             bool         `json:"is_bot"`
	IsAdmin           bool         `json:"is_admin"`
	IsOwner           bool         `json:"is_owner"`
	IsPrimaryOwner    bool         `json:"is_primary_owner"`
	IsRestricted      bool         `json:"is_restricted"`
	IsUltraRestricted bool         `json:"is_ultra_restricted"`
	Has2FA            bool         `json:"has_2fa"`
	TwoFactorType     string       `json:"two_factor_type"`
	HasFiles          bool         `json:"has_files"`
}

// UserProfile represents additional information about a Slack user.
type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RealName  string `json:"real_name"`
	Email     string `json:"email"`
	Skype     string `json:"skype"`
	Phone     string `json:"phone"`
	Image24   string `json:"image_24"`
	Image32   string `json:"image_32"`
	Image48   string `json:"image_48"`
	Image72   string `json:"image_72"`
	Image192  string `json:"image_192"`
}

// Channel is the struct that represents a Slack channel.
type Channel struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	IsChannel          bool      `json:"is_channel"`
	Created            Timestamp `json:"created"`
	Creator            string    `json:"creator"`
	IsArchived         bool      `json:"is_archived"`
	IsGeneral          bool      `json:"is_general"`
	Members            []string  `json:"members"`
	Topic              *Topic    `json:"topic"`
	Purpose            *Topic    `json:"purpose"`
	IsMember           bool      `json:"is_member"`
	LastRead           Timestamp `json:"last_read"`
	UnreadCount        int       `json:"unread_count"`
	UnreadCountDisplay int       `json:"unread_count_display"`
	Latest             Message   `json:"latest"`
}

// Topic represents a Slack topic.
type Topic struct {
	Value   string    `json:"value"`
	Creator string    `json:"creator"`
	LastSet Timestamp `json:"last_set"`
}

// Group represents a Slack group.
type Group struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	IsGroup            bool      `json:"is_group"`
	Created            Timestamp `json:"created"`
	Creator            string    `json:"creator"`
	IsArchived         bool      `json:"is_archived"`
	IsMPIM             bool      `json:"is_mpim"`
	Members            []string  `json:"members"`
	Topic              *Topic    `json:"topic"`
	Purpose            *Topic    `json:"purpose"`
	LastRead           Timestamp `json:"last_read"`
	UnreadCount        int       `json:"unread_count"`
	UnreadCountDisplay int       `json:"unread_count_display"`
	Latest             Message   `json:"latest"`
}

// InstantMessage represents a Slack instant message.
type InstantMessage struct {
	ID            string    `json:"id"`
	IsIM          bool      `json:"is_im"`
	User          string    `json:"user"`
	Created       Timestamp `json:"created"`
	IsUserDeleted bool      `json:"is_user_deleted"`
	Latest        Message   `json:"latest"`
}

// Icon represents a Slack icon.
type Icon struct {
	Image24  string `json:"image_24"`
	Image32  string `json:"image_32"`
	Image48  string `json:"image_48"`
	Image72  string `json:"image_72"`
	Image192 string `json:"image_192"`
}

// Bot represents a Slack bot.
type Bot struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icons Icon   `json:"icons"`
}

// BareMessage is an intermediate type used to figure out what final type to deserialize a message as.
type BareMessage struct {
	Type Event `json:"type"`
}

// Message is a basic final message type that encapsulates the most commonly used fields.
type Message struct {
	ID        string    `json:"id"`
	Type      Event     `json:"type"`
	SubType   string    `json:"subtype,omitempty"`
	Hidden    bool      `json:"hidden,omitempty"`
	Timestamp Timestamp `json:"ts,omitempty"`
	Channel   string    `json:"channel,omitempty"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
}

// ChannelJoinedMessage is a final message type for the EventChannelJoined event type.
type ChannelJoinedMessage struct {
	Type    Event   `json:"type"`
	Channel Channel `json:"channel,omitempty"`
}

// Self represents information about the bot itself.
type Self struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Created        Timestamp `json:"created"`
	ManualPresense string    `json:"manual_presence"`
}

// Team represents information about a Slack team.
type Team struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	EmailDomain       string `json:"email_domain"`
	MsgEditWindowMins int    `json:"msg_edit_window_mins"`
	OverStorageLimit  bool   `json:"over_storage_limit"`
}

// Session represents information about a Slack session and is returned by APIStart.
type Session struct {
	OK       bool             `json:"ok"`
	URL      string           `json:"url"`
	Self     *Self            `json:"self"`
	Team     *Team            `json:"team"`
	Users    []User           `json:"users"`
	Channels []Channel        `json:"channels"`
	Groups   []Group          `json:"groups"`
	IMs      []InstantMessage `json:"ims"`
	Error    string           `json:"error,omitempty"`
}

// basicResponse is a utility intermediate type.
type basicResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// ChatMessage is a struct that represents an outgoing chat message for the Slack chat message api.
type ChatMessage struct {
	Token       string                  `json:"token"`
	Channel     string                  `json:"channel"`
	Text        string                  `json:"text"`
	Username    *string                 `json:"username,omitempty"`
	AsUser      *bool                   `json:"as_user,omitempty"`
	Parse       *string                 `json:"parse,omitempty"`
	LinkNames   *bool                   `json:"link_names,omitempty"`
	UnfurlLinks *bool                   `json:"unfurl_links,omitempty"`
	UnfurlMedia *bool                   `json:"unfurl_media,omitempty"`
	IconURL     *string                 `json:"icon_url,omitempty"`
	IconEmoji   *string                 `json:"icon_emoji,omitempty"`
	Attachments []ChatMessageAttachment `json:"attachments,omitempty"`
}

// ChatMessageAttachment is a struct that represents an attachment to a chat message for the Slack chat message api.
type ChatMessageAttachment struct {
	Fallback      string  `json:"fallback"`
	Color         *string `json:"color"`
	Pretext       *string `json:"pretext,omitempty"`
	AuthorName    *string `json:"author_name,omitempty"`
	AuthorLink    *string `json:"author_link,omitempty"`
	AuthorIcon    *string `json:"author_icon,omitempty"`
	Title         *string `json:"title,omitempty"`
	TitleLink     *string `json:"title_link,omitempty"`
	Text          *string `json:"text,omitempty"`
	Fields        []Field `json:"fields,omitempty"`
	ImageURL      *string `json:"image_url,omitempty"`
	ImageThumbURL *string `json:"thumb_url,omitempty"`
}

// Field represents a field on a Slack ChatMessageAttachment.
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Listener is a function that recieves messages from a client.
type Listener func(message *Message, client *Client)

// --------------------------------------------------------------------------------
// Connect(token string) is the main slack entry point.
// It returns a `*Client`, which is the root struct for interacting with slack.
// Register "Listeners" to hook into incoming messages and filter out ones you
// dont' care about.
// --------------------------------------------------------------------------------

// Connect creates a Client with a given token.
func Connect(token string) *Client {
	c := &Client{Token: token, EventListeners: map[Event][]Listener{}, ActiveChannels: []string{}, activeLock: &sync.Mutex{}}
	c.Listen(EventChannelJoined, c.handleChannelJoined)
	c.Listen(EventChannelDeleted, c.handleChannelDeleted)
	c.Listen(EventChannelUnArchive, c.handleChannelUnarchive)
	c.Listen(EventChannelLeft, c.handleChannelLeft)
	return c
}

// Client is the mechanism with which the package consumer interacts with Slack.
type Client struct {
	Token          string
	EventListeners map[Event][]Listener
	ActiveChannels []string

	activeLock       *sync.Mutex
	socketConnection *websocket.Conn
}

// Listen attaches a new Listener to the given event.
// There can be multiple listeners to an event.
// If an event is already being listened for, calling Listen will add a new listener to that event.
func (rtm *Client) Listen(event Event, handler Listener) {
	if listeners, handlesEvent := rtm.EventListeners[event]; handlesEvent {
		rtm.EventListeners[event] = append(listeners, handler)
	} else {
		rtm.EventListeners[event] = []Listener{handler}
	}
}

// StopListening removes all listeners for an event.
func (rtm *Client) StopListening(event Event) {
	delete(rtm.EventListeners, event)
}

// Start begins a session with Slack.
func (rtm *Client) Start() (*Session, error) {
	res := Session{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/rtm.start").
		WithPostData("token", rtm.Token).
		WithPostData("no_unreads", "true").
		WithPostData("mpim_aware", "true").
		FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	//start socket connection
	u, urlErr := url.Parse(res.URL)
	if urlErr != nil {
		return nil, urlErr
	}

	var socketErr error
	rtm.socketConnection, _, socketErr = websocket.DefaultDialer.Dial(u.String(), nil)

	if socketErr != nil {
		return nil, socketErr
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
	return nil
}

// SendMessage sends a basic message over the open web socket connection to slack.
func (rtm Client) SendMessage(m *Message) error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	return rtm.socketConnection.WriteJSON(m)
}

// Say sends a basic message to a given channelID.
func (rtm Client) Say(channelID string, messageComponents ...interface{}) error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	m := &Message{Type: "message", Text: fmt.Sprint(messageComponents...), Channel: channelID}
	return rtm.SendMessage(m)
}

// Sayf is an overload that uses Printf style replacements for a basic message to a given channelID.
func (rtm Client) Sayf(channelID, format string, messageComponents ...interface{}) error {
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

	m := &Message{ID: UUIDv4().ToShortString(), Type: "ping"}
	return rtm.SendMessage(m)
}

//--------------------------------------------------------------------------------
// INTERNAL METHODS
//--------------------------------------------------------------------------------

func (rtm *Client) listenLoop() error {
	for {
		if rtm.socketConnection == nil {
			return nil
		}
		_, messageBytes, err := rtm.socketConnection.ReadMessage()
		if err != nil {
			return err
		}

		var bm BareMessage
		jsonErr := DeserializeJSON(&bm, messageBytes)
		if bm.Type == EventChannelJoined {
			var cm ChannelJoinedMessage
			jsonErr = DeserializeJSON(&cm, messageBytes)
			if jsonErr == nil {
				rtm.dispatch(&Message{Type: EventChannelJoined, Channel: cm.Channel.ID})
			}
		} else {
			var m Message
			jsonErr = DeserializeJSON(&m, messageBytes)
			if jsonErr == nil {
				rtm.dispatch(&m)
			}
		}
	}
}

func (rtm *Client) dispatch(m *Message) {
	if listeners, hasListeners := rtm.EventListeners[m.Type]; hasListeners {
		for x := 0; x < len(listeners); x++ {
			listener := listeners[x]
			go func() {
				listener(m, rtm)
			}()
		}
	}
}

func (rtm *Client) handleChannelJoined(message *Message, client *Client) {
	rtm.activeLock.Lock()
	defer rtm.activeLock.Unlock()
	rtm.ActiveChannels = append(rtm.ActiveChannels, message.Channel)
}

func (rtm *Client) handleChannelUnarchive(message *Message, client *Client) {
	rtm.activeLock.Lock()
	defer rtm.activeLock.Unlock()

	channel, channelErr := rtm.ChannelsInfo(message.Channel)
	if channelErr != nil {
		return
	}
	if channel.IsMember {
		rtm.ActiveChannels = append(rtm.ActiveChannels, message.Channel)
	}
}

func (rtm *Client) handleChannelLeft(message *Message, client *Client) {
	rtm.removeActiveChannel(message.Channel)
}

func (rtm *Client) handleChannelDeleted(message *Message, client *Client) {
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

//--------------------------------------------------------------------------------
// API METHODS
//--------------------------------------------------------------------------------

type channelsListResponse struct {
	Ok       bool      `json:"ok"`
	Error    string    `json:"error"`
	Channels []Channel `json:"channels"`
}

// ChannelsList returns the list of channels available to the bot.
func (rtm *Client) ChannelsList(excludeArchived bool) ([]Channel, error) {
	res := channelsListResponse{}
	req := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.list").
		WithPostData("token", rtm.Token)

	if excludeArchived {
		req = req.WithPostData("exclude_archived", "1")
	}

	resErr := req.FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.Channels, nil
}

type channelsInfoResponse struct {
	Ok      bool     `json:"ok"`
	Error   string   `json:"error"`
	Channel *Channel `json:"channel"`
}

// ChannelsInfo returns information about a given channelID.
func (rtm *Client) ChannelsInfo(channelID string) (*Channel, error) {
	res := channelsInfoResponse{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.info").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.Channel, nil
}

// ChannelsSetTopic sets the topic for a given Slack channel.
func (rtm *Client) ChannelsSetTopic(channelID, topic string) error {
	res := basicResponse{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.leave").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("topic", topic).
		FetchJsonToObject(&res)

	if resErr != nil {
		return resErr
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}

	return nil
}

// ChannelsSetPurpose sets the purpose for a given Slack channel.
func (rtm *Client) ChannelsSetPurpose(channelID, purpose string) error {
	res := basicResponse{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.leave").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("purpose", purpose).
		FetchJsonToObject(&res)

	if resErr != nil {
		return resErr
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}

	return nil
}

type usersListResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Users []User `json:"members"`
}

// UsersList returns all users for a given Slack organization.
func (rtm *Client) UsersList() ([]User, error) {
	res := usersListResponse{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/users.list").
		WithPostData("token", rtm.Token).
		FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.Users, nil
}

type usersInfoResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	User  *User  `json:"users"`
}

// UsersInfo returns an User object for a given userID.
func (rtm *Client) UsersInfo(userID string) (*User, error) {
	res := usersInfoResponse{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/users.info").
		WithPostData("token", rtm.Token).
		WithPostData("user", userID).
		FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.User, nil
}

// NewChatMessage instantiates a ChatMessage for use with ChatPostMessage.
func NewChatMessage(channelID, text string) *ChatMessage {
	return &ChatMessage{Channel: channelID, Text: text, Parse: OptionalString("full")}
}

type chatPostMessageResponse struct {
	Ok        bool         `json:"ok"`
	Timestamp Timestamp    `json:"timestamp"`
	Message   *ChatMessage `json:"message"`
	Error     string       `json:"error"`
}

// ChatPostMessage posts a message to Slack using the chat api.
func (rtm *Client) ChatPostMessage(m *ChatMessage) (*ChatMessage, error) { //the response version of the message is returned for verification
	m.Token = rtm.Token

	res := chatPostMessageResponse{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/chat.postMessage").
		WithPostDataFromObject(m).
		FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.Message, nil
}

// A Timestamp is a special time.Time alias that parses Slack timestamps better.
type Timestamp time.Time

// ParseTimestamp parses a given Slack timestamp.
func ParseTimestamp(strValue string) *Timestamp {
	if integerValue, integerErr := strconv.ParseInt(strValue, 10, 64); integerErr == nil {
		t := Timestamp(time.Unix(integerValue, 0))
		return &t
	}
	if _, floatErr := strconv.ParseFloat(strValue, 64); floatErr == nil {
		components := strings.Split(strValue, ".")
		if integerValue, integerErr := strconv.ParseInt(components[0], 10, 64); integerErr == nil {
			t := Timestamp(time.Unix(integerValue, 0))
			return &t
		}
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshal for the Timestamp struct.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	strValue := string(data)
	t = ParseTimestamp(strValue)
	return nil
}

// DateTime returns a regular golang time.Time for the Timestamp instance.
func (t Timestamp) DateTime() time.Time {
	return time.Time(t)
}

func OptionalUInt8(value uint8) *uint8 {
	return &value
}

func OptionalUInt16(value uint16) *uint16 {
	return &value
}

func OptionalUInt(value uint) *uint {
	return &value
}

func OptionalUInt32(value uint32) *uint32 {
	return &value
}

func OptionalUInt64(value uint64) *uint64 {
	return &value
}

func OptionalInt16(value int16) *int16 {
	return &value
}

func OptionalInt(value int) *int {
	return &value
}

func OptionalInt32(value int32) *int32 {
	return &value
}

func OptionalInt64(value int64) *int64 {
	return &value
}

func OptionalFloat32(value float32) *float32 {
	return &value
}

func OptionalFloat64(value float64) *float64 {
	return &value
}

func OptionalString(value string) *string {
	return &value
}

func OptionalTime(value time.Time) *time.Time {
	return &value
}

func SerializeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func DeserializeJSON(v interface{}, d []byte) error {
	return json.Unmarshal(d, v)
}

func IsEmpty(s string) bool {
	return len(s) == 0
}

type UUID []byte

func (uuid UUID) ToFullString() string {
	b := []byte(uuid)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (uuid UUID) ToShortString() string {
	b := []byte(uuid)
	return fmt.Sprintf("%x", b[:])
}

func (uuid UUID) Version() byte {
	return uuid[6] >> 4
}

func UUIDv4() UUID {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}
