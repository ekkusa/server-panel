package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
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

// ListDir lists the immediate contents of dirPath inside the container
// using `ls -1p` via exec.
func (dc *DockerClient) ListDir(ctx context.Context, dirPath string) (*DirListing, error) {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return nil, err
	}

	exec, err := dc.cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh", "-c", "ls -1p '" + strings.ReplaceAll(dirPath, "'", "'\\''") + "' 2>&1"},
	})
	if err != nil {
		return nil, fmt.Errorf("exec create: %w", err)
	}

	resp, err := dc.cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, fmt.Errorf("exec attach: %w", err)
	}
	defer resp.Close()

	var buf bytes.Buffer
	readDockerStream(resp.Reader, &buf)

	listing := &DirListing{Path: dirPath, Entries: []FileEntry{}}
	for _, line := range strings.Split(buf.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		isDir := strings.HasSuffix(line, "/")
		listing.Entries = append(listing.Entries, FileEntry{
			Name:  strings.TrimSuffix(line, "/"),
			IsDir: isDir,
		})
	}
	return listing, nil
}

// ReadFile reads a single file from inside the container via CopyFromContainer.
func (dc *DockerClient) ReadFile(ctx context.Context, filePath string) (*FileContent, error) {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return nil, err
	}

	if !isTextFile(path.Base(filePath)) {
		return &FileContent{Path: filePath, Binary: true}, nil
	}

	reader, _, err := dc.cli.CopyFromContainer(ctx, id, filePath)
	if err != nil {
		return nil, fmt.Errorf("copy from container: %w", err)
	}
	defer reader.Close()

	tr := tar.NewReader(reader)
	if _, err := tr.Next(); err != nil {
		return nil, fmt.Errorf("reading tar header: %w", err)
	}
	data, err := io.ReadAll(tr)
	if err != nil {
		return nil, fmt.Errorf("reading content: %w", err)
	}
	return &FileContent{Path: filePath, Content: string(data)}, nil
}

// WriteFile writes content to a file inside the container via CopyToContainer.
func (dc *DockerClient) WriteFile(ctx context.Context, filePath string, content []byte) error {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := tw.WriteHeader(&tar.Header{
		Name: path.Base(filePath),
		Size: int64(len(content)),
		Mode: 0644,
	}); err != nil {
		return err
	}
	if _, err := tw.Write(content); err != nil {
		return err
	}
	tw.Close()

	return dc.cli.CopyToContainer(ctx, id, path.Dir(filePath), &buf, types.CopyToContainerOptions{})
}

// safePath ensures the requested path is under dataPath to prevent traversal.
func safePath(dataPath, requested string) (string, error) {
	// If the request is already an absolute path under dataPath, use it.
	// Otherwise join it with dataPath.
	var full string
	if strings.HasPrefix(requested, "/") {
		full = path.Clean(requested)
	} else {
		full = path.Clean(path.Join(dataPath, requested))
	}
	if !strings.HasPrefix(full+"/", path.Clean(dataPath)+"/") {
		return "", fmt.Errorf("path outside data directory")
	}
	return full, nil
}
