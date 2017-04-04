package slack

import (
	"time"

	"github.com/blendlabs/go-exception"
)

//--------------------------------------------------------------------------------
// API METHODS
//--------------------------------------------------------------------------------

// AuthTest tests if the token works for a client.
func (rtm *Client) AuthTest() (*AuthTestResponse, error) {
	res := AuthTestResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/auth.test").
		WithPostData("token", rtm.Token).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if len(res.Error) != 0 {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}

	return &res, nil
}

// ChannelsHistory returns the messages in a channel.
func (rtm *Client) ChannelsHistory(channelID string, latest, oldest *time.Time, count int, unreads bool) (*ChannelsHistoryResponse, error) {
	unreadsValue := "0"
	if unreads {
		unreadsValue = "1"
	}

	if count == -1 {
		count = 1000
	} else if count < 1 {
		count = 1
	} else if count > 1000 {
		count = 1000
	}

	res := ChannelsHistoryResponse{}
	req := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.history").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("count", string(count)).
		WithPostData("unreads", unreadsValue)

	if latest != nil {
		req = req.WithPostData("latest", Timestamp{time: *latest}.String())
	}

	if oldest != nil {
		req = req.WithPostData("oldest", Timestamp{time: *latest}.String())
	}

	err := req.JSON(&res)
	if err != nil {
		return nil, err
	}

	if len(res.Error) != 0 {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}

	return &res, nil
}

// ChannelsInfo returns information about a given channelID.
func (rtm *Client) ChannelsInfo(channelID string) (*Channel, error) {
	res := channelsInfoResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.info").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}

	return res.Channel, nil
}

// ChannelsList returns the list of channels available to the bot.
func (rtm *Client) ChannelsList(excludeArchived bool) ([]Channel, error) {
	res := channelsListResponse{}
	req := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.list").
		WithPostData("token", rtm.Token)

	if excludeArchived {
		req = req.WithPostData("exclude_archived", "1")
	}

	err := req.JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}

	return res.Channels, nil
}

// ChannelsMark marks a message.
func (rtm *Client) ChannelsMark(channelID string, ts Timestamp) error {
	res := basicResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/chat.mark").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("ts", ts.String()).
		JSON(&res)

	if err != nil {
		return err
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}
	if !res.OK {
		return exception.New("slack response `ok` is false.")
	}

	return nil
}

// ChannelsSetPurpose sets the purpose for a given Slack channel.
func (rtm *Client) ChannelsSetPurpose(channelID, purpose string) error {
	res := basicResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.setPurpose").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("purpose", purpose).
		JSON(&res)

	if err != nil {
		return err
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}

	if !res.OK {
		return exception.New("slack response `ok` is false.")
	}

	return nil
}

// ChannelsSetTopic sets the topic for a given Slack channel.
func (rtm *Client) ChannelsSetTopic(channelID, topic string) error {
	res := basicResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.setTopic").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("topic", topic).
		JSON(&res)

	if err != nil {
		return err
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}

	if !res.OK {
		return exception.New("slack response `ok` is false.")
	}

	return nil
}

// ChatDelete deletes a message.
func (rtm *Client) ChatDelete(channelID string, ts Timestamp) error {
	res := basicResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/chat.delete").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("ts", ts.String()).
		JSON(&res)

	if err != nil {
		return err
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}
	if !res.OK {
		return exception.New("slack response `ok` is false.")
	}

	return nil
}

// ChatPostMessage posts a message to Slack using the chat api.
func (rtm *Client) ChatPostMessage(m *ChatMessage) (*ChatMessageResponse, error) { //the response version of the message is returned for verification
	res := ChatMessageResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/chat.postMessage").
		WithPostData("token", rtm.Token).
		WithPostDataFromObject(m).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}

	return &res, nil
}

// ChatUpdate updates a chat message.
func (rtm *Client) ChatUpdate(ts Timestamp, m *ChatMessage) (*ChatMessageResponse, error) { //the response version of the message is returned for verification
	res := ChatMessageResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/chat.update").
		WithPostData("token", rtm.Token).
		WithPostData("ts", ts.String()).
		WithPostDataFromObject(m).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}

	return &res, nil
}

// EmojiList returns a list of current emoji's for a slack.
func (rtm *Client) EmojiList() (map[string]string, error) {
	res := emojiResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/emoji.list").
		WithPostData("token", rtm.Token).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}
	return res.Emoji, nil
}

// ReactionsAdd adds a reaction.
func (rtm *Client) ReactionsAdd(name string, fileID, fileCommentID, channelID *string, ts *Timestamp) error {
	res := basicResponse{}
	req := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/reactions.add").
		WithPostData("token", rtm.Token).
		WithPostData("name", name)

	if fileID != nil {
		req = req.WithPostData("file", *fileID)
	} else if fileCommentID != nil {
		req = req.WithPostData("file_comment", *fileCommentID)
	} else if channelID != nil && ts != nil {
		req = req.WithPostData("channel", *channelID)
		req = req.WithPostData("timestamp", ts.String())
	} else {
		return exception.New("`fileId` or `fileCommentID` or (`channelID` and `ts`) must be not be nil.")
	}

	err := req.JSON(&res)

	if err != nil {
		return err
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}

	if !res.OK {
		return exception.New("slack response `ok` is false.")
	}
	return nil
}

// ReactionsGet gets reactions.
func (rtm *Client) ReactionsGet(fileID, fileCommentID, channelID *string, ts *Timestamp) (*ChatMessageResponse, error) {
	res := ChatMessageResponse{}
	req := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/reactions.get").
		WithPostData("token", rtm.Token)

	if fileID != nil {
		req = req.WithPostData("file", *fileID)
	} else if fileCommentID != nil {
		req = req.WithPostData("file_comment", *fileCommentID)
	} else if channelID != nil && ts != nil {
		req = req.WithPostData("channel", *channelID)
		req = req.WithPostData("timestamp", ts.String())
	} else {
		return nil, exception.New("`fileId` or `fileCommentID` or (`channelID` and `ts`) must be not be nil.")
	}

	err := req.JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	if !res.OK {
		return nil, exception.New("slack response `ok` is false.")
	}
	return &res, nil
}

// ReactionsRemove removes a reaction.
func (rtm *Client) ReactionsRemove(name string, fileID, fileCommentID, channelID *string, ts *Timestamp) error {
	res := basicResponse{}
	req := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/reactions.remove").
		WithPostData("token", rtm.Token).
		WithPostData("name", name)

	if fileID != nil {
		req = req.WithPostData("file", *fileID)
	} else if fileCommentID != nil {
		req = req.WithPostData("file_comment", *fileCommentID)
	} else if channelID != nil && ts != nil {
		req = req.WithPostData("channel", *channelID)
		req = req.WithPostData("timestamp", ts.String())
	} else {
		return exception.New("`fileId` or `fileCommentID` or (`channelID` and `ts`) must be not be nil.")
	}

	err := req.JSON(&res)

	if err != nil {
		return err
	}

	if !IsEmpty(res.Error) {
		return exception.New(res.Error)
	}

	if !res.OK {
		return exception.New("slack response `ok` is false.")
	}
	return nil
}

// UsersList returns all users for a given Slack organization.
func (rtm *Client) UsersList() ([]User, error) {
	res := usersListResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/users.list").
		WithPostData("token", rtm.Token).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.Users, nil
}

// UsersInfo returns an User object for a given userID.
func (rtm *Client) UsersInfo(userID string) (*User, error) {
	res := usersInfoResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/users.info").
		WithPostData("token", rtm.Token).
		WithPostData("user", userID).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.User, nil
}

// InviteUser invites a user to a channel.
func (rtm *Client) InviteUser(channelID, userID string) (*Channel, error) {
	res := channelsInfoResponse{}
	err := NewExternalRequest().
		AsPost().
		WithScheme(APIScheme).
		WithHost(APIEndpoint).
		WithPath("api/channels.invite").
		WithPostData("token", rtm.Token).
		WithPostData("channel", channelID).
		WithPostData("user", userID).
		JSON(&res)

	if err != nil {
		return nil, err
	}

	if !IsEmpty(res.Error) {
		return nil, exception.New(res.Error)
	}

	return res.Channel, nil
}
