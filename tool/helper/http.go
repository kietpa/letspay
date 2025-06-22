package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// Prepare request body
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

	// Create request
	req, err := http.NewRequest(config.Method, config.URL, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Set basic auth (base64 user:pass)
	if config.BasicAuth != nil {
		req.SetBasicAuth(config.BasicAuth.Username, config.BasicAuth.Password)
	}

	// Set timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// Validate status code if expected status is set
	if config.ExpectedStatus != 0 && resp.StatusCode != config.ExpectedStatus {
		return responseBody, resp.StatusCode, fmt.Errorf("unexpected status code: %d (expected %d)", resp.StatusCode, config.ExpectedStatus)
	}

	return responseBody, resp.StatusCode, nil
}

// // SendJSONRequest sends a JSON request and unmarshals the response
// func SendJSONRequest(config RequestConfig, response any) (int, error) {
// 	// Ensure Content-Type is set for JSON
// 	if config.Headers == nil {
// 		config.Headers = make(map[string]string)
// 	}
// 	if _, exists := config.Headers["Content-Type"]; !exists {
// 		config.Headers["Content-Type"] = "application/json"
// 	}

// 	body, statusCode, err := SendRequest(config)
// 	if err != nil {
// 		return statusCode, err
// 	}

// 	if response != nil {
// 		if err := json.Unmarshal(body, response); err != nil {
// 			return statusCode, fmt.Errorf("failed to unmarshal response: %w", err)
// 		}
// 	}

// 	return statusCode, nil
// }
