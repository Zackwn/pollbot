package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	env "github.com/joho/godotenv"
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

}
