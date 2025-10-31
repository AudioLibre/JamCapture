#!/bin/bash

set -euxo pipefail

# Configuration parameters
SCARLETT_INPUT="alsa_input.usb-Focusrite_Scarlett_2i2_USB_Y814JK8264026F-00.analog-stereo"
SAMPLE_RATE="48000"
CHANNELS="2"

# Get default system audio output monitor
DEFAULT_SINK=$(pactl get-default-sink)
DEFAULT_OUTPUT_MONITOR="${DEFAULT_SINK}.monitor"

echo "Using Scarlett input: $SCARLETT_INPUT"
echo "Using default system output monitor: $DEFAULT_OUTPUT_MONITOR"

# Function to display usage
show_usage() {
    echo "Usage: $0 [song_name] [output_dir]"
    echo "Example: $0 \"My Awesome Song\""
    echo "         $0 \"My Awesome Song\" \"/path/to/custom/dir\""
    echo "         JAMCAPTURE_SONG=\"My Song\" $0"
    echo ""
    echo "Options:"
    echo "  song_name   Name of the song (can be overridden by JAMCAPTURE_SONG env var)"
    echo "  output_dir  Directory to save the recording (default: \$HOME/Audio/JamCapture)"
    exit 1
}

# Get song name from argument or environment variable
if [ -n "${JAMCAPTURE_SONG:-}" ]; then
    # JAMCAPTURE_SONG is set, use it as song name
    SONG_NAME="$JAMCAPTURE_SONG"
    OUTPUT_DIR="${1:-$HOME/Audio/JamCapture}"
elif [ $# -gt 0 ] && [ -n "${1:-}" ]; then
    # No JAMCAPTURE_SONG, use first argument as song name
    SONG_NAME="$1"
    OUTPUT_DIR="${2:-$HOME/Audio/JamCapture}"
else
    # Neither exists
    echo "Error: Please provide a song name or set JAMCAPTURE_SONG environment variable"
    show_usage
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Clean song name for filename (remove special characters, replace spaces with underscores)
CLEAN_NAME=$(echo "$SONG_NAME" | sed 's/[^a-zA-Z0-9 ]//g' | tr ' ' '_')
MKV_FILE="$OUTPUT_DIR/${CLEAN_NAME}.mkv"

echo "Song: $SONG_NAME"
echo "MKV file: $MKV_FILE"

# Remove existing file
rm -f "$MKV_FILE"

echo "Starting recording..."
echo "Press Ctrl+C to stop recording"

# Start FFmpeg recording in background
ffmpeg \
    -f pulse -i "$SCARLETT_INPUT" \
    -f pulse -i "$DEFAULT_OUTPUT_MONITOR" \
    -filter_complex "[0]pan=stereo|c0=0.5*c0+0.5*c1|c1=0.5*c0+0.5*c1[guitar];[1]acopy[backing]" \
    -map "[guitar]" -map "[backing]" -c:a flac -ar "$SAMPLE_RATE" \
    "$MKV_FILE" &

FFMPEG_PID=$!

# Function to cleanup on exit
cleanup() {
    echo "Stopping recording..."
    kill $FFMPEG_PID 2>/dev/null || true
    wait $FFMPEG_PID 2>/dev/null || true

    if [ -f "$MKV_FILE" ]; then
        echo "Recording saved to: $MKV_FILE"
    fi
}

# Set trap to handle Ctrl+C
trap cleanup EXIT INT TERM

echo "🎤 RECORDING..."

# Wait for user to stop recording
while kill -0 $FFMPEG_PID 2>/dev/null; do
    sleep 1
done