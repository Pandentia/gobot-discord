package gobot // gobot is a command framework for discordgo-based bots.

import (
	"github.com/Pandentia/gobot-discord/gobot/statekeeping"
	"github.com/bwmarrin/discordgo"
)

// Bot represents a Discord bot with commands.
type Bot struct {
	Session *discordgo.Session // The underlying, preconfigured discordgo session.
	State   statekeeping.State // Our bot's statekeeping. May be nil to disable.

	Prefixes    []string // The prefixes the bot listens to.
	Description string   // The bot's self-description, used in help.

	Commands map[string]*Command // All the bot's commands. Do not write to this map directly.
	Me       *discordgo.User     // Our bot's user, populated after READY.
}

// Init initializes a new Bot instance and registers event handlers.
func (bot *Bot) Init() {
	bot.Commands = make(map[string]*Command)
	bot.RegisterCommand(DefaultHelper())
	// finally, register our handlers
	bot.Session.AddHandler(bot.handleReady)
	bot.Session.AddHandler(bot.handleMessage)
	// register state handlers
	if bot.State != nil {
		bot.Session.AddHandler(bot.State.OnEvent)
	}
}

// RegisterCommand registers (and overwrites) a command in the bot.
func (bot *Bot) RegisterCommand(cmd *Command) {
	if cmd.Category == "" {
		cmd.Category = "Generic"
	}
	bot.Commands[cmd.Name] = cmd
}

// RegisterCommands registers multiple commands at once.
// This is convenient for, for instance, importing modules created by others.
func (bot *Bot) RegisterCommands(cmds []*Command) {
	for _, cmd := range cmds {
		bot.RegisterCommand(cmd)
	}
}
