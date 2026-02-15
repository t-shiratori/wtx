package git

import (
	"fmt"
	"os/exec"
	"strings"

	"wtx/internal/domain/worktree"
)

// ListWorktreesTui executes git worktree list --porcelain and returns a list of worktrees
func ListWorktreesTui() ([]worktree.Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list failed: %w", err)
	}

	return parsePorcelain(out)
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
