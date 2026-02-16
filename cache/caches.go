package cache

import (
	"iter"
	"slices"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

type SelfUserCache interface {
	SelfUser() (fluxer.OAuth2User, bool)
	SetSelfUser(selfUser fluxer.OAuth2User)
}

func NewSelfUserCache() SelfUserCache {
	return &selfUserCacheImpl{}
}

type selfUserCacheImpl struct {
	selfUserMu sync.Mutex
	selfUser   *fluxer.OAuth2User
}

func (c *selfUserCacheImpl) SelfUser() (fluxer.OAuth2User, bool) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	if c.selfUser == nil {
		return fluxer.OAuth2User{}, false
	}
	return *c.selfUser, true
}

func (c *selfUserCacheImpl) SetSelfUser(user fluxer.OAuth2User) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	c.selfUser = &user
}

type GuildCache interface {
	GuildCache() Cache[fluxer.Guild]

	IsGuildUnready(guildID snowflake.ID) bool
	SetGuildUnready(guildID snowflake.ID, unready bool)
	UnreadyGuildIDs() []snowflake.ID

	IsGuildUnavailable(guildID snowflake.ID) bool
	SetGuildUnavailable(guildID snowflake.ID, unavailable bool)
	UnavailableGuildIDs() []snowflake.ID

	Guild(guildID snowflake.ID) (fluxer.Guild, bool)
	Guilds() iter.Seq[fluxer.Guild]
	GuildsLen() int
	AddGuild(guild fluxer.Guild)
	RemoveGuild(guildID snowflake.ID) (fluxer.Guild, bool)
}

func NewGuildCache(cache Cache[fluxer.Guild], unreadyGuilds Set[snowflake.ID], unavailableGuilds Set[snowflake.ID]) GuildCache {
	return &guildCacheImpl{
		cache:             cache,
		unreadyGuilds:     unreadyGuilds,
		unavailableGuilds: unavailableGuilds,
	}
}

type guildCacheImpl struct {
	cache             Cache[fluxer.Guild]
	unreadyGuilds     Set[snowflake.ID]
	unavailableGuilds Set[snowflake.ID]
}

func (c *guildCacheImpl) GuildCache() Cache[fluxer.Guild] {
	return c.cache
}

func (c *guildCacheImpl) IsGuildUnready(guildID snowflake.ID) bool {
	return c.unreadyGuilds.Has(guildID)
}

func (c *guildCacheImpl) SetGuildUnready(guildID snowflake.ID, unready bool) {
	if c.unreadyGuilds.Has(guildID) && !unready {
		c.unreadyGuilds.Remove(guildID)
	} else if !c.unreadyGuilds.Has(guildID) && unready {
		c.unreadyGuilds.Add(guildID)
	}
}

func (c *guildCacheImpl) UnreadyGuildIDs() []snowflake.ID {
	var guilds []snowflake.ID
	for guildID := range c.unreadyGuilds.All() {
		guilds = append(guilds, guildID)
	}
	return guilds
}

func (c *guildCacheImpl) IsGuildUnavailable(guildID snowflake.ID) bool {
	return c.unavailableGuilds.Has(guildID)
}

func (c *guildCacheImpl) SetGuildUnavailable(guildID snowflake.ID, unavailable bool) {
	if c.unavailableGuilds.Has(guildID) && !unavailable {
		c.unavailableGuilds.Remove(guildID)
	} else if !c.unavailableGuilds.Has(guildID) && unavailable {
		c.unavailableGuilds.Add(guildID)
	}
}

func (c *guildCacheImpl) UnavailableGuildIDs() []snowflake.ID {
	var guilds []snowflake.ID
	for guildId := range c.unavailableGuilds.All() {
		guilds = append(guilds, guildId)
	}
	return guilds
}

func (c *guildCacheImpl) Guild(guildID snowflake.ID) (fluxer.Guild, bool) {
	return c.cache.Get(guildID)
}

func (c *guildCacheImpl) Guilds() iter.Seq[fluxer.Guild] {
	return c.cache.All()
}

func (c *guildCacheImpl) GuildsLen() int {
	return c.cache.Len()
}

func (c *guildCacheImpl) AddGuild(guild fluxer.Guild) {
	c.cache.Put(guild.ID, guild)
}

func (c *guildCacheImpl) RemoveGuild(guildID snowflake.ID) (fluxer.Guild, bool) {
	return c.cache.Remove(guildID)
}

type ChannelCache interface {
	ChannelCache() Cache[fluxer.GuildChannel]

	Channel(channelID snowflake.ID) (fluxer.GuildChannel, bool)
	Channels() iter.Seq[fluxer.GuildChannel]
	ChannelsForGuild(guildID snowflake.ID) iter.Seq[fluxer.GuildChannel]
	ChannelsLen() int
	AddChannel(channel fluxer.GuildChannel)
	RemoveChannel(channelID snowflake.ID) (fluxer.GuildChannel, bool)
	RemoveChannelsByGuildID(guildID snowflake.ID)
}

func NewChannelCache(cache Cache[fluxer.GuildChannel]) ChannelCache {
	return &channelCacheImpl{
		cache: cache,
	}
}

type channelCacheImpl struct {
	cache Cache[fluxer.GuildChannel]
}

func (c *channelCacheImpl) ChannelCache() Cache[fluxer.GuildChannel] {
	return c.cache
}

func (c *channelCacheImpl) Channel(channelID snowflake.ID) (fluxer.GuildChannel, bool) {
	return c.cache.Get(channelID)
}

func (c *channelCacheImpl) Channels() iter.Seq[fluxer.GuildChannel] {
	return c.cache.All()
}

func (c *channelCacheImpl) ChannelsForGuild(guildID snowflake.ID) iter.Seq[fluxer.GuildChannel] {
	return func(yield func(fluxer.GuildChannel) bool) {
		for channel := range c.Channels() {
			if channel.GuildID() == guildID {
				if !yield(channel) {
					return
				}
			}
		}
	}
}

func (c *channelCacheImpl) ChannelsLen() int {
	return c.cache.Len()
}

func (c *channelCacheImpl) AddChannel(channel fluxer.GuildChannel) {
	c.cache.Put(channel.ID(), channel)
}

func (c *channelCacheImpl) RemoveChannel(channelID snowflake.ID) (fluxer.GuildChannel, bool) {
	return c.cache.Remove(channelID)
}

func (c *channelCacheImpl) RemoveChannelsByGuildID(guildID snowflake.ID) {
	c.cache.RemoveIf(func(channel fluxer.GuildChannel) bool {
		return channel.GuildID() == guildID
	})
}

type GuildScheduledEventCache interface {
	GuildScheduledEventCache() GroupedCache[fluxer.GuildScheduledEvent]

	GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (fluxer.GuildScheduledEvent, bool)
	GuildScheduledEvents(guildID snowflake.ID) iter.Seq[fluxer.GuildScheduledEvent]
	GuildScheduledEventsAllLen() int
	GuildScheduledEventsLen(guildID snowflake.ID) int
	AddGuildScheduledEvent(guildScheduledEvent fluxer.GuildScheduledEvent)
	RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (fluxer.GuildScheduledEvent, bool)
	RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID)
}

func NewGuildScheduledEventCache(cache GroupedCache[fluxer.GuildScheduledEvent]) GuildScheduledEventCache {
	return &guildScheduledEventCacheImpl{
		cache: cache,
	}
}

type guildScheduledEventCacheImpl struct {
	cache GroupedCache[fluxer.GuildScheduledEvent]
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventCache() GroupedCache[fluxer.GuildScheduledEvent] {
	return c.cache
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (fluxer.GuildScheduledEvent, bool) {
	return c.cache.Get(guildID, guildScheduledEventID)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvents(guildID snowflake.ID) iter.Seq[fluxer.GuildScheduledEvent] {
	return c.cache.GroupAll(guildID)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsAllLen() int {
	return c.cache.Len()
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *guildScheduledEventCacheImpl) AddGuildScheduledEvent(guildScheduledEvent fluxer.GuildScheduledEvent) {
	c.cache.Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (fluxer.GuildScheduledEvent, bool) {
	return c.cache.Remove(guildID, guildScheduledEventID)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type RoleCache interface {
	RoleCache() GroupedCache[fluxer.Role]

	Role(guildID snowflake.ID, roleID snowflake.ID) (fluxer.Role, bool)
	Roles(guildID snowflake.ID) iter.Seq[fluxer.Role]
	RolesAllLen() int
	RolesLen(guildID snowflake.ID) int
	AddRole(role fluxer.Role)
	RemoveRole(guildID snowflake.ID, roleID snowflake.ID) (fluxer.Role, bool)
	RemoveRolesByGuildID(guildID snowflake.ID)
}

func NewRoleCache(cache GroupedCache[fluxer.Role]) RoleCache {
	return &roleCacheImpl{
		cache: cache,
	}
}

type roleCacheImpl struct {
	cache GroupedCache[fluxer.Role]
}

func (c *roleCacheImpl) RoleCache() GroupedCache[fluxer.Role] {
	return c.cache
}

func (c *roleCacheImpl) Role(guildID snowflake.ID, roleID snowflake.ID) (fluxer.Role, bool) {
	return c.cache.Get(guildID, roleID)
}

func (c *roleCacheImpl) Roles(guildID snowflake.ID) iter.Seq[fluxer.Role] {
	return c.cache.GroupAll(guildID)
}

func (c *roleCacheImpl) RolesAllLen() int {
	return c.cache.Len()
}

func (c *roleCacheImpl) RolesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *roleCacheImpl) AddRole(role fluxer.Role) {
	c.cache.Put(role.GuildID, role.ID, role)
}

func (c *roleCacheImpl) RemoveRole(guildID snowflake.ID, roleID snowflake.ID) (fluxer.Role, bool) {
	return c.cache.Remove(guildID, roleID)
}

func (c *roleCacheImpl) RemoveRolesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type MemberCache interface {
	MemberCache() GroupedCache[fluxer.Member]

	Member(guildID snowflake.ID, userID snowflake.ID) (fluxer.Member, bool)
	Members(guildID snowflake.ID) iter.Seq[fluxer.Member]
	MembersAllLen() int
	MembersLen(guildID snowflake.ID) int
	AddMember(member fluxer.Member)
	RemoveMember(guildID snowflake.ID, userID snowflake.ID) (fluxer.Member, bool)
	RemoveMembersByGuildID(guildID snowflake.ID)
}

func NewMemberCache(cache GroupedCache[fluxer.Member]) MemberCache {
	return &memberCacheImpl{
		cache: cache,
	}
}

type memberCacheImpl struct {
	cache GroupedCache[fluxer.Member]
}

func (c *memberCacheImpl) MemberCache() GroupedCache[fluxer.Member] {
	return c.cache
}

func (c *memberCacheImpl) Member(guildID snowflake.ID, userID snowflake.ID) (fluxer.Member, bool) {
	return c.cache.Get(guildID, userID)
}

func (c *memberCacheImpl) Members(guildID snowflake.ID) iter.Seq[fluxer.Member] {
	return c.cache.GroupAll(guildID)
}

func (c *memberCacheImpl) MembersAllLen() int {
	return c.cache.Len()
}

func (c *memberCacheImpl) MembersLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *memberCacheImpl) AddMember(member fluxer.Member) {
	c.cache.Put(member.GuildID, member.User.ID, member)
}

func (c *memberCacheImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID) (fluxer.Member, bool) {
	return c.cache.Remove(guildID, userID)
}

func (c *memberCacheImpl) RemoveMembersByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type PresenceCache interface {
	PresenceCache() GroupedCache[fluxer.Presence]

	Presence(guildID snowflake.ID, userID snowflake.ID) (fluxer.Presence, bool)
	Presences(guildID snowflake.ID) iter.Seq[fluxer.Presence]
	PresencesAllLen() int
	PresencesLen(guildID snowflake.ID) int
	AddPresence(presence fluxer.Presence)
	RemovePresence(guildID snowflake.ID, userID snowflake.ID) (fluxer.Presence, bool)
	RemovePresencesByGuildID(guildID snowflake.ID)
}

func NewPresenceCache(cache GroupedCache[fluxer.Presence]) PresenceCache {
	return &presenceCacheImpl{
		cache: cache,
	}
}

type presenceCacheImpl struct {
	cache GroupedCache[fluxer.Presence]
}

func (c *presenceCacheImpl) PresenceCache() GroupedCache[fluxer.Presence] {
	return c.cache
}

func (c *presenceCacheImpl) Presence(guildID snowflake.ID, userID snowflake.ID) (fluxer.Presence, bool) {
	return c.cache.Get(guildID, userID)
}

func (c *presenceCacheImpl) Presences(guildID snowflake.ID) iter.Seq[fluxer.Presence] {
	return c.cache.GroupAll(guildID)
}

func (c *presenceCacheImpl) PresencesAllLen() int {
	return c.cache.Len()
}

func (c *presenceCacheImpl) PresencesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *presenceCacheImpl) AddPresence(presence fluxer.Presence) {
	c.cache.Put(presence.GuildID, presence.PresenceUser.ID, presence)
}

func (c *presenceCacheImpl) RemovePresence(guildID snowflake.ID, userID snowflake.ID) (fluxer.Presence, bool) {
	return c.cache.Remove(guildID, userID)
}

func (c *presenceCacheImpl) RemovePresencesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type VoiceStateCache interface {
	VoiceStateCache() GroupedCache[fluxer.VoiceState]

	VoiceState(guildID snowflake.ID, userID snowflake.ID) (fluxer.VoiceState, bool)
	VoiceStates(guildID snowflake.ID) iter.Seq[fluxer.VoiceState]
	VoiceStatesAllLen() int
	VoiceStatesLen(guildID snowflake.ID) int
	AddVoiceState(voiceState fluxer.VoiceState)
	RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID) (fluxer.VoiceState, bool)
	RemoveVoiceStatesByGuildID(guildID snowflake.ID)
}

func NewVoiceStateCache(cache GroupedCache[fluxer.VoiceState]) VoiceStateCache {
	return &voiceStateCacheImpl{
		cache: cache,
	}
}

type voiceStateCacheImpl struct {
	cache GroupedCache[fluxer.VoiceState]
}

func (c *voiceStateCacheImpl) VoiceStateCache() GroupedCache[fluxer.VoiceState] {
	return c.cache
}

func (c *voiceStateCacheImpl) VoiceState(guildID snowflake.ID, userID snowflake.ID) (fluxer.VoiceState, bool) {
	return c.cache.Get(guildID, userID)
}

func (c *voiceStateCacheImpl) VoiceStates(guildID snowflake.ID) iter.Seq[fluxer.VoiceState] {
	return c.cache.GroupAll(guildID)
}

func (c *voiceStateCacheImpl) VoiceStatesAllLen() int {
	return c.cache.Len()
}

func (c *voiceStateCacheImpl) VoiceStatesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *voiceStateCacheImpl) AddVoiceState(voiceState fluxer.VoiceState) {
	c.cache.Put(voiceState.GuildID, voiceState.UserID, voiceState)
}

func (c *voiceStateCacheImpl) RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID) (fluxer.VoiceState, bool) {
	return c.cache.Remove(guildID, userID)
}

func (c *voiceStateCacheImpl) RemoveVoiceStatesByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type MessageCache interface {
	MessageCache() GroupedCache[fluxer.Message]

	Message(channelID snowflake.ID, messageID snowflake.ID) (fluxer.Message, bool)
	Messages(channelID snowflake.ID) iter.Seq[fluxer.Message]
	MessagesAllLen() int
	MessagesLen(guildID snowflake.ID) int
	AddMessage(message fluxer.Message)
	RemoveMessage(channelID snowflake.ID, messageID snowflake.ID) (fluxer.Message, bool)
	RemoveMessagesByChannelID(channelID snowflake.ID)
	RemoveMessagesByGuildID(guildID snowflake.ID)
}

func NewMessageCache(cache GroupedCache[fluxer.Message]) MessageCache {
	return &messageCacheImpl{
		cache: cache,
	}
}

type messageCacheImpl struct {
	cache GroupedCache[fluxer.Message]
}

func (c *messageCacheImpl) MessageCache() GroupedCache[fluxer.Message] {
	return c.cache
}

func (c *messageCacheImpl) Message(channelID snowflake.ID, messageID snowflake.ID) (fluxer.Message, bool) {
	return c.cache.Get(channelID, messageID)
}

func (c *messageCacheImpl) Messages(channelID snowflake.ID) iter.Seq[fluxer.Message] {
	return c.cache.GroupAll(channelID)
}

func (c *messageCacheImpl) MessagesAllLen() int {
	return c.cache.Len()
}

func (c *messageCacheImpl) MessagesLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *messageCacheImpl) AddMessage(message fluxer.Message) {
	c.cache.Put(message.ChannelID, message.ID, message)
}

func (c *messageCacheImpl) RemoveMessage(channelID snowflake.ID, messageID snowflake.ID) (fluxer.Message, bool) {
	return c.cache.Remove(channelID, messageID)
}

func (c *messageCacheImpl) RemoveMessagesByChannelID(channelID snowflake.ID) {
	c.cache.GroupRemove(channelID)
}

func (c *messageCacheImpl) RemoveMessagesByGuildID(guildID snowflake.ID) {
	c.cache.RemoveIf(func(_ snowflake.ID, message fluxer.Message) bool {
		return message.GuildID != nil && *message.GuildID == guildID
	})
}

type EmojiCache interface {
	EmojiCache() GroupedCache[fluxer.Emoji]

	Emoji(guildID snowflake.ID, emojiID snowflake.ID) (fluxer.Emoji, bool)
	Emojis(guildID snowflake.ID) iter.Seq[fluxer.Emoji]
	EmojisAllLen() int
	EmojisLen(guildID snowflake.ID) int
	AddEmoji(emoji fluxer.Emoji)
	RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID) (fluxer.Emoji, bool)
	RemoveEmojisByGuildID(guildID snowflake.ID)
}

func NewEmojiCache(cache GroupedCache[fluxer.Emoji]) EmojiCache {
	return &emojiCacheImpl{
		cache: cache,
	}
}

type emojiCacheImpl struct {
	cache GroupedCache[fluxer.Emoji]
}

func (c *emojiCacheImpl) EmojiCache() GroupedCache[fluxer.Emoji] {
	return c.cache
}

func (c *emojiCacheImpl) Emoji(guildID snowflake.ID, emojiID snowflake.ID) (fluxer.Emoji, bool) {
	return c.cache.Get(guildID, emojiID)
}

func (c *emojiCacheImpl) Emojis(guildID snowflake.ID) iter.Seq[fluxer.Emoji] {
	return c.cache.GroupAll(guildID)
}

func (c *emojiCacheImpl) EmojisAllLen() int {
	return c.cache.Len()
}

func (c *emojiCacheImpl) EmojisLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *emojiCacheImpl) AddEmoji(emoji fluxer.Emoji) {
	c.cache.Put(emoji.GuildID, emoji.ID, emoji)
}

func (c *emojiCacheImpl) RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID) (fluxer.Emoji, bool) {
	return c.cache.Remove(guildID, emojiID)
}

func (c *emojiCacheImpl) RemoveEmojisByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

type StickerCache interface {
	StickerCache() GroupedCache[fluxer.Sticker]

	Sticker(guildID snowflake.ID, stickerID snowflake.ID) (fluxer.Sticker, bool)
	Stickers(guildID snowflake.ID) iter.Seq[fluxer.Sticker]
	StickersAllLen() int
	StickersLen(guildID snowflake.ID) int
	AddSticker(sticker fluxer.Sticker)
	RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID) (fluxer.Sticker, bool)
	RemoveStickersByGuildID(guildID snowflake.ID)
}

func NewStickerCache(cache GroupedCache[fluxer.Sticker]) StickerCache {
	return &stickerCacheImpl{
		cache: cache,
	}
}

type stickerCacheImpl struct {
	cache GroupedCache[fluxer.Sticker]
}

func (c *stickerCacheImpl) StickerCache() GroupedCache[fluxer.Sticker] {
	return c.cache
}

func (c *stickerCacheImpl) Sticker(guildID snowflake.ID, stickerID snowflake.ID) (fluxer.Sticker, bool) {
	return c.cache.Get(guildID, stickerID)
}

func (c *stickerCacheImpl) Stickers(guildID snowflake.ID) iter.Seq[fluxer.Sticker] {
	return c.cache.GroupAll(guildID)
}

func (c *stickerCacheImpl) StickersAllLen() int {
	return c.cache.Len()
}

func (c *stickerCacheImpl) StickersLen(guildID snowflake.ID) int {
	return c.cache.GroupLen(guildID)
}

func (c *stickerCacheImpl) AddSticker(sticker fluxer.Sticker) {
	if sticker.GuildID == nil {
		return
	}
	c.cache.Put(*sticker.GuildID, sticker.ID, sticker)
}

func (c *stickerCacheImpl) RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID) (fluxer.Sticker, bool) {
	return c.cache.Remove(guildID, stickerID)
}

func (c *stickerCacheImpl) RemoveStickersByGuildID(guildID snowflake.ID) {
	c.cache.GroupRemove(guildID)
}

// Caches combines all different entity caches into one with some utility methods.
type Caches interface {
	SelfUserCache
	GuildCache
	ChannelCache
	GuildScheduledEventCache
	RoleCache
	MemberCache
	PresenceCache
	VoiceStateCache
	MessageCache
	EmojiCache
	StickerCache

	// CacheFlags returns the current configured FLags of the caches.
	CacheFlags() Flags

	// MemberPermissions returns the calculated permissions of the given member.
	// This requires the FlagRoles to be set.
	MemberPermissions(member fluxer.Member) fluxer.Permissions

	// MemberPermissionsInChannel returns the calculated permissions of the given member in the given channel.
	// This requires the FlagRoles and FlagChannels to be set.
	MemberPermissionsInChannel(channel fluxer.GuildChannel, member fluxer.Member) fluxer.Permissions

	// MemberRoles returns all roles of the given member.
	// This requires the FlagRoles to be set.
	MemberRoles(member fluxer.Member) []fluxer.Role

	// AudioChannelMembers returns all members which are in the given audio channel.
	// This requires the FlagVoiceStates to be set.
	AudioChannelMembers(channel fluxer.GuildAudioChannel) []fluxer.Member

	// SelfMember returns the current bot member from the given guildID.
	// This is only available after we received the gateway.EventTypeGuildCreate event for the given guildID.
	SelfMember(guildID snowflake.ID) (fluxer.Member, bool)

	// GuildMessageChannel returns a fluxer.GuildMessageChannel from the ChannelCache and a bool indicating if it exists.
	GuildMessageChannel(channelID snowflake.ID) (fluxer.GuildMessageChannel, bool)

	// GuildAudioChannel returns a fluxer.GetGuildAudioChannel from the ChannelCache and a bool indicating if it exists.
	GuildAudioChannel(channelID snowflake.ID) (fluxer.GuildAudioChannel, bool)

	// GuildTextChannel returns a fluxer.GuildTextChannel from the ChannelCache and a bool indicating if it exists.
	GuildTextChannel(channelID snowflake.ID) (fluxer.GuildTextChannel, bool)

	// GuildVoiceChannel returns a fluxer.GuildVoiceChannel from the ChannelCache and a bool indicating if it exists.
	GuildVoiceChannel(channelID snowflake.ID) (fluxer.GuildVoiceChannel, bool)

	// GuildCategoryChannel returns a fluxer.GuildCategoryChannel from the ChannelCache and a bool indicating if it exists.
	GuildCategoryChannel(channelID snowflake.ID) (fluxer.GuildCategoryChannel, bool)
}

// New returns a new default Caches instance with the given ConfigOpt(s) applied.
func New(opts ...ConfigOpt) Caches {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &cachesImpl{
		config:                   cfg,
		selfUserCache:            cfg.SelfUserCache,
		guildCache:               cfg.GuildCache,
		channelCache:             cfg.ChannelCache,
		guildScheduledEventCache: cfg.GuildScheduledEventCache,
		roleCache:                cfg.RoleCache,
		memberCache:              cfg.MemberCache,
		presenceCache:            cfg.PresenceCache,
		voiceStateCache:          cfg.VoiceStateCache,
		messageCache:             cfg.MessageCache,
		emojiCache:               cfg.EmojiCache,
		stickerCache:             cfg.StickerCache,
	}
}

// these type aliases are needed to allow having the GuildCache, ChannelCache, etc. as methods on the cachesImpl struct
type (
	guildCache               = GuildCache
	channelCache             = ChannelCache
	guildScheduledEventCache = GuildScheduledEventCache
	roleCache                = RoleCache
	memberCache              = MemberCache
	presenceCache            = PresenceCache
	voiceStateCache          = VoiceStateCache
	messageCache             = MessageCache
	emojiCache               = EmojiCache
	stickerCache             = StickerCache
	selfUserCache            = SelfUserCache
)

type cachesImpl struct {
	config config

	guildCache
	channelCache
	guildScheduledEventCache
	roleCache
	memberCache
	presenceCache
	voiceStateCache
	messageCache
	emojiCache
	stickerCache
	selfUserCache
}

func (c *cachesImpl) CacheFlags() Flags {
	return c.config.CacheFlags
}

func (c *cachesImpl) MemberPermissions(member fluxer.Member) fluxer.Permissions {
	if guild, ok := c.Guild(member.GuildID); ok && guild.OwnerID == member.User.ID {
		return fluxer.PermissionsAll
	}

	var permissions fluxer.Permissions
	if publicRole, ok := c.Role(member.GuildID, member.GuildID); ok {
		permissions = publicRole.Permissions
	}

	for _, role := range c.MemberRoles(member) {
		permissions = permissions.Add(role.Permissions)
		if permissions.Has(fluxer.PermissionAdministrator) {
			return fluxer.PermissionsAll
		}
	}
	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= fluxer.PermissionViewChannel | fluxer.PermissionReadMessageHistory
	}
	return permissions
}

func (c *cachesImpl) MemberPermissionsInChannel(channel fluxer.GuildChannel, member fluxer.Member) fluxer.Permissions {
	permissions := c.MemberPermissions(member)
	if permissions.Has(fluxer.PermissionAdministrator) {
		return fluxer.PermissionsAll
	}

	var (
		allow fluxer.Permissions
		deny  fluxer.Permissions
	)

	if overwrite, ok := channel.PermissionOverwrites().Role(channel.GuildID()); ok {
		permissions |= overwrite.Allow
		permissions &= ^overwrite.Deny
	}

	for _, roleID := range member.RoleIDs {
		if roleID == channel.GuildID() {
			continue
		}

		if overwrite, ok := channel.PermissionOverwrites().Role(roleID); ok {
			allow |= overwrite.Allow
			deny |= overwrite.Deny
		}
	}

	if overwrite, ok := channel.PermissionOverwrites().Member(member.User.ID); ok {
		allow |= overwrite.Allow
		deny |= overwrite.Deny
	}

	permissions &= ^deny
	permissions |= allow

	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= fluxer.PermissionViewChannel | fluxer.PermissionReadMessageHistory
	}

	return permissions
}

func (c *cachesImpl) MemberRoles(member fluxer.Member) []fluxer.Role {
	var roles []fluxer.Role

	for role := range c.Roles(member.GuildID) {
		if slices.Contains(member.RoleIDs, role.ID) {
			roles = append(roles, role)
		}
	}
	return roles
}

func (c *cachesImpl) AudioChannelMembers(channel fluxer.GuildAudioChannel) []fluxer.Member {
	var members []fluxer.Member
	for state := range c.VoiceStates(channel.GuildID()) {
		if member, ok := c.Member(channel.GuildID(), state.UserID); ok && state.ChannelID != nil && *state.ChannelID == channel.ID() {
			members = append(members, member)
		}
	}
	return members
}

func (c *cachesImpl) SelfMember(guildID snowflake.ID) (fluxer.Member, bool) {
	selfUser, ok := c.SelfUser()
	if !ok {
		return fluxer.Member{}, false
	}
	return c.Member(guildID, selfUser.ID)
}

func (c *cachesImpl) GuildMessageChannel(channelID snowflake.ID) (fluxer.GuildMessageChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if chM, ok := ch.(fluxer.GuildMessageChannel); ok {
			return chM, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildAudioChannel(channelID snowflake.ID) (fluxer.GuildAudioChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(fluxer.GuildAudioChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildTextChannel(channelID snowflake.ID) (fluxer.GuildTextChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(fluxer.GuildTextChannel); ok {
			return cCh, true
		}
	}
	return fluxer.GuildTextChannel{}, false
}

func (c *cachesImpl) GuildVoiceChannel(channelID snowflake.ID) (fluxer.GuildVoiceChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(fluxer.GuildVoiceChannel); ok {
			return cCh, true
		}
	}
	return fluxer.GuildVoiceChannel{}, false
}

func (c *cachesImpl) GuildCategoryChannel(channelID snowflake.ID) (fluxer.GuildCategoryChannel, bool) {
	if ch, ok := c.Channel(channelID); ok {
		if cCh, ok := ch.(fluxer.GuildCategoryChannel); ok {
			return cCh, true
		}
	}
	return fluxer.GuildCategoryChannel{}, false
}
