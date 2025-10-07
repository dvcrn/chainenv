//go:build darwin

package backend

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type KeychainBackend struct{}

func NewKeychainBackend() (Backend, error) {
	return &KeychainBackend{}, nil
}

func (k *KeychainBackend) GetPassword(account string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", account, "-s", fmt.Sprintf("chainenv-%s", account), "-w")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error retrieving password: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (k *KeychainBackend) SetPassword(account, password string, update bool) error {
	args := []string{"add-generic-password", "-a", account, "-s", fmt.Sprintf("chainenv-%s", account), "-w", password, "-j", "Set by chainenv"}
	if update {
		args = append(args, "-U")
	}

	cmd := exec.Command("security", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error setting password: %v: %s", err, output)
	}
	return nil
}

var keychainServiceRegex = regexp.MustCompile(`"svce"<blob>="chainenv-(.+)"`)

func (k *KeychainBackend) List() ([]string, error) {
	cmd := exec.Command("security", "dump-keychain")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe for security command: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting security command: %w", err)
	}

	var accounts []string
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		matches := keychainServiceRegex.FindStringSubmatch(scanner.Text())
		if len(matches) > 1 {
			account := matches[1]
			if !seen[account] {
				accounts = append(accounts, account)
				seen[account] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading security command output: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("error listing keychain items: %v", err)
	}

	return accounts, nil
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
