package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericGuildChannel is called upon receiving GuildChannelCreate , GuildChannelUpdate or GuildChannelDelete
type GenericGuildChannel struct {
	*GenericEvent
	ChannelID snowflake.ID
	Channel   fluxer.GuildChannel
	GuildID   snowflake.ID
}

// Guild returns the fluxer.Guild the event happened in.
// This will only check cached guilds!
func (e *GenericGuildChannel) Guild() (fluxer.Guild, bool) {
	return e.Client().Caches.Guild(e.GuildID)
}

// GuildChannelCreate indicates that a new Channel got created in a fluxer.Guild
type GuildChannelCreate struct {
	*GenericGuildChannel
}

// GuildChannelUpdate indicates that a Channel got updated in a fluxer.Guild
type GuildChannelUpdate struct {
	*GenericGuildChannel
	OldChannel fluxer.GuildChannel
}

// GuildChannelDelete indicates that a Channel got deleted in a fluxer.Guild
type GuildChannelDelete struct {
	*GenericGuildChannel
}

// GuildChannelPinsUpdate indicates a fluxer.Message got pinned or unpinned in a fluxer.GuildMessageChannel
type GuildChannelPinsUpdate struct {
	*GenericEvent
	GuildID             snowflake.ID
	ChannelID           snowflake.ID
	NewLastPinTimestamp *time.Time
	OldLastPinTimestamp *time.Time
}
