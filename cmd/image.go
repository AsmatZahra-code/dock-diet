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