# LogGob - Clash Royale Battle Logger

LogGob is a Clash Royale battle logger that fetches battle data from the Clash Royale API and stores it in a SQLite database for analysis and tracking.

## Features

- Fetches battle logs from Clash Royale API
- Stores battle data in SQLite database
- Filters for Ladder game mode battles only
- Supports PvP battle tracking
- Environment variable configuration
- Structured data models for battles, players, cards, arenas, and game modes

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

The application fetches battle logs for player `#PLY2Q2LL` and stores them in the database. Modify the player tag in `cmd/main.go` if you want to track a different player.

## Project Structure

```
log-gob/
├── cmd/
│   └── main.go           # Main application entry point
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
