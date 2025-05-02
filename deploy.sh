#!/bin/bash

# Simple deploy script for a Go application binary

# Variables
APP_NAME="word-dict"
APP_BINARY="./$APP_NAME"
LOG_FILE="deploy.log"

# Check if the binary exists
if [ ! -f "$APP_BINARY" ]; then
    echo "Error: Binary '$APP_BINARY' not found!"
    exit 1
fi

# Make the binary executable
chmod +x "$APP_BINARY"

# Run the binary and log output
echo "Starting $APP_NAME at $(date)" >> "$LOG_FILE"
nohup "$APP_BINARY" >> "$LOG_FILE" 2>&1 &

echo "Deployment successful! $APP_NAME is running."
echo "Logs are being written to $LOG_FILE"