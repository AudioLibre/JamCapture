package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

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

func Load() (*Config, error) {
	viper.SetConfigName("jamcapture")
	viper.SetConfigType("yaml")

	// Add config search paths
	homeDir, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(homeDir, ".config"))
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix("JAMCAPTURE")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create one
			return createDefaultConfig()
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Expand tilde in output directory
	config.Output.Directory = expandPath(config.Output.Directory)

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("audio.sample_rate", defaultConfig.Audio.SampleRate)
	viper.SetDefault("audio.channels", defaultConfig.Audio.Channels)
	viper.SetDefault("record.guitar_input", defaultConfig.Record.GuitarInput)
	viper.SetDefault("record.monitor_input", defaultConfig.Record.MonitorInput)
	viper.SetDefault("mix.guitar_volume", defaultConfig.Mix.GuitarVolume)
	viper.SetDefault("mix.backing_volume", defaultConfig.Mix.BackingVolume)
	viper.SetDefault("mix.delay_ms", defaultConfig.Mix.DelayMs)
	viper.SetDefault("output.directory", defaultConfig.Output.Directory)
	viper.SetDefault("output.format", defaultConfig.Output.Format)
}

func createDefaultConfig() (*Config, error) {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "jamcapture.yaml")

	// Write default config
	viper.SetConfigFile(configPath)
	if err := viper.WriteConfig(); err != nil {
		return nil, fmt.Errorf("error writing default config: %w", err)
	}

	fmt.Printf("Created default config at: %s\n", configPath)
	return &defaultConfig, nil
}

func (c *Config) Save() error {
	return viper.WriteConfig()
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}