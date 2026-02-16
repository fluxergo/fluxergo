package handlers

import (
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/gateway"
)

// GetGatewayHandlers returns the default gateway.Gateway event handlers for processing the raw payload which gets passed into the bot.EventManager
func GetGatewayHandlers() map[gateway.EventType]bot.GatewayEventHandler {
	handlers := make(map[gateway.EventType]bot.GatewayEventHandler, len(allEventHandlers))
	for _, handler := range allEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var allEventHandlers = []bot.GatewayEventHandler{
	bot.NewGatewayEventHandler(gateway.EventTypeRaw, gatewayHandlerRaw),
	bot.NewGatewayEventHandler(gateway.EventTypeHeartbeatAck, gatewayHandlerHeartbeatAck),
	bot.NewGatewayEventHandler(gateway.EventTypeReady, gatewayHandlerReady),
	bot.NewGatewayEventHandler(gateway.EventTypeResumed, gatewayHandlerResumed),

	bot.NewGatewayEventHandler(gateway.EventTypeChannelCreate, gatewayHandlerChannelCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelUpdate, gatewayHandlerChannelUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelDelete, gatewayHandlerChannelDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelPinsUpdate, gatewayHandlerChannelPinsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildCreate, gatewayHandlerGuildCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildUpdate, gatewayHandlerGuildUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildDelete, gatewayHandlerGuildDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildBanAdd, gatewayHandlerGuildBanAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildBanRemove, gatewayHandlerGuildBanRemove),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildEmojisUpdate, gatewayHandlerGuildEmojisUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildStickersUpdate, gatewayHandlerGuildStickersUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildIntegrationsUpdate, gatewayHandlerGuildIntegrationsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberAdd, gatewayHandlerGuildMemberAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberRemove, gatewayHandlerGuildMemberRemove),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberUpdate, gatewayHandlerGuildMemberUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleCreate, gatewayHandlerGuildRoleCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleUpdate, gatewayHandlerGuildRoleUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleDelete, gatewayHandlerGuildRoleDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventCreate, gatewayHandlerGuildScheduledEventCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUpdate, gatewayHandlerGuildScheduledEventUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventDelete, gatewayHandlerGuildScheduledEventDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeInviteCreate, gatewayHandlerInviteCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeInviteDelete, gatewayHandlerInviteDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeMessageCreate, gatewayHandlerMessageCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageUpdate, gatewayHandlerMessageUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageDelete, gatewayHandlerMessageDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageDeleteBulk, gatewayHandlerMessageDeleteBulk),

	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionAdd, gatewayHandlerMessageReactionAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionRemove, gatewayHandlerMessageReactionRemove),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionRemoveAll, gatewayHandlerMessageReactionRemoveAll),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionRemoveEmoji, gatewayHandlerMessageReactionRemoveEmoji),

	bot.NewGatewayEventHandler(gateway.EventTypePresenceUpdate, gatewayHandlerPresenceUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeTypingStart, gatewayHandlerTypingStart),
	bot.NewGatewayEventHandler(gateway.EventTypeUserUpdate, gatewayHandlerUserUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeVoiceStateUpdate, gatewayHandlerVoiceStateUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeVoiceServerUpdate, gatewayHandlerVoiceServerUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeWebhooksUpdate, gatewayHandlerWebhooksUpdate),
}
