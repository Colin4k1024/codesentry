package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goreview/goreview/internal/engine"
	"github.com/goreview/goreview/internal/rules"
	"github.com/goreview/goreview/internal/types"
	"github.com/spf13/cobra"
)

var (
	securityFlag    bool
	performanceFlag bool
	noAIFlag        bool
	outputFlag      string
	excludeFlag     []string
)

var scanCmd = &cobra.Command{
	Use:   "scan [paths...]",
	Short: "Scan code for issues",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		paths := args
		if len(paths) == 0 {
			paths = []string{"."}
		}

		// Load rules
		rulesDir := "rules"
		loadedRules := rules.LoadRules(rulesDir)
		if len(loadedRules) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: No rules loaded from %s\n", rulesDir)
		}

		// Filter rules by category if flags are set
		var filteredRules []rules.Rule
		for _, r := range loadedRules {
			if securityFlag && r.Category == "security" {
				filteredRules = append(filteredRules, r)
			}
			if performanceFlag && r.Category == "performance" {
				filteredRules = append(filteredRules, r)
			}
		}
		if !securityFlag && !performanceFlag {
			filteredRules = loadedRules
		}

		// Create engine and scan
		eng := engine.New(filteredRules)
		cfg := &engine.Config{
			Security:    securityFlag,
			Performance: performanceFlag,
			NoAI:        noAIFlag,
			Output:      outputFlag,
			Exclude:     excludeFlag,
		}

		start := time.Now()
		result, err := eng.Scan(paths, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during scan: %v\n", err)
			os.Exit(1)
		}
		result.Duration = time.Since(start)

		// Output results
		outputResults(result)
	},
}

func outputResults(result *types.Result) {
	fmt.Printf("\n=== GoReview Scan Results ===\n")
	fmt.Printf("Files scanned: %d\n", result.FilesScanned)
	fmt.Printf("Total issues: %d\n", result.TotalIssues)
	fmt.Printf("  SEVERE:   %d\n", result.Severe)
	fmt.Printf("  WARNING:  %d\n", result.Warning)
	fmt.Printf("  INFO:     %d\n", result.Info)
	fmt.Printf("Duration: %v\n", result.Duration)

	if len(result.Issues) > 0 {
		fmt.Printf("\n=== Issues Found ===\n")
		for _, issue := range result.Issues {
			severity := issue.Severity
			if severity == "" {
				severity = "WARNING"
			}
			fmt.Printf("[%s] %s\n", severity, issue.Title)
			fmt.Printf("  File: %s:%d:%d\n", issue.File, issue.Line, issue.Column)
			fmt.Printf("  %s\n", issue.Message)
			if issue.Suggestion != "" {
				fmt.Printf("  Suggestion: %s\n", issue.Suggestion)
			}
			fmt.Println()
		}
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolVar(&securityFlag, "security", false, "Enable security rules")
	scanCmd.Flags().BoolVar(&performanceFlag, "performance", false, "Enable performance rules")
	scanCmd.Flags().BoolVar(&noAIFlag, "no-ai", false, "Disable AI-powered analysis")
	scanCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file for results")
	scanCmd.Flags().StringArrayVar(&excludeFlag, "exclude", []string{}, "Paths to exclude from scanning")
}
