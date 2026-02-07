package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func loadConfigFile(path string) (*Config, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 無ければ nil
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

	// 1. ローカルを優先
	localPath := LocalConfigPath(repoRoot)
	localCfg, err := loadConfigFile(localPath)
	if err != nil {
		return nil, err
	}
	if localCfg != nil {
		return localCfg, nil
	}

	// 2. グローバル
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

	// 3. どちらも無い
	return &Config{}, nil
}
