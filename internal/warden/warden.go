package warden

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/jacobmichels/warden/internal/discord"
)

func Start(ctx context.Context, client discord.Discord, commands []*discordgo.ApplicationCommand, handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) error {
	// register handlers before we open the connection so we don't miss anything
	client.AddReadyHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Discord client connected. Logged in as %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	client.AddCommandHandlers(handlers)

	// open the websocket connection
	err := client.Start()
	if err != nil {
		return fmt.Errorf("failed to start discord client: %w", err)
	}
	defer client.Close()

	ids, err := client.RegisterSlashCommands(commands)
	if err != nil {
		return fmt.Errorf("failed to register slash commands: %w", err)
	}

	// wait for ctrl+c
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)

	log.Println("Commands registered successfully, blocking for interrupt")
	<-ctx.Done()
	log.Println("Interrupt signal receieved, gracefully shutting down")

	err = client.UnregisterSlashCommands(ids)
	if err != nil {
		return fmt.Errorf("failed to unregister slash commands: %w", err)
	}

	log.Println("Commands unregistered")

	return nil
}
