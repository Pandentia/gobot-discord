package statekeeping

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mediocregopher/radix/v3"
)

// Define our marshalers.
var (
	Marshal   = json.Marshal
	Unmarshal = json.Unmarshal
)

var _ State = &RedisState{}

// RedisState is a State implementation using Redis as a backend.
type RedisState struct {
	Redis *radix.Pool
}

func marshal(v interface{}) []byte {
	data, err := Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func unmarshal(data []byte, v interface{}) {
	err := Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

// formatters
const fmtMsg = "msg:%s:%s"
const fmtGuild = "guild:%s"
const fmtChannel = "channel:%s:%s"
const fmtRole = "role:%s:%s"
const fmtEmoji = "emoji:%s:%s"
const fmtUser = "user:%s"
const fmtPresence = "presence:%s:%s"
const fmtMember = "member:%s:%s"

// private handlers

func (r *RedisState) handleReady(ready *discordgo.Ready) {
	for _, guild := range ready.Guilds {
		r.handleGuildUpdate(guild)
	}
}

func (r *RedisState) handleUserUpdate(user *discordgo.User) {
	key := fmt.Sprintf(fmtUser, user.ID)
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(user)))
}

func (r *RedisState) handlePresenceUpdate(presenceUpdate *discordgo.PresenceUpdate) {
	presence := presenceUpdate.Presence
	key := fmt.Sprintf(fmtPresence, presenceUpdate.GuildID, presence.User.ID)
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(presence)))
}

func (r *RedisState) handleGuildUpdate(guild *discordgo.Guild) {
	key := fmt.Sprintf(fmtGuild, guild.ID)
	// consume channels
	channels := guild.Channels
	guild.Channels = nil
	// consume roles
	roles := guild.Roles
	guild.Roles = nil
	// consume emojis
	emojis := guild.Emojis
	guild.Emojis = nil
	// consume members
	members := guild.Members
	guild.Members = nil
	// consume presences
	presences := guild.Presences
	guild.Presences = nil
	// set guild
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(guild)))
	// set channels
	for _, channel := range channels {
		r.handleChannelUpdate(channel)
	}
	// set roles
	for _, role := range roles {
		r.handleGuildRoleUpdate(&discordgo.GuildRole{
			Role:    role,
			GuildID: guild.ID,
		})
	}
	// set emojis
	r.handleGuildEmojisUpdate(&discordgo.GuildEmojisUpdate{
		Emojis:  emojis,
		GuildID: guild.ID,
	})
	// set members
	for _, member := range members {
		r.handleGuildMemberUpdate(member)
	}
	// set presences
	for _, presence := range presences {
		r.handlePresenceUpdate(&discordgo.PresenceUpdate{
			Presence: *presence,
			GuildID:  guild.ID,
		})
	}
}

func (r *RedisState) handleGuildDelete(guild *discordgo.Guild) {
	key := fmt.Sprintf(fmtGuild, guild.ID)
	r.Redis.Do(radix.Cmd(nil, "DEL", key))

	// recursive delete
	keys := make([]string, 0)
	// delete associated channels
	r.Redis.Do(radix.Cmd(&keys, "KEYS", fmt.Sprintf(fmtChannel, guild.ID, "*")))
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
	// delete associated roles
	r.Redis.Do(radix.Cmd(&keys, "KEYS", fmt.Sprintf(fmtRole, guild.ID, "*")))
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
	// delete associated emoji
	r.Redis.Do(radix.Cmd(&keys, "KEYS", fmt.Sprintf(fmtEmoji, guild.ID, "*")))
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
	// delete associated presences
	r.Redis.Do(radix.Cmd(&keys, "KEYS", fmt.Sprintf(fmtPresence, guild.ID, "*")))
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
	// delete associated members
	r.Redis.Do(radix.Cmd(&keys, "KEYS", fmt.Sprintf(fmtMember, guild.ID, "*")))
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
}

func (r *RedisState) handleGuildMemberUpdate(member *discordgo.Member) {
	key := fmt.Sprintf(fmtMember, member.GuildID, member.User.ID)
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(member)))
	r.handleUserUpdate(member.User)
}

func (r *RedisState) handleGuildMemberUpdateChunked(chunk *discordgo.GuildMembersChunk) {
	for _, member := range chunk.Members {
		r.handleGuildMemberUpdate(member)
	}
}

func (r *RedisState) handleGuildMemberDelete(member *discordgo.Member) {
	key := fmt.Sprintf(fmtMember, member.GuildID, member.User.ID)
	r.Redis.Do(radix.Cmd(nil, "DEL", key))
}

func (r *RedisState) handleGuildRoleUpdate(role *discordgo.GuildRole) {
	key := fmt.Sprintf(fmtRole, role.GuildID, role.Role.ID)
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(role.Role)))
}

func (r *RedisState) handleGuildRoleDelete(role *discordgo.GuildRole) {
	key := fmt.Sprintf(fmtRole, role.GuildID, role.Role.ID)
	r.Redis.Do(radix.Cmd(nil, "DEL", key))
}

func (r *RedisState) handleGuildEmojisUpdate(event *discordgo.GuildEmojisUpdate) {
	for _, emoji := range event.Emojis {
		key := fmt.Sprintf(fmtEmoji, event.GuildID, emoji.ID)
		r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(emoji)))
	}
}

func (r *RedisState) handleChannelUpdate(channel *discordgo.Channel) {
	key := fmt.Sprintf(fmtChannel, channel.GuildID, channel.ID)
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(channel)))
}

func (r *RedisState) handleChannelDelete(channel *discordgo.Channel) {
	key := fmt.Sprintf(fmtChannel, channel.GuildID, channel.ID)
	r.Redis.Do(radix.Cmd(nil, "DEL", key))
	// recursive delete
	keys := make([]string, 0)
	r.Redis.Do(radix.Cmd(&keys, "KEYS", fmt.Sprintf(fmtMsg, channel.ID, "*")))
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
}

func (r *RedisState) handleMessageUpdate(message *discordgo.Message) {
	key := fmt.Sprintf(fmtChannel, message.ChannelID, message.ID)
	r.Redis.Do(radix.FlatCmd(nil, "SET", key, marshal(message)))
	r.Redis.Do(radix.Cmd(nil, "EXPIRES", key, "3600"))
}

func (r *RedisState) handleMessageDelete(message *discordgo.Message) {
	key := fmt.Sprintf(fmtChannel, message.ChannelID, message.ID)
	r.Redis.Do(radix.Cmd(nil, "DEL", key))
}

func (r *RedisState) handleMessageDeleteBulk(bulk *discordgo.MessageDeleteBulk) {
	keys := make([]string, 0)
	for _, msg := range bulk.Messages {
		keys = append(keys, fmt.Sprintf(fmtMsg, bulk.ChannelID, msg))
	}
	r.Redis.Do(radix.Cmd(nil, "DEL", keys...))
}

// OnEvent processes all events coming from Discord.
func (r *RedisState) OnEvent(session *discordgo.Session, event interface{}) {
	switch typedEvent := event.(type) {
	case *discordgo.Event:
		r.Redis.Do(radix.Cmd(nil, "INCR", "stats:events"))
	case *discordgo.Ready:
		r.handleReady(typedEvent)

	case *discordgo.GuildCreate:
		r.handleGuildUpdate(typedEvent.Guild)
	case *discordgo.GuildUpdate:
		r.handleGuildUpdate(typedEvent.Guild)
	case *discordgo.GuildDelete:
		r.handleGuildDelete(typedEvent.Guild)

	case *discordgo.ChannelCreate:
		r.handleChannelUpdate(typedEvent.Channel)
	case *discordgo.ChannelUpdate:
		r.handleChannelUpdate(typedEvent.Channel)
	case *discordgo.ChannelDelete:
		r.handleChannelDelete(typedEvent.Channel)

	case *discordgo.MessageCreate:
		r.handleMessageUpdate(typedEvent.Message)
	case *discordgo.MessageUpdate:
		r.handleMessageUpdate(typedEvent.Message)
	case *discordgo.MessageDelete:
		r.handleMessageDelete(typedEvent.Message)
	case *discordgo.MessageDeleteBulk:
		r.handleMessageDeleteBulk(typedEvent)

	case *discordgo.GuildEmojisUpdate:
		r.handleGuildEmojisUpdate(typedEvent)

	case *discordgo.UserUpdate:
		r.handleUserUpdate(typedEvent.User)
	case *discordgo.PresenceUpdate:
		r.handlePresenceUpdate(typedEvent)
	case *discordgo.GuildMemberAdd:
		r.handleGuildMemberUpdate(typedEvent.Member)
	case *discordgo.GuildMemberUpdate:
		r.handleGuildMemberUpdate(typedEvent.Member)
	case *discordgo.GuildMemberRemove:
		r.handleGuildMemberDelete(typedEvent.Member)

	case *discordgo.Connect:
		r.Redis.Do(radix.Cmd(nil, "INCR", "stats:connects"))

	case *discordgo.TypingStart, *discordgo.Disconnect:
		// ignore

	default:
		fmt.Printf("State: Unrecognized event: %T\n", event)
	}
}

// Public functions

// Channels gets channels by Guild ID.
func (r *RedisState) Channels(guildID string) []*discordgo.Channel {
	return nil
}

// Guilds gets a slice of Guilds the bot is a part of.
func (r *RedisState) Guilds() []*discordgo.Guild {
	return nil
}

// Members gets a slice of Members in a given guild.
func (r *RedisState) Members(guildID string) []*discordgo.Member {
	return nil
}

func (r RedisState) Messages(channelID string) []*discordgo.Message {
	return nil
}

func (r *RedisState) User(userID string) *discordgo.User {
	return nil
}

func (r *RedisState) Presence(userID string) *discordgo.Presence {
	return nil
}

func (r *RedisState) Member(guildID, userID string) *discordgo.Member {
	return nil
}

func (r *RedisState) Message(channelID, messageID string) *discordgo.Message {
	return nil
}

func (r *RedisState) Role(guildID, roleID string) *discordgo.Role {
	return nil
}

// Events returns the number of events processed by RedisState.
func (r *RedisState) Events() string {
	events := ""
	r.Redis.Do(radix.Cmd(&events, "GET", "stats:events"))
	return fmt.Sprintf("%s events", events)
}

// Size returns the size of the data set.
func (r *RedisState) Size() string {
	size := ""
	r.Redis.Do(radix.Cmd(&size, "DBSIZE"))
	return fmt.Sprintf("%s keys", size)
}
