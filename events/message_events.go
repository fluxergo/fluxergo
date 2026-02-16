package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericMessage generic fluxer.Message event
type GenericMessage struct {
	*GenericEvent
	MessageID snowflake.ID
	Message   fluxer.Message
	ChannelID snowflake.ID
	GuildID   *snowflake.ID
}

// Channel returns the fluxer.GuildMessageChannel where the GenericMessage happened
func (e *GenericMessage) Channel() (fluxer.GuildMessageChannel, bool) {
	return e.Client().Caches.GuildMessageChannel(e.ChannelID)
}

// Guild returns the fluxer.Guild where the GenericMessage happened or nil if it happened in DMs
func (e *GenericMessage) Guild() (fluxer.Guild, bool) {
	if e.GuildID == nil {
		return fluxer.Guild{}, false
	}
	return e.Client().Caches.Guild(*e.GuildID)
}

// MessageCreate indicates that a fluxer.Message got received
type MessageCreate struct {
	*GenericMessage
}

// MessageUpdate indicates that a fluxer.Message got update
type MessageUpdate struct {
	*GenericMessage
	OldMessage fluxer.Message
}

// MessageDelete indicates that a fluxer.Message got deleted
type MessageDelete struct {
	*GenericMessage
}
