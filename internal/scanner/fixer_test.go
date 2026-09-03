// fixer_test.go contains tests for AutoFix, verifying that it correctly
// patches apt-get instructions and injects a non-root user when needed,
// without producing duplicates when those fixes are already present.
package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildScanResultFromLines constructs a ScanResult directly from a slice of
// Dockerfile lines. This avoids touching the filesystem for fixer input.
func buildScanResultFromLines(lines []string) ScanResult {
	return ScanResult{
		Lines:    lines,
		NeedsFix: true,
	}
}

// readOptimized reads the .optimized file that AutoFix writes next to originalPath.
func readOptimized(t *testing.T, originalPath string) string {
	t.Helper()
	content, err := os.ReadFile(originalPath + ".optimized")
	if err != nil {
		t.Fatalf("expected .optimized file to exist at %s: %v", originalPath+".optimized", err)
	}
	return string(content)
}

// ── USER injection tests ──────────────────────────────────────────────────────

func TestAutoFix_MissingUser_InjectsUserAppuser(t *testing.T) {
	lines := []string{
		"FROM alpine:3.20",
		"COPY app /app",
		"ENTRYPOINT [\"/app\"]",
	}
	result := buildScanResultFromLines(lines)
	tmpPath := filepath.Join(t.TempDir(), "Dockerfile")

	if err := AutoFix(result, tmpPath); err != nil {
		t.Fatalf("AutoFix returned error: %v", err)
	}

	content := readOptimized(t, tmpPath)
	if !strings.Contains(content, "USER appuser") {
		t.Errorf("expected 'USER appuser' in .optimized content, got:\n%s", content)
	}
}

func TestAutoFix_UserAlreadyPresent_NoDuplicate(t *testing.T) {
	lines := []string{
		"FROM alpine:3.20",
		"COPY app /app",
		"USER nobody",
		"ENTRYPOINT [\"/app\"]",
	}
	result := buildScanResultFromLines(lines)
	tmpPath := filepath.Join(t.TempDir(), "Dockerfile")

	if err := AutoFix(result, tmpPath); err != nil {
		t.Fatalf("AutoFix returned error: %v", err)
	}

	content := readOptimized(t, tmpPath)

	// "USER appuser" must NOT appear — the existing USER should be kept as-is.
	if strings.Contains(content, "USER appuser") {
		t.Errorf("AutoFix must not add 'USER appuser' when a USER instruction already exists:\n%s", content)
	}
	// The original USER instruction must still be present.
	if !strings.Contains(content, "USER nobody") {
		t.Errorf("original USER instruction was unexpectedly removed:\n%s", content)
	}
}

// ── apt-get cleanup tests ─────────────────────────────────────────────────────

func TestAutoFix_AptNoCleanup_AppendsRmRf(t *testing.T) {
	lines := []string{
		"FROM debian:bookworm-slim",
		"RUN apt-get update && apt-get install -y curl",
		"USER nobody",
	}
	result := buildScanResultFromLines(lines)
	tmpPath := filepath.Join(t.TempDir(), "Dockerfile")

	if err := AutoFix(result, tmpPath); err != nil {
		t.Fatalf("AutoFix returned error: %v", err)
	}

	content := readOptimized(t, tmpPath)
	if !strings.Contains(content, "rm -rf /var/lib/apt/lists/*") {
		t.Errorf("expected 'rm -rf /var/lib/apt/lists/*' to be appended, got:\n%s", content)
	}
}

func TestAutoFix_AptAlreadyClean_NoDuplicateRmRf(t *testing.T) {
	cleanup := "RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*"
	lines := []string{
		"FROM debian:bookworm-slim",
		cleanup,
		"USER nobody",
	}
	result := buildScanResultFromLines(lines)
	tmpPath := filepath.Join(t.TempDir(), "Dockerfile")

	if err := AutoFix(result, tmpPath); err != nil {
		t.Fatalf("AutoFix returned error: %v", err)
	}

	content := readOptimized(t, tmpPath)

	// The cleanup string should appear exactly once, not twice.
	count := strings.Count(content, "rm -rf /var/lib/apt/lists/*")
	if count != 1 {
		t.Errorf("expected 'rm -rf /var/lib/apt/lists/*' to appear exactly once, got %d occurrences:\n%s", count, content)
	}
}

// ── Output file tests ─────────────────────────────────────────────────────────

func TestAutoFix_OutputFileIsCreated(t *testing.T) {
	lines := []string{"FROM alpine:3.20", "COPY app /app"}
	result := buildScanResultFromLines(lines)
	tmpPath := filepath.Join(t.TempDir(), "Dockerfile")

	if err := AutoFix(result, tmpPath); err != nil {
		t.Fatalf("AutoFix returned error: %v", err)
	}

	optimizedPath := tmpPath + ".optimized"
	if _, err := os.Stat(optimizedPath); os.IsNotExist(err) {
		t.Errorf("expected .optimized file to exist at %s", optimizedPath)
	}
}

func TestAutoFix_InvalidOutputPath_ReturnsError(t *testing.T) {
	lines := []string{"FROM alpine:3.20"}
	result := buildScanResultFromLines(lines)

	// Pass a path inside a non-existent directory — WriteFile must fail.
	err := AutoFix(result, "/nonexistent/path/Dockerfile")
	if err == nil {
		t.Error("expected AutoFix to return an error for an invalid output path, got nil")
	}
}
