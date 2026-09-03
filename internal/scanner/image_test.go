// image_test.go tests ScanRemoteImage for error handling and (optionally) live
// registry behaviour.
//
// Real registry calls require an internet connection, so those tests are skipped
// automatically via t.Skip in standard CI. Only the error-path test (invalid
// image name) is always executed since it never touches the network.
package scanner

import (
	"os"
	"testing"
)

// ── Error handling (no network required) ─────────────────────────────────────

func TestScanRemoteImage_InvalidName_ReturnsError(t *testing.T) {
	// A completely malformed image reference must return a non-nil error
	// and must not panic. No network call is made because name.ParseReference
	// fails synchronously before any I/O.
	err := ScanRemoteImage(":::invalid:::name:::")
	if err == nil {
		t.Error("expected an error for an invalid image name, got nil")
	}
}

func TestScanRemoteImage_EmptyName_ReturnsError(t *testing.T) {
	err := ScanRemoteImage("")
	if err == nil {
		t.Error("expected an error for an empty image name, got nil")
	}
}

// ── Live registry tests (skipped in offline/CI environments) ─────────────────

// requiresNetwork skips the test if the CI environment variable is set, or if
// the caller passes the -short flag (go test -short), indicating a fast/offline run.
func requiresNetwork(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping network test: -short flag provided")
	}
	if os.Getenv("CI") != "" {
		t.Skip("skipping network test in CI environment")
	}
}

func TestScanRemoteImage_Live_AlpineLatest(t *testing.T) {
	requiresNetwork(t)
	// alpine:latest is a tiny image — it should not trigger size or layer warnings.
	err := ScanRemoteImage("alpine:latest")
	if err != nil {
		t.Errorf("unexpected error scanning alpine:latest: %v", err)
	}
}
