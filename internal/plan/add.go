package plan

import (
	"fmt"
	"go-worktree-cli/internal/ui"
)

type AddPlan struct {
	Branch       string
	BaseBranch   string
	WorktreePath string

	CreateBranch bool

	PreHook      []string
	PostHook     []string
	CopyFiles    []string
	PostCopyHook []string
}

func (p *AddPlan) Print() {
	ui.Info("Dry run: add worktree")

	fmt.Println("  - branch:  ", p.Branch)
	fmt.Println("  - path:    ", p.WorktreePath)
	fmt.Println("  - from:", p.BaseBranch)

	if p.CreateBranch {
		fmt.Println("  - create branch: yes")
	}

	if len(p.PreHook) > 0 {
		fmt.Println("  - pre-hook:")
		for _, h := range p.PreHook {
			fmt.Println("    -", h)
		}
	}

	if len(p.PostHook) > 0 {
		fmt.Println("  - post-hook:")
		for _, h := range p.PostHook {
			fmt.Println("    -", h)
		}
	}

	if len(p.CopyFiles) > 0 {
		fmt.Println("  - copy:")
		for _, f := range p.CopyFiles {
			fmt.Println("    -", f)
		}
	}

	if len(p.PostCopyHook) > 0 {
		fmt.Println("  - post-copy-hook:")
		for _, h := range p.PostCopyHook {
			fmt.Println("    -", h)
		}
	}
}
