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

			// Determine the max number of cards to display
			// maxCards := min(len(battle.Team[0].Cards), len(battle.Opponent[0].Cards))

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

	s.WriteString("\n")
	s.WriteString(statusStyle.Render(m.status))
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Controls: [J/K] Navigate | [R] Refresh | [Q] Quit"))

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
