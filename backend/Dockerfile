FROM golang:1.21.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o backend .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata inotify-tools

WORKDIR /app

# Binary vom builder stage kopieren
COPY --from=builder /app/backend .

# Copy file watcher script
COPY filewatcher.sh .
RUN chmod +x filewatcher.sh

RUN mkdir -p /app/videos/pending \
             /app/videos/uploading \
             /app/videos/uploaded \
             /app/videos/failed \
             /app/context/pending \
             /app/context/uploaded \
             /app/context/failed \
             /app/context/uploading


RUN chmod -R 755 /app/videos
RUN chmod -R 755 /app/context

EXPOSE 8000

CMD ["./filewatcher.sh"]