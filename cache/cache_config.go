package cache

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

func defaultConfig() config {
	return config{
		GuildCachePolicy:               PolicyAll[fluxer.Guild],
		ChannelCachePolicy:             PolicyAll[fluxer.GuildChannel],
		GuildScheduledEventCachePolicy: PolicyAll[fluxer.GuildScheduledEvent],
		RoleCachePolicy:                PolicyAll[fluxer.Role],
		MemberCachePolicy:              PolicyAll[fluxer.Member],
		PresenceCachePolicy:            PolicyAll[fluxer.Presence],
		VoiceStateCachePolicy:          PolicyAll[fluxer.VoiceState],
		MessageCachePolicy:             PolicyAll[fluxer.Message],
		EmojiCachePolicy:               PolicyAll[fluxer.Emoji],
		StickerCachePolicy:             PolicyAll[fluxer.Sticker],
	}
}

type config struct {
	CacheFlags Flags

	SelfUserCache SelfUserCache

	GuildCache       GuildCache
	GuildCachePolicy Policy[fluxer.Guild]

	ChannelCache       ChannelCache
	ChannelCachePolicy Policy[fluxer.GuildChannel]

	GuildScheduledEventCache       GuildScheduledEventCache
	GuildScheduledEventCachePolicy Policy[fluxer.GuildScheduledEvent]

	RoleCache       RoleCache
	RoleCachePolicy Policy[fluxer.Role]

	MemberCache       MemberCache
	MemberCachePolicy Policy[fluxer.Member]

	PresenceCache       PresenceCache
	PresenceCachePolicy Policy[fluxer.Presence]

	VoiceStateCache       VoiceStateCache
	VoiceStateCachePolicy Policy[fluxer.VoiceState]

	MessageCache       MessageCache
	MessageCachePolicy Policy[fluxer.Message]

	EmojiCache       EmojiCache
	EmojiCachePolicy Policy[fluxer.Emoji]

	StickerCache       StickerCache
	StickerCachePolicy Policy[fluxer.Sticker]
}

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Caches.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.SelfUserCache == nil {
		c.SelfUserCache = NewSelfUserCache()
	}
	if c.GuildCache == nil {
		c.GuildCache = NewGuildCache(NewCache[fluxer.Guild](c.CacheFlags, FlagGuilds, c.GuildCachePolicy), NewSet[snowflake.ID](), NewSet[snowflake.ID]())
	}
	if c.ChannelCache == nil {
		c.ChannelCache = NewChannelCache(NewCache[fluxer.GuildChannel](c.CacheFlags, FlagChannels, c.ChannelCachePolicy))
	}
	if c.GuildScheduledEventCache == nil {
		c.GuildScheduledEventCache = NewGuildScheduledEventCache(NewGroupedCache[fluxer.GuildScheduledEvent](c.CacheFlags, FlagGuildScheduledEvents, c.GuildScheduledEventCachePolicy))
	}
	if c.RoleCache == nil {
		c.RoleCache = NewRoleCache(NewGroupedCache[fluxer.Role](c.CacheFlags, FlagRoles, c.RoleCachePolicy))
	}
	if c.MemberCache == nil {
		c.MemberCache = NewMemberCache(NewGroupedCache[fluxer.Member](c.CacheFlags, FlagMembers, c.MemberCachePolicy))
	}
	if c.PresenceCache == nil {
		c.PresenceCache = NewPresenceCache(NewGroupedCache[fluxer.Presence](c.CacheFlags, FlagPresences, c.PresenceCachePolicy))
	}
	if c.VoiceStateCache == nil {
		c.VoiceStateCache = NewVoiceStateCache(NewGroupedCache[fluxer.VoiceState](c.CacheFlags, FlagVoiceStates, c.VoiceStateCachePolicy))
	}
	if c.MessageCache == nil {
		c.MessageCache = NewMessageCache(NewGroupedCache[fluxer.Message](c.CacheFlags, FlagMessages, c.MessageCachePolicy))
	}
	if c.EmojiCache == nil {
		c.EmojiCache = NewEmojiCache(NewGroupedCache[fluxer.Emoji](c.CacheFlags, FlagEmojis, c.EmojiCachePolicy))
	}
	if c.StickerCache == nil {
		c.StickerCache = NewStickerCache(NewGroupedCache[fluxer.Sticker](c.CacheFlags, FlagStickers, c.StickerCachePolicy))
	}
}

// WithCaches sets the Flags of the config.
func WithCaches(flags ...Flags) ConfigOpt {
	return func(config *config) {
		config.CacheFlags = config.CacheFlags.Add(flags...)
	}
}

// WithSelfUserCache sets the SelfUserCache of the config.
func WithSelfUserCache(cache SelfUserCache) ConfigOpt {
	return func(config *config) {
		config.SelfUserCache = cache
	}
}

// WithGuildCachePolicy sets the Policy[fluxer.Guild] of the config.
func WithGuildCachePolicy(policy Policy[fluxer.Guild]) ConfigOpt {
	return func(config *config) {
		config.GuildCachePolicy = policy
	}
}

// WithGuildCache sets the GuildCache of the config.
func WithGuildCache(guildCache GuildCache) ConfigOpt {
	return func(config *config) {
		config.GuildCache = guildCache
	}
}

// WithChannelCachePolicy sets the Policy[fluxer.Channel] of the config.
func WithChannelCachePolicy(policy Policy[fluxer.GuildChannel]) ConfigOpt {
	return func(config *config) {
		config.ChannelCachePolicy = policy
	}
}

// WithChannelCache sets the ChannelCache of the config.
func WithChannelCache(channelCache ChannelCache) ConfigOpt {
	return func(config *config) {
		config.ChannelCache = channelCache
	}
}

// WithGuildScheduledEventCachePolicy sets the Policy[fluxer.GuildScheduledEvent] of the config.
func WithGuildScheduledEventCachePolicy(policy Policy[fluxer.GuildScheduledEvent]) ConfigOpt {
	return func(config *config) {
		config.GuildScheduledEventCachePolicy = policy
	}
}

// WithGuildScheduledEventCache sets the GuildScheduledEventCache of the config.
func WithGuildScheduledEventCache(guildScheduledEventCache GuildScheduledEventCache) ConfigOpt {
	return func(config *config) {
		config.GuildScheduledEventCache = guildScheduledEventCache
	}
}

// WithRoleCachePolicy sets the Policy[fluxer.Role] of the config.
func WithRoleCachePolicy(policy Policy[fluxer.Role]) ConfigOpt {
	return func(config *config) {
		config.RoleCachePolicy = policy
	}
}

// WithRoleCache sets the RoleCache of the config.
func WithRoleCache(roleCache RoleCache) ConfigOpt {
	return func(config *config) {
		config.RoleCache = roleCache
	}
}

// WithMemberCachePolicy sets the Policy[fluxer.Member] of the config.
func WithMemberCachePolicy(policy Policy[fluxer.Member]) ConfigOpt {
	return func(config *config) {
		config.MemberCachePolicy = policy
	}
}

// WithMemberCache sets the MemberCache of the config.
func WithMemberCache(memberCache MemberCache) ConfigOpt {
	return func(config *config) {
		config.MemberCache = memberCache
	}
}

// WithPresenceCachePolicy sets the Policy[fluxer.Presence] of the config.
func WithPresenceCachePolicy(policy Policy[fluxer.Presence]) ConfigOpt {
	return func(config *config) {
		config.PresenceCachePolicy = policy
	}
}

// WithPresenceCache sets the PresenceCache of the config.
func WithPresenceCache(presenceCache PresenceCache) ConfigOpt {
	return func(config *config) {
		config.PresenceCache = presenceCache
	}
}

// WithVoiceStateCachePolicy sets the Policy[fluxer.VoiceState] of the config.
func WithVoiceStateCachePolicy(policy Policy[fluxer.VoiceState]) ConfigOpt {
	return func(config *config) {
		config.VoiceStateCachePolicy = policy
	}
}

// WithVoiceStateCache sets the VoiceStateCache of the config.
func WithVoiceStateCache(voiceStateCache VoiceStateCache) ConfigOpt {
	return func(config *config) {
		config.VoiceStateCache = voiceStateCache
	}
}

// WithMessageCachePolicy sets the Policy[fluxer.Message] of the config.
func WithMessageCachePolicy(policy Policy[fluxer.Message]) ConfigOpt {
	return func(config *config) {
		config.MessageCachePolicy = policy
	}
}

// WithMessageCache sets the MessageCache of the config.
func WithMessageCache(messageCache MessageCache) ConfigOpt {
	return func(config *config) {
		config.MessageCache = messageCache
	}
}

// WithEmojiCachePolicy sets the Policy[fluxer.Emoji] of the config.
func WithEmojiCachePolicy(policy Policy[fluxer.Emoji]) ConfigOpt {
	return func(config *config) {
		config.EmojiCachePolicy = policy
	}
}

// WithEmojiCache sets the EmojiCache of the config.
func WithEmojiCache(emojiCache EmojiCache) ConfigOpt {
	return func(config *config) {
		config.EmojiCache = emojiCache
	}
}

// WithStickerCachePolicy sets the Policy[fluxer.Sticker] of the config.
func WithStickerCachePolicy(policy Policy[fluxer.Sticker]) ConfigOpt {
	return func(config *config) {
		config.StickerCachePolicy = policy
	}
}

// WithStickerCache sets the StickerCache of the config.
func WithStickerCache(stickerCache StickerCache) ConfigOpt {
	return func(config *config) {
		config.StickerCache = stickerCache
	}
}
