package fluxer

import (
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
)

type ChannelUpdate interface {
	channelUpdate()
}

type GuildChannelUpdate interface {
	ChannelUpdate
	guildChannelUpdate()
}

type GuildTextChannelUpdate struct {
	Name                          *string                `json:"name,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Topic                         *string                `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites          *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      *snowflake.ID          `json:"parent_id,omitempty"`
	DefaultAutoArchiveDuration    *AutoArchiveDuration   `json:"default_auto_archive_duration,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
}

func (GuildTextChannelUpdate) channelUpdate()      {}
func (GuildTextChannelUpdate) guildChannelUpdate() {}

type GuildVoiceChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	RateLimitPerUser     *int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate              *int                   `json:"bitrate,omitempty"`
	UserLimit            *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             *snowflake.ID          `json:"parent_id,omitempty"`
	RTCRegion            *string                `json:"rtc_region,omitempty"`
	NSFW                 *bool                  `json:"nsfw,omitempty"`
	VideoQualityMode     *VideoQualityMode      `json:"video_quality_mode,omitempty"`
}

func (GuildVoiceChannelUpdate) channelUpdate()      {}
func (GuildVoiceChannelUpdate) guildChannelUpdate() {}

type GuildCategoryChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
}

func (GuildCategoryChannelUpdate) channelUpdate()      {}
func (GuildCategoryChannelUpdate) guildChannelUpdate() {}

type GuildLinkExtendedChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	URL                  *string                `json:"url,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             *snowflake.ID          `json:"parent_id,omitempty"`
}

func (GuildLinkExtendedChannelUpdate) channelUpdate()      {}
func (GuildLinkExtendedChannelUpdate) guildChannelUpdate() {}

type GuildChannelPositionUpdate struct {
	ID              snowflake.ID     `json:"id"`
	Position        omit.Omit[*int]  `json:"position,omitzero"`
	LockPermissions omit.Omit[*bool] `json:"lock_permissions,omitzero"`
	ParentID        *snowflake.ID    `json:"parent_id,omitempty"`
}
