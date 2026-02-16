package fluxer

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/internal/flags"
)

// The MessageType indicates the Message type
type MessageType int

// Constants for the MessageType
const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeUserJoin
	MessageTypeGuildBoost
	MessageTypeGuildBoostTier1
	MessageTypeGuildBoostTier2
	MessageTypeGuildBoostTier3
	MessageTypeChannelFollowAdd
	_
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
	MessageTypeGuildDiscoveryGracePeriodInitialWarning
	MessageTypeGuildDiscoveryGracePeriodFinalWarning
	MessageTypeThreadCreated
	MessageTypeReply
	MessageTypeSlashCommand
	MessageTypeThreadStarterMessage
	MessageTypeGuildInviteReminder
	MessageTypeContextMenuCommand
	MessageTypeAutoModerationAction
	MessageTypeRoleSubscriptionPurchase
	MessageTypeInteractionPremiumUpsell
	MessageTypeStageStart
	MessageTypeStageEnd
	MessageTypeStageSpeaker
	_
	MessageTypeStageTopic
	MessageTypeGuildApplicationPremiumSubscription
	_
	_
	_
	MessageTypeGuildIncidentAlertModeEnabled
	MessageTypeGuildIncidentAlertModeDisabled
	MessageTypeGuildIncidentReportRaid
	MessageTypeGuildIncidentReportFalseAlarm
	_
	_
	_
	_
	MessageTypePurchaseNotification
	_
	MessageTypePollResult
)

func (t MessageType) System() bool {
	switch t {
	case MessageTypeDefault, MessageTypeReply, MessageTypeSlashCommand, MessageTypeThreadStarterMessage, MessageTypeContextMenuCommand:
		return false

	default:
		return true
	}
}

func (t MessageType) Deleteable() bool {
	switch t {
	case MessageTypeRecipientAdd, MessageTypeRecipientRemove, MessageTypeCall,
		MessageTypeChannelNameChange, MessageTypeChannelIconChange, MessageTypeThreadStarterMessage:
		return false

	default:
		return true
	}
}

const MessageURLFmt = "https://fluxer.com/channels/%s/%d/%d"

func MessageURL(guildID snowflake.ID, channelID snowflake.ID, messageID snowflake.ID) string {
	return fmt.Sprintf(MessageURLFmt, guildID, channelID, messageID)
}

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	ID                snowflake.ID      `json:"id"`
	GuildID           *snowflake.ID     `json:"guild_id"`
	ChannelID         snowflake.ID      `json:"channel_id"`
	Author            User              `json:"author"`
	WebhookID         *snowflake.ID     `json:"webhook_id,omitempty"`
	Type              MessageType       `json:"type"`
	Flags             MessageFlags      `json:"flags"`
	Content           string            `json:"content,omitempty"`
	CreatedAt         time.Time         `json:"timestamp"`
	EditedAt          *time.Time        `json:"edited_timestamp"`
	Pinned            bool              `json:"pinned"`
	MentionEveryone   bool              `json:"mention_everyone"`
	TTS               bool              `json:"tts"`
	Mentions          []User            `json:"mentions"`
	MentionRoles      []snowflake.ID    `json:"mention_roles"`
	Embeds            []Embed           `json:"embeds,omitempty"`
	Attachments       []Attachment      `json:"attachments"`
	Stickers          []MessageSticker  `json:"sticker,omitempty"`
	Reactions         []MessageReaction `json:"reactions"`
	MessageReference  *MessageReference `json:"message_reference,omitempty"`
	ReferencedMessage *Message          `json:"referenced_message,omitempty"`
	Nonce             Nonce             `json:"nonce,omitempty"`
}

// JumpURL returns the URL which can be used to jump to the message in the discord client.
func (m Message) JumpURL() string {
	guildID := "@me"
	if m.GuildID != nil {
		guildID = m.GuildID.String()
	}
	return fmt.Sprintf(MessageURLFmt, guildID, m.ChannelID, m.ID) // duplicate code, but there isn't a better way without sacrificing user convenience
}

type MentionChannel struct {
	ID      snowflake.ID `json:"id"`
	GuildID snowflake.ID `json:"guild_id"`
	Type    ChannelType  `json:"type"`
	Name    string       `json:"name"`
}

type MessageSticker struct {
	ID          snowflake.ID `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Tags        []string     `json:"tags"`
	Animated    bool         `json:"animated"`
}

// MessageReaction contains information about the reactions of a message
type MessageReaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details"`
	Me           bool                 `json:"me"`
	MeBurst      bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []string             `json:"burst_colors"`
}

type ReactionCountDetails struct {
	Burst  int `json:"burst"`
	Normal int `json:"normal"`
}

type MessageReactionType int

const (
	MessageReactionTypeNormal MessageReactionType = iota
	MessageReactionTypeBurst
)

// MessageActivityType is the type of MessageActivity https://fluxer.com/developers/docs/resources/message#message-object-message-activity-types
type MessageActivityType int

// Constants for MessageActivityType
const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	_
	MessageActivityTypeJoinRequest
)

// MessageActivity is used for rich presence-related chat embeds in a Message
type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID *string             `json:"party_id,omitempty"`
}

// MessageApplication is used for rich presence-related chat embeds in a Message
type MessageApplication struct {
	ID          snowflake.ID `json:"id"`
	CoverImage  *string      `json:"cover_image,omitempty"`
	Description string       `json:"description"`
	Icon        *string      `json:"icon,omitempty"`
	Name        string       `json:"name"`
}

// MessageReference is a reference to another message
type MessageReference struct {
	Type            MessageReferenceType `json:"type,omitempty"`
	MessageID       *snowflake.ID        `json:"message_id"`
	ChannelID       *snowflake.ID        `json:"channel_id,omitempty"`
	GuildID         *snowflake.ID        `json:"guild_id,omitempty"`
	FailIfNotExists bool                 `json:"fail_if_not_exists,omitempty"`
}

type MessageReferenceType int

const (
	MessageReferenceTypeDefault MessageReferenceType = iota
	MessageReferenceTypeForward
)

type MessageSnapshot struct {
	Message PartialMessage `json:"message"`
}

type PartialMessage struct {
	Type            MessageType      `json:"type"`
	Content         string           `json:"content,omitempty"`
	Embeds          []Embed          `json:"embeds,omitempty"`
	Attachments     []Attachment     `json:"attachments"`
	CreatedAt       time.Time        `json:"timestamp"`
	EditedTimestamp *time.Time       `json:"edited_timestamp"`
	Flags           MessageFlags     `json:"flags"`
	Mentions        []User           `json:"mentions"`
	MentionRoles    []snowflake.ID   `json:"mention_roles"`
	Stickers        []Sticker        `json:"stickers"`
	StickerItems    []MessageSticker `json:"sticker_items,omitempty"`
}

type MessageBulkDelete struct {
	Messages []snowflake.ID `json:"messages"`
}

// The MessageFlags of a Message
type MessageFlags int

// Constants for MessageFlags
const (
	MessageFlagCrossposted MessageFlags = 1 << iota
	MessageFlagIsCrosspost
	MessageFlagSuppressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	MessageFlagHasThread
	MessageFlagEphemeral
	MessageFlagLoading // Message is an interaction of type 5, awaiting further response
	MessageFlagFailedToMentionSomeRolesInThread
	_
	_
	_
	MessageFlagSuppressNotifications
	MessageFlagIsVoiceMessage
	MessageFlagHasSnapshot
	// MessageFlagIsComponentsV2 should be set when you want to send v2 components.
	// After setting this, you will not be allowed to send message content and embeds anymore.
	// Once a message with the flag has been sent, it cannot be removed by editing the message.
	MessageFlagIsComponentsV2
	MessageFlagsNone MessageFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f MessageFlags) Add(bits ...MessageFlags) MessageFlags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f MessageFlags) Remove(bits ...MessageFlags) MessageFlags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f MessageFlags) Has(bits ...MessageFlags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f MessageFlags) Missing(bits ...MessageFlags) bool {
	return flags.Missing(f, bits...)
}

type RoleSubscriptionData struct {
	RoleSubscriptionListingID snowflake.ID `json:"role_subscription_listing_id"`
	TierName                  string       `json:"tier_name"`
	TotalMonthsSubscribed     int          `json:"total_months_subscribed"`
	IsRenewal                 bool         `json:"is_renewal"`
}

type MessageCall struct {
	Participants   []snowflake.ID `json:"participants"`
	EndedTimestamp *time.Time     `json:"ended_timestamp"`
}

// Nonce is a string or int used when sending a message to fluxer.
type Nonce string

// UnmarshalJSON unmarshals the Nonce from a string or int.
func (n *Nonce) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		i, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		*n = Nonce(strconv.FormatInt(i, 10))
	} else {
		*n = Nonce(unquoted)
	}

	return nil
}

type ChannelPins struct {
	Items   []MessagePin `json:"items"`
	HasMore bool         `json:"has_more"`
}

type MessagePin struct {
	PinnedAt time.Time `json:"pinned_at"`
	Message  Message   `json:"message"`
}
