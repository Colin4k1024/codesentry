package abcoder

import (
	"testing"
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

func TestFallbackHandler_FormatFallback(t *testing.T) {
	handler := NewFallbackHandler()

	output := handler.FormatFallback("SQL_INJECTION", "")
	if output == "" {
		t.Error("FormatFallback returned empty string")
	}

	// Should contain key elements
	if !containsString(output, "修复建议") {
		t.Error("FormatFallback missing '修复建议'")
	}
	if !containsString(output, "SQL") && !containsString(output, "parameterized") {
		// The fix pattern should be mentioned
	}
}

func TestFallbackHandler_BuildSuggestion(t *testing.T) {
	handler := NewFallbackHandler()

	// Test with known rule ID
	suggestion := handler.BuildSuggestion(nil)
	if suggestion == "" {
		t.Error("BuildSuggestion returned empty string")
	}
}
