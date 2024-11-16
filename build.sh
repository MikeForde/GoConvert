#!/bin/bash

# Output directory
OUTPUT_DIR="./bin"
mkdir -p $OUTPUT_DIR

# Build for Linux AMD64
echo "Building for linux/amd64..."
env GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/linux_amd64/hl7_to_fhir" .

# Check if the binary was created
if [ -f "$OUTPUT_DIR/linux_amd64/hl7_to_fhir" ]; then
    echo "Linux binary successfully created!"
else
    echo "Failed to create Linux binary."
fi
