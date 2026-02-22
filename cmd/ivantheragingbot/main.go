package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/ihribernik/ivantheragingbot/internal/bootstrap"
	"github.com/ihribernik/ivantheragingbot/internal/config"
	"github.com/ihribernik/ivantheragingbot/internal/logging"
)

func main() {
	if err := run(); err != nil {
		logging.Error("fatal: %v", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working dir: %w", err)
	}

	cfg, err := config.Load(filepath.Clean(workingDir))
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	chatBot, err := bootstrap.NewBot(cfg)
	if err != nil {
		return fmt.Errorf("build app: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logging.Info("starting ivantheragingbot as %s in channel %s", cfg.Username, cfg.Channel)
	if err := chatBot.Connect(ctx); err != nil && !errors.Is(err, context.Canceled) && ctx.Err() == nil {
		return fmt.Errorf("run bot: %w", err)
	}

	return nil
}
