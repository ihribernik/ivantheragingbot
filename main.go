package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/ihribernik/ivantheragingbot/internal/bot"
	"github.com/ihribernik/ivantheragingbot/internal/config"
	"github.com/ihribernik/ivantheragingbot/internal/logging"
)

func main() {
	_ = godotenv.Load()

	workingDir, err := os.Getwd()
	if err != nil {
		logging.Fatal("failed to resolve working dir: %v", err)
	}

	cfg, err := config.Load(filepath.Clean(workingDir))
	if err != nil {
		logging.Fatal("config error: %v", err)
	}

	chatBot, err := bot.New(cfg)
	if err != nil {
		logging.Fatal("bot init error: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logging.Info("starting ivantheragingbot as %s in channel %s", cfg.Username, cfg.Channel)
	if err := chatBot.Start(ctx); err != nil && ctx.Err() == nil {
		logging.Fatal("bot stopped with error: %v", err)
	}
}
