package git

import (
	"bytes"
	"os"
	"os/exec"
)

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

// RemoveWorktree removes the worktree at the specified path. If force is true, it performs a force removal
func RemoveWorktree(path string, force bool) error {
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

// DeleteBranch deletes the specified branch. If force is true, it uses -D flag for force deletion
func DeleteBranch(branch string, force bool) error {
	flag := "-d"

	if force {
		flag = "-D"
	}

	cmd := exec.Command("git", "branch", flag, branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
