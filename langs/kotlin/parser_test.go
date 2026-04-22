package kotlin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestKotlinParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.kt")

	// Kotlin code with hardcoded credentials
	code := `val password = "hardcoded123"
val apiKey = "sk-abcdefgh123456"
val token = "Bearer xxx.yyy.zzz"`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Kotlin rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var kotlinRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "KT_HARDCODED_SECRET" {
			kotlinRules = append(kotlinRules, r)
		}
	}

	if len(kotlinRules) == 0 {
		t.Skip("No hardcoded_secret rules found for Kotlin")
	}

	parser := &KotlinParser{}
	findings, err := parser.Parse(testFile, content, kotlinRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestKotlinParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.kt")

	// Kotlin code with SQL injection vulnerability - must use query/execute/rawQuery with +
	code := `db.query("SELECT * FROM users WHERE id=" + userId)
db.execute("DELETE FROM users WHERE id=" + userId)`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Kotlin rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var kotlinRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "KT_SQL_INJECTION" {
			kotlinRules = append(kotlinRules, r)
		}
	}

	if len(kotlinRules) == 0 {
		t.Skip("No SQL injection rules found for Kotlin")
	}

	parser := &KotlinParser{}
	findings, err := parser.Parse(testFile, content, kotlinRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestKotlinParser_Extensions(t *testing.T) {
	parser := &KotlinParser{}
	exts := parser.Extensions()

	expected := []string{".kt", ".kts"}
	if len(exts) != len(expected) {
		t.Errorf("Extensions() returned %d extensions, want %d", len(exts), len(expected))
	}
}

func TestKotlinParser_Language(t *testing.T) {
	parser := &KotlinParser{}
	if lang := parser.Language(); lang != "kotlin" {
		t.Errorf("Language() returned %q, want %q", lang, "kotlin")
	}
}
