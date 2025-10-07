//go:build linux

package backend

import (
	"errors"
	"fmt"
	"sync"

	"github.com/99designs/keyring"
)

type KeychainBackend struct {
	ring keyring.Keyring
}

func NewKeychainBackend() (Backend, error) {
	// Restrict to Secret Service (most common) but allow KWallet if available
	cfg := keyring.Config{
		ServiceName:     "chainenv",
		AllowedBackends: []keyring.BackendType{keyring.SecretServiceBackend, keyring.KWalletBackend},
	}

	r, err := keyring.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to open Linux keyring (Secret Service/KWallet): %w", err)
	}

	return &KeychainBackend{ring: r}, nil
}

func (k *KeychainBackend) GetPassword(account string) (string, error) {
	item, err := k.ring.Get(account)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return "", fmt.Errorf("the item '%s' does not exist in the keyring", account)
		}
		return "", fmt.Errorf("error retrieving password: %w", err)
	}
	return string(item.Data), nil
}

func (k *KeychainBackend) SetPassword(account, password string, update bool) error {
	_, err := k.ring.Get(account)
	exists := err == nil
	if exists && !update {
		return fmt.Errorf("item '%s' already exists. use 'update' to update.", account)
	}

	item := keyring.Item{
		Key:         account,
		Data:        []byte(password),
		Label:       fmt.Sprintf("chainenv-%s", account),
		Description: "Set by chainenv",
	}

	if err := k.ring.Set(item); err != nil {
		return fmt.Errorf("error setting password: %w", err)
	}
	return nil
}

func (k *KeychainBackend) List() ([]string, error) {
	keys, err := k.ring.Keys()
	if err != nil {
		return nil, fmt.Errorf("error listing keyring items: %w", err)
	}
	return keys, nil
}

func (k *KeychainBackend) GetMultiplePasswords(accounts []string) (map[string]string, error) {
	results := make(map[string]string)
	resultsChan := make(chan struct {
		account  string
		password string
	}, len(accounts))

	var wg sync.WaitGroup
	for _, account := range accounts {
		wg.Add(1)
		go func(acc string) {
			defer wg.Done()
			if password, err := k.GetPassword(acc); err == nil {
				resultsChan <- struct {
					account  string
					password string
				}{acc, password}
			}
		}(account)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		results[result.account] = result.password
	}

	return results, nil
}
