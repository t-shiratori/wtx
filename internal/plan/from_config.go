package plan

import "go-worktree-cli/internal/config"

func NewAddPlanFromConfig(
	cfg *config.Config,
	branch string,
	worktreePath string,
) *AddPlan {

	plan := &AddPlan{
		Branch:       branch,
		WorktreePath: worktreePath,
	}

	// --- hooks ---
	plan.PreHook = cfg.Hooks.PreCreate
	plan.PostHook = cfg.Hooks.PostCreate
	plan.AfterCopyHook = cfg.Hooks.PostCopy

	// --- copy ---
	for _, c := range cfg.Copy {
		// 表示用なので from -> to を 1 行にまとめる
		plan.CopyFiles = append(
			plan.CopyFiles,
			c.From+" -> "+c.To,
		)
	}

	return plan
}
