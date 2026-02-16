package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// GenericGuild represents a generic guild event
type GenericGuild struct {
	*GenericEvent
	GuildID snowflake.ID
}

// GuildUpdate is called upon receiving fluxer.Guild updates
type GuildUpdate struct {
	*GenericGuild
	Guild    fluxer.Guild
	OldGuild fluxer.Guild // the old cached guild
}

// GuildAvailable is called when an unavailable fluxer.Guild becomes available
type GuildAvailable struct {
	*GenericGuild
	Guild fluxer.GatewayGuild
}

// GuildUnavailable is called when an available fluxer.Guild becomes unavailable
type GuildUnavailable struct {
	*GenericGuild
	Guild fluxer.Guild // the old cached guild
}

// GuildJoin is called when the bot joins a fluxer.Guild
type GuildJoin struct {
	*GenericGuild
	Guild fluxer.GatewayGuild
}

// GuildLeave is called when the bot leaves a fluxer.Guild
type GuildLeave struct {
	*GenericGuild
	Guild fluxer.Guild // the old cached guild
}

// GuildReady is called when a fluxer.Guild becomes loaded for the first time
type GuildReady struct {
	*GenericGuild
	Guild fluxer.GatewayGuild
}

// GuildsReady is called when all fluxer.Guild(s) are loaded after logging in
type GuildsReady struct {
	*GenericEvent
}

// GuildBan is called when a fluxer.Member/fluxer.User is banned from the fluxer.Guild
type GuildBan struct {
	*GenericGuild
	User fluxer.User
}

// GuildUnban is called when a fluxer.Member/fluxer.User is unbanned from the fluxer.Guild
type GuildUnban struct {
	*GenericGuild
	User fluxer.User
}

// GuildAuditLogEntryCreate is called when a new fluxer.AuditLogEntry is created
type GuildAuditLogEntryCreate struct {
	*GenericGuild
	AuditLogEntry fluxer.AuditLogEntry
}
