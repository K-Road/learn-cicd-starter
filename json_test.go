package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		message    string
		wantBody   string
		wantHeader string
	}{
		{
			name:       "Client Error 400",
			code:       400,
			message:    "Bad Request",
			wantBody:   `{"error":"Bad Request"}`,
			wantHeader: "application/json",
		},
		{
			name:       "Server Error 500",
			code:       500,
			message:    "Internal Server Error",
			wantBody:   `{"error":"Internal Server Error"}`,
			wantHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture logs
			var logOutput strings.Builder
			log.SetOutput(&logOutput)

			// Set up response recorder
			recorder := httptest.NewRecorder()

			// Call the function
			respondWithError(recorder, tt.code, tt.message)

			// Check status code
			if status := recorder.Code; status != tt.code {
				t.Errorf("expected status code %d, got %d", tt.code, status)
			}

			// Check response body
			gotBody := strings.TrimSpace(recorder.Body.String())
			if gotBody != tt.wantBody {
				t.Errorf("expected body %q, got %q", tt.wantBody, gotBody)
			}

			// Check Content-Type header
			if contentType := recorder.Header().Get("Content-Type"); contentType != tt.wantHeader {
				t.Errorf("expected Content-Type header %q, got %q", tt.wantHeader, contentType)
			}

			// Check for log output on 5XX errors
			if tt.code >= 500 && !strings.Contains(logOutput.String(), "5XX error") {
				t.Errorf("expected log output for 5XX error, but got none")
			}
		})
	}
}

func TestRespondWithJSON(t *testing.T) {
	// Prepare test data and expected result
	payload := map[string]string{"message": "Success"}
	expectedBody, _ := json.Marshal(payload)
	expectedHeader := "application/json"

	// Set up response recorder
	recorder := httptest.NewRecorder()

	// Call the function
	respondWithJSON(recorder, http.StatusOK, payload)

	// Check status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	// Check response body
	gotBody := strings.TrimSpace(recorder.Body.String())
	if gotBody != string(expectedBody) {
		t.Errorf("expected body %q, got %q", string(expectedBody), gotBody)
	}

	// Check Content-Type header
	if contentType := recorder.Header().Get("Content-Type"); contentType != expectedHeader {
		t.Errorf("expected Content-Type header %q, got %q", expectedHeader, contentType)
	}
}
