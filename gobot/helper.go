package gobot

import (
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func helpMenu(ctx *Context) {
	embed := &discordgo.MessageEmbed{}
	embed.Title = "Command Help for " + ctx.Me.Username
	embed.Description = ctx.Bot.Description
	embed.Fields = make([]*discordgo.MessageEmbedField, 0, 1)
	categories, commandMap := aggregateCommands(ctx.Bot.Commands)
	for _, category := range categories {
		commands := make([]string, 0, 1)
		for _, command := range commandMap[category] {
			commands = append(commands, command.Name)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   category,
			Value:  "`" + strings.Join(commands, ", ") + "`",
			Inline: true,
		})
	}
	ctx.ReplyWithEmbed(embed)
}

func commandDescription(ctx *Context) {
	command, ok := ctx.Bot.Commands[ctx.Args[0]]
	if !ok {
		ctx.Reply("Command not found.")
		return
	}

	embed := &discordgo.MessageEmbed{}
	embed.Title = "Command Reference"
	embed.Description = "`" + ctx.Prefix + command.Name + "`"
	embed.Fields = make([]*discordgo.MessageEmbedField, 0, 1)
	if command.Description != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Description",
			Value:  command.Description,
			Inline: false,
		})
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Category",
		Value:  command.Category,
		Inline: true,
	})
	ctx.ReplyWithEmbed(embed)
}

// DefaultHelper returns the default help command.
func DefaultHelper() *Command {
	runner := func(ctx *Context) {
		if len(ctx.Args) == 0 {
			helpMenu(ctx)
		} else {
			commandDescription(ctx)
		}
	}

	return &Command{
		Name:        "help",
		Description: "Provides a list of all bot commands.",
		Category:    "Generic",
		Runner:      runner,
	}
}

func aggregateCommands(commands map[string]*Command) (categories []string, commandMap map[string][]*Command) {
	// aggregate commands first, easier
	commandMap = make(map[string][]*Command)
	for _, cmd := range commands {
		cat := cmd.Category
		if _, ok := commandMap[cat]; !ok {
			commandMap[cat] = make([]*Command, 0, 1)
		}
		commandMap[cat] = append(commandMap[cat], cmd)
	}

	// sort and aggregate categories
	categories = make([]string, 0, len(commandMap))
	for category := range commandMap {
		sort.Slice(commandMap[category], func(i, j int) bool {
			return strings.Compare(commandMap[category][i].Name, commandMap[category][j].Name) == -1
		})
		categories = append(categories, category)
	}
	sort.Strings(categories)

	return
}
