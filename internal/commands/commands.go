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
			Description: "Adds a new to-do item",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "The title of the to-do item",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "Lists your to-do items",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "filter by",
					Description: "Which items to list",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "all",
							Value: todo.ListFilterAll,
						},
						{
							Name:  "pending",
							Value: todo.ListFilterPending,
						},
						{
							Name:  "completed",
							Value: todo.ListFilterCompleted,
						},
					},
				},
			},
		},
		{
			Name:        "done",
			Description: "Marks a to-do item as done",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "The ID of the to-do item to mark as done",
					Required:    true,
				},
			},
		},
		{
			Name:        "edit",
			Description: "Edits a to-do item title",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "The ID of the to-do item to edit",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "The new title of the to-do item",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete",
			Description: "Deletes a to-do item",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "The ID of the to-do item to delete",
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
						Content: "Error adding to-do item. Please try again.",
					},
				})
				return
			}

			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:         session,
				interaction:     interaction,
				listFilter:      todo.ListFilterPending,
				highlightTodoID: createdTodo.ID,
			})
		},
		"list": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			listFilter := todo.ListFilterPending

			options := interaction.ApplicationCommandData().Options

			if len(options) > 0 {
				switch int(options[0].IntValue()) {
				case int(todo.ListFilterAll):
					listFilter = todo.ListFilterAll
				case int(todo.ListFilterCompleted):
					listFilter = todo.ListFilterCompleted
				}
			}

			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:     session,
				interaction: interaction,
				listFilter:  listFilter,
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

			todoID := uint32(idOption.IntValue())

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

			_, err = todo.MarkTodoDone(userID, todoID)

			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Error marking to-do item as done: %s", err.Error()),
					},
				})
				return
			}

			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:         session,
				interaction:     interaction,
				listFilter:      todo.ListFilterPending,
				highlightTodoID: todoID,
			})
		},
		"edit": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

			todoID := uint32(idOption.IntValue())
			title := titleOption.StringValue()

			editedTodo, err := todo.EditTodo(userID, todoID, title)

			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error editing to-do item. Please try again.",
					},
				})
				return
			}

			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:         session,
				interaction:     interaction,
				listFilter:      todo.ListFilterPending,
				highlightTodoID: editedTodo.ID,
			})
		},
		"delete": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

			todoID := uint32(idOption.IntValue())

			deletedTodo, err := todo.DeleteTodo(userID, todoID)

			if err != nil {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error deleting to-do item. Please try again.",
					},
				})
				return
			}

			listFilter := todo.ListFilterPending

			if deletedTodo.Done {
				listFilter = todo.ListFilterCompleted
			}

			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:     session,
				interaction: interaction,
				listFilter:  listFilter,
			})
		},
	}

	componentHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"list_completed": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:     session,
				interaction: interaction,
				listFilter:  todo.ListFilterCompleted,
			})
		},
		"list_pending": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			respondInteractionWithTodoList(RespondInteractionWithTodoListParams{
				session:     session,
				interaction: interaction,
				listFilter:  todo.ListFilterPending,
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

type RespondInteractionWithTodoListParams struct {
	session         *discordgo.Session
	interaction     *discordgo.InteractionCreate
	highlightTodoID uint32
	listFilter      todo.ListFilter
}

func respondInteractionWithTodoList(p RespondInteractionWithTodoListParams) {
	emptyDefaultResponse := "No to-do items. Add one with `/add`!"

	if p.listFilter == todo.ListFilterCompleted {
		emptyDefaultResponse = "No completed to-do items. Complete a to-do with `/done`!"
	}

	userID, err := getUserID(p.interaction)

	if err != nil {
		p.session.InteractionRespond(p.interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error: Unable to retrieve user information.",
			},
		})
		return
	}

	todosList := todo.GetFilteredTodos(userID, p.listFilter, p.highlightTodoID)

	responseContent := ""

	for _, todoItem := range todosList {
		if todoItem.ID == p.highlightTodoID {
			responseContent += fmt.Sprintf("__%s__\n\n", todoItem.String())
		} else {
			responseContent += fmt.Sprintf("%s\n\n", todoItem.String())
		}
	}

	if responseContent == "" {
		responseContent = emptyDefaultResponse
	}

	var button *discordgo.Button
	var embedColor int

	switch p.listFilter {
	case todo.ListFilterPending:
		button = &discordgo.Button{
			Style: discordgo.SuccessButton,
			Emoji: &discordgo.ComponentEmoji{
				Name: "✅",
			},
			Label:    "Show Completed",
			CustomID: "list_completed",
		}
		embedColor = 0x000370
	default:
		button = &discordgo.Button{
			Style: discordgo.SecondaryButton,
			Emoji: &discordgo.ComponentEmoji{
				Name: "📜",
			},
			Label:    "Show Pending",
			CustomID: "list_pending",
		}
		embedColor = 0x3DC13C
	}

	p.session.InteractionRespond(p.interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "TO-DO",
					Description: responseContent,
					Color:       embedColor,
				},
			},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						button,
					},
				},
			},
		},
	})

}

func RegisterCommands(session *discordgo.Session) error {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
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
