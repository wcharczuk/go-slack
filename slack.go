package slack

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/blendlabs/go-exception"
	"github.com/blendlabs/go-request"
	"github.com/blendlabs/go-util"
	"github.com/gorilla/websocket"
)

type Timestamp time.Time

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

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	strValue := string(data)
	t = ParseTimestamp(strValue)
	return nil
}

func (t Timestamp) DateTime() time.Time {
	return time.Time(t)
}

type Event string

const (
	API_SCHEME   = "https"
	API_ENDPOINT = "slack.com"

	EVENT_HELLO                   Event = "hello"
	EVENT_MESSAGE                 Event = "message"
	EVENT_USER_TYPING             Event = "user_typing"
	EVENT_CHANNEL_MARKED          Event = "channel_marked"
	EVENT_CHANNEL_JOINED          Event = "channel_joined"
	EVENT_CHANNEL_LEFT            Event = "channel_left"
	EVENT_CHANNEL_DELETED         Event = "channel_deleted"
	EVENT_CHANNEL_RENAME          Event = "channel_rename"
	EVENT_CHANNEL_ARCHIVE         Event = "channel_archive"
	EVENT_CHANNEL_UNARCHIVE       Event = "channel_unarchive"
	EVENT_CHANNEL_HISTORY_CHANGED Event = "channel_history_changed"
	EVENT_DND_UPDATED             Event = "dnd_updated"
	EVENT_DND_UPDATED_USER        Event = "dnd_updated_user"
	EVENT_IM_CREATED              Event = "im_created"
	EVENT_IM_OPEN                 Event = "im_open"
	EVENT_IM_CLOSE                Event = "im_close"
	EVENT_IM_MARKED               Event = "im_marked"
	EVENT_IM_HISTORY_CHANGED      Event = "im_history_changed"
	EVENT_GROUP_JOINED            Event = "group_joined"
	EVENT_GROUP_LEFT              Event = "group_left"
	EVENT_GROUP_OPEN              Event = "group_open"
	EVENT_GROUP_CLOSE             Event = "group_close"
	EVENT_GROUP_ARCHIVE           Event = "group_archive"
	EVENT_GROUP_UNARCHIVE         Event = "group_unarchive"
	EVENT_GROUP_RENAME            Event = "group_rename"
	EVENT_GROUP_MARKED            Event = "group_marked"
	EVENT_GROUP_HISTORY_CHANGED   Event = "group_history_changed"
	EVENT_FILE_CREATED            Event = "file_created"
	EVENT_FILE_SHARED             Event = "file_shared"
	EVENT_FILE_UNSHARED           Event = "file_unshared"
	EVENT_FILE_PUBLIC             Event = "file_public"
	EVENT_FILE_PRIVATE            Event = "file_private"
	EVENT_FILE_CHANGE             Event = "file_change"
	EVENT_FILE_DELETED            Event = "file_deleted"
	EVENT_FILE_COMMENT_ADDED      Event = "file_comment_added"
	EVENT_FILE_COMMENT_EDITED     Event = "file_comment_edited"
	EVENT_FILE_COMMENT_DELETED    Event = "file_comment_deleted"
	EVENT_PIN_ADDED               Event = "pin_added"
	EVENT_PIN_REMOVED             Event = "pin_removed"
	EVENT_PRESENCE_CHANGE         Event = "presence_change"
	EVENT_MANUAL_PRESENCE_CHANGE  Event = "manual_presence_change"
	EVENT_PREF_CHANGE             Event = "pref_change"
	EVENT_USER_CHANGE             Event = "user_change"
	EVENT_TEAM_JOIN               Event = "team_join"
	EVENT_STAR_ADDED              Event = "star_added"
	EVENT_STAR_REMOVED            Event = "star_removed"
	EVENT_REACTION_ADDED          Event = "reaction_added"
	EVENT_REACTION_REMOVED        Event = "reaction_removed"
	EVENT_EMOJI_CHANGED           Event = "emoji_changed"
	EVENT_COMMANDS_CHANGED        Event = "commands_changed"
	EVENT_TEAM_PLAN_CHANGED       Event = "team_plan_changed"
	EVENT_TEAM_PREF_CHANGED       Event = "team_pref_changed"
	EVENT_EMAIL_DOMAIN_CHANGED    Event = "email_domain_changed"
	EVENT_TEAM_PROFILE_CHANGE     Event = "team_profile_change"
	EVENT_TEAM_PROFILE_DELETE     Event = "team_profile_delete"
	EVENT_TEAM_PROFILE_REORDER    Event = "team_profile_reorder"
	EVENT_BOT_ADDED               Event = "bot_added"
	EVENT_BOT_CHANGED             Event = "bot_changed"
	EVENT_ACCOUNTS_CHANGED        Event = "accounts_changed"
	EVENT_TEAM_MIGRATION_STARTED  Event = "team_migration_started"
)

type User struct {
	Id                string       `json:"id"`
	Name              string       `json:"name"`
	Deleted           bool         `json:"deletd"`
	Color             string       `json:"color"`
	Profile           *UserProfile `json:"profile"`
	IsAdmin           bool         `json:"is_admin"`
	IsOwner           bool         `json:"is_owner"`
	IsPrimaryOwner    bool         `json:"is_primary_owner"`
	IsRestricted      bool         `json:"is_restricted"`
	IsUltraRestricted bool         `json:"is_ultra_restricted"`
	Has2FA            bool         `json:"has_2fa"`
	TwoFactorType     string       `json:"two_factor_type"`
	HasFiles          bool         `json:"has_files"`
}

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

type Channel struct {
	Id                 string    `json:"id"`
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

	Latest []Message `json:"latest"`
}

type Topic struct {
	Value   string    `json:"value"`
	Creator string    `json:"creator"`
	LastSet Timestamp `json:"last_set"`
}

type Group struct {
	Id                 string    `json:"id"`
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
	Latest             []Message `json:"latest"`
}

type InstantMessage struct {
	Id            string    `json:"id"`
	IsIM          bool      `json:"is_im"`
	User          string    `json:"user"`
	Created       Timestamp `json:"created"`
	IsUserDeleted bool      `json:"is_user_deleted"`
	Latest        []Message `json:"latest"`
}

type Message struct {
	Type      Event     `json:"type"`
	SubType   string    `json:"subtype,omitempty"`
	Hidden    bool      `json:"hidden,omitemptyH"`
	Timestamp Timestamp `json:"ts"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
}

type Self struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Created        Timestamp `json:"created"`
	ManualPresense string    `json:"manual_presence"`
}

type Team struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	EmailDomain       string `json:"email_domain"`
	MsgEditWindowMins int    `json:"msg_edit_window_mins"`
	OverStorageLimit  bool   `json:"over_storage_limit"`
}

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

type BasicResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type Listener func(message *Message, client *Client)

func Connect(token string) *Client {
	return &Client{Token: token, EventListeners: map[Event][]Listener{}}
}

type Client struct {
	Token          string
	EventListeners map[Event][]Listener

	allListeners     []Listener
	socketConnection *websocket.Conn
}

func (rtm *Client) Listen(event Event, handler Listener) {
	if listeners, handlesEvent := rtm.EventListeners[event]; handlesEvent {
		rtm.EventListeners[event] = append(listeners, handler)
	} else {
		rtm.EventListeners[event] = []Listener{handler}
	}
}

func (rtm *Client) ListenAll(handler Listener) {
	rtm.allListeners = append(rtm.allListeners, handler)
}

func (rtm *Client) StopListening(event Event) {
	delete(rtm.EventListeners, event)
}

func (rtm *Client) Start() (*Session, error) {
	res := Session{}
	resErr := request.NewRequest().
		AsPost().
		WithScheme(API_SCHEME).
		WithHost(API_ENDPOINT).
		WithPath("api/rtm.start").
		WithPostData("token", rtm.Token).
		FetchJsonToObject(&res)

	if resErr != nil {
		return nil, resErr
	}

	if !util.IsEmpty(res.Error) {
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

	for event, listeners := range rtm.EventListeners {
		fmt.Printf("Event: `%s` has %d listeners.\n", event, len(listeners))
	}

	go func() {
		rtm.listenLoop()
	}()

	return &res, nil
}

func (rtm Client) SendMessage(m *Message) error {
	if rtm.socketConnection == nil {
		return exception.New("Connection is closed.")
	}

	return rtm.socketConnection.WriteJSON(m)
}

func (rtm *Client) listenLoop() error {
	for {
		if rtm.socketConnection == nil {
			return nil
		}
		_, messageBytes, err := rtm.socketConnection.ReadMessage()
		if err != nil {
			return err
		}
		var m Message
		jsonErr := util.DeserializeJson(&m, string(messageBytes))
		if jsonErr != nil {
			return jsonErr
		}

		rtm.dispatch(&m)
	}
	return nil
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

	for x := 0; x < len(rtm.allListeners); x++ {
		listener := rtm.allListeners[x]
		go func() {
			listener(m, rtm)
		}()
	}
}

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
