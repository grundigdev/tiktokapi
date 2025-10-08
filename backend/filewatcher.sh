#!/bin/sh

echo "Starting file watcher for /app/videos/pending..."

# Process any existing files first
for file in /app/videos/pending/*; do
    if [ -f "$file" ]; then
        echo "Processing existing file: $file"
        ./backend -filepath "$file"
    fi
done

# Watch for new files
inotifywait -m -e close_write -e moved_to /app/videos/pending |
while read -r directory events filename; do
    filepath="${directory}${filename}"
    
    # Check if it's a video file
    case "$filename" in
        *.mp4|*.MP4|*.avi|*.AVI|*.mov|*.MOV|*.mkv|*.MKV)
            echo "New video detected: $filepath"
            ./backend -filepath "$filepath"
            ;;
        *)
            echo "Ignoring non-video file: $filename"
            ;;
    esac
done