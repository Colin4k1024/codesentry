package typescript

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestTypeScriptParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	// TypeScript code with hardcoded credentials
	code := `const password = "hardcoded123";
const apiKey = "sk-abcdefgh123456";
const token = "Bearer xxx.yyy.zzz";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load TypeScript rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var tsRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "TS_HARDCODED_SECRET" {
			tsRules = append(tsRules, r)
		}
	}

	if len(tsRules) == 0 {
		t.Skip("No hardcoded_secret rules found for TypeScript")
	}

	parser := &TSParser{}
	findings, err := parser.Parse(testFile, content, tsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestTypeScriptParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	// TypeScript code with SQL injection vulnerability
	code := `db.query("SELECT * FROM users WHERE id=" + userId);
connection.query("DELETE FROM users WHERE id=" + userId);`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load TypeScript rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var tsRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "TS_SQL_INJECTION" {
			tsRules = append(tsRules, r)
		}
	}

	if len(tsRules) == 0 {
		t.Skip("No SQL injection rules found for TypeScript")
	}

	parser := &TSParser{}
	findings, err := parser.Parse(testFile, content, tsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestTypeScriptParser_EVAL(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	// TypeScript code with dangerous eval usage
	code := `eval(userInput);
const result = eval("(" + userData + ")")`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load TypeScript rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var tsRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "TS_EVAL" {
			tsRules = append(tsRules, r)
		}
	}

	if len(tsRules) == 0 {
		t.Skip("No eval rules found for TypeScript")
	}

	parser := &TSParser{}
	findings, err := parser.Parse(testFile, content, tsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find dangerous eval, found none")
	}
}

func TestTypeScriptParser_PROTOTYPE_POLLUTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	// TypeScript code with prototype pollution vulnerability - must use merge with spread
	code := `merge({}, ...req.body);
lodash.merge({}, ...userData);`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load TypeScript rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var tsRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "TSPrototype_POLLUTION" {
			tsRules = append(tsRules, r)
		}
	}

	if len(tsRules) == 0 {
		t.Skip("No prototype pollution rules found for TypeScript")
	}

	parser := &TSParser{}
	findings, err := parser.Parse(testFile, content, tsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find prototype pollution, found none")
	}
}

func TestTypeScriptParser_Extensions(t *testing.T) {
	parser := &TSParser{}
	exts := parser.Extensions()

	if len(exts) != 4 {
		t.Errorf("Extensions() returned %d extensions, want 4", len(exts))
	}
}

func TestTypeScriptParser_Language(t *testing.T) {
	parser := &TSParser{}
	if lang := parser.Language(); lang != "typescript" {
		t.Errorf("Language() returned %q, want %q", lang, "typescript")
	}
}
