package cmd

import (
	"fmt"
	"github.com/audiolibre/jamcapture/internal/play"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play [song-name]",
	Short: "Play the mixed audio file",
	Long: `Play the mixed FLAC file using the system's default audio player.
Will attempt to use VLC if available, otherwise falls back to system default.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		songName := args[0]

		fmt.Printf("Playing song: %s\n", songName)

		player := play.New(cfg)
		err := player.Play(songName)
		if err != nil {
			return fmt.Errorf("playback failed: %w", err)
		}

		return nil
	},
}