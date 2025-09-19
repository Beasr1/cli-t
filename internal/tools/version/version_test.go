package versioncmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"cli-t/internal/command"
	versioncmd "cli-t/internal/tools/version"
	"cli-t/pkg/version"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion_Metadata(t *testing.T) {
	cmd := versioncmd.New()

	assert.Equal(t, "version", cmd.Name())
	assert.Equal(t, "[--short|-s] [--json]", cmd.Usage())
	assert.Equal(t, "Show version information", cmd.Description())
}

func TestVersion_ValidateArgs(t *testing.T) {
	cmd := versioncmd.New()

	// Version command accepts any arguments
	assert.NoError(t, cmd.ValidateArgs([]string{}))
	assert.NoError(t, cmd.ValidateArgs([]string{"extra", "args"}))
}

func TestVersion_DefineFlags(t *testing.T) {
	cmd := versioncmd.New()

	// Check if command implements FlagDefiner
	flagDefiner, ok := cmd.(command.FlagDefiner)
	require.True(t, ok, "version command should implement FlagDefiner")

	flags := flagDefiner.DefineFlags()

	// Should have 2 flags
	assert.Len(t, flags, 2)

	// Check each flag
	flagMap := make(map[string]command.Flag)
	for _, f := range flags {
		flagMap[f.Name] = f
	}

	// Check short flag
	shortFlag, ok := flagMap["short"]
	assert.True(t, ok)
	assert.Equal(t, "s", shortFlag.Shorthand)
	assert.Equal(t, "bool", shortFlag.Type)
	assert.Equal(t, false, shortFlag.Default)

	// Check json flag
	jsonFlag, ok := flagMap["json"]
	assert.True(t, ok)
	assert.Equal(t, "", jsonFlag.Shorthand)
	assert.Equal(t, "bool", jsonFlag.Type)
	assert.Equal(t, false, jsonFlag.Default)
}

func TestVersion_Execute(t *testing.T) {
	// Save original values
	origVersion := version.Version
	origCommit := version.Commit
	origBuildTime := version.BuildTime

	// Set test values
	version.Version = "1.2.3"
	version.Commit = "abc123def456"
	version.BuildTime = "2024-01-15T10:00:00Z"

	defer func() {
		// Restore original values
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildTime = origBuildTime
	}()

	tests := []struct {
		name      string
		flags     map[string]interface{}
		checkFunc func(t *testing.T, output string)
	}{
		{
			name:  "default output",
			flags: map[string]interface{}{},
			checkFunc: func(t *testing.T, output string) {
				// Should contain detailed version string
				assert.Contains(t, output, "CLI-T version 1.2.3")
				assert.Contains(t, output, "Git commit: abc123def456")
				assert.Contains(t, output, "Built: 2024-01-15T10:00:00Z")
				assert.Contains(t, output, "OS/Arch:")
				assert.Contains(t, output, "Go version:")
			},
		},
		{
			name: "short output",
			flags: map[string]interface{}{
				"short": true,
			},
			checkFunc: func(t *testing.T, output string) {
				// Should only contain version number
				assert.Equal(t, "1.2.3\n", output)
			},
		},
		{
			name: "json output",
			flags: map[string]interface{}{
				"json": true,
			},
			checkFunc: func(t *testing.T, output string) {
				// Should be valid JSON
				var info version.Info
				err := json.Unmarshal([]byte(output), &info)
				assert.NoError(t, err)

				assert.Equal(t, "1.2.3", info.Version)
				assert.Equal(t, "abc123def456", info.Commit)
				assert.Equal(t, "2024-01-15T10:00:00Z", info.BuildTime)
				assert.NotEmpty(t, info.OS)
				assert.NotEmpty(t, info.Arch)
				assert.NotEmpty(t, info.GoVersion)
			},
		},
		{
			name: "both short and json (json takes precedence)",
			flags: map[string]interface{}{
				"short": true,
				"json":  true,
			},
			checkFunc: func(t *testing.T, output string) {
				// JSON should take precedence
				var info version.Info
				err := json.Unmarshal([]byte(output), &info)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := versioncmd.New()
			args := helpers.TestArgs()
			args.Flags = tt.flags

			err := cmd.Execute(context.Background(), args)
			assert.NoError(t, err)

			stdout := args.Stdout.(*bytes.Buffer).String()
			tt.checkFunc(t, stdout)
		})
	}
}

func TestVersion_DevBuild(t *testing.T) {
	// Save original values
	origVersion := version.Version
	origCommit := version.Commit
	origBuildTime := version.BuildTime

	// Set dev values
	version.Version = "dev"
	version.Commit = "unknown"
	version.BuildTime = "unknown"

	defer func() {
		// Restore original values
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildTime = origBuildTime
	}()

	cmd := versioncmd.New()
	args := helpers.TestArgs()

	err := cmd.Execute(context.Background(), args)
	assert.NoError(t, err)

	stdout := args.Stdout.(*bytes.Buffer).String()
	assert.Contains(t, stdout, "CLI-T version dev")
	// Unknown values might not be shown in DetailedString()
}

// Benchmark
func BenchmarkVersion_Default(b *testing.B) {
	cmd := versioncmd.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs()
		cmd.Execute(context.Background(), args)
	}
}

func BenchmarkVersion_JSON(b *testing.B) {
	cmd := versioncmd.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs()
		args.Flags = map[string]interface{}{"json": true}
		cmd.Execute(context.Background(), args)
	}
}
