// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

// ElixirStats - Leak analysis
type ElixirStats struct {
	AvgLeakWins     float64
	AvgLeakLosses   float64
	MaxLeakInLoss   float64
	HighLeakLosses  int    // count of losses with leak > 2.0
	LeakImprovement string // e.g., "Losses leak 1.4 more than wins"
}

// CrownStats - Aggression/defense breakdown
type CrownStats struct {
	AvgCrownsTaken    float64
	AvgCrownsConceded float64
	WinTypes          map[string]int // "3-0": 120, "2-1": 80, etc.
	LossTypes         map[string]int // "0-1": 50, "0-3": 20
	AggressionScore   float64        // e.g., 3-crown rate normalized
}
