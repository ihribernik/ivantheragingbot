package bootstrap

import (
	"fmt"
	"net/http"
	"os"

	chatapp "github.com/ihribernik/ivantheragingbot/internal/chat/application"
	"github.com/ihribernik/ivantheragingbot/internal/chat/application/terminal"
	chatbot "github.com/ihribernik/ivantheragingbot/internal/chat/infrastructure/twitch"
	"github.com/ihribernik/ivantheragingbot/internal/config"
	"github.com/ihribernik/ivantheragingbot/internal/voice/application/reader"
	"github.com/ihribernik/ivantheragingbot/internal/voice/application/soundboard"
	"github.com/ihribernik/ivantheragingbot/internal/voice/infrastructure/assets"
	"github.com/ihribernik/ivantheragingbot/internal/voice/infrastructure/audio"
	tts "github.com/ihribernik/ivantheragingbot/internal/voice/infrastructure/tts"
)

func NewBot(cfg config.Config) (chatapp.Client, error) {
	terminalStream, err := terminal.New(os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("init terminal stream: %w", err)
	}

	player := audio.NewBeepPlayer()

	resolver, err := assets.NewAssetCache(cfg.AssetsDir, map[string]string{
		"red":       "codec.mp3",
		"alerta":    "alerta.mp3",
		"categoria": "categoria.mp3",
	})
	if err != nil {
		return nil, fmt.Errorf("init asset cache: %w", err)
	}

	httpClient := &http.Client{Timeout: cfg.TTSHTTPTimeout}
	ttsClient, err := tts.NewClient(cfg.Lang, cfg.TLD, httpClient)
	if err != nil {
		return nil, fmt.Errorf("init tts client: %w", err)
	}

	voiceSvc, err := reader.New(cfg.AssetsDir, ttsClient, player)
	if err != nil {
		return nil, fmt.Errorf("init voice service: %w", err)
	}

	soundboardSvc := soundboard.New(resolver, player)

	chatBot, err := chatbot.New(cfg, terminalStream, voiceSvc, soundboardSvc)
	if err != nil {
		return nil, fmt.Errorf("init twitch bot: %w", err)
	}

	return chatBot, nil
}
