package mix

import (
	"fmt"
	"github.com/audiolibre/jamcapture/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Mixer struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Mixer {
	return &Mixer{cfg: cfg}
}

func (m *Mixer) Mix(songName string) error {
	cleanName := m.cleanFileName(songName)
	inputFile := filepath.Join(m.cfg.Output.Directory, cleanName+".mkv")
	outputFile := filepath.Join(m.cfg.Output.Directory, cleanName+"."+m.cfg.Output.Format)

	// Check if input file exists
	if _, err := os.Stat(inputFile); err != nil {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	// Remove existing output file
	os.Remove(outputFile)

	// Build FFmpeg filter based on delay configuration
	var mixFilter string
	if m.cfg.Mix.DelayMs > 0 {
		// Delay backing track to sync with guitar played late due to Bluetooth
		mixFilter = fmt.Sprintf(
			"[0:0]volume=%.1f[guitar];[0:1]volume=%.1f,adelay=%d|%d[other];[guitar][other]amix=inputs=2:normalize=0",
			m.cfg.Mix.GuitarVolume,
			m.cfg.Mix.BackingVolume,
			m.cfg.Mix.DelayMs,
			m.cfg.Mix.DelayMs,
		)
	} else {
		// No delay
		mixFilter = fmt.Sprintf(
			"[0:0]volume=%.1f[guitar];[0:1]volume=%.1f[other];[guitar][other]amix=inputs=2:normalize=0",
			m.cfg.Mix.GuitarVolume,
			m.cfg.Mix.BackingVolume,
		)
	}

	fmt.Printf("Input file: %s\n", inputFile)
	fmt.Printf("Output file: %s\n", outputFile)
	fmt.Printf("Mix filter: %s\n", mixFilter)
	fmt.Println("Creating mixed audio file...")

	// Prepare FFmpeg command
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-filter_complex", mixFilter,
		"-ac", fmt.Sprintf("%d", m.cfg.Audio.Channels),
		"-ar", fmt.Sprintf("%d", m.cfg.Audio.SampleRate),
		"-c:a", m.cfg.Output.Format,
		"-y", // Overwrite output file
		outputFile,
	)

	// Run FFmpeg
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg mixing failed: %w\nOutput: %s", err, string(output))
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); err != nil {
		return fmt.Errorf("output file not created: %s", outputFile)
	}

	fmt.Printf("Mixed audio file saved to: %s\n", outputFile)
	return nil
}

func (m *Mixer) MixWithOptions(songName string, guitarVol, backingVol float64, delayMs int) error {
	// Temporarily override config values
	originalGuitarVol := m.cfg.Mix.GuitarVolume
	originalBackingVol := m.cfg.Mix.BackingVolume
	originalDelay := m.cfg.Mix.DelayMs

	if guitarVol > 0 {
		m.cfg.Mix.GuitarVolume = guitarVol
	}
	if backingVol > 0 {
		m.cfg.Mix.BackingVolume = backingVol
	}
	if delayMs >= 0 {
		m.cfg.Mix.DelayMs = delayMs
	}

	err := m.Mix(songName)

	// Restore original values
	m.cfg.Mix.GuitarVolume = originalGuitarVol
	m.cfg.Mix.BackingVolume = originalBackingVol
	m.cfg.Mix.DelayMs = originalDelay

	return err
}

func (m *Mixer) cleanFileName(name string) string {
	// Remove special characters and replace spaces with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	cleaned := reg.ReplaceAllString(name, "")
	return strings.ReplaceAll(strings.TrimSpace(cleaned), " ", "_")
}