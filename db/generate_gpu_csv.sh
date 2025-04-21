#!/bin/bash

# filepath: /Users/mac/Projects/gpu_request/scripts/generate_gpu_csv.sh

# Output CSV file
OUTPUT_FILE="gpus.csv"

# Check if nvidia-smi is installed
if ! command -v nvidia-smi &> /dev/null; then
    echo "Error: nvidia-smi is not installed or not in PATH."
    exit 1
fi

# Write the CSV header
echo "server_name,gpu_number,manufacturer,model_name,vram_size_mb" > "$OUTPUT_FILE"

# Get the server name (hostname)
SERVER_NAME=$(hostname)

# Parse nvidia-smi output
nvidia-smi --query-gpu=index,name,memory.total --format=csv,noheader | while IFS=',' read -r GPU_INDEX GPU_NAME VRAM_MB; do
    # Trim whitespace
    GPU_INDEX=$(echo "$GPU_INDEX" | xargs)
    GPU_NAME=$(echo "$GPU_NAME" | xargs)
    VRAM_MB=$(echo "$VRAM_MB" | xargs | sed 's/ MiB//')

    # Extract manufacturer (assume NVIDIA for nvidia-smi)
    MANUFACTURER="NVIDIA"

    # Write GPU details to the CSV file
    echo "$SERVER_NAME,$GPU_INDEX,$MANUFACTURER,$GPU_NAME,$VRAM_MB" >> "$OUTPUT_FILE"
done

echo "GPU details have been written to $OUTPUT_FILE"