# Minecraft Panel

A lightweight web panel for managing your Docker Minecraft server.  
Go backend + terminal-themed HTML frontend.

## Features

- **Start / Stop / Restart** your Minecraft container with one click
- **Live log streaming** via Server-Sent Events (no WebSocket dependency)
- **Resource stats**: CPU % and memory usage with live progress bars
- **Console commands** via `rcon-cli` (works with `itzg/minecraft-server`)
- **5-second auto-refresh** of container status

---

## Prerequisites

| Tool | Version |
|------|---------|
| Go   | 1.22+   |
| Docker | any recent version |
| `itzg/minecraft-server` image | (or any image with `rcon-cli` for commands) |

---

## Quick Start

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Run

```bash
# Defaults: PORT=8080, CONTAINER_NAME=minecraft
go run .

# Custom container name
CONTAINER_NAME=my-mc-server go run .

# Custom port
PORT=9090 CONTAINER_NAME=mc go run .
```

### 3. Open

Visit **http://localhost:8080** in your browser.

---

## API Reference

| Method | Path | Description |
|--------|------|-------------|
| GET  | `/api/status`  | Container state + resource usage |
| POST | `/api/start`   | Start the container |
| POST | `/api/stop`    | Graceful stop (30s timeout) |
| POST | `/api/restart` | Restart the container |
| GET  | `/api/logs`    | SSE stream of container logs (`?tail=100`) |
| POST | `/api/command` | Send command: `{"command": "say hi"}` |

---

## Docker Compose Example

If you want to run the panel itself in Docker alongside your Minecraft server:

```yaml
services:
  minecraft:
    image: itzg/minecraft-server
    container_name: minecraft
    environment:
      EULA: "TRUE"
      ENABLE_RCON: "true"
      RCON_PASSWORD: "changeme"
    ports:
      - "25565:25565"
    volumes:
      - ./data:/data

  panel:
    build: .
    ports:
      - "8080:8080"
    environment:
      CONTAINER_NAME: minecraft
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  # Required!
    depends_on:
      - minecraft
```

> **Important**: The panel needs access to the Docker socket. If running outside Docker, it uses your local socket automatically.

---

## Notes on Commands

- Commands are sent via `rcon-cli` inside the container — this requires `ENABLE_RCON=true` in your `itzg/minecraft-server` setup.
- If you use a different image without `rcon-cli`, you can modify `SendCommand` in `docker.go` to write to the container's stdin instead.

---

## Project Structure

```
minecraft-panel/
├── main.go       — HTTP server + middleware
├── docker.go     — Docker SDK wrapper (status, start/stop, logs, commands)
├── handlers.go   — HTTP handlers for all API endpoints
├── go.mod
└── static/
    └── index.html — Frontend panel (self-contained, no build step)
```
