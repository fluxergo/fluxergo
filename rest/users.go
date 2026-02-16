package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Users = (*userImpl)(nil)

func NewUsers(client Client) Users {
	return &userImpl{client: client}
}

type Users interface {
	GetUser(userID snowflake.ID, opts ...RequestOpt) (*fluxer.User, error)
	UpdateCurrentUser(userUpdate fluxer.UserUpdate, opts ...RequestOpt) (*fluxer.OAuth2User, error)
	LeaveGuild(guildID snowflake.ID, opts ...RequestOpt) error
	CreateDMChannel(userID snowflake.ID, opts ...RequestOpt) (*fluxer.DMChannel, error)
}

type userImpl struct {
	client Client
}

func (s *userImpl) GetUser(userID snowflake.ID, opts ...RequestOpt) (user *fluxer.User, err error) {
	err = s.client.Do(GetUser.Compile(nil, userID), nil, &user, opts...)
	return
}

func (s *userImpl) UpdateCurrentUser(userUpdate fluxer.UserUpdate, opts ...RequestOpt) (selfUser *fluxer.OAuth2User, err error) {
	err = s.client.Do(UpdateCurrentUser.Compile(nil), userUpdate, &selfUser, opts...)
	return
}

func (s *userImpl) LeaveGuild(guildID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(LeaveGuild.Compile(nil, guildID), nil, nil, opts...)
}

func (s *userImpl) CreateDMChannel(userID snowflake.ID, opts ...RequestOpt) (channel *fluxer.DMChannel, err error) {
	err = s.client.Do(CreateDMChannel.Compile(nil), fluxer.DMChannelCreate{RecipientID: userID}, &channel, opts...)
	return
}
