#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 [-f | -b]"
    echo "  -f    Build the frontend container (frontend:latest)"
    echo "  -b    Build the backend container (backend:latest)"
    exit 1
}

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    usage
fi

# Process command-line options
while getopts ":fb" opt; do
  case $opt in
    f)
      echo "Building frontend container..."
      sudo docker build -f Dockerfile.frontend -t frontend:latest .
      ;;
    b)
      echo "Building backend container..."
      sudo docker build -f Dockerfile.backend -t backend:latest .
      ;;
    *)
      usage
      ;;
  esac
  exit 0

done

# If no valid options were provided, display usage
usage

