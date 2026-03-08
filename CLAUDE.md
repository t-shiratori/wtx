# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o wtx .

# Install locally
go install .

# Run all tests
go test ./...

# Run a specific test
go test ./internal/usecase/... -run TestAddWorktree_Execute

# Run tests in a specific package
go test ./internal/usecase/...
go test ./cmd/...
```

## Architecture

`wtx` follows a clean/hexagonal architecture. The layers are:

```
cmd/          → Cobra CLI commands (entry points)
internal/
  app/        → App-wide context type (Config + RepoRoot) passed via cobra context
  domain/     → Core domain types (Worktree, AddOptions, errors)
  port/       → Interface definitions (GitRepository, FileSystem, HookRunner)
  adapter/    → Concrete implementations of ports
    git/      → git CLI wrapper (implements port.GitRepository)
    fs/       → OS filesystem ops (implements port.FileSystem)
    tui/      → Bubbletea TUI for interactive worktree selection
  usecase/    → Business logic (AddWorktree, RemoveWorktree, ListWorktrees)
    plan/     → Dry-run plan building (mirrors usecase logic without side effects)
  config/     → Config struct, loading (TOML), path resolution
  shared/     → Logger
```

### Request flow

1. `cmd/root.go` `PersistentPreRunE` runs first: detects repo root, loads config (local + global merged), stores `app.Context` in cobra's `context.Context`
2. Each command retrieves `app.Context`, instantiates adapters, constructs the usecase, and calls `Execute()`
3. Dry-run (`--dry-run`) uses the corresponding `plan/` package instead of the usecase

### Configuration

- Local: `.wtx/config.toml` (project root)
- Global: `~/.config/wtx/config.toml`
- Both are merged on load; local values take priority for scalar fields, slice fields are combined

### Testing

Tests use hand-written mocks (no mock generation framework). Mock types for `port.GitRepository`, `port.FileSystem`, and `port.HookRunner` are defined inline in `*_test.go` files.
