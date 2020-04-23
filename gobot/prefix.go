package gobot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// SinglePrefixHandler returns a simple single-prefix handler for the bot.
func SinglePrefixHandler(prefix string) func(*discordgo.Message) string {
	return func(msg *discordgo.Message) string {
		if strings.HasPrefix(msg.Content, prefix) {
			return prefix
		}
		return ""
	}
}
