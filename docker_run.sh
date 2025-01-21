#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 [-f | -b]"
    echo "  -f    Run the frontend container (frontend:latest)"
    echo "  -b    Run the backend container (backend:latest)"
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
      echo "Running frontend container..."
      sudo docker run -it --rm -p 8081:80 frontend:latest
      ;;
    b)
      echo "Running backend container..."
      sudo docker run -it --rm -p 3000:3000 backend:latest
      ;;
    *)
      usage
      ;;
  esac
  exit 0

done

# If no valid options were provided, display usage
usage

