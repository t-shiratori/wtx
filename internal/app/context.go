package app

import "wtx/internal/config"

// context に載せるアプリ全体の状態
type Context struct {
	Config   *config.Config
	RepoRoot string
}

// context key 用（衝突防止）
type contextKey struct{}

var Key = contextKey{}
