package usecase

import (
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

	// Copy files
	for _, c := range u.config.Copy {
		if c.From == "" || c.To == "" {
			continue
		}
		src := filepath.Join(u.repoRoot, c.From)
		dst := filepath.Join(worktreePath, c.To)
		if err := u.fs.CopyFile(src, dst); err != nil {
			return nil, err
		}
	}

	// Execute post-copy hook
	if err := u.hookRunner.Run(u.config.Hooks.PostCopy, worktreePath); err != nil {
		return nil, err
	}

	return output, nil
}
