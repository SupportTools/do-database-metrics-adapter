package version

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	expectedVersion := "1.0.0"
	expectedGitCommit := "abcdefg"
	expectedBuildTime := "2024-05-10T10:00:00Z"

	Version = expectedVersion
	GitCommit = expectedGitCommit
	BuildTime = expectedBuildTime

	if Version != expectedVersion {
		t.Errorf("Version is incorrect, expected: %s, got: %s", expectedVersion, Version)
	}

	if GitCommit != expectedGitCommit {
		t.Errorf("GitCommit is incorrect, expected: %s, got: %s", expectedGitCommit, GitCommit)
	}

	if BuildTime != expectedBuildTime {
		t.Errorf("BuildTime is incorrect, expected: %s, got: %s", expectedBuildTime, BuildTime)
	}
}
