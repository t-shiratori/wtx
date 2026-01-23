package plan

import (
	"fmt"
	"go-worktree-cli/internal/ui"
)

type AddPlan struct {
	Branch       string
	WorktreePath string

	PreHook       []string
	PostHook      []string
	CopyFiles     []string
	AfterCopyHook []string
}

func (p *AddPlan) Print() {
	ui.Info("Dry run: add worktree")

	fmt.Println("  - branch:  ", p.Branch)
	fmt.Println("  - path:    ", p.WorktreePath)

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

	if len(p.AfterCopyHook) > 0 {
		fmt.Println("  - after-copy-hook:")
		for _, h := range p.AfterCopyHook {
			fmt.Println("    -", h)
		}
	}
}
