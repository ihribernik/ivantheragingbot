package reader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type synthStub struct {
	dir     string
	message string
	path    string
	err     error
}

func (s *synthStub) Download(dir, message string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	s.dir = dir
	s.message = message
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	f, err := os.CreateTemp(dir, "tts-*.mp3")
	if err != nil {
		return "", err
	}
	_ = f.Close()
	s.path = f.Name()
	return s.path, nil
}

type outputStub struct {
	path string
	err  error
}

func (o *outputStub) Play(path string) error {
	o.path = path
	return o.err
}

func TestSpeakMessageParsesURLsAndDelegates(t *testing.T) {
	assetsDir := t.TempDir()
	synth := &synthStub{}
	output := &outputStub{}

	svc, err := New(assetsDir, synth, output)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	content := "mira https://example.com/test ahora"
	if err := svc.SpeakMessage("ivan", content); err != nil {
		t.Fatalf("SpeakMessage() error = %v", err)
	}

	if !strings.Contains(synth.message, "ivan dice:") {
		t.Fatalf("unexpected synthesized message: %q", synth.message)
	}
	if strings.Contains(synth.message, "https://example.com/test") {
		t.Fatalf("url should have been replaced in message: %q", synth.message)
	}
	if !strings.Contains(synth.message, "[Enlace...]") {
		t.Fatalf("expected url placeholder in message: %q", synth.message)
	}
	if synth.dir != filepath.Join(assetsDir, "tts-cache") {
		t.Fatalf("unexpected download dir: %q", synth.dir)
	}
	if output.path == "" {
		t.Fatal("expected audio output to be called")
	}
	if _, err := os.Stat(output.path); !os.IsNotExist(err) {
		t.Fatalf("expected temp file to be removed, stat err = %v", err)
	}
}
