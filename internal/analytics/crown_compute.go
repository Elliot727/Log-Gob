// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"fmt"

	"github.com/elliot727/log-gob/internal/types"
)

// computeCrowns analyzes crown-related stats.
func computeCrowns(battles []types.Battle, myTag string) CrownStats {
	var cs CrownStats
	cs.WinTypes = make(map[string]int)
	cs.LossTypes = make(map[string]int)

	var totalCrownsTaken, totalCrownsConceded int

	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		opponent := b.Opponent[0]
		if me == nil {
			continue
		}

		totalCrownsTaken += int(me.Crowns)
		totalCrownsConceded += int(opponent.Crowns)

		outcome := fmt.Sprintf("%d-%d", me.Crowns, opponent.Crowns)
		if isWin(me.Crowns, opponent.Crowns) {
			cs.WinTypes[outcome]++
		} else if me.Crowns < opponent.Crowns {
			cs.LossTypes[outcome]++
		}
		// Ties are ignored for win/loss types
	}

	if len(battles) > 0 {
		cs.AvgCrownsTaken = float64(totalCrownsTaken) / float64(len(battles))
		cs.AvgCrownsConceded = float64(totalCrownsConceded) / float64(len(battles))
	}

	threeCrownWins := cs.WinTypes["3-0"] + cs.WinTypes["3-1"] + cs.WinTypes["3-2"]
	totalWins, _ := countWinsLosses(battles, myTag)
	if totalWins > 0 {
		// Aggression score: % of wins that are 3-crowns
		cs.AggressionScore = float64(threeCrownWins) / float64(totalWins) * 100
	}

	return cs
}

// countWinsLosses is a helper to get total wins and losses.
func countWinsLosses(battles []types.Battle, myTag string) (wins, losses int) {
	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}
		if isWin(me.Crowns, b.Opponent[0].Crowns) {
			wins++
		} else if me.Crowns < b.Opponent[0].Crowns {
			losses++
		}
	}
	return
}
