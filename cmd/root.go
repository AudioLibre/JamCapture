package cmd

import (
	"fmt"
	"os"

	"github.com/audiolibre/jamcapture/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	cfgFile   string
	pipeline  string
	profile   string
)

var rootCmd = &cobra.Command{
	Use:   "jamcapture [song-name]",
	Short: "Audio recording and mixing tool for jam sessions",
	Long: `JamCapture is a CLI tool for recording audio with support for
Bluetooth latency compensation and multi-track mixing.

It supports recording guitar input and system audio simultaneously,
then mixing them with customizable volume levels and delay compensation.

When a song name is provided, it acts as 'jamcapture run [song-name]'.`,
	Args: cobra.MaximumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for sources command
		if cmd.Name() == "sources" {
			return nil
		}

		if cfgFile == "" {
			return fmt.Errorf("config file required, use --config flag")
		}

		var err error
		cfg, err = config.LoadWithProfile(cfgFile, profile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate pipeline if provided
		if err := validatePipeline(); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// If a song name is provided, delegate to run command
		if len(args) == 1 {
			return runCmd.RunE(cmd, args)
		}
		// Otherwise show help
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/jamcapture.yaml)")
	rootCmd.PersistentFlags().StringVarP(&pipeline, "pipeline", "p", "", "pipeline steps: r=record, m=mix, p=play (e.g., 'rmp', 'mp', 'rm')")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "configuration profile to use (overrides active_config from file)")

	// Add flags for direct song execution
	rootCmd.Flags().Float64P("guitar-volume", "g", 0, "guitar volume (overrides config)")
	rootCmd.Flags().Float64P("backing-volume", "b", 0, "backing volume (overrides config)")
	rootCmd.Flags().IntP("delay", "d", -1, "backing track delay in ms (overrides config)")
	rootCmd.Flags().StringP("output", "o", "", "output directory (overrides config)")

	// Add subcommands
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(recordCmd)
	rootCmd.AddCommand(mixCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(sourcesCmd)
}