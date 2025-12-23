package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Keys []KeyEntry `toml:"keys"`
}

type KeyEntry struct {
	Name     string  `toml:"name"`
	Provider string  `toml:"provider,omitempty"`
	Default  *string `toml:"default,omitempty"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseConfig(data)
}

func LoadOrEmpty(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	return parseConfig(data)
}

func Save(path string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	return atomicWriteFile(path, data, 0o644)
}

func (c *Config) FindKey(name string) (*KeyEntry, bool) {
	for i := range c.Keys {
		if c.Keys[i].Name == name {
			return &c.Keys[i], true
		}
	}
	return nil, false
}

func (c *Config) UpsertKey(entry KeyEntry) {
	for i := range c.Keys {
		if c.Keys[i].Name == entry.Name {
			c.Keys[i] = entry
			return
		}
	}
	c.Keys = append(c.Keys, entry)
}

func parseConfig(data []byte) (*Config, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return &Config{}, nil
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}
