package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericDMMessage is called upon receiving DMMessageCreate , DMMessageUpdate , DMMessageDelete , GenericDMMessageReaction , DMMessageReactionAdd , DMMessageReactionRemove , DMMessageReactionRemoveEmoji or DMMessageReactionRemoveAll (requires gateway.IntentsDirectMessage)
type GenericDMMessage struct {
	*GenericEvent
	MessageID snowflake.ID
	Message   fluxer.Message
	ChannelID snowflake.ID
}

// DMMessageCreate is called upon receiving a fluxer.Message in a Channel (requires gateway.IntentsDirectMessage)
type DMMessageCreate struct {
	*GenericDMMessage
}

// DMMessageUpdate is called upon editing a fluxer.Message in a Channel (requires gateway.IntentsDirectMessage)
type DMMessageUpdate struct {
	*GenericDMMessage
	OldMessage fluxer.Message
}

// DMMessageDelete is called upon deleting a fluxer.Message in a Channel (requires gateway.IntentsDirectMessage)
type DMMessageDelete struct {
	*GenericDMMessage
}
