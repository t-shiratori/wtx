package cmd

import (
	"context"
	"testing"

	"wtx/internal/app"
	"wtx/internal/config"
)

func testAppContext(t *testing.T) context.Context {
	t.Helper()

	cfg := &config.Config{
		Worktree: config.WorktreeConfig{
			RootDir:           ".wt/worktrees",
			DefaultBaseBranch: "main",
		},
	}

	appCtx := &app.Context{
		Config:   cfg,
		RepoRoot: "/tmp/repo",
	}

	return context.WithValue(context.Background(), app.Key, appCtx)
}

func TestAddCommand_DryRun(t *testing.T) {
	// フラグの後始末（超重要）
	t.Cleanup(func() {
		addDryRun = false
		createBranch = false
		fromBranch = ""
	})

	// Arrange
	ctx := testAppContext(t)

	cmd := addCmd
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"feature/test"})
	addDryRun = true

	// Act
	err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("add command failed: %v", err)
	}
}

func TestAddCommand_FromBranchOverridesConfig(t *testing.T) {
	t.Cleanup(func() {
		addDryRun = false
		fromBranch = ""
	})

	ctx := testAppContext(t)

	cmd := addCmd
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"feature/test"})
	addDryRun = true
	fromBranch = "develop"

	if err := cmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}
}
