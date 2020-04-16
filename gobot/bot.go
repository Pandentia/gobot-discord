package gobot // gobot is a command framework for discordgo-based bots.

import "github.com/bwmarrin/discordgo"

// Bot represents a Discord bot with commands.
type Bot struct {
	Session *discordgo.Session // The underlying, preconfigured discordgo session.

	Prefixes    []string // The prefixes the bot listens to.
	Description string   // The bot's self-description, used in help.

	Commands map[string]Command // All the bot's commands. Do not write to this map directly.
}

// Init initializes a new Bot instance and registers event handlers.
func (bot *Bot) Init() {
	bot.Commands = make(map[string]Command)
	bot.RegisterCommand(DefaultHelper())
	// finally, register our handler
	bot.Session.AddHandler(bot.handleMessage)
}

// RegisterCommand registers (and overwrites) a command in the bot.
func (bot *Bot) RegisterCommand(cmd Command) {
	bot.Commands[cmd.Name] = cmd
}
