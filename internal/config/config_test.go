package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"cli-t/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "cli-t-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set HOME to temp directory
	os.Setenv("HOME", tmpDir)

	tests := []struct {
		name      string
		setupFunc func(t *testing.T)
		envVars   map[string]string
		wantErr   bool
		checkFunc func(t *testing.T, cfg *config.Config)
	}{
		{
			name: "load with defaults (no config file)",
			setupFunc: func(t *testing.T) {
				// No config file
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "1", cfg.Version)
				assert.Equal(t, "default", cfg.Global.Theme)
				assert.Equal(t, "auto", cfg.Global.OutputFormat)
				assert.Equal(t, "auto", cfg.Global.Color)
				assert.True(t, cfg.UI.ShowHints)
				assert.True(t, cfg.UI.ShowIcons)
				assert.Equal(t, "bar", cfg.UI.ProgressStyle)
				assert.Equal(t, "rounded", cfg.UI.TableStyle)
				assert.Equal(t, "relative", cfg.UI.TimeFormat)
			},
		},
		{
			name: "load with config file",
			setupFunc: func(t *testing.T) {
				// Create config directory and file
				configDir := filepath.Join(tmpDir, ".config", "cli-t")
				err := os.MkdirAll(configDir, 0755)
				require.NoError(t, err)

				configFile := filepath.Join(configDir, "config.yaml")
				configContent := `
version: "1"
global:
  theme: monokai
  output_format: json
  color: always
ui:
  show_hints: false
  show_icons: false
  progress_style: dots
tools:
  wc:
    default_flags: ["-l", "-w"]
aliases:
  g: grep
  j: json
`
				err = os.WriteFile(configFile, []byte(configContent), 0644)
				require.NoError(t, err)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "1", cfg.Version)
				assert.Equal(t, "monokai", cfg.Global.Theme)
				assert.Equal(t, "json", cfg.Global.OutputFormat)
				assert.Equal(t, "always", cfg.Global.Color)
				assert.False(t, cfg.UI.ShowHints)
				assert.False(t, cfg.UI.ShowIcons)
				assert.Equal(t, "dots", cfg.UI.ProgressStyle)

				// Check tools config
				assert.NotNil(t, cfg.Tools["wc"])

				// Check aliases
				assert.Equal(t, "grep", cfg.Aliases["g"])
				assert.Equal(t, "json", cfg.Aliases["j"])
			},
		},
		{
			name: "environment variables override config",
			setupFunc: func(t *testing.T) {
				// Create a config file
				configDir := filepath.Join(tmpDir, ".config", "cli-t")
				err := os.MkdirAll(configDir, 0755)
				require.NoError(t, err)

				configFile := filepath.Join(configDir, "config.yaml")
				configContent := `
version: "1"
global:
  theme: monokai
`
				err = os.WriteFile(configFile, []byte(configContent), 0644)
				require.NoError(t, err)
			},
			envVars: map[string]string{
				"CLI_T_GLOBAL_THEME": "dracula",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *config.Config) {
				// Environment variable should override config file
				assert.Equal(t, "dracula", cfg.Global.Theme)
			},
		},
		{
			name: "invalid yaml in config file",
			setupFunc: func(t *testing.T) {
				// Create config directory and file with invalid YAML
				configDir := filepath.Join(tmpDir, ".config", "cli-t")
				err := os.MkdirAll(configDir, 0755)
				require.NoError(t, err)

				configFile := filepath.Join(configDir, "config.yaml")
				configContent := `
version: "1"
global:
  theme: monokai
  invalid yaml here
`
				err = os.WriteFile(configFile, []byte(configContent), 0644)
				require.NoError(t, err)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			for _, key := range []string{"CLI_T_GLOBAL_THEME", "CLI_T_GLOBAL_OUTPUT_FORMAT"} {
				os.Unsetenv(key)
			}

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Setup test
			if tt.setupFunc != nil {
				tt.setupFunc(t)
			}

			// Load config
			cfg, err := config.Load()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if tt.checkFunc != nil {
					tt.checkFunc(t, cfg)
				}
			}
		})
	}
}

func TestConfig_Get(t *testing.T) {
	// This tests the Get function after Load
	cfg, err := config.Load()
	require.NoError(t, err)

	got := config.Get()
	assert.Equal(t, cfg, got)
}

func TestConfig_ConfigDir(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name     string
		homeDir  string
		expected string
	}{
		{
			name:     "standard home directory",
			homeDir:  "/home/user",
			expected: "/home/user/.config/cli-t",
		},
		{
			name:     "root user",
			homeDir:  "/root",
			expected: "/root/.config/cli-t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.homeDir)
			dir := config.ConfigDir()
			assert.Equal(t, tt.expected, dir)
		})
	}
}

func TestConfig_Save(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "cli-t-config-save-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)

	// Create config directory
	configDir := filepath.Join(tmpDir, ".config", "cli-t")
	err = os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// Create initial config file
	configFile := filepath.Join(configDir, "config.yaml")
	initialContent := `
version: "1"
global:
  theme: default
`
	err = os.WriteFile(configFile, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := config.Load()
	require.NoError(t, err)

	// Modify config
	cfg.Global.Theme = "monokai"
	cfg.Global.OutputFormat = "json"
	cfg.UI.ShowHints = false

	// Save config
	err = config.Save(cfg)
	require.NoError(t, err)

	// Load again to verify
	newCfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "monokai", newCfg.Global.Theme)
	assert.Equal(t, "json", newCfg.Global.OutputFormat)
	assert.False(t, newCfg.UI.ShowHints)
}

func TestToolConfig(t *testing.T) {
	toolConfig := config.ToolConfig{
		"port":        6379,
		"persistence": true,
		"data_dir":    "~/.cli-t/redis",
		"flags":       []string{"-l", "-w"},
	}

	// Test accessing various types
	assert.Equal(t, 6379, toolConfig["port"])
	assert.Equal(t, true, toolConfig["persistence"])
	assert.Equal(t, "~/.cli-t/redis", toolConfig["data_dir"])

	// Test type assertion
	flags, ok := toolConfig["flags"].([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"-l", "-w"}, flags)
}
