package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/Pandentia/gobot-discord/gobot"
	"github.com/bwmarrin/discordgo"
)

func main() {
	token := flag.String("token", "", "The bot token")
	flag.Parse()
	if *token == "" {
		return
	}

	session, err := discordgo.New("Bot " + *token)
	if err != nil {
		fmt.Println(err)
		return
	}
	bot := gobot.Bot{
		Session:     session,
		Prefix:      gobot.SinglePrefixHandler("?"),
		Description: "A testing discord bot",
	}
	bot.Init()
	bot.RegisterCommand(&gobot.Command{
		Name:        "ping",
		Description: "A ping command",
		Category:    "Generic",
		Runner: func(ctx *gobot.Context) {
			ctx.Reply("Pong!")
		},
	})

	session.Open()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	session.Close()
}
