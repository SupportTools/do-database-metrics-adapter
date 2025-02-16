package digitalocean

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// databasesResponse represents the JSON structure for responses from database listings.
type databasesResponse struct {
	Databases []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"databases"`
}

// MetricsEndpoint defines the structure for the metrics endpoint details.
type MetricsEndpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// DatabaseMetricsResponse encapsulates the list of available metrics endpoints.
type DatabaseMetricsResponse struct {
	MetricsEndpoints []MetricsEndpoint `json:"metrics_endpoints"`
}

// Database represents the core details of a database including its metrics endpoints.
type Database struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	MetricsEndpoints []MetricsEndpoint `json:"metrics_endpoints"`
	// Include other fields as necessary
}

// DatabaseResponse wraps the Database object, matching the JSON response structure.
type DatabaseResponse struct {
	Database Database `json:"database"`
}

// GetDatabaseMetricsEndpoint fetches the first available metrics endpoint for a given database.
func GetDatabaseMetricsEndpoint(token, url string) (string, error) {
	log.Printf("Starting request for metrics endpoint at URL: %s", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var dbResponse DatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&dbResponse); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if len(dbResponse.Database.MetricsEndpoints) == 0 {
		log.Println("No metrics endpoints found in response")
		return "", fmt.Errorf("no metrics endpoints found")
	}

	endpoint := dbResponse.Database.MetricsEndpoints[0]
	log.Printf("Metrics endpoint found: %s:%d", endpoint.Host, endpoint.Port)
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port), nil
}

// GetMetricsCredentials retrieves the basic auth credentials for accessing metrics.
func GetMetricsCredentials(token, url string) (string, string, error) {
	log.Printf("Requesting metrics credentials from URL: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return "", "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var creds map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&creds); err != nil {
		log.Printf("Error decoding response: %v", err)
		return "", "", fmt.Errorf("decoding response: %w", err)
	}

	username := creds["credentials"]["basic_auth_username"]
	password := creds["credentials"]["basic_auth_password"]
	log.Printf("Metrics credentials retrieved: username=%s, password=[REDACTED]", username)
	return username, password, nil
}

// GetDatabaseUUIDByName searches for a database UUID by its name.
func GetDatabaseUUIDByName(token, dbName, url string) (string, error) {
	log.Printf("Searching for UUID of database named '%s' at URL: %s", dbName, url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var result databasesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding response: %v", err)
		return "", fmt.Errorf("decoding response: %w", err)
	}

	for _, db := range result.Databases {
		if db.Name == dbName {
			log.Printf("Database UUID found: %s", db.ID)
			return db.ID, nil
		}
	}

	log.Printf("Database with name '%s' not found", dbName)
	return "", fmt.Errorf("database with name %s not found", dbName)
}
