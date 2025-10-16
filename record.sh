#!/bin/bash

set -euo pipefail

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
SCRIPT_DIR="$(dirname "$0")"

# Record audio then mix it
echo "Starting recording phase..."
if "$SCRIPT_DIR/record-audio.sh" "$SONG_NAME"; then
    echo "Recording completed successfully"
else
    echo "Recording interrupted by user (Ctrl+C), checking if file was created..."
fi

# Check if recording file exists before mixing
CLEAN_NAME=$(echo "$SONG_NAME" | sed 's/[^a-zA-Z0-9 ]//g' | tr ' ' '_')
MKV_FILE="$HOME/${CLEAN_NAME}.mkv"

if [ -f "$MKV_FILE" ]; then
    echo "Recording file found, proceeding to mix..."
    "$SCRIPT_DIR/mix-audio.sh" "$SONG_NAME"
else
    echo "Error: No recording file found at $MKV_FILE"
    echo "Recording may have failed or been cancelled too early"
    exit 1
fi