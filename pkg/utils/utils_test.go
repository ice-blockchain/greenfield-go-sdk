package utils

import (
	"net/url"
	"testing"
)

func TestCheckEndpointUrl(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name: "valid domain no path",
			url:  "https://sp.example.com",
		},
		{
			name: "valid domain with trailing slash",
			url:  "https://sp.example.com/",
		},
		{
			name: "valid domain with base path",
			url:  "https://sp.example.com/base/path",
		},
		{
			name: "valid IP address",
			url:  "https://192.168.1.1",
		},
		{
			name: "valid IP with port",
			url:  "https://192.168.1.1:9090",
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid host",
			url:     "https://-invalid.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u url.URL
			if tt.url != "" {
				parsed, err := url.Parse(tt.url)
				if err != nil {
					t.Fatalf("invalid test URL: %v", err)
				}
				u = *parsed
			}
			err := checkEndpointUrl(u)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetEndpointURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
		secure   bool
		wantURL  string
		wantErr  bool
	}{
		{
			name:     "simple domain secure",
			endpoint: "sp.example.com",
			secure:   true,
			wantURL:  "https://sp.example.com",
		},
		{
			name:     "simple domain insecure",
			endpoint: "sp.example.com",
			secure:   false,
			wantURL:  "http://sp.example.com",
		},
		{
			name:     "domain with port",
			endpoint: "sp.example.com:9090",
			secure:   true,
			wantURL:  "https://sp.example.com:9090",
		},
		{
			name:     "domain with http prefix stripped",
			endpoint: "http://sp.example.com",
			secure:   true,
			wantURL:  "https://sp.example.com",
		},
		{
			name:     "domain with base path",
			endpoint: "sp.example.com/base/path",
			secure:   true,
			wantURL:  "https://sp.example.com/base/path",
		},
		{
			name:     "IP address",
			endpoint: "192.168.1.1",
			secure:   true,
			wantURL:  "https://192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEndpointURL(tt.endpoint, tt.secure)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.wantURL {
				t.Errorf("got %q, want %q", got.String(), tt.wantURL)
			}
		})
	}
}
