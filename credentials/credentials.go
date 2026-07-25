package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const appName = "forge"

func DataDir() (string, error) {
	xdgData := os.Getenv("XDG_DATA_HOME")
	if xdgData != "" {
		return filepath.Join(xdgData, appName), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", appName), nil
}

func credentialPath(name string) (string, error) {
	dataDir, err := DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "credentials", name), nil
}

func Read(name string) (string, error) {
	path, err := credentialPath(name)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("reading credential: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func Write(name, value string) error {
	path, err := credentialPath(name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating credentials directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(value), 0600); err != nil {
		return fmt.Errorf("writing credential: %w", err)
	}

	return nil
}

func Delete(name string) error {
	path, err := credentialPath(name)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("deleting credential: %w", err)
	}

	return nil
}
