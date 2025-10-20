#!/bin/bash

set -euxo pipefail

# Configuration parameters
SAMPLE_RATE="48000"
CHANNELS="2"

# Default mix levels
DEFAULT_GUITAR_VOLUME="2"    # Input channel volume (guitar/instrument)
DEFAULT_OTHER_VOLUME="0.8"     # Output monitor volume (backing track/computer audio)

# Function to display usage
show_usage() {
    echo "Usage: $0 [song_name] [guitar_volume] [backing_volume] [output_dir]"
    echo "Example: $0 \"My Awesome Song\""
    echo "         $0 \"My Awesome Song\" 2.0 0.5"
    echo "         $0 \"My Awesome Song\" 2.0 0.5 \"/path/to/custom/dir\""
    echo "         JAMCAPTURE_SONG=\"My Song\" $0"
    echo ""
    echo "Options:"
    echo "  song_name       Name of the song (can be overridden by JAMCAPTURE_SONG env var)"
    echo "  guitar_volume   Volume level for guitar/input channel (default: $DEFAULT_GUITAR_VOLUME)"
    echo "  backing_volume  Volume level for backing track/computer audio (default: $DEFAULT_OTHER_VOLUME)"
    echo "  output_dir      Directory containing the MKV and where to save FLAC (default: \$HOME/Audio/JamCapture)"
    exit 1
}

# Get song name from argument or environment variable
if [ -n "${JAMCAPTURE_SONG:-}" ]; then
    # JAMCAPTURE_SONG is set, use it as song name
    SONG_NAME="$JAMCAPTURE_SONG"
    GUITAR_VOLUME="${1:-$DEFAULT_GUITAR_VOLUME}"
    OTHER_VOLUME="${2:-$DEFAULT_OTHER_VOLUME}"
    OUTPUT_DIR="${3:-$HOME/Audio/JamCapture}"
elif [ $# -gt 0 ] && [ -n "${1:-}" ]; then
    # No JAMCAPTURE_SONG, use first argument as song name
    SONG_NAME="$1"
    GUITAR_VOLUME="${2:-$DEFAULT_GUITAR_VOLUME}"
    OTHER_VOLUME="${3:-$DEFAULT_OTHER_VOLUME}"
    OUTPUT_DIR="${4:-$HOME/Audio/JamCapture}"
else
    # Neither exists
    echo "Error: Please provide a song name or set JAMCAPTURE_SONG environment variable"
    show_usage
fi

MIX_FILTER="[0:0]volume=${GUITAR_VOLUME}[guitar];[0:1]volume=${OTHER_VOLUME}[other];[guitar][other]amix=inputs=2:normalize=0"
# Clean song name for filename (remove special characters, replace spaces with underscores)
CLEAN_NAME=$(echo "$SONG_NAME" | sed 's/[^a-zA-Z0-9 ]//g' | tr ' ' '_')

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

MKV_FILE="$OUTPUT_DIR/${CLEAN_NAME}.mkv"
FLAC_FILE="$OUTPUT_DIR/${CLEAN_NAME}.flac"

echo "Song: $SONG_NAME"
echo "Guitar volume: $GUITAR_VOLUME"
echo "Backing volume: $OTHER_VOLUME"
echo "Output directory: $OUTPUT_DIR"
echo "MKV file: $MKV_FILE"
echo "FLAC file: $FLAC_FILE"

# Check if MKV file exists
if [ ! -f "$MKV_FILE" ]; then
    echo "Error: MKV file not found: $MKV_FILE"
    exit 1
fi

# Remove existing FLAC file
rm -f "$FLAC_FILE"

echo "Creating mixed FLAC file..."

# Create mixed FLAC file
ffmpeg \
    -i "$MKV_FILE" \
    -filter_complex "$MIX_FILTER" \
    -ac "$CHANNELS" -ar "$SAMPLE_RATE" -c:a flac \
    "$FLAC_FILE" -y

echo "Mixed FLAC file saved to: $FLAC_FILE"
