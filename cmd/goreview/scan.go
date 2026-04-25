package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Colin4k1024/codesentry/internal/engine"
	"github.com/Colin4k1024/codesentry/internal/output"
	"github.com/Colin4k1024/codesentry/internal/rules"
	"github.com/spf13/cobra"
)

var (
	securityFlag    bool
	performanceFlag bool
	noAIFlag        bool
	noColorFlag     bool
	outputFlag      string
	formatFlag      string
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

		// Set color mode: only override auto-detection if --no-color is passed
		if noColorFlag {
			output.SetColorEnabled(false)
		}
		// Output results
		if outputFlag != "" || formatFlag != "" {
			outFmt := output.FormatText
			if formatFlag != "" {
				outFmt = output.ParseFormat(formatFlag)
			}
			if err := output.Write(result, outFmt, outputFlag); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
				os.Exit(1)
			}
		} else {
			output.Write(result, output.FormatText, "")
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolVar(&securityFlag, "security", false, "Enable security rules")
	scanCmd.Flags().BoolVar(&performanceFlag, "performance", false, "Enable performance rules")
	scanCmd.Flags().BoolVar(&noAIFlag, "no-ai", false, "Disable AI-powered analysis")
	scanCmd.Flags().BoolVar(&noColorFlag, "no-color", false, "Disable colored output")
	scanCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file for results")
	scanCmd.Flags().StringVar(&formatFlag, "format", "", "Output format: text, json, sarif, ghsl, clang")
	scanCmd.Flags().StringArrayVar(&excludeFlag, "exclude", []string{}, "Paths to exclude from scanning")
}
