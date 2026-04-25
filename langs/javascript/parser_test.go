package javascript

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func testRule(rulesDir, ruleID string) []rules.Rule {
	allRules := rules.LoadRules(rulesDir)
	var result []rules.Rule
	for _, r := range allRules {
		for _, lang := range r.Languages {
			if lang == "javascript" && r.ID == ruleID {
				result = append(result, r)
				break
			}
		}
	}
	return result
}

func TestJSParser_HARDCODED_SECRET(t *testing.T) {
	rulesDir := filepath.Join("..", "..", "rules")
	jsRules := testRule(rulesDir, "TS_HARDCODED_SECRET")
	if len(jsRules) == 0 {
		t.Skip("No TS_HARDCODED_SECRET rule found for JavaScript")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.js")
	code := `const password = "hardcoded12345678";
const api_key = "sk_test_abcdef1234567890";
const token = "ghp_abcdefghijklmnopqrstuvwxyz1234567890";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	parser := &JSParser{}
	findings, err := parser.Parse(testFile, content, jsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestJSParser_SQL_INJECTION(t *testing.T) {
	rulesDir := filepath.Join("..", "..", "rules")
	jsRules := testRule(rulesDir, "TS_SQL_INJECTION")
	if len(jsRules) == 0 {
		t.Skip("No TS_SQL_INJECTION rule found for JavaScript")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.js")
	code := `db.query("SELECT * FROM users WHERE id=" + userId);
result = await conn.query("SELECT * FROM orders WHERE user_id=" + uid + " AND status='" + status + "'");`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	parser := &JSParser{}
	findings, err := parser.Parse(testFile, content, jsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestJSParser_SSRF(t *testing.T) {
	rulesDir := filepath.Join("..", "..", "rules")
	jsRules := testRule(rulesDir, "TS_SSRF")
	if len(jsRules) == 0 {
		t.Skip("No TS_SSRF rule found for JavaScript")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.js")
	code := "fetch(userInputUrl + \"/api/data\");\n" +
		"const response = await axios.get(baseUrl + \"/\" + userPath);\n" +
		"fetch(\"https://example.com/users/\" + userId + \"/profile\");"
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	parser := &JSParser{}
	findings, err := parser.Parse(testFile, content, jsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SSRF vulnerabilities, found none")
	}
}

func TestJSParser_EVAL(t *testing.T) {
	rulesDir := filepath.Join("..", "..", "rules")
	jsRules := testRule(rulesDir, "TS_EVAL")
	if len(jsRules) == 0 {
		t.Skip("No TS_EVAL rule found for JavaScript")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.js")
	code := `eval(userInput);
const fn = new Function("x", "return x + " + userExpr);
element.innerHTML = "<b>" + userContent + "</b>";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	parser := &JSParser{}
	findings, err := parser.Parse(testFile, content, jsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find eval vulnerabilities, found none")
	}
}

func TestJSParser_CleanNoFindings(t *testing.T) {
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var jsRules []rules.Rule
	for _, r := range allRules {
		for _, lang := range r.Languages {
			if lang == "javascript" {
				jsRules = append(jsRules, r)
				break
			}
		}
	}
	if len(jsRules) == 0 {
		t.Skip("No rules found for JavaScript")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "clean.js")
	code := `const x = 42;
function greet(name) { return "Hello, " + name; }
console.log(greet("World"));`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	parser := &JSParser{}
	findings, err := parser.Parse(testFile, content, jsRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) > 0 {
		t.Errorf("expected no findings for clean code, got %d", len(findings))
	}
}
