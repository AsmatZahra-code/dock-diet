// config_test.go tests LoadConfig for correct default values, valid YAML parsing,
// graceful fallback on invalid or missing files, and correct field unmarshalling.
//
// LoadConfig reads ".dock-diet.yaml" from the current working directory, so each
// test that needs a custom config file changes the working directory to a temporary
// directory and restores it with defer after the test completes.
package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// chdir temporarily changes the working directory to dir and returns a function
// that restores the original directory. Call it with defer.
func chdir(t *testing.T, dir string) func() {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("could not chdir to %s: %v", dir, err)
	}
	return func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("could not restore working directory: %v", err)
		}
	}
}

// writeConfig writes content to .dock-diet.yaml in dir.
func writeConfig(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, ".dock-diet.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("could not write config file: %v", err)
	}
}

// ── Default config ────────────────────────────────────────────────────────────

func TestLoadConfig_Defaults_WhenFileAbsent(t *testing.T) {
	// Use a temp directory that has no .dock-diet.yaml file.
	tmpDir := t.TempDir()
	defer chdir(t, tmpDir)()

	cfg := LoadConfig()

	if cfg.FailUnder != 100 {
		t.Errorf("expected default FailUnder=100, got %d", cfg.FailUnder)
	}
	if len(cfg.Ignore) != 0 {
		t.Errorf("expected empty Ignore slice, got %v", cfg.Ignore)
	}
}

// ── Valid YAML ────────────────────────────────────────────────────────────────

func TestLoadConfig_ValidFile_ReadsFailUnder(t *testing.T) {
	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "fail_under: 80\n")
	defer chdir(t, tmpDir)()

	cfg := LoadConfig()

	if cfg.FailUnder != 80 {
		t.Errorf("expected FailUnder=80, got %d", cfg.FailUnder)
	}
}

func TestLoadConfig_ValidFile_ReadsIgnoreRules(t *testing.T) {
	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "fail_under: 70\nignore_rules:\n  - Tag\n  - Cache\n")
	defer chdir(t, tmpDir)()

	cfg := LoadConfig()

	if cfg.FailUnder != 70 {
		t.Errorf("expected FailUnder=70, got %d", cfg.FailUnder)
	}
	if len(cfg.Ignore) != 2 {
		t.Fatalf("expected 2 ignore rules, got %d: %v", len(cfg.Ignore), cfg.Ignore)
	}
	if cfg.Ignore[0] != "Tag" || cfg.Ignore[1] != "Cache" {
		t.Errorf("expected ignore_rules=[Tag Cache], got %v", cfg.Ignore)
	}
}

// ── Fallback on bad input ─────────────────────────────────────────────────────

func TestLoadConfig_InvalidYAML_FallsBackToDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "this: is: not: valid: yaml: :::\n")
	defer chdir(t, tmpDir)()

	cfg := LoadConfig()

	if cfg.FailUnder != 100 {
		t.Errorf("expected fallback FailUnder=100 on invalid YAML, got %d", cfg.FailUnder)
	}
}

func TestLoadConfig_EmptyFile_FallsBackToDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "")
	defer chdir(t, tmpDir)()

	cfg := LoadConfig()

	if cfg.FailUnder != 100 {
		t.Errorf("expected fallback FailUnder=100 on empty file, got %d", cfg.FailUnder)
	}
}
