package plan

import (
	"fmt"
	"go-worktree-cli/internal/ui"
	"io"
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

func (p *AddPlan) Print(w io.Writer) {
	ui.Info(w, "Dry run: add worktree")

	fmt.Fprintln(w, "  - branch:  ", p.Branch)
	fmt.Fprintln(w, "  - path:    ", p.WorktreePath)
	fmt.Fprintln(w, "  - from:", p.BaseBranch)
	if p.CreateBranch {
		ui.Info(w, "  - create-branch: %s", p.Branch)
	}

	if len(p.PreHook) > 0 {
		ui.Info(w, "  - pre-hook:")
		for _, h := range p.PreHook {
			ui.Info(w, "    - %s", h)
		}
	}

	if len(p.PostHook) > 0 {
		ui.Info(w, "  - post-hook:")
		for _, h := range p.PostHook {
			ui.Info(w, "    - %s", h)
		}
	}

	if len(p.CopyFiles) > 0 {
		ui.Info(w, "  - copy:")
		for _, f := range p.CopyFiles {
			ui.Info(w, "    - %s", f)
		}
	}

	if len(p.PostCopyHook) > 0 {
		ui.Info(w, "  - post-copy-hook:")
		for _, h := range p.PostCopyHook {
			ui.Info(w, "    - %s", h)
		}
	}
}
