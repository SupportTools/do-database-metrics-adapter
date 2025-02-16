package logging

import (
	"os"
	"testing"

	"github.com/supporttools/do-database-metrics-adapter/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Use a mock for Getwd which can be toggled in tests
var osGetwd = os.Getwd

func TestLogFile(t *testing.T) {
	originalDebug := config.CFG.Debug

	// Test logging in debug mode
	config.CFG.Debug = true
	log := LogFile()
	assert.NotNil(t, log)
	assert.NotEmpty(t, log.Data["line"])
	if log.Data["filename"] == "" {
		t.Error("Expected filename to be populated in debug mode")
	}

	// Reset logrus fields for the next test and restore config
	logrus.StandardLogger().WithFields(logrus.Fields{})
	config.CFG.Debug = false

	// Test logging in non-debug mode
	log = LogFile()
	assert.NotNil(t, log)
	assert.Empty(t, log.Data["filename"], "filename should be empty in non-debug mode")
	assert.NotEmpty(t, log.Data["line"])

	// Restore original debug setting
	config.CFG.Debug = originalDebug
}
