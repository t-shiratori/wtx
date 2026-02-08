package app

import "wtx/internal/config"

// Context holds the application-wide state stored in context
type Context struct {
	Config   *config.Config
	RepoRoot string
}

// contextKey is used to avoid key collisions in context
type contextKey struct{}

var Key = contextKey{}
