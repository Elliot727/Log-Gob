// Package types defines the data structures used throughout the application for Clash Royale data.
package types

// GameMode represents a Clash Royale game mode.
type GameMode struct {
	ID   int32  `json:"id"`   // The unique identifier for the game mode
	Name string `json:"name"` // The name of the game mode
}
