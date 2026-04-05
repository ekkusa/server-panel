FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o panel .

# ── Runtime image ──────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/panel .
COPY static/ ./static/

EXPOSE 8080
ENV PORT=8080
ENV CONTAINER_NAME=minecraft

CMD ["./panel"]
