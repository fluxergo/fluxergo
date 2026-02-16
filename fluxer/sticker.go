package fluxer

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// Sticker is a sticker sent with a Message
type Sticker struct {
	ID          snowflake.ID  `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        []string      `json:"tags"`
	Animated    bool          `json:"animated"`
	GuildID     *snowflake.ID `json:"guild_id"`
}

func (s Sticker) URL(opts ...CDNOpt) string {
	format := FileFormatWebP
	if s.Animated {
		format = FileFormatGIF
	}
	return formatAssetURL(CustomSticker, append(opts, WithFormat(format)), s.ID)
}

func (s Sticker) CreatedAt() time.Time {
	return s.ID.Time()
}

type StickerType int

const (
	StickerTypeStandard StickerType = iota + 1
	StickerTypeGuild
)

// StickerFormatType is the Format type of Sticker
type StickerFormatType int

// Constants for StickerFormatType
const (
	StickerFormatTypePNG StickerFormatType = iota + 1
	StickerFormatTypeAPNG
	StickerFormatTypeLottie
	StickerFormatTypeGIF
)

type StickerCreate struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags"`
	File        *File    `json:"-"`
}

// ToBody returns the MessageCreate ready for body
func (c StickerCreate) ToBody() (any, error) {
	if c.File != nil {
		return PayloadWithFiles(c, c.File)
	}
	return c, nil
}

type StickerUpdate struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}
