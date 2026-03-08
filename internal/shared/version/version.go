package version

// Version is the version of the wtx CLI.
// This can be overridden at build time using -ldflags:
// go build -ldflags "-X wtx/internal/shared/version.Version=v1.0.0"
var Version = "v0.0.0"
