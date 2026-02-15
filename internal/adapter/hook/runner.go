package hook

import (
	oldhook "wtx/internal/hook"
)

// Runner implements port.HookRunner
type Runner struct{}

// NewRunner creates a new Runner
func NewRunner() *Runner {
	return &Runner{}
}

// Run executes shell commands in the specified directory
func (r *Runner) Run(commands []string, dir string) error {
	return oldhook.Run(commands, dir)
}
