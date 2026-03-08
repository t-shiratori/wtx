package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/domain/worktree"
)

func newPurgeCtx(t *testing.T) context.Context {
	t.Helper()
	return context.WithValue(
		context.Background(),
		app.Key,
		&app.Context{
			Config:   &config.Config{},
			RepoRoot: t.TempDir(),
		},
	)
}

// resetPurgeFlags prevents flag state leakage between tests
func resetPurgeFlags() {
	purgeCmd.Flags().Set("dry-run", "false")
	purgeCmd.Flags().Set("branch", "false")
	purgeCmd.Flags().Set("force", "false")
	purgeCmd.Flags().Set("yes", "false")
}

// TestPurgeCmd_NoWorktrees: does nothing when only the main worktree exists
func TestPurgeCmd_NoWorktrees(t *testing.T) {
	resetPurgeFlags()

	orig := listWorktreesPurge
	defer func() { listWorktreesPurge = orig }()
	listWorktreesPurge = func() ([]worktree.Worktree, error) {
		// only the main worktree (first element)
		return []worktree.Worktree{
			{Path: "/repo", Branch: "main"},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"purge", "--yes"})

	err := rootCmd.ExecuteContext(newPurgeCtx(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No worktrees to purge") {
		t.Errorf("expected 'No worktrees to purge', got: %s", buf.String())
	}
}

// TestPurgeCmd_DryRun: only prints what would be removed when --dry-run is set
func TestPurgeCmd_DryRun(t *testing.T) {
	resetPurgeFlags()

	orig := listWorktreesPurge
	defer func() { listWorktreesPurge = orig }()
	listWorktreesPurge = func() ([]worktree.Worktree, error) {
		return []worktree.Worktree{
			{Path: "/repo", Branch: "main"},
			{Path: "/repo/worktrees/feature-a", Branch: "feature-a"},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"purge", "--dry-run"})

	err := rootCmd.ExecuteContext(newPurgeCtx(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Dry run") {
		t.Errorf("expected dry-run output, got: %s", output)
	}
	if !strings.Contains(output, "feature-a") {
		t.Errorf("expected worktree path in output, got: %s", output)
	}
}

// TestPurgeCmd_Confirm_Abort: aborts when "n" is entered at the confirmation prompt
func TestPurgeCmd_Confirm_Abort(t *testing.T) {
	resetPurgeFlags()

	orig := listWorktreesPurge
	defer func() { listWorktreesPurge = orig }()
	listWorktreesPurge = func() ([]worktree.Worktree, error) {
		return []worktree.Worktree{
			{Path: "/repo", Branch: "main"},
			{Path: "/repo/worktrees/feature-a", Branch: "feature-a"},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetIn(strings.NewReader("n\n"))
	rootCmd.SetArgs([]string{"purge"})

	err := rootCmd.ExecuteContext(newPurgeCtx(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Aborted") {
		t.Errorf("expected 'Aborted', got: %s", buf.String())
	}
}

// TestPurgeCmd_Confirm_Proceed: proceeds past the confirmation prompt when "y" is entered
func TestPurgeCmd_Confirm_Proceed(t *testing.T) {
	resetPurgeFlags()

	orig := listWorktreesPurge
	defer func() { listWorktreesPurge = orig }()
	listWorktreesPurge = func() ([]worktree.Worktree, error) {
		return []worktree.Worktree{
			{Path: "/repo", Branch: "main"},
			{Path: "/repo/worktrees/feature-a", Branch: "feature-a"},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetIn(strings.NewReader("y\n"))
	rootCmd.SetArgs([]string{"purge"})

	// The usecase will fail since there is no real git repo,
	// but we verify the prompt was passed by checking that "Aborted" is not in the output.
	_ = rootCmd.ExecuteContext(newPurgeCtx(t))

	if strings.Contains(buf.String(), "Aborted") {
		t.Errorf("expected to proceed past confirmation, but got Aborted: %s", buf.String())
	}
}
