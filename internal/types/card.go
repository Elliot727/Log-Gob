// Package types defines the data structures used throughout the application for Clash Royale data.
package types

// Card represents a Clash Royale card in a player's deck.
type Card struct {
	ID         int32  `json:"id"`         // The unique identifier for the card
	Name       string `json:"name"`       // The name of the card
	Level      int32  `json:"level"`      // The current level of the card
	MaxLevel   uint32 `json:"maxLevel"`   // The maximum possible level for this card in battles
	Rarity     string `json:"rarity"`     // The rarity of the card (e.g., "Common", "Rare", "Epic", "Legendary")
	ElixirCost int32  `json:"elixirCost"` // The elixir cost to play this card
}
