// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

// Analytics represents the complete set of analytics computed for a player.
type Analytics struct {
	Overall    OverallStats
	Recent     RecentForm
	Arenas     []ArenaStats // ordered by progress
	Projection TrophyProjection
	Elixir     ElixirStats
	Crowns     CrownStats
	Cards      []CardImpact // one per card
	Losses     LossInsights
	Challenge  ChallengeProof
}

// LevelPerformance holds performance statistics for a specific card level
type LevelPerformance struct {
	Battles int
	Wins    int
	WinRate float64
}
