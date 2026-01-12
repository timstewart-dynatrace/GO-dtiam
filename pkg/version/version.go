// Package version provides version information for the CLI.
package version

// These variables are set at build time using ldflags.
var (
	// Version is the semantic version of the CLI.
	Version = "dev"
	// Commit is the git commit hash.
	Commit = "unknown"
	// Date is the build date.
	Date = "unknown"
)
