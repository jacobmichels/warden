package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Discord interface {
	RegisterSlashCommands(commands []*discordgo.ApplicationCommand) ([]string, error)
	UnregisterSlashCommands(ids []string) error
	AddReadyHandler(func(s *discordgo.Session, r *discordgo.Ready))
	AddCommandHandlers(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	Start() error
	Close() error
}
