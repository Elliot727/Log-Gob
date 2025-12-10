// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"github.com/elliot727/log-gob/internal/types"
)

// computeCardImpact analyzes the impact of card levels on win rates.
func computeCardImpact(battles []types.Battle, myTag string) []CardImpact {
	// This is a complex analysis that requires tracking card levels over time.
	// For this mock, we will assume a static deck and analyze win rates per card level found.

	cardStats := make(map[string]*CardImpact)

	// Initialize with the deck from the most recent battle
	if len(battles) > 0 {
		latestBattle := battles[len(battles)-1]
		me := findMyParticipant(latestBattle, myTag)
		if me != nil {
			for _, card := range me.Cards {
				cardStats[card.Name] = &CardImpact{
					CardName:       card.Name,
					CurrentLevel:   int(card.Level),
					BattlesAtLevel: make(map[int]LevelPerformance),
				}
			}
		}
	}

	// Process all battles to populate stats
	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}
		won := isWin(me.Crowns, b.Opponent[0].Crowns)

		for _, card := range me.Cards {
			stat, exists := cardStats[card.Name]
			if !exists {
				// Card not in the final deck, skip for simplicity
				continue
			}

			// Update stats for the level this card was at during this battle
			levelPerf := stat.BattlesAtLevel[int(card.Level)]
			levelPerf.Battles++
			if won {
				levelPerf.Wins++
			}
			stat.BattlesAtLevel[int(card.Level)] = levelPerf
		}
	}

	var result []CardImpact
	for _, stat := range cardStats {
		// Finalize win rates
		for level, perf := range stat.BattlesAtLevel {
			if perf.Battles > 0 {
				perf.WinRate = float64(perf.Wins) / float64(perf.Battles) * 100
				stat.BattlesAtLevel[level] = perf
			}
		}
		// Compute stats since last upgrade
		stat.SinceLastUpgrade = calculateSinceLastUpgrade(battles, myTag, stat.CardName, stat.CurrentLevel)
		result = append(result, *stat)
	}

	return result
}

// calculateSinceLastUpgrade finds win rate since a card was upgraded to its current level.
func calculateSinceLastUpgrade(battles []types.Battle, myTag, cardName string, currentLevel int) struct {
	Battles int
	WinRate float64
} {
	var sinceUpgradeBattles, sinceUpgradeWins int
	var foundUpgradePoint bool

	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}

		cardInBattle, found := findCardInDeck(me.Cards, cardName)
		if !found {
			continue
		}

		// Start counting from the first battle with the current level
		if int(cardInBattle.Level) == currentLevel {
			foundUpgradePoint = true
		}

		if foundUpgradePoint && int(cardInBattle.Level) == currentLevel {
			sinceUpgradeBattles++
			if isWin(me.Crowns, b.Opponent[0].Crowns) {
				sinceUpgradeWins++
			}
		}
	}

	var winRate float64
	if sinceUpgradeBattles > 0 {
		winRate = float64(sinceUpgradeWins) / float64(sinceUpgradeBattles) * 100
	}

	return struct {
		Battles int
		WinRate float64
	}{
		Battles: sinceUpgradeBattles,
		WinRate: winRate,
	}
}

// findCardInDeck is a helper to find a specific card in a player's deck.
func findCardInDeck(cards []types.Card, cardName string) (types.Card, bool) {
	for _, c := range cards {
		if c.Name == cardName {
			return c, true
		}
	}
	return types.Card{}, false
}
