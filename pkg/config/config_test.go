package config

import (
	"os"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	tests := []struct {
		name                 string
		envVariables         map[string]string
		expectedDebug        bool
		expectedServerPort   int
		expectedDoToken      string
		expectedDatabaseName string
	}{
		{
			name: "All Environment Variables Set",
			envVariables: map[string]string{
				"DEBUG":         "true",
				"PORT":          "8080",
				"DO_TOKEN":      "test-token",
				"DATABASE_NAME": "test-db",
			},
			expectedDebug:        true,
			expectedServerPort:   8080,
			expectedDoToken:      "test-token",
			expectedDatabaseName: "test-db",
		},
		{
			name: "Only Required Environment Variables Set",
			envVariables: map[string]string{
				"DEBUG": "false",
				"PORT":  "9000",
			},
			expectedDebug:        false,
			expectedServerPort:   9000,
			expectedDoToken:      "",
			expectedDatabaseName: "db-test-for-metrics.c.db.ondigitalocean.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key, value := range test.envVariables {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			LoadConfiguration()

			if CFG.Debug != test.expectedDebug {
				t.Errorf("Debug: got %t, want %t", CFG.Debug, test.expectedDebug)
			}
			if CFG.ServerPort != test.expectedServerPort {
				t.Errorf("ServerPort: got %d, want %d", CFG.ServerPort, test.expectedServerPort)
			}
			if CFG.DoToken != test.expectedDoToken {
				t.Errorf("DoToken: got %s, want %s", CFG.DoToken, test.expectedDoToken)
			}
			if CFG.DatabaseName != test.expectedDatabaseName {
				t.Errorf("DatabaseName: got %s, want %s", CFG.DatabaseName, test.expectedDatabaseName)
			}
		})
	}
}
