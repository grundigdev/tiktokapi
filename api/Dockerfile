FROM golang:1.21.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app/

COPY --from=builder /app/api .

COPY --from=builder /app/.env .

COPY --from=builder /app/mailer/templates ./mailer/templates

EXPOSE 8080

CMD ["./api"]
