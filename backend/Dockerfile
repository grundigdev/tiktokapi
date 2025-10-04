FROM golang:1.21.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o backend .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

# Binary vom builder stage kopieren
COPY --from=builder /app/backend .

EXPOSE 8000

CMD ["./backend"]