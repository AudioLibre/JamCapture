#!/bin/bash

set -euo pipefail

# Function to display usage
show_usage() {
    echo "Usage: $0 [song_name] [output_dir]"
    echo "Example: $0 \"My Awesome Song\""
    echo "         $0 \"My Awesome Song\" \"/path/to/custom/dir\""
    echo "         JAMCAPTURE_SONG=\"My Song\" $0"
    echo ""
    echo "Options:"
    echo "  song_name   Name of the song (can be overridden by JAMCAPTURE_SONG env var)"
    echo "  output_dir  Directory to save recordings (default: \$HOME/Audio/JamCapture)"
    exit 1
}

# Get song name from argument or environment variable
if [ -n "${JAMCAPTURE_SONG:-}" ]; then
    SONG_NAME="$JAMCAPTURE_SONG"
    OUTPUT_DIR="${1:-$HOME/Audio/JamCapture}"
elif [ $# -eq 0 ]; then
    echo "Error: Please provide a song name or set JAMCAPTURE_SONG environment variable"
    show_usage
else
    SONG_NAME="$1"
    OUTPUT_DIR="${2:-$HOME/Audio/JamCapture}"
fi

SCRIPT_DIR="$(dirname "$0")"

# Record audio then mix it
echo "Starting recording phase..."
if JAMCAPTURE_SONG="$SONG_NAME" "$SCRIPT_DIR/record-audio.sh" "$OUTPUT_DIR"; then
    echo "Recording completed successfully"
else
    echo "Recording interrupted by user (Ctrl+C), checking if file was created..."
fi

# Check if recording file exists before mixing
CLEAN_NAME=$(echo "$SONG_NAME" | sed 's/[^a-zA-Z0-9 ]//g' | tr ' ' '_')
MKV_FILE="$OUTPUT_DIR/${CLEAN_NAME}.mkv"

if [ -f "$MKV_FILE" ]; then
    echo "Recording file found, proceeding to mix..."
    JAMCAPTURE_SONG="$SONG_NAME" "$SCRIPT_DIR/mix-audio.sh" "" "" "$OUTPUT_DIR"
else
    echo "Error: No recording file found at $MKV_FILE"
    echo "Recording may have failed or been cancelled too early"
    exit 1
fi