package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/internal/slicehelper"
)

var _ Applications = (*applicationsImpl)(nil)

func NewApplications(client Client) Applications {
	return &applicationsImpl{client: client}
}

type Applications interface {
	GetCurrentApplication(opts ...RequestOpt) (*fluxer.Application, error)
	UpdateCurrentApplication(applicationUpdate fluxer.ApplicationUpdate, opts ...RequestOpt) (*fluxer.Application, error)

	GetActivityInstance(applicationID snowflake.ID, instanceID string, opts ...RequestOpt) (*fluxer.ActivityInstance, error)
}

// GetEntitlementsParams holds query parameters for Applications.GetEntitlements (https://fluxer.com/developers/docs/resources/entitlement#list-entitlements)
type GetEntitlementsParams struct {
	UserID         snowflake.ID
	SkuIDs         []snowflake.ID
	Before         snowflake.ID
	After          snowflake.ID
	Limit          int
	GuildID        snowflake.ID
	ExcludeEnded   bool
	ExcludeDeleted bool
}

func (p GetEntitlementsParams) ToQueryValues() fluxer.QueryValues {
	queryValues := fluxer.QueryValues{
		"exclude_ended":   p.ExcludeEnded,
		"exclude_deleted": p.ExcludeDeleted,
	}
	if len(p.SkuIDs) > 0 {
		queryValues["sku_ids"] = slicehelper.JoinSnowflakes(p.SkuIDs)
	}
	if p.UserID != 0 {
		queryValues["user_id"] = p.UserID
	}
	if p.Before != 0 {
		queryValues["before"] = p.Before
	}
	if p.After != 0 {
		queryValues["after"] = p.After
	}
	if p.Limit != 0 {
		queryValues["limit"] = p.Limit
	}
	if p.GuildID != 0 {
		queryValues["guild_id"] = p.GuildID
	}
	return queryValues
}

type applicationsImpl struct {
	client Client
}

func (s *applicationsImpl) GetCurrentApplication(opts ...RequestOpt) (application *fluxer.Application, err error) {
	err = s.client.Do(GetCurrentApplication.Compile(nil), nil, &application, opts...)
	return
}

func (s *applicationsImpl) UpdateCurrentApplication(applicationUpdate fluxer.ApplicationUpdate, opts ...RequestOpt) (application *fluxer.Application, err error) {
	err = s.client.Do(UpdateCurrentApplication.Compile(nil), applicationUpdate, &application, opts...)
	return
}

func (s *applicationsImpl) GetActivityInstance(applicationID snowflake.ID, instanceID string, opts ...RequestOpt) (instance *fluxer.ActivityInstance, err error) {
	err = s.client.Do(GetActivityInstance.Compile(nil, applicationID, instanceID), nil, &instance, opts...)
	return
}
