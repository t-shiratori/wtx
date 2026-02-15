package git

import (
	"wtx/internal/domain/worktree"
	oldgit "wtx/internal/git"
)

// Repository implements port.GitRepository
type Repository struct{}

// NewRepository creates a new Repository
func NewRepository() *Repository {
	return &Repository{}
}

// AddWorktree creates a new worktree
func (r *Repository) AddWorktree(opts worktree.AddOptions) error {
	return oldgit.AddWorktree(opts.Path, opts.Branch, opts.CreateBranch, opts.BaseBranch)
}

// RemoveWorktree removes an existing worktree
func (r *Repository) RemoveWorktree(path string, force bool) error {
	return oldgit.RemoveWorktree(path, force)
}

// DeleteBranch deletes a git branch
func (r *Repository) DeleteBranch(branch string, force bool) error {
	return oldgit.DeleteBranch(branch, force)
}

// ListWorktrees returns all worktrees
func (r *Repository) ListWorktrees() ([]worktree.Worktree, error) {
	return oldgit.ListWorktreesTui()
}

// BranchFromWorktree returns the branch name for a worktree path
func (r *Repository) BranchFromWorktree(path string) (string, error) {
	return oldgit.BranchFromWorktree(path)
}

// RepoRoot returns the repository root path
func (r *Repository) RepoRoot() (string, error) {
	return oldgit.RepoRoot()
}
