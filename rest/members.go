package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Members = (*memberImpl)(nil)

func NewMembers(client Client) Members {
	return &memberImpl{client: client}
}

type Members interface {
	GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*fluxer.Member, error)
	GetMembers(guildID snowflake.ID, limit int, after snowflake.ID, opts ...RequestOpt) ([]fluxer.Member, error)
	SearchMembers(guildID snowflake.ID, query string, limit int, opts ...RequestOpt) ([]fluxer.Member, error)
	AddMember(guildID snowflake.ID, userID snowflake.ID, memberAdd fluxer.MemberAdd, opts ...RequestOpt) (*fluxer.Member, error)
	RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error
	UpdateMember(guildID snowflake.ID, userID snowflake.ID, memberUpdate fluxer.MemberUpdate, opts ...RequestOpt) (*fluxer.Member, error)

	AddMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error
	RemoveMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error

	UpdateCurrentMember(guildID snowflake.ID, memberUpdate fluxer.CurrentMemberUpdate, opts ...RequestOpt) (*fluxer.Member, error)
}

type memberImpl struct {
	client Client
}

func (s *memberImpl) GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (member *fluxer.Member, err error) {
	err = s.client.Do(GetMember.Compile(nil, guildID, userID), nil, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) GetMembers(guildID snowflake.ID, limit int, after snowflake.ID, opts ...RequestOpt) (members []fluxer.Member, err error) {
	values := fluxer.QueryValues{
		"limit": limit,
		"after": after,
	}
	err = s.client.Do(GetMembers.Compile(values, guildID), nil, &members, opts...)
	if err == nil {
		for i := range members {
			members[i].GuildID = guildID
		}
	}
	return
}

func (s *memberImpl) SearchMembers(guildID snowflake.ID, query string, limit int, opts ...RequestOpt) (members []fluxer.Member, err error) {
	values := fluxer.QueryValues{}
	if query != "" {
		values["query"] = query
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(SearchMembers.Compile(values, guildID), nil, &members, opts...)
	if err == nil {
		for i := range members {
			members[i].GuildID = guildID
		}
	}
	return
}

func (s *memberImpl) AddMember(guildID snowflake.ID, userID snowflake.ID, memberAdd fluxer.MemberAdd, opts ...RequestOpt) (member *fluxer.Member, err error) {
	err = s.client.Do(AddMember.Compile(nil, guildID, userID), memberAdd, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveMember.Compile(nil, guildID, userID), nil, nil, opts...)
}

func (s *memberImpl) UpdateMember(guildID snowflake.ID, userID snowflake.ID, memberUpdate fluxer.MemberUpdate, opts ...RequestOpt) (member *fluxer.Member, err error) {
	err = s.client.Do(UpdateMember.Compile(nil, guildID, userID), memberUpdate, &member, opts...)
	if err == nil {
		member.GuildID = guildID
	}
	return
}

func (s *memberImpl) AddMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(AddMemberRole.Compile(nil, guildID, userID, roleID), nil, nil, opts...)
}

func (s *memberImpl) RemoveMemberRole(guildID snowflake.ID, userID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveMemberRole.Compile(nil, guildID, userID, roleID), nil, nil, opts...)
}

func (s *memberImpl) UpdateCurrentMember(guildID snowflake.ID, memberUpdate fluxer.CurrentMemberUpdate, opts ...RequestOpt) (member *fluxer.Member, err error) {
	err = s.client.Do(UpdateCurrentMember.Compile(nil, guildID), memberUpdate, &member, opts...)
	return
}
