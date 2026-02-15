// Tests for the remove command
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

// TestRemoveCmd_DryRun verifies that --dry-run flag works correctly
func TestRemoveCmd_DryRun(t *testing.T) {
	repoRoot := t.TempDir()

	cfg := &config.Config{}
	ctx := context.WithValue(
		context.Background(),
		app.Key,
		&app.Context{
			Config:   cfg,
			RepoRoot: repoRoot,
		},
	)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"remove", "--dry-run", "test-branch"})

	err := rootCmd.ExecuteContext(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify dry-run output
	output := out.String()
	if !strings.Contains(output, "Dry run") {
		t.Errorf("expected dry-run output, got: %s", output)
	}
}

// TestRemoveCmd_NoArgs_NoSelection verifies TUI selection flow
func TestRemoveCmd_NoArgs_NoSelection(t *testing.T) {
	// Mock TUI functions to simulate no selection
	origList := listWorktreesTui
	origSelect := selectWorktrees
	defer func() {
		listWorktreesTui = origList
		selectWorktrees = origSelect
	}()

	listWorktreesTui = func() ([]worktree.Worktree, error) {
		return []worktree.Worktree{
			{Path: "wt1"},
			{Path: "wt2"},
		}, nil
	}

	selectWorktrees = func(_ []worktree.Worktree) ([]string, error) {
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
