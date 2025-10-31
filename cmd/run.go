package cmd

import (
	"fmt"
	"jamcapture/internal/mix"
	"jamcapture/internal/play"
	"jamcapture/internal/record"
	"strings"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [song-name]",
	Short: "Execute pipeline steps on a song",
	Long: `Execute the specified pipeline steps on a song. Use -p to specify which steps to run.
If no pipeline is specified, will try to infer the action based on existing files.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		songName := args[0]

		if pipeline == "" {
			return fmt.Errorf("no pipeline specified, use -p flag (e.g., -p rmp)")
		}

		steps := []rune(strings.ToLower(pipeline))

		for i, step := range steps {
			fmt.Printf("Pipeline: executing step %d/%d: '%c'...\n", i+1, len(steps), step)

			switch step {
			case 'r':
				fmt.Printf("Recording song: %s\n", songName)
				fmt.Printf("Guitar input: %s\n", cfg.Record.GuitarInput)
				fmt.Printf("Sample rate: %d Hz\n", cfg.Audio.SampleRate)
				fmt.Printf("Output directory: %s\n", cfg.Output.Directory)

				recorder := record.New(cfg)
				if err := recorder.Record(songName); err != nil {
					return fmt.Errorf("pipeline record failed: %w", err)
				}
				fmt.Println("Pipeline: recording completed")

			case 'm':
				// Get command line overrides for mix
				guitarVol, _ := cmd.Flags().GetFloat64("guitar-volume")
				backingVol, _ := cmd.Flags().GetFloat64("backing-volume")
				delay, _ := cmd.Flags().GetInt("delay")

				// Display effective values
				effectiveGuitarVol := cfg.Mix.GuitarVolume
				effectiveBackingVol := cfg.Mix.BackingVolume
				effectiveDelay := cfg.Mix.DelayMs

				if guitarVol > 0 {
					effectiveGuitarVol = guitarVol
				}
				if backingVol > 0 {
					effectiveBackingVol = backingVol
				}
				if delay >= 0 {
					effectiveDelay = delay
				}

				fmt.Printf("Mixing song: %s\n", songName)
				fmt.Printf("Guitar volume: %.1f\n", effectiveGuitarVol)
				fmt.Printf("Backing volume: %.1f\n", effectiveBackingVol)
				fmt.Printf("Backing track delay: %dms\n", effectiveDelay)

				mixer := mix.New(cfg)
				var err error
				if guitarVol > 0 || backingVol > 0 || delay >= 0 {
					err = mixer.MixWithOptions(songName, guitarVol, backingVol, delay)
				} else {
					err = mixer.Mix(songName)
				}

				if err != nil {
					return fmt.Errorf("pipeline mix failed: %w", err)
				}
				fmt.Println("Pipeline: mixing completed")

			case 'p':
				fmt.Printf("Playing song: %s\n", songName)
				player := play.New(cfg)
				if err := player.Play(songName); err != nil {
					return fmt.Errorf("pipeline play failed: %w", err)
				}
				fmt.Println("Pipeline: playback completed")

			default:
				return fmt.Errorf("unknown pipeline step: '%c' (valid: r=record, m=mix, p=play)", step)
			}
		}

		return nil
	},
}

func init() {
	// Add mix-specific flags to run command
	runCmd.Flags().Float64P("guitar-volume", "g", 0, "guitar volume (overrides config)")
	runCmd.Flags().Float64P("backing-volume", "b", 0, "backing volume (overrides config)")
	runCmd.Flags().IntP("delay", "d", -1, "backing track delay in ms (overrides config)")
	runCmd.Flags().StringP("output", "o", "", "output directory (overrides config)")
}