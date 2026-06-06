# Architecture

## Directory Structure

```
wtx/
├── main.go                        # Entry point — calls cmd.Execute() only
├── cmd/                           # CLI command definitions
│   ├── root.go                    # Root command and shared pre-run logic
│   ├── add.go                     # wtx add
│   ├── remove.go                  # wtx remove
│   ├── list.go                    # wtx list
│   ├── purge.go                   # wtx purge
│   └── init.go                    # wtx init
└── internal/                      # Packages not importable from outside this module
    ├── app/                       # Shared application state passed across commands
    ├── config/                    # Config loading and path resolution
    ├── domain/                    # Domain types and business rules
    │   ├── worktree/              # Worktree entity and AddOptions
    │   └── errors/                # Domain error definitions
    ├── port/                      # Interface definitions (abstractions)
    ├── adapter/                   # Concrete implementations of ports
    │   ├── git/                   # git CLI wrapper
    │   ├── fs/                    # OS filesystem operations
    │   ├── hook/                  # Hook execution
    │   └── tui/                   # Bubbletea TUI for interactive selection
    ├── usecase/                   # Business logic
    │   └── plan/                  # Dry-run plan building (no side effects)
    └── shared/                    # Cross-cutting utilities
        ├── logger/
        ├── spinner/
        ├── prompt/
        └── version/
```

## Architecture Overview

This project follows **Clean Architecture** (also known as Hexagonal Architecture). Dependencies always point inward — outer layers depend on inner layers, never the reverse.

```
┌─────────────────────────────────────────────────────┐
│  cmd/  (outermost: entry points)                    │
│                                                     │
│  add.go / remove.go / list.go / purge.go            │
│   ↓ assembles adapters and injects them into usecase│
└──────────────────────┬──────────────────────────────┘
                       │ depends on
┌──────────────────────▼──────────────────────────────┐
│  usecase/  (business logic)                         │
│                                                     │
│  AddWorktree / RemoveWorktree / ListWorktrees        │
│   ↓ accesses the outside world via port interfaces  │
└──────────────────────┬──────────────────────────────┘
                       │ depends on (via interfaces)
┌──────────────────────▼──────────────────────────────┐
│  port/  (abstractions: interface definitions)       │
│                                                     │
│  GitRepository / FileSystem / HookRunner            │
└────────┬──────────────────────────┬─────────────────┘
         │ implements                │ implements
┌────────▼────────┐      ┌──────────▼──────────┐
│  adapter/git/   │      │  adapter/fs/         │
│  adapter/hook/  │      │  adapter/tui/        │
│  (concrete)     │      │  (concrete)          │
└─────────────────┘      └──────────────────────┘

* domain/ holds the core types referenced by all layers
```

## Layer Responsibilities

### `cmd/` — Entry point and wiring

- Defines subcommands using Cobra
- Reads flags, instantiates adapters, and injects them into the usecase
- Contains no business logic — only calls the usecase

```go
// Typical pattern in cmd/add.go
gitRepo := git.NewRepository()           // create adapter
uc := usecase.NewAddWorktree(gitRepo, .) // inject into usecase
uc.Execute(input)                        // run
```

### `internal/domain/` — What the app deals with

- Core types such as `Worktree` (path, branch, commit)
- Business rules such as `ResolveBaseBranch`
- No external dependencies; referenced by all other layers

### `internal/port/` — Contracts between layers

- Defines interfaces: `GitRepository`, `FileSystem`, `HookRunner`
- The usecase layer only knows these interfaces, not the concrete implementations

```go
// port/git.go
type GitRepository interface {
    AddWorktree(opts worktree.AddOptions) error
    ListWorktrees() ([]worktree.Worktree, error)
    // ...
}
```

### `internal/adapter/` — How things are done

| Package | Implements | Mechanism |
|---|---|---|
| `adapter/git/` | `GitRepository` | Shells out to the `git` CLI |
| `adapter/fs/` | `FileSystem` | Uses the `os` package |
| `adapter/hook/` | `HookRunner` | Executes shell scripts |
| `adapter/tui/` | — (standalone) | Bubbletea TUI |

### `internal/usecase/` — What the app does

- Implements `AddWorktree`, `RemoveWorktree`, and `ListWorktrees`
- Receives port interfaces, so concrete adapters can be swapped with mocks in tests
- `plan/` provides dry-run output without executing any side effects

### `internal/app/` — Shared state across commands

- `PersistentPreRunE` in `cmd/root.go` builds an `app.Context` (config + repo root) and stores it in Cobra's context so every subcommand can retrieve it

## Request Flow: `wtx add <branch>`

```
main.go
  └─ cmd.Execute()
       └─ rootCmd.PersistentPreRunE     ← resolve repo root, load config, store app.Context
            └─ cmd/add.go runAdd()
                 ├─ retrieve app.Context
                 ├─ instantiate adapters (git, fs, hook)
                 ├─ usecase.NewAddWorktree(...) — dependency injection
                 └─ uc.Execute(input)
                      ├─ domain: resolve base branch
                      ├─ port.HookRunner: pre-create hook
                      ├─ port.GitRepository: git worktree add
                      ├─ port.HookRunner: post-create hook
                      ├─ port.FileSystem: copy files
                      └─ port.HookRunner: post-copy hook
```

## Why This Design

The primary motivation is **testability**. Because the usecase layer only depends on port interfaces, tests can swap in mocks without running any real `git` commands.

```go
// Typical test setup in usecase/add_worktree_test.go
uc := usecase.NewAddWorktree(
    &mockGitRepo{}, // no real git execution
    &mockFS{},
    &mockHook{},
    cfg,
    "/repo",
)
```

---

## Port and Adapter in Depth

### port — Contracts ("what can be done")

`port` defines only interfaces — no implementation. There are three interfaces in this project:

**`port.GitRepository`** (`port/git.go`)

```go
type GitRepository interface {
    AddWorktree(opts worktree.AddOptions) error
    RemoveWorktree(path string, force bool) error
    DeleteBranch(branch string, force bool) error
    ListWorktrees() ([]worktree.Worktree, error)
    BranchFromWorktree(path string) (string, error)
    RepoRoot() (string, error)
}
```

**`port.FileSystem`** (`port/filesystem.go`)

```go
type FileSystem interface {
    CopyFile(src, dst string) error
    RemoveIfEmpty(dir string) error
    RemoveEmptyParents(dir, stopDir string) error
    EnsureDir(path string) error
}
```

**`port.HookRunner`** (`port/hook.go`)

```go
type HookRunner interface {
    Run(commands []string, dir string) error
}
```

---

### adapter — Implementations ("how it's done")

| Package | Implements | Mechanism |
|---|---|---|
| `adapter/git/` | `GitRepository` | `exec.Command("git", ...)` |
| `adapter/fs/` | `FileSystem` | `os.ReadFile`, `os.WriteFile`, etc. |
| `adapter/hook/` | `HookRunner` | `exec.Command("sh", "-c", command)` |
| `adapter/tui/` | — (standalone) | Bubbletea TUI for interactive selection |

**`adapter/git`** shells out to the `git` CLI for every operation. `ListWorktrees` parses the output of `git worktree list --porcelain` into `[]worktree.Worktree`.

**`adapter/fs`** wraps the standard `os` package. `CopyFile` reads the source file and writes it to the destination, creating intermediate directories automatically with `os.MkdirAll`.

**`adapter/hook`** executes each configured shell command via `sh -c` in a specified working directory. Used for `pre_create`, `post_create`, and `post_copy` hooks defined in config.

---

### Relationship diagram

```
port (interfaces)              adapter (implementations)
┌─────────────────────┐       ┌──────────────────────────┐
│ GitRepository       │◄──────│ adapter/git/             │
│  AddWorktree()      │       │  exec.Command("git", .)  │
│  ListWorktrees()    │       └──────────────────────────┘
│  RepoRoot()  ...    │
└─────────────────────┘
┌─────────────────────┐       ┌──────────────────────────┐
│ FileSystem          │◄──────│ adapter/fs/              │
│  CopyFile()         │       │  os.ReadFile()           │
│  RemoveIfEmpty()    │       │  os.WriteFile()          │
│  EnsureDir()  ...   │       └──────────────────────────┘
└─────────────────────┘
┌─────────────────────┐       ┌──────────────────────────┐
│ HookRunner          │◄──────│ adapter/hook/            │
│  Run()              │       │  exec.Command("sh", "-c")│
└─────────────────────┘       └──────────────────────────┘
         ▲
         │ usecase only knows these interfaces
┌────────┴────────────┐
│ usecase/            │
│  AddWorktree        │
│  RemoveWorktree     │
└─────────────────────┘
```

---

### Effect on testing

Because the usecase receives port interfaces rather than concrete types, tests replace adapters with hand-written mocks — no real `git` commands, no filesystem writes.

```go
type mockGitRepo struct{}

func (m *mockGitRepo) AddWorktree(opts worktree.AddOptions) error { return nil }
func (m *mockGitRepo) ListWorktrees() ([]worktree.Worktree, error) {
    return []worktree.Worktree{{Path: "/fake", Branch: "main"}}, nil
}
// ... remaining methods

uc := usecase.NewAddWorktree(&mockGitRepo{}, &mockFS{}, &mockHook{}, cfg, "/repo")
```

| | port | adapter |
|---|---|---|
| What it contains | interfaces only | concrete implementations |
| Dependencies | domain only | os, exec, external libs |
| Change frequency | rarely | when the underlying mechanism changes |
| Role in tests | used as the mock's type | used only in production |
