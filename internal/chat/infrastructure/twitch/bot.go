package twitch

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"

	chatapp "github.com/ihribernik/ivantheragingbot/internal/chat/application"
	"github.com/ihribernik/ivantheragingbot/internal/config"
	"github.com/ihribernik/ivantheragingbot/internal/logging"
	"github.com/ihribernik/ivantheragingbot/internal/voice/application/reader"
	"github.com/ihribernik/ivantheragingbot/internal/voice/application/soundboard"
)

type Bot struct {
	cfg        config.Config
	client     *twitch.Client
	terminal   chatapp.Stream
	voice      *reader.Service
	commands   map[string]func(twitch.PrivateMessage, string)
	cooldowns  map[string]map[string]time.Time
	cooldownMu sync.Mutex
	ignored    map[string]struct{}
	soundboard *soundboard.Service
}

func New(cfg config.Config, terminalStream chatapp.Stream, voiceSvc *reader.Service, soundboardSvc *soundboard.Service) (*Bot, error) {
	if terminalStream == nil {
		return nil, fmt.Errorf("terminal stream is required")
	}
	if voiceSvc == nil {
		return nil, fmt.Errorf("voice service is required")
	}
	if soundboardSvc == nil {
		return nil, fmt.Errorf("soundboard service is required")
	}

	client := twitch.NewClient(cfg.Username, cfg.OAuthToken)
	client.Join(cfg.Channel)

	bot := &Bot{
		cfg:       cfg,
		client:    client,
		terminal:  terminalStream,
		voice:     voiceSvc,
		commands:  make(map[string]func(twitch.PrivateMessage, string)),
		cooldowns: make(map[string]map[string]time.Time),
		ignored: map[string]struct{}{
			strings.ToLower(cfg.Username): {},
			"nightbot":                    {},
			"streamelements":              {},
		},
		soundboard: soundboardSvc,
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

func (b *Bot) Connect(ctx context.Context) error {
	return b.Start(ctx)
}

func (b *Bot) Say(channel, message string) {
	b.client.Say(channel, message)
}

func (b *Bot) registerHandlers() {
	b.client.OnConnect(func() {
		msg := fmt.Sprintf("Bot has connected to Twitch as %s", b.cfg.Username)
		b.terminal.System(msg)
		go func() {
			if err := b.voice.SpeakText(msg); err != nil {
				logging.Error("connect tts error: %v", err)
			}
		}()
	})

	b.commands["speak"] = b.commandSpeak
	b.commands["help"] = b.commandHelp
	b.commands["red"] = b.commandRed
	b.commands["alerta"] = b.commandAlerta
	b.commands["categoria"] = b.commandCategoria

	b.client.OnPrivateMessage(b.handleMessage)
}

func (b *Bot) handleMessage(message twitch.PrivateMessage) {
	if b.shouldIgnore(message.User.Name) && !b.cfg.ReadAuthorMessages {
		return
	}

	content := strings.TrimSpace(message.Message)
	if content == "" {
		return
	}
	author := message.User.DisplayName
	if author == "" {
		author = message.User.Name
	}
	b.terminal.Message(author, message.User.Color, content)

	if strings.HasPrefix(content, "!") {
		b.handleCommand(message, content)
		return
	}

	if b.cfg.AutoRead {
		if err := b.voice.SpeakMessage(message.User.Name, content); err != nil {
			logging.Error("tts error: %v", err)
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

func (b *Bot) commandSpeak(msg twitch.PrivateMessage, phrase string) {
	text := strings.TrimSpace(phrase)
	if text == "" {
		b.client.Say(b.cfg.Channel, "Uso: !speak <frase>")
		return
	}

	if err := b.voice.SpeakMessage(msg.User.Name, text); err != nil {
		logging.Error("speak tts error: %v", err)
	}
}

func (b *Bot) commandHelp(msg twitch.PrivateMessage, _ string) {
	b.client.Say(b.cfg.Channel, "Comandos disponibles >> !speak, !help, !red, !alerta, !categoria")
}

func (b *Bot) commandRed(msg twitch.PrivateMessage, _ string) {
	if err := b.soundboard.Play("red"); err != nil {
		logging.Error("red command error: %v", err)
		b.client.Say(b.cfg.Channel, "No se pudo reproducir la alerta de red.")
		return
	}
	b.client.Say(b.cfg.Channel, "Notificacion de red baja enviada...")
}

func (b *Bot) commandAlerta(msg twitch.PrivateMessage, _ string) {
	if err := b.soundboard.Play("alerta"); err != nil {
		logging.Error("alerta command error: %v", err)
		b.client.Say(b.cfg.Channel, "No se pudo reproducir la alerta.")
		return
	}
	b.client.Say(b.cfg.Channel, "ya se alerto al streamer")
}

func (b *Bot) commandCategoria(msg twitch.PrivateMessage, _ string) {
	if err := b.soundboard.Play("categoria"); err != nil {
		logging.Error("categoria command error: %v", err)
		b.client.Say(b.cfg.Channel, "No se pudo reproducir el aviso de categoria.")
		return
	}
	b.client.Say(b.cfg.Channel, "Se le aviso al streamer que cambie la categoria...")
}
