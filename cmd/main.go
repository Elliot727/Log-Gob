// Package main provides the entry point for the Clash Royale battle logger application.
// It connects to the Clash Royale API, fetches battle logs for a specific player,
// and stores the battle data in a SQLite database.
package main

import (
	"database/sql"
	"log"
	"net/url"
	"os"

	"github.com/elliot727/log-gob/internal/api"
	"github.com/elliot727/log-gob/internal/storage"
	"github.com/elliot727/log-gob/internal/types"
	_ "github.com/glebarez/go-sqlite"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Open SQLite database connection
	db, err := sql.Open("sqlite", "battles.db")
	if err != nil {
		log.Fatal("Failed to open db", err)
	}
	defer db.Close()

	// Initialize storage with database connection
	s := storage.NewStorage(db)

	// Create required database tables if they don't exist
	err = s.Init()
	if err != nil {
		log.Fatal("Failed to init storage", err)
	}

	// Get API key from environment variables
	apiKey := os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal("APIKEY environment variable not set")
	} else {
		// Create API client with Clash Royale API base URL
		apiClient := api.New("https://api.clashroyale.com", apiKey)

		// Get player tag from environment variables
		playerTag := os.Getenv("PLAYERTAG")
		if playerTag == "" {
			log.Fatal("PLAYERTAG environment variable not set. Please set PLAYERTAG in your .env file with your Clash Royale player tag (e.g., #ABC123)")
		}

		// Validate player tag format - should start with # and have valid characters
		if len(playerTag) < 2 || playerTag[0] != '#' {
			log.Fatalf("Invalid player tag format: %s. Player tag should start with # (e.g., #ABC123)", playerTag)
		}

		// Fetch battle log for the specified player
		var battleLog []types.Battle
		encodedPlayerTag := url.PathEscape(playerTag)
		err := apiClient.Get("/v1/players/"+encodedPlayerTag+"/battlelog", &battleLog)
		if err != nil {
			log.Printf("Error fetching battle log for player %s: %v", playerTag, err)
		} else {
			log.Printf("Successfully fetched %d battles from API for player %s", len(battleLog), playerTag)

			// Process each battle and store it in the database
			for _, battle := range battleLog {
				// Save the battle to storage
				err = s.InsertBattle(&battle)
				if err != nil {
					log.Printf("Error saving battle to storage: %v", err)
					continue
				}

				log.Printf("Successfully saved battle: %s", battle.BattleTime)
			}
		}
	}
}
