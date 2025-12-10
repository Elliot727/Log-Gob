// Package analytics provides functionality for analyzing Clash Royale battle data.
package analytics

import (
	"time"

	"github.com/elliot727/log-gob/internal/types"
)

// computeRecent calculates stats for the last 10, 20, 50 battles and today.
func computeRecent(battles []types.Battle, myTag string) RecentForm {
	var rf RecentForm
	rf.Last10 = computeSession(battles, myTag, 10)
	rf.Last20 = computeSession(battles, myTag, 20)
	rf.Last50 = computeSession(battles, myTag, 50)
	rf.Today = computeTodaySession(battles, myTag)
	return rf
}

// computeSession helper for last N battles
func computeSession(battles []types.Battle, myTag string, n int) SessionStats {
	var s SessionStats
	if len(battles) < n {
		n = len(battles)
	}
	sessionBattles := battles[len(battles)-n:]

	for _, b := range sessionBattles {
		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}
		s.Battles++
		if isWin(me.Crowns, b.Opponent[0].Crowns) {
			s.Wins++
		}
		s.TrophyChange += int(me.TrophyChange)
	}

	if s.Battles > 0 {
		s.WinRate = float64(s.Wins) / float64(s.Battles) * 100
		s.AvgTrophyPerBattle = float64(s.TrophyChange) / float64(s.Battles)
	}
	return s
}

// computeTodaySession helper for today's battles
func computeTodaySession(battles []types.Battle, myTag string) SessionStats {
	var s SessionStats
	today := time.Now().UTC().Truncate(24 * time.Hour)

	for i := len(battles) - 1; i >= 0; i-- {
		b := battles[i]
		battleTime, err := time.Parse(time.RFC3339, b.BattleTime)
		if err != nil {
			continue // skip if time parsing fails
		}
		if battleTime.UTC().Before(today) {
			break // past today's battles
		}

		me := findMyParticipant(b, myTag)
		if me == nil {
			continue
		}
		s.Battles++
		if isWin(me.Crowns, b.Opponent[0].Crowns) {
			s.Wins++
		}
		s.TrophyChange += int(me.TrophyChange)
	}

	if s.Battles > 0 {
		s.WinRate = float64(s.Wins) / float64(s.Battles) * 100
		s.AvgTrophyPerBattle = float64(s.TrophyChange) / float64(s.Battles)
	}
	return s
}
