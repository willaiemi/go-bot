package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/willaiemi/go-bot/internal/todo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with Pong!",
		},
		{
			Name:        "add",
			Description: "Adds a new TO-DO item",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "The title of the TO-DO item",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "Lists all your TO-DO items",
		},
		{
			Name:        "done",
			Description: "Marks a TO-DO item as done",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "The ID of the TO-DO item to mark as done",
					Required:    true,
				},
			},
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
		"add": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			options := interaction.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			titleOption, ok := optionMap["title"]
			if !ok {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: Title option is required.",
					},
				})
				return
			}

			userID, err := getUserID(interaction)
			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: Unable to retrieve user information.",
					},
				})
				return
			}

			title := titleOption.StringValue()

			createdTodo, err := todo.AddTodo(userID, title)

			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error adding TO-DO item. Please try again.",
					},
				})
				return
			}

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added TO-DO item: **%s** (ID: %d)", createdTodo.Title, createdTodo.ID),
				},
			})
		},
		"list": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			userID, err := getUserID(interaction)
			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: Unable to retrieve user information.",
					},
				})
				return
			}

			todosList := todo.GetTodos(userID)

			if len(todosList) == 0 {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You have no TO-DO items. Create one with `/add`!",
					},
				})
				return
			}

			responseContent := ""
			for _, todoItem := range todosList {
				responseContent += fmt.Sprintf("> **%s** (ID: %d)\n", todoItem.Title, todoItem.ID)
			}
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "TO-DO",
							Description: responseContent,
							Color:       0x000370,
						},
					},
				},
			})
		},
		"done": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			options := interaction.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			idOption, ok := optionMap["id"]

			if !ok {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: ID option is required.",
					},
				})
				return
			}
			todoID := idOption.IntValue()
			userID, err := getUserID(interaction)
			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: Unable to retrieve user information.",
					},
				})
				return
			}

			markedTodo, err := todo.MarkTodoDone(userID, uint32(todoID))
			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Error marking TO-DO item as done: %s", err.Error()),
					},
				})
				return
			}
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Marked TO-DO item as done: ~~**%s**~~ (ID: %d)", markedTodo.Title, markedTodo.ID),
				},
			})
		},
	}

	registeredCommandsIds = make([]string, len(commands))
)

func getUserID(interaction *discordgo.InteractionCreate) (string, error) {
	if interaction.Member != nil {
		return interaction.Member.User.ID, nil
	} else if interaction.User != nil {
		return interaction.User.ID, nil
	}
	return "", fmt.Errorf("unable to retrieve user information")
}

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
