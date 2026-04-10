package runner

import "context"

// Runner is the interface for fetching live service state from a container runtime.
type Runner interface {
	ListContainers(ctx context.Context) ([]ContainerInfo, error)
}

// RuntimeType identifies the supported container runtimes.
type RuntimeType string

const (
	RuntimeDocker RuntimeType = "docker"
)

// New returns a Runner for the given runtime type.
func New(rt RuntimeType) (Runner, error) {
	switch rt {
	case RuntimeDocker:
		return NewDockerRunner(), nil
	default:
		return nil, &UnsupportedRuntimeError{Runtime: string(rt)}
	}
}

// UnsupportedRuntimeError is returned when an unknown runtime is requested.
type UnsupportedRuntimeError struct {
	Runtime string
}

func (e *UnsupportedRuntimeError) Error() string {
	return "unsupported runtime: " + e.Runtime
}
