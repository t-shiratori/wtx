package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"wtx/internal/domain/worktree"
)

// CommandError wraps a command execution error with its stderr output
type CommandError struct {
	Err     error
	Message string
}

func (e *CommandError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return e.Err.Error()
}

// Repository implements port.GitRepository
type Repository struct{}

// NewRepository creates a new Repository
func NewRepository() *Repository {
	return &Repository{}
}

// AddWorktree creates a new worktree
func (r *Repository) AddWorktree(opts worktree.AddOptions) error {
	args := []string{"worktree", "add"}

	if opts.CreateBranch {
		args = append(args, "-b", opts.Branch)
		args = append(args, opts.Path, opts.BaseBranch)
	} else {
		args = append(args, opts.Path, opts.Branch)
	}

	cmd := exec.Command("git", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return &CommandError{
			Err:     err,
			Message: stderr.String(),
		}
	}

	return nil
}

// RemoveWorktree removes an existing worktree
func (r *Repository) RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}

	if force {
		args = append(args, "--force")
	}

	args = append(args, path)

	cmd := exec.Command("git", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return &CommandError{
			Err:     err,
			Message: stderr.String(),
		}
	}

	return nil
}

// DeleteBranch deletes a git branch
func (r *Repository) DeleteBranch(branch string, force bool) error {
	flag := "-d"

	if force {
		flag = "-D"
	}

	cmd := exec.Command("git", "branch", flag, branch)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return &CommandError{
			Err:     err,
			Message: stderr.String(),
		}
	}

	return nil
}

// ListWorktrees returns all worktrees
func (r *Repository) ListWorktrees() ([]worktree.Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list failed: %w", err)
	}

	return parsePorcelain(out)
}

// BranchFromWorktree returns the branch name for a worktree path
func (r *Repository) BranchFromWorktree(path string) (string, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")

	var currentPath string
	for _, line := range lines {
		if rest, ok := strings.CutPrefix(line, "worktree "); ok {
			currentPath = rest
		}

		if currentPath == path {
			if rest, ok := strings.CutPrefix(line, "branch "); ok {
				branch := strings.TrimPrefix(rest, "refs/heads/")
				return branch, nil
			}
		}
	}

	return "", errors.New("branch not found for worktree")
}

// RepoRoot returns the repository root path
func (r *Repository) RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("not a git repository")
	}

	return strings.TrimSpace(string(out)), nil
}

// parsePorcelain parses the output of git worktree list --porcelain
func parsePorcelain(data []byte) ([]worktree.Worktree, error) {
	var result []worktree.Worktree
	var current *worktree.Worktree

	for line := range strings.SplitSeq(string(data), "\n") {
		if len(line) == 0 {
			continue
		}

		fields := strings.SplitN(line, " ", 2)
		key := fields[0]

		switch key {
		case "worktree":
			if current != nil {
				result = append(result, *current)
			}
			if len(fields) < 2 {
				continue
			}
			current = &worktree.Worktree{
				Path: fields[1],
			}
		case "branch":
			if current != nil && len(fields) >= 2 {
				current.Branch = strings.TrimPrefix(fields[1], "refs/heads/")
			}
		case "HEAD":
			if current != nil && len(fields) >= 2 {
				current.Commit = fields[1]
			}
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	return result, nil
}
