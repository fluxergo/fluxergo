package handlers

import (
	"slices"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/cache"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

type updatedSticker struct {
	old fluxer.Sticker
	new fluxer.Sticker
}

func gatewayHandlerGuildStickersUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildStickersUpdate) {
	client.EventManager.DispatchEvent(&events.StickersUpdate{
		GenericEvent:             events.NewGenericEvent(client, sequenceNumber, shardID),
		EventGuildStickersUpdate: event,
	})

	if client.Caches.CacheFlags().Missing(cache.FlagStickers) {
		return
	}

	createdStickers := map[snowflake.ID]fluxer.Sticker{}
	deletedStickers := map[snowflake.ID]fluxer.Sticker{}
	updatedStickers := map[snowflake.ID]updatedSticker{}

	for sticker := range client.Caches.Stickers(event.GuildID) {
		deletedStickers[sticker.ID] = sticker
	}

	for _, newSticker := range event.Stickers {
		oldSticker, ok := deletedStickers[newSticker.ID]
		if ok {
			delete(deletedStickers, newSticker.ID)
			if isStickerUpdated(oldSticker, newSticker) {
				updatedStickers[newSticker.ID] = updatedSticker{new: newSticker, old: oldSticker}
			}
			continue
		}
		createdStickers[newSticker.ID] = newSticker
	}

	for _, emoji := range createdStickers {
		client.EventManager.DispatchEvent(&events.StickerCreate{
			GenericSticker: &events.GenericSticker{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Sticker:      emoji,
			},
		})
	}

	for _, emoji := range updatedStickers {
		client.EventManager.DispatchEvent(&events.StickerUpdate{
			GenericSticker: &events.GenericSticker{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Sticker:      emoji.new,
			},
			OldSticker: emoji.old,
		})
	}

	for _, emoji := range deletedStickers {
		client.EventManager.DispatchEvent(&events.StickerDelete{
			GenericSticker: &events.GenericSticker{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Sticker:      emoji,
			},
		})
	}
}

func isStickerUpdated(old fluxer.Sticker, new fluxer.Sticker) bool {
	if old.Name != new.Name {
		return true
	}
	if old.Description != new.Description {
		return true
	}
	if !slices.Equal(old.Tags, new.Tags) {
		return true
	}
	return false
}
