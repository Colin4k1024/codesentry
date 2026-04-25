package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// GoldenFile represents the expected findings for a rule test
type GoldenFile struct {
	RuleID   string          `json:"rule_id"`
	Findings []GoldenFinding `json:"findings"`
}

// GoldenFinding represents a single expected finding in a golden file
type GoldenFinding struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	EndLine  int    `json:"end_line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// loadGoldenFile loads a golden file from testdata/rules.
//
//nolint:unused // kept for future golden file testing
func loadGoldenFile(t *testing.T, ruleID string) *GoldenFile {
	t.Helper()

	// Validate ruleID doesn't contain path separators (security check)
	if strings.ContainsAny(ruleID, "/\\") {
		t.Fatalf("ruleID must not contain path separators: %q", ruleID)
	}

	// Look in testdata/rules/
	path := filepath.Join("..", "..", "testdata", "rules", ruleID+".golden.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("golden file not found: %s", path)
	}

	var golden GoldenFile
	if err := json.Unmarshal(data, &golden); err != nil {
		t.Fatalf("failed to parse golden file %s: %v", path, err)
	}

	return &golden
}

// CompareFindings compares actual findings with golden file expected findings
func CompareFindings(t *testing.T, actual []Finding, golden *GoldenFile) {
	t.Helper()

	if len(actual) != len(golden.Findings) {
		t.Errorf("found %d issues, want %d", len(actual), len(golden.Findings))
		return
	}

	for i := range actual {
		if actual[i].Line != golden.Findings[i].Line {
			t.Errorf("finding[%d].Line = %d, want %d", i, actual[i].Line, golden.Findings[i].Line)
		}
		if actual[i].Severity != golden.Findings[i].Severity {
			t.Errorf("finding[%d].Severity = %q, want %q", i, actual[i].Severity, golden.Findings[i].Severity)
		}
		if actual[i].Message != golden.Findings[i].Message {
			t.Errorf("finding[%d].Message = %q, want %q", i, actual[i].Message, golden.Findings[i].Message)
		}
	}
}

// DiffFindings returns a detailed diff between actual and expected findings using go-cmp
func DiffFindings(actual []Finding, golden *GoldenFile) string {
	return cmp.Diff(golden.Findings, actual,
		cmp.Transformer("GoldenFinding", func(f Finding) GoldenFinding {
			return GoldenFinding{
				Line:     f.Line,
				Column:   f.Column,
				EndLine:  f.EndLine,
				Severity: f.Severity,
				Message:  f.Message,
			}
		}))
}
