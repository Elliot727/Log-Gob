// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

// OverallStats - Top-level summary shown on main dashboard
type OverallStats struct {
	TotalBattles     int
	Wins             int
	Losses           int
	WinRate          float64 // percentage
	ThreeCrownWins   int
	ThreeCrownRate   float64
	CurrentStreak    int // positive = wins, negative = losses
	LongestWinStreak int
	TotalTrophyGain  int
	CurrentTrophies  int // latest known
	PeakTrophies     int
}

// SessionStats holds statistics for a specific session or time period
type SessionStats struct {
	Battles            int
	Wins               int
	WinRate            float64
	TrophyChange       int
	AvgTrophyPerBattle float64
}

// RecentForm - Last N battles performance
type RecentForm struct {
	Last10 SessionStats
	Last20 SessionStats
	Last50 SessionStats
	Today  SessionStats // or last session
}
