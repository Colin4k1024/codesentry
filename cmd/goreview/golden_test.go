package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/engine"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

// GoldenFile represents expected findings
type GoldenFile struct {
	RuleID   string          `json:"rule_id"`
	Findings []GoldenFinding `json:"findings"`
}

type GoldenFinding struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	EndLine  int    `json:"end_line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func TestGoldenFile_HARDCODED_SECRET(t *testing.T) {
	// Create temp input file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.py")
	inputContent := `password = "hardcoded123"
api_key = "secret_key_456"
token = "token_789"
`
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create golden file
	goldenFile := filepath.Join(tmpDir, "golden.json")
	goldenContent := `{
  "rule_id": "HARDCODED_SECRET",
  "findings": [
    {"line": 1, "column": 1, "end_line": 1, "severity": "SEVERE", "message": "Possible hardcoded secret"},
    {"line": 2, "column": 1, "end_line": 2, "severity": "SEVERE", "message": "Possible hardcoded secret"},
    {"line": 3, "column": 1, "end_line": 3, "severity": "SEVERE", "message": "Possible hardcoded secret"}
  ]
}
`
	if err := os.WriteFile(goldenFile, []byte(goldenContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Load rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var hardcodedRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "HARDCODED_SECRET" {
			hardcodedRules = append(hardcodedRules, r)
		}
	}

	if len(hardcodedRules) == 0 {
		t.Fatal("HARDCODED_SECRET rule not found")
	}

	// Create engine with HARDCODED_SECRET rule
	e := engine.New(hardcodedRules)
	cfg := &engine.Config{}

	result, err := e.Scan([]string{inputFile}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Load golden file
	goldenData, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	// Parse golden file
	var golden GoldenFile
	if err := json.Unmarshal(goldenData, &golden); err != nil {
		t.Fatalf("failed to parse golden file: %v", err)
	}

	// Compare
	if len(result.Issues) != len(golden.Findings) {
		t.Errorf("found %d issues, want %d", len(result.Issues), len(golden.Findings))
		return
	}

	for i, want := range golden.Findings {
		got := result.Issues[i]
		if got.Line != want.Line {
			t.Errorf("issue[%d].Line = %d, want %d", i, got.Line, want.Line)
		}
		if got.Severity != want.Severity {
			t.Errorf("issue[%d].Severity = %q, want %q", i, got.Severity, want.Severity)
		}
	}
}

func TestGoldenFile_SQL_INJECTION(t *testing.T) {
	// Create temp input file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.py")
	inputContent := `query = "SELECT * FROM users WHERE id=" + user_id
sql = "INSERT INTO logs VALUES ('" + data + "')"
`
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Load rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var sqlRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "SQL_INJECTION" {
			sqlRules = append(sqlRules, r)
		}
	}

	if len(sqlRules) == 0 {
		t.Fatal("SQL_INJECTION rule not found")
	}

	// Create engine
	e := engine.New(sqlRules)
	cfg := &engine.Config{}

	result, err := e.Scan([]string{inputFile}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should find 2 SQL injection issues
	if len(result.Issues) != 2 {
		t.Errorf("found %d issues, want 2", len(result.Issues))
	}
}

func TestGoldenFile_GOROUTINE_LEAK(t *testing.T) {
	// Create temp input file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.go")
	inputContent := `package main

func main() {
	go func() {
		println("leaked")
	}()
}
`
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Load rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var goroutineRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "GOROUTINE_LEAK" {
			goroutineRules = append(goroutineRules, r)
		}
	}

	if len(goroutineRules) == 0 {
		t.Fatal("GOROUTINE_LEAK rule not found")
	}

	// Create engine
	e := engine.New(goroutineRules)
	cfg := &engine.Config{}

	result, err := e.Scan([]string{inputFile}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should find at least 1 goroutine issue (AST-based check)
	if len(result.Issues) < 1 {
		t.Errorf("found %d issues, want at least 1", len(result.Issues))
	}
}
