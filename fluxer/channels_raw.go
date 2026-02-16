package fluxer

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type dmChannel struct {
	ID               snowflake.ID  `json:"id"`
	Type             ChannelType   `json:"type"`
	LastMessageID    *snowflake.ID `json:"last_message_id"`
	Recipients       []User        `json:"recipients"`
	LastPinTimestamp *time.Time    `json:"last_pin_timestamp"`
}

type groupDMChannel struct {
	ID               snowflake.ID  `json:"id"`
	Type             ChannelType   `json:"type"`
	OwnerID          *snowflake.ID `json:"owner_id"`
	Name             string        `json:"name"`
	LastPinTimestamp *time.Time    `json:"last_pin_timestamp"`
	LastMessageID    *snowflake.ID `json:"last_message_id"`
	Icon             *string       `json:"icon"`
}

type guildTextChannel struct {
	ID                         snowflake.ID          `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildID                    snowflake.ID          `json:"guild_id"`
	Position                   int                   `json:"position"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name"`
	Topic                      *string               `json:"topic"`
	NSFW                       bool                  `json:"nsfw"`
	LastMessageID              *snowflake.ID         `json:"last_message_id"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user"`
	ParentID                   *snowflake.ID         `json:"parent_id"`
	LastPinTimestamp           *time.Time            `json:"last_pin_timestamp"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
}

func (t *guildTextChannel) UnmarshalJSON(data []byte) error {
	type guildTextChannelAlias guildTextChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildTextChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildTextChannel(v.guildTextChannelAlias)
	t.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
}

type guildLinkExtended struct {
	ID                   snowflake.ID          `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.ID          `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	URL                  string                `json:"url"`
	ParentID             *snowflake.ID         `json:"parent_id"`
}

func (t *guildLinkExtended) UnmarshalJSON(data []byte) error {
	type guildLinkExtendedAlias guildLinkExtended
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildLinkExtendedAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildLinkExtended(v.guildLinkExtendedAlias)
	t.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
}

type guildCategoryChannel struct {
	ID                   snowflake.ID          `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.ID          `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
}

func (t *guildCategoryChannel) UnmarshalJSON(data []byte) error {
	type guildCategoryChannelAlias guildCategoryChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildCategoryChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildCategoryChannel(v.guildCategoryChannelAlias)
	t.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
}

type guildVoiceChannel struct {
	ID                   snowflake.ID          `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.ID          `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Bitrate              int                   `json:"bitrate"`
	UserLimit            int                   `json:"user_limit"`
	ParentID             *snowflake.ID         `json:"parent_id"`
	RTCRegion            string                `json:"rtc_region"`
	VideoQualityMode     VideoQualityMode      `json:"video_quality_mode"`
	LastMessageID        *snowflake.ID         `json:"last_message_id"`
	NSFW                 bool                  `json:"nsfw"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user"`
}

func (t *guildVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildVoiceChannelAlias guildVoiceChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildVoiceChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildVoiceChannel(v.guildVoiceChannelAlias)
	t.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
}

func parsePermissionOverwrites(overwrites []UnmarshalPermissionOverwrite) []PermissionOverwrite {
	if len(overwrites) == 0 {
		return nil
	}
	permOverwrites := make([]PermissionOverwrite, len(overwrites))
	for i := range overwrites {
		permOverwrites[i] = overwrites[i].PermissionOverwrite
	}
	return permOverwrites
}
