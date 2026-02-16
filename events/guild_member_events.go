package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericGuildMember generic fluxer.Member event
type GenericGuildMember struct {
	*GenericEvent
	GuildID snowflake.ID
	Member  fluxer.Member
}

// GuildMemberJoin indicates that a fluxer.Member joined the fluxer.Guild
type GuildMemberJoin struct {
	*GenericGuildMember
}

// GuildMemberUpdate indicates that a fluxer.Member updated
type GuildMemberUpdate struct {
	*GenericGuildMember
	OldMember fluxer.Member
}

// GuildMemberLeave indicates that a fluxer.Member left the fluxer.Guild
type GuildMemberLeave struct {
	*GenericEvent
	GuildID snowflake.ID
	User    fluxer.User
	Member  fluxer.Member
}

// GuildMemberTypingStart indicates that a fluxer.Member started typing in a fluxer.BaseGuildMessageChannel(requires gateway.IntentGuildMessageTyping)
// Member will be empty when event is triggered by [Clyde bot]
//
// [Clyde bot]: https://support.fluxer.com/hc/en-us/articles/13066317497239-Clyde-Discord-s-AI-Chatbot
type GuildMemberTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	GuildID   snowflake.ID
	Timestamp time.Time
	Member    fluxer.Member
}

// Channel returns the fluxer.GuildMessageChannel the GuildMemberTypingStart happened in
func (e *GuildMemberTypingStart) Channel() (fluxer.GuildMessageChannel, bool) {
	return e.Client().Caches.GuildMessageChannel(e.ChannelID)
}
