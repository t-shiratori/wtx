# wtx

Git worktreeを簡単に管理するためのCLIツール。

## 概要

wtxは、Gitのworktree機能を効率的に管理するための小型CLIツールです。複数のブランチで並行作業を行う際に、worktreeの作成・削除・一覧表示を簡単に行えます。

主な機能:
- worktreeの作成・削除・一覧表示
- TUI（ターミナルUI）によるworktree選択
- フック機能（pre_create, post_create, post_copy）
- ファイルの自動コピー機能

## インストール

### Homebrew

```bash
brew install t-shiratori/tap/wtx
```

### GitHubリリースから直接ダウンロード

[Releases](https://github.com/t-shiratori/wtx/releases)から`wtx_darwin_amd64`（Intel Mac）または`wtx_darwin_arm64`（Apple Silicon）をダウンロードしてください。

### ソースからビルド

```bash
git clone https://github.com/t-shiratori/wtx.git
cd wtx
go build -o wtx .
```

## 使い方

### worktreeの追加

```bash
wtx add <branch>
```

**オプション:**

| オプション | 説明 |
|-----------|------|
| `-b, --create-branch` | ブランチが存在しない場合は新規作成 |
| `--from <branch>` | ベースブランチまたはコミットを指定 |
| `--dry-run` | 実行内容を表示するだけで実行しない |

**例:**

```bash
# 既存ブランチからworktreeを作成
wtx add feature/new-feature

# 新しいブランチを作成してworktreeを追加
wtx add -b feature/new-feature

# mainブランチをベースに新しいブランチを作成
wtx add -b feature/new-feature --from main

# 実行内容を確認（dry-run）
wtx add -b feature/new-feature --dry-run
```

### worktreeの一覧表示

```bash
wtx list
```

現在のリポジトリに存在するworktreeの一覧を表示します。

### worktreeの削除

```bash
wtx remove [worktree ...]
```

**オプション:**

| オプション | 説明 |
|-----------|------|
| `-b, --branch` | worktreeと一緒にブランチも削除 |
| `-f, --force` | 強制削除 |
| `--dry-run` | 実行内容を表示するだけで実行しない |

**例:**

```bash
# TUIで対話的に選択して削除
wtx remove

# 指定したworktreeを削除
wtx remove feature/old-feature

# ブランチも一緒に削除
wtx remove -b feature/old-feature

# 複数のworktreeを一度に削除
wtx remove feature/a feature/b feature/c
```

## 設定

プロジェクトルートに`.wtx/config.toml`を作成して設定をカスタマイズできます。

```toml
[worktree]
root_dir = ""                    # worktreeディレクトリ（空の場合は.wtx/worktrees）
default_base_branch = "main"     # デフォルトのベースブランチ

[[copy]]                         # ファイルコピー設定（複数指定可）
from = ".env.example"            # コピー元（リポジトリルート相対）
to = ".env"                      # コピー先（worktree相対）

[hooks]                          # フック設定
pre_create = ["echo pre"]        # worktree作成前に実行
post_create = ["echo post"]      # worktree作成後に実行
post_copy = ["echo copied"]      # ファイルコピー後に実行
```

### 設定項目

#### `[worktree]`

| 項目 | 説明 | デフォルト |
|------|------|----------|
| `root_dir` | worktreeを作成するディレクトリ | `.wtx/worktrees` |
| `default_base_branch` | `--from`を省略した場合のベースブランチ | `main` |

#### `[[copy]]`

worktree作成時に自動でファイルをコピーします。複数指定可能です。

| 項目 | 説明 |
|------|------|
| `from` | コピー元ファイル（リポジトリルートからの相対パス） |
| `to` | コピー先ファイル（worktreeからの相対パス） |

**複数ファイルをコピーする場合:**

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

各タイミングで実行するコマンドを指定できます。

| 項目 | 実行タイミング |
|------|---------------|
| `pre_create` | worktree作成前 |
| `post_create` | worktree作成後 |
| `post_copy` | ファイルコピー後 |

## ライセンス

MIT
