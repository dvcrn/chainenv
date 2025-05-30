package backend

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type KeychainBackend struct{}

func NewKeychainBackend() *KeychainBackend {
	return &KeychainBackend{}
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

func (k *KeychainBackend) List() ([]string, error) {
	// Use security dump-keychain to get all keychain items
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error listing keychain items: %v", err)
	}

	// Parse the output to find chainenv items
	var accounts []string
	seen := make(map[string]bool)
	
	// Look for service attributes that match "chainenv-*"
	lines := strings.Split(string(output), "\n")
	serviceRegex := regexp.MustCompile(`"svce"<blob>="chainenv-(.+)"`)
	
	for _, line := range lines {
		matches := serviceRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			account := matches[1]
			if !seen[account] {
				accounts = append(accounts, account)
				seen[account] = true
			}
		}
	}
	
	return accounts, nil
}
