package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericUser is called upon receiving UserUpdate or UserTypingStart
type GenericUser struct {
	*GenericEvent
	UserID snowflake.ID
	User   fluxer.User
}

// UserUpdate  indicates that a fluxer.User updated
type UserUpdate struct {
	*GenericUser
	OldUser fluxer.User
}

// UserTypingStart indicates that a fluxer.User started typing in a fluxer.DMChannel or fluxer.MessageChanel(requires the gateway.IntentDirectMessageTyping and/or gateway.IntentGuildMessageTyping)
type UserTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	GuildID   *snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}

// Channel returns the fluxer.GuildMessageChannel the fluxer.User started typing in
func (e *UserTypingStart) Channel() (fluxer.GuildMessageChannel, bool) {
	return e.Client().Caches.GuildMessageChannel(e.ChannelID)
}
