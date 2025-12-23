package cmd

import (
	"os"

	"github.com/dvcrn/chainenv/config"
)

func loadConfig() (*config.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configPath, ok, err := config.FindConfig(cwd)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return config.Load(configPath)
}

func resolveKeyConfig(cfg *config.Config, name, fallbackProvider string) (string, *string) {
	provider := fallbackProvider
	var defaultValue *string
	if cfg == nil {
		return provider, nil
	}

	if entry, ok := cfg.FindKey(name); ok {
		if entry.Provider != "" {
			provider = entry.Provider
		}
		defaultValue = entry.Default
	}

	return provider, defaultValue
}
