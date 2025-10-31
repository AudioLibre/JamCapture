# JamCapture Go

JamCapture reimplemented in Go with Cobra CLI framework.

## Features

- **Record**: Capture guitar input and system audio simultaneously
- **Mix**: Combine tracks with volume control and Bluetooth latency compensation
- **Play**: Play the mixed results
- **Chain**: Execute commands in sequence (record → mix → play)
- **Configuration**: YAML-based configuration with sane defaults

## Installation

```bash
go build -o jamcapture
```

## Configuration

Configuration is stored in `~/.config/jamcapture.yaml`:

```yaml
audio:
  sample_rate: 48000
  channels: 2

record:
  guitar_input: "alsa_input.usb-Focusrite_Scarlett_2i2_USB_Y814JK8264026F-00.analog-stereo"
  monitor_input: ""  # Auto-detect if empty

mix:
  guitar_volume: 4.0
  backing_volume: 0.8
  delay_ms: 200  # Bluetooth compensation delay

output:
  directory: "~/Audio/JamCapture"
  format: "flac"
```

### Pre-configured Templates

The project includes several configuration templates in `jamcapture-go/`:

- **`jamcapture-scarlett.yaml`**: Focusrite Scarlett setup with 200ms Bluetooth delay
- **`jamcapture-wired.yaml`**: Scarlett setup with no delay (wired speakers)
- **`jamcapture-quiet.yaml`**: Lower volumes for quiet practice sessions

Use with: `jamcapture --config jamcapture-go/jamcapture-scarlett.yaml -p rmp "song"`

## Usage

### Simplified Syntax (Recommended)

Use `-p` to specify pipeline steps directly on song name:

```bash
# Record, mix, and play in one command
jamcapture -p rmp "My Song"

# Mix and play with custom delay
jamcapture -p mp -d 180 "My Song"

# Just mix with custom settings
jamcapture -p m -g 3.0 -b 0.6 -d 150 "My Song"

# Record and mix only (no playback)
jamcapture -p rm "My Song"
```

Pipeline steps:
- `r` = record
- `m` = mix
- `p` = play

### Traditional Commands (Still Available)

```bash
# List available audio sources
jamcapture sources

# Individual commands
jamcapture record "My Song"
jamcapture mix "My Song" -d 150
jamcapture play "My Song"

# With pipeline chaining
jamcapture -p mp mix "My Song"
```

### Configuration Management

```bash
# View current configuration
jamcapture config show

# Edit configuration (opens in $EDITOR)
jamcapture config edit
```

## Bluetooth Latency Compensation

When using Bluetooth speakers/headphones:

1. You hear the backing track with ~200ms delay
2. You naturally play guitar late to match what you hear
3. In the recording, guitar is behind the backing track
4. The `delay_ms` setting delays the backing track in the mix to sync with your guitar

Common Bluetooth delays:
- **Standard A2DP**: 180-250ms
- **aptX Low Latency**: 40-80ms
- **aptX Standard**: 100-150ms

## Examples

```bash
# Quick jam session with full pipeline
jamcapture -p rmp "blues_jam"

# Re-mix existing recording with different settings
jamcapture -p m -d 180 -g 3.5 "blues_jam"

# Test different delay values with immediate playback
jamcapture -p mp -d 150 "blues_jam"
jamcapture -p mp -d 200 "blues_jam"
jamcapture -p mp -d 250 "blues_jam"
```

## Files

- Input recordings: `~/Audio/JamCapture/{song}.mkv`
- Mixed output: `~/Audio/JamCapture/{song}.flac`
- Configuration: `~/.config/jamcapture.yaml`

## Requirements

- FFmpeg
- PulseAudio/PipeWire
- VLC or other audio player for playback