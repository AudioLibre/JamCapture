package cmd

import (
	"fmt"
	"github.com/audiolibre/jamcapture/internal/mix"
	"github.com/audiolibre/jamcapture/internal/play"
	"github.com/audiolibre/jamcapture/internal/record"
	"strings"
)

func executePipeline(songName string, startStep rune) error {
	if pipeline == "" {
		return nil
	}

	steps := []rune(strings.ToLower(pipeline))

	// Find the starting position in the pipeline
	startIndex := -1
	for i, step := range steps {
		if step == startStep {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return fmt.Errorf("step '%c' not found in pipeline '%s'", startStep, pipeline)
	}

	// Execute remaining steps in the pipeline
	for i := startIndex + 1; i < len(steps); i++ {
		step := steps[i]
		fmt.Printf("Pipeline: executing step '%c'...\n", step)

		switch step {
		case 'r':
			recorder := record.New(cfg)
			if err := recorder.Record(songName); err != nil {
				return fmt.Errorf("pipeline record failed: %w", err)
			}
			fmt.Println("Pipeline: recording completed")

		case 'm':
			mixer := mix.New(cfg)
			if err := mixer.Mix(songName); err != nil {
				return fmt.Errorf("pipeline mix failed: %w", err)
			}
			fmt.Println("Pipeline: mixing completed")

		case 'p':
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
}

func validatePipeline() error {
	if pipeline == "" {
		return nil
	}

	validSteps := map[rune]bool{
		'r': true, // record
		'm': true, // mix
		'p': true, // play
	}

	steps := []rune(strings.ToLower(pipeline))
	for _, step := range steps {
		if !validSteps[step] {
			return fmt.Errorf("invalid pipeline step: '%c' (valid: r=record, m=mix, p=play)", step)
		}
	}

	return nil
}