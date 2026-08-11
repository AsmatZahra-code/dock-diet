package scanner

import (
	"os"
	"strings"
)

// CI/CD JSON Output ke liye Structs
type Issue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ScanResult struct {
	Score      int      `json:"score"`
	Grade      string   `json:"grade"`
	Issues     []Issue  `json:"issues"`
	Lines      []string `json:"-"` // "-" ka matlab hai JSON mein isay hide rakhna
	NeedsFix   bool     `json:"needs_fix"`
}

// AnalyzeDockerfile asal scanning karta hai
func AnalyzeDockerfile(filePath string) (ScanResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ScanResult{}, err
	}

	lines := strings.Split(string(content), "\n")
	result := ScanResult{Score: 100, Lines: lines}

	fromCount := 0
	runCount := 0
	hasUser := false

	for _, line := range lines {
		upperLine := strings.ToUpper(strings.TrimSpace(line))

		// 1. Fat Base Image & Latest Tag
		if strings.HasPrefix(upperLine, "FROM ") {
			fromCount++
			if strings.Contains(upperLine, ":LATEST") {
				result.Issues = append(result.Issues, Issue{"Tag", "Avoid using ':latest' tag in base images."})
				result.Score -= 10
			}
			if !strings.Contains(strings.ToLower(line), "alpine") && !strings.Contains(strings.ToLower(line), "slim") {
				result.Issues = append(result.Issues, Issue{"BaseImage", "'Fat' base image detected. Use -alpine or -slim."})
				result.Score -= 20
			}
		}

		// 2. Layer Consolidation (Counting RUNs)
		if strings.HasPrefix(upperLine, "RUN ") {
			runCount++
		}

		// 3. Apt-Get Cache Cleanup
		if strings.Contains(upperLine, "APT-GET") && !strings.Contains(upperLine, "RM -RF /VAR/LIB/APT/LISTS") {
			result.Issues = append(result.Issues, Issue{"Cache", "apt-get used without clearing cache."})
			result.Score -= 15
		}

		// 4. Root User Check
		if strings.HasPrefix(upperLine, "USER ") {
			hasUser = true
		}
	}

	// Layer Consolidation Rule: Agar 2 se zyada RUN hain
	if runCount > 2 {
		result.Issues = append(result.Issues, Issue{"Layers", "Too many RUN instructions. Consolidate them using '&& \\' to save layers."})
		result.Score -= 15
	}

	// Multi-stage & User Rules
	if fromCount < 2 {
		result.Issues = append(result.Issues, Issue{"MultiStage", "No multi-stage build detected."})
		result.Score -= 10
	}
	if !hasUser {
		result.Issues = append(result.Issues, Issue{"Security", "Container runs as root. No 'USER' instruction found."})
		result.Score -= 20
	}

	// Grade Calculation
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score >= 90 {
		result.Grade = "A 🏆"
	} else if result.Score >= 75 {
		result.Grade = "B 🥈"
	} else if result.Score >= 60 {
		result.Grade = "C 🥉"
	} else {
		result.Grade = "D ⚠️"
	}

	result.NeedsFix = len(result.Issues) > 0
	return result, nil
}