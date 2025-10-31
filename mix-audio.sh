#!/bin/bash

set -euxo pipefail

# Configuration parameters
SAMPLE_RATE="48000"
CHANNELS="2"

# Default mix levels
DEFAULT_GUITAR_VOLUME="4"    # Input channel volume (guitar/instrument)
DEFAULT_OTHER_VOLUME="0.8"     # Output monitor volume (backing track/computer audio)
DEFAULT_DELAY_MS="0"           # Delay in milliseconds to compensate for Bluetooth latency

# Function to display usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [song_name]"
    echo "       JAMCAPTURE_SONG=\"My Song\" $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -g, --guitar-volume VOLUME   Volume level for guitar/input channel (default: $DEFAULT_GUITAR_VOLUME)"
    echo "  -b, --backing-volume VOLUME  Volume level for backing track/computer audio (default: $DEFAULT_OTHER_VOLUME)"
    echo "  -d, --delay MILLISECONDS     Delay in milliseconds for backing track to sync with guitar played late due to Bluetooth (default: $DEFAULT_DELAY_MS)"
    echo "  -o, --output-dir DIRECTORY   Directory containing the MKV and where to save FLAC (default: \$HOME/Audio/JamCapture)"
    echo "  -h, --help                   Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 \"My Awesome Song\""
    echo "  $0 -g 2.0 -b 0.5 \"My Song\""
    echo "  $0 --delay 150 --guitar-volume 3.0 \"My Song\""
    echo "  $0 -d 150 -g 2.0 -b 0.5 -o \"/custom/path\" \"My Song\""
    echo "  JAMCAPTURE_SONG=\"My Song\" $0 -d 150"
    exit 1
}

# Initialize variables with defaults
GUITAR_VOLUME="$DEFAULT_GUITAR_VOLUME"
OTHER_VOLUME="$DEFAULT_OTHER_VOLUME"
DELAY_MS="$DEFAULT_DELAY_MS"
OUTPUT_DIR="$HOME/Audio/JamCapture"
SONG_NAME=""

# Parse command line options
while [[ $# -gt 0 ]]; do
    case $1 in
        -g|--guitar-volume)
            GUITAR_VOLUME="$2"
            shift 2
            ;;
        -b|--backing-volume)
            OTHER_VOLUME="$2"
            shift 2
            ;;
        -d|--delay)
            DELAY_MS="$2"
            shift 2
            ;;
        -o|--output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            ;;
        -*)
            echo "Unknown option $1"
            show_usage
            ;;
        *)
            # This should be the song name
            if [ -z "$SONG_NAME" ]; then
                SONG_NAME="$1"
            else
                echo "Error: Multiple song names provided: '$SONG_NAME' and '$1'"
                show_usage
            fi
            shift
            ;;
    esac
done

# Get song name from environment variable if not provided as argument
if [ -z "$SONG_NAME" ] && [ -n "${JAMCAPTURE_SONG:-}" ]; then
    SONG_NAME="$JAMCAPTURE_SONG"
fi

# Check if we have a song name
if [ -z "$SONG_NAME" ]; then
    echo "Error: Please provide a song name or set JAMCAPTURE_SONG environment variable"
    show_usage
fi

# Apply delay if specified - delay backing track to sync with guitar played late due to Bluetooth
if [ "$DELAY_MS" -gt 0 ]; then
    MIX_FILTER="[0:0]volume=${GUITAR_VOLUME}[guitar];[0:1]volume=${OTHER_VOLUME},adelay=${DELAY_MS}|${DELAY_MS}[other];[guitar][other]amix=inputs=2:normalize=0"
else
    MIX_FILTER="[0:0]volume=${GUITAR_VOLUME}[guitar];[0:1]volume=${OTHER_VOLUME}[other];[guitar][other]amix=inputs=2:normalize=0"
fi
# Clean song name for filename (remove special characters, replace spaces with underscores)
CLEAN_NAME=$(echo "$SONG_NAME" | sed 's/[^a-zA-Z0-9 ]//g' | tr ' ' '_')

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

MKV_FILE="$OUTPUT_DIR/${CLEAN_NAME}.mkv"
FLAC_FILE="$OUTPUT_DIR/${CLEAN_NAME}.flac"

echo "Song: $SONG_NAME"
echo "Guitar volume: $GUITAR_VOLUME"
echo "Backing volume: $OTHER_VOLUME"
echo "Backing track delay: ${DELAY_MS}ms"
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
