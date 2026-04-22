package java

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestJavaParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.java")

	// Java code with hardcoded credentials
	code := `String password = "hardcoded123";
String apiKey = "sk-abcdefgh123456";
String token = "Bearer xxx.yyy.zzz";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Java rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var javaRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "JAVA_HARDCODED_SECRET" {
			javaRules = append(javaRules, r)
		}
	}

	if len(javaRules) == 0 {
		t.Skip("No hardcoded_secret rules found for Java")
	}

	parser := &JavaParser{}
	findings, err := parser.Parse(testFile, content, javaRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestJavaParser_SQL_INJECTION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.java")

	// Java code with SQL injection vulnerability - must use executeQuery/execute with +
	code := `stmt.executeQuery("SELECT * FROM users WHERE id=" + userId);
stmt.execute("DELETE FROM users WHERE id=" + userId);`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Java rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var javaRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "JAVA_SQL_INJECTION" {
			javaRules = append(javaRules, r)
		}
	}

	if len(javaRules) == 0 {
		t.Skip("No SQL injection rules found for Java")
	}

	parser := &JavaParser{}
	findings, err := parser.Parse(testFile, content, javaRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find SQL injection, found none")
	}
}

func TestJavaParser_DESERIALIZATION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.java")

	// Java code with unsafe deserialization
	code := `ObjectInputStream ois = new ObjectInputStream(input);
Object obj = ois.readObject();`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load Java rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var javaRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "JAVA_DESERIALIZATION" {
			javaRules = append(javaRules, r)
		}
	}

	if len(javaRules) == 0 {
		t.Skip("No deserialization rules found for Java")
	}

	parser := &JavaParser{}
	findings, err := parser.Parse(testFile, content, javaRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find unsafe deserialization, found none")
	}
}

func TestJavaParser_Extensions(t *testing.T) {
	parser := &JavaParser{}
	exts := parser.Extensions()

	if len(exts) != 1 || exts[0] != ".java" {
		t.Errorf("Extensions() returned %v, want [.java]", exts)
	}
}

func TestJavaParser_Language(t *testing.T) {
	parser := &JavaParser{}
	if lang := parser.Language(); lang != "java" {
		t.Errorf("Language() returned %q, want %q", lang, "java")
	}
}
