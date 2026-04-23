package output

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

func testResult() *types.Result {
	return &types.Result{
		FilesScanned: 2,
		TotalIssues:  1,
		Severe:       1,
		Duration:     25 * time.Millisecond,
		Issues: []types.Issue{
			{
				RuleID:     "HARDCODED_SECRET",
				Title:      "Hardcoded Secret",
				Severity:   types.SEVERE,
				File:       "app.go",
				Line:       12,
				Column:     3,
				Message:    "Possible hardcoded secret",
				Suggestion: "Use environment variables",
			},
		},
	}
}

func TestFormatForPath(t *testing.T) {
	tests := []struct {
		path string
		want Format
	}{
		{"", FormatText},
		{"results.txt", FormatText},
		{"results.out", FormatText},
		{"results.json", FormatJSON},
		{"RESULTS.JSON", FormatJSON},
		{"results.sarif", FormatSarif},
		{"RESULTS.SARIF", FormatSarif},
	}

	for _, tt := range tests {
		if got := FormatForPath(tt.path); got != tt.want {
			t.Errorf("FormatForPath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestWriteTextToStdout(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	os.Stdout = w

	writeErr := Write(testResult(), FormatText, "")

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close stdout writer: %v", err)
	}
	os.Stdout = oldStdout
	if writeErr != nil {
		t.Fatalf("Write returned error: %v", writeErr)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "CodeSentry Scan Results") {
		t.Error("text output missing header")
	}
	if !strings.Contains(output, "app.go:12:3") {
		t.Error("text output missing issue location")
	}
}

func TestWriteJSONFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "results.json")
	if err := Write(testResult(), FormatJSON, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var result types.Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON output is invalid: %v", err)
	}
	if result.FilesScanned != 2 {
		t.Errorf("FilesScanned = %d, want 2", result.FilesScanned)
	}
}

func TestWriteSARIFFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "results.sarif")
	if err := Write(testResult(), FormatSarif, outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var sarif map[string]interface{}
	if err := json.Unmarshal(data, &sarif); err != nil {
		t.Fatalf("SARIF output is invalid JSON: %v", err)
	}
	if sarif["version"] != "2.1.0" {
		t.Errorf("SARIF version = %v, want 2.1.0", sarif["version"])
	}
	runs, ok := sarif["runs"].([]interface{})
	if !ok || len(runs) != 1 {
		t.Fatalf("SARIF runs = %#v, want one run", sarif["runs"])
	}
	run := runs[0].(map[string]interface{})
	tool := run["tool"].(map[string]interface{})
	driver := tool["driver"].(map[string]interface{})
	if driver["version"] != toolVersion {
		t.Errorf("SARIF tool version = %v, want %s", driver["version"], toolVersion)
	}
}

func TestWriteUnknownFormatFallsBackToText(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "results.unknown")
	if err := Write(testResult(), Format("unknown"), outputPath); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "CodeSentry Scan Results") {
		t.Error("unknown format should fall back to text")
	}
}
