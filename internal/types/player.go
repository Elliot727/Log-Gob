// Package types defines the data structures used throughout the application for Clash Royale data.
package types

// Player represents a player in a Clash Royale battle, including their tag, name, performance stats, and cards.
type Player struct {
	Tag              string  `json:"tag"`              // The player's unique tag identifier
	Name             string  `json:"name"`             // The player's name
	StartingTrophies int32   `json:"startingTrophies"` // Trophies the player had at the start of the battle
	TrophyChange     int32   `json:"trophyChange"`     // Change in trophies after the battle
	Crowns           int32   `json:"crowns"`           // Number of crowns earned in the battle
	Cards            []Card  `json:"cards"`            // The cards in the player's deck
	ElixirLeaked     float64 `json:"elixirLeaked"`     // Amount of elixir leaked to the opponent (0-1 scale)
	SupportCards     []Card  `json:"supportCards"`     // Support cards in the player's deck (if any)
}
