package git

import (
	"errors"
	"os/exec"
	"strings"
)

// RepoRoot returns the path to the repository root directory
func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("not a git repository")
	}

	return strings.TrimSpace(string(out)), nil
}

// BranchFromWorktree returns the branch name corresponding to the specified worktree path
func BranchFromWorktree(path string) (string, error) {

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
