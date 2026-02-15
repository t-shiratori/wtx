package usecase

import (
	"wtx/internal/domain/worktree"
	"wtx/internal/port"
)

// ListWorktreesOutput contains the result of ListWorktrees usecase
type ListWorktreesOutput struct {
	Worktrees []worktree.Worktree
}

// ListWorktrees handles the list worktrees usecase
type ListWorktrees struct {
	gitRepo port.GitRepository
}

// NewListWorktrees creates a new ListWorktrees usecase
func NewListWorktrees(gitRepo port.GitRepository) *ListWorktrees {
	return &ListWorktrees{
		gitRepo: gitRepo,
	}
}

// Execute runs the list worktrees usecase
func (u *ListWorktrees) Execute() (*ListWorktreesOutput, error) {
	wts, err := u.gitRepo.ListWorktrees()
	if err != nil {
		return nil, err
	}

	return &ListWorktreesOutput{
		Worktrees: wts,
	}, nil
}
