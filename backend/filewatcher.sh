#!/bin/sh

#PROD 
#basepath="/home/marcel/app/backend"

#CONTAINER
basepath="/app"

#DEV
#basepath="/home/marcel/dev/scripts/go/backend"

echo "Starting file watcher for ${basepath}/videos/pending and ${basepath}/context/pending..."

# Function to check if both directories have files and process them
check_and_process() {
    video_dir="${basepath}/videos/pending"
    context_dir="${basepath}/context/pending"
    
    # Get the first video file
    video_file=$(find "$video_dir" -maxdepth 1 -type f \( -iname "*.mp4" -o -iname "*.avi" -o -iname "*.mov" -o -iname "*.mkv" \) | head -n 1)
    
    # Get the first JSON file from context directory
    context_file=$(find "$context_dir" -maxdepth 1 -type f -iname "*.json" | head -n 1)
    
    # Check if both files exist
    if [ -n "$video_file" ] && [ -n "$context_file" ]; then
        echo "Processing video: $video_file"
        echo "Using context: $context_file"
        ./backend -video "$video_file" -context "$context_file"
    else
        if [ -z "$video_file" ]; then
            echo "Waiting for video file in ${video_dir}"
        fi
        if [ -z "$context_file" ]; then
            echo "Waiting for context file in ${context_dir}"
        fi
    fi
}

# Process existing files if both directories have content
check_and_process

# Watch both directories for new files
(inotifywait -m -e close_write -e moved_to ${basepath}/videos/pending/ & 
 inotifywait -m -e close_write -e moved_to ${basepath}/context/pending/) |
while read -r directory events filename; do
    filepath="${directory}${filename}"
    
    # Determine which directory triggered the event
    case "$directory" in
        *videos/pending*)
            case "$filename" in
                *.mp4|*.MP4|*.avi|*.AVI|*.mov|*.MOV|*.mkv|*.MKV)
                    echo "New video detected: $filepath"
                    check_and_process
                    ;;
                *)
                    echo "Ignoring non-video file: $filename"
                    ;;
            esac
            ;;
        *context/pending*)
            case "$filename" in
                *.json|*.JSON)
                    echo "New context file detected: $filepath"
                    check_and_process
                    ;;
                *)
                    echo "Ignoring non-JSON file: $filename"
                    ;;
            esac
            ;;
    esac
done