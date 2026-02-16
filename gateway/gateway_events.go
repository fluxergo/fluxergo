package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// EventData is the base interface for all data types sent by discord
type EventData interface {
	MessageData
	eventData()
}

// EventUnknown is an event that is not known to fluxergo
type EventUnknown json.RawMessage

func (e EventUnknown) MarshalJSON() ([]byte, error) {
	return json.RawMessage(e).MarshalJSON()
}

func (e *EventUnknown) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(e).UnmarshalJSON(data)
}

func (EventUnknown) messageData() {}
func (EventUnknown) eventData()   {}

// EventReady is the event sent by discord when you successfully Identify
type EventReady struct {
	Version          int                       `json:"v"`
	User             fluxer.OAuth2User         `json:"user"`
	Guilds           []fluxer.UnavailableGuild `json:"guilds"`
	SessionID        string                    `json:"session_id"`
	ResumeGatewayURL string                    `json:"resume_gateway_url"`
	Shard            [2]int                    `json:"shard,omitempty"`
	Application      fluxer.PartialApplication `json:"application"`
}

func (EventReady) messageData() {}
func (EventReady) eventData()   {}

// EventResumed is the event sent by discord when you successfully resume
type EventResumed struct{}

func (EventResumed) messageData() {}
func (EventResumed) eventData()   {}

// RateLimitedMetadata is an interface for all ratelimited metadatas.
// [RateLimitedMetadataRequestGuildMembers]
// [RateLimitedMetadataUnknown]
type RateLimitedMetadata interface {
	// ratelimitedmetadata is a marker to simulate unions.
	ratelimitedmetadata()
}

type RateLimitedMetadataRequestGuildMembers struct {
	GuildID snowflake.ID `json:"guild_id"`
	Nonce   string       `json:"nonce"`
}

func (RateLimitedMetadataRequestGuildMembers) ratelimitedmetadata() {}

type RateLimitedMetadataUnknown json.RawMessage

func (RateLimitedMetadataUnknown) ratelimitedmetadata() {}

type EventRateLimited struct {
	Opcode     Opcode              `json:"opcode"`
	RetryAfter float64             `json:"retry_after"`
	Meta       RateLimitedMetadata `json:"meta"`
}

func (e *EventRateLimited) UnmarshalJSON(data []byte) error {
	var event struct {
		Opcode     Opcode          `json:"opcode"`
		RetryAfter float64         `json:"retry_after"`
		Meta       json.RawMessage `json:"meta"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	var (
		meta RateLimitedMetadata
		err  error
	)

	switch event.Opcode {
	case OpcodeRequestGuildMembers:
		var v RateLimitedMetadataRequestGuildMembers
		err = json.Unmarshal(event.Meta, &v)
		meta = v

	default:
		meta = RateLimitedMetadataUnknown{}
	}

	if err != nil {
		return fmt.Errorf("failed to deserialize metadata payload for opcode %d: %w", event.Opcode, err)
	}

	e.Opcode = event.Opcode
	e.RetryAfter = event.RetryAfter
	e.Meta = meta

	return nil
}

func (EventRateLimited) messageData() {}
func (EventRateLimited) eventData()   {}

type EventChannelCreate struct {
	fluxer.GuildChannel
}

func (e *EventChannelCreate) UnmarshalJSON(data []byte) error {
	var v fluxer.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(fluxer.GuildChannel)
	return nil
}

func (EventChannelCreate) messageData() {}
func (EventChannelCreate) eventData()   {}

type EventChannelUpdate struct {
	fluxer.GuildChannel
}

func (e *EventChannelUpdate) UnmarshalJSON(data []byte) error {
	var v fluxer.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(fluxer.GuildChannel)
	return nil
}

func (EventChannelUpdate) messageData() {}
func (EventChannelUpdate) eventData()   {}

type EventChannelDelete struct {
	fluxer.GuildChannel
}

func (e *EventChannelDelete) UnmarshalJSON(data []byte) error {
	var v fluxer.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(fluxer.GuildChannel)
	return nil
}

func (EventChannelDelete) messageData() {}
func (EventChannelDelete) eventData()   {}

type EventGuildCreate struct {
	fluxer.GatewayGuild
}

func (EventGuildCreate) messageData() {}
func (EventGuildCreate) eventData()   {}

type EventGuildUpdate struct {
	fluxer.GatewayGuild
}

func (EventGuildUpdate) messageData() {}
func (EventGuildUpdate) eventData()   {}

type EventGuildDelete struct {
	fluxer.UnavailableGuild
}

func (EventGuildDelete) messageData() {}
func (EventGuildDelete) eventData()   {}

type EventGuildAuditLogEntryCreate struct {
	fluxer.AuditLogEntry
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildAuditLogEntryCreate) messageData() {}
func (EventGuildAuditLogEntryCreate) eventData()   {}

type EventMessageReactionAdd struct {
	UserID          snowflake.ID               `json:"user_id"`
	ChannelID       snowflake.ID               `json:"channel_id"`
	MessageID       snowflake.ID               `json:"message_id"`
	GuildID         *snowflake.ID              `json:"guild_id"`
	Member          *fluxer.Member             `json:"member"`
	Emoji           fluxer.PartialEmoji        `json:"emoji"`
	MessageAuthorID *snowflake.ID              `json:"message_author_id"`
	BurstColors     []string                   `json:"burst_colors"`
	Burst           bool                       `json:"burst"`
	Type            fluxer.MessageReactionType `json:"type"`
}

func (e *EventMessageReactionAdd) UnmarshalJSON(data []byte) error {
	type eventMessageReactionAdd EventMessageReactionAdd
	var v eventMessageReactionAdd
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventMessageReactionAdd(v)
	if e.Member != nil && e.GuildID != nil {
		e.Member.GuildID = *e.GuildID
	}
	return nil
}

func (EventMessageReactionAdd) messageData() {}
func (EventMessageReactionAdd) eventData()   {}

type EventMessageReactionRemove struct {
	UserID      snowflake.ID               `json:"user_id"`
	ChannelID   snowflake.ID               `json:"channel_id"`
	MessageID   snowflake.ID               `json:"message_id"`
	GuildID     *snowflake.ID              `json:"guild_id"`
	Emoji       fluxer.PartialEmoji        `json:"emoji"`
	BurstColors []string                   `json:"burst_colors"`
	Burst       bool                       `json:"burst"`
	Type        fluxer.MessageReactionType `json:"type"`
}

func (EventMessageReactionRemove) messageData() {}
func (EventMessageReactionRemove) eventData()   {}

type EventMessageReactionRemoveEmoji struct {
	ChannelID snowflake.ID        `json:"channel_id"`
	MessageID snowflake.ID        `json:"message_id"`
	GuildID   *snowflake.ID       `json:"guild_id"`
	Emoji     fluxer.PartialEmoji `json:"emoji"`
}

func (EventMessageReactionRemoveEmoji) messageData() {}
func (EventMessageReactionRemoveEmoji) eventData()   {}

type EventMessageReactionRemoveAll struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	MessageID snowflake.ID  `json:"message_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
}

func (EventMessageReactionRemoveAll) messageData() {}
func (EventMessageReactionRemoveAll) eventData()   {}

type EventChannelPinsUpdate struct {
	GuildID          *snowflake.ID `json:"guild_id"`
	ChannelID        snowflake.ID  `json:"channel_id"`
	LastPinTimestamp *time.Time    `json:"last_pin_timestamp"`
}

func (EventChannelPinsUpdate) messageData() {}
func (EventChannelPinsUpdate) eventData()   {}

type EventGuildMembersChunk struct {
	GuildID    snowflake.ID      `json:"guild_id"`
	Members    []fluxer.Member   `json:"members"`
	ChunkIndex int               `json:"chunk_index"`
	ChunkCount int               `json:"chunk_count"`
	NotFound   []snowflake.ID    `json:"not_found"`
	Presences  []fluxer.Presence `json:"presences"`
	Nonce      string            `json:"nonce"`
}

func (EventGuildMembersChunk) messageData() {}
func (EventGuildMembersChunk) eventData()   {}

type EventGuildBanAdd struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    fluxer.User  `json:"user"`
}

func (EventGuildBanAdd) messageData() {}
func (EventGuildBanAdd) eventData()   {}

type EventGuildBanRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    fluxer.User  `json:"user"`
}

func (EventGuildBanRemove) messageData() {}
func (EventGuildBanRemove) eventData()   {}

type EventGuildEmojisUpdate struct {
	GuildID snowflake.ID   `json:"guild_id"`
	Emojis  []fluxer.Emoji `json:"emojis"`
}

func (e *EventGuildEmojisUpdate) UnmarshalJSON(data []byte) error {
	type eventGuildEmojisUpdate EventGuildEmojisUpdate
	var v eventGuildEmojisUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildEmojisUpdate(v)
	for i := range e.Emojis {
		e.Emojis[i].GuildID = e.GuildID
	}
	return nil
}

func (EventGuildEmojisUpdate) messageData() {}
func (EventGuildEmojisUpdate) eventData()   {}

type EventGuildStickersUpdate struct {
	GuildID  snowflake.ID     `json:"guild_id"`
	Stickers []fluxer.Sticker `json:"stickers"`
}

func (EventGuildStickersUpdate) messageData() {}
func (EventGuildStickersUpdate) eventData()   {}

type EventGuildIntegrationsUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildIntegrationsUpdate) messageData() {}
func (EventGuildIntegrationsUpdate) eventData()   {}

type EventGuildMemberAdd struct {
	fluxer.Member
}

func (EventGuildMemberAdd) messageData() {}
func (EventGuildMemberAdd) eventData()   {}

type EventGuildMemberUpdate struct {
	fluxer.Member
}

func (EventGuildMemberUpdate) messageData() {}
func (EventGuildMemberUpdate) eventData()   {}

type EventGuildMemberRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    fluxer.User  `json:"user"`
}

func (EventGuildMemberRemove) messageData() {}
func (EventGuildMemberRemove) eventData()   {}

type EventGuildRoleCreate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    fluxer.Role  `json:"role"`
}

func (e *EventGuildRoleCreate) UnmarshalJSON(data []byte) error {
	type eventGuildRoleCreate EventGuildRoleCreate
	var v eventGuildRoleCreate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildRoleCreate(v)
	e.Role.GuildID = e.GuildID
	return nil
}

func (e *EventGuildRoleCreate) MarshalJSON() ([]byte, error) {
	type eventGuildRoleCreate EventGuildRoleCreate
	e.GuildID = e.Role.GuildID
	return json.Marshal(eventGuildRoleCreate(*e))
}

func (EventGuildRoleCreate) messageData() {}
func (EventGuildRoleCreate) eventData()   {}

type EventGuildRoleDelete struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
}

func (EventGuildRoleDelete) messageData() {}
func (EventGuildRoleDelete) eventData()   {}

type EventGuildRoleUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    fluxer.Role  `json:"role"`
}

func (e *EventGuildRoleUpdate) UnmarshalJSON(data []byte) error {
	type eventGuildRoleUpdate EventGuildRoleUpdate
	var v eventGuildRoleUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildRoleUpdate(v)
	e.Role.GuildID = e.GuildID
	return nil
}

func (e *EventGuildRoleUpdate) MarshalJSON() ([]byte, error) {
	type eventGuildRoleUpdate EventGuildRoleUpdate
	e.GuildID = e.Role.GuildID
	return json.Marshal(eventGuildRoleUpdate(*e))
}

func (EventGuildRoleUpdate) messageData() {}
func (EventGuildRoleUpdate) eventData()   {}

type EventGuildScheduledEventCreate struct {
	fluxer.GuildScheduledEvent
}

func (EventGuildScheduledEventCreate) messageData() {}
func (EventGuildScheduledEventCreate) eventData()   {}

type EventGuildScheduledEventUpdate struct {
	fluxer.GuildScheduledEvent
}

func (EventGuildScheduledEventUpdate) messageData() {}
func (EventGuildScheduledEventUpdate) eventData()   {}

type EventGuildScheduledEventDelete struct {
	fluxer.GuildScheduledEvent
}

func (EventGuildScheduledEventDelete) messageData() {}
func (EventGuildScheduledEventDelete) eventData()   {}

type EventGuildScheduledEventUserAdd struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

func (EventGuildScheduledEventUserAdd) messageData() {}
func (EventGuildScheduledEventUserAdd) eventData()   {}

type EventGuildScheduledEventUserRemove struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

func (EventGuildScheduledEventUserRemove) messageData() {}
func (EventGuildScheduledEventUserRemove) eventData()   {}

type EventInviteCreate struct {
	ChannelID         snowflake.ID               `json:"channel_id"`
	Code              string                     `json:"code"`
	CreatedAt         time.Time                  `json:"created_at"`
	GuildID           *snowflake.ID              `json:"guild_id"`
	Inviter           *fluxer.User               `json:"inviter"`
	MaxAge            int                        `json:"max_age"`
	MaxUses           int                        `json:"max_uses"`
	TargetType        fluxer.InviteTargetType    `json:"target_type"`
	TargetUser        *fluxer.User               `json:"target_user"`
	TargetApplication *fluxer.PartialApplication `json:"target_application"`
	Temporary         bool                       `json:"temporary"`
	Uses              int                        `json:"uses"`
	ExpiresAt         *time.Time                 `json:"expires_at"`
	RoleIDs           *[]snowflake.ID            `json:"role_ids"`
}

func (EventInviteCreate) messageData() {}
func (EventInviteCreate) eventData()   {}

type EventInviteDelete struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
	Code      string        `json:"code"`
}

func (EventInviteDelete) messageData() {}
func (EventInviteDelete) eventData()   {}

type EventMessageCreate struct {
	fluxer.Message
}

func (EventMessageCreate) messageData() {}
func (EventMessageCreate) eventData()   {}

type EventMessageUpdate struct {
	fluxer.Message
}

func (EventMessageUpdate) messageData() {}
func (EventMessageUpdate) eventData()   {}

type EventMessageDelete struct {
	ID        snowflake.ID  `json:"id"`
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id,omitempty"`
}

func (EventMessageDelete) messageData() {}
func (EventMessageDelete) eventData()   {}

type EventMessageDeleteBulk struct {
	IDs       []snowflake.ID `json:"ids"`
	ChannelID snowflake.ID   `json:"channel_id"`
	GuildID   *snowflake.ID  `json:"guild_id,omitempty"`
}

func (EventMessageDeleteBulk) messageData() {}
func (EventMessageDeleteBulk) eventData()   {}

type EventPresenceUpdate struct {
	fluxer.Presence
}

func (EventPresenceUpdate) messageData() {}
func (EventPresenceUpdate) eventData()   {}

type EventTypingStart struct {
	ChannelID snowflake.ID   `json:"channel_id"`
	GuildID   *snowflake.ID  `json:"guild_id,omitempty"`
	UserID    snowflake.ID   `json:"user_id"`
	Timestamp time.Time      `json:"timestamp"`
	Member    *fluxer.Member `json:"member,omitempty"`
	User      fluxer.User    `json:"user"`
}

func (e *EventTypingStart) UnmarshalJSON(data []byte) error {
	type typingStartEvent EventTypingStart
	var v struct {
		Timestamp int64 `json:"timestamp"`
		typingStartEvent
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventTypingStart(v.typingStartEvent)
	e.Timestamp = time.Unix(v.Timestamp, 0)
	return nil
}

func (EventTypingStart) messageData() {}
func (EventTypingStart) eventData()   {}

type EventUserUpdate struct {
	fluxer.OAuth2User
}

func (EventUserUpdate) messageData() {}
func (EventUserUpdate) eventData()   {}

type EventVoiceStateUpdate struct {
	fluxer.VoiceState
	Member fluxer.Member `json:"member"`
}

func (EventVoiceStateUpdate) messageData() {}
func (EventVoiceStateUpdate) eventData()   {}

type EventVoiceServerUpdate struct {
	ChannelID    snowflake.ID `json:"channel_id"`
	ConnectionID string       `json:"connection_id"`
	Endpoint     *string      `json:"endpoint"`
	GuildID      snowflake.ID `json:"guild_id"`
	Token        string       `json:"token"`
}

func (EventVoiceServerUpdate) messageData() {}
func (EventVoiceServerUpdate) eventData()   {}

type EventWebhooksUpdate struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

func (EventWebhooksUpdate) messageData() {}
func (EventWebhooksUpdate) eventData()   {}

type EventIntegrationCreate struct {
	fluxer.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationCreate) UnmarshalJSON(data []byte) error {
	var integration fluxer.UnmarshalIntegration
	if err := json.Unmarshal(data, &integration); err != nil {
		return err
	}

	var v struct {
		GuildID snowflake.ID `json:"guild_id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Integration = integration.Integration
	e.GuildID = v.GuildID
	return nil
}

func (EventIntegrationCreate) messageData() {}
func (EventIntegrationCreate) eventData()   {}

type EventIntegrationUpdate struct {
	fluxer.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationUpdate) UnmarshalJSON(data []byte) error {
	var integration fluxer.UnmarshalIntegration
	if err := json.Unmarshal(data, &integration); err != nil {
		return err
	}

	var v struct {
		GuildID snowflake.ID `json:"guild_id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Integration = integration.Integration
	e.GuildID = v.GuildID
	return nil
}

func (EventIntegrationUpdate) messageData() {}
func (EventIntegrationUpdate) eventData()   {}

type EventIntegrationDelete struct {
	ID            snowflake.ID  `json:"id"`
	GuildID       snowflake.ID  `json:"guild_id"`
	ApplicationID *snowflake.ID `json:"application_id"`
}

func (EventIntegrationDelete) messageData() {}
func (EventIntegrationDelete) eventData()   {}

type EventRaw struct {
	EventType EventType
	Payload   io.Reader
}

func (EventRaw) messageData() {}
func (EventRaw) eventData()   {}

type EventHeartbeatAck struct {
	LastHeartbeat time.Time
	NewHeartbeat  time.Time
}

func (EventHeartbeatAck) messageData() {}
func (EventHeartbeatAck) eventData()   {}
