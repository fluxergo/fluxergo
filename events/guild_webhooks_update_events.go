package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// WebhooksUpdate indicates that a guilds webhooks were updated.
type WebhooksUpdate struct {
	*GenericEvent
	GuildId   snowflake.ID
	ChannelID snowflake.ID
}

// Guild returns the Guild the webhook was updated in.
// This will only return cached guilds!
func (e *WebhooksUpdate) Guild() (fluxer.Guild, bool) {
	return e.Client().Caches.Guild(e.GuildId)
}

// Channel returns the fluxer.GuildMessageChannel the webhook was updated in.
// This will only return cached channels!
func (e *WebhooksUpdate) Channel() (fluxer.GuildMessageChannel, bool) {
	return e.Client().Caches.GuildMessageChannel(e.ChannelID)
}
