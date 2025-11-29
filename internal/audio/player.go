package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Player wraps speaker playback with a fixed sample rate to ensure all MP3 assets
// and downloaded TTS clips can be reproduced without overlapping.
type Player struct {
	sampleRate beep.SampleRate
	initOnce   sync.Once
	mu         sync.Mutex
}

// NewPlayer initializes the audio speaker with a target sample rate.
func NewPlayer() *Player {
	return &Player{}
}

// Play synchronously plays the MP3 file at the given path.
func (p *Player) Play(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("open audio file: %w", err)
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return fmt.Errorf("decode mp3: %w", err)
	}
	defer streamer.Close()

	if err := p.ensureSpeaker(format.SampleRate); err != nil {
		return err
	}

	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))
	<-done
	return nil
}

func (p *Player) ensureSpeaker(sampleRate beep.SampleRate) error {
	var initErr error
	p.initOnce.Do(func() {
		bufferSize := sampleRate.N(time.Second / 10) // 100ms buffer helps avoid stutter.
		initErr = speaker.Init(sampleRate, bufferSize)
		p.sampleRate = sampleRate
	})
	if initErr != nil {
		return fmt.Errorf("init speaker: %w", initErr)
	}

	if sampleRate != p.sampleRate {
		speaker.Clear()
		bufferSize := sampleRate.N(time.Second / 10)
		if err := speaker.Init(sampleRate, bufferSize); err != nil {
			return fmt.Errorf("reinit speaker: %w", err)
		}
		p.sampleRate = sampleRate
	}
	return nil
}
