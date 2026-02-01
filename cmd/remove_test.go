package cmd

import (
	"bytes"
	"context"
	"wtx/internal/app"
	"wtx/internal/config"
	"testing"
)

func TestRemoveCommand_DryRun(t *testing.T) {

	cfg := &config.Config{}
	appCtx := &app.Context{
		Config:   cfg,
		RepoRoot: "/repo",
	}

	ctx := context.WithValue(context.Background(), app.Key, appCtx)

	buf := new(bytes.Buffer)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{
		"remove",
		"feature/test",
		"--dry-run",
	})

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !contains(out, "Dry run") {
		t.Fatalf("expected dry run output, got: \n%s", out)
	}

}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
