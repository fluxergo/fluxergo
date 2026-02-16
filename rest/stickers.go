package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Stickers = (*stickerImpl)(nil)

func NewStickers(client Client) Stickers {
	return &stickerImpl{client: client}
}

type Stickers interface {
	GetSticker(stickerID snowflake.ID, opts ...RequestOpt) (*fluxer.Sticker, error)
	GetStickers(guildID snowflake.ID, opts ...RequestOpt) ([]fluxer.Sticker, error)
	CreateSticker(guildID snowflake.ID, createSticker fluxer.StickerCreate, opts ...RequestOpt) (*fluxer.Sticker, error)
	UpdateSticker(guildID snowflake.ID, stickerID snowflake.ID, stickerUpdate fluxer.StickerUpdate, opts ...RequestOpt) (*fluxer.Sticker, error)
	DeleteSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) error
}

type stickerImpl struct {
	client Client
}

func (s *stickerImpl) GetSticker(stickerID snowflake.ID, opts ...RequestOpt) (sticker *fluxer.Sticker, err error) {
	err = s.client.Do(GetSticker.Compile(nil, stickerID), nil, &sticker, opts...)
	return
}

func (s *stickerImpl) GetStickers(guildID snowflake.ID, opts ...RequestOpt) (stickers []fluxer.Sticker, err error) {
	err = s.client.Do(GetGuildStickers.Compile(nil, guildID), nil, &stickers, opts...)
	return
}

func (s *stickerImpl) CreateSticker(guildID snowflake.ID, createSticker fluxer.StickerCreate, opts ...RequestOpt) (sticker *fluxer.Sticker, err error) {
	body, err := createSticker.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(CreateGuildSticker.Compile(nil, guildID), body, &sticker, opts...)
	return
}

func (s *stickerImpl) UpdateSticker(guildID snowflake.ID, stickerID snowflake.ID, stickerUpdate fluxer.StickerUpdate, opts ...RequestOpt) (sticker *fluxer.Sticker, err error) {
	err = s.client.Do(UpdateGuildSticker.Compile(nil, guildID, stickerID), stickerUpdate, &sticker, opts...)
	return
}

func (s *stickerImpl) DeleteSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGuildSticker.Compile(nil, guildID, stickerID), nil, nil, opts...)
}
