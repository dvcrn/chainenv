package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, DotConfigName)
	defaultVal := "fallback"

	cfg := &Config{
		Keys: []KeyEntry{
			{Name: "FOO", Provider: "keychain", Default: &defaultVal},
			{Name: "BAR", Provider: "1password"},
		},
	}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if !reflect.DeepEqual(cfg, loaded) {
		t.Fatalf("round-trip mismatch: %#v != %#v", cfg, loaded)
	}
}

func TestUpsertKey(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Keys: []KeyEntry{
			{Name: "FOO", Provider: "keychain"},
		},
	}

	cfg.UpsertKey(KeyEntry{Name: "FOO", Provider: "1password"})
	if len(cfg.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(cfg.Keys))
	}
	if cfg.Keys[0].Provider != "1password" {
		t.Fatalf("expected provider updated, got %s", cfg.Keys[0].Provider)
	}
}

func TestLoadOrEmptyMissing(t *testing.T) {
	t.Parallel()

	cfg, err := LoadOrEmpty(filepath.Join(t.TempDir(), DotConfigName))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg == nil || len(cfg.Keys) != 0 {
		t.Fatalf("expected empty config, got %#v", cfg)
	}
}

func TestLoadEmptyFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, DotConfigName)
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg == nil || len(cfg.Keys) != 0 {
		t.Fatalf("expected empty config, got %#v", cfg)
	}
}
