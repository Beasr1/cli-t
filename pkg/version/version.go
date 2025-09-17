package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Version is the semantic version of CLI-T
	Version = "dev"

	// Commit is the git commit SHA
	Commit = "unknown"

	// BuildTime is when the binary was built
	BuildTime = "unknown"

	// These are set by goreleaser or build scripts
	OS   = runtime.GOOS
	Arch = runtime.GOARCH
)

// String returns the version information as a string
func String() string {
	return fmt.Sprintf("CLI-T %s (commit: %s, built: %s)",
		Version,
		ShortCommit(),
		BuildTime)
}

// ShortCommit returns the first 7 characters of the commit hash
func ShortCommit() string {
	if len(Commit) > 7 {
		return Commit[:7]
	}
	return Commit
}

// Info returns structured version information
type Info struct {
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit" yaml:"commit"`
	BuildTime string `json:"buildTime" yaml:"build_time"`
	OS        string `json:"os" yaml:"os"`
	Arch      string `json:"arch" yaml:"arch"`
	GoVersion string `json:"goVersion" yaml:"go_version"`
}

// GetInfo returns version info struct
func GetInfo() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		OS:        OS,
		Arch:      Arch,
		GoVersion: runtime.Version(),
	}
}

// DetailedString returns detailed version information
func DetailedString() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("CLI-T version %s", Version))

	if Commit != "unknown" {
		parts = append(parts, fmt.Sprintf("Git commit: %s", Commit))
	}

	if BuildTime != "unknown" {
		parts = append(parts, fmt.Sprintf("Built: %s", BuildTime))
	}

	parts = append(parts, fmt.Sprintf("OS/Arch: %s/%s", OS, Arch))
	parts = append(parts, fmt.Sprintf("Go version: %s", runtime.Version()))

	return strings.Join(parts, "\n")
}

// IsDevBuild returns true if this is a development build
func IsDevBuild() bool {
	return Version == "dev"
}

// IsReleaseBuild returns true if this is a release build
func IsReleaseBuild() bool {
	return !IsDevBuild() && !strings.Contains(Version, "-")
}

// CompareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareVersions(v1, v2 string) int {
	// Simple implementation - you might want to use a library like
	// github.com/Masterminds/semver for production
	if v1 == v2 {
		return 0
	}
	// Add proper semantic version comparison logic here
	return strings.Compare(v1, v2)
}
