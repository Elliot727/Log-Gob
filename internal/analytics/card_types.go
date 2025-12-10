// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

// CardImpact - Since deck unchanged, track level-up impact
type CardImpact struct {
	CardName         string
	CurrentLevel     int
	BattlesAtLevel   map[int]LevelPerformance // level -> stats
	SinceLastUpgrade struct {
		Battles int
		WinRate float64
	}
}
