package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/willaiemi/go-bot/internal/commands"
	"github.com/willaiemi/go-bot/internal/database"
)

var session *discordgo.Session
var botToken string

func init() {
	log.Println("Retrieving BOT_TOKEN from environment variables...")
	envBotToken := os.Getenv("BOT_TOKEN")

	if envBotToken == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	botToken = envBotToken
}

func init() {
	var err error
	log.Println("Creating Discord session...")
	session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
}

func init() {
	log.Println("Try to connect to DB...")
	database.RunDb()
}

func main() {
	log.Println("Starting bot...")

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := session.Open()

	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	defer session.Close()

	log.Println("Registering commands...")
	commands.RegisterCommands(session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Bot is now operational...")
	<-stop

	log.Println("Shutting down bot...")
	err = commands.RemoveCommands(session)

	if err != nil {
		log.Panicf("Error removing commands: %v", err)
	} else {
		log.Println("Commands removed successfully.")
	}
}
