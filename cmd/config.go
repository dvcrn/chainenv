package cmd

import (
	"fmt"
	"os"

	"github.com/dvcrn/chainenv/backend"
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

func ensureOpServiceAccountToken(cfg *config.Config) error {
	if os.Getenv("OP_SERVICE_ACCOUNT_TOKEN") != "" {
		return nil
	}
	if cfg == nil || cfg.OnePassword == nil {
		return nil
	}
	tokenKey := cfg.OnePassword.ServiceAccountTokenKey
	if tokenKey == "" {
		return nil
	}

	keychain, err := backend.NewKeychainBackend()
	if err != nil {
		return fmt.Errorf("keychain backend unavailable: %w", err)
	}

	token, err := keychain.GetPassword(tokenKey)
	if err != nil {
		return fmt.Errorf("failed to load %s from keychain: %w", tokenKey, err)
	}

	if err := os.Setenv("OP_SERVICE_ACCOUNT_TOKEN", token); err != nil {
		return fmt.Errorf("failed to set OP_SERVICE_ACCOUNT_TOKEN: %w", err)
	}
	return nil
}
