package worktree

// Worktree represents a git worktree entity
type Worktree struct {
	Path   string
	Branch string
	Commit string
}

// NewWorktree creates a new Worktree entity
func NewWorktree(path, branch, commit string) *Worktree {
	return &Worktree{
		Path:   path,
		Branch: branch,
		Commit: commit,
	}
}
