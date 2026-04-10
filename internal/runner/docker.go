package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// ContainerInfo holds the relevant runtime state of a Docker container.
type ContainerInfo struct {
	Name  string
	Image string
	ID    string
}

// DockerRunner fetches running container information from the local Docker daemon.
type DockerRunner struct{}

// NewDockerRunner creates a new DockerRunner.
func NewDockerRunner() *DockerRunner {
	return &DockerRunner{}
}

type dockerPsEntry struct {
	ID    string `json:"ID"`
	Image string `json:"Image"`
	Names string `json:"Names"`
}

// ListContainers returns the currently running containers.
func (d *DockerRunner) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	cmd := exec.CommandContext(ctx, "docker", "ps", "--format", "{{json .}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker ps failed: %w", err)
	}

	var containers []ContainerInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var entry dockerPsEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to parse docker ps output: %w", err)
		}
		containers = append(containers, ContainerInfo{
			ID:    entry.ID,
			Image: entry.Image,
			Name:  strings.TrimPrefix(entry.Names, "/"),
		})
	}
	return containers, nil
}
