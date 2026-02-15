package usecase

import (
	"wtx/internal/config"
	"wtx/internal/port"
)

// RemoveWorktreeInput contains the input parameters for RemoveWorktree usecase
type RemoveWorktreeInput struct {
	Paths        []string
	DeleteBranch bool
	Force        bool
}

// RemoveWorktreeResult contains the result for a single worktree removal
type RemoveWorktreeResult struct {
	Path    string
	Branch  string
	Success bool
	Error   error
}

// RemoveWorktreeOutput contains the results of RemoveWorktree usecase
type RemoveWorktreeOutput struct {
	Results []RemoveWorktreeResult
}

// RemoveWorktree handles the remove worktree usecase
type RemoveWorktree struct {
	gitRepo  port.GitRepository
	fs       port.FileSystem
	config   *config.Config
	repoRoot string
}

// NewRemoveWorktree creates a new RemoveWorktree usecase
func NewRemoveWorktree(
	gitRepo port.GitRepository,
	fs port.FileSystem,
	cfg *config.Config,
	repoRoot string,
) *RemoveWorktree {
	return &RemoveWorktree{
		gitRepo:  gitRepo,
		fs:       fs,
		config:   cfg,
		repoRoot: repoRoot,
	}
}

// Execute runs the remove worktree usecase
func (u *RemoveWorktree) Execute(input RemoveWorktreeInput) (*RemoveWorktreeOutput, error) {
	output := &RemoveWorktreeOutput{
		Results: make([]RemoveWorktreeResult, 0, len(input.Paths)),
	}

	worktreesRoot := config.ResolveWorktreesDir(u.repoRoot, u.config)

	for _, inputPath := range input.Paths {
		result := u.removeOne(inputPath, input, worktreesRoot)
		output.Results = append(output.Results, result)
	}

	return output, nil
}

func (u *RemoveWorktree) removeOne(inputPath string, input RemoveWorktreeInput, worktreesRoot string) RemoveWorktreeResult {
	result := RemoveWorktreeResult{Path: inputPath}

	// Resolve path
	path, err := config.ResolveInputWorktreePath(u.repoRoot, u.config, inputPath)
	if err != nil {
		result.Error = err
		return result
	}
	result.Path = path

	// Get branch if needed
	var branch string
	if input.DeleteBranch {
		branch, err = u.gitRepo.BranchFromWorktree(path)
		if err != nil {
			result.Error = err
			return result
		}
		result.Branch = branch
	}

	// Remove worktree
	if err := u.gitRepo.RemoveWorktree(path, input.Force); err != nil {
		result.Error = err
		return result
	}

	// Delete branch
	if input.DeleteBranch && branch != "" {
		if err := u.gitRepo.DeleteBranch(branch, input.Force); err != nil {
			result.Error = err
			return result
		}
	}

	// Cleanup empty directories
	_ = u.fs.RemoveIfEmpty(path)
	_ = u.fs.RemoveEmptyParents(path, worktreesRoot)

	result.Success = true
	return result
}

// FailedCount returns the number of failed removals
func (o *RemoveWorktreeOutput) FailedCount() int {
	count := 0
	for _, r := range o.Results {
		if !r.Success {
			count++
		}
	}
	return count
}
