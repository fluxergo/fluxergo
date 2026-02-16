package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Webhooks = (*webhookImpl)(nil)

func NewWebhooks(client Client, defaultAllowedMentions fluxer.AllowedMentions) Webhooks {
	return &webhookImpl{client: client, defaultAllowedMentions: defaultAllowedMentions}
}

type Webhooks interface {
	GetWebhook(webhookID snowflake.ID, opts ...RequestOpt) (fluxer.Webhook, error)
	UpdateWebhook(webhookID snowflake.ID, webhookUpdate fluxer.WebhookUpdate, opts ...RequestOpt) (fluxer.Webhook, error)
	DeleteWebhook(webhookID snowflake.ID, opts ...RequestOpt) error

	GetWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) (fluxer.Webhook, error)
	UpdateWebhookWithToken(webhookID snowflake.ID, webhookToken string, webhookUpdate fluxer.WebhookUpdateWithToken, opts ...RequestOpt) (fluxer.Webhook, error)
	DeleteWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) error

	GetWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, opts ...RequestOpt) (*fluxer.Message, error)
	CreateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.WebhookMessageCreate, params CreateWebhookMessageParams, opts ...RequestOpt) (*fluxer.Message, error)
	CreateWebhookMessageSlack(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.Payload, params CreateWebhookMessageParams, opts ...RequestOpt) (*fluxer.Message, error)
	CreateWebhookMessageGitHub(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.Payload, params CreateWebhookMessageParams, opts ...RequestOpt) (*fluxer.Message, error)
	UpdateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, messageUpdate fluxer.WebhookMessageUpdate, opts ...RequestOpt) (*fluxer.Message, error)
	DeleteWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, opts ...RequestOpt) error
}

type CreateWebhookMessageParams struct {
	Wait bool
}

func (p CreateWebhookMessageParams) ToQueryValues() fluxer.QueryValues {
	queryValues := fluxer.QueryValues{}
	if p.Wait {
		queryValues["wait"] = true
	}
	return queryValues
}

type webhookImpl struct {
	client                 Client
	defaultAllowedMentions fluxer.AllowedMentions
}

func (s *webhookImpl) GetWebhook(webhookID snowflake.ID, opts ...RequestOpt) (webhook fluxer.Webhook, err error) {
	var unmarshalWebhook fluxer.UnmarshalWebhook
	err = s.client.Do(GetWebhook.Compile(nil, webhookID), nil, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) UpdateWebhook(webhookID snowflake.ID, webhookUpdate fluxer.WebhookUpdate, opts ...RequestOpt) (webhook fluxer.Webhook, err error) {
	var unmarshalWebhook fluxer.UnmarshalWebhook
	err = s.client.Do(UpdateWebhook.Compile(nil, webhookID), webhookUpdate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) DeleteWebhook(webhookID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteWebhook.Compile(nil, webhookID), nil, nil, opts...)
}

func (s *webhookImpl) GetWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) (webhook fluxer.Webhook, err error) {
	var unmarshalWebhook fluxer.UnmarshalWebhook
	err = s.client.Do(GetWebhookWithToken.Compile(nil, webhookID, webhookToken), nil, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) UpdateWebhookWithToken(webhookID snowflake.ID, webhookToken string, webhookUpdate fluxer.WebhookUpdateWithToken, opts ...RequestOpt) (webhook fluxer.Webhook, err error) {
	var unmarshalWebhook fluxer.UnmarshalWebhook
	err = s.client.Do(UpdateWebhookWithToken.Compile(nil, webhookID, webhookToken), webhookUpdate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) DeleteWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) error {
	return s.client.Do(DeleteWebhookWithToken.Compile(nil, webhookID, webhookToken), nil, nil, opts...)
}

func (s *webhookImpl) GetWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, opts ...RequestOpt) (message *fluxer.Message, err error) {
	err = s.client.Do(GetWebhookMessage.Compile(nil, webhookID, webhookToken, messageID), nil, &message, opts...)
	return
}

func (s *webhookImpl) createWebhookMessage(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.Payload, params CreateWebhookMessageParams, endpoint *Endpoint, opts []RequestOpt) (message *fluxer.Message, err error) {
	compiledEndpoint := endpoint.Compile(params.ToQueryValues(), webhookID, webhookToken)

	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}

	if params.Wait {
		err = s.client.Do(compiledEndpoint, body, &message, opts...)
	} else {
		err = s.client.Do(compiledEndpoint, body, nil, opts...)
	}
	return
}

func (s *webhookImpl) CreateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.WebhookMessageCreate, params CreateWebhookMessageParams, opts ...RequestOpt) (*fluxer.Message, error) {
	if messageCreate.AllowedMentions == nil {
		messageCreate.AllowedMentions = &s.defaultAllowedMentions
	}

	return s.createWebhookMessage(webhookID, webhookToken, messageCreate, params, CreateWebhookMessage, opts)
}

func (s *webhookImpl) CreateWebhookMessageSlack(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.Payload, params CreateWebhookMessageParams, opts ...RequestOpt) (*fluxer.Message, error) {
	return s.createWebhookMessage(webhookID, webhookToken, messageCreate, params, CreateWebhookMessageSlack, opts)
}

func (s *webhookImpl) CreateWebhookMessageGitHub(webhookID snowflake.ID, webhookToken string, messageCreate fluxer.Payload, params CreateWebhookMessageParams, opts ...RequestOpt) (*fluxer.Message, error) {
	return s.createWebhookMessage(webhookID, webhookToken, messageCreate, params, CreateWebhookMessageGitHub, opts)
}

func (s *webhookImpl) UpdateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, messageUpdate fluxer.WebhookMessageUpdate, opts ...RequestOpt) (message *fluxer.Message, err error) {
	if messageUpdate.AllowedMentions == nil && (messageUpdate.Content != nil || (messageUpdate.Flags != nil && messageUpdate.Flags.Has(fluxer.MessageFlagIsComponentsV2))) {
		messageUpdate.AllowedMentions = &s.defaultAllowedMentions
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}

	err = s.client.Do(UpdateWebhookMessage.Compile(nil, webhookID, webhookToken, messageID), body, &message, opts...)
	return
}

func (s *webhookImpl) DeleteWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteWebhookMessage.Compile(nil, webhookID, webhookToken, messageID), nil, nil, opts...)
}
