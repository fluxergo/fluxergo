package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericGuildMessageReaction is called upon receiving GuildMessageReactionAdd or GuildMessageReactionRemove
type GenericGuildMessageReaction struct {
	*GenericEvent
	UserID      snowflake.ID
	ChannelID   snowflake.ID
	MessageID   snowflake.ID
	GuildID     snowflake.ID
	Emoji       fluxer.PartialEmoji
	BurstColors []string
	Burst       bool
}

// Member returns the Member that reacted to the fluxer.Message from the cache.
func (e *GenericGuildMessageReaction) Member() (fluxer.Member, bool) {
	return e.Client().Caches.Member(e.GuildID, e.UserID)
}

// GuildMessageReactionAdd indicates that a fluxer.Member added a fluxer.PartialEmoji to a fluxer.Message in a fluxer.GuildMessageChannel(requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionAdd struct {
	*GenericGuildMessageReaction
	Member          fluxer.Member
	MessageAuthorID *snowflake.ID
}

// GuildMessageReactionRemove indicates that a fluxer.Member removed a fluxer.MessageReaction from a fluxer.Message in a Channel (requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionRemove struct {
	*GenericGuildMessageReaction
}

// GuildMessageReactionRemoveEmoji indicates someone removed all fluxer.MessageReaction of a specific fluxer.Emoji from a fluxer.Message in a Channel (requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   snowflake.ID
	Emoji     fluxer.PartialEmoji
}

// GuildMessageReactionRemoveAll indicates someone removed all fluxer.MessageReaction(s) from a fluxer.Message in a Channel (requires the gateway.IntentGuildMessageReactions)
type GuildMessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   snowflake.ID
}
