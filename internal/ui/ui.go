// Package ui provides a Bubble Tea TUI for the Clash Royale battle logger.
package ui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elliot727/log-gob/internal/storage"
	"github.com/elliot727/log-gob/internal/types"
)

var (
	// Styling for the UI
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")). // Gold
			Background(lipgloss.Color("#8B4513")). // Saddle brown
			Padding(0, 1).
			Bold(true).
			MarginBottom(1)

	battleHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#40E0D0")). // Turquoise
				Bold(true).
				MarginBottom(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA07A")) // Light salmon

	teamStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98FB98")). // Pale green
			Bold(true)

	opponentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6347")). // Tomato
			Bold(true)

	playerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F0E68C")) // Khaki

	crownStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")). // Gold
			Bold(true)

	trophyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4169E1")). // Royal blue
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB")). // Sky blue
			Italic(true).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D3D3D3")). // Light gray
			Italic(true).
			MarginTop(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF69B4")). // Hot pink
			Bold(true)

	cardStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDA0DD")) // Plum
)

type model struct {
	storage     *storage.Storage
	playerTag   string
	battles     []types.Battle
	currentIdx  int
	status      string
	initialized bool
	showStats   bool // Toggle between detail view and stats view
}

type fetchMsg struct {
	battles []types.Battle
	err     error
}

type statusMsg struct {
	text string
}

func InitialModel(s *storage.Storage, playerTag string) model {
	return model{
		storage:     s,
		playerTag:   playerTag,
		battles:     []types.Battle{},
		currentIdx:  0,
		status:      "Loading battles...",
		initialized: false,
		showStats:   false,
	}
}

func (m model) Init() tea.Cmd {
	return fetchBattles(m.storage, m.playerTag)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			if len(m.battles) > 0 {
				m.currentIdx = (m.currentIdx + 1) % len(m.battles)
			}
		case "k", "up":
			if len(m.battles) > 0 {
				if m.currentIdx == 0 {
					m.currentIdx = len(m.battles) - 1
				} else {
					m.currentIdx--
				}
			}
		case "r", "R":
			return m, fetchBattles(m.storage, m.playerTag)
		case "s", "S":
			// Toggle between stats view and detail view
			m.showStats = !m.showStats
			if m.showStats {
				m.status = "Switched to stats view"
			} else {
				m.status = "Switched to detail view"
			}
		}

	case fetchMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Error fetching battles: %v", msg.err)
		} else {
			m.battles = msg.battles
			if len(m.battles) > 0 {
				m.status = fmt.Sprintf("Fetched %d battles for player %s", len(m.battles), m.playerTag)
			} else {
				m.status = fmt.Sprintf("No battles found for player %s", m.playerTag)
			}
		}
		m.initialized = true

	case statusMsg:
		m.status = msg.text
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("=== Clash Royale Battle Logger ==="))
	s.WriteString("\n\n")

	if !m.initialized {
		s.WriteString(statusStyle.Render(m.status))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press 'R' to refresh, 'Q' to quit"))
		return s.String()
	}

	if len(m.battles) == 0 {
		s.WriteString(statusStyle.Render(m.status))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press 'R' to refresh, 'Q' to quit"))
		return s.String()
	}

	if m.showStats {
		// Stats view
		stats := CalculateStats(m.battles)

		s.WriteString(battleHeaderStyle.Render("BATTLE STATISTICS"))
		s.WriteString("\n\n")

		// Basic stats
		s.WriteString(headerStyle.Bold(true).Render("Overall Performance"))
		s.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("─", 40)))

		winRateStyle := infoStyle
		if stats.WinRate >= 60 {
			winRateStyle = teamStyle // Green for good win rate
		} else if stats.WinRate >= 50 {
			winRateStyle = infoStyle // Light salmon for OK win rate
		} else {
			winRateStyle = opponentStyle // Red for poor win rate
		}

		s.WriteString(fmt.Sprintf("Win Rate:       %-6s %s %.1f%%\n",
			winRateStyle.Bold(true).Render(fmt.Sprintf("%.1f%%", stats.WinRate)),
			makeProgressBar(stats.WinRate, 15),
			stats.WinRate))
		s.WriteString(fmt.Sprintf("Wins:            %s\n",
			teamStyle.Render(fmt.Sprintf("%d", stats.TotalWins))))
		s.WriteString(fmt.Sprintf("Losses:          %s\n",
			opponentStyle.Render(fmt.Sprintf("%d", stats.TotalLosses))))
		s.WriteString(fmt.Sprintf("Total Battles:   %s\n",
			infoStyle.Render(fmt.Sprintf("%d", stats.TotalBattles))))
		s.WriteString("\n")

		// Performance metrics
		s.WriteString(headerStyle.Bold(true).Render("Performance Metrics"))
		s.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("─", 40)))
		s.WriteString(fmt.Sprintf("Avg Crowns Won:  %s\n",
			teamStyle.Render(fmt.Sprintf("%.2f", stats.AvgCrownsWon))))
		s.WriteString(fmt.Sprintf("Avg Crowns Lost: %s\n",
			opponentStyle.Render(fmt.Sprintf("%.2f", stats.AvgCrownsLost))))
		s.WriteString(fmt.Sprintf("Avg Trophy Gain: %s\n",
			teamStyle.Render(fmt.Sprintf("%.1f", stats.AvgTrophyGain))))
		s.WriteString(fmt.Sprintf("Avg Trophy Loss: %s\n",
			opponentStyle.Render(fmt.Sprintf("%.1f", stats.AvgTrophyLoss))))
		s.WriteString("\n")

		// Arena stats (UPDATED)
		s.WriteString(headerStyle.Bold(true).Render("Arena Performance (Win Rate)"))
		s.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("─", 65)))
		s.WriteString(fmt.Sprintf("%-20s %-12s %-8s %s\n", "Arena", "Win Rate", "%", "Record (W-L)"))

		for arena, stat := range stats.ArenaStats {
			winRate := 0.0
			if stat.Total > 0 {
				winRate = float64(stat.Wins) / float64(stat.Total) * 100
			}

			// Color code the win rate
			var style lipgloss.Style
			if winRate >= 60 {
				style = teamStyle
			} else if winRate >= 50 {
				style = infoStyle
			} else {
				style = opponentStyle
			}

			record := fmt.Sprintf("%d-%d", stat.Wins, stat.Losses)

			s.WriteString(fmt.Sprintf("%-20s %s %-8s %s\n",
				infoStyle.Render(arena),
				makeProgressBar(winRate, 10),
				style.Render(fmt.Sprintf("%.1f%%", winRate)),
				infoStyle.Render(record)))
		}
	} else {
		// Show current battle in detail
		if len(m.battles) > 0 && m.currentIdx < len(m.battles) {
			battle := m.battles[m.currentIdx]

			// Format time to be human-readable: 2025-10-11 08:23
			timeDisplay := battle.BattleTime
			if len(timeDisplay) >= 16 {
				// Extract YYYY-MM-DD and HH:MM from "20251011T082308.000Z"
				if len(timeDisplay) >= 19 {
					year := timeDisplay[0:4]
					month := timeDisplay[4:6]
					day := timeDisplay[6:8]
					hour := timeDisplay[9:11]
					minute := timeDisplay[11:13]
					timeDisplay = fmt.Sprintf("%s-%s-%s %s:%s", year, month, day, hour, minute)
				}
			}

			// Calculate battle result
			result := "Unknown"
			var teamCrowns, opponentCrowns int32 = 0, 0

			if len(battle.Team) > 0 {
				teamCrowns = battle.Team[0].Crowns
			}
			if len(battle.Opponent) > 0 {
				opponentCrowns = battle.Opponent[0].Crowns
			}

			if teamCrowns > opponentCrowns {
				result = "Victory"
			} else if teamCrowns < opponentCrowns {
				result = "Loss"
			} else {
				result = "Draw"
			}

			// Determine result color
			var resultStyle lipgloss.Style
			if result == "Victory" {
				resultStyle = teamStyle
			} else if result == "Loss" {
				resultStyle = opponentStyle
			} else {
				resultStyle = headerStyle
			}

			// Battle header
			s.WriteString(battleHeaderStyle.Render(fmt.Sprintf("Battle %d of %d", m.currentIdx+1, len(m.battles))))
			s.WriteString("\n")
			s.WriteString(resultStyle.Bold(true).Render(fmt.Sprintf("Result: %s", result)))
			s.WriteString("\n")
			s.WriteString(infoStyle.Render(fmt.Sprintf("Time: %s", timeDisplay)))
			s.WriteString("\n")
			s.WriteString(infoStyle.Render(fmt.Sprintf("Type: %s", battle.BattleType)))
			s.WriteString("\n")
			s.WriteString(infoStyle.Render(fmt.Sprintf("Arena: %s", battle.Arena.Name)))
			s.WriteString("\n")
			s.WriteString(infoStyle.Render(fmt.Sprintf("Game Mode: %s", battle.GameMode.Name)))
			s.WriteString("\n\n")

			// Show team
			s.WriteString(teamStyle.Render("Team:"))
			s.WriteString("\n")
			for _, player := range battle.Team {
				playerInfo := fmt.Sprintf("  - %s (%s) | Crowns: %s | Trophies: %s",
					playerStyle.Render(player.Name),
					infoStyle.Render(player.Tag),
					crownStyle.Render(fmt.Sprintf("%d", player.Crowns)),
					trophyStyle.Render(fmt.Sprintf("+%d", player.TrophyChange)))
				s.WriteString(playerInfo)
				s.WriteString("\n")
			}

			// Show opponent
			s.WriteString("\n")
			s.WriteString(opponentStyle.Render("Opponent:"))
			s.WriteString("\n")
			for _, player := range battle.Opponent {
				playerInfo := fmt.Sprintf("  - %s (%s) | Crowns: %s | Trophies: %s",
					playerStyle.Render(player.Name),
					infoStyle.Render(player.Tag),
					crownStyle.Render(fmt.Sprintf("%d", player.Crowns)),
					trophyStyle.Render(fmt.Sprintf("+%d", player.TrophyChange)))
				s.WriteString(playerInfo)
				s.WriteString("\n")
			}

			// Show cards for team vs opponent side by side (if available)
			if len(battle.Team) > 0 && len(battle.Opponent) > 0 && len(battle.Team[0].Cards) > 0 && len(battle.Opponent[0].Cards) > 0 {
				s.WriteString("\n")
				s.WriteString(headerStyle.Render("Deck Matchup:"))
				s.WriteString("\n")

				// Display cards side by side
				for i := range battle.Team[0].Cards {
					var teamCardName, teamCardLevel, oppCardName, oppCardLevel string
					if i < len(battle.Team[0].Cards) {
						teamCardName = battle.Team[0].Cards[i].Name
						teamCardLevel = fmt.Sprintf("Lvl %d", battle.Team[0].Cards[i].Level)
					} else {
						teamCardName = ""
						teamCardLevel = ""
					}

					if i < len(battle.Opponent[0].Cards) {
						oppCardName = battle.Opponent[0].Cards[i].Name
						oppCardLevel = fmt.Sprintf("Lvl %d", battle.Opponent[0].Cards[i].Level)
					} else {
						oppCardName = ""
						oppCardLevel = ""
					}

					// Format the side-by-side display with proper alignment
					cardInfo := fmt.Sprintf("%-20s %-10s │ %-20s %-10s\n",
						cardStyle.Render(teamCardName),
						cardStyle.Render(teamCardLevel),
						opponentStyle.Render(oppCardName),
						opponentStyle.Render(oppCardLevel))
					s.WriteString(cardInfo)
				}
			} else {
				// Fallback to individual display if side-by-side isn't possible
				if len(battle.Team) > 0 && len(battle.Team[0].Cards) > 0 {
					s.WriteString("\n")
					s.WriteString(headerStyle.Render("Team Cards:"))
					s.WriteString("\n")
					for i, card := range battle.Team[0].Cards {
						if i < 8 { // Limit to first 8 cards to keep it readable
							cardInfo := fmt.Sprintf("  • %s (Lvl %d)",
								cardStyle.Render(card.Name),
								card.Level)
							s.WriteString(cardInfo)
							s.WriteString("\n")
						}
					}
				}

				if len(battle.Opponent) > 0 && len(battle.Opponent[0].Cards) > 0 {
					s.WriteString("\n")
					s.WriteString(headerStyle.Render("Opponent Cards:"))
					s.WriteString("\n")
					for i, card := range battle.Opponent[0].Cards {
						if i < 8 { // Limit to first 8 cards to keep it readable
							cardInfo := fmt.Sprintf("  • %s (Lvl %d)",
								cardStyle.Render(card.Name),
								card.Level)
							s.WriteString(cardInfo)
							s.WriteString("\n")
						}
					}
				}
			}

			// Show support cards for the first team player (if available)
			if len(battle.Team) > 0 && len(battle.Team[0].SupportCards) > 0 {
				s.WriteString("\n")
				s.WriteString(headerStyle.Render("Team Support Cards:"))
				s.WriteString("\n")
				for i, card := range battle.Team[0].SupportCards {
					if i < 8 { // Limit to first 8 cards to keep it readable
						cardInfo := fmt.Sprintf("  • %s (Lvl %d)",
							cardStyle.Render(card.Name),
							card.Level)
						s.WriteString(cardInfo)
						s.WriteString("\n")
					}
				}
			}

			// Show support cards for the first opponent player (if available)
			if len(battle.Opponent) > 0 && len(battle.Opponent[0].SupportCards) > 0 {
				s.WriteString("\n")
				s.WriteString(headerStyle.Render("Opponent Support Cards:"))
				s.WriteString("\n")
				for i, card := range battle.Opponent[0].SupportCards {
					if i < 8 { // Limit to first 8 cards to keep it readable
						cardInfo := fmt.Sprintf("  • %s (Lvl %d)",
							cardStyle.Render(card.Name),
							card.Level)
						s.WriteString(cardInfo)
						s.WriteString("\n")
					}
				}
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(statusStyle.Render(m.status))
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Controls: [J/K] Navigate | [R] Refresh | [S] Stats | [Q] Quit"))

	return s.String()
}

func fetchBattles(s *storage.Storage, playerTag string) tea.Cmd {
	return func() tea.Msg {
		battles, err := s.GetBattlesForPlayer(playerTag)
		if err != nil {
			log.Printf("Error fetching battles: %v", err)
			return fetchMsg{battles: nil, err: err}
		}

		// Reverse the order to show most recent first
		for i, j := 0, len(battles)-1; i < j; i, j = i+1, j-1 {
			battles[i], battles[j] = battles[j], battles[i]
		}

		return fetchMsg{battles: battles, err: nil}
	}
}

// UpdateStatus updates the status message
func UpdateStatus(text string) tea.Cmd {
	return func() tea.Msg {
		return statusMsg{text: text}
	}
}

// Stats structure to hold calculated battle statistics
type Stats struct {
	TotalBattles  int
	TotalWins     int
	TotalLosses   int
	WinRate       float64
	AvgCrownsWon  float64
	AvgCrownsLost float64
	AvgTrophyGain float64
	AvgTrophyLoss float64
	ArenaStats    map[string]*ArenaPerformance // Changed to store complex struct
}

// New struct to hold detailed arena stats
type ArenaPerformance struct {
	Wins   int
	Losses int
	Total  int
}

// CalculateStats calculates battle statistics from the battles data
func CalculateStats(battles []types.Battle) Stats {
	stats := Stats{
		ArenaStats: make(map[string]*ArenaPerformance),
	}

	totalCrownsWon := 0
	totalCrownsLost := 0
	totalTrophyGain := 0
	totalTrophyLoss := 0

	for _, battle := range battles {
		stats.TotalBattles++

		// Initialize arena stat if not exists
		if _, exists := stats.ArenaStats[battle.Arena.Name]; !exists {
			stats.ArenaStats[battle.Arena.Name] = &ArenaPerformance{}
		}
		arenaStat := stats.ArenaStats[battle.Arena.Name]
		arenaStat.Total++

		// Count wins/losses/draws based on crowns
		var teamCrowns, opponentCrowns int32 = 0, 0
		if len(battle.Team) > 0 {
			teamCrowns = battle.Team[0].Crowns
		}
		if len(battle.Opponent) > 0 {
			opponentCrowns = battle.Opponent[0].Crowns
		}

		if teamCrowns > opponentCrowns {
			stats.TotalWins++
			arenaStat.Wins++ // Update Arena Stat
			totalCrownsWon += int(teamCrowns)
			if battle.Team[0].TrophyChange > 0 {
				totalTrophyGain += int(battle.Team[0].TrophyChange)
			}
		} else if teamCrowns < opponentCrowns {
			stats.TotalLosses++
			arenaStat.Losses++ // Update Arena Stat
			totalCrownsLost += int(teamCrowns)
			if battle.Team[0].TrophyChange < 0 {
				totalTrophyLoss += int(battle.Team[0].TrophyChange)
			}
		} else {

		}
	}

	if stats.TotalBattles > 0 {
		stats.WinRate = float64(stats.TotalWins) / float64(stats.TotalBattles) * 100
	}

	if stats.TotalWins > 0 {
		stats.AvgCrownsWon = float64(totalCrownsWon) / float64(stats.TotalWins)
		stats.AvgTrophyGain = float64(totalTrophyGain) / float64(stats.TotalWins)
	}

	if stats.TotalLosses > 0 {
		stats.AvgCrownsLost = float64(totalCrownsLost) / float64(stats.TotalLosses)
		stats.AvgTrophyLoss = float64(totalTrophyLoss) / float64(stats.TotalLosses)
	}

	return stats
}

// Helper function to create a progress bar
func makeProgressBar(percentage float64, width int) string {
	filled := int(percentage / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	if empty < 0 {
		empty = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return bar
}
