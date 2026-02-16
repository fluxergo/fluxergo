package fluxer

import (
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

var _ Mentionable = (*Emoji)(nil)

// Emoji allows you to interact with emojis & emotes
type Emoji struct {
	PartialEmoji
	GuildID snowflake.ID `json:"guild_id,omitempty"` // not present in the API but we need it
}

type PartialEmoji struct {
	ID       snowflake.ID `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"`
	Animated bool         `json:"animated,omitempty"`
}

// Reaction returns a string used for manipulating with reactions. May be empty if the Name is empty
func (e PartialEmoji) Reaction() string {
	if e.Name == "" {
		return ""
	}
	return reaction(e.Name, e.ID)
}

// Mention returns the string used to send the Emoji
func (e PartialEmoji) Mention() string {
	if e.Animated {
		return AnimatedEmojiMention(e.ID, e.Name)
	}
	return EmojiMention(e.ID, e.Name)
}

// String formats the Emoji as string
func (e PartialEmoji) String() string {
	return e.Mention()
}

func (e PartialEmoji) URL(opts ...CDNOpt) string {
	return formatAssetURL(CustomEmoji, opts, e.ID)
}

func (e PartialEmoji) CreatedAt() time.Time {
	if e.ID == 0 {
		return time.Time{}
	}
	return e.ID.Time()
}

type EmojiWithUser struct {
	Emoji
	User User `json:"user,omitempty"`
}

type EmojiCreate struct {
	Name  string `json:"name"`
	Image Icon   `json:"image"`
}

type EmojiUpdate struct {
	Name *string `json:"name,omitempty"`
}

func reaction(name string, id snowflake.ID) string {
	return fmt.Sprintf("%s:%s", name, id)
}
