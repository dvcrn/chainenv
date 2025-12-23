package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DotConfigName = ".chainenv.toml"
	ConfigName    = "chainenv.toml"
)

// FindConfig walks up from startDir looking for a chainenv config file.
// It prefers .chainenv.toml over chainenv.toml within the same directory.
func FindConfig(startDir string) (string, bool, error) {
	if startDir == "" {
		return "", false, fmt.Errorf("start dir is empty")
	}

	dir := filepath.Clean(startDir)
	for {
		dotPath := filepath.Join(dir, DotConfigName)
		if ok, err := isFile(dotPath); err != nil {
			return "", false, err
		} else if ok {
			return dotPath, true, nil
		}

		plainPath := filepath.Join(dir, ConfigName)
		if ok, err := isFile(plainPath); err != nil {
			return "", false, err
		} else if ok {
			return plainPath, true, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", false, nil
}

func isFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if info.IsDir() {
		return false, nil
	}

	return true, nil
}
