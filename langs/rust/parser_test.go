package rust

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestRustParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.rs")

	// Rust code with hardcoded credentials
	code := `let password = "hardcoded123";
let api_key = "sk-abcdefgh123456";
let token = "Bearer xxx.yyy.zzz";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Rust rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var rustRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RS_HARDCODED_SECRET" {
			rustRules = append(rustRules, r)
		}
	}

	if len(rustRules) == 0 {
		t.Skip("No hardcoded_secret rules found for Rust")
	}

	parser := &RustParser{}
	findings, err := parser.Parse(testFile, content, rustRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestRustParser_UNSAFE_BLOCK(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.rs")

	// Rust code with unsafe block
	code := `unsafe {
    let ptr = transmute::<_, *const u8>(42);
    println!("{:?}", ptr);
}`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Rust rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var rustRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RUST_UNSAFE_BLOCK" {
			rustRules = append(rustRules, r)
		}
	}

	if len(rustRules) == 0 {
		t.Skip("No unsafe_block rules found for Rust")
	}

	parser := &RustParser{}
	findings, err := parser.Parse(testFile, content, rustRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find unsafe block, found none")
	}
}

func TestRustParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.rs")

	// Rust code with SQL injection vulnerability - using format! with SELECT
	code := `let query = format!("SELECT * FROM users WHERE id={}", user_id);`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Rust rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var rustRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RS_SQL_INJECTION" {
			rustRules = append(rustRules, r)
		}
	}

	if len(rustRules) == 0 {
		t.Skip("No SQL injection rules found for Rust")
	}

	parser := &RustParser{}
	findings, err := parser.Parse(testFile, content, rustRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestRustParser_Extensions(t *testing.T) {
	parser := &RustParser{}
	exts := parser.Extensions()

	if len(exts) != 1 || exts[0] != ".rs" {
		t.Errorf("Extensions() returned %v, want [.rs]", exts)
	}
}

func TestRustParser_Language(t *testing.T) {
	parser := &RustParser{}
	if lang := parser.Language(); lang != "rust" {
		t.Errorf("Language() returned %q, want %q", lang, "rust")
	}
}
