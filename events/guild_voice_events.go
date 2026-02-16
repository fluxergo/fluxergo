package events

import (
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

// GenericGuildVoiceState is called upon receiving GuildVoiceJoin, GuildVoiceMove and GuildVoiceLeave
type GenericGuildVoiceState struct {
	*GenericEvent
	VoiceState fluxer.VoiceState
	Member     fluxer.Member
}

// GuildVoiceStateUpdate indicates that the fluxer.VoiceState of a fluxer.Member has updated (requires gateway.IntentGuildVoiceStates)
type GuildVoiceStateUpdate struct {
	*GenericGuildVoiceState
	OldVoiceState fluxer.VoiceState
}

// GuildVoiceJoin indicates that a fluxer.Member joined a fluxer.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceJoin struct {
	*GenericGuildVoiceState
}

// GuildVoiceMove indicates that a fluxer.Member was moved to a different fluxer.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceMove struct {
	*GenericGuildVoiceState
	OldVoiceState fluxer.VoiceState
}

// GuildVoiceLeave indicates that a fluxer.Member left a fluxer.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceLeave struct {
	*GenericGuildVoiceState
	OldVoiceState fluxer.VoiceState
}

// VoiceServerUpdate indicates that a voice server the bot is connected to has been changed
type VoiceServerUpdate struct {
	*GenericEvent
	gateway.EventVoiceServerUpdate
}
