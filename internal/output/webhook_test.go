package output

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func makeWebhookReport(drifted bool) drift.Report {
	result := drift.Result{
		Service: "api",
		Drifted: drifted,
	}
	if drifted {
		result.Fields = []drift.FieldDiff{{Field: "image", Expected: "nginx:1.24", Actual: "nginx:1.23"}}
	}
	return drift.Report{Results: []drift.Result{result}}
}

func TestWebhookWriter_Disabled_NoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	opts := DefaultWebhookOptions()
	opts.Enabled = false
	opts.URL = ts.URL
	w := NewWebhookWriter(opts)

	if err := w.Write(makeWebhookReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when disabled")
	}
}

func TestWebhookWriter_OnlyDrift_NoDrift_Suppressed(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	opts := DefaultWebhookOptions()
	opts.Enabled = true
	opts.OnlyDrift = true
	opts.URL = ts.URL
	w := NewWebhookWriter(opts)

	if err := w.Write(makeWebhookReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when no drift and only_drift=true")
	}
}

func TestWebhookWriter_WithDrift_PostsPayload(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	opts := DefaultWebhookOptions()
	opts.Enabled = true
	opts.URL = ts.URL
	w := NewWebhookWriter(opts)

	if err := w.Write(makeWebhookReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Drifted != 1 {
		t.Errorf("expected drifted=1, got %d", received.Drifted)
	}
	if received.Total != 1 {
		t.Errorf("expected total=1, got %d", received.Total)
	}
}

func TestWebhookWriter_SecretHeader(t *testing.T) {
	var gotSecret string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSecret = r.Header.Get("X-Driftwatch-Secret")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	opts := DefaultWebhookOptions()
	opts.Enabled = true
	opts.OnlyDrift = false
	opts.URL = ts.URL
	opts.Secret = "s3cr3t"
	w := NewWebhookWriter(opts)

	if err := w.Write(makeWebhookReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotSecret != "s3cr3t" {
		t.Errorf("expected secret header 's3cr3t', got %q", gotSecret)
	}
}

func TestWebhookWriter_Non2xx_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	opts := DefaultWebhookOptions()
	opts.Enabled = true
	opts.OnlyDrift = false
	opts.URL = ts.URL
	w := NewWebhookWriter(opts)

	if err := w.Write(makeWebhookReport(false)); err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func TestWebhookWriter_NoURL_ReturnsError(t *testing.T) {
	opts := DefaultWebhookOptions()
	opts.Enabled = true
	opts.URL = ""
	w := NewWebhookWriter(opts)

	if err := w.Write(makeWebhookReport(true)); err == nil {
		t.Error("expected error when URL is empty")
	}
}
