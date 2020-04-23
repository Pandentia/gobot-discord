package gobot // gobot is a command framework for discordgo-based bots.

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Bot represents a Discord bot with commands.
type Bot struct {
	// The underlying, preconfigured discordgo session.
	Session *discordgo.Session
	// The bot's prefix handler.
	Prefix func(msg *discordgo.Message) string
	// The bot's self-description, used in help.
	Description string

	commands map[string]*Command // All the bot's commands. Do not write to this map directly.
	me       *discordgo.User     // Our bot's user, populated after READY.
}

// Init initializes a new Bot instance and registers event handlers.
func (bot *Bot) Init() {
	bot.commands = make(map[string]*Command)
	bot.RegisterCommand(DefaultHelper())
	// finally, register our handlers
	bot.Session.AddHandler(bot.handleReady)
	bot.Session.AddHandler(bot.handleMessage)
}

// RegisterCommand registers (and overwrites) a command in the bot.
func (bot *Bot) RegisterCommand(cmd *Command) {
	if cmd.Category == "" {
		cmd.Category = "Generic"
	}
	bot.commands[strings.ToLower(cmd.Name)] = cmd
}

// RegisterCommands registers multiple commands at once.
// This is convenient for, for instance, importing modules created by others.
func (bot *Bot) RegisterCommands(cmds []*Command) {
	for _, cmd := range cmds {
		bot.RegisterCommand(cmd)
	}
}

// ListCommands returns a slice of all command names.
func (bot *Bot) ListCommands() []string {
	cmds := make([]string, 0, len(bot.commands))
	for cmd := range bot.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// GetCommand gets a command pointer.
// Returns nil if command does not exist.
func (bot *Bot) GetCommand(commandName string) *Command {
	if cmd, ok := bot.commands[strings.ToLower(commandName)]; ok {
		return cmd
	}
	return nil
}

// RemoveCommand removes a command from the bot.
// If the command does not exist, it does nothing.
func (bot *Bot) RemoveCommand(commandName string) {
	delete(bot.commands, strings.ToLower(commandName))
}

// Me returns the discord bot's own User instance.
// May return nil before the first READY.
func (bot *Bot) Me() *discordgo.User {
	return bot.me
}
