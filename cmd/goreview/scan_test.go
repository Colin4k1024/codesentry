package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func withOutputFlag(t *testing.T, value string) {
	t.Helper()
	old := outputFlag
	outputFlag = value
	t.Cleanup(func() {
		outputFlag = old
	})
}

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	os.Stdout = w

	fnErr := fn()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close stdout writer: %v", err)
	}
	os.Stdout = oldStdout
	if fnErr != nil {
		t.Fatalf("function returned error: %v", fnErr)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}
	return buf.String()
}

func TestWriteScanOutputStdoutText(t *testing.T) {
	withOutputFlag(t, "")

	output := captureStdout(t, func() error {
		return writeScanOutput(sampleScanResult())
	})

	if !strings.Contains(output, "CodeSentry Scan Results") {
		t.Error("missing scan results header")
	}
	if !strings.Contains(output, "Files scanned: 5") {
		t.Error("missing files scanned count")
	}
	if !strings.Contains(output, "config.go:10:5") {
		t.Error("missing file location")
	}
	if !strings.Contains(output, "Use environment variables instead") {
		t.Error("missing suggestion")
	}
}

func TestWriteScanOutputJSONFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.json")
	withOutputFlag(t, outputPath)

	if err := writeScanOutput(sampleScanResult()); err != nil {
		t.Fatalf("writeScanOutput returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read JSON output: %v", err)
	}

	var result types.Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result.TotalIssues != 1 {
		t.Errorf("TotalIssues = %d, want 1", result.TotalIssues)
	}
}

func TestWriteScanOutputSARIFFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.sarif")
	withOutputFlag(t, outputPath)

	if err := writeScanOutput(sampleScanResult()); err != nil {
		t.Fatalf("writeScanOutput returned error: %v", err)
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

func TestWriteScanOutputUnknownExtensionFallsBackToText(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "report.out")
	withOutputFlag(t, outputPath)

	if err := writeScanOutput(sampleScanResult()); err != nil {
		t.Fatalf("writeScanOutput returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read text output: %v", err)
	}
	if !strings.Contains(string(data), "CodeSentry Scan Results") {
		t.Error("unknown extension should write text output")
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
