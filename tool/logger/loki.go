package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LokiWriter struct {
	url    string
	labels map[string]string
	client *http.Client
}

func NewLokiWriter(url string, labels map[string]string) (io.Writer, error) {
	return &LokiWriter{
		url:    url,
		labels: labels,
		client: &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (l *LokiWriter) Write(p []byte) (int, error) {
	// Sanitize the log line (remove control chars, trim whitespace)
	logLine := string(p)

	// Get current time in nanoseconds since epoch (as string)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	// Build the Loki push payload
	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": l.labels,
				"values": [][]string{
					{timestamp, logLine},
				},
			},
		},
	}

	// Encode to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Send to Loki
	req, err := http.NewRequest("POST", l.url, bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("failed to create Loki request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to push log to Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("Loki push failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return len(p), nil
}
