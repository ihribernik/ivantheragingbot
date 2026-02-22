package tts

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultTLD       = "com"
	defaultUserAgent = "Mozilla/5.0 (compatible; ivantheragingbot/1.0)"
)

// Client downloads speech MP3 files using Google's public translate TTS endpoint.
type Client struct {
	lang       string
	tld        string
	httpClient *http.Client
}

// NewClient creates a TTS client with validated configuration.
func NewClient(lang, tld string, httpClient *http.Client) (*Client, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return nil, fmt.Errorf("lang is required")
	}

	tld = strings.TrimSpace(tld)
	if tld == "" {
		tld = defaultTLD
	}

	if httpClient == nil {
		return nil, fmt.Errorf("http client is required")
	}

	return &Client{
		lang:       lang,
		tld:        tld,
		httpClient: httpClient,
	}, nil
}

// Download saves a synthesized MP3 for the provided message into dir and returns the file path.
func (c *Client) Download(dir, message string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("ensure dir: %w", err)
	}

	reqURL := fmt.Sprintf(
		"https://translate.google.%s/translate_tts?ie=UTF-8&client=tw-ob&q=%s&tl=%s&ttsspeed=1",
		c.tld,
		url.QueryEscape(message),
		url.QueryEscape(c.lang),
	)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
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
