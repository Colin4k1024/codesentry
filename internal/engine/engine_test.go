package engine

import (
	"os"
	"path/filepath"
	"testing"

	// Import language parsers to trigger their init() functions
	_ "github.com/Colin4k1024/codesentry/langs/cpp"
	_ "github.com/Colin4k1024/codesentry/langs/golang"
	_ "github.com/Colin4k1024/codesentry/langs/java"
	_ "github.com/Colin4k1024/codesentry/langs/javascript"
	_ "github.com/Colin4k1024/codesentry/langs/kotlin"
	_ "github.com/Colin4k1024/codesentry/langs/php"
	_ "github.com/Colin4k1024/codesentry/langs/python"
	_ "github.com/Colin4k1024/codesentry/langs/ruby"
	_ "github.com/Colin4k1024/codesentry/langs/rust"
	_ "github.com/Colin4k1024/codesentry/langs/swift"
	_ "github.com/Colin4k1024/codesentry/langs/typescript"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestEngine_Scan(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Create a Go file with hardcoded secret
	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte(`password = "hardcoded123"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a Python file
	pyFile := filepath.Join(tmpDir, "test.py")
	if err := os.WriteFile(pyFile, []byte(`token = "secret456"`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go", "python"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)(password|token)\\s*[=:]\\s*[\"'][^\"']{8,}[\"']", Comment: "Possible hardcoded secret"},
			},
		},
	})

	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	if result.FilesScanned != 2 {
		t.Errorf("FilesScanned = %d, want 2", result.FilesScanned)
	}

	if result.TotalIssues != 2 {
		t.Errorf("TotalIssues = %d, want 2", result.TotalIssues)
	}

	if result.Severe != 2 {
		t.Errorf("Severe = %d, want 2", result.Severe)
	}
}

func TestEngine_Scan_SecurityFilter(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte(`password = "hardcoded123"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create rules: one security, one performance
	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
		{
			ID:       "GOROUTINE_LEAK",
			Name:     "Goroutine Leak",
			Severity: "WARNING",
			Category: "performance",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "go\\s+func\\s*\\(", Comment: "goroutine"},
			},
		},
	})

	// Security filter should only return security rules
	cfg := &Config{Security: true}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should find the security rule (password)
	if result.Severe != 1 {
		t.Errorf("Severe = %d, want 1", result.Severe)
	}

	// Should NOT find the performance rule (go func)
	if result.Warning != 0 {
		t.Errorf("Warning = %d, want 0 (goroutine rule should be filtered)", result.Warning)
	}
}

func TestEngine_Scan_PerformanceFilter(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	// File with both a security issue (password) and a performance issue (goroutine)
	if err := os.WriteFile(goFile, []byte(`package main
password = "hardcoded123"
func main() {
	go func() {}
}
`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
		{
			ID:       "GOROUTINE_LEAK",
			Name:     "Goroutine Leak",
			Severity: "WARNING",
			Category: "performance",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "go\\s+func\\s*\\(", Comment: "goroutine"},
			},
		},
	})

	// Performance filter should only return performance rules
	cfg := &Config{Performance: true}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should NOT find the security rule
	if result.Severe != 0 {
		t.Errorf("Severe = %d, want 0 (security rule should be filtered)", result.Severe)
	}

	// Should find the performance rule (go func)
	if result.Warning != 1 {
		t.Errorf("Warning = %d, want 1", result.Warning)
	}
}

func TestEngine_Scan_NoFilter(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	// File with both security and performance issues
	if err := os.WriteFile(goFile, []byte(`package main
password = "secret"
func main() {
	go func() {}
}
`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
		{
			ID:       "GOROUTINE_LEAK",
			Name:     "Goroutine Leak",
			Severity: "WARNING",
			Category: "performance",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "go\\s+func\\s*\\(", Comment: "goroutine"},
			},
		},
	})

	// No filter - should return all rules
	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	if result.Severe != 1 {
		t.Errorf("Severe = %d, want 1", result.Severe)
	}

	if result.Warning != 1 {
		t.Errorf("Warning = %d, want 1", result.Warning)
	}

	if result.TotalIssues != 2 {
		t.Errorf("TotalIssues = %d, want 2", result.TotalIssues)
	}
}

func TestEngine_Scan_SkipDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create node_modules with a .js file
	nodeModules := filepath.Join(tmpDir, "node_modules")
	if err := os.MkdirAll(nodeModules, 0755); err != nil {
		t.Fatal(err)
	}
	jsFile := filepath.Join(nodeModules, "evil.js")
	if err := os.WriteFile(jsFile, []byte(`password = "should_not_find"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .git with a .go file
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	goFile := filepath.Join(gitDir, "config.go")
	if err := os.WriteFile(goFile, []byte(`password = "should_not_find"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create vendor with a .py file
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}
	pyFile := filepath.Join(vendorDir, "evil.py")
	if err := os.WriteFile(pyFile, []byte(`password = "should_not_find"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a real source file
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	realFile := filepath.Join(srcDir, "main.go")
	if err := os.WriteFile(realFile, []byte(`password = "should_find"`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go", "python", "javascript"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	})

	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should only find the real source file, not node_modules/.git/vendor
	if result.FilesScanned != 1 {
		t.Errorf("FilesScanned = %d, want 1 (should skip node_modules, .git, vendor)", result.FilesScanned)
	}

	if result.TotalIssues != 1 {
		t.Errorf("TotalIssues = %d, want 1", result.TotalIssues)
	}
}

func TestEngine_Scan_ExcludePattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file in a directory that should be excluded
	excludeDir := filepath.Join(tmpDir, "internal")
	if err := os.MkdirAll(excludeDir, 0755); err != nil {
		t.Fatal(err)
	}
	excludeFile := filepath.Join(excludeDir, "secret.go")
	if err := os.WriteFile(excludeFile, []byte(`password = "should_not_find"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a real source file
	realFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(realFile, []byte(`password = "should_find"`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	})

	cfg := &Config{Exclude: []string{"internal"}}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should exclude the internal directory
	if result.FilesScanned != 1 {
		t.Errorf("FilesScanned = %d, want 1 (should exclude 'internal' dir)", result.FilesScanned)
	}

	if result.TotalIssues != 1 {
		t.Errorf("TotalIssues = %d, want 1", result.TotalIssues)
	}
}

func TestEngine_Scan_Dedup(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with the same issue on the same line (dedup should keep only one)
	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte(`password = "hardcoded123"
token = "hardcoded456"
`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)(password|token)", Comment: "hardcoded secret"},
			},
		},
	})

	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should find 2 issues (one for password, one for token) on different lines
	if result.TotalIssues != 2 {
		t.Errorf("TotalIssues = %d, want 2 (different lines)", result.TotalIssues)
	}
}

func TestEngine_Scan_UnknownExtension(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with unknown extension
	unknownFile := filepath.Join(tmpDir, "test.xyz")
	if err := os.WriteFile(unknownFile, []byte(`password = "hardcoded123"`), 0644); err != nil {
		t.Fatal(err)
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	})

	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should not find any issues (unknown extension)
	if result.TotalIssues != 0 {
		t.Errorf("TotalIssues = %d, want 0 (unknown extension)", result.TotalIssues)
	}
}

func TestEngine_Scan_ReadError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a symlink to a non-existent file
	brokenLink := filepath.Join(tmpDir, "broken.link")
	if err := os.Symlink("/nonexistent/file.go", brokenLink); err != nil {
		// Skip test if symlinks not supported
		t.Skip("symlinks not supported")
	}

	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	})

	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() should not return error for unreadable file, got: %v", err)
	}

	// Should handle read error gracefully
	_ = result
}

func TestEngine_Scan_NoRules(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte(`password = "secret"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create engine with no rules
	e := New([]rules.Rule{})
	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should scan file but find no issues
	if result.FilesScanned != 1 {
		t.Errorf("FilesScanned = %d, want 1", result.FilesScanned)
	}
	if result.TotalIssues != 0 {
		t.Errorf("TotalIssues = %d, want 0", result.TotalIssues)
	}
}

func TestEngine_Scan_NoMatchingRule(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte(`password = "secret"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create rule for Python only
	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			Severity: "SEVERE",
			Category: "security",
			Languages: []string{"python"}, // Go file should not match
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	})
	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Should not find issues (rule is for Python, not Go)
	if result.TotalIssues != 0 {
		t.Errorf("TotalIssues = %d, want 0 (no matching language)", result.TotalIssues)
	}
}

func TestEngine_Scan_EmptySeverity(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte(`password = "secret"`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create rule without severity
	e := New([]rules.Rule{
		{
			ID:       "HARDCODED_SECRET",
			Name:     "Hardcoded Secret",
			// Severity is empty - should default to WARNING
			Category: "security",
			Languages: []string{"go"},
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	})
	cfg := &Config{}
	result, err := e.Scan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// Severity should default to WARNING
	if result.Warning != 1 {
		t.Errorf("Warning = %d, want 1 (severity should default to WARNING)", result.Warning)
	}
}

func TestIsCodeFile(t *testing.T) {
	tests := []struct {
		ext    string
		isCode bool
	}{
		{".go", true},
		{".js", true},
		{".py", true},
		{".java", true},
		{".rb", true},
		{".rs", true},
		{".cpp", true},
		{".c", true},
		{".h", true},
		{".php", true},
		{".swift", true},
		{".kt", true},
		{".ts", true},
		{".tsx", true},
		{".txt", false},
		{".md", false},
		{".json", false},
		{".xml", false},
		{".yaml", false},
		{".yml", false},
		{"", false},
	}

	for _, tt := range tests {
		result := isCodeFile(tt.ext)
		if result != tt.isCode {
			t.Errorf("isCodeFile(%q) = %v, want %v", tt.ext, result, tt.isCode)
		}
	}
}

func TestEngine_New(t *testing.T) {
	rules := []rules.Rule{
		{
			ID:   "TEST_RULE",
			Name: "Test Rule",
		},
	}
	e := New(rules)
	if e == nil {
		t.Fatal("New() returned nil")
	}
	if len(e.rules) != 1 {
		t.Errorf("len(e.rules) = %d, want 1", len(e.rules))
	}
}
