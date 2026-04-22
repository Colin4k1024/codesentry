package php

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestPHPParser_HARDCODED_SECRET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.php")

	// PHP code with hardcoded credentials
	code := `<?php
$password = "hardcoded123";
$api_key = "sk-abcdefgh123456";
$token = "Bearer xxx.yyy.zzz";
$config['secret'] = "my_secret_value";`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load PHP rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var phpRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PHP_HARDCODED_SECRET" {
			phpRules = append(phpRules, r)
		}
	}

	if len(phpRules) == 0 {
		t.Skip("No hardcoded_secret rules found for PHP")
	}

	parser := &PHPParser{}
	findings, err := parser.Parse(testFile, content, phpRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find hardcoded secrets, found none")
	}
}

func TestPHPParser_DESERIALIZATION(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.php")

	// PHP code with unsafe deserialization
	code := `<?php
$data = unserialize($_GET['data']);
$obj = unserialize($userInput);`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load PHP rules
	rulesDir := filepath.Join("..", "..", "rules")
	allRules := rules.LoadRules(rulesDir)
	var phpRules []rules.Rule
	for _, r := range allRules {
		if r.ID == "PHP_DESERIALIZATION" {
			phpRules = append(phpRules, r)
		}
	}

	if len(phpRules) == 0 {
		t.Skip("No deserialization rules found for PHP")
	}

	parser := &PHPParser{}
	findings, err := parser.Parse(testFile, content, phpRules)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(findings) == 0 {
		t.Error("expected to find unsafe deserialization, found none")
	}
}

func TestPHPParser_Extensions(t *testing.T) {
	parser := &PHPParser{}
	exts := parser.Extensions()

	if len(exts) != 1 || exts[0] != ".php" {
		t.Errorf("Extensions() returned %v, want [.php]", exts)
	}
}

func TestPHPParser_Language(t *testing.T) {
	parser := &PHPParser{}
	if lang := parser.Language(); lang != "php" {
		t.Errorf("Language() returned %q, want %q", lang, "php")
	}
}
