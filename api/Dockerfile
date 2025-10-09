FROM golang:1.21.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o api ./cmd/api

# Build the migration binaries
RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o migrate_down ./database/migrate_down.go
RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o migrate_up ./database/migrate_up.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app/

# Copy all binaries
COPY --from=builder /app/api .
COPY --from=builder /app/migrate_down .
COPY --from=builder /app/migrate_up .

COPY --from=builder /app/.env .
COPY --from=builder /app/mailer/templates ./mailer/templates

EXPOSE 8080

# Create a startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo './migrate_down' >> /app/start.sh && \
    echo './migrate_up' >> /app/start.sh && \
    echo './api' >> /app/start.sh && \
    chmod +x /app/start.sh

CMD ["/app/start.sh"]