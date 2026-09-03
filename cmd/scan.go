// scan.go registers the "scan" sub-command for dock-diet.
//
// Usage:
//
//	dock-diet scan [file-path] [flags]
//
// Flags:
//
//	-o, --output string   Output format: "text" (default) or "json"
//	-f, --fix             Auto-fix detected issues and write a .optimized file
//
// The command loads the optional .dock-diet.yaml configuration to determine
// the fail_under score threshold. If the Dockerfile score falls below this
// threshold, the process exits with code 1, enabling CI/CD pipeline gating.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AsmatZahra-code/dock-diet/internal/scanner"
	"github.com/spf13/cobra"
)

var outputFormat string
var autoFix bool

var scanCmd = &cobra.Command{
	Use:   "scan [file-path]",
	Short: "Advanced Dockerfile optimization scanner",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Please provide a Dockerfile path.")
			return
		}

		filePath := args[0]
		
		// Load User Configuration
		config := scanner.LoadConfig()

		result, err := scanner.AnalyzeDockerfile(filePath)
		if err != nil {
			fmt.Printf("❌ Error: Could not read file '%s'.\n", filePath)
			os.Exit(1)
		}

		if outputFormat == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
			
			// Fail based on config threshold
			if result.Score < config.FailUnder {
				os.Exit(1)
			}
			return
		}

		fmt.Printf("🔍 Scanned: %s\n", filePath)
		fmt.Printf("⚙️  Target Score: %d (From Config)\n", config.FailUnder)
		fmt.Println("-------------------------------------------------")
		
		for i, issue := range result.Issues {
			fmt.Printf("⚠️  ISSUE %d [%s]: %s\n", i+1, issue.Type, issue.Description)
		}
		
		fmt.Println("-------------------------------------------------")
		fmt.Printf("📊 DIET SCORE: %d/100 (Grade: %s)\n", result.Score, result.Grade)

		if autoFix && result.NeedsFix {
			scanner.AutoFix(result, filePath)
		}

		// CI/CD Failure Logic based on Config
		if result.Score < config.FailUnder {
			fmt.Printf("❌ Pipeline Failed: Score %d is below the required threshold of %d.\n", result.Score, config.FailUnder)
			os.Exit(1)
		} else {
			fmt.Println("✅ Pipeline Passed: Score meets the required threshold.")
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text or json)")
	scanCmd.Flags().BoolVarP(&autoFix, "fix", "f", false, "Auto-fix issues and generate Dockerfile.optimized")
	
	rootCmd.AddCommand(scanCmd)
}