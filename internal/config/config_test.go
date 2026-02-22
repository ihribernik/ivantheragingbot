package config

import (
	"strings"
	"testing"
	"time"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("TWITCH_USERNAME", "bot")
	t.Setenv("IRC_TOKEN", "token")
	t.Setenv("TWITCH_CHANNEL", "chan")
}

func TestLoadParsesBooleanFlags(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("AUTO_READ", "true")
	t.Setenv("READ_AUTHOR_MESSAGE", "0")

	cfg, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.AutoRead {
		t.Fatal("expected AUTO_READ=true")
	}
	if cfg.ReadAuthorMessages {
		t.Fatal("expected READ_AUTHOR_MESSAGE=false")
	}
}

func TestLoadRejectsInvalidBoolean(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("AUTO_READ", "banana")

	_, err := Load(t.TempDir())
	if err == nil {
		t.Fatal("expected invalid boolean error")
	}
	if !strings.Contains(err.Error(), "AUTO_READ") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadParsesTTSHTTPTimeout(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("TTS_HTTP_TIMEOUT", "3s")

	cfg, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.TTSHTTPTimeout != 3*time.Second {
		t.Fatalf("unexpected timeout: %v", cfg.TTSHTTPTimeout)
	}
}
