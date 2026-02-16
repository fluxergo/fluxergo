package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericRole generic fluxer.Role event
type GenericRole struct {
	*GenericEvent
	GuildID snowflake.ID
	RoleID  snowflake.ID
	Role    fluxer.Role
}

// RoleCreate indicates that a fluxer.Role got created
type RoleCreate struct {
	*GenericRole
}

// RoleUpdate indicates that a fluxer.Role got updated
type RoleUpdate struct {
	*GenericRole
	OldRole fluxer.Role
}

// RoleDelete indicates that a fluxer.Role got deleted
type RoleDelete struct {
	*GenericRole
}
