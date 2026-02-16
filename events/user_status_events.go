package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// UserStatusUpdate generic Status event
type UserStatusUpdate struct {
	*GenericEvent
	UserID    snowflake.ID
	OldStatus fluxer.OnlineStatus
	Status    fluxer.OnlineStatus
}

// UserClientStatusUpdate generic client-specific Status event
type UserClientStatusUpdate struct {
	*GenericEvent
	UserID          snowflake.ID
	OldClientStatus fluxer.ClientStatus
	ClientStatus    fluxer.ClientStatus
}
