package main

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"os/signal"
)

func main() {
	// Read environment variables.
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("No token specified; use env DISCORD_TOKEN")
		return
	}

	// Create a Discord session.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.WithError(err).Fatal("Couldn't create Discord session")
		return
	}

	// Get app info, print invite link.
	app, err := s.Application("@me")
	if err != nil {
		log.WithError(err).Fatal("Couldn't get application info")
		return
	}
	perms := discordgo.PermissionSendMessages | discordgo.PermissionManageChannels | discordgo.PermissionManageRoles
	fmt.Printf("https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d\n", app.ID, perms)

	// Create the bot.
	b := New(s, TheEntireBeeMovieScript)

	// Set up a context that's canceled by any signal.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit)
		<-quit
		signal.Reset()
		cancel()
	}()

	// Run the bot.
	if err := b.Run(ctx); err != nil {
		log.WithError(err).Error("Error")
	}
}
