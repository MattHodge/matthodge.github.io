#!/bin/bash
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# Get the absolute path of the input file
DIAGRAM_PATH=$(realpath $1)

# Extract the post name from the directory structure
POST_NAME=$(basename $(dirname $DIAGRAM_PATH))

# Set the destination path and file
DESTINATION_PATH="_posts/diagrams/$POST_NAME"
DESTINATION_FILE="$(basename $DIAGRAM_PATH .d2).svg"

# Log some info to the console
echo "Generating diagrams for $POST_NAME"
echo "Destination path: $DESTINATION_PATH/$DESTINATION_FILE"

# Extract layout and sketch flags from the file name
D2_LAYOUT="dagre" # default to dagre
D2_SKETCH="false"
FILENAME=$(basename "$DIAGRAM_PATH" .d2)

if [[ $FILENAME == *".elk"* ]]; then
  D2_LAYOUT="elk"
fi

if [[ $FILENAME == *".sketch"* ]]; then
  D2_SKETCH="true"
fi

# Run the d2 command with the extracted flags
echo Running with D2_LAYOUT=$D2_LAYOUT D2_SKETCH=$D2_SKETCH
D2_LAYOUT=elk D2_SKETCH=$D2_SKETCH d2 --pad 0 "$DIAGRAM_PATH" "$DESTINATION_PATH/$DESTINATION_FILE"
