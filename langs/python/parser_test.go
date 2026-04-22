package python

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestPythonParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.py")

	// Python code with hardcoded credentials
	code := `password = "hardcoded123"
api_key = "sk-abcdefgh123456"
token = "Bearer xxx.yyy.zzz"`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Python rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var pythonRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PY_HARDCODED_SECRET" {
			pythonRules = append(pythonRules, r)
		}
	}

	if len(pythonRules) == 0 {
		t.Skip("No hardcoded_secret rules found for Python")
	}

	parser := &PythonParser{}
	findings, err := parser.Parse(testFile, content, pythonRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestPythonParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.py")

	// Python code with SQL injection vulnerability
	code := `cursor.execute("SELECT * FROM users WHERE id=" + user_id)
cursor.execute("INSERT INTO logs VALUES ('" + data + "')")`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Python rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var pythonRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PY_SQL_INJECTION" {
			pythonRules = append(pythonRules, r)
		}
	}

	if len(pythonRules) == 0 {
		t.Skip("No SQL injection rules found for Python")
	}

	parser := &PythonParser{}
	findings, err := parser.Parse(testFile, content, pythonRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestPythonParser_PICKLE(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.py")

	// Python code with unsafe pickle deserialization
	code := `data = pickle.loads(user_input)
obj = pickle.load(file_handle)`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Python rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var pythonRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PY_PICKLE" {
			pythonRules = append(pythonRules, r)
		}
	}

	if len(pythonRules) == 0 {
		t.Skip("No pickle rules found for Python")
	}

	parser := &PythonParser{}
	findings, err := parser.Parse(testFile, content, pythonRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find unsafe pickle, found none")
	}
}

func TestPythonParser_SUBPROCESS(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.py")

	// Python code with subprocess shell=True
	code := `subprocess.run("ls " + path, shell=True)
subprocess.call("rm -rf " + target, shell=True)`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Python rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var pythonRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PY_SUBPROCESS" {
			pythonRules = append(pythonRules, r)
		}
	}

	if len(pythonRules) == 0 {
		t.Skip("No subprocess rules found for Python")
	}

	parser := &PythonParser{}
	findings, err := parser.Parse(testFile, content, pythonRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find subprocess issue, found none")
	}
}

func TestPythonParser_YAML_LOAD(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.py")

	// Python code with unsafe YAML load
	code := `data = yaml.load(user_input)
data = yaml.load(input_file)`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Python rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var pythonRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PY_YAML_LOAD" {
			pythonRules = append(pythonRules, r)
		}
	}

	if len(pythonRules) == 0 {
		t.Skip("No yaml_load rules found for Python")
	}

	parser := &PythonParser{}
	findings, err := parser.Parse(testFile, content, pythonRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find unsafe YAML load, found none")
	}
}

func TestPythonParser_Extensions(t *testing.T) {
	parser := &PythonParser{}
	exts := parser.Extensions()

	if len(exts) != 3 {
		t.Errorf("Extensions() returned %d extensions, want 3", len(exts))
	}
}

func TestPythonParser_Language(t *testing.T) {
	parser := &PythonParser{}
	if lang := parser.Language(); lang != "python" {
		t.Errorf("Language() returned %q, want %q", lang, "python")
	}
}
