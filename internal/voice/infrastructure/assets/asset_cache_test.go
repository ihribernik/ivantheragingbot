package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewAssetCacheAndResolve(t *testing.T) {
	baseDir := t.TempDir()
	file := filepath.Join(baseDir, "alerta.mp3")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	cache, err := NewAssetCache(baseDir, map[string]string{"alerta": "alerta.mp3"})
	if err != nil {
		t.Fatalf("NewAssetCache() error = %v", err)
	}

	path, err := cache.Resolve("alerta")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if path != file {
		t.Fatalf("unexpected resolved path: got %q want %q", path, file)
	}
}

func TestResolveNotFound(t *testing.T) {
	baseDir := t.TempDir()
	cache, err := NewAssetCache(baseDir, map[string]string{})
	if err != nil {
		t.Fatalf("NewAssetCache() error = %v", err)
	}

	if _, err := cache.Resolve("missing"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestNewAssetCacheRejectsTraversal(t *testing.T) {
	baseDir := t.TempDir()
	parent := filepath.Dir(baseDir)
	outside := filepath.Join(parent, "outside.mp3")
	if err := os.WriteFile(outside, []byte("x"), 0o644); err != nil {
		t.Fatalf("write outside asset: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(outside) })

	_, err := NewAssetCache(baseDir, map[string]string{"bad": ".." + string(filepath.Separator) + "outside.mp3"})
	if err == nil {
		t.Fatal("expected path traversal error")
	}
}
