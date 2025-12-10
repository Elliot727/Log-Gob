// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"fmt"
	"math"
	"time"

	"github.com/elliot727/log-gob/internal/types"
)

// computeProjection estimates time to reach a trophy target.
func computeProjection(battles []types.Battle, targetTrophies int, myTag string) TrophyProjection {
	if len(battles) == 0 {
		return TrophyProjection{}
	}

	currentTrophies := 0
	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me != nil {
			currentTrophies += int(me.TrophyChange)
		}
	}

	if targetTrophies <= currentTrophies {
		return TrophyProjection{TargetTrophies: targetTrophies} // Already there
	}

	trophiesNeeded := targetTrophies - currentTrophies

	// Get win rates
	careerWR := float64(computeSession(battles, myTag, len(battles)).Wins) / float64(len(battles))
	last50WR := computeSession(battles, myTag, 50).WinRate / 100
	last20WR := computeSession(battles, myTag, 20).WinRate / 100

	// Estimate battles needed
	avgGainPerWin := 30.0   // Simplified assumption
	avgLossPerLoss := -25.0 // Simplified assumption

	calculateBattles := func(wr float64) int {
		if wr <= 0.5 { // Unlikely to climb
			return -1
		}
		// Expected gain per battle: (WR * Gain) + ((1-WR) * Loss)
		netGainPerBattle := (wr * avgGainPerWin) + ((1 - wr) * avgLossPerLoss)
		if netGainPerBattle <= 0 {
			return -1
		}
		return int(math.Ceil(float64(trophiesNeeded) / netGainPerBattle))
	}

	battlesRealistic := calculateBattles(last50WR)

	// Estimate days needed
	battlesPerDay := calculateBattlesPerDay(battles)
	daysNeeded := -1
	if battlesRealistic > 0 && battlesPerDay > 0 {
		daysNeeded = int(math.Ceil(float64(battlesRealistic) / battlesPerDay))
	}

	return TrophyProjection{
		TargetTrophies: targetTrophies,
		BattlesNeeded:  battlesRealistic,
		DaysNeeded:     daysNeeded,
		RealisticWR:    last50WR * 100,
		OptimisticWR:   last20WR * 100,
		PessimisticWR:  careerWR * 100,
		EstimatedDate:  formatDays(daysNeeded),
	}
}

// calculateBattlesPerDay determines avg battles per day from history.
func calculateBattlesPerDay(battles []types.Battle) float64 {
	if len(battles) < 2 {
		return float64(len(battles))
	}

	firstBattleTimeStr := battles[0].BattleTime
	lastBattleTimeStr := battles[len(battles)-1].BattleTime

	firstTime, err := time.Parse(time.RFC3339, firstBattleTimeStr)
	if err != nil {
		return 0 // Or handle error appropriately
	}
	lastTime, err := time.Parse(time.RFC3339, lastBattleTimeStr)
	if err != nil {
		return 0 // Or handle error appropriately
	}

	duration := lastTime.Sub(firstTime)
	days := duration.Hours() / 24.0

	if days < 1.0 {
		return float64(len(battles))
	}

	return float64(len(battles)) / days
}

// formatDays converts days into a human-readable string.
func formatDays(days int) string {
	if days < 0 {
		return "Never at this rate"
	}
	if days == 0 {
		return "Today"
	}
	if days == 1 {
		return "Tomorrow"
	}
	return fmt.Sprintf("In %d days", days)
}
