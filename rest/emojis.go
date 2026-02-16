package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Emojis = (*emojiImpl)(nil)

func NewEmojis(client Client) Emojis {
	return &emojiImpl{client: client}
}

type Emojis interface {
	GetEmojis(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.Emoji, error)
	GetEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (*fluxer.Emoji, error)
	CreateEmoji(guildID snowflake.ID, emojiCreate fluxer.EmojiCreate, opts ...RequestOpt) (*fluxer.Emoji, error)
	UpdateEmoji(guildID snowflake.ID, emojiID snowflake.ID, emojiUpdate fluxer.EmojiUpdate, opts ...RequestOpt) (*fluxer.Emoji, error)
	DeleteEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) error
}

type emojiImpl struct {
	client Client
}

func (s *emojiImpl) GetEmojis(guildID snowflake.ID, opts ...RequestOpt) (emojis []fluxer.Emoji, err error) {
	err = s.client.Do(GetEmojis.Compile(nil, guildID), nil, &emojis, opts...)
	for i := range emojis {
		emojis[i].GuildID = guildID
	}
	return
}

func (s *emojiImpl) GetEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (emoji *fluxer.Emoji, err error) {
	err = s.client.Do(GetEmoji.Compile(nil, guildID, emojiID), nil, &emoji, opts...)
	if emoji != nil {
		emoji.GuildID = guildID
	}
	return
}

func (s *emojiImpl) CreateEmoji(guildID snowflake.ID, emojiCreate fluxer.EmojiCreate, opts ...RequestOpt) (emoji *fluxer.Emoji, err error) {
	err = s.client.Do(CreateEmoji.Compile(nil, guildID), emojiCreate, &emoji, opts...)
	if emoji != nil {
		emoji.GuildID = guildID
	}
	return
}

func (s *emojiImpl) UpdateEmoji(guildID snowflake.ID, emojiID snowflake.ID, emojiUpdate fluxer.EmojiUpdate, opts ...RequestOpt) (emoji *fluxer.Emoji, err error) {
	err = s.client.Do(UpdateEmoji.Compile(nil, guildID, emojiID), emojiUpdate, &emoji, opts...)
	if emoji != nil {
		emoji.GuildID = guildID
	}
	return
}

func (s *emojiImpl) DeleteEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteEmoji.Compile(nil, guildID, emojiID), nil, nil, opts...)
}
