package cmd

import (
	"wtx/internal/adapter/fs"
	"wtx/internal/adapter/git"
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/domain/worktree"
	"wtx/internal/shared/logger"
	"wtx/internal/shared/spinner"
	"wtx/internal/usecase/plan"
	"wtx/internal/adapter/tui"
	"wtx/internal/usecase"

	"github.com/spf13/cobra"
)

// Function variables for testing (TUI only)
var (
	listWorktreesTui func() ([]worktree.Worktree, error)
	selectWorktrees  func([]worktree.Worktree) ([]string, error) = tui.SelectWorktrees
)

func init() {
	repo := git.NewRepository()
	listWorktreesTui = repo.ListWorktrees
}

var removeCmd = &cobra.Command{
	Use:   "remove [worktree ...]",
	Short: "Remove git worktrees",
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	appCtx := cmd.Context().Value(app.Key).(*app.Context)
	cfg := appCtx.Config
	repoRoot := appCtx.RepoRoot

	// Get flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	removeBranch, _ := cmd.Flags().GetBool("branch")
	force, _ := cmd.Flags().GetBool("force")

	// If no worktree argument is provided, use TUI
	if len(args) == 0 {
		worktrees, err := listWorktreesTui()
		if err != nil {
			return err
		}

		selected, err := selectWorktrees(worktrees)
		if err != nil {
			return err
		}

		if len(selected) == 0 {
			logger.Warn(cmd.OutOrStdout(), "No worktrees selected")
			return nil
		}

		args = selected
	}

	// Handle dry-run in cmd layer
	if dryRun {
		gitRepo := git.NewRepository()
		for _, inputPath := range args {
			path, err := config.ResolveInputWorktreePath(repoRoot, cfg, inputPath)
			if err != nil {
				logger.Error(cmd.ErrOrStderr(), "%v", err)
				continue
			}

			var branch string
			if removeBranch {
				branch, _ = gitRepo.BranchFromWorktree(path)
			}

			p := plan.RemovePlan{
				WorktreePath: path,
				Branch:       branch,
				Force:        force,
			}
			p.Print(cmd.OutOrStdout())
		}
		return nil
	}

	// Create dependencies
	gitRepo := git.NewRepository()
	fsAdapter := fs.NewFileSystem()

	// Create and execute usecase
	uc := usecase.NewRemoveWorktree(gitRepo, fsAdapter, cfg, repoRoot)

	sp := spinner.New(cmd.OutOrStdout(), "Removing worktrees...")
	sp.Start()

	output, err := uc.Execute(usecase.RemoveWorktreeInput{
		Paths:        args,
		DeleteBranch: removeBranch,
		Force:        force,
	})

	if err != nil {
		sp.StopWithError()
		return err
	}
	sp.StopWithSuccess()

	// Display results
	for _, result := range output.Results {
		if result.Error != nil {
			logger.Error(cmd.ErrOrStderr(), "%v", result.Error)
		} else {
			logger.Success(cmd.OutOrStdout(), "Removed worktree '%s'", result.Path)
			if result.Branch != "" {
				logger.Success(cmd.OutOrStdout(), "Deleted branch '%s'", result.Branch)
			}
		}
	}

	// Summary warning
	if failed := output.FailedCount(); failed > 0 {
		logger.Warn(
			cmd.ErrOrStderr(),
			"%d worktrees could not be removed (use --force to override)",
			failed,
		)
	}

	return nil
}

func init() {
	removeCmd.Flags().BoolP("branch", "b", false, "also delete branch")
	removeCmd.Flags().BoolP("force", "f", false, "force remove worktree and branch")
	removeCmd.Flags().Bool("dry-run", false, "show what would be done")

	rootCmd.AddCommand(removeCmd)
}
