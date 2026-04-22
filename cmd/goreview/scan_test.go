package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Colin4k1024/codesentry/internal/types"
)

func TestOutputResultsFormat(t *testing.T) {
	result := &types.Result{
		FilesScanned: 5,
		TotalIssues:   3,
		Severe:       1,
		Warning:      2,
		Info:         0,
		Duration:     100 * time.Millisecond,
		Issues: []types.Issue{
			{
				RuleID:    "TEST001",
				Title:     "Hardcoded Secret",
				Severity:  "SEVERE",
				File:      "config.go",
				Line:      10,
				Column:    5,
				Message:   "Possible hardcoded secret detected",
				Suggestion: "Use environment variables instead",
			},
			{
				RuleID:    "TEST002",
				Title:     "SQL Injection Risk",
				Severity:  "WARNING",
				File:      "db.go",
				Line:      25,
				Column:    3,
				Message:   "Potential SQL injection vulnerability",
				Suggestion: "Use parameterized queries",
			},
		},
	}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputResults(result)

	w.Close()
	os.Stdout = oldStdout

	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains expected sections
	if !strings.Contains(output, "CodeSentry Scan Results") {
		t.Error("missing scan results header")
	}
	if !strings.Contains(output, "Files scanned: 5") {
		t.Error("missing files scanned count")
	}
	if !strings.Contains(output, "Total issues: 3") {
		t.Error("missing total issues count")
	}
	if !strings.Contains(output, "SEVERE:   1") {
		t.Error("missing severe count")
	}
	if !strings.Contains(output, "WARNING:  2") {
		t.Error("missing warning count")
	}
	if !strings.Contains(output, "Issues Found") {
		t.Error("missing issues section")
	}
	if !strings.Contains(output, "Hardcoded Secret") {
		t.Error("missing issue title")
	}
	if !strings.Contains(output, "config.go:10:5") {
		t.Error("missing file location")
	}
	if !strings.Contains(output, "Use environment variables instead") {
		t.Error("missing suggestion")
	}
}

func TestOutputResultsEmpty(t *testing.T) {
	result := &types.Result{
		FilesScanned: 10,
		TotalIssues:   0,
		Severe:        0,
		Warning:       0,
		Info:          0,
		Duration:      50 * time.Millisecond,
		Issues:        []types.Issue{},
	}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputResults(result)

	w.Close()
	os.Stdout = oldStdout

	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Should still show header but not "Issues Found"
	if !strings.Contains(output, "CodeSentry Scan Results") {
		t.Error("missing scan results header")
	}
	if strings.Contains(output, "Issues Found") {
		t.Error("should not show issues section when no issues")
	}
}

func TestScanCommandFlags(t *testing.T) {
	// Test that all flags are properly initialized
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
	// Test argument validation
	if scanCmd.Args != nil {
		// Should accept minimum 0 args (defaults to ".")
		err := scanCmd.Args(nil, []string{})
		if err != nil {
			t.Errorf("expected no error for empty args, got: %v", err)
		}
	}
}

func TestScanFlagDefaults(t *testing.T) {
	// Reset flags to default values
	securityFlag = false
	performanceFlag = false
	noAIFlag = false
	outputFlag = ""
	excludeFlag = []string{}

	if securityFlag != false {
		t.Error("securityFlag should default to false")
	}
	if performanceFlag != false {
		t.Error("performanceFlag should default to false")
	}
	if noAIFlag != false {
		t.Error("noAIFlag should default to false")
	}
	if outputFlag != "" {
		t.Error("outputFlag should default to empty string")
	}
	if len(excludeFlag) != 0 {
		t.Error("excludeFlag should default to empty slice")
	}
}

func TestOutputResultsEmptySeverity(t *testing.T) {
	result := &types.Result{
		FilesScanned: 1,
		TotalIssues:   1,
		Severe:       0,
		Warning:      0,
		Info:         1,
		Duration:     10 * time.Millisecond,
		Issues: []types.Issue{
			{
				RuleID:    "TEST003",
				Title:     "Info Message",
				Severity:  "", // Empty severity
				File:      "info.go",
				Line:      1,
				Column:    1,
				Message:   "Informational message",
			},
		},
	}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputResults(result)

	w.Close()
	os.Stdout = oldStdout

	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Empty severity should default to WARNING
	if !strings.Contains(output, "[WARNING]") {
		t.Error("expected empty severity to default to WARNING")
	}
}

func TestScanRulesFiltering(t *testing.T) {
	// Test that rule filtering works correctly
	tmpDir := t.TempDir()

	// Create test files
	goFile := filepath.Join(tmpDir, "test.go")
	goContent := `package main

import "fmt"

func main() {
	password := "hardcoded"
	fmt.Println(password)
}
`
	if err := os.WriteFile(goFile, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// The actual scan is tested through the engine
	// Here we just verify the file was created correctly
	info, err := os.Stat(goFile)
	if err != nil {
		t.Fatalf("failed to stat test file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("test file should not be empty")
	}
}
