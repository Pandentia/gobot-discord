// Package statekeeping implements state-keeping for the Bot.
// Currently, there is only one implementation, RedisState.
package statekeeping

import "github.com/bwmarrin/discordgo"

// State defines a state-keeping interface for the Bot.
type State interface {
	// Guilds returns a generously-populated slice of all guilds.
	Guilds() []*discordgo.Guild
	// Channels returns a list of all the channels in a guild.
	Channels(guildID string) []*discordgo.Channel
	// Members returns all members of a guild.
	Members(guildID string) []*discordgo.Member
	// Messages returns a slice of messages in the cache matching the given channel ID.
	Messages(channelID string) []*discordgo.Message

	// User returns a user.
	User(userID string) *discordgo.User
	// Presence returns a user presence, if allowed.
	Presence(userID string) *discordgo.Presence
	// Member returns a guild member, if allowed.
	Member(guildID, userID string) *discordgo.Member
	// Message returns a message if still present in the cache.
	Message(channelID, messageID string) *discordgo.Message
	// Role returns a guild role.
	Role(guildID, roleID string) *discordgo.Role

	// Events returns the number of events processed.
	Events() string
	// Size returns the size of the state cache in human-readable form.
	Size() string

	// OnEvent is the discordgo event handler for all events.
	OnEvent(session *discordgo.Session, event interface{})
}
