package git

import (
	"os"
	"os/exec"
)

// AddWorktree
// git worktree add <path> <branch>
func AddWorktree(path string, branch string, createBranch bool, baseBranch string) error {
	args := []string{"worktree", "add"}

	if createBranch {
		args = append(args, "-b", branch)
		args = append(args, path, baseBranch)
	} else {
		args = append(args, path, branch)
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
