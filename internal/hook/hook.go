package hook

import (
	"os"
	"os/exec"
)

func Run(cmd string, dir string) error {
	if cmd == "" {
		return nil
	}

	c := exec.Command("sh", "-c", cmd)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
