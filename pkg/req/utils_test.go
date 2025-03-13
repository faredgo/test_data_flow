package req_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"test_data_flow/pkg/req"
	"testing"
)

func TestGetIPAddress(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		expectedIP string
	}{
		{
			name:       "X-Forwarded-For header present",
			req:        httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil),
			expectedIP: "192.168.1.1",
		},
		{
			name:       "X-Real-IP header present",
			req:        httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil),
			expectedIP: "192.168.1.1",
		},
		{
			name:       "RemoteAddr with IPv4",
			req:        httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil),
			expectedIP: "192.168.1.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "X-Forwarded-For header present" {
				tt.req.Header.Add("X-Forwarded-For", "192.168.1.1")
			} else if tt.name == "X-Real-IP header present" {
				tt.req.Header.Add("X-Real-IP", "192.168.1.1")
			} else if tt.name == "RemoteAddr with IPv4" {
				tt.req.RemoteAddr = "192.168.1.2:8080"
			} else if tt.name == "IPv6 localhost address" {
				tt.req.RemoteAddr = "::1:8080"
			} else if tt.name == "Empty headers, use RemoteAddr" {
				tt.req.RemoteAddr = "192.0.2.1:8080"
			}

			gotIP := req.GetIPAddress(tt.req)
			if gotIP != tt.expectedIP {
				t.Errorf("GetIPAddress() = %v, want %v", gotIP, tt.expectedIP)
			}
		})
	}
}
