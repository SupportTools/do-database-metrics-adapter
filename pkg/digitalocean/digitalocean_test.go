package digitalocean

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDatabaseMetricsEndpoint(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
		w.WriteHeader(http.StatusOK) // Ensure to set status code to 200 OK.
		w.Write([]byte(`{"database": {"metrics_endpoints": [{"host": "example.com", "port": 1234}]}}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	token := "test-token"
	endpoint, err := GetDatabaseMetricsEndpoint(token, server.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedEndpoint := "example.com:1234"
	if endpoint != expectedEndpoint {
		t.Errorf("Expected endpoint '%s', got '%s'", expectedEndpoint, endpoint)
	}
}

func TestGetMetricsCredentials(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
		w.Write([]byte(`{"credentials": {"basic_auth_username": "user", "basic_auth_password": "pass"}}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	token := "test-token"

	username, password, err := GetMetricsCredentials(token, server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedUsername := "user"
	if username != expectedUsername {
		t.Errorf("Expected username '%s', got '%s'", expectedUsername, username)
	}

	expectedPassword := "pass"
	if password != expectedPassword {
		t.Errorf("Expected password '%s', got '%s'", expectedPassword, password)
	}
}

func TestGetDatabaseUUIDByName(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"databases": [{"id": "some-uuid", "name": "test-db"}]}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	token := "test-token"
	dbName := "test-db"

	uuid, err := GetDatabaseUUIDByName(token, dbName, server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedUUID := "some-uuid"
	if uuid != expectedUUID {
		t.Errorf("Expected UUID '%s', got '%s'", expectedUUID, uuid)
	}
}
