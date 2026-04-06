package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// DockerClient wraps the Docker SDK client for Minecraft container management.
type DockerClient struct {
	cli           *client.Client
	containerName string
}

// ServerStatus holds the current state of the Minecraft container.
type ServerStatus struct {
	Running     bool      `json:"running"`
	Status      string    `json:"status"`
	ContainerID string    `json:"container_id"`
	Image       string    `json:"image"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	Uptime      string    `json:"uptime,omitempty"`
	CPUPercent  float64   `json:"cpu_percent"`
	MemUsageMB  float64   `json:"mem_usage_mb"`
	MemLimitMB  float64   `json:"mem_limit_mb"`
}

// NewDockerClient creates a new DockerClient connected to the local Docker daemon.
func NewDockerClient(containerName string) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("connecting to Docker daemon: %w", err)
	}
	return &DockerClient{cli: cli, containerName: containerName}, nil
}

// Close releases the Docker client connection.
func (dc *DockerClient) Close() {
	dc.cli.Close()
}

// findContainer returns the container ID by name.
func (dc *DockerClient) findContainer(ctx context.Context) (string, error) {
	f := filters.NewArgs(filters.Arg("name", dc.containerName))
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return "", fmt.Errorf("listing containers: %w", err)
	}
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+dc.containerName || name == dc.containerName {
				return c.ID, nil
			}
		}
	}
	return "", fmt.Errorf("container %q not found", dc.containerName)
}

// GetStatus returns the current status of the Minecraft container.
func (dc *DockerClient) GetStatus(ctx context.Context) (*ServerStatus, error) {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return &ServerStatus{Running: false, Status: "not found"}, nil
	}

	info, err := dc.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("inspecting container: %w", err)
	}

	status := &ServerStatus{
		Running:     info.State.Running,
		Status:      info.State.Status,
		ContainerID: id[:12],
		Image:       info.Config.Image,
	}

	if info.State.Running {
		startedAt, err := time.Parse(time.RFC3339Nano, info.State.StartedAt)
		if err == nil {
			status.StartedAt = startedAt
			status.Uptime = formatUptime(time.Since(startedAt))
		}

		statsResp, err := dc.cli.ContainerStats(ctx, id, false)
		if err == nil {
			defer statsResp.Body.Close()
			var statsJSON types.StatsJSON
			if err := json.NewDecoder(statsResp.Body).Decode(&statsJSON); err == nil {
				status.CPUPercent = calcCPUPercent(&statsJSON)
				status.MemUsageMB = float64(statsJSON.MemoryStats.Usage) / 1024 / 1024
				status.MemLimitMB = float64(statsJSON.MemoryStats.Limit) / 1024 / 1024
			}
		}
	}

	return status, nil
}

// Start starts the Minecraft container.
func (dc *DockerClient) Start(ctx context.Context) error {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return err
	}
	return dc.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// Stop gracefully stops the Minecraft container with a 30-second timeout.
func (dc *DockerClient) Stop(ctx context.Context) error {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return err
	}
	timeout := 30
	return dc.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

// Restart restarts the Minecraft container.
func (dc *DockerClient) Restart(ctx context.Context) error {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return err
	}
	timeout := 30
	return dc.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

// StreamLogs streams container logs to the provided writer.
// Docker's multiplexed log format prefixes each frame with an 8-byte header.
func (dc *DockerClient) StreamLogs(ctx context.Context, w io.Writer, tail string) error {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return err
	}
	if tail == "" {
		tail = "100"
	}

	reader, err := dc.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
		Tail:       tail,
	})
	if err != nil {
		return fmt.Errorf("fetching logs: %w", err)
	}
	defer reader.Close()

	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		header := make([]byte, 8)
		if _, err := io.ReadFull(reader, header); err != nil {
			return nil
		}
		size := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
		if size == 0 {
			continue
		}
		if size > len(buf) {
			buf = make([]byte, size)
		}
		if _, err := io.ReadFull(reader, buf[:size]); err != nil {
			return nil
		}
		line := strings.TrimRight(string(buf[:size]), "\n")
		fmt.Fprintf(w, "%s\n", line)
	}
}

// SendCommand sends a command to the Minecraft server via rcon-cli and captures output
func (dc *DockerClient) SendCommand(ctx context.Context, command string) (string, error) {
	id, err := dc.findContainer(ctx)
	if err != nil {
		return "", err
	}

	exec, err := dc.cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"rcon-cli", command},
	})
	if err != nil {
		return "", fmt.Errorf("creating exec: %w", err)
	}

	resp, err := dc.cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("attaching to exec: %w", err)
	}
	defer resp.Close()

	var output strings.Builder
	buf := make([]byte, 4096)
	for {
		// Read 8-byte header
		header := make([]byte, 8)
		_, err := io.ReadFull(resp.Conn, header)
		if err != nil {
			break
		}
		
		// Parse frame size from header
		size := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
		if size == 0 {
			continue
		}
		
		if size > len(buf) {
			buf = make([]byte, size)
		}
		
		n, err := io.ReadFull(resp.Conn, buf[:size])
		if n > 0 {
			output.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	return strings.TrimSpace(output.String()), nil
}

func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func calcCPUPercent(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	numCPU := float64(stats.CPUStats.OnlineCPUs)
	if numCPU == 0 {
		numCPU = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0 && cpuDelta > 0 {
		return (cpuDelta / systemDelta) * numCPU * 100.0
	}
	return 0.0
}
