package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"

	"github.com/ihribernik/ivantheragingbot/internal/audio"
	"github.com/ihribernik/ivantheragingbot/internal/config"
	"github.com/ihribernik/ivantheragingbot/internal/tts"
)

type Bot struct {
	cfg        config.Config
	client     *twitch.Client
	player     *audio.Player
	ttsClient  *tts.Client
	urlRe      *regexp.Regexp
	commands   map[string]func(twitch.PrivateMessage, string)
	cooldowns  map[string]map[string]time.Time
	cooldownMu sync.Mutex
	ignored    map[string]struct{}
}

func New(cfg config.Config) (*Bot, error) {
	player := audio.NewPlayer()

	client := twitch.NewClient(cfg.Username, cfg.OAuthToken)
	client.Join(cfg.Channel)

	bot := &Bot{
		cfg:       cfg,
		client:    client,
		player:    player,
		ttsClient: &tts.Client{Lang: cfg.Lang, TLD: cfg.TLD},
		urlRe:     regexp.MustCompile(`https?://(?:www\.)?[^\s/$.?#].[^\s]*`),
		commands:  make(map[string]func(twitch.PrivateMessage, string)),
		cooldowns: make(map[string]map[string]time.Time),
		ignored: map[string]struct{}{
			strings.ToLower(cfg.Username): {},
			"nightbot":                    {},
			"streamelements":              {},
		},
	}

	bot.registerHandlers()
	return bot, nil
}

func (b *Bot) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- b.client.Connect()
	}()

	select {
	case <-ctx.Done():
		b.client.Disconnect()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (b *Bot) registerHandlers() {
	b.client.OnConnect(func() {
		msg := fmt.Sprintf("Bot has connected to Twitch as %s", b.cfg.Username)
		go b.tts(msg)
	})

	b.commands["speak"] = b.commandSpeak
	b.commands["help"] = b.commandHelp
	b.commands["red"] = b.commandRed
	b.commands["alerta"] = b.commandAlerta
	b.commands["categoria"] = b.commandCategoria

	b.client.OnPrivateMessage(b.handleMessage)
}

func (b *Bot) handleMessage(message twitch.PrivateMessage) {
	log.Printf("chat message from %s: %s", message.User.Name, strings.TrimSpace(message.Message))

	if b.shouldIgnore(message.User.Name) && !b.cfg.ReadAuthorMessages {
		return
	}

	content := strings.TrimSpace(message.Message)
	if content == "" {
		return
	}

	if strings.HasPrefix(content, "!") {
		b.handleCommand(message, content)
		return
	}

	if b.cfg.AutoRead {
		text := b.parseMessage(message.User.Name, content)
		if err := b.tts(text); err != nil {
			log.Printf("tts error: %v", err)
		}
	}
}

func (b *Bot) handleCommand(message twitch.PrivateMessage, content string) {
	parts := strings.Fields(strings.TrimPrefix(content, "!"))
	if len(parts) == 0 {
		return
	}

	cmd := strings.ToLower(parts[0])
	args := strings.TrimSpace(strings.TrimPrefix(content, "!"+parts[0]))

	commandHandler, ok := b.commands[cmd]
	if !ok {
		b.client.Say(b.cfg.Channel, "El comando no existe.... !help para ver los commandos disponibles")
		return
	}

	if b.isOnCooldown(cmd, message.User.Name, b.cooldownFor(cmd)) {
		b.client.Say(b.cfg.Channel, fmt.Sprintf("%s espera un poco antes de usar %s de nuevo", message.User.DisplayName, cmd))
		return
	}

	commandHandler(message, args)
}

func (b *Bot) cooldownFor(cmd string) time.Duration {
	switch cmd {
	case "speak", "help", "alerta", "categoria":
		return 30 * time.Second
	default:
		return 0
	}
}

func (b *Bot) isOnCooldown(cmd, user string, duration time.Duration) bool {
	if duration == 0 {
		return false
	}

	b.cooldownMu.Lock()
	defer b.cooldownMu.Unlock()

	if _, ok := b.cooldowns[cmd]; !ok {
		b.cooldowns[cmd] = make(map[string]time.Time)
	}

	nextAllowed, ok := b.cooldowns[cmd][strings.ToLower(user)]
	if !ok || time.Now().After(nextAllowed) {
		b.cooldowns[cmd][strings.ToLower(user)] = time.Now().Add(duration)
		return false
	}

	return true
}

func (b *Bot) shouldIgnore(username string) bool {
	_, skip := b.ignored[strings.ToLower(username)]
	return skip
}

func (b *Bot) parseMessage(author, content string) string {
	parsed := b.urlRe.ReplaceAllString(content, "[Enlace...]")
	return fmt.Sprintf("%s dice: %s", author, parsed)
}

func (b *Bot) tts(message string) error {
	path, err := b.ttsClient.Download(filepath.Join(b.cfg.AssetsDir, "tts-cache"), message)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(path)
	}()

	return b.player.Play(path)
}

func (b *Bot) playAsset(filename string) error {
	path := filepath.Join(b.cfg.AssetsDir, filename)
	return b.player.Play(path)
}

func (b *Bot) commandSpeak(msg twitch.PrivateMessage, phrase string) {
	text := strings.TrimSpace(phrase)
	if text == "" {
		b.client.Say(b.cfg.Channel, "Uso: !speak <frase>")
		return
	}

	message := fmt.Sprintf("%s dice: %s", msg.User.Name, text)
	if err := b.tts(message); err != nil {
		log.Printf("speak tts error: %v", err)
	}
}

func (b *Bot) commandHelp(msg twitch.PrivateMessage, _ string) {
	b.client.Say(b.cfg.Channel, "Comandos disponibles >> !speak, !help, !red, !alerta, !categoria")
}

func (b *Bot) commandRed(msg twitch.PrivateMessage, _ string) {
	if err := b.playAsset("codec.mp3"); err != nil {
		log.Printf("red command error: %v", err)
	}
	b.client.Say(b.cfg.Channel, "Notificacion de red baja enviada...")
}

func (b *Bot) commandAlerta(msg twitch.PrivateMessage, _ string) {
	if err := b.playAsset("alerta.mp3"); err != nil {
		log.Printf("alerta command error: %v", err)
	}
	b.client.Say(b.cfg.Channel, "ya se alerto al streamer")
}

func (b *Bot) commandCategoria(msg twitch.PrivateMessage, _ string) {
	if err := b.playAsset("categoria.mp3"); err != nil {
		log.Printf("categoria command error: %v", err)
	}
	b.client.Say(b.cfg.Channel, "Se le aviso al streamer que cambie la categoria...")
}
