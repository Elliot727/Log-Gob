// Package main provides the Bubble Tea TUI for the Clash Royale battle logger.
package main

import (
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elliot727/log-gob/internal/config"
	"github.com/elliot727/log-gob/internal/storage"
	"github.com/elliot727/log-gob/internal/ui"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if cfg.PlayerTag == "" {
		log.Fatal("PLAYERTAG environment variable not set. Please set PLAYERTAG in your .env file with your Clash Royale player tag (e.g., #ABC123)")
	}

	// Validate player tag format - should start with # after our potential addition and have valid characters
	if len(cfg.PlayerTag) < 2 {
		log.Fatalf("Invalid player tag format: %s. Player tag should not be empty", cfg.PlayerTag)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	s := storage.NewStorage(db)

	// Initialize storage (create tables if they don't exist)
	err = s.Init()
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	p := tea.NewProgram(ui.InitialModel(s, cfg.PlayerTag))
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
