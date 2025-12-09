// Package types defines the data structures used throughout the application for Clash Royale data.
package types

// Arena represents a Clash Royale arena where battles take place.
type Arena struct {
	ID   int32  `json:"id"`   // The unique identifier for the arena
	Name string `json:"name"` // The name of the arena
}
