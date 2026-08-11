package scanner

import (
	"fmt"
	"os"
	"strings"
)

// AutoFix ghaltiyon ko theek kar ke nayi file banata hai
func AutoFix(result ScanResult, originalPath string) error {
	var optimizedLines []string
	hasUser := false

	for _, line := range result.Lines {
		trimmedLine := strings.TrimSpace(line)
		upperLine := strings.ToUpper(trimmedLine)

		// Fix 1: Add cache clear to apt-get
		if strings.HasPrefix(upperLine, "RUN APT-GET") && !strings.Contains(upperLine, "RM -RF") {
			line = line + " && rm -rf /var/lib/apt/lists/*"
		}

		// Track USER
		if strings.HasPrefix(upperLine, "USER ") {
			hasUser = true
		}

		optimizedLines = append(optimizedLines, line)
	}

	// Fix 2: Add non-root user at the end if missing
	if !hasUser {
		optimizedLines = append(optimizedLines, "\n# Auto-Fixed: Added non-root user for security")
		optimizedLines = append(optimizedLines, "RUN useradd -m appuser")
		optimizedLines = append(optimizedLines, "USER appuser")
	}

	// Nayi file create karna
	newFilePath := originalPath + ".optimized"
	finalContent := strings.Join(optimizedLines, "\n")
	
	err := os.WriteFile(newFilePath, []byte(finalContent), 0644)
	if err == nil {
		fmt.Printf("🛠️  Auto-Fix applied! Optimized file created at: %s\n", newFilePath)
		fmt.Println("⚠️  Note: Base image changes and Multi-stage builds must be fixed manually to prevent breaking your app.")
	}
	return err
}