package main

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
)

// PlayerList holds the result of a /list query.
type PlayerList struct {
	Online []string `json:"online"`
	Count  int      `json:"count"`
	Max    int      `json:"max"`
}

// rconListRegex matches: "There are 2 of a max of 20 players online: Player1, Player2"
var rconListRegex = regexp.MustCompile(`There are (\d+) of a max(?: of)? (\d+) players online:(.*)`)

// GetPlayers runs `rcon-cli list` inside the container and parses the response.
// Uses the container name directly — no findContainer needed.
func (dc *DockerClient) GetPlayers(ctx context.Context) (*PlayerList, error) {
	exec, err := dc.cli.ContainerExecCreate(ctx, dc.containerName, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"rcon-cli", "list"},
	})
	if err != nil {
		// Container not running or not found — return empty list silently.
		return &PlayerList{Online: []string{}}, nil
	}

	resp, err := dc.cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, fmt.Errorf("attaching exec: %w", err)
	}
	defer resp.Close()

	var buf bytes.Buffer
	readDockerStream(resp.Reader, &buf)

	output := strings.TrimSpace(buf.String())
	return parsePlayerList(output), nil
}

// parsePlayerList extracts player count and names from the rcon-cli list output.
func parsePlayerList(output string) *PlayerList {
	pl := &PlayerList{Online: []string{}}

	m := rconListRegex.FindStringSubmatch(output)
	if m == nil {
		return pl
	}

	count, _ := strconv.Atoi(m[1])
	max, _ := strconv.Atoi(m[2])
	pl.Count = count
	pl.Max = max

	namesPart := strings.TrimSpace(m[3])
	if namesPart != "" {
		for _, name := range strings.Split(namesPart, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				pl.Online = append(pl.Online, name)
			}
		}
	}

	return pl
}
