package record

import (
	"fmt"
	"github.com/audiolibre/jamcapture/internal/audio"
	"github.com/audiolibre/jamcapture/internal/config"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

type Recorder struct {
	cfg   *config.Config
	pulse *audio.PulseAudio
}

func New(cfg *config.Config) *Recorder {
	return &Recorder{
		cfg:   cfg,
		pulse: audio.NewPulseAudio(),
	}
}

func (r *Recorder) Record(songName string) error {
	// Validate and prepare inputs
	guitarInput := r.cfg.Record.GuitarInput
	if r.cfg.Record.MonitorInput == "" {
		monitor, err := r.pulse.GetDefaultSinkMonitor()
		if err != nil {
			return fmt.Errorf("failed to get default sink monitor: %w", err)
		}
		r.cfg.Record.MonitorInput = monitor
	}

	// Validate guitar input
	if err := r.pulse.ValidateSource(guitarInput); err != nil {
		return fmt.Errorf("invalid guitar input: %w", err)
	}

	// Validate monitor input
	if err := r.pulse.ValidateSource(r.cfg.Record.MonitorInput); err != nil {
		return fmt.Errorf("invalid monitor input: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(r.cfg.Output.Directory, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Clean song name for filename
	cleanName := r.cleanFileName(songName)
	outputFile := filepath.Join(r.cfg.Output.Directory, cleanName+".mkv")

	// Remove existing file
	os.Remove(outputFile)

	fmt.Printf("Guitar input: %s\n", guitarInput)
	fmt.Printf("Monitor input: %s\n", r.cfg.Record.MonitorInput)
	fmt.Printf("Output file: %s\n", outputFile)
	fmt.Println("Starting recording... Press Ctrl+C to stop")

	// Prepare FFmpeg command
	cmd := exec.Command("ffmpeg",
		"-f", "pulse", "-i", guitarInput,
		"-f", "pulse", "-i", r.cfg.Record.MonitorInput,
		"-filter_complex", "[0]pan=stereo|c0=0.5*c0+0.5*c1|c1=0.5*c0+0.5*c1[guitar];[1]acopy[backing]",
		"-map", "[guitar]", "-map", "[backing]",
		"-c:a", "flac",
		"-ar", fmt.Sprintf("%d", r.cfg.Audio.SampleRate),
		outputFile,
	)

	// Start recording
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}

	// Handle interruption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	fmt.Println("🎤 RECORDING...")

	select {
	case <-sigChan:
		fmt.Println("\nStopping recording...")
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			cmd.Process.Kill()
		}
		<-done // Wait for process to finish
	case err := <-done:
		if err != nil {
			return fmt.Errorf("recording process failed: %w", err)
		}
	}

	// Check if file was created
	if _, err := os.Stat(outputFile); err != nil {
		return fmt.Errorf("recording file not found: %s", outputFile)
	}

	fmt.Printf("Recording saved to: %s\n", outputFile)
	return nil
}

func (r *Recorder) cleanFileName(name string) string {
	// Remove special characters and replace spaces with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	cleaned := reg.ReplaceAllString(name, "")
	return strings.ReplaceAll(strings.TrimSpace(cleaned), " ", "_")
}