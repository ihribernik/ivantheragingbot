package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/ihribernik/ivantheragingbot/internal/bot"
	"github.com/ihribernik/ivantheragingbot/internal/config"
)

func main() {
	_ = godotenv.Load()

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to resolve working dir: %v", err)
	}

	cfg, err := config.Load(filepath.Clean(workingDir))
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	chatBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("bot init error: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("starting ivantheragingbot as %s in channel %s", cfg.Username, cfg.Channel)
	if err := chatBot.Start(ctx); err != nil && ctx.Err() == nil {
		log.Fatalf("bot stopped with error: %v", err)
	}
}
