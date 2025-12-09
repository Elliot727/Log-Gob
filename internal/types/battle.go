// Package types defines the data structures used throughout the application for Clash Royale data.
package types

// Battle represents a single battle in Clash Royale, including the time, type, arena, game mode, and participants.
type Battle struct {
	BattleTime string   `json:"battleTime"` // The time when the battle occurred in ISO 8601 format
	BattleType string   `json:"type"`       // The type of battle (e.g., "PvP", "friendly", "tournament")
	Arena      Arena    `json:"arena"`      // The arena where the battle took place
	GameMode   GameMode `json:"gameMode"`   // The game mode of the battle
	Team       []Player `json:"team"`       // The players on the user's team
	Opponent   []Player `json:"opponent"`   // The players on the opposing team
}
