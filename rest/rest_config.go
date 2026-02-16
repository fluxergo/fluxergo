package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

func defaultConfig() config {
	return config{
		DefaultAllowedMentions: fluxer.AllowedMentions{
			Parse:       []fluxer.AllowedMentionType{fluxer.AllowedMentionTypeUsers, fluxer.AllowedMentionTypeRoles, fluxer.AllowedMentionTypeEveryone},
			Roles:       []snowflake.ID{},
			Users:       []snowflake.ID{},
			RepliedUser: true,
		},
	}
}

type config struct {
	DefaultAllowedMentions fluxer.AllowedMentions
}

// ConfigOpt can be used to supply optional parameters to New
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithDefaultAllowedMentions(mentions fluxer.AllowedMentions) ConfigOpt {
	return func(config *config) {
		config.DefaultAllowedMentions = mentions
	}
}
