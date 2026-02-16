package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericReaction is called upon receiving MessageReactionAdd or MessageReactionRemove
type GenericReaction struct {
	*GenericEvent
	UserID      snowflake.ID
	ChannelID   snowflake.ID
	MessageID   snowflake.ID
	GuildID     *snowflake.ID
	Emoji       fluxer.PartialEmoji
	BurstColors []string
	Burst       bool
}

// MessageReactionAdd indicates that a fluxer.User added a fluxer.MessageReaction to a fluxer.Message in a fluxer.Channel(this+++ requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionAdd struct {
	*GenericReaction
	Member *fluxer.Member
}

// MessageReactionRemove indicates that a fluxer.User removed a fluxer.MessageReaction from a fluxer.Message in a fluxer.GetChannel(requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionRemove struct {
	*GenericReaction
}

// MessageReactionRemoveEmoji indicates someone removed all fluxer.MessageReaction of a specific fluxer.Emoji from a fluxer.Message in a fluxer.Channel(requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
	Emoji     fluxer.PartialEmoji
}

// MessageReactionRemoveAll indicates someone removed all fluxer.MessageReaction(s) from a fluxer.Message in a fluxer.Channel(requires the gateway.IntentGuildMessageReactions and/or gateway.IntentDirectMessageReactions)
type MessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
}
