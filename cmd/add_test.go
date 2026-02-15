package cmd

import (
	"bytes"
	"context"
	"strings"
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
	ctx := testAppContext(t)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"add", "--dry-run", "feature/test"})

	err := rootCmd.ExecuteContext(ctx)

	if err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	// Verify dry-run output contains expected information
	output := out.String()
	if !strings.Contains(output, "Dry run") {
		t.Errorf("expected dry-run output, got: %s", output)
	}
}

func TestAddCommand_DryRunWithFromBranch(t *testing.T) {
	ctx := testAppContext(t)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"add", "--dry-run", "--from", "develop", "feature/test"})

	err := rootCmd.ExecuteContext(ctx)

	if err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "develop") {
		t.Errorf("expected 'develop' in output, got: %s", output)
	}
}
