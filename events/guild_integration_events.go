package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericIntegration is called upon receiving IntegrationCreate, IntegrationUpdate or IntegrationDelete(requires the gateway.IntentGuildIntegrations)
type GenericIntegration struct {
	*GenericEvent
	GuildID     snowflake.ID
	Integration fluxer.Integration
}

// Guild returns the Guild this Integration was created in.
// This will only check cached guilds!
func (e *GenericIntegration) Guild() (fluxer.Guild, bool) {
	return e.Client().Caches.Guild(e.GuildID)
}

// IntegrationCreate indicates that a new Integration was created in a Guild
type IntegrationCreate struct {
	*GenericIntegration
}

// IntegrationUpdate indicates that an integration was updated in a Guild
type IntegrationUpdate struct {
	*GenericIntegration
}

// IntegrationDelete indicates that an Integration was deleted from a Guild
type IntegrationDelete struct {
	*GenericEvent
	ID            snowflake.ID
	GuildID       snowflake.ID
	ApplicationID *snowflake.ID
}

// GuildIntegrationsUpdate indicates that a Guild's integrations were updated
type GuildIntegrationsUpdate struct {
	*GenericEvent
	GuildID snowflake.ID
}
