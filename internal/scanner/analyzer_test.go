// analyzer_test.go contains the full table-driven test suite for AnalyzeDockerfile.
//
// Each test case loads a purpose-built fixture from the testdata/ directory,
// runs the analyser, and verifies the expected score, grade, issue types, and
// NeedsFix flag. False-positive guard tests confirm that correct Dockerfiles
// are not penalised by mistake.
package scanner

import (
	"os"
	"strings"
	"testing"
)

// testdataPath returns the absolute path to a fixture file inside testdata/.
func testdataPath(name string) string {
	return "testdata/" + name
}

// hasIssueType reports whether any issue in the result has the given type.
func hasIssueType(result ScanResult, issueType string) bool {
	for _, issue := range result.Issues {
		if issue.Type == issueType {
			return true
		}
	}
	return false
}

// ── Positive rule tests (each fixture triggers exactly one rule) ──────────────

func TestAnalyze_PerfectDockerfile(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("perfect.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 100 {
		t.Errorf("expected Score=100, got %d", result.Score)
	}
	if result.Grade != "A 🏆" {
		t.Errorf("expected Grade=A, got %q", result.Grade)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d: %v", len(result.Issues), result.Issues)
	}
	if result.NeedsFix {
		t.Error("expected NeedsFix=false for a perfect Dockerfile")
	}
}

func TestAnalyze_LatestTag(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("latest_tag.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 90 {
		t.Errorf("expected Score=90 (Tag penalty -10), got %d", result.Score)
	}
	if !hasIssueType(result, "Tag") {
		t.Errorf("expected Tag issue, got %v", result.Issues)
	}
	if !result.NeedsFix {
		t.Error("expected NeedsFix=true")
	}
}

func TestAnalyze_FatImage(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("fat_image.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 80 {
		t.Errorf("expected Score=80 (BaseImage penalty -20), got %d", result.Score)
	}
	if !hasIssueType(result, "BaseImage") {
		t.Errorf("expected BaseImage issue, got %v", result.Issues)
	}
}

func TestAnalyze_NoAptCleanup(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("no_apt_cleanup.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 85 {
		t.Errorf("expected Score=85 (Cache penalty -15), got %d", result.Score)
	}
	if !hasIssueType(result, "Cache") {
		t.Errorf("expected Cache issue, got %v", result.Issues)
	}
}

func TestAnalyze_NoUser(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("no_user.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 80 {
		t.Errorf("expected Score=80 (Security penalty -20), got %d", result.Score)
	}
	if !hasIssueType(result, "Security") {
		t.Errorf("expected Security issue, got %v", result.Issues)
	}
}

func TestAnalyze_NoMultiStage(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("no_multistage.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 90 {
		t.Errorf("expected Score=90 (MultiStage penalty -10), got %d", result.Score)
	}
	if !hasIssueType(result, "MultiStage") {
		t.Errorf("expected MultiStage issue, got %v", result.Issues)
	}
}

func TestAnalyze_TooManyRuns(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("too_many_runs.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 85 {
		t.Errorf("expected Score=85 (Layers penalty -15), got %d", result.Score)
	}
	if !hasIssueType(result, "Layers") {
		t.Errorf("expected Layers issue, got %v", result.Issues)
	}
}

// ── All-issues and score-floor test ──────────────────────────────────────────

func TestAnalyze_AllIssues(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("all_issues.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}

	// Score must be clamped to 0, never negative.
	if result.Score != 0 {
		t.Errorf("expected Score=0 (floored from negative), got %d", result.Score)
	}
	if result.Grade != "D ⚠️" {
		t.Errorf("expected Grade=D, got %q", result.Grade)
	}

	// All 6 distinct issue types must be present.
	requiredTypes := []string{"Tag", "BaseImage", "Cache", "Layers", "MultiStage", "Security"}
	for _, issueType := range requiredTypes {
		if !hasIssueType(result, issueType) {
			t.Errorf("expected issue type %q to be present, but it was not. Issues: %v", issueType, result.Issues)
		}
	}
}

func TestAnalyze_ScoreFloor_NeverNegative(t *testing.T) {
	// The all_issues fixture accumulates more than 100 points of penalties.
	// The score must be clamped to exactly 0, not a negative number.
	result, err := AnalyzeDockerfile(testdataPath("all_issues.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score < 0 {
		t.Errorf("score must never be negative, got %d", result.Score)
	}
}

// ── False-positive guard tests ────────────────────────────────────────────────

func TestAnalyze_MultiStageAlias_NoFalsePositive(t *testing.T) {
	// A stage alias used as a subsequent FROM base must NOT be flagged as fat.
	result, err := AnalyzeDockerfile(testdataPath("edge_multistage.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if hasIssueType(result, "BaseImage") {
		t.Error("multi-stage alias should NOT be flagged as a fat BaseImage")
	}
	if result.Score != 100 {
		t.Errorf("expected perfect score=100, got %d. Issues: %v", result.Score, result.Issues)
	}
}

func TestAnalyze_AptCleanup_NoFalsePositive(t *testing.T) {
	// apt-get WITH "rm -rf /var/lib/apt/lists/*" must NOT produce a Cache issue.
	result, err := AnalyzeDockerfile(testdataPath("apt_already_clean.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if hasIssueType(result, "Cache") {
		t.Error("apt-get WITH cache cleanup should NOT be flagged as a Cache issue")
	}
}

func TestAnalyze_SlimImage_NoFalsePositive(t *testing.T) {
	// A "-slim" image variant must NOT be flagged as a fat BaseImage.
	result, err := AnalyzeDockerfile(testdataPath("slim_variant.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if hasIssueType(result, "BaseImage") {
		t.Error("slim image should NOT be flagged as a fat BaseImage")
	}
}

// ── Error handling ────────────────────────────────────────────────────────────

func TestAnalyze_FileNotFound(t *testing.T) {
	_, err := AnalyzeDockerfile("testdata/nonexistent.Dockerfile")
	if err == nil {
		t.Error("expected an error for a non-existent file, got nil")
	}
}

func TestAnalyze_EmptyFile(t *testing.T) {
	// An empty Dockerfile should not crash — it will trigger MultiStage and Security.
	tmp, err := os.CreateTemp("", "empty.Dockerfile")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	result, err := AnalyzeDockerfile(tmp.Name())
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed on empty file: %v", err)
	}
	if result.Score < 0 {
		t.Errorf("score must not be negative, got %d", result.Score)
	}
}

// ── Grade boundary tests ──────────────────────────────────────────────────────

func TestAnalyze_GradeA_Boundary(t *testing.T) {
	// latest_tag.Dockerfile scores exactly 90 — the Grade A boundary.
	result, err := AnalyzeDockerfile(testdataPath("latest_tag.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 90 {
		t.Errorf("expected Score=90, got %d", result.Score)
	}
	if !strings.Contains(result.Grade, "A") {
		t.Errorf("expected Grade A at score=90, got %q", result.Grade)
	}
}

func TestAnalyze_GradeB_Boundary(t *testing.T) {
	// grade_b_boundary.Dockerfile scores exactly 75 — the Grade B boundary.
	result, err := AnalyzeDockerfile(testdataPath("grade_b_boundary.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 75 {
		t.Errorf("expected Score=75, got %d", result.Score)
	}
	if !strings.Contains(result.Grade, "B") {
		t.Errorf("expected Grade B at score=75, got %q", result.Grade)
	}
}

func TestAnalyze_GradeC_Boundary(t *testing.T) {
	// grade_c_boundary.Dockerfile scores exactly 60 — the Grade C boundary.
	result, err := AnalyzeDockerfile(testdataPath("grade_c_boundary.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.Score != 60 {
		t.Errorf("expected Score=60, got %d", result.Score)
	}
	if !strings.Contains(result.Grade, "C") {
		t.Errorf("expected Grade C at score=60, got %q", result.Grade)
	}
}

func TestAnalyze_GradeD(t *testing.T) {
	// all_issues.Dockerfile scores 0 — well within Grade D territory.
	result, err := AnalyzeDockerfile(testdataPath("all_issues.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if !strings.Contains(result.Grade, "D") {
		t.Errorf("expected Grade D at score=%d, got %q", result.Score, result.Grade)
	}
}

// ── NeedsFix flag ─────────────────────────────────────────────────────────────

func TestAnalyze_NeedsFix_TrueWhenIssuesExist(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("all_issues.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if !result.NeedsFix {
		t.Error("expected NeedsFix=true when issues are present")
	}
}

func TestAnalyze_NeedsFix_FalseWhenPerfect(t *testing.T) {
	result, err := AnalyzeDockerfile(testdataPath("perfect.Dockerfile"))
	if err != nil {
		t.Fatalf("AnalyzeDockerfile failed: %v", err)
	}
	if result.NeedsFix {
		t.Error("expected NeedsFix=false when no issues are present")
	}
}