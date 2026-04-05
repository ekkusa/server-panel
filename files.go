package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os"
)

// FileEntry is one item in a directory listing.
type FileEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

// DirListing is the response for a directory listing request.
type DirListing struct {
	Path    string      `json:"path"`
	Entries []FileEntry `json:"entries"`
}

// FileContent is the response for a file read request.
type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Binary  bool   `json:"binary"`
}

// textExtensions are file types we'll read and display as text.
var textExtensions = map[string]bool{
	".properties": true, ".json": true, ".yml": true, ".yaml": true,
	".txt": true, ".conf": true, ".cfg": true, ".toml": true,
	".ini": true, ".xml": true, ".sh": true, ".md": true,
	".env": true, ".log": true, ".csv": true,
}

func isTextFile(name string) bool {
	idx := strings.LastIndex(name, ".")
	if idx < 0 {
		return false
	}
	return textExtensions[strings.ToLower(name[idx:])]
}

// ListDir lists the immediate contents of dirPath from the host filesystem.
func ListDir(dirPath string) (*DirListing, error) {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	listing := &DirListing{Path: dirPath, Entries: []FileEntry{}}
	for _, entry := range entries {
		listing.Entries = append(listing.Entries, FileEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  entry.Size(),
		})
	}
	return listing, nil
}

// ReadFile reads a single file from the host filesystem.
func ReadFile(filePath string) (*FileContent, error) {
	if !isTextFile(filepath.Base(filePath)) {
		return &FileContent{Path: filePath, Binary: true}, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return &FileContent{Path: filePath, Content: string(data)}, nil
}

// WriteFile writes content to a file on the host filesystem.
func WriteFile(filePath string, content []byte) error {
	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// safePath ensures the requested path is under dataPath to prevent traversal.
func safePath(dataPath, requested string) (string, error) {
	// Convert to absolute paths for comparison
	absData, err := filepath.Abs(dataPath)
	if err != nil {
		return "", fmt.Errorf("invalid data path: %w", err)
	}

	var full string
	if filepath.IsAbs(requested) {
		full = filepath.Clean(requested)
	} else {
		full = filepath.Clean(filepath.Join(absData, requested))
	}

	// Ensure full path is under dataPath
	if !strings.HasPrefix(full+"/", absData+"/") && full != absData {
		return "", fmt.Errorf("path outside data directory")
	}
	return full, nil
}

// disableMod moves a mod file to a disabled folder
func disableMod(modPath, dataPath string) error {
	disabledDir := filepath.Join(dataPath, "mods", "disabled")
	if err := os.MkdirAll(disabledDir, 0755); err != nil {
		return fmt.Errorf("creating disabled directory: %w", err)
	}
	
	fileName := filepath.Base(modPath)
	disabledPath := filepath.Join(disabledDir, fileName)
	
	if err := os.Rename(modPath, disabledPath); err != nil {
		return fmt.Errorf("moving mod to disabled: %w", err)
	}
	return nil
}
