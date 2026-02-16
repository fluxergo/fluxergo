package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericGuildMessage is called upon receiving GuildMessageCreate , GuildMessageUpdate or GuildMessageDelete
type GenericGuildMessage struct {
	*GenericEvent
	MessageID snowflake.ID
	Message   fluxer.Message
	ChannelID snowflake.ID
	GuildID   snowflake.ID
}

// Guild returns the fluxer.Guild the GenericGuildMessage happened in.
// This will only check cached guilds!
func (e *GenericGuildMessage) Guild() (fluxer.Guild, bool) {
	return e.Client().Caches.Guild(e.GuildID)
}

// Channel returns the fluxer.GuildMessageChannel where the GenericGuildMessage happened
func (e *GenericGuildMessage) Channel() (fluxer.GuildMessageChannel, bool) {
	return e.Client().Caches.GuildMessageChannel(e.ChannelID)
}

// GuildMessageCreate is called upon receiving a fluxer.Message in a Channel
type GuildMessageCreate struct {
	*GenericGuildMessage
}

// GuildMessageUpdate is called upon editing a fluxer.Message in a Channel
type GuildMessageUpdate struct {
	*GenericGuildMessage
	OldMessage fluxer.Message
}

// GuildMessageDelete is called upon deleting a fluxer.Message in a Channel
type GuildMessageDelete struct {
	*GenericGuildMessage
}
