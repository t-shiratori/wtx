package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig(repoRoot string) (*Config, error) {
	path := filepath.Join(repoRoot, ".wt", "config.toml")

	configSrc, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil // 無ければ空
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(configSrc, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
