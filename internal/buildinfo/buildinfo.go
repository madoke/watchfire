// Package buildinfo holds version information injected at build time via ldflags.
package buildinfo

// Build-time variables injected via ldflags.
var (
	Version    = "dev"
	Codename   = "unknown"
	CommitHash = "unknown"
	BuildDate  = "unknown"
	PostHogKey = ""
)
