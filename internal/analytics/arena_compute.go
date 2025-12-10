// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import "github.com/elliot727/log-gob/internal/types"

// computeArenas analyzes performance per arena.
func computeArenas(battles []types.Battle, myTag string) []ArenaStats {
	arenaMap := make(map[string]*ArenaStats)
	var arenaOrder []string

	for _, b := range battles {
		me := findMyParticipant(b, myTag)
		if me == nil || b.Arena.Name == "" {
			continue
		}

		// Get or create arena stats
		stats, exists := arenaMap[b.Arena.Name]
		if !exists {
			stats = &ArenaStats{ArenaName: b.Arena.Name}
			arenaMap[b.Arena.Name] = stats
			arenaOrder = append(arenaOrder, b.Arena.Name)
		}

		stats.Battles++
		if isWin(me.Crowns, b.Opponent[0].Crowns) {
			stats.Wins++
		}
		stats.AvgTrophyGain += float64(me.TrophyChange) // Will be averaged later
		if me.Crowns == 3 {
			stats.ThreeCrownRate++ // Will be averaged later
		}
	}

	// Finalize calculations and create the result slice
	var result []ArenaStats
	currentArenaName := ""
	if len(battles) > 0 {
		currentArenaName = battles[len(battles)-1].Arena.Name
	}

	for _, arenaName := range arenaOrder {
		stats := arenaMap[arenaName]
		if stats.Battles > 0 {
			stats.WinRate = float64(stats.Wins) / float64(stats.Battles) * 100
			stats.AvgTrophyGain /= float64(stats.Battles)
			stats.ThreeCrownRate = (stats.ThreeCrownRate / float64(stats.Wins)) * 100
		}
		stats.IsCurrent = (arenaName == currentArenaName)
		result = append(result, *stats)
	}
	return result
}
