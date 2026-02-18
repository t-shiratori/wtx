package cmd

import (
	"wtx/internal/adapter/fs"
	"wtx/internal/adapter/git"
	"wtx/internal/adapter/hook"
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/domain/worktree"
	"wtx/internal/shared/logger"
	"wtx/internal/shared/spinner"
	"wtx/internal/usecase/plan"
	"wtx/internal/usecase"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "add git worktree",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	appCtx := cmd.Context().Value(app.Key).(*app.Context)
	cfg := appCtx.Config
	repoRoot := appCtx.RepoRoot

	// Get flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	createBranch, _ := cmd.Flags().GetBool("create-branch")
	fromBranch, _ := cmd.Flags().GetString("from")

	branch := args[0]

	// Handle dry-run in cmd layer
	if dryRun {
		baseBranch := worktree.ResolveBaseBranch(fromBranch, cfg.Worktree.DefaultBaseBranch)
		worktreeDir := config.ResolveWorktreePath(repoRoot, cfg, branch)

		addPlan := plan.NewAddPlanFromConfig(
			cfg,
			branch,
			worktreeDir,
			createBranch,
			baseBranch,
		)
		addPlan.Print(cmd.OutOrStdout())
		return nil
	}

	// Create dependencies
	gitRepo := git.NewRepository()
	fsAdapter := fs.NewFileSystem()
	hookRunner := hook.NewRunner()

	sp := spinner.New(cmd.OutOrStdout(), "Adding worktree...")

	// Create and execute usecase
	uc := usecase.NewAddWorktree(
		gitRepo,
		fsAdapter,
		newHookWithSpinner(hookRunner, sp),
		cfg,
		repoRoot,
	)
	sp.Start()

	output, err := uc.Execute(usecase.AddWorktreeInput{
		Branch:       branch,
		CreateBranch: createBranch,
		FromBranch:   fromBranch,
	})

	if err != nil {
		sp.StopWithError()
		return err
	}
	sp.StopWithSuccess()

	logger.Success(cmd.OutOrStdout(), "Added worktree '%s'", output.WorktreePath)
	return nil
}

func init() {
	addCmd.Flags().BoolP("create-branch", "b", false, "create branch if it does not exist")
	addCmd.Flags().String("from", "", "base branch or commit")
	addCmd.Flags().Bool("dry-run", false, "show what would be done")
	rootCmd.AddCommand(addCmd)
}
