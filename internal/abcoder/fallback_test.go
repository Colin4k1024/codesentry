package abcoder

import (
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestNewFallbackHandler(t *testing.T) {
	handler := NewFallbackHandler()
	if handler == nil {
		t.Fatal("NewFallbackHandler returned nil")
	}
	if handler.ruleFixes == nil {
		t.Error("ruleFixes is nil")
	}
	if len(handler.ruleFixes) == 0 {
		t.Error("ruleFixes is empty")
	}
}

func TestFallbackHandler_GetFix(t *testing.T) {
	handler := NewFallbackHandler()

	tests := []struct {
		ruleID   string
		expected bool
	}{
		{"SQL_INJECTION", true},
		{"HARDCODED_SECRET", true},
		{"EXECUTION", true},
		{"XSS", true},
		{"DANGEROUS_EVAL", true},
		{"PATH_TRAVERSAL", true},
		{"DESERIALIZATION", true},
		{"UNKNOWN_RULE", false},
	}

	for _, tt := range tests {
		t.Run(tt.ruleID, func(t *testing.T) {
			fix := handler.GetFix(tt.ruleID)
			if tt.expected && fix == nil {
				t.Errorf("GetFix(%q) returned nil, want non-nil", tt.ruleID)
			}
			if !tt.expected && fix != nil {
				t.Errorf("GetFix(%q) returned non-nil, want nil", tt.ruleID)
			}
		})
	}
}

func TestFallbackHandler_GetFix_Values(t *testing.T) {
	handler := NewFallbackHandler()

	fix := handler.GetFix("SQL_INJECTION")
	if fix == nil {
		t.Fatal("GetFix(SQL_INJECTION) returned nil")
	}

	if fix.Pattern == "" {
		t.Error("SQL_INJECTION Pattern is empty")
	}
	if fix.Template == "" {
		t.Error("SQL_INJECTION Template is empty")
	}
}

func TestIsFallbackNeeded(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		err      error
		expected bool
	}{
		{"Go file, no error", "test.go", nil, false},
		{"Non-Go file, no error", "test.py", nil, true},
		{"Go file, with error", "test.go", errSomeError, true},
		{"Non-Go file, with error", "test.py", errSomeError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFallbackNeeded(tt.file, tt.err)
			if result != tt.expected {
				t.Errorf("IsFallbackNeeded(%q, %v) = %v, want %v", tt.file, tt.err, result, tt.expected)
			}
		})
	}
}

var errSomeError = &testError{msg: "test error"}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestFallbackHandler_BuildSuggestion(t *testing.T) {
	handler := NewFallbackHandler()

	tests := []struct {
		name       string
		rule       *rules.Rule
		wantPrefix string
	}{
		{
			name:       "nil rule returns default",
			rule:       nil,
			wantPrefix: "Review and fix",
		},
		{
			name: "rule with suggestion",
			rule: &rules.Rule{
				ID:         "CUSTOM_RULE",
				Name:       "Custom Rule",
				Suggestion: "Custom suggestion",
				Category:   "security",
			},
			wantPrefix: "Custom suggestion",
		},
		{
			name: "rule with known fallback fix",
			rule: &rules.Rule{
				ID:       "SQL_INJECTION",
				Name:     "SQL Injection",
				Category: "security",
			},
			wantPrefix: "Use parameterized queries",
		},
		{
			name: "rule with category security no suggestion",
			rule: &rules.Rule{
				ID:       "UNKNOWN_SECURITY",
				Name:     "Unknown Security Issue",
				Category: "security",
			},
			wantPrefix: "Security issue detected",
		},
		{
			name: "rule with category performance no suggestion",
			rule: &rules.Rule{
				ID:       "UNKNOWN_PERF",
				Name:     "Unknown Performance Issue",
				Category: "performance",
			},
			wantPrefix: "Performance issue detected",
		},
		{
			name: "rule with unknown category",
			rule: &rules.Rule{
				ID:       "UNKNOWN_OTHER",
				Name:     "Unknown Issue",
				Category: "other",
			},
			wantPrefix: "Code issue detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.BuildSuggestion(tt.rule)
			if result == "" {
				t.Fatal("BuildSuggestion returned empty string")
			}
			if !containsSubstr(result, tt.wantPrefix) {
				t.Errorf("BuildSuggestion() = %q, want to contain %q", result, tt.wantPrefix)
			}
		})
	}
}

func TestFallbackHandler_FormatFallback(t *testing.T) {
	handler := NewFallbackHandler()

	tests := []struct {
		name       string
		ruleID     string
		suggestion string
		wantParts  []string
	}{
		{
			name:       "known rule with suggestion",
			ruleID:     "SQL_INJECTION",
			suggestion: "Additional info",
			wantParts:  []string{"修复建议", "String concatenation in SQL query", "parameterized"},
		},
		{
			name:       "known rule with same suggestion as template",
			ruleID:     "HARDCODED_SECRET",
			suggestion: "Use environment variables: os.Getenv(\"KEY\")",
			wantParts:  []string{"修复建议", "Hardcoded sensitive value"},
		},
		{
			name:       "unknown rule returns original suggestion",
			ruleID:     "UNKNOWN_RULE_XYZ",
			suggestion: "Custom suggestion",
			wantParts:  []string{"Custom suggestion"},
		},
		{
			name:       "unknown rule with empty suggestion",
			ruleID:     "UNKNOWN_RULE_EMPTY",
			suggestion: "",
			wantParts:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.FormatFallback(tt.ruleID, tt.suggestion)
			if result == "" && len(tt.wantParts) > 0 && tt.wantParts[0] != "" {
				t.Error("FormatFallback returned empty string")
			}
			for _, part := range tt.wantParts {
				if part != "" && !containsSubstr(result, part) {
					t.Errorf("FormatFallback(%q, %q) = %q, want to contain %q", tt.ruleID, tt.suggestion, result, part)
				}
			}
		})
	}
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
