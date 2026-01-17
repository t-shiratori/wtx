package main

import (
	"fmt"
	"os"
)

func main() {
	ensureGitRepo()

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "list":
		runList()

	case "add":
		if len(args) < 2 {
			fmt.Println("branch name required")
			os.Exit(1)
		}
		runAdd(args[1])

	case "remove":
		if len(args) < 2 {
			fmt.Println("path required")
			os.Exit(1)
		}
		runRemove(args[1])

	default:
		fmt.Println("unknown command:", args[0])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`usage: wt <command> [arguments]

commands:
  list              List all worktrees
  add <branch>      Add a new worktree for the specified branch
  remove <path>     Remove the worktree at the specified path`)
}
