package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type AssetCache struct {
	baseDir string
	mu      sync.RWMutex
	cache   map[string]string
}

func NewAssetCache(baseDir string, items map[string]string) (*AssetCache, error) {
	c := &AssetCache{
		baseDir: filepath.Clean(baseDir),
		cache:   make(map[string]string, len(items)),
	}

	for key, name := range items {
		fullPath := filepath.Clean(filepath.Join(c.baseDir, name))
		rel, err := filepath.Rel(c.baseDir, fullPath)
		if err != nil {
			return nil, fmt.Errorf("resolve asset %q: %w", key, err)
		}

		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("asset path escapes base dir for key %q", key)
		}

		if _, err := os.Stat(fullPath); err != nil {
			return nil, fmt.Errorf("invalid asset for key %q (%s): %w", key, fullPath, err)
		}

		c.cache[key] = fullPath
	}

	return c, nil
}

func (c *AssetCache) Resolve(key string) (string, error) {
	c.mu.RLock()
	path, ok := c.cache[key]
	c.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("asset not found: %s", key)
	}

	return path, nil
}
