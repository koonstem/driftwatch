package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/driftwatch/internal/drift"
)

// WebhookOptions configures outbound webhook delivery.
type WebhookOptions struct {
	Enabled  bool
	URL      string
	Secret   string // sent as X-Driftwatch-Secret header
	OnlyDrift bool  // only POST when drift is detected
	Timeout  time.Duration
}

// DefaultWebhookOptions returns safe defaults.
func DefaultWebhookOptions() WebhookOptions {
	return WebhookOptions{
		Enabled:   false,
		OnlyDrift: true,
		Timeout:   10 * time.Second,
	}
}

// webhookPayload is the JSON body sent to the webhook endpoint.
type webhookPayload struct {
	Timestamp string            `json:"timestamp"`
	Drifted   int               `json:"drifted"`
	Total     int               `json:"total"`
	Results   []drift.Result    `json:"results"`
}

// WebhookWriter posts drift results to an HTTP endpoint.
type WebhookWriter struct {
	opts   WebhookOptions
	client *http.Client
}

// NewWebhookWriter creates a WebhookWriter with the given options.
func NewWebhookWriter(opts WebhookOptions) *WebhookWriter {
	return &WebhookWriter{
		opts: opts,
		client: &http.Client{Timeout: opts.Timeout},
	}
}

// Write sends the report to the configured webhook URL if enabled.
func (w *WebhookWriter) Write(report drift.Report) error {
	if !w.opts.Enabled {
		return nil
	}
	if w.opts.URL == "" {
		return fmt.Errorf("webhook enabled but no URL configured")
	}
	if w.opts.OnlyDrift && !HasDrift(report) {
		return nil
	}

	drifted := 0
	for _, r := range report.Results {
		if r.Drifted {
			drifted++
		}
	}

	payload := webhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Drifted:   drifted,
		Total:     len(report.Results),
		Results:   report.Results,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.opts.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if w.opts.Secret != "" {
		req.Header.Set("X-Driftwatch-Secret", w.opts.Secret)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
