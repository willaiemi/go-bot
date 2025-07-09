package commands

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with Pong!",
		},
	}

	commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"ping": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		},
	}

	registeredCommandsIds = make([]string, len(commands))
)

func RegisterCommands(session *discordgo.Session) error {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			return err
		}
		registeredCommandsIds[i] = cmd.ID
	}

	return nil
}

func RemoveCommands(session *discordgo.Session) error {
	for _, cmdID := range registeredCommandsIds {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", cmdID)
		if err != nil {
			return err
		}
	}

	return nil
}
