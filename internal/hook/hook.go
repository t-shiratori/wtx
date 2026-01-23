package hook

import (
	"os"
	"os/exec"
)

func Run(cmd []string, dir string) error {
	if len(cmd) == 0 {
		return nil
	}

	for _, c := range cmd {
		c := exec.Command("sh", "-c", c)
		c.Dir = dir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return err
		}
	}

	return nil
}
