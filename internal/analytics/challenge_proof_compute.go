// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"fmt"
	"time"

	"github.com/elliot727/log-gob/internal/types"
)

// computeChallengeProof creates a summary of the player's journey.
func computeChallengeProof(battles []types.Battle, myTag string) ChallengeProof {
	if len(battles) == 0 {
		return ChallengeProof{}
	}

	firstBattle := battles[0]
	lastBattle := battles[len(battles)-1]

	meFirst := findMyParticipant(firstBattle, myTag)
	meLast := findMyParticipant(lastBattle, myTag)

	if meFirst == nil || meLast == nil {
		return ChallengeProof{} // Not enough data
	}

	startTrophies := meFirst.StartingTrophies

	// Calculate current trophies from starting trophies + trophy changes
	// Use the starting trophies from the first battle as base, and final trophy change as indicator of current state
	currentTrophies := meLast.StartingTrophies + meLast.TrophyChange
	totalTrophiesGained := currentTrophies - startTrophies

	wins, _ := countWinsLosses(battles, myTag)
	winRate := float64(wins) / float64(len(battles)) * 100

	// Check for deck changes (simple check: compare first and last battle decks)
	deckUnchanged := areDecksSame(meFirst.Cards, meLast.Cards)
	deckUnchangedSince := firstBattle.BattleTime
	if !deckUnchanged {
		deckUnchangedSince = "Deck has changed" // Or find the actual last change
	} else {
		t, _ := time.Parse(time.RFC3339, deckUnchangedSince)
		deckUnchangedSince = t.Format("Jan 2, 2006")
	}

	return ChallengeProof{
		StartArena:         firstBattle.Arena.Name,
		StartTrophies:      int(startTrophies),
		BattlesSinceStart:  len(battles),
		TrophiesGained:     int(totalTrophiesGained),
		WinRateSinceStart:  winRate,
		UniqueCardsUsed:    countUniqueCards(battles, myTag),
		DeckUnchangedSince: deckUnchangedSince,
		MilestoneMessage:   fmt.Sprintf("From %s to %d trophies: a %+d journey.", firstBattle.Arena.Name, currentTrophies, totalTrophiesGained),
	}
}

// areDecksSame checks if two decks (slices of cards) are identical.
func areDecksSame(deck1, deck2 []types.Card) bool {
	if len(deck1) != len(deck2) {
		return false
	}
	cardNames1 := make(map[string]bool)
	for _, card := range deck1 {
		cardNames1[card.Name] = true
	}
	for _, card := range deck2 {
		if !cardNames1[card.Name] {
			return false
		}
	}
	return true
}

// countUniqueCards counts how many unique cards were used across all battles.
func countUniqueCards(battles []types.Battle, myTag string) int {
	uniqueCards := make(map[string]bool)
	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}
		for _, card := range me.Cards {
			uniqueCards[card.Name] = true
		}
	}
	return len(uniqueCards)
}
