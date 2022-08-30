package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Discordgo struct {
	session *discordgo.Session
}

func NewClient(token string) (*Discordgo, error) {
	session, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	return &Discordgo{
		session,
	}, nil
}

func (d *Discordgo) Start() error {
	return d.session.Open()
}

func (d *Discordgo) Close() error {
	return d.session.Close()
}

// Registers the passed commands as global, returns a slice of their IDs.
func (d *Discordgo) RegisterSlashCommands(commands []*discordgo.ApplicationCommand) ([]string, error) {
	ids := make([]string, 0, len(commands))

	for _, cmd := range commands {
		res, err := d.session.ApplicationCommandCreate(d.session.State.User.ID, "", cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to register %s application command: %w", cmd.Name, err)
		}

		ids = append(ids, res.ID)
	}

	return ids, nil
}

func (d *Discordgo) UnregisterSlashCommands(ids []string) error {
	for _, id := range ids {
		err := d.session.ApplicationCommandDelete(d.session.State.User.ID, "", id)
		if err != nil {
			return fmt.Errorf("failed to unregister %s application command: %w", id, err)
		}
	}

	return nil
}

func (d *Discordgo) AddReadyHandler(f func(s *discordgo.Session, r *discordgo.Ready)) {
	d.session.AddHandler(f)
}

func (d *Discordgo) AddCommandHandlers(handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	d.session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		data := interaction.ApplicationCommandData()
		if len(data.Options) != 1 {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Invalid command",
				},
			})
			if err != nil {
				log.Printf("interaction respond failed: %s", err)
			}
		}

		if handler, ok := handlers[fmt.Sprintf("%s-%s", data.Name, data.Options[0].Name)]; ok {
			handler(session, interaction)
		}
	})
}
