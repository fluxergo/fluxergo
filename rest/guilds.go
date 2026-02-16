package rest

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/internal/slicehelper"
)

var _ Guilds = (*guildImpl)(nil)

func NewGuilds(client Client) Guilds {
	return &guildImpl{client: client}
}

type Guilds interface {
	GetGuild(guildID snowflake.ID, withCounts bool, opts ...RequestOpt) (*fluxer.RestGuild, error)
	UpdateGuild(guildID snowflake.ID, guildUpdate fluxer.GuildUpdate, opts ...RequestOpt) (*fluxer.RestGuild, error)

	GetGuildVanityURL(guildID snowflake.ID, opts ...RequestOpt) (*fluxer.PartialInvite, error)

	CreateGuildChannel(guildID snowflake.ID, guildChannelCreate fluxer.GuildChannelCreate, opts ...RequestOpt) (fluxer.GuildChannel, error)
	GetGuildChannels(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.GuildChannel, error)
	UpdateChannelPositions(guildID snowflake.ID, guildChannelPositionUpdates []fluxer.GuildChannelPositionUpdate, opts ...RequestOpt) error

	GetRoles(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.Role, error)
	GetRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (*fluxer.Role, error)
	CreateRole(guildID snowflake.ID, createRole fluxer.RoleCreate, opts ...RequestOpt) (*fluxer.Role, error)
	UpdateRole(guildID snowflake.ID, roleID snowflake.ID, roleUpdate fluxer.RoleUpdate, opts ...RequestOpt) (*fluxer.Role, error)
	UpdateRolePositions(guildID snowflake.ID, rolePositionUpdates []fluxer.RolePositionUpdate, opts ...RequestOpt) ([]fluxer.Role, error)
	DeleteRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error
	GetRoleMemberCounts(guildID snowflake.ID, opts ...RequestOpt) (map[snowflake.ID]int, error)

	GetBans(guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) ([]fluxer.Ban, error)
	GetBansPage(guildID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[fluxer.Ban]
	GetBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*fluxer.Ban, error)
	AddBan(guildID snowflake.ID, userID snowflake.ID, deleteMessageDuration time.Duration, opts ...RequestOpt) error
	DeleteBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	BulkBan(guildID snowflake.ID, ban fluxer.BulkBan, opts ...RequestOpt) (*fluxer.BulkBanResult, error)

	GetIntegrations(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.Integration, error)
	DeleteIntegration(guildID snowflake.ID, integrationID snowflake.ID, opts ...RequestOpt) error

	GetGuildPruneCount(guildID snowflake.ID, days int, includeRoles []snowflake.ID, opts ...RequestOpt) (*fluxer.GuildPruneResult, error)
	BeginGuildPrune(guildID snowflake.ID, guildPrune fluxer.GuildPrune, opts ...RequestOpt) (*fluxer.GuildPruneResult, error)

	GetAllWebhooks(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.Webhook, error)

	GetGuildVoiceRegions(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.VoiceRegion, error)

	GetAuditLog(guildID snowflake.ID, userID snowflake.ID, actionType fluxer.AuditLogEvent, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (*fluxer.AuditLog, error)
	GetAuditLogPage(guildID snowflake.ID, userID snowflake.ID, actionType fluxer.AuditLogEvent, startID snowflake.ID, limit int, opts ...RequestOpt) AuditLogPage
}

type guildImpl struct {
	client Client
}

func (s *guildImpl) GetGuild(guildID snowflake.ID, withCounts bool, opts ...RequestOpt) (guild *fluxer.RestGuild, err error) {
	values := fluxer.QueryValues{
		"with_counts": withCounts,
	}
	err = s.client.Do(GetGuild.Compile(values, guildID), nil, &guild, opts...)
	return
}

func (s *guildImpl) UpdateGuild(guildID snowflake.ID, guildUpdate fluxer.GuildUpdate, opts ...RequestOpt) (guild *fluxer.RestGuild, err error) {
	err = s.client.Do(UpdateGuild.Compile(nil, guildID), guildUpdate, &guild, opts...)
	return
}

func (s *guildImpl) GetGuildVanityURL(guildID snowflake.ID, opts ...RequestOpt) (partialInvite *fluxer.PartialInvite, err error) {
	err = s.client.Do(GetGuildVanityURL.Compile(nil, guildID), nil, &partialInvite, opts...)
	return
}

func (s *guildImpl) CreateGuildChannel(guildID snowflake.ID, guildChannelCreate fluxer.GuildChannelCreate, opts ...RequestOpt) (guildChannel fluxer.GuildChannel, err error) {
	var ch fluxer.UnmarshalChannel
	err = s.client.Do(CreateGuildChannel.Compile(nil, guildID), guildChannelCreate, &ch, opts...)
	if err == nil {
		guildChannel = ch.Channel.(fluxer.GuildChannel)
	}
	return
}

func (s *guildImpl) GetGuildChannels(guildID snowflake.ID, opts ...RequestOpt) (channels []fluxer.GuildChannel, err error) {
	var chs []fluxer.UnmarshalChannel
	err = s.client.Do(GetGuildChannels.Compile(nil, guildID), nil, &chs, opts...)
	if err == nil {
		channels = make([]fluxer.GuildChannel, len(chs))
		for i := range chs {
			channels[i] = chs[i].Channel.(fluxer.GuildChannel)
		}
	}
	return
}

func (s *guildImpl) UpdateChannelPositions(guildID snowflake.ID, guildChannelPositionUpdates []fluxer.GuildChannelPositionUpdate, opts ...RequestOpt) error {
	return s.client.Do(UpdateChannelPositions.Compile(nil, guildID), guildChannelPositionUpdates, nil, opts...)
}

func (s *guildImpl) GetRoles(guildID snowflake.ID, opts ...RequestOpt) (roles []fluxer.Role, err error) {
	err = s.client.Do(GetRoles.Compile(nil, guildID), nil, &roles, opts...)
	if err == nil {
		for i := range roles {
			roles[i].GuildID = guildID
		}
	}
	return
}

func (s *guildImpl) GetRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (role *fluxer.Role, err error) {
	err = s.client.Do(GetRole.Compile(nil, guildID, roleID), nil, &role, opts...)
	if err == nil {
		role.GuildID = guildID
	}
	return
}

func (s *guildImpl) CreateRole(guildID snowflake.ID, createRole fluxer.RoleCreate, opts ...RequestOpt) (role *fluxer.Role, err error) {
	err = s.client.Do(CreateRole.Compile(nil, guildID), createRole, &role, opts...)
	if err == nil {
		role.GuildID = guildID
	}
	return
}

func (s *guildImpl) UpdateRole(guildID snowflake.ID, roleID snowflake.ID, roleUpdate fluxer.RoleUpdate, opts ...RequestOpt) (role *fluxer.Role, err error) {
	err = s.client.Do(UpdateRole.Compile(nil, guildID, roleID), roleUpdate, &role, opts...)
	if err == nil {
		role.GuildID = guildID
	}
	return
}

func (s *guildImpl) UpdateRolePositions(guildID snowflake.ID, rolePositionUpdates []fluxer.RolePositionUpdate, opts ...RequestOpt) (roles []fluxer.Role, err error) {
	err = s.client.Do(UpdateRolePositions.Compile(nil, guildID), rolePositionUpdates, &roles, opts...)
	if err == nil {
		for i := range roles {
			roles[i].GuildID = guildID
		}
	}
	return
}

func (s *guildImpl) DeleteRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteRole.Compile(nil, guildID, roleID), nil, nil, opts...)
}

func (s *guildImpl) GetRoleMemberCounts(guildID snowflake.ID, opts ...RequestOpt) (memberCounts map[snowflake.ID]int, err error) {
	err = s.client.Do(GetRoleMemberCounts.Compile(nil, guildID), nil, &memberCounts, opts...)
	return
}

func (s *guildImpl) GetBans(guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (bans []fluxer.Ban, err error) {
	values := fluxer.QueryValues{}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetBans.Compile(values, guildID), nil, &bans, opts...)
	return
}

func (s *guildImpl) GetBansPage(guildID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[fluxer.Ban] {
	return Page[fluxer.Ban]{
		getItemsFunc: func(before snowflake.ID, after snowflake.ID) (bans []fluxer.Ban, err error) {
			return s.GetBans(guildID, before, after, limit, opts...)
		},
		getIDFunc: func(ban fluxer.Ban) snowflake.ID {
			return ban.User.ID
		},
		ID: startID,
	}
}

func (s *guildImpl) GetBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (ban *fluxer.Ban, err error) {
	err = s.client.Do(GetBan.Compile(nil, guildID, userID), nil, &ban, opts...)
	return
}

func (s *guildImpl) AddBan(guildID snowflake.ID, userID snowflake.ID, deleteMessageDuration time.Duration, opts ...RequestOpt) error {
	return s.client.Do(AddBan.Compile(nil, guildID, userID), fluxer.AddBan{DeleteMessageSeconds: int(deleteMessageDuration.Seconds())}, nil, opts...)
}

func (s *guildImpl) DeleteBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteBan.Compile(nil, guildID, userID), nil, nil, opts...)
}

func (s *guildImpl) BulkBan(guildID snowflake.ID, ban fluxer.BulkBan, opts ...RequestOpt) (result *fluxer.BulkBanResult, err error) {
	err = s.client.Do(BulkBan.Compile(nil, guildID), ban, &result, opts...)
	return
}

func (s *guildImpl) GetIntegrations(guildID snowflake.ID, opts ...RequestOpt) (integrations []fluxer.Integration, err error) {
	err = s.client.Do(GetIntegrations.Compile(nil, guildID), nil, &integrations, opts...)
	return
}

func (s *guildImpl) DeleteIntegration(guildID snowflake.ID, integrationID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteIntegration.Compile(nil, guildID, integrationID), nil, nil, opts...)
}

func (s *guildImpl) GetGuildPruneCount(guildID snowflake.ID, days int, includeRoles []snowflake.ID, opts ...RequestOpt) (result *fluxer.GuildPruneResult, err error) {
	values := fluxer.QueryValues{
		"days": days,
	}
	if len(includeRoles) > 0 {
		values["include_roles"] = slicehelper.JoinSnowflakes(includeRoles)
	}
	err = s.client.Do(GetGuildPruneCount.Compile(values, guildID), nil, &result, opts...)
	return
}

func (s *guildImpl) BeginGuildPrune(guildID snowflake.ID, guildPrune fluxer.GuildPrune, opts ...RequestOpt) (result *fluxer.GuildPruneResult, err error) {
	err = s.client.Do(BeginGuildPrune.Compile(nil, guildID), guildPrune, &result, opts...)
	return
}

func (s *guildImpl) GetAllWebhooks(guildID snowflake.ID, opts ...RequestOpt) (webhooks []fluxer.Webhook, err error) {
	var whs []fluxer.UnmarshalWebhook
	err = s.client.Do(GetGuildWebhooks.Compile(nil, guildID), nil, &whs, opts...)
	if err == nil {
		webhooks = make([]fluxer.Webhook, len(whs))
		for i := range whs {
			webhooks[i] = whs[i].Webhook
		}
	}
	return
}

func (s *guildImpl) GetGuildVoiceRegions(guildID snowflake.ID, opts ...RequestOpt) (regions []fluxer.VoiceRegion, err error) {
	err = s.client.Do(GetGuildVoiceRegions.Compile(nil, guildID), nil, &regions, opts...)
	return
}

func (s *guildImpl) GetAuditLog(guildID snowflake.ID, userID snowflake.ID, actionType fluxer.AuditLogEvent, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (auditLog *fluxer.AuditLog, err error) {
	values := fluxer.QueryValues{}
	if userID != 0 {
		values["user_id"] = userID
	}
	if actionType != 0 {
		values["action_type"] = actionType
	}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetAuditLogs.Compile(values, guildID), nil, &auditLog, opts...)
	return
}

func (s *guildImpl) GetAuditLogPage(guildID snowflake.ID, userID snowflake.ID, actionType fluxer.AuditLogEvent, startID snowflake.ID, limit int, opts ...RequestOpt) AuditLogPage {
	return AuditLogPage{
		getItems: func(before snowflake.ID, after snowflake.ID) (fluxer.AuditLog, error) {
			log, err := s.GetAuditLog(guildID, userID, actionType, before, after, limit, opts...)
			var finalLog fluxer.AuditLog
			if log != nil {
				finalLog = *log
			}
			return finalLog, err
		},
		ID: startID,
	}
}
