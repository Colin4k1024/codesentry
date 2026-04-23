package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Colin4k1024/codesentry/internal/engine"
	"github.com/Colin4k1024/codesentry/internal/output"
	"github.com/Colin4k1024/codesentry/internal/rules"
	"github.com/Colin4k1024/codesentry/internal/types"
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
		if err := writeScanOutput(result); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing scan output: %v\n", err)
			os.Exit(1)
		}
	},
}

func writeScanOutput(result *types.Result) error {
	return output.Write(result, output.FormatForPath(outputFlag), outputFlag)
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolVar(&securityFlag, "security", false, "Enable security rules")
	scanCmd.Flags().BoolVar(&performanceFlag, "performance", false, "Enable performance rules")
	scanCmd.Flags().BoolVar(&noAIFlag, "no-ai", false, "Disable AI-powered analysis")
	scanCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file for results")
	scanCmd.Flags().StringArrayVar(&excludeFlag, "exclude", []string{}, "Paths to exclude from scanning")
}
