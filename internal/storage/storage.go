// Package storage provides database operations for saving and retrieving Clash Royale battle data.
package storage

import (
	"database/sql"

	"github.com/elliot727/log-gob/internal/types"
)

// Storage represents a database storage handler with a SQL database connection.
type Storage struct {
	DB *sql.DB
}

// NewStorage creates a new storage instance with the provided database connection.
func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

// Init creates the required database tables if they don't exist.
// This includes tables for arenas, game modes, players, cards, battles, participants, and decks.
func (s *Storage) Init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS arenas (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS gamemodes (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS players (
		tag TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		maxLevel INTEGER NOT NULL,
		rarity TEXT NOT NULL,
		elixirCost INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS battles (
		battleTime TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		arena_id INTEGER NOT NULL,
		gamemode_id INTEGER NOT NULL,
		FOREIGN KEY (arena_id) REFERENCES arenas(id),
		FOREIGN KEY (gamemode_id) REFERENCES gamemodes(id)
	);
	CREATE TABLE IF NOT EXISTS battle_participants (
		battleTime TEXT NOT NULL,
		player_tag TEXT NOT NULL,
		role TEXT NOT NULL,
		crowns INTEGER NOT NULL,
		startingTrophies INTEGER NOT NULL,
		trophyChange INTEGER NOT NULL,
		elixirLeaked REAL NOT NULL,
		PRIMARY KEY (battleTime, player_tag),
		FOREIGN KEY (battleTime) REFERENCES battles(battleTime) ON DELETE CASCADE,
		FOREIGN KEY (player_tag) REFERENCES players(tag)
	);
	CREATE TABLE IF NOT EXISTS battle_decks (
		battleTime TEXT NOT NULL,
		player_tag TEXT NOT NULL,
		card_id INTEGER NOT NULL,
		card_level INTEGER NOT NULL,
		PRIMARY KEY (battleTime, player_tag, card_id),
		FOREIGN KEY (battleTime, player_tag) REFERENCES battle_participants(battleTime, player_tag) ON DELETE CASCADE,
		FOREIGN KEY (card_id) REFERENCES cards(id)
	);
	`

	_, err := s.DB.Exec(schema)
	return err
}

// InsertBattle saves a battle and its related data to the database.
// It only processes PvP battles in the Ladder game mode and ignores other battle types/game modes.
// The function handles inserting or updating data for arenas, game modes, players, cards,
// battle records, participants, and decks.
func (s *Storage) InsertBattle(b *types.Battle) error {
	if b.BattleType != "PvP" {
		return nil
	}

	// Only store battles in the Ladder game mode (ID 72000006)
	if b.GameMode.ID != 72000006 {
		return nil
	}

	_, err := s.DB.Exec(
		"INSERT OR IGNORE INTO arenas (id, name) VALUES (?, ?)",
		b.Arena.ID, b.Arena.Name,
	)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		"INSERT OR IGNORE INTO gamemodes (id, name) VALUES (?, ?)",
		b.GameMode.ID, b.GameMode.Name,
	)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		"INSERT OR IGNORE INTO battles (battleTime, type, arena_id, gamemode_id) VALUES (?, ?, ?, ?)",
		b.BattleTime, b.BattleType, b.Arena.ID, b.GameMode.ID,
	)
	if err != nil {
		return err
	}

	allPlayers := append(b.Team, b.Opponent...)

	for _, p := range allPlayers {

		_, err := s.DB.Exec(
			"INSERT OR IGNORE INTO players (tag, name) VALUES (?, ?)",
			p.Tag, p.Name,
		)
		if err != nil {
			return err
		}

		for _, c := range p.Cards {
			_, err = s.DB.Exec(
				"INSERT OR IGNORE INTO cards (id, name, maxLevel, rarity, elixirCost) VALUES (?, ?, ?, ?, ?)",
				c.ID, c.Name, c.MaxLevel, c.Rarity, c.ElixirCost,
			)
			if err != nil {
				return err
			}
		}

		role := "opponent"
		for _, t := range b.Team {
			if t.Tag == p.Tag {
				role = "team"
				break
			}
		}

		_, err = s.DB.Exec(
			`INSERT OR REPLACE INTO battle_participants
			 (battleTime, player_tag, role, crowns, startingTrophies, trophyChange, elixirLeaked)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			b.BattleTime,
			p.Tag,
			role,
			p.Crowns,
			int64(p.StartingTrophies),
			int64(p.TrophyChange),
			p.ElixirLeaked,
		)
		if err != nil {
			return err
		}

		for _, c := range p.Cards {
			_, err := s.DB.Exec(
				`INSERT OR REPLACE INTO battle_decks
				 (battleTime, player_tag, card_id, card_level)
				 VALUES (?, ?, ?, ?)`,
				b.BattleTime,
				p.Tag,
				c.ID,
				c.Level,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetBattlesForPlayer retrieves all battles for a specific player from the database.
// Results are ordered by battle time in descending order (most recent first).
func (s *Storage) GetBattlesForPlayer(tag string) ([]types.Battle, error) {
	rows, err := s.DB.Query(`
		SELECT b.battleTime, b.type, a.id, a.name, g.id, g.name
		FROM battles b
		JOIN arenas a ON b.arena_id = a.id
		JOIN gamemodes g ON b.gamemode_id = g.id
		JOIN battle_participants bp ON bp.battleTime = b.battleTime
		WHERE bp.player_tag = ?
		ORDER BY b.battleTime DESC
	`, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []types.Battle

	for rows.Next() {
		var b types.Battle
		err := rows.Scan(
			&b.BattleTime,
			&b.BattleType,
			&b.Arena.ID,
			&b.Arena.Name,
			&b.GameMode.ID,
			&b.GameMode.Name,
		)
		if err != nil {
			return nil, err
		}

		team, err := s.loadParticipants(b.BattleTime, "team")
		if err != nil {
			return nil, err
		}
		opponent, err := s.loadParticipants(b.BattleTime, "opponent")
		if err != nil {
			return nil, err
		}

		b.Team = team
		b.Opponent = opponent
		battles = append(battles, b)
	}

	return battles, nil
}

// loadParticipants retrieves all participants for a specific battle with the given role (team or opponent).
func (s *Storage) loadParticipants(battleTime string, role string) ([]types.Player, error) {
	rows, err := s.DB.Query(`
		SELECT p.tag, p.name, bp.crowns, bp.startingTrophies, bp.trophyChange, bp.elixirLeaked
		FROM battle_participants bp
		JOIN players p ON p.tag = bp.player_tag
		WHERE bp.battleTime = ? AND bp.role = ?
	`, battleTime, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []types.Player

	for rows.Next() {
		var p types.Player
		var startingTrophies sql.NullInt64
		var trophyChange sql.NullInt32

		err := rows.Scan(
			&p.Tag,
			&p.Name,
			&p.Crowns,
			&startingTrophies,
			&trophyChange,
			&p.ElixirLeaked,
		)
		if err != nil {
			return nil, err
		}

		// Handle NULL values
		if startingTrophies.Valid {
			p.StartingTrophies = int32(startingTrophies.Int64)
		} else {
			p.StartingTrophies = 0
		}

		if trophyChange.Valid {
			p.TrophyChange = int32(trophyChange.Int32)
		} else {
			p.TrophyChange = 0
		}

		cards, err := s.loadDeck(p.Tag, battleTime)
		if err != nil {
			return nil, err
		}
		p.Cards = cards

		players = append(players, p)
	}

	return players, nil
}

// loadDeck retrieves all cards in a player's deck for a specific battle.
func (s *Storage) loadDeck(playerTag string, battleTime string) ([]types.Card, error) {
	rows, err := s.DB.Query(`
		SELECT c.id, c.name, c.maxLevel, c.rarity, c.elixirCost, bd.card_level
		FROM battle_decks bd
		JOIN cards c ON c.id = bd.card_id
		WHERE bd.battleTime = ? AND bd.player_tag = ?
		ORDER BY c.id
	`, battleTime, playerTag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []types.Card

	for rows.Next() {
		var c types.Card
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.MaxLevel,
			&c.Rarity,
			&c.ElixirCost,
			&c.Level,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	return cards, nil
}
