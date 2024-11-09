package auth_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bootdotdev/learn-cicd-starter/internal/auth"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expectedKey string
		expectError error
	}{
		{
			name:        "Valid API Key",
			headers:     http.Header{"Authorization": []string{"ApiKey my-api-key"}},
			expectedKey: "my-api-key",
			expectError: nil,
		},
		{
			name:        "Missing Authorization Header",
			headers:     http.Header{},
			expectedKey: "",
			expectError: auth.ErrNoAuthHeaderIncluded,
		},
		{
			name:        "Malformed Authorization Header",
			headers:     http.Header{"Authorization": []string{"Bearer my-api-key"}},
			expectedKey: "",
			expectError: errors.New("malformed authorization header"),
		},
		{
			name:        "Empty Authorization Header",
			headers:     http.Header{"Authorization": []string{""}},
			expectedKey: "",
			expectError: auth.ErrNoAuthHeaderIncluded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey, err := auth.GetAPIKey(tt.headers)

			if apiKey != tt.expectedKey {
				t.Errorf("expected API key: %v, got: %v", tt.expectedKey, apiKey)
			}
			if err != nil && err.Error() != tt.expectError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
