package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"strings"
	"os"
)

// Handler holds all HTTP handler methods.
type Handler struct {
	dc       *DockerClient
	dataPath string
}

// NewHandler creates a Handler.
func NewHandler(dc *DockerClient, dataPath string) *Handler {
	return &Handler{dc: dc, dataPath: dataPath}
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
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "server restarting"})
}

// GET /api/logs — SSE stream
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

	tail := r.URL.Query().Get("tail")
	sseWriter := &sseLineWriter{w: w, flusher: flusher}
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
	if err := h.dc.SendCommand(r.Context(), body.Command); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
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

// GET /api/files?path=/data/...  — list directory
func (h *Handler) FilesDir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	reqPath := r.URL.Query().Get("path")
	if reqPath == "" {
		reqPath = "/data"
	}
	// Map /data to the actual dataPath
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

// GET /api/files/content?path=/data/...  — read file
func (h *Handler) FileContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	reqPath := r.URL.Query().Get("path")
	// Map /data to the actual dataPath
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

// POST /api/files/write  — write file body: {path, content}
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
	// Map /data to the actual dataPath
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

// GET /api/mods  — list mods directory
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

// POST /api/mods/enable  — move mod from disabled back to mods folder
func (h *Handler) ModEnable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	// Map /data to actual dataPath
	modPath := body.Path
	if strings.HasPrefix(modPath, "/data/") {
		modPath = filepath.Join(h.dataPath, strings.TrimPrefix(modPath, "/data"))
	}
	fileName := filepath.Base(modPath)
	enabledPath := filepath.Join(h.dataPath, "mods", fileName)
	
	if err := os.Rename(modPath, enabledPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "mod enabled"})
}

// POST /api/mods/disable  — move mod to disabled folder
func (h *Handler) ModDisable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	// Map /data to actual dataPath
	modPath := body.Path
	if strings.HasPrefix(modPath, "/data/") {
		modPath = filepath.Join(h.dataPath, strings.TrimPrefix(modPath, "/data"))
	}
	if err := disableMod(modPath, h.dataPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "mod disabled"})
}

// POST /api/mods/remove  — delete mod
func (h *Handler) ModRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}
	var body struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON"})
		return
	}
	// Map /data to actual dataPath
	modPath := body.Path
	if strings.HasPrefix(modPath, "/data/") {
		modPath = filepath.Join(h.dataPath, strings.TrimPrefix(modPath, "/data"))
	}
	if err := os.Remove(modPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{OK: true, Message: "mod removed"})
}

// GET /api/config  — read docker-compose file from host
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

// POST /api/config  — write docker-compose file body: {content}
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

// sseLineWriter emits each write as an SSE data event.
type sseLineWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func (s *sseLineWriter) Write(p []byte) (int, error) {
	fmt.Fprintf(s.w, "data: %s\n\n", p)
	s.flusher.Flush()
	return len(p), nil
}
