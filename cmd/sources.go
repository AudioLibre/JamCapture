package cmd

import (
	"fmt"
	"github.com/audiolibre/jamcapture/internal/audio"

	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "List available audio sources",
	Long:  `List all available PulseAudio sources that can be used for recording.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pulse := audio.NewPulseAudio()

		sources, err := pulse.ListSources()
		if err != nil {
			return fmt.Errorf("failed to list sources: %w", err)
		}

		fmt.Println("Available audio sources:")
		for i, source := range sources {
			fmt.Printf("%d. %s\n", i+1, source)
		}

		// Show current default monitor
		monitor, err := pulse.GetDefaultSinkMonitor()
		if err == nil {
			fmt.Printf("\nDefault system monitor: %s\n", monitor)
		}

		return nil
	},
}