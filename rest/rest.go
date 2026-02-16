package rest

// Rest is a manager for all of fluxergo's HTTP requests
type Rest interface {
	Client

	Applications
	OAuth2
	Gateway
	Guilds
	Members
	Channels
	Invites
	Users
	Webhooks
	Emojis
	Stickers
	GuildScheduledEvents
}

var _ Rest = (*restImpl)(nil)

// New returns a new default Rest
func New(client Client, opts ...ConfigOpt) Rest {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &restImpl{
		Client:               client,
		Applications:         NewApplications(client),
		OAuth2:               NewOAuth2(client),
		Gateway:              NewGateway(client),
		Guilds:               NewGuilds(client),
		Members:              NewMembers(client),
		Channels:             NewChannels(client, cfg.DefaultAllowedMentions),
		Invites:              NewInvites(client),
		Users:                NewUsers(client),
		Webhooks:             NewWebhooks(client, cfg.DefaultAllowedMentions),
		Emojis:               NewEmojis(client),
		Stickers:             NewStickers(client),
		GuildScheduledEvents: NewGuildScheduledEvents(client),
	}
}

type restImpl struct {
	Client

	Applications
	OAuth2
	Gateway
	Guilds
	Members
	Channels
	Invites
	Users
	Webhooks
	Emojis
	Stickers
	GuildScheduledEvents
}
