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
"$SCRIPT_DIR/record-audio.sh" "$SONG_NAME"
"$SCRIPT_DIR/mix-audio.sh" "$SONG_NAME"