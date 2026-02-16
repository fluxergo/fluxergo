package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

// StickersUpdate is dispatched when a guild's stickers are updated.
// This event does not depend on a cache like StickerCreate, StickerUpdate or StickerDelete.
type StickersUpdate struct {
	*GenericEvent
	gateway.EventGuildStickersUpdate
}

// GenericSticker is called upon receiving StickerCreate , StickerUpdate or StickerDelete (requires gateway.IntentGuildExpressions)
type GenericSticker struct {
	*GenericEvent
	GuildID snowflake.ID
	Sticker fluxer.Sticker
}

// StickerCreate indicates that a new fluxer.Sticker got created in a fluxer.Guild (requires gateway.IntentGuildExpressions)
type StickerCreate struct {
	*GenericSticker
}

// StickerUpdate indicates that a fluxer.Sticker got updated in a fluxer.Guild (requires gateway.IntentGuildExpressions)
type StickerUpdate struct {
	*GenericSticker
	OldSticker fluxer.Sticker
}

// StickerDelete indicates that a fluxer.Sticker got deleted in a fluxer.Guild (requires gateway.IntentGuildExpressions)
type StickerDelete struct {
	*GenericSticker
}
