// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"fmt"

	"github.com/elliot727/log-gob/internal/types"
)

// computeElixir analyzes elixir leak patterns.
func computeElixir(battles []types.Battle, myTag string) ElixirStats {
	var es ElixirStats
	var totalLeakWins, totalLeakLosses float64
	var winCount, lossCount int

	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}

		myLeak := me.ElixirLeaked
		won := isWin(me.Crowns, b.Opponent[0].Crowns)

		if won {
			totalLeakWins += myLeak
			winCount++
		} else {
			totalLeakLosses += myLeak
			lossCount++
			if myLeak > es.MaxLeakInLoss {
				es.MaxLeakInLoss = myLeak
			}
			if myLeak > 2.0 {
				es.HighLeakLosses++
			}
		}
	}

	if winCount > 0 {
		es.AvgLeakWins = totalLeakWins / float64(winCount)
	}
	if lossCount > 0 {
		es.AvgLeakLosses = totalLeakLosses / float64(lossCount)
	}

	if es.AvgLeakLosses > es.AvgLeakWins {
		diff := es.AvgLeakLosses - es.AvgLeakWins
		es.LeakImprovement = fmt.Sprintf("You leak %.1f more elixir in losses than in wins.", diff)
	} else {
		es.LeakImprovement = "Your elixir management is consistent across wins and losses."
	}

	return es
}
