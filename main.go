package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	env "github.com/joho/godotenv"
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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
		mc := strings.SplitN(message.Content, " ", 2)
		if len(mc) <= 1 {
			return
		}

		if mc[0] == "!poll" {
			poll := poll.NewPoll([]rune(mc[1]))

			embed := poll.BuildEmbed()
			author := &discordgo.MessageEmbedAuthor{
				Name:    message.Author.Username,
				IconURL: message.Author.AvatarURL("48x48"),
			}
			embed.Author = author

			poll.Start(embed, session, message.Message)
			return
		}
	}
}
