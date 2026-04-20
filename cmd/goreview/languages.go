package main

import (
	"fmt"

	"github.com/goreview/goreview/internal/parser"
	"github.com/spf13/cobra"
)

var languagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "List supported programming languages",
	Run: func(cmd *cobra.Command, args []string) {
		langs := parser.List()
		fmt.Println("Supported languages:")
		for _, lang := range langs {
			fmt.Printf("  - %s\n", lang)
		}
	},
}

func init() {
	rootCmd.AddCommand(languagesCmd)
}
