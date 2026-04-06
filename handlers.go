package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Handler struct {
	dc       *DockerClient
	dataPath string
	lb       *LogBuffer
}

func NewHandler(dc *DockerClient, dataPath string, lb *LogBuffer) *Handler {
	return &Handler{dc: dc, dataPath: dataPath, lb: lb}
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

// GET /api/status
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

// POST /api/start
func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	if err := h.dc.Start(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	h.lb.Push("[panel] Server starting...", "panel")
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server starting"})
}

// POST /api/stop
func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	if err := h.dc.Stop(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	h.lb.Push("[panel] Server stopping...", "panel")
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server stopping"})
}

// POST /api/restart
func (h *Handler) Restart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	if err := h.dc.Restart(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	h.lb.Push("[panel] Server restarting...", "panel")
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server restarting"})
}

// GET /api/logs — SSE stream.
// On connect: replay the full buffer, then stream new lines in real-time.
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
	w.Header().Set("X-Accel-Buffering", "no")

	// Subscribe BEFORE snapshotting so we don't miss lines that arrive
	// between the snapshot and the subscribe.
	ch := h.lb.Subscribe()
	defer h.lb.Unsubscribe(ch)

	// Replay history.
	snapshot := h.lb.Snapshot()
	for _, line := range snapshot {
		fmt.Fprintf(w, "data: %s\n\n", line.Text)
	}
	flusher.Flush()

	// Stream live updates.
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case line, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", line.Text)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		}
	}
}

// POST /api/command
func (h *Handler) Command(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	if body.Command == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "command must not be empty"})
		return
	}

	output, err := h.dc.SendCommand(r.Context(), body.Command)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	// Push rcon output into the buffer so it persists across refreshes.
	// The buffer deduplicates within 2s, so if the same line also appears
	// in Docker logs it will only be shown once.
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			h.lb.Push(line, "rcon")
		}
	}

	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "command sent"})
}

// GET /api/players
func (h *Handler) Players(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	pl, err := h.dc.GetPlayers(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, pl)
}

// GET /api/files
func (h *Handler) FilesDir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	reqPath := r.URL.Query().Get("path")
	if reqPath == "" {
		reqPath = "/data"
	}
	if reqPath == "/data" {
		reqPath = h.dataPath
	} else if strings.HasPrefix(reqPath, "/data/") {
		reqPath = filepath.Join(h.dataPath, strings.TrimPrefix(reqPath, "/data"))
	}
	safe, err := safePath(h.dataPath, reqPath)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: err.Error()})
		return
	}
	listing, err := ListDir(safe)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, listing)
}

// GET /api/files/content
func (h *Handler) FileContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	reqPath := r.URL.Query().Get("path")
	if reqPath == "/data" {
		reqPath = h.dataPath
	} else if strings.HasPrefix(reqPath, "/data/") {
		reqPath = filepath.Join(h.dataPath, strings.TrimPrefix(reqPath, "/data"))
	}
	safe, err := safePath(h.dataPath, reqPath)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: err.Error()})
		return
	}
	fc, err := ReadFile(safe)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, fc)
}

// POST /api/files/write
func (h *Handler) FileWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	reqPath := body.Path
	if reqPath == "/data" {
		reqPath = h.dataPath
	} else if strings.HasPrefix(reqPath, "/data/") {
		reqPath = filepath.Join(h.dataPath, strings.TrimPrefix(reqPath, "/data"))
	}
	safe, err := safePath(h.dataPath, reqPath)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: err.Error()})
		return
	}
	if err := WriteFile(safe, []byte(body.Content)); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "file saved"})
}

// GET /api/mods
func (h *Handler) Mods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	modsPath := filepath.Join(h.dataPath, "mods")
	listing, err := ListDir(modsPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, listing)
}

// POST /api/mods/enable
func (h *Handler) ModEnable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct{ Path string `json:"path"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	modPath := mapDataPath(body.Path, h.dataPath)
	enabledPath := filepath.Join(h.dataPath, "mods", filepath.Base(modPath))
	if err := os.Rename(modPath, enabledPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "mod enabled"})
}

// POST /api/mods/disable
func (h *Handler) ModDisable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct{ Path string `json:"path"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	modPath := mapDataPath(body.Path, h.dataPath)
	if err := disableMod(modPath, h.dataPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "mod disabled"})
}

// POST /api/mods/remove
func (h *Handler) ModRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct{ Path string `json:"path"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	modPath := mapDataPath(body.Path, h.dataPath)
	if err := os.Remove(modPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "mod removed"})
}

// GET /api/config
func (h *Handler) ConfigGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	filename, content, err := ReadComposeFile()
	if err != nil {
		writeJSON(w, http.StatusNotFound, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"filename": filename, "content": content})
}

// POST /api/config
func (h *Handler) ConfigSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	filename, err := WriteComposeFile(body.Content)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "saved " + filename})
}

// mapDataPath translates frontend /data/... paths to real host paths.
func mapDataPath(p, dataPath string) string {
	if strings.HasPrefix(p, "/data/") {
		return filepath.Join(dataPath, strings.TrimPrefix(p, "/data"))
	}
	return p
}

// sseLineWriter kept for compatibility with docker.go StreamLogs signature.
type sseLineWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (s *sseLineWriter) Write(p []byte) (int, error) {
	fmt.Fprintf(s.w, "data: %s\n\n", p)
	s.flusher.Flush()
	return len(p), nil
}
