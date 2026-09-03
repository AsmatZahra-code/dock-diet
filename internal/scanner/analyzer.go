// Package scanner provides Dockerfile static analysis and automated fix capabilities.
//
// analyzer.go implements AnalyzeDockerfile, the core scanning engine.
// It parses a Dockerfile line-by-line and applies the following rule set:
//
//	Rule         Penalty   Trigger
//	Tag            -10    :latest tag on any FROM instruction
//	BaseImage      -20    Non-alpine/slim image (fat base image)
//	Cache          -15    apt-get used without clearing /var/lib/apt/lists/*
//	Layers         -15    More than 2 RUN instructions (layer bloat)
//	MultiStage     -10    Fewer than 2 FROM instructions (no multi-stage build)
//	Security       -20    No USER instruction (container runs as root)
//
// The final score is clamped to [0, 100] and mapped to a letter grade:
//
//	Grade A  score >= 90
//	Grade B  score >= 75
//	Grade C  score >= 60
//	Grade D  score <  60
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
	Score    int      `json:"score"`
	Grade    string   `json:"grade"`
	Issues   []Issue  `json:"issues"`
	Lines    []string `json:"-"` // "-" ka matlab hai JSON mein isay hide rakhna
	NeedsFix bool     `json:"needs_fix"`
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

	// Multi-stage stage names track karne ke liye map
	validStages := make(map[string]bool)

	for _, line := range lines {
		upperLine := strings.ToUpper(strings.TrimSpace(line))

		// Skip Dockerfile comment lines. Without this guard, comment text that
		// contains keywords (e.g. "apt-get", "USER", "FROM") would incorrectly
		// trigger analysis rules against the comments rather than the instructions.
		if strings.HasPrefix(upperLine, "#") {
			continue
		}

		if strings.HasPrefix(upperLine, "FROM ") {
			fromCount++
			if strings.Contains(upperLine, ":LATEST") {
				result.Issues = append(result.Issues, Issue{"Tag", "Avoid using ':latest' tag in base images."})
				result.Score -= 10
			}

			// --- SMART MULTI-STAGE PARSING FIX ---
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				imageName := strings.ToLower(parts[1]) // Yahan se exact image name mil jayega (jaise node:22-alpine ya base)

				// Agar image name mein alpine/slim hai YA wo kisi pehle declare ki gayi stage ka naam hai (jaise base)
				if !strings.Contains(imageName, "alpine") && !strings.Contains(imageName, "slim") && !validStages[imageName] {
					result.Issues = append(result.Issues, Issue{"BaseImage", "'Fat' base image detected. Use -alpine or -slim."})
					result.Score -= 20
				}
			}

			// Agar is line mein AS alias hai (jaise "AS base" ya "AS dev"), to usay valid stages mein save kar lo
			for i, part := range parts {
				if strings.ToUpper(part) == "AS" && i+1 < len(parts) {
					stageName := strings.ToLower(parts[i+1])
					validStages[stageName] = true
				}
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