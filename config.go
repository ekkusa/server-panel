package main

import (
	"fmt"
	"os"
)

var composeFilenames = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

// ReadComposeFile finds and reads the docker-compose file in the current directory.
func ReadComposeFile() (filename string, content string, err error) {
	for _, name := range composeFilenames {
		data, e := os.ReadFile(name)
		if e == nil {
			return name, string(data), nil
		}
	}
	return "", "", fmt.Errorf("no docker-compose file found in current directory")
}

// WriteComposeFile writes content back to whichever compose file was found,
// or creates docker-compose.yml if none exists.
func WriteComposeFile(content string) (string, error) {
	for _, name := range composeFilenames {
		if _, err := os.Stat(name); err == nil {
			return name, os.WriteFile(name, []byte(content), 0644)
		}
	}
	return "docker-compose.yml", os.WriteFile("docker-compose.yml", []byte(content), 0644)
}
