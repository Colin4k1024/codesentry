package cpp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestCppParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.cpp")

	// C++ code with hardcoded credentials
	code := `std::string password = "hardcoded123";
std::string api_key = "sk-abcdefgh123456";
const char* token = "Bearer xxx.yyy.zzz";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load C++ rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var cppRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "CPP_HARDCODED_SECRET" {
			cppRules = append(cppRules, r)
		}
	}

	if len(cppRules) == 0 {
		t.Skip("No hardcoded_secret rules found for C++")
	}

	parser := &CppParser{}
	findings, err := parser.Parse(testFile, content, cppRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestCppParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.cpp")

	// C++ code with SQL injection vulnerability - must use query/execute/prepare with +
	code := `stmt->execute("SELECT * FROM users WHERE id=" + userId);
query("DELETE FROM users WHERE id=" + userId);`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load C++ rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var cppRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "CPP_SQL_INJECTION" {
			cppRules = append(cppRules, r)
		}
	}

	if len(cppRules) == 0 {
		t.Skip("No SQL injection rules found for C++")
	}

	parser := &CppParser{}
	findings, err := parser.Parse(testFile, content, cppRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestCppParser_Extensions(t *testing.T) {
	parser := &CppParser{}
	exts := parser.Extensions()

	expected := []string{".cpp", ".cc", ".cxx", ".c++", ".h", ".hpp"}
	if len(exts) != len(expected) {
		t.Errorf("Extensions() returned %d extensions, want %d", len(exts), len(expected))
	}
}

func TestCppParser_Language(t *testing.T) {
	parser := &CppParser{}
	if lang := parser.Language(); lang != "cpp" {
		t.Errorf("Language() returned %q, want %q", lang, "cpp")
	}
}
