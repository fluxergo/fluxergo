package handlers

import (
	"context"
	"log/slog"

	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

func gatewayHandlerGuildCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildCreate) {
	wasUnready := client.Caches.IsGuildUnready(event.ID)
	wasUnavailable := client.Caches.IsGuildUnavailable(event.ID)

	client.Caches.AddGuild(event.Guild)

	for _, channel := range event.Channels {
		channel = fluxer.ApplyGuildIDToChannel(channel, event.ID) // populate unset field
		client.Caches.AddChannel(channel)
	}

	for _, role := range event.Roles {
		role.GuildID = event.ID // populate unset field
		client.Caches.AddRole(role)
	}

	for _, member := range event.Members {
		member.GuildID = event.ID // populate unset field
		client.Caches.AddMember(member)
	}

	for _, voiceState := range event.VoiceStates {
		voiceState.GuildID = event.ID // populate unset field
		client.Caches.AddVoiceState(voiceState)
	}

	for _, emoji := range event.Emojis {
		emoji.GuildID = event.ID // populate unset field
		client.Caches.AddEmoji(emoji)
	}

	for _, sticker := range event.Stickers {
		sticker.GuildID = &event.ID // populate unset field
		client.Caches.AddSticker(sticker)
	}

	for _, guildScheduledEvent := range event.GuildScheduledEvents {
		client.Caches.AddGuildScheduledEvent(guildScheduledEvent)
	}

	for _, presence := range event.Presences {
		presence.GuildID = event.ID // populate unset field
		client.Caches.AddPresence(presence)
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
	}

	if wasUnready {
		client.Caches.SetGuildUnready(event.ID, false)
		client.EventManager.DispatchEvent(&events.GuildReady{
			GenericGuild: genericGuildEvent,
			Guild:        event.GatewayGuild,
		})
		if len(client.Caches.UnreadyGuildIDs()) == 0 {
			client.EventManager.DispatchEvent(&events.GuildsReady{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			})
		}
		if client.MemberChunkingManager.MemberChunkingFilter()(event.ID) {
			go func() {
				if _, err := client.MemberChunkingManager.RequestMembersWithQuery(context.Background(), event.ID, "", 0); err != nil {
					client.Logger.Error("failed to chunk guild on guild_create", slog.Any("err", err))
				}
			}()
		}

		return
	}
	if wasUnavailable {
		client.Caches.SetGuildUnavailable(event.ID, false)
		client.EventManager.DispatchEvent(&events.GuildAvailable{
			GenericGuild: genericGuildEvent,
			Guild:        event.GatewayGuild,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildJoin{
			GenericGuild: genericGuildEvent,
			Guild:        event.GatewayGuild,
		})
	}
}

func gatewayHandlerGuildUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildUpdate) {
	oldGuild, _ := client.Caches.Guild(event.ID)
	client.Caches.AddGuild(event.Guild)

	client.EventManager.DispatchEvent(&events.GuildUpdate{
		GenericGuild: &events.GenericGuild{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.ID,
		},
		Guild:    event.Guild,
		OldGuild: oldGuild,
	})
}

func gatewayHandlerGuildDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildDelete) {
	if event.Unavailable {
		client.Caches.SetGuildUnavailable(event.ID, true)
	}

	guild, _ := client.Caches.RemoveGuild(event.ID)
	client.Caches.RemoveVoiceStatesByGuildID(event.ID)
	client.Caches.RemovePresencesByGuildID(event.ID)
	client.Caches.RemoveChannelsByGuildID(event.ID)
	client.Caches.RemoveEmojisByGuildID(event.ID)
	client.Caches.RemoveStickersByGuildID(event.ID)
	client.Caches.RemoveRolesByGuildID(event.ID)
	client.Caches.RemoveMembersByGuildID(event.ID)
	client.Caches.RemoveGuildScheduledEventsByGuildID(event.ID)
	client.Caches.RemoveMessagesByGuildID(event.ID)

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
	}

	if event.Unavailable {
		client.EventManager.DispatchEvent(&events.GuildUnavailable{
			GenericGuild: genericGuildEvent,
			Guild:        guild,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildLeave{
			GenericGuild: genericGuildEvent,
			Guild:        guild,
		})
	}
}

func gatewayHandlerGuildAuditLogEntryCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildAuditLogEntryCreate) {
	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
	}

	client.EventManager.DispatchEvent(&events.GuildAuditLogEntryCreate{
		GenericGuild:  genericGuildEvent,
		AuditLogEntry: event.AuditLogEntry,
	})
}
