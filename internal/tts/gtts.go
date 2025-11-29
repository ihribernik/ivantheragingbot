package tts

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Client downloads speech MP3 files using Google's public translate TTS endpoint.
type Client struct {
	Lang string
	TLD  string
}

// Download saves a synthesized MP3 for the provided message into dir and returns the file path.
func (c *Client) Download(dir, message string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("ensure dir: %w", err)
	}

	reqURL := fmt.Sprintf(
		"https://translate.google.%s/translate_tts?ie=UTF-8&client=tw-ob&q=%s&tl=%s&ttsspeed=1",
		c.TLD,
		url.QueryEscape(message),
		url.QueryEscape(c.Lang),
	)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ivantheragingbot/1.0)")

	httpClient := http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download tts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tts request failed with status %s", resp.Status)
	}

	file, err := os.CreateTemp(dir, "tts-*.mp3")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("write mp3: %w", err)
	}

	return filepath.Clean(file.Name()), nil
}
