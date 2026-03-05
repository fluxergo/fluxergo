package fluxer

import (
	"github.com/disgoorg/snowflake/v2"
)

// AllowedMentions are used for avoiding mentioning users in Message and Interaction
type AllowedMentions struct {
	Parse       []AllowedMentionType `json:"parse,omitzero"`
	Roles       []snowflake.ID       `json:"roles,omitzero"`
	Users       []snowflake.ID       `json:"users,omitzero"`
	RepliedUser bool                 `json:"replied_user"`
}

// AllowedMentionType ?
type AllowedMentionType string

// All AllowedMentionType(s)
const (
	AllowedMentionTypeRoles    AllowedMentionType = "roles"
	AllowedMentionTypeUsers    AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)