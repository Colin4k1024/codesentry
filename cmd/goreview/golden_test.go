package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/engine"
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"

	// Import language parsers to trigger their init() functions
	_ "github.com/Colin4k1024/codesentry/langs/cpp"
	_ "github.com/Colin4k1024/codesentry/langs/golang"
	_ "github.com/Colin4k1024/codesentry/langs/java"
	_ "github.com/Colin4k1024/codesentry/langs/javascript"
	_ "github.com/Colin4k1024/codesentry/langs/kotlin"
	_ "github.com/Colin4k1024/codesentry/langs/php"
	_ "github.com/Colin4k1024/codesentry/langs/python"
	_ "github.com/Colin4k1024/codesentry/langs/ruby"
	_ "github.com/Colin4k1024/codesentry/langs/rust"
	_ "github.com/Colin4k1024/codesentry/langs/swift"
	_ "github.com/Colin4k1024/codesentry/langs/typescript"
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
	inputContent := `password = "hardcoded123456"
api_key = "sk_test_abcdef1234567890"
token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
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

	// Debug: check parser detection and rules
	t.Logf("DEBUG: tmpDir=%s, inputFile=%s", tmpDir, inputFile)
	t.Logf("DEBUG: total rules loaded: %d", len(allRules))
	for _, r := range allRules {
		if r.ID == "HARDCODED_SECRET" {
			t.Logf("DEBUG: HARDCODED_SECRET rule found, langs=%v, patterns=%d", r.Languages, len(r.Patterns))
			for i, p := range r.Patterns {
				t.Logf("DEBUG: rule pattern[%d] = %q", i, p.Pattern)
			}
		}
	}

	// Create engine with HARDCODED_SECRET rule
	e := engine.New(hardcodedRules)
	cfg := &engine.Config{}

	result, err := e.Scan([]string{inputFile}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Debug: directly check parser
	importedParser, ok := parserpkg.Get("python")
	t.Logf("DEBUG: parser.Get(python) = %v, ok=%v", importedParser, ok)
	if importedParser != nil {
		content, _ := os.ReadFile(inputFile)
		findings, _ := importedParser.Parse(inputFile, content, hardcodedRules)
		t.Logf("DEBUG: direct parser.Findings = %d: %+v", len(findings), findings)
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
		t.Logf("DEBUG: found %d issues, want %d. Issues: %+v", len(result.Issues), len(golden.Findings), result.Issues)
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
	inputContent := `query = "SELECT * FROM users" + " WHERE id=" + user_input
sql = "DELETE FROM logs" + " WHERE id=" + record_id
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
