package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/Pandentia/gobot-discord/gobot"
	"github.com/Pandentia/gobot-discord/gobot/statekeeping"
	"github.com/bwmarrin/discordgo"
	"github.com/mediocregopher/radix/v3"
)

func main() {
	token := flag.String("token", "", "The bot token to log in with.")
	redis := flag.String("redis", "", "Specifies the redis to use for statekeeping.")
	flag.Parse()
	if *token == "" {
		return
	}

	var state statekeeping.State
	if *redis != "" {
		pool, err := radix.NewPool("tcp", *redis, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		state = &statekeeping.RedisState{Redis: pool}
	}

	session, err := discordgo.New("Bot " + *token)
	if err != nil {
		fmt.Println(err)
		return
	}
	bot := gobot.Bot{
		Session:     session,
		Prefixes:    []string{"?"},
		Description: "A testing discord bot",
		State:       state,
	}
	bot.Init()
	bot.RegisterCommand(&gobot.Command{
		Name:        "ping",
		Description: "Tests if the bot is working",
		Runner: func(ctx *gobot.Context) {
			ctx.Reply("Pong!")
		},
	})
	bot.RegisterCommand(&gobot.Command{
		Name:        "stats",
		Description: "Returns bot statistics",
		Runner: func(ctx *gobot.Context) {
			if ctx.Bot.State == nil {
				ctx.Reply("Bot does not have statekeeping enabled.")
				return
			}
			embed := &discordgo.MessageEmbed{
				Title: "Bot statistics",
			}
			embed.Fields = []*discordgo.MessageEmbedField{
				{
					Name:   "Events processed",
					Value:  ctx.State.Events(),
					Inline: true,
				},
				{
					Name:   "State cache size",
					Value:  ctx.State.Size(),
					Inline: true,
				},
			}
			ctx.ReplyWithEmbed(embed)
		},
	})

	session.Open()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	session.Close()
}
