package version_test

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"

	"cli-t/pkg/version"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion_String(t *testing.T) {
	// Save original values
	origVersion := version.Version
	origCommit := version.Commit
	origBuildTime := version.BuildTime

	tests := []struct {
		name      string
		version   string
		commit    string
		buildTime string
		want      string
	}{
		{
			name:      "dev version",
			version:   "dev",
			commit:    "unknown",
			buildTime: "unknown",
			want:      "CLI-T dev (commit: unknown, built: unknown)",
		},
		{
			name:      "release version",
			version:   "1.2.3",
			commit:    "abc123def456",
			buildTime: "2024-01-15T10:00:00Z",
			want:      "CLI-T 1.2.3 (commit: abc123d, built: 2024-01-15T10:00:00Z)",
		},
		{
			name:      "pre-release version",
			version:   "1.2.3-beta.1",
			commit:    "def456abc789",
			buildTime: "2024-01-15T10:00:00Z",
			want:      "CLI-T 1.2.3-beta.1 (commit: def456a, built: 2024-01-15T10:00:00Z)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test values
			version.Version = tt.version
			version.Commit = tt.commit
			version.BuildTime = tt.buildTime

			got := version.String()
			assert.Equal(t, tt.want, got)
		})
	}

	// Restore original values
	version.Version = origVersion
	version.Commit = origCommit
	version.BuildTime = origBuildTime
}

func TestVersion_ShortCommit(t *testing.T) {
	// Save original
	origCommit := version.Commit

	tests := []struct {
		name   string
		commit string
		want   string
	}{
		{
			name:   "long commit hash",
			commit: "abc123def456789012345",
			want:   "abc123d",
		},
		{
			name:   "short commit hash",
			commit: "abc123",
			want:   "abc123",
		},
		{
			name:   "empty commit",
			commit: "",
			want:   "",
		},
		{
			name:   "exactly 7 chars",
			commit: "1234567",
			want:   "1234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version.Commit = tt.commit
			got := version.ShortCommit()
			assert.Equal(t, tt.want, got)
		})
	}

	// Restore
	version.Commit = origCommit
}

func TestVersion_GetInfo(t *testing.T) {
	// Save originals
	origVersion := version.Version
	origCommit := version.Commit
	origBuildTime := version.BuildTime
	origOS := version.OS
	origArch := version.Arch

	// Set test values
	version.Version = "1.0.0"
	version.Commit = "testcommit"
	version.BuildTime = "2024-01-15"
	version.OS = "linux"
	version.Arch = "amd64"

	info := version.GetInfo()

	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, "testcommit", info.Commit)
	assert.Equal(t, "2024-01-15", info.BuildTime)
	assert.Equal(t, "linux", info.OS)
	assert.Equal(t, "amd64", info.Arch)
	assert.Equal(t, runtime.Version(), info.GoVersion)

	// Test JSON serialization
	data, err := json.Marshal(info)
	require.NoError(t, err)

	var decoded version.Info
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, info, decoded)

	// Restore
	version.Version = origVersion
	version.Commit = origCommit
	version.BuildTime = origBuildTime
	version.OS = origOS
	version.Arch = origArch
}

func TestVersion_DetailedString(t *testing.T) {
	// Save originals
	origVersion := version.Version
	origCommit := version.Commit
	origBuildTime := version.BuildTime
	origOS := version.OS
	origArch := version.Arch

	tests := []struct {
		name      string
		version   string
		commit    string
		buildTime string
		os        string
		arch      string
		checkFunc func(t *testing.T, output string)
	}{
		{
			name:      "all fields present",
			version:   "1.2.3",
			commit:    "abc123def",
			buildTime: "2024-01-15T10:00:00Z",
			os:        "darwin",
			arch:      "arm64",
			checkFunc: func(t *testing.T, output string) {
				assert.Contains(t, output, "CLI-T version 1.2.3")
				assert.Contains(t, output, "Git commit: abc123def")
				assert.Contains(t, output, "Built: 2024-01-15T10:00:00Z")
				assert.Contains(t, output, "OS/Arch: darwin/arm64")
				assert.Contains(t, output, "Go version: "+runtime.Version())
			},
		},
		{
			name:      "unknown fields",
			version:   "dev",
			commit:    "unknown",
			buildTime: "unknown",
			os:        runtime.GOOS,
			arch:      runtime.GOARCH,
			checkFunc: func(t *testing.T, output string) {
				assert.Contains(t, output, "CLI-T version dev")
				assert.NotContains(t, output, "Git commit") // Should skip unknown
				assert.NotContains(t, output, "Built:")     // Should skip unknown
				assert.Contains(t, output, "OS/Arch: "+runtime.GOOS+"/"+runtime.GOARCH)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version.Version = tt.version
			version.Commit = tt.commit
			version.BuildTime = tt.buildTime
			version.OS = tt.os
			version.Arch = tt.arch

			output := version.DetailedString()
			tt.checkFunc(t, output)

			// Check it has multiple lines
			lines := strings.Split(output, "\n")
			assert.Greater(t, len(lines), 1)
		})
	}

	// Restore
	version.Version = origVersion
	version.Commit = origCommit
	version.BuildTime = origBuildTime
	version.OS = origOS
	version.Arch = origArch
}

func TestVersion_IsDevBuild(t *testing.T) {
	// Save original
	origVersion := version.Version

	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "dev build",
			version: "dev",
			want:    true,
		},
		{
			name:    "release build",
			version: "1.0.0",
			want:    false,
		},
		{
			name:    "pre-release build",
			version: "1.0.0-beta",
			want:    false,
		},
		{
			name:    "dev prefix but not dev",
			version: "develop",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version.Version = tt.version
			got := version.IsDevBuild()
			assert.Equal(t, tt.want, got)
		})
	}

	// Restore
	version.Version = origVersion
}

func TestVersion_IsReleaseBuild(t *testing.T) {
	// Save original
	origVersion := version.Version

	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "release build",
			version: "1.0.0",
			want:    true,
		},
		{
			name:    "major release",
			version: "2.0.0",
			want:    true,
		},
		{
			name:    "patch release",
			version: "1.0.1",
			want:    true,
		},
		{
			name:    "dev build",
			version: "dev",
			want:    false,
		},
		{
			name:    "pre-release alpha",
			version: "1.0.0-alpha",
			want:    false,
		},
		{
			name:    "pre-release beta",
			version: "1.0.0-beta.1",
			want:    false,
		},
		{
			name:    "release candidate",
			version: "1.0.0-rc.1",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version.Version = tt.version
			got := version.IsReleaseBuild()
			assert.Equal(t, tt.want, got)
		})
	}

	// Restore
	version.Version = origVersion
}

func TestVersion_CompareVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{
			name: "equal versions",
			v1:   "1.0.0",
			v2:   "1.0.0",
			want: 0,
		},
		{
			name: "v1 less than v2",
			v1:   "1.0.0",
			v2:   "2.0.0",
			want: -1,
		},
		{
			name: "v1 greater than v2",
			v1:   "2.0.0",
			v2:   "1.0.0",
			want: 1,
		},
		{
			name: "dev versions",
			v1:   "dev",
			v2:   "dev",
			want: 0,
		},
		// Note: This is a simple string comparison
		// A real implementation would need semantic version parsing
		{
			name: "string comparison quirk",
			v1:   "10.0.0",
			v2:   "2.0.0",
			want: -1, // String comparison: "1" < "2"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := version.CompareVersions(tt.v1, tt.v2)
			if tt.want == 0 {
				assert.Equal(t, 0, got)
			} else if tt.want < 0 {
				assert.Less(t, got, 0)
			} else {
				assert.Greater(t, got, 0)
			}
		})
	}
}

// Test runtime info accuracy
func TestVersion_RuntimeInfo(t *testing.T) {
	// When not explicitly set, should use runtime values
	assert.Equal(t, runtime.GOOS, version.OS)
	assert.Equal(t, runtime.GOARCH, version.Arch)

	info := version.GetInfo()
	assert.Equal(t, runtime.Version(), info.GoVersion)
}
