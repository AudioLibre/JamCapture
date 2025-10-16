#!/bin/bash

set -euxo pipefail

# Configuration parameters
SCARLETT_INPUT="alsa_input.usb-Focusrite_Scarlett_2i2_USB_Y814JK8264026F-00.analog-stereo"
SCARLETT_OUTPUT="alsa_output.usb-Focusrite_Scarlett_2i2_USB_Y814JK8264026F-00.analog-stereo.monitor"
SAMPLE_RATE="48000"
CHANNELS="2"
MIX_FILTER="amix=inputs=2:normalize=0:duration=longest,pan=stereo|c0=.5*c0+.5*c1|c1=.5*c0+.5*c1,volume=1.0"

# Function to display usage
show_usage() {
    echo "Usage: $0 <song_name>"
    echo "Example: $0 \"My Awesome Song\""
    exit 1
}

# Check if song name is provided
if [ $# -eq 0 ]; then
    echo "Error: Please provide a song name"
    show_usage
fi

SONG_NAME="$1"
# Clean song name for filename (remove special characters, replace spaces with underscores)
CLEAN_NAME=$(echo "$SONG_NAME" | sed 's/[^a-zA-Z0-9 ]//g' | tr ' ' '_')
MKV_FILE="$HOME/${CLEAN_NAME}.mkv"
FLAC_FILE="$HOME/${CLEAN_NAME}.flac"

echo "Song: $SONG_NAME"
echo "MKV file: $MKV_FILE"
echo "FLAC file: $FLAC_FILE"

# Remove existing files
rm -f "$MKV_FILE" "$FLAC_FILE"

echo "Starting recording..."
echo "Press Ctrl+C to stop recording"

# Start FFmpeg recording in background
ffmpeg \
    -f pulse -i "$SCARLETT_INPUT" \
    -f pulse -i "$SCARLETT_OUTPUT" \
    -map 0:a -map 1:a -c:a flac -ar "$SAMPLE_RATE" -ac "$CHANNELS" \
    "$MKV_FILE" &

FFMPEG_PID=$!

# Function to cleanup on exit
cleanup() {
    echo "Stopping recording..."
    kill $FFMPEG_PID 2>/dev/null || true
    wait $FFMPEG_PID 2>/dev/null || true

    if [ -f "$MKV_FILE" ]; then
        echo "Recording saved to: $MKV_FILE"
        echo "Creating mixed FLAC file..."

        # Create mixed FLAC file
        ffmpeg \
            -i "$MKV_FILE" \
            -filter_complex "$MIX_FILTER" \
            -ac "$CHANNELS" -ar "$SAMPLE_RATE" -c:a flac \
            "$FLAC_FILE" -y

        echo "Mixed FLAC file saved to: $FLAC_FILE"
    fi
}

# Set trap to handle Ctrl+C
trap cleanup EXIT INT TERM

echo "🎤 RECORDING..."

# Wait for user to stop recording
while kill -0 $FFMPEG_PID 2>/dev/null; do
    sleep 1
done