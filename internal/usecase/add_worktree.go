package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"wtx/internal/config"
	"wtx/internal/domain/worktree"
	"wtx/internal/port"
)

// AddWorktreeInput contains the input parameters for AddWorktree usecase
type AddWorktreeInput struct {
	Branch       string
	CreateBranch bool
	FromBranch   string
}

// AddWorktreeOutput contains the result of AddWorktree usecase
type AddWorktreeOutput struct {
	WorktreePath string
	Branch       string
	BaseBranch   string
}

// AddWorktree handles the add worktree usecase
type AddWorktree struct {
	gitRepo    port.GitRepository
	fs         port.FileSystem
	hookRunner port.HookRunner
	config     *config.Config
	repoRoot   string
}

// NewAddWorktree creates a new AddWorktree usecase
func NewAddWorktree(
	gitRepo port.GitRepository,
	fs port.FileSystem,
	hookRunner port.HookRunner,
	cfg *config.Config,
	repoRoot string,
) *AddWorktree {
	return &AddWorktree{
		gitRepo:    gitRepo,
		fs:         fs,
		hookRunner: hookRunner,
		config:     cfg,
		repoRoot:   repoRoot,
	}
}

// Execute runs the add worktree usecase
func (u *AddWorktree) Execute(input AddWorktreeInput) (*AddWorktreeOutput, error) {
	// Resolve base branch
	baseBranch := worktree.ResolveBaseBranch(input.FromBranch, u.config.Worktree.DefaultBaseBranch)

	// Resolve worktree path
	worktreePath := config.ResolveWorktreePath(u.repoRoot, u.config, input.Branch)

	output := &AddWorktreeOutput{
		WorktreePath: worktreePath,
		Branch:       input.Branch,
		BaseBranch:   baseBranch,
	}

	// Execute pre-create hook
	if err := u.hookRunner.Run(u.config.Hooks.PreCreate, u.repoRoot); err != nil {
		return nil, err
	}

	// Add worktree
	opts := worktree.AddOptions{
		Branch:       input.Branch,
		Path:         worktreePath,
		CreateBranch: input.CreateBranch,
		BaseBranch:   baseBranch,
	}
	if err := u.gitRepo.AddWorktree(opts); err != nil {
		return nil, err
	}

	// Execute post-create hook
	if err := u.hookRunner.Run(u.config.Hooks.PostCreate, worktreePath); err != nil {
		return nil, err
	}

	// Copy files - patterns
	for _, pattern := range u.config.Copy.Patterns {
		if err := u.expandAndCopyPattern(pattern, u.repoRoot, worktreePath); err != nil {
			return nil, fmt.Errorf("copy pattern %s: %w", pattern, err)
		}
	}

	// Copy files - explicit files (for renaming)
	for _, file := range u.config.Copy.Files {
		if file.From == "" || file.To == "" {
			continue
		}
		src := filepath.Join(u.repoRoot, file.From)
		dst := filepath.Join(worktreePath, file.To)
		if err := u.fs.CopyFile(src, dst); err != nil {
			return nil, fmt.Errorf("copy %s to %s: %w", file.From, file.To, err)
		}
	}

	// Execute post-copy hook
	if err := u.hookRunner.Run(u.config.Hooks.PostCopy, worktreePath); err != nil {
		return nil, err
	}

	return output, nil
}

// expandAndCopyPattern expands glob pattern and copies matched files
func (u *AddWorktree) expandAndCopyPattern(pattern, repoRoot, worktreePath string) error {
	// Expand glob pattern based on repoRoot
	fullPattern := filepath.Join(repoRoot, pattern)
	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return fmt.Errorf("invalid glob pattern: %w", err)
	}

	// If no matches, just return (not an error - file might not exist in all environments)
	if len(matches) == 0 {
		return nil
	}

	for _, match := range matches {
		// Skip directories
		info, err := os.Stat(match)
		if err != nil {
			return err
		}
		if info.IsDir() {
			continue
		}

		// Get relative path from repoRoot
		relPath, err := filepath.Rel(repoRoot, match)
		if err != nil {
			return err
		}

		// Copy to the same relative path in worktree
		dst := filepath.Join(worktreePath, relPath)
		if err := u.fs.CopyFile(match, dst); err != nil {
			return fmt.Errorf("copy %s: %w", relPath, err)
		}
	}

	return nil
}
