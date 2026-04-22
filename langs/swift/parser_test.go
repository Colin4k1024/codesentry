package swift

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestSwiftParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.swift")

	// Swift code with hardcoded credentials
	code := `let password = "hardcoded123"
let apiKey = "sk-abcdefgh123456"
let token = "Bearer xxx.yyy.zzz"`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Swift rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var swiftRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "SWIFT_HARDCODED_SECRET" {
			swiftRules = append(swiftRules, r)
		}
	}

	if len(swiftRules) == 0 {
		t.Skip("No hardcoded_secret rules found for Swift")
	}

	parser := &SwiftParser{}
	findings, err := parser.Parse(testFile, content, swiftRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestSwiftParser_SENSITIVE_LOG(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.swift")

	// Swift code with sensitive data logging
	code := `print("User password: \(password)")
NSLog("API token: \(token)")
logger.debug("Secret key: \(apiKey)")`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Swift rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var swiftRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "SWIFT_SENSITIVE_LOG" {
			swiftRules = append(swiftRules, r)
		}
	}

	if len(swiftRules) == 0 {
		t.Skip("No sensitive_log rules found for Swift")
	}

	parser := &SwiftParser{}
	findings, err := parser.Parse(testFile, content, swiftRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find sensitive logging, found none")
	}
}

func TestSwiftParser_Extensions(t *testing.T) {
	parser := &SwiftParser{}
	exts := parser.Extensions()

	if len(exts) != 1 || exts[0] != ".swift" {
		t.Errorf("Extensions() returned %v, want [.swift]", exts)
	}
}

func TestSwiftParser_Language(t *testing.T) {
	parser := &SwiftParser{}
	if lang := parser.Language(); lang != "swift" {
		t.Errorf("Language() returned %q, want %q", lang, "swift")
	}
}
