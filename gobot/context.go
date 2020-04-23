package gobot

import "github.com/bwmarrin/discordgo"

// Context provides command context.
type Context struct {
	Bot     *Bot     // The bot instance.
	Prefix  string   // The prefix this command was invoked with.
	Command *Command // The invoked command.
	Args    []string // The command arguments.

	Message   *discordgo.MessageCreate // The message that triggered this command.
	Author    *discordgo.Member        // Shorthand for Message.Author.
	ChannelID string                   // Shorthand for Message.ChannelID.
	GuildID   string                   // Shorthand for Message.GuildID.
}

// Me is a shorthand helper to Bot.Me().
func (c *Context) Me() *discordgo.User {
	return c.Bot.Me()
}

// Reply is shorthand to send a message to the channel the command was invoked
// from.
func (c *Context) Reply(msg string) (*discordgo.Message, error) {
	return c.Bot.Session.ChannelMessageSend(c.Message.ChannelID, msg)
}

// ReplyWithEmbed is shorthand to send a message to the channel the command was
// invoked from, with an embed.
func (c *Context) ReplyWithEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Bot.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
}
