package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfig(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	makeFile := func(path string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	t.Run("no config", func(t *testing.T) {
		path, ok, err := FindConfig(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatalf("expected no config, got %s", path)
		}
	})

	t.Run("dotfile in current dir wins", func(t *testing.T) {
		dir := filepath.Join(tempDir, "dotfile")
		makeFile(filepath.Join(dir, DotConfigName))
		makeFile(filepath.Join(dir, ConfigName))

		path, ok, err := FindConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatalf("expected config to be found")
		}
		if want := filepath.Join(dir, DotConfigName); path != want {
			t.Fatalf("expected %s, got %s", want, path)
		}
	})

	t.Run("plain config in current dir", func(t *testing.T) {
		dir := filepath.Join(tempDir, "plain")
		makeFile(filepath.Join(dir, ConfigName))

		path, ok, err := FindConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatalf("expected config to be found")
		}
		if want := filepath.Join(dir, ConfigName); path != want {
			t.Fatalf("expected %s, got %s", want, path)
		}
	})

	t.Run("nearest directory wins", func(t *testing.T) {
		parent := filepath.Join(tempDir, "parent")
		child := filepath.Join(parent, "child")
		makeFile(filepath.Join(parent, DotConfigName))
		makeFile(filepath.Join(child, ConfigName))

		path, ok, err := FindConfig(child)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatalf("expected config to be found")
		}
		if want := filepath.Join(child, ConfigName); path != want {
			t.Fatalf("expected %s, got %s", want, path)
		}
	})
}
