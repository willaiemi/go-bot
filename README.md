# Go Discord Bot – TODO Manager

A lightweight Discord bot written in Go that helps users manage personal TODO items through simple slash commands. Built using [discordgo](https://github.com/bwmarrin/discordgo).

---

## Features

- `/ping`: Check if the bot is responsive.
- `/add [title]`: Add a new TODO item.
- `/list [filter]`: View TODO items (all, pending, or completed).
- `/done [id]`: Mark a TODO item as done.
- `/edit [id] [title]`: Edit the title of a TODO.
- `/delete [id]`: Delete a TODO item.
- Buttons: Toggle between pending and completed items interactively.

Each user's TODOs are currently stored in-memory and scoped by Discord user ID.

---

## Getting Started

### Prerequisites

- Go 1.20+
- A Discord bot token

### Local Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/willaiemi/go-bot
   cd go-bot
   ```

2. Create a `.env` file with the following variables:

   ```env
   BOT_TOKEN=your_discord_bot_token

   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=tododb
   DB_HOST=postgres
   DB_PORT=5432
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Run the bot:

   ```bash
   go run main.go
   ```

---

## Docker Setup (Optional)

The project includes a working `docker-compose.yml` for spinning up both the bot and a PostgreSQL database.

1. Start the services:

   ```bash
   docker-compose up --build
   ```

2. This setup builds the Go app and waits for the Postgres container to become healthy before launching the bot.

> **Note:** The bot currently stores TODOs in memory, but PostgreSQL integration is in progress.

---

## Project Structure

- `internal/commands`: Handles slash command registration and interactions.
- `internal/todo`: Manages in-memory TODO lists.
- `internal/database`: PostgreSQL connection pool and DB setup (coming soon).
- `main.go`: Application entry point.

---

## Example Usage

```
/add title:Buy milk
/list filter:pending
> ▪ **Buy milk** (ID: 1)
/done id:1
> ✅ ~~**Buy milk**~~ (ID: 1)
/edit id:1 title:Buy chocolate milk
/delete id:1
```

---

## Roadmap

- [x] Slash command support
- [x] Button-based navigation
- [ ] PostgreSQL persistence
- [ ] Add due dates or timestamps
- [ ] Pagination for long lists
- [ ] Unit tests and structured logging
