package main

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jacobmichels/warden/internal/config"
	"github.com/jacobmichels/warden/internal/discord"
	"github.com/jacobmichels/warden/internal/warden"
)

func main() {
	ctx := context.Background()

	err := config.Init()
	if err != nil {
		log.Fatalf("failed to initialize viper config: %s", err)
	}

	client, err := discord.NewClient("Bot " + config.GetStr("discord.token"))
	if err != nil {
		log.Fatalf("failed to create discord client: %s", err)
	}

	commands, handlers := createApplicationCommands()

	err = warden.Start(ctx, client, commands, handlers)
	if err != nil {
		log.Fatal(err)
	}
}

func createApplicationCommands() ([]*discordgo.ApplicationCommand, map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "whitelist",
			Description: "Manage your minecraft server's whitelist from the comfort of Discord",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Add a player to the whitelist",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The minecraft username of the player to whitelist",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove a player from the whitelist",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The minecraft username of the player to whitelist",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	}

	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"whitelist-add": func(_ *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Printf("got request to add %s to the whitelist", i.ApplicationCommandData().Options[0].Options[0].StringValue())
		},
		"whitelist-remove": func(_ *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Printf("got request to remove %s from the whitelist", i.ApplicationCommandData().Options[0].Options[0].StringValue())
		},
	}

	return commands, handlers
}
