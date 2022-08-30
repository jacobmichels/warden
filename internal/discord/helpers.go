package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func InteractionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Printf("interaction respond failed: %s", err)
	}
}
