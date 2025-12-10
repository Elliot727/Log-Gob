package main

import (
	"database/sql"
	"log"
	"net/url"

	"github.com/elliot727/log-gob/internal/api"
	"github.com/elliot727/log-gob/internal/config"
	"github.com/elliot727/log-gob/internal/storage"
	"github.com/elliot727/log-gob/internal/types"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if cfg.APIKey == "" {
		log.Fatal("APIKEY not set in environment variables")
	}

	if cfg.PlayerTag == "" {
		log.Fatal("PLAYERTAG not set in environment variables")
	}

	if len(cfg.PlayerTag) < 2 || cfg.PlayerTag[0] != '#' {
		log.Fatalf("invalid tag: %s", cfg.PlayerTag)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	s := storage.NewStorage(db)

	if err := s.Init(); err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)

	var battleLog []types.Battle

	escaped := url.PathEscape(cfg.PlayerTag)
	if err := client.Get("/v1/players/"+escaped+"/battlelog", &battleLog); err != nil {
		log.Fatalf("failed to fetch battles: %v", err)
	}

	log.Printf("Fetched %d battles for %s", len(battleLog), cfg.PlayerTag)

	for _, b := range battleLog {
		if err := s.InsertBattle(&b); err != nil {
			log.Printf("storage error: %v", err)
			continue
		}
		log.Printf("Saved battle: %s", b.BattleTime)
	}
}
