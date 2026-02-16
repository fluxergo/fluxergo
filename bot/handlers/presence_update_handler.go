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

func gatewayHandlerPresenceUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventPresenceUpdate) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager.DispatchEvent(&events.PresenceUpdate{
		GenericEvent:        genericEvent,
		EventPresenceUpdate: event,
	})

	if client.Caches.CacheFlags().Missing(cache.FlagPresences) {
		return
	}

	var (
		oldStatus       fluxer.OnlineStatus
		oldClientStatus fluxer.ClientStatus
		oldActivities   []fluxer.Activity
	)

	if oldPresence, ok := client.Caches.Presence(event.GuildID, event.PresenceUser.ID); ok {
		oldStatus = oldPresence.Status
		oldClientStatus = oldPresence.ClientStatus
		oldActivities = oldPresence.Activities
	}

	client.Caches.AddPresence(event.Presence)

	if oldStatus != event.Status {
		client.EventManager.DispatchEvent(&events.UserStatusUpdate{
			GenericEvent: genericEvent,
			UserID:       event.PresenceUser.ID,
			OldStatus:    oldStatus,
			Status:       event.Status,
		})
	}

	if oldClientStatus.Desktop != event.ClientStatus.Desktop || oldClientStatus.Mobile != event.ClientStatus.Mobile || oldClientStatus.Web != event.ClientStatus.Web {
		client.EventManager.DispatchEvent(&events.UserClientStatusUpdate{
			GenericEvent:    genericEvent,
			UserID:          event.PresenceUser.ID,
			OldClientStatus: oldClientStatus,
			ClientStatus:    event.ClientStatus,
		})
	}

	genericUserActivityEvent := events.GenericUserActivity{
		GenericEvent: genericEvent,
		UserID:       event.PresenceUser.ID,
		GuildID:      event.GuildID,
	}

	for _, oldActivity := range oldActivities {
		var found bool
		for _, newActivity := range event.Activities {
			if oldActivity.ID == newActivity.ID {
				found = true
				break
			}
		}
		if !found {
			genericUserActivityEvent.Activity = oldActivity
			client.EventManager.DispatchEvent(&events.UserActivityStop{
				GenericUserActivity: &genericUserActivityEvent,
			})
		}
	}

	for _, newActivity := range event.Activities {
		var found bool
		for _, oldActivity := range oldActivities {
			if newActivity.ID == oldActivity.ID {
				found = true
				break
			}
		}
		if !found {
			genericUserActivityEvent.Activity = newActivity
			client.EventManager.DispatchEvent(&events.UserActivityStart{
				GenericUserActivity: &genericUserActivityEvent,
			})
		}
	}

	for _, newActivity := range event.Activities {
		var oldActivity *fluxer.Activity
		for _, activity := range oldActivities {
			if newActivity.ID == activity.ID {
				oldActivity = &activity
				break
			}
		}
		if oldActivity != nil && isActivityUpdated(*oldActivity, newActivity) {
			genericUserActivityEvent.Activity = newActivity
			client.EventManager.DispatchEvent(&events.UserActivityUpdate{
				GenericUserActivity: &genericUserActivityEvent,
				OldActivity:         *oldActivity,
			})
		}
	}
}

func isActivityUpdated(old fluxer.Activity, new fluxer.Activity) bool {
	return old.Name != new.Name ||
		old.Type != new.Type ||
		compareStringPtr(old.URL, new.URL) ||
		old.CreatedAt.Equal(new.CreatedAt) ||
		compareActivityTimestampsPtr(old.Timestamps, new.Timestamps) ||
		compareStringPtr(old.SyncID, new.SyncID) ||
		old.ApplicationID != new.ApplicationID ||
		compareStatusDisplayTypePtr(old.StatusDisplayType, new.StatusDisplayType) ||
		compareStringPtr(old.Details, new.Details) ||
		compareStringPtr(old.DetailsURL, new.DetailsURL) ||
		compareStringPtr(old.State, new.State) ||
		compareStringPtr(old.StateURL, new.StateURL) ||
		comparePartialEmojiPtr(old.Emoji, new.Emoji) ||
		compareActivityPartyPtr(old.Party, new.Party) ||
		compareActivityAssetsPtr(old.Assets, new.Assets) ||
		compareActivitySecretsPtr(old.Secrets, new.Secrets) ||
		compareBoolPtr(old.Instance, new.Instance) ||
		old.Flags != new.Flags ||
		slices.Equal(old.Buttons, new.Buttons)
}

func compareActivityTimestampsPtr(old *fluxer.ActivityTimestamps, new *fluxer.ActivityTimestamps) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return old.Start.Equal(new.Start) && old.End.Equal(new.End)
}

func compareBoolPtr(old *bool, new *bool) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return *old != *new
}

func compareStringPtr(old *string, new *string) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return *old != *new
}

func comparePartialEmojiPtr(old *fluxer.PartialEmoji, new *fluxer.PartialEmoji) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	if old.ID != new.ID {
		return true
	}
	if old.Name != new.Name {
		return true
	}
	return old.Animated != new.Animated
}

func compareSnowflakePtr(old *snowflake.ID, new *snowflake.ID) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return *old != *new
}

func compareActivityPartyPtr(old *fluxer.ActivityParty, new *fluxer.ActivityParty) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return old.ID != new.ID || old.Size != new.Size
}

func compareActivityAssetsPtr(old *fluxer.ActivityAssets, new *fluxer.ActivityAssets) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return old.LargeText != new.LargeText || old.LargeImage != new.LargeImage || old.LargeURL != new.LargeURL || old.SmallText != new.SmallText || old.SmallImage != new.SmallImage || old.SmallURL != new.SmallURL
}

func compareActivitySecretsPtr(old *fluxer.ActivitySecrets, new *fluxer.ActivitySecrets) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return old.Join != new.Join || old.Spectate != new.Spectate || old.Match != new.Match
}

func compareStatusDisplayTypePtr(old *fluxer.ActivityStatusDisplayType, new *fluxer.ActivityStatusDisplayType) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return *old != *new
}
