// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"github.com/elliot727/log-gob/internal/types"
)

// computeOverall - career summary
func computeOverall(battles []types.Battle, myTag string) OverallStats {
	os := OverallStats{
		PeakTrophies:    0,
		CurrentTrophies: 0,
	}

	runningTrophies := 0
	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}

		os.TotalBattles++
		if isWin(me.Crowns, b.Opponent[0].Crowns) {
			os.Wins++
			if me.Crowns == 3 {
				os.ThreeCrownWins++
			}
		} else {
			os.Losses++
		}

		runningTrophies += int(me.TrophyChange)
		if runningTrophies > os.PeakTrophies {
			os.PeakTrophies = runningTrophies
		}
	}

	os.CurrentTrophies = runningTrophies
	os.TotalTrophyGain = runningTrophies

	if os.TotalBattles > 0 {
		os.WinRate = float64(os.Wins) / float64(os.TotalBattles) * 100
		os.ThreeCrownRate = float64(os.ThreeCrownWins) / float64(os.Wins) * 100
	}

	// Streaks
	os.CurrentStreak, os.LongestWinStreak = computeStreaks(battles, myTag)

	return os
}

// computeStreaks helper used by Overall
func computeStreaks(battles []types.Battle, myTag string) (current int, longest int) {
	if len(battles) == 0 {
		return 0, 0
	}

	currentStreak := 0
	maxWinStreak := 0
	currentWinStreak := 0

	// Walk backwards from latest battle
	for i := len(battles) - 1; i >= 0; i-- {
		b := battles[i]
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}

		won := isWin(me.Crowns, b.Opponent[0].Crowns)

		if won {
			// Current streak is positive (winning)
			if currentStreak >= 0 {
				currentStreak++
			} else {
				currentStreak = 1 // Switch from losing to winning
			}

			currentWinStreak++
			if currentWinStreak > maxWinStreak {
				maxWinStreak = currentWinStreak
			}
		} else {
			// Current streak is negative (losing)
			if currentStreak <= 0 {
				currentStreak--
			} else {
				currentStreak = -1 // Switch from winning to losing
			}

			currentWinStreak = 0
		}
	}

	return currentStreak, maxWinStreak
}
