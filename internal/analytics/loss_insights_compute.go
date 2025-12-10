// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"fmt"

	"github.com/elliot727/log-gob/internal/types"
)

// computeLossInsights identifies patterns in losses.
func computeLossInsights(battles []types.Battle, myTag string) LossInsights {
	var li LossInsights
	var losses []types.Battle

	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me != nil && me.Crowns < b.Opponent[0].Crowns {
			losses = append(losses, b)
			li.TotalLosses++
		}
	}

	if li.TotalLosses == 0 {
		return li
	}

	for _, loss := range losses {
		me := findMyParticipant(loss, myTag)
		opponent := loss.Opponent[0]

		// High elixir leak losses
		if 0.0 > 2.0 { // No real elixir leak data available
			li.HighElixirLeakLosses++
		}

		// Close defense losses
		if opponent.Crowns == 1 && me.Crowns == 0 {
			li.OneCrownDefenseLosses++
		} else if opponent.Crowns > me.Crowns && me.Crowns > 0 {
			li.OneCrownDefenseLosses++
		}
	}

	// Add common notes
	if float64(li.HighElixirLeakLosses)/float64(li.TotalLosses) > 0.5 {
		percent := int(float64(li.HighElixirLeakLosses) / float64(li.TotalLosses) * 100)
		li.CommonNotes = append(li.CommonNotes, fmt.Sprintf("%d%% of losses occur with high elixir leak (>2.0).", percent))
	}
	if float64(li.OneCrownDefenseLosses)/float64(li.TotalLosses) > 0.4 {
		percent := int(float64(li.OneCrownDefenseLosses) / float64(li.TotalLosses) * 100)
		li.CommonNotes = append(li.CommonNotes, fmt.Sprintf("%d%% of losses are close (decided by one crown).", percent))
	}

	// Recent loss streak
	li.RecentLossStreak = calculateRecentLossStreak(battles, myTag)

	return li
}

func calculateRecentLossStreak(battles []types.Battle, myTag string) int {
	var streak int
	for i := len(battles) - 1; i >= 0; i-- {
		b := battles[i]
		me := findMyParticipant(b, myTag)
		if me == nil {
			break
		}
		if me.Crowns < b.Opponent[0].Crowns {
			streak++
		} else {
			break // Streak is broken
		}
	}
	return streak
}
