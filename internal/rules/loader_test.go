package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRules(t *testing.T) {
	// Use the actual rules directory
	dir := filepath.Join("..", "..", "rules")

	rules := LoadRules(dir)

	if len(rules) == 0 {
		t.Error("LoadRules returned empty slice, expected rules")
	}

	// Verify each rule has required fields
	for _, rule := range rules {
		if rule.ID == "" {
			t.Error("rule has empty ID")
		}
		if rule.Name == "" {
			t.Errorf("rule %s has empty Name", rule.ID)
		}
		if rule.Severity == "" {
			t.Errorf("rule %s has empty Severity", rule.ID)
		}
		if rule.Category == "" {
			t.Errorf("rule %s has empty Category", rule.ID)
		}
		if len(rule.Languages) == 0 {
			t.Errorf("rule %s has no languages", rule.ID)
		}
	}
}

func TestLoadRules_FileByFile(t *testing.T) {
	// Test that individual rule files can be loaded
	tests := []struct {
		name string
		want string // partial content to verify
	}{
		{"hardcoded_secret", "HARDCODED_SECRET"},
		{"sql_injection", "SQL_INJECTION"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join("..", "..", "rules", "security")
			rules := LoadRules(dir)

			found := false
			for _, rule := range rules {
				if rule.ID == tt.want {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("LoadRules did not find rule %q", tt.want)
			}
		})
	}
}

func TestFilterByLanguage(t *testing.T) {
	rules := []Rule{
		{ID: "A", Languages: []string{"go", "python"}},
		{ID: "B", Languages: []string{"java"}},
		{ID: "C", Languages: []string{"go"}},
	}

	tests := []struct {
		lang     string
		wantIDs  []string
		wantCount int
	}{
		{"go", []string{"A", "C"}, 2},
		{"python", []string{"A"}, 1},
		{"java", []string{"B"}, 1},
		{"unknown", []string{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			filtered := FilterByLanguage(tt.lang, rules)

			if len(filtered) != tt.wantCount {
				t.Errorf("FilterByLanguage(%q) returned %d rules, want %d", tt.lang, len(filtered), tt.wantCount)
			}

			gotIDs := make([]string, len(filtered))
			for i, r := range filtered {
				gotIDs[i] = r.ID
			}

			for _, want := range tt.wantIDs {
				found := false
				for _, got := range gotIDs {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("FilterByLanguage(%q) missing rule %q", tt.lang, want)
				}
			}
		})
	}
}

func TestLoadRules_InvalidYAML(t *testing.T) {
	// Create temp directory with invalid YAML
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(invalidFile, []byte("invalid: [yaml: content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not panic, just skip invalid file
	rules := LoadRules(tmpDir)

	if len(rules) != 0 {
		t.Errorf("LoadRules should skip invalid file, got %d rules", len(rules))
	}
}

func TestLoadRules_NoID(t *testing.T) {
	// Create temp directory with rule without ID
	tmpDir := t.TempDir()
	noIDFile := filepath.Join(tmpDir, "noid.yaml")
	// YAML without ID field
	if err := os.WriteFile(noIDFile, []byte("name: Test\nseverity: WARNING\n"), 0644); err != nil {
		t.Fatal(err)
	}

	rules := LoadRules(tmpDir)

	if len(rules) != 0 {
		t.Errorf("LoadRules should skip rule without ID, got %d rules", len(rules))
	}
}
