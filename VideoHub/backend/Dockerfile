FROM golang:1.20-alpine AS backend-builder

WORKDIR /app/backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o media-server ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ffmpeg nginx sqlite3

WORKDIR /app
COPY --from=backend-builder /app/backend/media-server ./
COPY config/config.yaml ./config/
COPY start.sh ./

RUN chmod +x /app/start.sh

VOLUME ["/app/data", "/app/media", "/app/cache"]

EXPOSE 8080

CMD ["/app/start.sh"]
