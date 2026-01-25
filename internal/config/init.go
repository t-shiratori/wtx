package config

import (
	"os"
)

func EnsureWTDir(wtDir string) error {
	return os.MkdirAll(wtDir, 0755)
}
