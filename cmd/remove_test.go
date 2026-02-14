// Tests for the remove command
package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/git"

	"github.com/spf13/cobra"
)

// TestRemoveCmd_DryRun verifies that the actual removal is not executed
// when the --dry-run flag is specified
func TestRemoveCmd_DryRun(t *testing.T) {
	// Create a temporary directory for testing
	repoRoot := t.TempDir()

	// Set up the test context
	cfg := &config.Config{}
	ctx := context.WithValue(
		context.Background(),
		app.Key,
		&app.Context{
			Config:   cfg,
			RepoRoot: repoRoot,
		},
	)

	// Set up a buffer to capture command output
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetContext(ctx)
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	// Enable dry-run mode and disable other flags
	dryRun = true
	removeBranch = false
	force = false

	// Mock gitRemoveWorktree to track if it gets called
	called := false
	gitRemoveWorktree = func(path string, force bool) error {
		called = true
		return nil
	}

	// Restore the original function after the test
	defer func() {
		gitRemoveWorktree = git.RemoveWorktree
	}()

	// Execute the command
	err := removeCmd.RunE(cmd, []string{"test-branch"})

	// Verify no error occurred
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that the actual removal was not called in dry-run mode
	if called {
		t.Fatalf("gitRemoveWorktree should not be called in dry-run")
	}

	// Verify that output was generated
	if out.Len() == 0 {
		t.Fatalf("expected output, got empty")
	}
}

// TestRemoveCmd_RemoveWorktreeCalled verifies that gitRemoveWorktree
// is correctly called during normal execution
func TestRemoveCmd_RemoveWorktreeCalled(t *testing.T) {
	// Create a temporary directory for testing
	repoRoot := t.TempDir()

	// Set up the test context
	cfg := &config.Config{}
	ctx := context.WithValue(
		context.Background(),
		app.Key,
		&app.Context{
			Config:   cfg,
			RepoRoot: repoRoot,
		},
	)

	// Set up a buffer to capture command output
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetContext(ctx)
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	// Test in normal mode (dry-run disabled)
	dryRun = false
	removeBranch = false
	force = false

	// Mock gitRemoveWorktree to track if it gets called
	called := false
	gitRemoveWorktree = func(path string, force bool) error {
		called = true
		return nil
	}

	// Restore the original function after the test
	defer func() {
		gitRemoveWorktree = git.RemoveWorktree
	}()

	// Execute the command
	err := removeCmd.RunE(cmd, []string{"test-branch"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that gitRemoveWorktree was called
	if !called {
		t.Fatalf("gitRemoveWorktree was not called")
	}
}

func TestRemoveCmd_NoArgs_NoSelection(t *testing.T) {
	// Mock TUI functions to simulate no selection
	origList := listWorktreesTui
	origSelect := selectWorktrees
	defer func() {
		listWorktreesTui = origList
		selectWorktrees = origSelect
	}()

	listWorktreesTui = func() ([]git.Worktree, error) {
		return []git.Worktree{
			{Path: "wt1"},
			{Path: "wt2"},
		}, nil
	}

	selectWorktrees = func(_ []git.Worktree) ([]string, error) {
		return []string{}, nil
	}

	ctx := context.WithValue(
		context.Background(),
		app.Key,
		&app.Context{
			Config:   &config.Config{},
			RepoRoot: t.TempDir(),
		},
	)

	buf := new(bytes.Buffer)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"remove"})

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "No worktrees selected") {
		t.Fatalf("expected warning message, got: %s", output)
	}
}
