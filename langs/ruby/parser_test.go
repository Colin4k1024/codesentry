package ruby

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestRubyParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.rb")

	// Ruby code with hardcoded credentials
	code := `password = "hardcoded123"
api_key = "sk-abcdefgh123456"
token = "Bearer xxx.yyy.zzz"
SECRET_KEY = "my_secret_value"`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Ruby rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var rubyRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RB_HARDCODED_SECRET" {
			rubyRules = append(rubyRules, r)
		}
	}

	if len(rubyRules) == 0 {
		t.Skip("No hardcoded_secret rules found for Ruby")
	}

	parser := &RubyParser{}
	findings, err := parser.Parse(testFile, content, rubyRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestRubyParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.rb")

	// Ruby code with SQL injection vulnerability - must use string interpolation #{}
	code := `query("SELECT * FROM users WHERE id=#{user_id}")
exec("DELETE FROM users WHERE id=#{user_id}")`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Ruby rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var rubyRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RB_SQL_INJECTION" {
			rubyRules = append(rubyRules, r)
		}
	}

	if len(rubyRules) == 0 {
		t.Skip("No SQL injection rules found for Ruby")
	}

	parser := &RubyParser{}
	findings, err := parser.Parse(testFile, content, rubyRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestRubyParser_YAML_LOAD(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.rb")

	// Ruby code with unsafe YAML load
	code := `data = YAML.load(user_input)
config = YAML.load(File.read("config.yml"))`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Ruby rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var rubyRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "RB_YAML_LOAD" {
			rubyRules = append(rubyRules, r)
		}
	}

	if len(rubyRules) == 0 {
		t.Skip("No yaml_load rules found for Ruby")
	}

	parser := &RubyParser{}
	findings, err := parser.Parse(testFile, content, rubyRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find unsafe YAML load, found none")
	}
}

func TestRubyParser_Extensions(t *testing.T) {
	parser := &RubyParser{}
	exts := parser.Extensions()

	if len(exts) != 1 || exts[0] != ".rb" {
		t.Errorf("Extensions() returned %v, want [.rb]", exts)
	}
}

func TestRubyParser_Language(t *testing.T) {
	parser := &RubyParser{}
	if lang := parser.Language(); lang != "ruby" {
		t.Errorf("Language() returned %q, want %q", lang, "ruby")
	}
}
