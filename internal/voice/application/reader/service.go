package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var urlRe = regexp.MustCompile(`https?://(?:www\.)?[^\s/$.?#].[^\s]*`)

type Synthesizer interface {
	Download(dir, message string) (string, error)
}

type AudioOutput interface {
	Play(path string) error
}

type Service struct {
	assetsDir   string
	synthesizer Synthesizer
	output      AudioOutput
}

func New(assetsDir string, synthesizer Synthesizer, output AudioOutput) (*Service, error) {
	if assetsDir == "" {
		return nil, fmt.Errorf("assets dir is required")
	}
	if synthesizer == nil {
		return nil, fmt.Errorf("synthesizer is required")
	}
	if output == nil {
		return nil, fmt.Errorf("audio output is required")
	}

	return &Service{
		assetsDir:   assetsDir,
		synthesizer: synthesizer,
		output:      output,
	}, nil
}

func (s *Service) SpeakMessage(author, content string) error {
	parsed := urlRe.ReplaceAllString(content, "[Enlace...]")
	message := fmt.Sprintf("%s dice: %s", author, parsed)
	return s.SpeakText(message)
}

func (s *Service) SpeakText(message string) error {
	path, err := s.synthesizer.Download(filepath.Join(s.assetsDir, "tts-cache"), message)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(path)
	}()

	return s.output.Play(path)
}
