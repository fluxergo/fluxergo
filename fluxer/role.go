package fluxer

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

var _ Mentionable = (*Role)(nil)

// Role is a Guild Role object
type Role struct {
	ID            snowflake.ID `json:"id"`
	Name          string       `json:"name"`
	Color         int          `json:"color"`
	Position      int          `json:"position"`
	HoistPosition *int         `json:"hoist_position"`
	Permissions   Permissions  `json:"permissions"`
	Hoist         bool         `json:"hoist"`
	Mentionable   bool         `json:"mentionable"`
	UnicodeEmoji  *string      `json:"unicode_emoji,omitempty"`
	GuildID       snowflake.ID `json:"guild_id,omitempty"` // not present in the API but we need it
}

func (Role) isMentionableValue() {}

func (r Role) String() string {
	return RoleMention(r.ID)
}

func (r Role) Mention() string {
	return r.String()
}

func (r Role) CreatedAt() time.Time {
	return r.ID.Time()
}

// RoleCreate is the payload to create a Role
type RoleCreate struct {
	Name         string       `json:"name,omitempty"`
	Color        int          `json:"color,omitempty"`
	Permissions  *Permissions `json:"permissions,omitempty"`
	Hoist        bool         `json:"hoist,omitempty"`
	Mentionable  bool         `json:"mentionable,omitempty"`
	UnicodeEmoji *string      `json:"unicode_emoji,omitempty"`
}

// RoleUpdate is the payload to update a Role
type RoleUpdate struct {
	Name         *string      `json:"name,omitempty"`
	Permissions  *Permissions `json:"permissions,omitempty"`
	Color        *int         `json:"color,omitempty"`
	Hoist        *bool        `json:"hoist,omitempty"`
	Mentionable  *bool        `json:"mentionable,omitempty"`
	UnicodeEmoji *string      `json:"unicode_emoji,omitempty"`
}

// RolePositionUpdate is the payload to update a Role(s) position
type RolePositionUpdate struct {
	ID       snowflake.ID `json:"id"`
	Position *int         `json:"position,omitempty"`
}
