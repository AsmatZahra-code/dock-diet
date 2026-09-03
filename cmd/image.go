// image.go registers the "image" sub-command for dock-diet.
//
// Usage:
//
//	dock-diet image [image-name]
//
// This command pulls metadata (compressed size, layer count) for a published
// Docker image directly from the remote container registry and reports
// optimization warnings. It delegates all analysis to scanner.ScanRemoteImage.
//
// Example:
//
//	dock-diet image nginx:alpine
package cmd

import (
	"fmt"
	"os"

	"github.com/AsmatZahra-code/dock-diet/internal/scanner"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image [image-name]",
	Short: "Scan a remote Docker image (e.g., ubuntu:latest)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Please provide a Docker image name (e.g., dock-diet image nginx:alpine).")
			os.Exit(1)
		}

		imageName := args[0]
		err := scanner.ScanRemoteImage(imageName)
		if err != nil {
			fmt.Printf("❌ Error scanning image: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add this new command to the root command
	rootCmd.AddCommand(imageCmd)
}