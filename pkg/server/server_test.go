package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestExtractSearchTerm(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/", ""},
		{"/ ", ""},
		{"/  ", ""},
		{"/   ", ""},
		{"/hello", "hello"},
		{"/hello world", "hello world"},
		{"/言葉", "言葉"},
		{"hello", ""},
	}

	for _, test := range tests {
		if output := extractSearchTerm(test.input); output != test.expected {
			t.Errorf("Test Failed: %s inputted, %s expected, received: %s", test.input, test.expected, output)
		}
	}
}

func TestHandleValidRequest(t *testing.T) {
	tests := []struct {
		Path         string
		ExpectedWord string
	}{
		{"/hello", "hello"},
		{"/言葉", "言葉"},
		{"/hello%20world", "hello world"},
	}

	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, test := range tests {
		go func(path, expectedWord string) {
			defer wg.Done()

			handler := http.HandlerFunc(handleRequest)

			// Create a new test HTTP request
			req := httptest.NewRequest("GET", path, nil)

			// Create a response recorder to capture the server's response
			rec := httptest.NewRecorder()

			// Serve the HTTP request and record the response
			handler.ServeHTTP(rec, req)

			// Check the response status code
			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}

			// Decode the response body
			var response map[string]string
			err := json.NewDecoder(rec.Body).Decode(&response)
			if err != nil {
				t.Errorf("failed to decode response body: %v", err)
			}

			// Check the response content
			expected := map[string]string{
				"word": expectedWord,
			}
			if !compareMaps(response, expected) {
				t.Errorf("expected response %v, got %v", expected, response)
			}
		}(test.Path, test.ExpectedWord)
	}

	wg.Wait()
}

func TestHandleInvalidRequest(t *testing.T) {
	tests := []string{
		"/%20%20",
		"/%20%20%20/%20/",
		"/",
	}

	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, test := range tests {
		go func(path string) {
			defer wg.Done()

			handler := http.HandlerFunc(handleRequest)

			// Create a new test HTTP request
			req := httptest.NewRequest("GET", path, nil)

			// Create a response recorder to capture the server's response
			rec := httptest.NewRecorder()

			// Serve the HTTP request and record the response
			handler.ServeHTTP(rec, req)

			// Check the response status code
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rec.Code)
			}
		}(test)
	}

	wg.Wait()
}

// compareMaps compares two maps for equality
func compareMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valueA := range a {
		if valueB, ok := b[key]; !ok || valueA != valueB {
			return false
		}
	}
	return true
}
