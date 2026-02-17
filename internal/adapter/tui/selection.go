package tui

// Selection represents a single worktree removal target.
// This is used by both TUI and non-interactive CLI flows.
type Selection struct {
	Path         string
	DeleteBranch bool
	Force        bool
}
