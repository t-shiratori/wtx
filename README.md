# wtx

A CLI tool for easily managing Git worktrees.

## Overview

wtx is a lightweight CLI tool for efficiently managing Git's worktree feature. It makes it easy to create, delete, and list worktrees when working on multiple branches in parallel.

Key features:
- Create, delete, and list worktrees
- Interactive worktree selection via TUI (Terminal UI)
- Hook functionality (pre_create, post_create, post_copy)
- Automatic file copying
- Configuration file initialization (local/global)
- Automatic cleanup of empty directories after worktree removal

## Installation

### Homebrew

```bash
brew install t-shiratori/tap/wtx
```

### Direct download from GitHub releases

Download `wtx_darwin_amd64` (Intel Mac) or `wtx_darwin_arm64` (Apple Silicon) from [Releases](https://github.com/t-shiratori/wtx/releases).

### Go install (local)

```bash
git clone https://github.com/t-shiratori/wtx.git
cd wtx
go install .
```

### Build from source

```bash
git clone https://github.com/t-shiratori/wtx.git
cd wtx
go build -o wtx .
```

## Usage

### Add a worktree

```bash
wtx add <branch>
```

**Options:**

| Option | Description |
|--------|-------------|
| `-b, --create-branch` | Create a new branch if it doesn't exist |
| `--from <branch>` | Specify the base branch |
| `--dry-run` | Show what would be done without executing |

**Examples:**

```bash
# Create a worktree from an existing branch
wtx add feature/new-feature

# Create a new branch and add a worktree
wtx add -b feature/new-feature

# Create a new branch based on main
wtx add -b feature/new-feature --from main

# Preview what would be done (dry-run)
wtx add -b feature/new-feature --dry-run
```

### List worktrees

```bash
wtx list
```

Displays a list of worktrees in the current repository.

### Remove a worktree

```bash
wtx remove [worktree ...]
```

After removing a worktree, empty directories are automatically cleaned up (including empty parent directories up to the worktrees root).

**Options:**

| Option | Description |
|--------|-------------|
| `-b, --branch` | Also delete the branch along with the worktree |
| `-f, --force` | Force deletion |
| `--dry-run` | Show what would be done without executing |

**Examples:**

```bash
# Interactively select and delete using TUI
wtx remove

# Delete a specific worktree
wtx remove feature/old-feature

# Also delete the branch
wtx remove -b feature/old-feature

# Delete multiple worktrees at once
wtx remove feature/a feature/b feature/c
```

### Initialize configuration

```bash
wtx init
```

Initializes the configuration file (`config.toml`).

**Options:**

| Option | Description |
|--------|-------------|
| `-f, --force` | Overwrite existing configuration file |
| `--global` | Create global configuration in `~/.config/wtx/` |

**Examples:**

```bash
# Create project-local configuration (.wtx/config.toml)
wtx init

# Overwrite existing configuration
wtx init --force

# Create global configuration (~/.config/wtx/config.toml)
wtx init --global
```

## Configuration

You can customize settings by creating `.wtx/config.toml` in your project root.

```toml
[worktree]
root_dir = ""                    # Worktree directory (defaults to .wtx/worktrees if empty)
default_base_branch = "main"     # Default base branch

[[copy]]                         # File copy settings (multiple entries allowed)
from = ".env.example"            # Source file (relative to repository root)
to = ".env"                      # Destination file (relative to worktree)

[hooks]                          # Hook settings
pre_create = ["echo pre"]        # Run before worktree creation
post_create = ["echo post"]      # Run after worktree creation
post_copy = ["echo copied"]      # Run after file copying
```

### Configuration options

#### `[worktree]`

| Option | Description | Default |
|--------|-------------|---------|
| `root_dir` | Directory where worktrees are created | `.wtx/worktrees` |
| `default_base_branch` | Base branch used when `--from` is omitted | `main` |

#### `[[copy]]`

Automatically copies files when creating a worktree. Multiple entries can be specified.

| Option | Description |
|--------|-------------|
| `from` | Source file (relative path from repository root) |
| `to` | Destination file (relative path from worktree) |

**Copying multiple files:**

```toml
[[copy]]
from = ".env.example"
to = ".env"

[[copy]]
from = "config/local.example.json"
to = "config/local.json"

[[copy]]
from = ".npmrc.example"
to = ".npmrc"
```

#### `[hooks]`

You can specify commands to run at each timing.

| Option | Execution timing |
|--------|------------------|
| `pre_create` | Before worktree creation |
| `post_create` | After worktree creation |
| `post_copy` | After file copying |

## License

MIT
