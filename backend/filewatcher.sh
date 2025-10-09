#!/bin/sh



#PROD 
#basepath="/home/marcel/app/backend"

#CONTAINER
basepath="/app"

#DEV
#basepath="/home/marcel/dev/scripts/go/backend"


echo "Starting file watcher for ${basepath}/videos/pending..."

# Process any existing files first
for file in ${basepath}/videos/pending/*; do
    if [ -f "$file" ]; then
        echo "Processing existing file: $file"
        ./backend -video "$file" -context "${basepath}/title.json"
    fi
done

# Watch for new files
inotifywait -m -e close_write -e moved_to ${basepath}/videos/pending/ |
while read -r directory events filename; do
    filepath="${directory}${filename}"
    
    # Check if it's a video file
    case "$filename" in
        *.mp4|*.MP4|*.avi|*.AVI|*.mov|*.MOV|*.mkv|*.MKV)
            echo "New video detected: $filepath"
            ./backend -video "$filepath" -context "${basepath}/title.json"
            ;;
        *)
            echo "Ignoring non-video file: $filename"
            ;;
    esac
done