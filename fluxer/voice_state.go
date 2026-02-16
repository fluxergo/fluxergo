package fluxer

import (
	"github.com/disgoorg/snowflake/v2"
)

// VoiceState from Discord
type VoiceState struct {
	ChannelID       *snowflake.ID `json:"channel_id"`
	ConnectionID    *string       `json:"connection_id"`
	GuildDeaf       bool          `json:"deaf"`
	GuildID         snowflake.ID  `json:"guild_id,omitempty"`
	GuildMute       bool          `json:"mute"`
	SelfDeaf        bool          `json:"self_deaf"`
	SelfMute        bool          `json:"self_mute"`
	SelfStream      bool          `json:"self_stream"`
	SelfVideo       bool          `json:"self_video"`
	SessionID       string        `json:"session_id"`
	UserID          snowflake.ID  `json:"user_id"`
	Version         int           `json:"version"`
	ViewerStreamKey []string      `json:"viewer_stream_key"`
}
