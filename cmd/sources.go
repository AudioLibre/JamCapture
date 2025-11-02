package cmd

import (
	"fmt"
	"strings"
	"github.com/audiolibre/jamcapture/internal/audio"

	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "List available audio sources",
	Long:  `List all available PulseAudio and PipeWire sources that can be used for recording.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pulse := audio.NewPulseAudio()

		sources, err := pulse.ListSources()
		if err != nil {
			return fmt.Errorf("failed to list sources: %w", err)
		}

		fmt.Println("Available audio sources:")
		carlaFound := false
		for i, source := range sources {
			fmt.Printf("%d. %s\n", i+1, source)

			// Highlight Carla sources
			if strings.Contains(source, "Carla") {
				fmt.Printf("   → Carla audio output (can be used as guitar source)\n")
				carlaFound = true
			}
		}

		// Show current default monitor
		monitor, err := pulse.GetDefaultSinkMonitor()
		if err == nil {
			fmt.Printf("\nDefault system monitor: %s\n", monitor)
		}

		// Inform about Carla if not found
		if !carlaFound {
			fmt.Printf("\nNote: No Carla sources detected. If you're using Carla, make sure it's running with 'pw-jack carla'.\n")
		}

		return nil
	},
}