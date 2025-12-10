// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

// ArenaStats - Performance per arena
type ArenaStats struct {
	ArenaName      string
	Battles        int
	Wins           int
	WinRate        float64
	AvgTrophyGain  float64
	ThreeCrownRate float64
	IsCurrent      bool // highlight current arena
}

// TrophyProjection - Forward-looking estimates
type TrophyProjection struct {
	TargetTrophies int     // e.g., 7000
	BattlesNeeded  int     // at current/recent WR
	DaysNeeded     int     // assuming avg battles/day
	RealisticWR    float64 // e.g., last 50 WR
	OptimisticWR   float64 // e.g., last 20 WR
	PessimisticWR  float64 // e.g., career WR
	EstimatedDate  string  // formatted "in 5-7 days"
}
