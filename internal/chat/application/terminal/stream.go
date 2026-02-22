package terminal

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Stream struct {
	out io.Writer
	now func() time.Time
	mu  sync.Mutex
}

const (
	ansiReset  = "\x1b[0m"
	ansiCyan   = "\x1b[36m"
	ansiYellow = "\x1b[33m"
)

func New(out io.Writer) (*Stream, error) {
	if out == nil {
		return nil, fmt.Errorf("output writer is required")
	}

	return &Stream{
		out: out,
		now: time.Now,
	}, nil
}

func (s *Stream) Message(author, colorHex, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := s.now().Format("15:04:05")
	_, _ = fmt.Fprintf(s.out, "[%s] %s: %s\n", timestamp, colorize(author, colorHex), styleContent(content))
}

func (s *Stream) System(content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := s.now().Format("15:04:05")
	_, _ = fmt.Fprintf(s.out, "[%s] %s[system]%s: %s\n", timestamp, ansiCyan, ansiReset, content)
}

func colorize(text, colorHex string) string {
	colorHex = strings.TrimSpace(colorHex)
	if colorHex == "" {
		return text
	}

	colorHex = strings.TrimPrefix(colorHex, "#")
	if len(colorHex) != 6 {
		return text
	}

	r, err := strconv.ParseUint(colorHex[0:2], 16, 8)
	if err != nil {
		return text
	}
	g, err := strconv.ParseUint(colorHex[2:4], 16, 8)
	if err != nil {
		return text
	}
	b, err := strconv.ParseUint(colorHex[4:6], 16, 8)
	if err != nil {
		return text
	}

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

func styleContent(content string) string {
	if strings.HasPrefix(strings.TrimSpace(content), "!") {
		return ansiYellow + content + ansiReset
	}
	return content
}
