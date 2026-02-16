package handlers

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/cache"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

type updatedEmoji struct {
	old fluxer.Emoji
	new fluxer.Emoji
}

func gatewayHandlerGuildEmojisUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildEmojisUpdate) {
	client.EventManager.DispatchEvent(&events.EmojisUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventGuildEmojisUpdate: event,
	})

	if client.Caches.CacheFlags().Missing(cache.FlagEmojis) {
		return
	}

	createdEmojis := map[snowflake.ID]fluxer.Emoji{}
	deletedEmojis := map[snowflake.ID]fluxer.Emoji{}
	updatedEmojis := map[snowflake.ID]updatedEmoji{}

	for emoji := range client.Caches.Emojis(event.GuildID) {
		deletedEmojis[emoji.ID] = emoji
	}

	for _, newEmoji := range event.Emojis {
		oldEmoji, ok := deletedEmojis[newEmoji.ID]
		if ok {
			delete(deletedEmojis, newEmoji.ID)
			if isEmojiUpdated(oldEmoji, newEmoji) {
				updatedEmojis[newEmoji.ID] = updatedEmoji{new: newEmoji, old: oldEmoji}
			}
			continue
		}
		createdEmojis[newEmoji.ID] = newEmoji
	}

	for _, emoji := range createdEmojis {
		client.Caches.AddEmoji(emoji)
		client.EventManager.DispatchEvent(&events.EmojiCreate{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		client.Caches.AddEmoji(emoji.new)
		client.EventManager.DispatchEvent(&events.EmojiUpdate{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji.new,
			},
			OldEmoji: emoji.old,
		})
	}

	for _, emoji := range deletedEmojis {
		client.Caches.RemoveEmoji(event.GuildID, emoji.ID)
		client.EventManager.DispatchEvent(&events.EmojiDelete{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji,
			},
		})
	}

}

func isEmojiUpdated(old fluxer.Emoji, new fluxer.Emoji) bool {
	if old.Name != new.Name {
		return true
	}
	return false
}
