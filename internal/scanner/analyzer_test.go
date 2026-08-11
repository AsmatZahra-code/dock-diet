package scanner

import (
	"os"
	"testing"
)

func TestAnalyzeDockerfile(t *testing.T) {
	// 1. ARRANGE: Create a temporary bad Dockerfile for testing
	badDockerfileContent := `FROM ubuntu:latest
RUN apt-get update
COPY . /app`

	// Create a temporary file that will be automatically deleted after the test
	tmpFile, err := os.CreateTemp("", "Dockerfile.test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Cleanup after test execution completes

	tmpFile.WriteString(badDockerfileContent)
	tmpFile.Close()

	// 2. ACT: Execute the core scanning function
	result, err := AnalyzeDockerfile(tmpFile.Name())
	
	if err != nil {
		t.Fatalf("AnalyzeDockerfile execution failed: %v", err)
	}

	// 3. ASSERT: Verify the expected outcomes
	// Our dummy file has known issues (latest tag, no user, no cache cleanup).
	// Therefore, the score should be penalized and not equal 100.
	if result.Score == 100 {
		t.Errorf("Expected score to be penalized, but got %d", result.Score)
	}

	// We expect at least 3 issues to be flagged based on our rule set
	if len(result.Issues) < 3 {
		t.Errorf("Expected multiple issues, but found only %d", len(result.Issues))
	}
}