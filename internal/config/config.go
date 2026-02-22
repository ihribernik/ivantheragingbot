package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultTTSHTTPTimeout = 10 * time.Second

// Config holds runtime configuration for the Twitch bot.
type Config struct {
	Username           string
	OAuthToken         string
	Channel            string
	AssetsDir          string
	Lang               string
	TLD                string
	TTSHTTPTimeout     time.Duration
	AutoRead           bool
	ReadAuthorMessages bool
}

// Load reads configuration values from environment variables and fills in defaults
// so the bot can start with minimal setup.
func Load(baseDir string) (Config, error) {
	ttsHTTPTimeout, err := durationEnvWithDefault("TTS_HTTP_TIMEOUT", defaultTTSHTTPTimeout)
	if err != nil {
		return Config{}, err
	}
	autoRead, err := boolEnvWithDefault("AUTO_READ", false)
	if err != nil {
		return Config{}, err
	}
	readAuthorMessages, err := boolEnvWithDefault("READ_AUTHOR_MESSAGE", false)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Username:           strings.TrimSpace(os.Getenv("TWITCH_USERNAME")),
		OAuthToken:         strings.TrimSpace(os.Getenv("IRC_TOKEN")),
		Channel:            strings.TrimSpace(os.Getenv("TWITCH_CHANNEL")),
		AssetsDir:          filepath.Join(baseDir, "assets"),
		Lang:               envWithDefault("TTS_LANG", "es"),
		TLD:                envWithDefault("TTS_TLD", "com.ar"),
		TTSHTTPTimeout:     ttsHTTPTimeout,
		AutoRead:           autoRead,
		ReadAuthorMessages: readAuthorMessages,
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

func durationEnvWithDefault(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration (e.g. 10s): %w", key, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return d, nil
}

func boolEnvWithDefault(key string, fallback bool) (bool, error) {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback, nil
	}

	switch value {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("%s must be a valid boolean value", key)
	}
}
