package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Handler holds all HTTP handler methods.
type Handler struct {
	dc *DockerClient
}

// NewHandler creates a Handler with the given DockerClient.
func NewHandler(dc *DockerClient) *Handler {
	return &Handler{dc: dc}
}

type apiResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// Status godoc
// GET /api/status — returns container state and resource usage.
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	status, err := h.dc.GetStatus(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, status)
}

// Start godoc
// POST /api/start — starts the container.
func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	if err := h.dc.Start(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server starting"})
}

// Stop godoc
// POST /api/stop — stops the container gracefully.
func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	if err := h.dc.Stop(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server stopping"})
}

// Restart godoc
// POST /api/restart — restarts the container.
func (h *Handler) Restart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	if err := h.dc.Restart(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server restarting"})
}

// StreamLogs godoc
// GET /api/logs — streams logs as Server-Sent Events (SSE).
// Query params:
//   - tail: number of past lines to include (default "100")
func (h *Handler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: "streaming not supported"})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering if behind proxy

	tail := r.URL.Query().Get("tail")

	// Wrap the ResponseWriter so log lines are formatted as SSE events.
	sseWriter := &sseLineWriter{w: w, flusher: flusher}

	// Send a heartbeat comment every 15s to keep the connection alive.
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		h.dc.StreamLogs(r.Context(), sseWriter, tail)
	}()

	for {
		select {
		case <-done:
			return
		case <-r.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		}
	}
}

// Command godoc
// POST /api/command — sends a command to the Minecraft server console.
// Body: { "command": "say Hello!" }
func (h *Handler) Command(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}

	var body struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON body"})
		return
	}
	if body.Command == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "command must not be empty"})
		return
	}

	if err := h.dc.SendCommand(r.Context(), body.Command); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "command sent"})
}

// sseLineWriter wraps http.ResponseWriter to emit each written line as an SSE data event.
type sseLineWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (s *sseLineWriter) Write(p []byte) (int, error) {
	fmt.Fprintf(s.w, "data: %s\n\n", p)
	s.flusher.Flush()
	return len(p), nil
}
