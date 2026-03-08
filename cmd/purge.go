package cmd

import (
	"fmt"

	"wtx/internal/adapter/fs"
	"wtx/internal/adapter/git"
	"wtx/internal/app"
	"wtx/internal/domain/worktree"
	"wtx/internal/shared/logger"
	"wtx/internal/shared/prompt"
	"wtx/internal/shared/spinner"
	"wtx/internal/usecase"
	"wtx/internal/usecase/plan"

	"github.com/spf13/cobra"
)

var listWorktreesPurge func() ([]worktree.Worktree, error)

func init() {
	repo := git.NewRepository()
	listWorktreesPurge = repo.ListWorktrees
}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove all worktrees at once",
	RunE:  runPurge,
}

func runPurge(cmd *cobra.Command, args []string) error {
	appCtx := cmd.Context().Value(app.Key).(*app.Context)
	cfg := appCtx.Config
	repoRoot := appCtx.RepoRoot

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	removeBranch, _ := cmd.Flags().GetBool("branch")
	force, _ := cmd.Flags().GetBool("force")
	yes, _ := cmd.Flags().GetBool("yes")

	// Fetch all worktrees and exclude the main one (always at index 0)
	worktrees, err := listWorktreesPurge()
	if err != nil {
		return err
	}

	// git worktree list always returns the main worktree first
	var targets []worktree.Worktree
	if len(worktrees) > 1 {
		targets = worktrees[1:]
	}

	if len(targets) == 0 {
		logger.Warn(cmd.OutOrStdout(), "No worktrees to purge")
		return nil
	}

	// dry-run: print what would be removed and exit
	if dryRun {
		gitRepo := git.NewRepository()
		for _, wt := range targets {
			var branch string
			if removeBranch {
				branch, _ = gitRepo.BranchFromWorktree(wt.Path)
			}
			p := plan.RemovePlan{
				WorktreePath: wt.Path,
				Branch:       branch,
				Force:        force,
			}
			p.Print(cmd.OutOrStdout())
		}
		return nil
	}

	// Show targets to be removed
	logger.Info(cmd.OutOrStdout(), "The following worktrees will be removed:")
	for _, wt := range targets {
		logger.Info(cmd.OutOrStdout(), "  - %s (%s)", wt.Path, wt.Branch)
	}

	// Confirmation prompt
	if !yes {
		confirmed, err := prompt.Confirm(cmd.OutOrStdout(), cmd.InOrStdin(), fmt.Sprintf("Remove %d worktree(s)?", len(targets)))
		if err != nil {
			return err
		}
		if !confirmed {
			logger.Warn(cmd.OutOrStdout(), "Aborted")
			return nil
		}
	}

	// Collect paths and pass to usecase
	paths := make([]string, len(targets))
	for i, wt := range targets {
		paths[i] = wt.Path
	}

	gitRepo := git.NewRepository()
	fsAdapter := fs.NewFileSystem()
	uc := usecase.NewRemoveWorktree(gitRepo, fsAdapter, cfg, repoRoot)

	sp := spinner.New(cmd.OutOrStdout(), "Removing worktrees...")
	sp.Start()

	output, err := uc.Execute(usecase.RemoveWorktreeInput{
		Paths:        paths,
		DeleteBranch: removeBranch,
		Force:        force,
	})
	if err != nil {
		sp.StopWithError()
		return err
	}
	sp.StopWithSuccess()

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
	purgeCmd.Flags().BoolP("branch", "b", false, "also delete branch")
	purgeCmd.Flags().BoolP("force", "f", false, "force remove worktree and branch")
	purgeCmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")
	purgeCmd.Flags().Bool("dry-run", false, "show what would be done")

	rootCmd.AddCommand(purgeCmd)
}
