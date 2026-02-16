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
- Version information display

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

### Check version

```bash
wtx --version
```

Displays the current version of wtx.

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

$ wtx remove            
Select worktrees to remove

  [ ] /Users/user/repo1 (main)
  [ ] /Users/user/repo1/worktrees/feature/task-1 (feature/task-1)
  [ ] /Users/user/repo1/worktrees/feature/task-2 (feature/task-2)
  [x] /Users/user/repo1/worktrees/feature/task-3 (feature/task-3)
> [x] /Users/user/repo1/worktrees/feature/task-4 (feature/task-4)

↑/↓ move • space select • enter confirm • q cancel

```

````bash
# Delete a specific worktree
wtx remove feature/old-feature

# Also delete the branch
wtx remove -b feature/old-feature

# Delete multiple worktrees at once
wtx remove feature/a feature/b feature/c
````

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

[copy]
patterns = ["*.env", "config/*.yaml"]  # Glob patterns for file copying

# For renaming files (optional)
[[copy.files]]
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

#### `[copy]`

Automatically copies files when creating a worktree.

##### Glob Patterns: `patterns`

Use glob patterns to specify multiple files at once:

```toml
[copy]
patterns = [
    "*.env",
    "config/*.yaml",
    ".npmrc.example",
    "secrets/*.json"
]
```

Files matching the patterns are copied from the repository root to the same relative path in the worktree.

**Supported patterns:**
- `*` - Matches any string (excluding directory separators)
- `?` - Matches any single character
- `[abc]` - Matches any character in the set

**Note:** Recursive glob patterns (`**`) are not supported in the initial version.

##### Explicit Files: `[[copy.files]]` (For renaming)

Use this format when you need to rename files or copy to a different path:

```toml
[[copy.files]]
from = ".env.example"
to = ".env"

[[copy.files]]
from = "config/dev.yaml"
to = "config/local.yaml"
```

| Option | Description |
|--------|-------------|
| `from` | Source file (relative path from repository root) |
| `to` | Destination file (relative path from worktree) |

##### Using Both Formats

You can use both patterns and explicit files together:

```toml
[copy]
patterns = ["*.env", "config/*.yaml"]

# For renaming specific files
[[copy.files]]
from = ".env.example"
to = ".env"
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
