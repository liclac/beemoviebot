package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"sync"
)

const (
	MaxMessageLength = 2000
	ChannelName      = "the-bee-movie"
)

type Bot struct {
	Session *discordgo.Session

	chunks []string
	wg     sync.WaitGroup
}

func New(s *discordgo.Session, msg string) *Bot {
	return &Bot{
		Session: s,
		chunks:  MakeChunks(msg, MaxMessageLength),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	// Register handlers, unregister them again upon returning.
	defer b.Session.AddHandler(b.HandleGuildCreate)()
	defer b.Session.AddHandler(b.HandleGuildDelete)()

	// Connect to Discord.
	if err := b.Session.Open(); err != nil {
		return err
	}

	// Log info.
	log.Infof("Chunks: %d", len(b.chunks))

	// Wait for the context to terminate.
	<-ctx.Done()

	// Wait for ongoing tasks to finish.
	b.wg.Wait()

	// Disconnect from Discord.
	if err := b.Session.Close(); err != nil {
		return err
	}

	return nil
}

func (b *Bot) HandleGuildCreate(_ *discordgo.Session, g *discordgo.GuildCreate) {
	log.WithField("gid", g.ID).Info("Joined guild!")

	// Ensure we don't get shut down in the middle.
	b.wg.Add(1)
	defer b.wg.Done()

	// Ensure we leave cleanly when we're done.
	defer func() {
		if err := b.Session.GuildLeave(g.ID); err != nil {
			log.WithError(err).Error("Couldn't leave guild!")
		}
	}()

	// Quit if there's already a bee movie channel.
	for _, c := range g.Channels {
		if c.Name == ChannelName {
			log.WithField("gid", g.ID).Warn("Channel already exists, quitting")
			return
		}
	}

	// Create the bee movie channel.
	c, err := b.Session.GuildChannelCreate(g.ID, ChannelName, "text")
	if err != nil {
		log.WithError(err).Error("Couldn't create channel")
		return
	}

	// SHITPOST!
	for i, chunk := range b.chunks {
		if _, err := b.Session.ChannelMessageSend(c.ID, chunk); err != nil {
			log.WithError(err).Error("Couldn't send message")

			if _, err := b.Session.ChannelDelete(c.ID); err != nil {
				log.WithError(err).Error("Couldn't delete channel")
			}
			return
		}
		log.WithField("gid", g.ID).Infof("Sent chunk: %d / %d", i+1, len(b.chunks))
	}
}

func (b *Bot) HandleGuildDelete(_ *discordgo.Session, g *discordgo.GuildDelete) {
	log.WithField("gid", g.ID).Info("Left guild")
}
