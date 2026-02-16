package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// DMChannelPinsUpdate indicates that a fluxer.Message got pinned or unpinned.
type DMChannelPinsUpdate struct {
	*GenericEvent
	ChannelID           snowflake.ID
	NewLastPinTimestamp *time.Time
}

// DMUserTypingStart indicates that a fluxer.User started typing in a fluxer.DMChannel(requires gateway.IntentDirectMessageTyping).
type DMUserTypingStart struct {
	*GenericEvent
	ChannelID snowflake.ID
	UserID    snowflake.ID
	Timestamp time.Time
}
