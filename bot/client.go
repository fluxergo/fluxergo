package bot

import (
	"context"
	"log/slog"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/cache"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
	"github.com/fluxergo/fluxergo/rest"
	"github.com/fluxergo/fluxergo/voice"
)

// Client is a high level struct for interacting with Discord's API.
// It combines the functionality of the rest, gateway/sharding, httpserver and cache into one easy to use package.
// Create a new client with fluxergo.New.
type Client struct {
	Token                 string
	ApplicationID         snowflake.ID
	Logger                *slog.Logger
	Rest                  rest.Rest
	EventManager          EventManager
	Gateway               gateway.Gateway
	VoiceManager          voice.Manager
	Caches                cache.Caches
	MemberChunkingManager MemberChunkingManager
}

func (c *Client) Close(ctx context.Context) {
	if c.VoiceManager != nil {
		c.VoiceManager.Close(ctx)
	}
	if c.Gateway != nil {
		c.Gateway.Close(ctx)
	}
	if c.Rest != nil {
		c.Rest.Close(ctx)
	}
}

func (c *Client) ID() snowflake.ID {
	if selfUser, ok := c.Caches.SelfUser(); ok {
		return selfUser.ID
	}
	return 0
}

func (c *Client) AddEventListeners(listeners ...EventListener) {
	c.EventManager.AddEventListeners(listeners...)
}

func (c *Client) RemoveEventListeners(listeners ...EventListener) {
	c.EventManager.RemoveEventListeners(listeners...)
}

func (c *Client) OpenGateway(ctx context.Context) error {
	if c.Gateway == nil {
		return fluxer.ErrNoGateway
	}
	return c.Gateway.Open(ctx)
}

func (c *Client) HasGateway() bool {
	return c.Gateway != nil
}

func (c *Client) shard() (gateway.Gateway, error) {
	if c.HasGateway() {
		return c.Gateway, nil
	}
	return nil, fluxer.ErrNoGateway
}

func (c *Client) UpdateVoiceState(ctx context.Context, data gateway.MessageDataVoiceStateUpdate) error {
	shard, err := c.shard()
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodeVoiceStateUpdate, data)
}

func (c *Client) RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error {
	shard, err := c.shard()
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodeRequestGuildMembers, gateway.MessageDataRequestGuildMembers{
		GuildID:   guildID,
		Presences: presence,
		UserIDs:   userIDs,
		Nonce:     nonce,
	})
}

func (c *Client) RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error {
	shard, err := c.shard()
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodeRequestGuildMembers, gateway.MessageDataRequestGuildMembers{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presence,
		Nonce:     nonce,
	})
}

func (c *Client) SetPresence(ctx context.Context, opts ...gateway.PresenceOpt) error {
	shard, err := c.shard()
	if err != nil {
		return err
	}
	return shard.Send(ctx, gateway.OpcodePresenceUpdate, applyPresenceFromOpts(shard, opts...))
}

func applyPresenceFromOpts(g gateway.Gateway, opts ...gateway.PresenceOpt) gateway.MessageDataPresenceUpdate {
	presenceUpdate := g.Presence()
	if presenceUpdate == nil {
		presenceUpdate = &gateway.MessageDataPresenceUpdate{}
	}
	for _, opt := range opts {
		opt(presenceUpdate)
	}
	return *presenceUpdate
}
