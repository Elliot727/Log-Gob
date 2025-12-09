# LogGob - Clash Royale Battle Logger

LogGob is a Clash Royale battle logger that fetches battle data from the Clash Royale API and stores it in a SQLite database for analysis and tracking.

## Features

- Fetches battle logs from Clash Royale API
- Stores battle data in SQLite database
- Filters for Ladder game mode battles only
- Supports PvP battle tracking
- Environment variable configuration
- Structured data models for battles, players, cards, arenas, and game modes
- Interactive TUI for viewing battle history (using Bubble Tea and Lip Gloss for styling)
- Makefile for easy building and running
- Colorful terminal interface with detailed battle information and cards
- Battle result display (Victory/Loss/Draw) prominently shown
- Side-by-side deck comparison showing team vs opponent cards
- Properly aligned card levels in table format
- Human-readable time format (YYYY-MM-DD HH:MM)

## Prerequisites

- Go 1.25 or higher
- Clash Royale API key

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/elliot727/log-gob.git
   cd log-gob
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file with your API key (see `.env.example`)

4. Run the application:
   ```bash
   make run
   ```

## Configuration

Create a `.env` file in the root directory with your Clash Royale API key and player tag:

```
APIKEY=your_clash_royale_api_key_here
PLAYERTAG=your_player_tag_here
```

If you don't specify a `PLAYERTAG`, the application will use a default tag and show a warning message.

## Database Schema

The application creates the following tables:

- `arenas` - Clash Royale arenas
- `gamemodes` - Game modes like Ladder, Friendly, etc.
- `players` - Player information
- `cards` - Card data (name, level, rarity, etc.)
- `battles` - Battle records
- `battle_participants` - Players in each battle
- `battle_decks` - Cards used in each battle

## Usage

### CLI Version
The main application fetches battle logs from the API and stores them in the database:
```bash
go run cmd/main.go
# Or using make:
make run
```

### TUI Version
The TUI version allows you to interactively view battles stored in the database:
```bash
go run cmd/tui/main.go
# Or using make:
make tui
```

Controls in TUI:
- `J` or `Down Arrow`: Navigate down through battles
- `K` or `Up Arrow`: Navigate up through battles
- `R`: Refresh battles from database
- `Q` or `Ctrl+C`: Quit the application

Make sure to run the CLI version first to populate the database with battle data before using the TUI.

The TUI displays:
- Battle result (Victory/Loss/Draw) prominently after the battle header
- Human-readable time format (YYYY-MM-DD HH:MM)
- Detailed battle information including teams, opponents, and crowns
- Side-by-side deck comparison showing team vs opponent cards
- Properly aligned card levels in table format for easy comparison

## Project Structure

```
log-gob/
├── cmd/
│   ├── main.go           # CLI application - fetches battles from API and stores to database
│   └── tui/
│       └── main.go       # TUI application - interactive terminal interface to view battles
├── internal/
│   ├── api/
│   │   └── client.go     # API client implementation
│   ├── storage/
│   │   └── storage.go    # Database operations
│   └── types/
│       ├── battle.go     # Battle data structure
│       ├── player.go     # Player data structure
│       ├── card.go       # Card data structure
│       ├── arena.go      # Arena data structure
│       └── gameMode.go   # Game mode data structure
├── .env.example          # Example environment file
├── .gitignore            # Git ignore rules
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── battles.db            # SQLite database (automatically created when the application runs)
└── README.md             # This file
```

## Environment Variables

- `APIKEY` - Your Clash Royale API key (required)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
