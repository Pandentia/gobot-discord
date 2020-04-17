package gobot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) handleReady(_ *discordgo.Session, ready *discordgo.Ready) {
	bot.Me = ready.User
}

func (bot *Bot) handleMessage(_ *discordgo.Session, msg *discordgo.MessageCreate) {
	// are we ready yet?
	if bot.Me == nil {
		return
	}
	// ignore ourselves
	if msg.Author.ID == bot.Session.State.User.ID {
		return
	}
	// ignore empty messages
	if msg.Content == "" {
		return
	}

	// detect our prefix
	var prefix string
	var validPrefix bool
	for _, prefix = range bot.Prefixes {
		if strings.HasPrefix(msg.Content, prefix) {
			validPrefix = true
			break
		}
	}
	if !validPrefix {
		return
	}

	// parse the command
	commandContent := strings.ToLower(msg.Content)
	commandContent = strings.TrimPrefix(commandContent, prefix)
	if commandContent == "" {
		return
	}
	commandArgs := strings.Split(commandContent, " ")
	commandArgs[0] = strings.ToLower(commandArgs[0])
	command, exists := bot.Commands[commandArgs[0]]
	if !exists {
		return
	}

	// create context and run
	context := &Context{
		Bot:     bot,
		Prefix:  prefix,
		Command: command,
		Args:    commandArgs[1:],

		Message:   msg,
		Author:    msg.Member,
		ChannelID: msg.ChannelID,
		GuildID:   msg.GuildID,

		Me:    bot.Me,
		State: bot.State,
	}
	command.Runner(context)
}
