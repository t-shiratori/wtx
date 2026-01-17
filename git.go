package main

import (
	"fmt"
	"os"
	"os/exec"
)

func git(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensureGitRepo() {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		fmt.Println("error: not a git repository")
		os.Exit(1)
	}
}
