package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// rconNoisePatterns are Docker log lines we suppress globally.
// These are RCON connection handshake messages, not actual command output.
var rconNoisePatterns = []string{
	"RCON Listener",
	"Thread RCON Client",
	"RCON Client #",
}

func isRconNoise(line string) bool {
	for _, p := range rconNoisePatterns {
		if strings.Contains(line, p) {
			return true
		}
	}
	return false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	containerName := os.Getenv("CONTAINER_NAME")
	if containerName == "" {
		containerName = "minecraft"
	}
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "/home/exer/Downloads/files/data"
	}
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "panel-console.log"
	}

	dc, err := NewDockerClient(containerName)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer dc.Close()

	lb := NewLogBuffer()
	defer lb.Close()

	// Load persisted log history from disk.
	if err := lb.LoadFromFile(logFile); err != nil {
		log.Printf("Warning: could not load log file %q: %v", logFile, err)
	} else {
		log.Printf("Loaded console history from %s", logFile)
	}

	// Background: pump Docker logs into the buffer.
	// First connect tails 500 lines; reconnects use tail=0 to avoid duplicates.
	go func() {
		firstConnect := true
		for {
			tail := "0"
			if firstConnect {
				tail = "500"
			}
			err := dc.StreamLogs(context.Background(), &bufferWriter{lb: lb}, tail)
			if err != nil {
				log.Printf("docker log stream error: %v — retrying in 5s", err)
			}
			firstConnect = false
			time.Sleep(5 * time.Second)
		}
	}()

	h := NewHandler(dc, dataPath, lb)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", h.Status)
	mux.HandleFunc("/api/start", h.Start)
	mux.HandleFunc("/api/stop", h.Stop)
	mux.HandleFunc("/api/restart", h.Restart)
	mux.HandleFunc("/api/logs", h.StreamLogs)
	mux.HandleFunc("/api/command", h.Command)
	mux.HandleFunc("/api/players", h.Players)
	mux.HandleFunc("/api/files", h.FilesDir)
	mux.HandleFunc("/api/files/content", h.FileContent)
	mux.HandleFunc("/api/files/write", h.FileWrite)
	mux.HandleFunc("/api/mods/enable", h.ModEnable)
	mux.HandleFunc("/api/mods/disable", h.ModDisable)
	mux.HandleFunc("/api/mods/remove", h.ModRemove)
	mux.HandleFunc("/api/mods", h.Mods)
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ConfigGet(w, r)
		} else {
			h.ConfigSet(w, r)
		}
	})
	mux.HandleFunc("/", serveFrontend)

	log.Printf("Minecraft Panel on http://localhost:%s  container=%s  data=%s", port, containerName, dataPath)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}

// bufferWriter splits incoming bytes on newlines and pushes each line into
// the LogBuffer, filtering RCON noise globally.
type bufferWriter struct {
	lb      *LogBuffer
	partial strings.Builder
}

func (bw *bufferWriter) Write(p []byte) (int, error) {
	bw.partial.Write(p)
	for {
		s := bw.partial.String()
		idx := strings.IndexByte(s, '\n')
		if idx < 0 {
			break
		}
		line := strings.TrimRight(s[:idx], "\r")
		bw.partial.Reset()
		bw.partial.WriteString(s[idx+1:])
		if line != "" && !isRconNoise(line) {
			bw.lb.Push(line, "docker")
		}
	}
	return len(p), nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
