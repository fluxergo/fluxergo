package fluxer

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
)

type ChannelCreate interface {
	json.Marshaler
	Type() ChannelType
	channelCreate()
}

type GuildChannelCreate interface {
	ChannelCreate
	guildChannelCreate()
}

var (
	_ ChannelCreate      = (*GuildTextChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildTextChannelCreate)(nil)
)

type GuildTextChannelCreate struct {
	Name                          string                `json:"name"`
	Topic                         string                `json:"topic,omitempty"`
	RateLimitPerUser              int                   `json:"rate_limit_per_user,omitempty"`
	Position                      int                   `json:"position,omitempty"`
	PermissionOverwrites          []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      snowflake.ID          `json:"parent_id,omitempty"`
	NSFW                          bool                  `json:"nsfw,omitempty"`
	DefaultAutoArchiveDuration    AutoArchiveDuration   `json:"default_auto_archive_days,omitempty"`
	DefaultThreadRateLimitPerUser int                   `json:"default_thread_rate_limit_per_user,omitempty"`
}

func (c GuildTextChannelCreate) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c GuildTextChannelCreate) MarshalJSON() ([]byte, error) {
	type guildTextChannelCreate GuildTextChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildTextChannelCreate
	}{
		Type:                   c.Type(),
		guildTextChannelCreate: guildTextChannelCreate(c),
	})
}

func (GuildTextChannelCreate) channelCreate()      {}
func (GuildTextChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildVoiceChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildVoiceChannelCreate)(nil)
)

type GuildVoiceChannelCreate struct {
	Name                 string                `json:"name"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             snowflake.ID          `json:"parent_id,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
	RTCRegion            string                `json:"rtc_region,omitempty"`
	VideoQualityMode     VideoQualityMode      `json:"video_quality_mode,omitempty"`
}

func (c GuildVoiceChannelCreate) Type() ChannelType {
	return ChannelTypeGuildVoice
}

func (c GuildVoiceChannelCreate) MarshalJSON() ([]byte, error) {
	type guildVoiceChannelCreate GuildVoiceChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildVoiceChannelCreate
	}{
		Type:                    c.Type(),
		guildVoiceChannelCreate: guildVoiceChannelCreate(c),
	})
}

func (GuildVoiceChannelCreate) channelCreate()      {}
func (GuildVoiceChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildCategoryChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildCategoryChannelCreate)(nil)
)

type GuildCategoryChannelCreate struct {
	Name                 string                `json:"name"`
	Topic                string                `json:"topic,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
}

func (c GuildCategoryChannelCreate) Type() ChannelType {
	return ChannelTypeGuildCategory
}

func (c GuildCategoryChannelCreate) MarshalJSON() ([]byte, error) {
	type guildCategoryChannelCreate GuildCategoryChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildCategoryChannelCreate
	}{
		Type:                       c.Type(),
		guildCategoryChannelCreate: guildCategoryChannelCreate(c),
	})
}

func (GuildCategoryChannelCreate) channelCreate()      {}
func (GuildCategoryChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildCategoryChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildCategoryChannelCreate)(nil)
)

type GuildLinkExtendedChannelCreate struct {
	Name                 string                `json:"name"`
	URL                  string                `json:"url"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
}

func (c GuildLinkExtendedChannelCreate) Type() ChannelType {
	return ChannelTypeGuildLinkExtended
}

func (c GuildLinkExtendedChannelCreate) MarshalJSON() ([]byte, error) {
	type guildLinkExtendedChannelCreate GuildLinkExtendedChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildLinkExtendedChannelCreate
	}{
		Type:                           c.Type(),
		guildLinkExtendedChannelCreate: guildLinkExtendedChannelCreate(c),
	})
}

func (GuildLinkExtendedChannelCreate) channelCreate()      {}
func (GuildLinkExtendedChannelCreate) guildChannelCreate() {}

type DMChannelCreate struct {
	RecipientID snowflake.ID `json:"recipient_id"`
}
