package config

import (
	"os"
)

func EnsureConfigRootDir(configRootDir string) error {
	return os.MkdirAll(configRootDir, 0755)
}
