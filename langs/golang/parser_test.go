package golang

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestGoParser_GOROUTINE_LEAK(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	// Go code with goroutine leak - go statement without errgroup
	code := `package main

func main() {
	go func() {
		println("leaked goroutine")
	}()
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var goRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "GOROUTINE_LEAK" {
			goRules = append(goRules, r)
		}
	}

	if len(goRules) == 0 {
		t.Fatal("GOROUTINE_LEAK rule not found")
	}

	parser := &GoParser{}
	findings, err := parser.Parse(testFile, content, goRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// Should find at least 1 goroutine leak
	if len(findings) == 0 {
		t.Error("expected to find goroutine leak, found none")
	}

	// Verify the finding is for GOROUTINE_LEAK
	found := false
	for _, f := range findings {
		if f.RuleID == "GOROUTINE_LEAK" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected GOROUTINE_LEAK finding, but found different rule")
	}
}

func TestGoParser_GOROUTINE_LEAK_WithErrgroup(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	// Go code with errgroup - should NOT trigger goroutine leak
	code := `package main

import "golang.org/x/sync/errgroup"

func main() {
	g := new(errgroup.Group)
	g.Go(func() error {
		return nil
	})
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var goRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "GOROUTINE_LEAK" {
			goRules = append(goRules, r)
		}
	}

	parser := &GoParser{}
	findings, err := parser.Parse(testFile, content, goRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// Should NOT find any GOROUTINE_LEAK when using errgroup
	for _, f := range findings {
		if f.RuleID == "GOROUTINE_LEAK" {
			t.Error("expected NO goroutine leak when using errgroup")
		}
	}
}

func TestGoParser_RESOUCE_LEAK(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	// Go code with resource leak - sql.Open without Close
	code := `package main

import "database/sql"

func query() {
	db, _ := sql.Open("postgres", "connection_string")
	rows, _ := db.Query("SELECT * FROM users")
	// missing db.Close()
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var goRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RESOURCE_LEAK" {
			goRules = append(goRules, r)
		}
	}

	if len(goRules) == 0 {
		t.Fatal("RESOURCE_LEAK rule not found")
	}

	parser := &GoParser{}
	findings, err := parser.Parse(testFile, content, goRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// Should find at least 1 resource leak
	if len(findings) == 0 {
		t.Error("expected to find resource leak, found none")
	}
}

func TestGoParser_Extensions(t *testing.T) {
	parser := &GoParser{}
	exts := parser.Extensions()

	if len(exts) != 1 || exts[0] != ".go" {
		t.Errorf("Extensions() returned %v, want [.go]", exts)
	}
}

func TestGoParser_Language(t *testing.T) {
	parser := &GoParser{}
	if lang := parser.Language(); lang != "go" {
		t.Errorf("Language() returned %q, want %q", lang, "go")
	}
}

func TestGoParser_ASTImportDetection(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	// Test that imports are correctly detected
	code := `package main

import (
	"golang.org/x/sync/errgroup"
	"github.com/golang-jwt/jwt/v5"
)
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)

	parser := &GoParser{}
	_, err = parser.Parse(testFile, content, allRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	_ = parser // use the parser
}
