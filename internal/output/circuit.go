package output

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// CircuitState represents the current state of the circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// DefaultCircuitOptions returns sensible defaults for the circuit breaker.
func DefaultCircuitOptions() CircuitOptions {
	return CircuitOptions{
		Enabled:      false,
		MaxFailures:  3,
		ResetTimeout: 30 * time.Second,
	}
}

// CircuitOptions configures the circuit breaker writer.
type CircuitOptions struct {
	Enabled      bool
	MaxFailures  int
	ResetTimeout time.Duration
}

// circuitWriter wraps a downstream Writer with circuit breaker logic.
type circuitWriter struct {
	opts      CircuitOptions
	next      Writer
	mu        sync.Mutex
	state     CircuitState
	failures  int
	lastOpen  time.Time
}

// NewCircuitWriter returns a Writer that trips open after MaxFailures consecutive
// downstream errors, and attempts recovery after ResetTimeout.
func NewCircuitWriter(opts CircuitOptions, next Writer) (Writer, error) {
	if next == nil {
		return nil, errors.New("circuit: next writer must not be nil")
	}
	return &circuitWriter{
		opts:  opts,
		next:  next,
		state: CircuitClosed,
	}, nil
}

func (c *circuitWriter) Write(results []drift.Result) error {
	if !c.opts.Enabled {
		return c.next.Write(results)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.state {
	case CircuitOpen:
		if time.Since(c.lastOpen) >= c.opts.ResetTimeout {
			c.state = CircuitHalfOpen
		} else {
			return fmt.Errorf("circuit: open — downstream unavailable, retry after %s", c.opts.ResetTimeout)
		}
	case CircuitHalfOpen:
		// allow one probe through
	}

	err := c.next.Write(results)
	if err != nil {
		c.failures++
		if c.failures >= c.opts.MaxFailures || c.state == CircuitHalfOpen {
			c.state = CircuitOpen
			c.lastOpen = time.Now()
		}
		return err
	}

	c.failures = 0
	c.state = CircuitClosed
	return nil
}
