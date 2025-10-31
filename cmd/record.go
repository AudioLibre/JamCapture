package cmd

import (
	"fmt"
	"jamcapture/internal/record"

	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record [song-name]",
	Short: "Record guitar input and system audio",
	Long: `Record audio from guitar input and system audio monitor simultaneously.
The recording will be saved as an MKV file with separate tracks for guitar and backing audio.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		songName := args[0]

		fmt.Printf("Recording song: %s\n", songName)
		fmt.Printf("Guitar input: %s\n", cfg.Record.GuitarInput)
		fmt.Printf("Sample rate: %d Hz\n", cfg.Audio.SampleRate)
		fmt.Printf("Output directory: %s\n", cfg.Output.Directory)

		recorder := record.New(cfg)
		err := recorder.Record(songName)
		if err != nil {
			return fmt.Errorf("recording failed: %w", err)
		}

		fmt.Println("Recording completed successfully")

		// Execute pipeline if specified
		return executePipeline(songName, 'r')
	},
}

func init() {
	recordCmd.Flags().StringP("output", "o", "", "output directory (overrides config)")
}