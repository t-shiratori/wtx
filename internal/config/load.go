package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func loadConfigFile(path string) (*Config, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // return nil if file does not exist
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(src, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadConfig(repoRoot string) (*Config, error) {

	// 1. Local config takes priority
	localPath := LocalConfigPath(repoRoot)
	localCfg, err := loadConfigFile(localPath)
	if err != nil {
		return nil, err
	}
	if localCfg != nil {
		return localCfg, nil
	}

	// 2. Global config
	globalPath, err := GlobalConfigPath()
	if err != nil {
		return nil, err
	}
	globalCfg, err := loadConfigFile(globalPath)
	if err != nil {
		return nil, err
	}
	if globalCfg != nil {
		return globalCfg, nil
	}

	// 3. Neither exists
	return &Config{}, nil
}
