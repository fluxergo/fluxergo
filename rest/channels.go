package rest

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Channels = (*channelImpl)(nil)

func NewChannels(client Client, defaultAllowedMentions fluxer.AllowedMentions) Channels {
	return &channelImpl{client: client, defaultAllowedMentions: defaultAllowedMentions}
}

type Channels interface {
	GetChannel(channelID snowflake.ID, opts ...RequestOpt) (fluxer.Channel, error)
	UpdateChannel(channelID snowflake.ID, channelUpdate fluxer.ChannelUpdate, opts ...RequestOpt) (fluxer.Channel, error)
	DeleteChannel(channelID snowflake.ID, opts ...RequestOpt) error

	GetWebhooks(channelID snowflake.ID, opts ...RequestOpt) ([]fluxer.Webhook, error)
	CreateWebhook(channelID snowflake.ID, webhookCreate fluxer.WebhookCreate, opts ...RequestOpt) (*fluxer.IncomingWebhook, error)

	UpdatePermissionOverwrite(channelID snowflake.ID, overwriteID snowflake.ID, permissionOverwrite fluxer.PermissionOverwriteUpdate, opts ...RequestOpt) error
	DeletePermissionOverwrite(channelID snowflake.ID, overwriteID snowflake.ID, opts ...RequestOpt) error

	SendTyping(channelID snowflake.ID, opts ...RequestOpt) error

	GetMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (*fluxer.Message, error)
	GetMessages(channelID snowflake.ID, around snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) ([]fluxer.Message, error)
	GetMessagesPage(channelID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[fluxer.Message]
	CreateMessage(channelID snowflake.ID, messageCreate fluxer.MessageCreate, opts ...RequestOpt) (*fluxer.Message, error)
	UpdateMessage(channelID snowflake.ID, messageID snowflake.ID, messageUpdate fluxer.MessageUpdate, opts ...RequestOpt) (*fluxer.Message, error)
	DeleteMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error
	BulkDeleteMessages(channelID snowflake.ID, messageIDs []snowflake.ID, opts ...RequestOpt) error
	CrosspostMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (*fluxer.Message, error)

	GetReactions(channelID snowflake.ID, messageID snowflake.ID, emoji string, reactionType fluxer.MessageReactionType, after int, limit int, opts ...RequestOpt) ([]fluxer.User, error)
	AddReaction(channelID snowflake.ID, messageID snowflake.ID, emoji string, opts ...RequestOpt) error
	RemoveOwnReaction(channelID snowflake.ID, messageID snowflake.ID, emoji string, opts ...RequestOpt) error
	RemoveUserReaction(channelID snowflake.ID, messageID snowflake.ID, emoji string, userID snowflake.ID, opts ...RequestOpt) error
	RemoveAllReactions(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error
	RemoveAllReactionsForEmoji(channelID snowflake.ID, messageID snowflake.ID, emoji string, opts ...RequestOpt) error

	GetChannelPins(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (*fluxer.ChannelPins, error)
	GetChannelPinsPage(channelID snowflake.ID, start time.Time, limit int, opts ...RequestOpt) ChannelPinsPage
	PinMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error
	UnpinMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error

	Follow(channelID snowflake.ID, targetChannelID snowflake.ID, opts ...RequestOpt) (*fluxer.FollowedChannel, error)
}

type channelImpl struct {
	client                 Client
	defaultAllowedMentions fluxer.AllowedMentions
}

func (s *channelImpl) GetChannel(channelID snowflake.ID, opts ...RequestOpt) (channel fluxer.Channel, err error) {
	var ch fluxer.UnmarshalChannel
	err = s.client.Do(GetChannel.Compile(nil, channelID), nil, &ch, opts...)
	if err == nil {
		channel = ch.Channel
	}
	return
}

func (s *channelImpl) UpdateChannel(channelID snowflake.ID, channelUpdate fluxer.ChannelUpdate, opts ...RequestOpt) (channel fluxer.Channel, err error) {
	var ch fluxer.UnmarshalChannel
	err = s.client.Do(UpdateChannel.Compile(nil, channelID), channelUpdate, &ch, opts...)
	if err == nil {
		channel = ch.Channel
	}
	return
}

func (s *channelImpl) DeleteChannel(channelID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteChannel.Compile(nil, channelID), nil, nil, opts...)
}

func (s *channelImpl) GetWebhooks(channelID snowflake.ID, opts ...RequestOpt) (webhooks []fluxer.Webhook, err error) {
	var whs []fluxer.UnmarshalWebhook
	err = s.client.Do(GetChannelWebhooks.Compile(nil, channelID), nil, &whs, opts...)
	if err == nil {
		webhooks = make([]fluxer.Webhook, len(whs))
		for i := range whs {
			webhooks[i] = whs[i].Webhook
		}
	}
	return
}

func (s *channelImpl) CreateWebhook(channelID snowflake.ID, webhookCreate fluxer.WebhookCreate, opts ...RequestOpt) (webhook *fluxer.IncomingWebhook, err error) {
	err = s.client.Do(CreateWebhook.Compile(nil, channelID), webhookCreate, &webhook, opts...)
	return
}

func (s *channelImpl) UpdatePermissionOverwrite(channelID snowflake.ID, overwriteID snowflake.ID, permissionOverwrite fluxer.PermissionOverwriteUpdate, opts ...RequestOpt) error {
	return s.client.Do(UpdatePermissionOverwrite.Compile(nil, channelID, overwriteID), permissionOverwrite, nil, opts...)
}

func (s *channelImpl) DeletePermissionOverwrite(channelID snowflake.ID, overwriteID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeletePermissionOverwrite.Compile(nil, channelID, overwriteID), nil, nil, opts...)
}

func (s *channelImpl) SendTyping(channelID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(SendTyping.Compile(nil, channelID), nil, nil, opts...)
}

func (s *channelImpl) GetMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (message *fluxer.Message, err error) {
	err = s.client.Do(GetMessage.Compile(nil, channelID, messageID), nil, &message, opts...)
	return
}

func (s *channelImpl) GetMessages(channelID snowflake.ID, around snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (messages []fluxer.Message, err error) {
	values := fluxer.QueryValues{}
	if around != 0 {
		values["around"] = around
	}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetMessages.Compile(values, channelID), nil, &messages, opts...)
	return
}

func (s *channelImpl) GetMessagesPage(channelID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[fluxer.Message] {
	return Page[fluxer.Message]{
		getItemsFunc: func(before snowflake.ID, after snowflake.ID) ([]fluxer.Message, error) {
			return s.GetMessages(channelID, 0, before, after, limit, opts...)
		},
		getIDFunc: func(msg fluxer.Message) snowflake.ID {
			return msg.ID
		},
		ID: startID,
	}
}

func (s *channelImpl) CreateMessage(channelID snowflake.ID, messageCreate fluxer.MessageCreate, opts ...RequestOpt) (message *fluxer.Message, err error) {
	if messageCreate.AllowedMentions == nil {
		messageCreate.AllowedMentions = &s.defaultAllowedMentions
	}
	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(CreateMessage.Compile(nil, channelID), body, &message, opts...)
	return
}

func (s *channelImpl) UpdateMessage(channelID snowflake.ID, messageID snowflake.ID, messageUpdate fluxer.MessageUpdate, opts ...RequestOpt) (message *fluxer.Message, err error) {
	if messageUpdate.AllowedMentions == nil && (messageUpdate.Content != nil || (messageUpdate.Flags != nil && messageUpdate.Flags.Has(fluxer.MessageFlagIsComponentsV2))) {
		messageUpdate.AllowedMentions = &s.defaultAllowedMentions
	}
	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(UpdateMessage.Compile(nil, channelID, messageID), body, &message, opts...)
	return
}

func (s *channelImpl) DeleteMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteMessage.Compile(nil, channelID, messageID), nil, nil, opts...)
}

func (s *channelImpl) BulkDeleteMessages(channelID snowflake.ID, messageIDs []snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(BulkDeleteMessages.Compile(nil, channelID), fluxer.MessageBulkDelete{Messages: messageIDs}, nil, opts...)
}

func (s *channelImpl) CrosspostMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (message *fluxer.Message, err error) {
	err = s.client.Do(CrosspostMessage.Compile(nil, channelID, messageID), nil, &message, opts...)
	return
}

func (s *channelImpl) GetReactions(channelID snowflake.ID, messageID snowflake.ID, emoji string, reactionType fluxer.MessageReactionType, after int, limit int, opts ...RequestOpt) (users []fluxer.User, err error) {
	values := fluxer.QueryValues{
		"type": reactionType,
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetReactions.Compile(values, channelID, messageID, emoji), nil, &users, opts...)
	return
}

func (s *channelImpl) AddReaction(channelID snowflake.ID, messageID snowflake.ID, emoji string, opts ...RequestOpt) error {
	return s.client.Do(AddReaction.Compile(nil, channelID, messageID, emoji), nil, nil, opts...)
}

func (s *channelImpl) RemoveOwnReaction(channelID snowflake.ID, messageID snowflake.ID, emoji string, opts ...RequestOpt) error {
	return s.client.Do(RemoveOwnReaction.Compile(nil, channelID, messageID, emoji), nil, nil, opts...)
}

func (s *channelImpl) RemoveUserReaction(channelID snowflake.ID, messageID snowflake.ID, emoji string, userID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveUserReaction.Compile(nil, channelID, messageID, emoji, userID), nil, nil, opts...)
}

func (s *channelImpl) RemoveAllReactions(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(RemoveAllReactions.Compile(nil, channelID, messageID), nil, nil, opts...)
}

func (s *channelImpl) RemoveAllReactionsForEmoji(channelID snowflake.ID, messageID snowflake.ID, emoji string, opts ...RequestOpt) error {
	return s.client.Do(RemoveAllReactionsForEmoji.Compile(nil, channelID, messageID, emoji), nil, nil, opts...)
}

func (s *channelImpl) GetChannelPins(channelID snowflake.ID, before time.Time, limit int, opts ...RequestOpt) (pins *fluxer.ChannelPins, err error) {
	values := fluxer.QueryValues{}
	if !before.IsZero() {
		values["before"] = before.Format(time.RFC3339)
	}
	if limit != 0 {
		values["limit"] = limit
	}
	err = s.client.Do(GetChannelPins.Compile(values, channelID), nil, &pins, opts...)
	return
}

func (s *channelImpl) GetChannelPinsPage(channelID snowflake.ID, start time.Time, limit int, opts ...RequestOpt) ChannelPinsPage {
	return ChannelPinsPage{
		getItems: func(before time.Time) ([]fluxer.MessagePin, error) {
			pins, err := s.GetChannelPins(channelID, before, limit, opts...)
			if err != nil {
				return nil, err
			}
			return pins.Items, nil
		},
		Time: start,
	}
}

func (s *channelImpl) PinMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(PinMessage.Compile(nil, channelID, messageID), nil, nil, opts...)
}

func (s *channelImpl) UnpinMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(UnpinMessage.Compile(nil, channelID, messageID), nil, nil, opts...)
}

func (s *channelImpl) Follow(channelID snowflake.ID, targetChannelID snowflake.ID, opts ...RequestOpt) (followedChannel *fluxer.FollowedChannel, err error) {
	err = s.client.Do(FollowChannel.Compile(nil, channelID), fluxer.FollowChannel{ChannelID: targetChannelID}, &followedChannel, opts...)
	return
}

type pollAnswerVotesResponse struct {
	Users []fluxer.User `json:"users"`
}
