package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

// InviteCreate is called upon creation of a new fluxer.Invite (requires gateway.IntentGuildInvites)
type InviteCreate struct {
	*GenericEvent

	gateway.EventInviteCreate
}

// Channel returns the fluxer.GuildChannel the GenericInvite happened in.
func (e *InviteCreate) Channel() (fluxer.GuildChannel, bool) {
	return e.Client().Caches.Channel(e.ChannelID)
}

func (e *InviteCreate) Guild() (fluxer.Guild, bool) {
	if e.GuildID == nil {
		return fluxer.Guild{}, false
	}
	return e.Client().Caches.Guild(*e.GuildID)
}

// InviteDelete is called upon deletion of a fluxer.Invite (requires gateway.IntentGuildInvites)
type InviteDelete struct {
	*GenericEvent

	GuildID   *snowflake.ID
	ChannelID snowflake.ID
	Code      string
}

// Channel returns the fluxer.GuildChannel the GenericInvite happened in.
func (e *InviteDelete) Channel() (fluxer.GuildChannel, bool) {
	return e.Client().Caches.Channel(e.ChannelID)
}

func (e *InviteDelete) Guild() (fluxer.Guild, bool) {
	if e.GuildID == nil {
		return fluxer.Guild{}, false
	}
	return e.Client().Caches.Guild(*e.GuildID)
}
