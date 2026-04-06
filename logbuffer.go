package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const bufferMaxLines = 2000

// LogLine is a single console entry.
type LogLine struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
	Kind string    `json:"kind"` // "docker" | "rcon" | "panel"
}

// LogBuffer is a thread-safe ring buffer that broadcasts new lines to all
// active SSE subscribers and optionally persists lines to a log file.
type LogBuffer struct {
	mu      sync.RWMutex
	buf     []LogLine
	subs    map[chan LogLine]struct{}
	logFile *os.File // nil if persistence is disabled
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{
		subs: make(map[chan LogLine]struct{}),
	}
}

// LoadFromFile opens (or creates) the given file, reads existing lines into
// the buffer, and keeps the file open for future appends.
func (lb *LogBuffer) LoadFromFile(path string) error {
	// Open for reading first to load history.
	rf, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}

	scanner := bufio.NewScanner(rf)
	var lines []LogLine
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// Format: "kind\ttext"
		idx := strings.IndexByte(line, '\t')
		if idx < 0 {
			continue
		}
		lines = append(lines, LogLine{
			Text: line[idx+1:],
			Kind: line[:idx],
			Time: time.Now(),
		})
	}
	rf.Close()

	// Trim to max lines.
	if len(lines) > bufferMaxLines {
		lines = lines[len(lines)-bufferMaxLines:]
	}

	lb.mu.Lock()
	lb.buf = lines
	lb.mu.Unlock()

	// If the file grew too large, rewrite it with only the kept lines.
	if len(lines) < bufferMaxLines {
		// No truncation needed yet.
	} else {
		lb.rewriteFile(path, lines)
	}

	// Re-open for appending.
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("opening log file for append: %w", err)
	}
	lb.logFile = wf
	return nil
}

// rewriteFile rewrites the log file from scratch with only the given lines.
func (lb *LogBuffer) rewriteFile(path string, lines []LogLine) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range lines {
		fmt.Fprintf(w, "%s\t%s\n", l.Kind, l.Text)
	}
	w.Flush()
}

// Push adds a line to the buffer, persists it, and fans it out to subscribers.
// Duplicate lines arriving within 2 seconds are dropped to prevent
// double-delivery when rcon output also appears in Docker logs.
func (lb *LogBuffer) Push(text, kind string) {
	now := time.Now()
	lb.mu.Lock()

	// Dedup: skip if identical text was pushed in the last 2 seconds.
	for i := len(lb.buf) - 1; i >= 0 && now.Sub(lb.buf[i].Time) < 2*time.Second; i-- {
		if lb.buf[i].Text == text {
			lb.mu.Unlock()
			return
		}
	}

	line := LogLine{Text: text, Time: now, Kind: kind}
	lb.buf = append(lb.buf, line)
	if len(lb.buf) > bufferMaxLines {
		lb.buf = lb.buf[len(lb.buf)-bufferMaxLines:]
	}

	// Persist to file.
	if lb.logFile != nil {
		fmt.Fprintf(lb.logFile, "%s\t%s\n", kind, text)
	}

	// Snapshot subscribers while holding lock.
	subs := make([]chan LogLine, 0, len(lb.subs))
	for ch := range lb.subs {
		subs = append(subs, ch)
	}
	lb.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- line:
		default:
		}
	}
}

// Snapshot returns a copy of all buffered lines.
func (lb *LogBuffer) Snapshot() []LogLine {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	out := make([]LogLine, len(lb.buf))
	copy(out, lb.buf)
	return out
}

// Subscribe returns a channel that receives future lines.
func (lb *LogBuffer) Subscribe() chan LogLine {
	ch := make(chan LogLine, 256)
	lb.mu.Lock()
	lb.subs[ch] = struct{}{}
	lb.mu.Unlock()
	return ch
}

// Unsubscribe removes a channel from the subscriber set.
func (lb *LogBuffer) Unsubscribe(ch chan LogLine) {
	lb.mu.Lock()
	delete(lb.subs, ch)
	lb.mu.Unlock()
	for len(ch) > 0 {
		<-ch
	}
}

// Close flushes and closes the log file.
func (lb *LogBuffer) Close() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if lb.logFile != nil {
		lb.logFile.Close()
		lb.logFile = nil
	}
}
