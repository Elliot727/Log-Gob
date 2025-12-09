// Package main provides the Bubble Tea TUI for the Clash Royale battle logger.
package main

import (
	"database/sql"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elliot727/log-gob/internal/storage"
	"github.com/elliot727/log-gob/internal/ui"
	_ "github.com/glebarez/go-sqlite"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get player tag from environment variables
	playerTag := os.Getenv("PLAYERTAG")
	if playerTag == "" {
		// Use a default tag if none is set, but warn the user
		playerTag = "PLY2Q2LL" // Default example tag
		log.Println("PLAYERTAG environment variable not set, using default tag. Please set PLAYERTAG in your .env file")
	}

	// Ensure player tag has proper format for API request
	if playerTag[0] != '#' {
		playerTag = "#" + playerTag
	}

	db, err := sql.Open("sqlite", "battles.db")
	if err != nil {
		log.Fatal("Failed to open db", err)
	}
	defer db.Close()

	s := storage.NewStorage(db)

	// Initialize storage (create tables if they don't exist)
	err = s.Init()
	if err != nil {
		log.Fatal("Failed to init storage", err)
	}

	p := tea.NewProgram(ui.InitialModel(s, playerTag))
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
