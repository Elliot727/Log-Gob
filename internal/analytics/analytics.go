// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"sort"

	"github.com/elliot727/log-gob/internal/storage"
	"github.com/elliot727/log-gob/internal/types"
)

// Compute builds the full Analytics struct by loading battles once
// and computing all stats from them. This is efficient for current scale (~200 battles).
func Compute(s *storage.Storage, myTag string, targetTrophies int) (Analytics, error) {
	var a Analytics

	// 1. Load all battles for the player (most recent first from storage)
	battles, err := s.GetBattlesForPlayer(myTag)
	if err != nil {
		return a, err
	}
	if len(battles) == 0 {
		return a, nil // empty but valid
	}

	// Reverse to chronological order (oldest → newest) for easier progression calculations
	// (streaks, trophy history, level ups, etc.)
	sort.Slice(battles, func(i, j int) bool {
		return battles[i].BattleTime < battles[j].BattleTime
	})

	// 2. Compute each section using the separate compute functions
	a.Overall = computeOverall(battles, myTag)
	a.Recent = computeRecent(battles, myTag)
	a.Arenas = computeArenas(battles, myTag)
	a.Projection = computeProjection(battles, targetTrophies, myTag)
	a.Elixir = computeElixir(battles, myTag)
	a.Crowns = computeCrowns(battles, myTag)
	a.Cards = computeCardImpact(battles, myTag)
	a.Losses = computeLossInsights(battles, myTag)
	a.Challenge = computeChallengeProof(battles, myTag)

	return a, nil
}

// Helper to find the player's participant record in a battle
func findMyParticipant(battle types.Battle, myTag string) *types.Player {
	for i := range battle.Team {
		if battle.Team[i].Tag == myTag {
			return &battle.Team[i]
		}
	}
	for i := range battle.Opponent {
		if battle.Opponent[i].Tag == myTag {
			return &battle.Opponent[i]
		}
	}
	return nil
}

// Helper to determine if the player won the battle
func isWin(me, opp int32) bool { return me > opp }
