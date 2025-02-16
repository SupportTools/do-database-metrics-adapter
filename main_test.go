package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForwardRequest(t *testing.T) {
	// Create a TLS config that accepts any certificate (since it's a test)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create a mock HTTPS server
	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock response"))
	}))
	mockServer.TLS = tlsConfig
	mockServer.StartTLS()
	defer mockServer.Close()

	// Call the forwardRequest function with the mock server endpoint
	req, err := http.NewRequest("GET", mockServer.URL+"/metrics", nil)
	if err != nil {
		t.Fatal("Creating request failed:", err)
	}
	w := httptest.NewRecorder()

	username := "testuser"
	password := "testpassword"
	forwardRequest(w, req, mockServer.URL, username, password)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "mock response", w.Body.String())
}
