package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds runtime configuration for the Twitch bot.
type Config struct {
	Username           string
	OAuthToken         string
	Channel            string
	AssetsDir          string
	Lang               string
	TLD                string
	AutoRead           bool
	ReadAuthorMessages bool
}

// Load reads configuration values from environment variables and fills in defaults
// so the bot can start with minimal setup.
func Load(baseDir string) (Config, error) {
	cfg := Config{
		Username:           strings.TrimSpace(os.Getenv("TWITCH_USERNAME")),
		OAuthToken:         strings.TrimSpace(os.Getenv("IRC_TOKEN")),
		Channel:            strings.TrimSpace(os.Getenv("TWITCH_CHANNEL")),
		AssetsDir:          filepath.Join(baseDir, "assets"),
		Lang:               envWithDefault("TTS_LANG", "es"),
		TLD:                envWithDefault("TTS_TLD", "com.ar"),
		AutoRead:           hasFlag("AUTO_READ"),
		ReadAuthorMessages: hasFlag("READ_AUTHOR_MESSAGE"),
	}

	if cfg.Channel == "" {
		cfg.Channel = "ivantheragingpython"
	}

	if cfg.Username == "" {
		return Config{}, fmt.Errorf("TWITCH_USERNAME is required")
	}

	if cfg.OAuthToken == "" {
		return Config{}, fmt.Errorf("IRC_TOKEN is required")
	}

	if !strings.HasPrefix(cfg.OAuthToken, "oauth:") {
		cfg.OAuthToken = "oauth:" + cfg.OAuthToken
	}

	return cfg, nil
}

func envWithDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func hasFlag(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}
