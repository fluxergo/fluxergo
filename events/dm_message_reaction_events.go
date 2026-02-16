package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericDMMessageReaction is called upon receiving DMMessageReactionAdd or DMMessageReactionRemove (requires the gateway.IntentDirectMessageReactions)
type GenericDMMessageReaction struct {
	*GenericEvent
	UserID      snowflake.ID
	ChannelID   snowflake.ID
	MessageID   snowflake.ID
	Emoji       fluxer.PartialEmoji
	BurstColors []string
	Burst       bool
}

// DMMessageReactionAdd indicates that a fluxer.User added a fluxer.MessageReaction to a fluxer.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionAdd struct {
	*GenericDMMessageReaction
	MessageAuthorID *snowflake.ID
}

// DMMessageReactionRemove indicates that a fluxer.User removed a fluxer.MessageReaction from a fluxer.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionRemove struct {
	*GenericDMMessageReaction
}

// DMMessageReactionRemoveEmoji indicates someone removed all fluxer.MessageReaction(s) of a specific fluxer.Emoji from a fluxer.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionRemoveEmoji struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
	Emoji     fluxer.PartialEmoji
}

// DMMessageReactionRemoveAll indicates someone removed all fluxer.MessageReaction(s) from a fluxer.Message in a Channel (requires the gateway.IntentDirectMessageReactions)
type DMMessageReactionRemoveAll struct {
	*GenericEvent
	ChannelID snowflake.ID
	MessageID snowflake.ID
}
