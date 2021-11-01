package main

import (
	"fmt"
	"image/color"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	env "github.com/joho/godotenv"
	"github.com/zackwn/pollbot/gol"
	"github.com/zackwn/pollbot/poll"
)

func main() {
	env.Load()
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating bot session:", err)
	}

	// add message handler
	dg.AddHandler(messageHandler)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is ready")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if message.Author.ID == session.State.User.ID {
		return
	}

	if len(message.Content) == 0 {
		return
	}

	if message.Content[0] == '!' {
		// message content
		mc := strings.Split(message.Content, " ")
		// mc := strings.SplitN(message.Content, " ", 2)
		if len(mc) < 1 {
			return
		}

		fmt.Println(len(mc), "message:", mc)
		switch mc[0] {
		case "!poll":
			text := strings.Join(mc[1:], " ")
			poll, err := poll.NewPoll(message.Author.ID, []rune(text))
			if err != nil {
				session.ChannelMessageSend(message.ChannelID, err.Error())
				return
			}

			embed := poll.BuildEmbed()
			author := &discordgo.MessageEmbedAuthor{
				Name:    message.Author.Username,
				IconURL: message.Author.AvatarURL("48x48"),
			}
			embed.Author = author

			poll.Start(embed, session, message.Message)
			session.ChannelMessageDelete(message.ChannelID, message.ID)
			return
		case "!lissajous":
			var err error
			// set up palette colors
			var palette, n = []color.Color{color.Black}, 1
			if len(mc) >= 2 {
				n, err = strconv.Atoi(mc[1])
				if err != nil || n <= 0 || n > 100 {
					session.ChannelMessageSend(message.ChannelID, "error: Invalid number of colors (valid from 1 to 100)")
					return
				}
			}
			for i := 0; i < n; i++ {
				palette = append(palette, RandomColor())
			}
			// get cycles
			var cycles float64 = 4
			if len(mc) >= 3 {
				cycles, err = strconv.ParseFloat(mc[2], 64)
				if err != nil || cycles <= 0 || cycles > 10 {
					session.ChannelMessageSend(message.ChannelID, "error: Invalid number of cycles (valid from 0.1 to 10.0)")
					return
				}
			}
			// generate and send lissajous
			r := Lissajous(cycles, palette)
			session.ChannelFileSend(message.ChannelID, "lissajous.gif", r)
		case "!life":
			gol.Run(session, message)
		}
	}
}
