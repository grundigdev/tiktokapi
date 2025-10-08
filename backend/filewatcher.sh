#!/bin/sh

echo "Starting file watcher for /home/marcel/dev/scripts/go/backend/videos/pending..."

# Process any existing files first
for file in /home/marcel/dev/scripts/go/backend/videos/pending/*; do
    if [ -f "$file" ]; then
        echo "Processing existing file: $file"
        ./backend -video "/home/marcel/dev/scripts/go/backend/videos/pending/video.mp4" -context "/home/marcel/dev/scripts/go/backend/title.json"
    fi
done

# Watch for new files
inotifywait -m -e close_write -e moved_to /home/marcel/dev/scripts/go/backend/videos/pending/ |
while read -r directory events filename; do
    filepath="${directory}${filename}"
    
    # Check if it's a video file
    case "$filename" in
        *.mp4|*.MP4|*.avi|*.AVI|*.mov|*.MOV|*.mkv|*.MKV)
            echo "New video detected: $filepath"
            ./backend -video "/home/marcel/dev/scripts/go/backend/videos/pending/video.mp4" -context "/home/marcel/dev/scripts/go/backend/title.json"
            ;;
        *)
            echo "Ignoring non-video file: $filename"
            ;;
    esac
done