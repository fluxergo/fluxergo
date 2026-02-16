package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

// EmojisUpdate is dispatched when a guild's emojis are updated.
// This event does not depend on a cache like EmojiCreate, EmojiUpdate or EmojiDelete.
type EmojisUpdate struct {
	*GenericEvent
	gateway.EventGuildEmojisUpdate
}

// GenericEmoji is called upon receiving EmojiCreate , EmojiUpdate or EmojiDelete (requires gateway.IntentGuildExpressions)
type GenericEmoji struct {
	*GenericEvent
	GuildID snowflake.ID
	Emoji   fluxer.Emoji
}

// EmojiCreate indicates that a new fluxer.Emoji got created in a fluxer.Guild (requires gateway.IntentGuildExpressions)
type EmojiCreate struct {
	*GenericEmoji
}

// EmojiUpdate indicates that a fluxer.Emoji got updated in a fluxer.Guild (requires gateway.IntentGuildExpressions)
type EmojiUpdate struct {
	*GenericEmoji
	OldEmoji fluxer.Emoji
}

// EmojiDelete indicates that a fluxer.Emoji got deleted in a fluxer.Guild (requires gateway.IntentGuildExpressions)
type EmojiDelete struct {
	*GenericEmoji
}
