package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jacobmichels/warden/internal/config"
	"github.com/jacobmichels/warden/internal/discord"
	"github.com/jacobmichels/warden/internal/rcon"
	"github.com/jacobmichels/warden/internal/warden"
)

func main() {
	ctx := context.Background()

	err := config.Init()
	if err != nil {
		log.Panicf("failed to initialize viper config: %s", err)
	}

	discordClient, err := discord.NewClient("Bot " + config.GetStr("discord.token"))
	if err != nil {
		log.Panicf("failed to create discord client: %s", err)
	}

	rconClient, err := rcon.NewClient(config.GetStr("rcon.address"), config.GetStr("rcon.password"))
	if err != nil {
		log.Panicf("failed to create rcon client: %s", err)
	}
	defer rconClient.Close()

	commands, handlers := createApplicationCommands(rconClient)

	err = warden.Start(ctx, discordClient, commands, handlers)
	if err != nil {
		log.Panic(err)
	}
}

func createApplicationCommands(rcon *rcon.Client) ([]*discordgo.ApplicationCommand, map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
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
				{
					Name:        "list",
					Description: "List all whitelisted players",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}

	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"whitelist-add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].Options[0].StringValue()

			log.Printf("Request receieved: ADD %s", user)

			cmd := fmt.Sprintf("whitelist add %s", user)

			output, err := rcon.SendCommand(cmd)
			if err != nil {
				log.Printf("whitelist-add command failed: %s", err)
				discord.InteractionRespond(s, i, "Internal error, please try again later. If error persists, please contact the bot owner.")

				return
			}

			if strings.Contains(output, "That player does not exist") {
				discord.InteractionRespond(s, i, fmt.Sprintf("User %s does not exist. The whitelist was not updated.", user))

				return
			}

			discord.InteractionRespond(s, i, fmt.Sprintf("User %s added to the whitelist", user))
		},
		"whitelist-remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].Options[0].StringValue()

			log.Printf("Request receieved: REMOVE %s", user)

			cmd := fmt.Sprintf("whitelist remove %s", user)

			output, err := rcon.SendCommand(cmd)
			if err != nil {
				log.Printf("whitelist-remove command failed: %s", err)
				discord.InteractionRespond(s, i, "Internal error, please try again later. If error persists, please contact the bot owner.")

				return
			}

			if strings.Contains(output, "That player does not exist") {
				discord.InteractionRespond(s, i, fmt.Sprintf("User %s does not exist. The whitelist was not updated.", user))

				return
			}

			discord.InteractionRespond(s, i, fmt.Sprintf("User %s removed from the whitelist", user))
		},
		"whitelist-list": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Request receieved: LIST")

			output, err := rcon.SendCommand("whitelist list")
			if err != nil {
				log.Printf("whitelist-list command failed: %s", err)
				discord.InteractionRespond(s, i, "Internal error, please try again later. If error persists, please contact the bot owner.")

				return
			}

			discord.InteractionRespond(s, i, output)
		},
	}

	return commands, handlers
}
