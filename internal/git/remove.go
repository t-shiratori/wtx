package git

import (
	"os"
	"os/exec"
)

// RemoveWorktree
// 指定されたパスの worktree を削除する。force が true の場合は強制削除を行う
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}

	if force {
		args = append(args, "--force")
	}

	args = append(args, path)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DeleteBranch
// 指定されたブランチを削除する。force が true の場合は -D フラグで強制削除を行う
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
