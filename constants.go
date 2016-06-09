package slack

// Event is a type alias for string to differentiate Slack event types.
type Event string

// Slack constants
const (
	// APIScheme is the protocol used to communicate with slack.
	APIScheme = "https"
	// APIEndpoint is the host used to communicate with slack.
	APIEndpoint = "slack.com"

	// ErrorNotAuthed : No authentication token provided.
	ErrorNotAuthed = "not_authed"
	// ErrorInvalidAuth : Invalid authentication token
	ErrorInvalidAuth = "invalid_auth"
	// ErrorAccountInactive : Authentication token is for a deleted user or team.
	ErrorAccountInactive = "account_inactive"
	// ErrorInvalidArrayArg : The method was passed a PHP-style array argument (e.g. with a name like foo[7]). These are never valid with the Slack API.
	ErrorInvalidArrayArg = "invalid_array_arg"
	// ErrorInvalidCharset : The method was called via a POST request, but the charset specified in the Content-Type header was invalid. Valid charset names are: utf-8 or iso-8859-1.
	ErrorInvalidCharset = "invalid_charset"
	// ErrorInvalidFormData : The method was called via a POST request with Content-Type application/x-www-form-urlencoded or multipart/form-data, but the form data was either missing or syntactically invalid.
	ErrorInvalidFormData = "invalid_form_data"
	// ErrorInvalidPostType : The method was called via a POST request, but the specified Content-Type was invalid. Valid types are: application/json application/x-www-form-urlencoded multipart/form-data text/plain.
	ErrorInvalidPostType = "invalid_post_type"
	// ErrorMissingPostType : The method was called via a POST request and included a data payload, but the request did not include a Content-Type header.
	ErrorMissingPostType = "missing_post_type"
	// ErrorRequestTimeout : The method was called via a POST request, but the POST data was either missing or truncated.
	ErrorRequestTimeout = "request_timeout"
	// ErrorMessageNotFound No message exists with the requested timestamp.
	ErrorMessageNotFound = "message_not_found"
	// ErrorChannelNotFound : Value passed for channel was invalid.
	ErrorChannelNotFound = "channel_not_found"
	// ErrorCantDeleteMessage : Authenticated user does not have permission to delete this message.
	ErrorCantDeleteMessage = "cant_delete_message"
	// ErrorUserIsBot : Authenticated user is a bot and is restricted in it's use of certain api endpoints.
	ErrorUserIsBot = "user_is_bot"
	// ErrorBadTimestamp : Value passed for timestamp was invalid.
	ErrorBadTimestamp = "bad_timestamp"
	// ErrorFileNotFound : File specified by file does not exist.
	ErrorFileNotFound = "file_not_found"
	// ErrorFileCommentNotFound : File comment specified by file_comment does not exist.
	ErrorFileCommentNotFound = "file_comment_not_found"
	// ErrorNoItemSpecified : file, file_comment, or combination of channel and timestamp was not specified.
	ErrorNoItemSpecified = "no_item_specified"
	// ErrorInvalidName : Value passed for name was invalid.
	ErrorInvalidName = "invalid_name"
	// ErrorAlreadyReacted : The specified item already has the user/reaction combination.
	ErrorAlreadyReacted = "already_reacted"
	// ErrorTooManyEmoji : The limit for distinct reactions (i.e emoji) on the item has been reached.
	ErrorTooManyEmoji = "too_many_emoji"
	// ErrorTooManyReactions : 	The limit for reactions a person may add to the item has been reached.
	ErrorTooManyReactions = "too_many_reactions"

	// EventHello is an enumerated event.
	EventHello Event = "hello"
	// EventPing is an enumerated event.
	EventPing Event = "ping"
	// EventPong is an enumerated event.
	EventPong Event = "pong"
	// EventMessage is an enumerated event.
	EventMessage Event = "message"
	// EventMessageACK is a new enumerated event.
	EventMessageACK Event = "message_ack"
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
