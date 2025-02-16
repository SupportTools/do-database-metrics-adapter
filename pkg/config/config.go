package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig structure for environment-based configurations.
type AppConfig struct {
	Debug        bool   `json:"debug"`
	ServerPort   int    `json:"serverPort"`
	DoToken      string `json:"doToken"`
	DatabaseName string `json:"databaseName"`
}

var CFG AppConfig

// LoadConfiguration loads the configuration from the environment variables.
func LoadConfiguration() {
	CFG.Debug = parseEnvBool("DEBUG", false)
	CFG.ServerPort = parseEnvInt("PORT", 9000)
	CFG.DoToken = getEnvOrDefault("DO_TOKEN", "")
	CFG.DatabaseName = getEnvOrDefault("DATABASE_NAME", "db-test-for-metrics.c.db.ondigitalocean.com")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing %s as int: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

func parseEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing %s as bool: %v. Using default value: %t", key, err, defaultValue)
		return defaultValue
	}
	return boolValue
}
