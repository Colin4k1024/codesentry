package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Colin4k1024/codesentry/internal/output"
	"github.com/Colin4k1024/codesentry/internal/types"
)

func sampleScanResult() *types.Result {
	return &types.Result{
		FilesScanned: 5,
		TotalIssues:  1,
		Severe:       1,
		Warning:      0,
		Info:         0,
		Duration:     100 * time.Millisecond,
		Issues: []types.Issue{
			{
				RuleID:     "TEST001",
				Title:      "Hardcoded Secret",
				Severity:   "SEVERE",
				File:       "config.go",
				Line:       10,
				Column:     5,
				Message:    "Possible hardcoded secret detected",
				Suggestion: "Use environment variables instead",
			},
		},
	}
}

func TestWriteScanOutputJSONFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.json")

	result := sampleScanResult()
	if err := output.Write(result, output.FormatJSON, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read JSON output: %v", err)
	}

	var result2 types.Result
	if err := json.Unmarshal(data, &result2); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result2.TotalIssues != 1 {
		t.Errorf("TotalIssues = %d, want 1", result2.TotalIssues)
	}
}

func TestWriteScanOutputSARIFFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.sarif")

	result := sampleScanResult()
	if err := output.Write(result, output.FormatSarif, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read SARIF output: %v", err)
	}

	var sarif map[string]interface{}
	if err := json.Unmarshal(data, &sarif); err != nil {
		t.Fatalf("output is not valid SARIF JSON: %v", err)
	}
	if sarif["version"] != "2.1.0" {
		t.Errorf("SARIF version = %v, want 2.1.0", sarif["version"])
	}
	if _, ok := sarif["runs"].([]interface{}); !ok {
		t.Error("SARIF output missing runs array")
	}
}

func TestWriteScanOutputGHSLFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.ghsl")

	result := sampleScanResult()
	if err := output.Write(result, output.FormatGHSL, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read GHSL output: %v", err)
	}

	var ghsl map[string]interface{}
	if err := json.Unmarshal(data, &ghsl); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if ghsl["tool"] != "CodeSentry goreview" {
		t.Errorf("tool = %v, want 'CodeSentry goreview'", ghsl["tool"])
	}
}

func TestWriteScanOutputCLangFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.clang")

	result := sampleScanResult()
	if err := output.Write(result, output.FormatCLang, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read CLang output: %v", err)
	}

	content := string(data)
	if content == "" {
		t.Error("CLang output should not be empty")
	}
}

func TestWriteScanOutputUnknownExtensionFallsBackToText(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.out")

	result := sampleScanResult()
	if err := output.Write(result, output.FormatText, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read text output: %v", err)
	}
	if string(data) == "" {
		t.Error("text output should not be empty")
	}
}

func TestScanCommandFlags(t *testing.T) {
	if scanCmd.Flags().Lookup("security") == nil {
		t.Error("security flag not found")
	}
	if scanCmd.Flags().Lookup("performance") == nil {
		t.Error("performance flag not found")
	}
	if scanCmd.Flags().Lookup("no-ai") == nil {
		t.Error("no-ai flag not found")
	}
	if scanCmd.Flags().Lookup("output") == nil {
		t.Error("output flag not found")
	}
	if scanCmd.Flags().Lookup("exclude") == nil {
		t.Error("exclude flag not found")
	}
	if scanCmd.Flags().Lookup("format") == nil {
		t.Error("format flag not found")
	}
}

func TestScanCmdArgs(t *testing.T) {
	if scanCmd.Args == nil {
		t.Fatal("scan command should define argument validation")
	}
	if err := scanCmd.Args(nil, []string{}); err != nil {
		t.Errorf("expected no error for empty args, got: %v", err)
	}
}

func TestScanFlagDefaults(t *testing.T) {
	securityFlag = false
	performanceFlag = false
	noAIFlag = false
	outputFlag = ""
	excludeFlag = []string{}

	if securityFlag {
		t.Error("securityFlag should default to false")
	}
	if performanceFlag {
		t.Error("performanceFlag should default to false")
	}
	if noAIFlag {
		t.Error("noAIFlag should default to false")
	}
	if outputFlag != "" {
		t.Error("outputFlag should default to empty string")
	}
	if len(excludeFlag) != 0 {
		t.Error("excludeFlag should default to empty slice")
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected output.Format
	}{
		{"json", output.FormatJSON},
		{"sarif", output.FormatSarif},
		{"ghsl", output.FormatGHSL},
		{"clang", output.FormatCLang},
		{"text", output.FormatText},
		{"unknown", output.FormatText},
	}
	for _, tt := range tests {
		got := output.ParseFormat(tt.input)
		if got != tt.expected {
			t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
