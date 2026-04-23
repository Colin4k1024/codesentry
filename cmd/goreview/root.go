package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "goreview",
	Short: "CodeSentry - AI-powered multi-language code review tool",
	Long: `CodeSentry is a static analysis and AI-powered code review tool that supports
multiple programming languages including Go, JavaScript/TypeScript, Python, and more.

Examples:
  goreview scan ./...
  goreview scan --security --performance
  goreview languages
  goreview version
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'goreview --help' for more information")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of CodeSentry",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("CodeSentry version %s\n", version)
		},
	})
}
