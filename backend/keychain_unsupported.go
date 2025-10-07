//go:build !darwin && !linux

package backend

import (
	"fmt"
	"runtime"
)

// NewKeychainBackend returns an error on unsupported platforms.
func NewKeychainBackend() (Backend, error) {
	return nil, fmt.Errorf("keychain backend is unsupported on %s", runtime.GOOS)
}
