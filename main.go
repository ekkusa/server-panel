package main

import (
	"log"
	"net/http"
	"os"
)

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
		dataPath = "/data"
	}

	dc, err := NewDockerClient(containerName)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer dc.Close()

	h := NewHandler(dc, dataPath)

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
