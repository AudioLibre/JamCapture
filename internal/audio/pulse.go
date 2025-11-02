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
	// Get PulseAudio sources
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

	// Add PipeWire/JACK sources (like Carla)
	pwSources, err := p.listPipeWireSources()
	if err == nil {
		sources = append(sources, pwSources...)
	}

	return sources, nil
}

func (p *PulseAudio) listPipeWireSources() ([]string, error) {
	cmd := exec.Command("pw-cli", "list-objects")
	output, err := cmd.Output()
	if err != nil {
		// PipeWire not available or pw-cli not found
		return nil, err
	}

	var sources []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		// Look for Carla or JACK audio outputs that can be used as sources
		if strings.Contains(line, "object.path") &&
		   (strings.Contains(line, "Carla") || strings.Contains(line, "JACK")) {
			// Extract the port name from object.path = "..."
			if start := strings.Index(line, `"`); start != -1 {
				if end := strings.Index(line[start+1:], `"`); end != -1 {
					portName := line[start+1 : start+1+end]
					// Add both Carla formats: Carla:output_X and Carla-Patchbay_X:output_X
					// Also check for audio-out patterns for Patchbay mode
					if (strings.Contains(portName, "output") && !strings.Contains(portName, "events")) ||
					   strings.Contains(portName, "audio-out") {
						sources = append(sources, portName)
					}
				}
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