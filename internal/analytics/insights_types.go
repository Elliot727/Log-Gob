// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

// LossInsights - Actionable loss patterns
type LossInsights struct {
	TotalLosses           int
	HighElixirLeakLosses  int
	OneCrownDefenseLosses int      // 0-1 or 1-2/3 losses
	CommonNotes           []string // e.g., "70% of losses leaked >2.0"
	RecentLossStreak      int
}

// ChallengeProof - Your pride module
type ChallengeProof struct {
	StartArena         string
	StartTrophies      int
	BattlesSinceStart  int
	TrophiesGained     int
	WinRateSinceStart  float64
	UniqueCardsUsed    int    // should always be 8
	DeckUnchangedSince string // date or arena
	MilestoneMessage   string // e.g., "Arena 6 → 17: +4700 trophies, no changes"
}
