package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type RequestConfig struct {
	URL            string
	Method         string
	Headers        map[string]string
	Body           any
	Timeout        time.Duration
	BasicAuth      *BasicAuthConfig
	ExpectedStatus int // If non-zero, validates response status code
}

type BasicAuthConfig struct {
	Username string
	Password string
}

func SendRequest(config RequestConfig) ([]byte, int, error) {
	// make request body
	var bodyReader io.Reader
	if config.Body != nil {
		switch body := config.Body.(type) {
		case []byte:
			bodyReader = bytes.NewBuffer(body)
		case string:
			bodyReader = bytes.NewBufferString(body)
		default:
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonBody)
		}
	}

	// create request
	req, err := http.NewRequest(config.Method, config.URL, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// set basic auth (base64 user:pass)
	if config.BasicAuth != nil {
		req.SetBasicAuth(config.BasicAuth.Username, config.BasicAuth.Password)
	}

	// set timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// validate status code if expected status is set
	if config.ExpectedStatus != 0 && resp.StatusCode != config.ExpectedStatus {
		return responseBody, resp.StatusCode, fmt.Errorf("unexpected status code: %d (expected %d)", resp.StatusCode, config.ExpectedStatus)
	}

	return responseBody, resp.StatusCode, nil
}

func GetIP(r *http.Request) string {
	// Check X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// This may contain multiple IPs; the first is the client
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON: %v", err)
		}
	}
}

func RespondWithError(w http.ResponseWriter, statusCode int, errors any) {
	RespondWithJSON(w, statusCode, errors)
}
