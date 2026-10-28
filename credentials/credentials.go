package credentials

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nanoteck137/forge"
	"github.com/pelletier/go-toml/v2"
)

const fileName = "credentials.toml"

type Store struct {
	Token       string `toml:"token"`
	GitHubToken string `toml:"github_token"`
}

func DefaultDir() (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, forge.AppName), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".local", "share", forge.AppName), nil
}

func Path() (string, error) {
	dir, err := DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fileName), nil
}

func Load() (*Store, error) {
	path, err := Path()
	if err != nil {
		return nil, fmt.Errorf("getting credentials path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Store{}, nil
		}
		return nil, fmt.Errorf("reading credentials: %w", err)
	}

	var store Store
	if err := toml.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parsing credentials: %w", err)
	}

	return &store, nil
}

func Save(store *Store) error {
	path, err := Path()
	if err != nil {
		return fmt.Errorf("getting credentials path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating credentials directory: %w", err)
	}

	data, err := toml.Marshal(store)
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing credentials: %w", err)
	}

	return nil
}
