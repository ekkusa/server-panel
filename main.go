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

	dc, err := NewDockerClient(containerName)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer dc.Close()

	h := NewHandler(dc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", h.Status)
	mux.HandleFunc("/api/start", h.Start)
	mux.HandleFunc("/api/stop", h.Stop)
	mux.HandleFunc("/api/restart", h.Restart)
	mux.HandleFunc("/api/logs", h.StreamLogs)
	mux.HandleFunc("/api/command", h.Command)
	mux.HandleFunc("/", serveFrontend) // HTML baked into binary — no static/ dir needed

	log.Printf("🟢 Minecraft Panel running on http://localhost:%s", port)
	log.Printf("📦 Managing container: %s", containerName)
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
