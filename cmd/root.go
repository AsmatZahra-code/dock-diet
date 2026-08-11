package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dock-diet",
	Short: "Dock-Diet is a container optimization CLI",
	Long:  `A fast and flexible CLI tool built in Go to analyze and optimize Docker images for cloud-native environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Dock-Diet! Use 'dock-diet --help' to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}