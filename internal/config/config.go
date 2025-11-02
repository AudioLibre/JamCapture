package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type RootConfig struct {
	ActiveConfig string                 `mapstructure:"active_config" yaml:"active_config"`
	Configs      map[string]*Config     `mapstructure:"configs" yaml:"configs"`
}

type Config struct {
	Audio  AudioConfig  `mapstructure:"audio" yaml:"audio"`
	Record RecordConfig `mapstructure:"record" yaml:"record"`
	Mix    MixConfig    `mapstructure:"mix" yaml:"mix"`
	Output OutputConfig `mapstructure:"output" yaml:"output"`
}

type AudioConfig struct {
	SampleRate int `mapstructure:"sample_rate" yaml:"sample_rate"`
	Channels   int `mapstructure:"channels" yaml:"channels"`
}

type RecordConfig struct {
	GuitarInput  string `mapstructure:"guitar_input" yaml:"guitar_input"`
	MonitorInput string `mapstructure:"monitor_input" yaml:"monitor_input"`
}

type MixConfig struct {
	GuitarVolume  float64 `mapstructure:"guitar_volume" yaml:"guitar_volume"`
	BackingVolume float64 `mapstructure:"backing_volume" yaml:"backing_volume"`
	DelayMs       int     `mapstructure:"delay_ms" yaml:"delay_ms"`
}

type OutputConfig struct {
	Directory string `mapstructure:"directory" yaml:"directory"`
	Format    string `mapstructure:"format" yaml:"format"`
}

var defaultConfig = Config{
	Audio: AudioConfig{
		SampleRate: 48000,
		Channels:   2,
	},
	Record: RecordConfig{
		GuitarInput:  "alsa_input.usb-Focusrite_Scarlett_2i2_USB_Y814JK8264026F-00.analog-stereo",
		MonitorInput: "", // Auto-detect if empty
	},
	Mix: MixConfig{
		GuitarVolume:  4.0,
		BackingVolume: 0.8,
		DelayMs:       0,
	},
	Output: OutputConfig{
		Directory: filepath.Join(os.Getenv("HOME"), "Audio", "JamCapture"),
		Format:    "flac",
	},
}

func Load(configFile string) (*Config, error) {
	return LoadWithProfile(configFile, "")
}

func LoadWithProfile(configFile, profile string) (*Config, error) {
	if configFile == "" {
		return nil, fmt.Errorf("no config file specified, use --config flag")
	}

	viper.SetConfigFile(configFile)

	// Set environment variable prefix
	viper.SetEnvPrefix("JAMCAPTURE")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configFile, err)
	}

	var rootConfig RootConfig
	if err := viper.Unmarshal(&rootConfig); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Determine which config to use
	configName := profile
	if configName == "" {
		configName = rootConfig.ActiveConfig
	}
	if configName == "" {
		configName = "default"
	}

	// Get the requested config
	selectedConfig, exists := rootConfig.Configs[configName]
	if !exists {
		return nil, fmt.Errorf("configuration profile '%s' not found", configName)
	}

	// Merge with default config if it exists and we're not already using default
	if configName != "default" {
		if defaultConfig, exists := rootConfig.Configs["default"]; exists {
			selectedConfig = mergeConfigs(defaultConfig, selectedConfig)
		}
	}

	// Expand tilde in output directory
	selectedConfig.Output.Directory = expandPath(selectedConfig.Output.Directory)

	return selectedConfig, nil
}


func (c *Config) Save() error {
	return viper.WriteConfig()
}

func mergeConfigs(base, override *Config) *Config {
	result := &Config{}

	// Start with base config
	if base != nil {
		*result = *base
	}

	// Override with non-zero values from override config
	if override != nil {
		// Audio config
		if override.Audio.SampleRate != 0 {
			result.Audio.SampleRate = override.Audio.SampleRate
		}
		if override.Audio.Channels != 0 {
			result.Audio.Channels = override.Audio.Channels
		}

		// Record config
		if override.Record.GuitarInput != "" {
			result.Record.GuitarInput = override.Record.GuitarInput
		}
		if override.Record.MonitorInput != "" {
			result.Record.MonitorInput = override.Record.MonitorInput
		}

		// Mix config - use different logic for DelayMs to allow 0 values
		if override.Mix.GuitarVolume != 0 {
			result.Mix.GuitarVolume = override.Mix.GuitarVolume
		}
		if override.Mix.BackingVolume != 0 {
			result.Mix.BackingVolume = override.Mix.BackingVolume
		}
		// For DelayMs, handle special case where -999 means force to 0
		if override.Mix.DelayMs == -999 {
			result.Mix.DelayMs = 0
		} else if override.Mix.DelayMs != 0 {
			result.Mix.DelayMs = override.Mix.DelayMs
		}

		// Output config
		if override.Output.Directory != "" {
			result.Output.Directory = override.Output.Directory
		}
		if override.Output.Format != "" {
			result.Output.Format = override.Output.Format
		}
	}

	return result
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}