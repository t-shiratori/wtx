package port

// HookRunner defines the interface for running hooks
type HookRunner interface {
	// Run executes shell commands in the specified directory
	Run(commands []string, dir string) error
}
