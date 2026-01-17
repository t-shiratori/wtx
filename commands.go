package main

import (
	"fmt"
)

func runList() {
	if err := git("worktree", "list"); err != nil {
		fmt.Println("failed to list worktrees")
	}
}

func runAdd(branch string) {
	path := "../" + branch

	if err := git("worktree", "add", path, branch); err != nil {
		fmt.Println("failed to add worktree")
	}
}

func runRemove(path string) {
	if err := git("worktree", "remove", path); err != nil {
		fmt.Println("failed to remove worktree")
	}
}
