package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Version string                `mapstructure:"version"`
	Global  GlobalConfig          `mapstructure:"global"`
	UI      UIConfig              `mapstructure:"ui"`
	Tools   map[string]ToolConfig `mapstructure:"tools"`
	Aliases map[string]string     `mapstructure:"aliases"`
	Hooks   HooksConfig           `mapstructure:"hooks"`
}

// GlobalConfig holds global settings
type GlobalConfig struct {
	Theme        string `mapstructure:"theme"`
	OutputFormat string `mapstructure:"output_format"`
	Color        string `mapstructure:"color"`
	Editor       string `mapstructure:"editor"`
	Pager        string `mapstructure:"pager"`
	Shell        string `mapstructure:"shell"`
}

// UIConfig holds UI customization settings
type UIConfig struct {
	ShowHints     bool   `mapstructure:"show_hints"`
	ShowIcons     bool   `mapstructure:"show_icons"`
	ProgressStyle string `mapstructure:"progress_style"`
	TableStyle    string `mapstructure:"table_style"`
	TimeFormat    string `mapstructure:"time_format"`
}

// ToolConfig holds tool-specific configuration
type ToolConfig map[string]interface{}

// HooksConfig holds hook scripts
type HooksConfig struct {
	PreExecute  string `mapstructure:"pre_execute"`
	PostExecute string `mapstructure:"post_execute"`
	OnError     string `mapstructure:"on_error"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

// Load loads the configuration
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		instance, err = load()
	})
	return instance, err
}

func load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Config search paths
	viper.AddConfigPath("$HOME/.config/cli-t")
	viper.AddConfigPath("$HOME/.cli-t")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("CLI_T")
	viper.AutomaticEnv()

	// Defaults
	setDefaults()

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		// Config not found, use defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("version", "1")
	viper.SetDefault("global.theme", "default")
	viper.SetDefault("global.output_format", "auto")
	viper.SetDefault("global.color", "auto")
	viper.SetDefault("global.editor", os.Getenv("EDITOR"))
	viper.SetDefault("global.pager", "less")
	viper.SetDefault("global.shell", os.Getenv("SHELL"))

	viper.SetDefault("ui.show_hints", true)
	viper.SetDefault("ui.show_icons", true)
	viper.SetDefault("ui.progress_style", "bar")
	viper.SetDefault("ui.table_style", "rounded")
	viper.SetDefault("ui.time_format", "relative")
}

// Get returns the current configuration
func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

// ConfigDir returns the configuration directory
func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "cli-t")
}

// Save saves the current configuration
func Save(config *Config) error {
	mu.Lock()
	defer mu.Unlock()

	instance = config
	return viper.WriteConfig()
}
