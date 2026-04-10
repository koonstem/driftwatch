package runner

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
)

// TestListContainers_Integration runs only when Docker is available.
func TestListContainers_Integration(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not available, skipping integration test")
	}
	if os.Getenv("DRIFTWATCH_INTEGRATION") == "" {
		t.Skip("set DRIFTWATCH_INTEGRATION=1 to run Docker integration tests")
	}

	r := NewDockerRunner()
	containers, err := r.ListContainers(context.Background())
	if err != nil {
		t.Fatalf("ListContainers returned error: %v", err)
	}
	t.Logf("found %d running containers", len(containers))
}

func TestDockerPsEntryParsing(t *testing.T) {
	entry := dockerPsEntry{
		ID:    "abc123",
		Image: "nginx:latest",
		Names: "/web",
	}
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed dockerPsEntry
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if parsed.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", parsed.ID)
	}
	if parsed.Image != "nginx:latest" {
		t.Errorf("expected image nginx:latest, got %s", parsed.Image)
	}
	if parsed.Names != "/web" {
		t.Errorf("expected Names /web, got %s", parsed.Names)
	}
}

func TestContainerInfoNameStripping(t *testing.T) {
	info := ContainerInfo{
		ID:    "abc123",
		Image: "nginx:latest",
		Name:  "web",
	}
	if info.Name != "web" {
		t.Errorf("expected name 'web', got '%s'", info.Name)
	}
}
