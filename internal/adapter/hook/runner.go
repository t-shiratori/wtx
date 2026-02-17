package hook

import (
	"os"
	"os/exec"
)

// Runner implements port.HookRunner
type Runner struct{}

// NewRunner creates a new Runner
func NewRunner() *Runner {
	return &Runner{}
}

// Run executes shell commands in the specified directory
func (r *Runner) Run(commands []string, dir string) error {
	if len(commands) == 0 {
		return nil
	}

	for _, c := range commands {
		cmd := exec.Command("sh", "-c", c)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
