package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the .codesentry.yaml configuration
type Config struct {
	RulesDir          string
	Include           []string
	Exclude           []string
	SeverityThreshold string
	NoColor           bool
}

// BindFlags binds command-line flags to viper (cobra pattern)
func BindFlags(v *viper.Viper, cmd *Config) {
	// Flags will override config file values
}

// Load reads and parses a .codesentry.yaml file using viper
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Search for config file
		v.SetConfigName(".codesentry")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./.config")
		v.AddConfigPath("$HOME/.codesentry")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return FromViper(v), nil
}

// LoadFromDefault loads config from default search paths
func LoadFromDefault() (*Config, error) {
	return Load("")
}

// FromViper creates a Config from viper instance
func FromViper(v *viper.Viper) *Config {
	return &Config{
		RulesDir:          v.GetString("rules-dir"),
		Include:           v.GetStringSlice("include"),
		Exclude:           v.GetStringSlice("exclude"),
		SeverityThreshold: v.GetString("severity-threshold"),
		NoColor:           v.GetBool("no-color"),
	}
}

// SetDefaults sets default values on viper
func SetDefaults(v *viper.Viper) {
	v.SetDefault("rules-dir", "rules")
	v.SetDefault("include", []string{})
	v.SetDefault("exclude", []string{})
	v.SetDefault("severity-threshold", "low")
	v.SetDefault("no-color", false)
}

// MergeConfig merges another config into viper (for cascading config)
func MergeConfig(v *viper.Viper, cfg *Config) {
	if cfg.RulesDir != "" {
		v.Set("rules-dir", cfg.RulesDir)
	}
	if len(cfg.Include) > 0 {
		v.Set("include", cfg.Include)
	}
	if len(cfg.Exclude) > 0 {
		v.Set("exclude", cfg.Exclude)
	}
	if cfg.SeverityThreshold != "" {
		v.Set("severity-threshold", cfg.SeverityThreshold)
	}
	v.Set("no-color", cfg.NoColor)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	validSeverities := []string{"critical", "high", "medium", "low", "info"}
	severity := strings.ToLower(c.SeverityThreshold)

	for _, valid := range validSeverities {
		if severity == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid severity-threshold: %s (allowed: %s)", c.SeverityThreshold, strings.Join(validSeverities, ", "))
}

// ToMap converts config to map for debugging
func (c *Config) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"rules-dir":          c.RulesDir,
		"include":            c.Include,
		"exclude":            c.Exclude,
		"severity-threshold": c.SeverityThreshold,
		"no-color":           c.NoColor,
	}
}

// Exists checks if a config file exists at the given path
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
