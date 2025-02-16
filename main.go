package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/supporttools/do-database-metrics-adapter/pkg/config"
	"github.com/supporttools/do-database-metrics-adapter/pkg/digitalocean"
	"github.com/supporttools/do-database-metrics-adapter/pkg/health"
	"github.com/supporttools/do-database-metrics-adapter/pkg/logging"
)

var log = logging.SetupLogging(config.CFG.Debug)

func main() {
	flag.Parse()
	config.LoadConfiguration()
	log := logging.SetupLogging(config.CFG.Debug)
	log.Debug("Debug logging enabled")

	log.Info("Starting do-database-metrics-adapter...")

	databaseUrl := "https://api.digitalocean.com/v2/databases"
	metricsCredsUrl := "https://api.digitalocean.com/v2/databases/metrics/credentials"
	log.Debug("Using database name: ", config.CFG.DatabaseName)

	uuid, err := digitalocean.GetDatabaseUUIDByName(config.CFG.DoToken, config.CFG.DatabaseName, databaseUrl)
	if err != nil {
		log.Fatalf("Error getting database UUID: %v", err)
	}
	log.Debug("Obtained UUID for database: ", uuid)

	metricsEndpointUrl := fmt.Sprintf("https://api.digitalocean.com/v2/databases/%s", uuid)
	metricsEndpoint, err := digitalocean.GetDatabaseMetricsEndpoint(config.CFG.DoToken, metricsEndpointUrl)
	if err != nil {
		log.Fatalf("Error getting metrics endpoint: %v", err)
	}
	log.Debug("Metrics endpoint URL: ", metricsEndpoint)

	username, password, err := digitalocean.GetMetricsCredentials(config.CFG.DoToken, metricsCredsUrl)
	if err != nil {
		log.Fatalf("Error getting metrics credentials: %v", err)
	}
	log.Debug("Retrieved metrics credentials")

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, metricsEndpoint, username, password)
	})
	http.HandleFunc("/healthz", health.HealthzHandler())
	http.HandleFunc("/readyz", health.ReadyzHandler())
	http.HandleFunc("/version", health.VersionHandler())

	log.Infof("Listening on port %d", config.CFG.ServerPort)
	port := fmt.Sprintf(":%d", config.CFG.ServerPort)
	log.Fatal(http.ListenAndServe(port, nil))
}

func forwardRequest(w http.ResponseWriter, r *http.Request, endpoint, username, password string) {
	log.Debug("Forwarding request to metrics endpoint")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	log.Debug("Final metrics endpoint: ", endpoint)

	req, err := http.NewRequest("GET", endpoint+"/metrics", nil)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		log.Error("Failed to create request: ", err)
		return
	}
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		log.Error("Failed to forward request: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body: "+err.Error(), http.StatusInternalServerError)
		log.Error("Failed to read response body: ", err)
		return
	}

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Debug("Successfully forwarded and responded with metrics")
}
