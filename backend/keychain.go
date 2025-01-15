package backend

import (
	"fmt"
	"os/exec"
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
