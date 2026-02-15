package port

import "wtx/internal/domain/worktree"

// GitRepository defines the interface for git operations
type GitRepository interface {
	// AddWorktree creates a new worktree
	AddWorktree(opts worktree.AddOptions) error

	// RemoveWorktree removes an existing worktree
	RemoveWorktree(path string, force bool) error

	// DeleteBranch deletes a git branch
	DeleteBranch(branch string, force bool) error

	// ListWorktrees returns all worktrees
	ListWorktrees() ([]worktree.Worktree, error)

	// BranchFromWorktree returns the branch name for a worktree path
	BranchFromWorktree(path string) (string, error)

	// RepoRoot returns the repository root path
	RepoRoot() (string, error)
}
