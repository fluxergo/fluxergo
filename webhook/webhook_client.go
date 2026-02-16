package webhook

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/rest"
)

var ErrInvalidWebhookURL = errors.New("invalid webhook URL")

// New creates a new Client with the given ID, Token and ConfigOpt(s).
func New(id snowflake.ID, token string, opts ...ConfigOpt) *Client {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &Client{
		ID:         id,
		Token:      token,
		Rest:       cfg.Webhooks,
		RestClient: cfg.RestClient,
		logger:     cfg.Logger,
	}
}

// NewWithURL creates a new Client by parsing the given webhookURL for the ID and Token.
func NewWithURL(webhookURL string, opts ...ConfigOpt) (*Client, error) {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}

	parts := strings.FieldsFunc(u.Path, func(r rune) bool { return r == '/' })
	if len(parts) != 3 {
		return nil, ErrInvalidWebhookURL
	}

	token := parts[2]
	id, err := snowflake.Parse(parts[1])
	if err != nil {
		return nil, err
	}

	return New(id, token, opts...), nil
}

type Client struct {
	ID         snowflake.ID
	Token      string
	Rest       rest.Webhooks
	RestClient rest.Client
	logger     *slog.Logger
}

func (c *Client) URL() string {
	return fluxer.WebhookURL(c.ID, c.Token)
}

func (c *Client) Close(ctx context.Context) {
	c.RestClient.Close(ctx)
}

func (c *Client) GetWebhook(opts ...rest.RequestOpt) (*fluxer.IncomingWebhook, error) {
	webhook, err := c.Rest.GetWebhookWithToken(c.ID, c.Token, opts...)
	if incomingWebhook, ok := webhook.(fluxer.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

func (c *Client) UpdateWebhook(webhookUpdate fluxer.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*fluxer.IncomingWebhook, error) {
	webhook, err := c.Rest.UpdateWebhookWithToken(c.ID, c.Token, webhookUpdate, opts...)
	if incomingWebhook, ok := webhook.(fluxer.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

func (c *Client) DeleteWebhook(opts ...rest.RequestOpt) error {
	return c.Rest.DeleteWebhookWithToken(c.ID, c.Token, opts...)
}

func (c *Client) CreateMessage(messageCreate fluxer.WebhookMessageCreate, params rest.CreateWebhookMessageParams, opts ...rest.RequestOpt) (*fluxer.Message, error) {
	return c.Rest.CreateWebhookMessage(c.ID, c.Token, messageCreate, params, opts...)
}

func (c *Client) CreateContent(content string, opts ...rest.RequestOpt) (*fluxer.Message, error) {
	return c.CreateMessage(fluxer.WebhookMessageCreate{Content: content}, rest.CreateWebhookMessageParams{}, opts...)
}

func (c *Client) CreateEmbeds(embeds []fluxer.Embed, opts ...rest.RequestOpt) (*fluxer.Message, error) {
	return c.CreateMessage(fluxer.WebhookMessageCreate{Embeds: embeds}, rest.CreateWebhookMessageParams{}, opts...)
}
