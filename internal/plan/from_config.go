package plan

import "wtx/internal/config"

func NewAddPlanFromConfig(
	cfg *config.Config,
	branch string,
	worktreePath string,
	createBranch bool,
	baseBranch string,
) *AddPlan {

	plan := &AddPlan{
		Branch:       branch,
		WorktreePath: worktreePath,
		CreateBranch: createBranch,
		BaseBranch:   baseBranch,
	}

	// --- hooks ---
	plan.PreHook = cfg.Hooks.PreCreate
	plan.PostHook = cfg.Hooks.PostCreate
	plan.PostCopyHook = cfg.Hooks.PostCopy

	// --- copy ---
	// Patterns
	for _, pattern := range cfg.Copy.Patterns {
		plan.CopyFiles = append(plan.CopyFiles, pattern)
	}

	// Explicit files (from -> to)
	for _, file := range cfg.Copy.Files {
		// Combine from -> to into a single line for display purposes
		plan.CopyFiles = append(
			plan.CopyFiles,
			file.From+" -> "+file.To,
		)
	}

	return plan
}
