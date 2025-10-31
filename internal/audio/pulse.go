package audio

import (
	"fmt"
	"os/exec"
	"strings"
)

type PulseAudio struct{}

func NewPulseAudio() *PulseAudio {
	return &PulseAudio{}
}

func (p *PulseAudio) GetDefaultSink() (string, error) {
	cmd := exec.Command("pactl", "get-default-sink")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default sink: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (p *PulseAudio) GetDefaultSinkMonitor() (string, error) {
	sink, err := p.GetDefaultSink()
	if err != nil {
		return "", err
	}
	return sink + ".monitor", nil
}

func (p *PulseAudio) ListSources() ([]string, error) {
	cmd := exec.Command("pactl", "list", "short", "sources")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var sources []string
	for _, line := range lines {
		if line != "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				sources = append(sources, parts[1])
			}
		}
	}
	return sources, nil
}

func (p *PulseAudio) ValidateSource(sourceName string) error {
	sources, err := p.ListSources()
	if err != nil {
		return err
	}

	for _, source := range sources {
		if source == sourceName {
			return nil
		}
	}
	return fmt.Errorf("source %s not found", sourceName)
}