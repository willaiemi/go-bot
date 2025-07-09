# Go Discord Bot – TODO Manager

A lightweight Discord bot written in Go that helps users manage personal TODO items through simple slash commands. Built using [discordgo](https://github.com/bwmarrin/discordgo).

## Features (so far)

- `/ping`: Check if the bot is responsive.
- `/add [title]`: Add a new TODO item with a title.
- `/list`: View your list of TODOs.
- `/done [id]`: Mark a TODO item as done by its ID.

Each user's TODOs are isolated by their Discord user ID.

## Getting Started

### Prerequisites

- Go 1.20+
- A Discord bot token
- A `.env` file with the following:
  ```env
  BOT_TOKEN=your_discord_bot_token
  ```

### Running the Bot

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/go-bot
   cd go-bot
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the bot:

   ```bash
   go run main.go
   ```

The bot will register its slash commands on startup and listen for interactions.

## Structure

- `internal/commands`: Registers slash commands and handles user interactions.
- `internal/todo`: Manages in-memory TODO lists scoped by user ID.

## Example

```
/add title:Buy milk
/list
> Buy milk (ID: 1)
/done 1
> ~~Buy milk~~ (ID: 1)
```

## Future Improvements

- Integrate with a persistent storage backend (e.g., PostgreSQL, Redis, Supabase);
- Add timestamps and list filtering (ex: completed vs active);
- Support reminders or due dates;
- Edit/delete a TODOs.
