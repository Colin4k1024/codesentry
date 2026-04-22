package abcoder

import (
	"testing"

	"github.com/Colin4k1024/codesentry/internal/types"
)

func TestIssueWithContext_ToIssue(t *testing.T) {
	originalIssue := types.Issue{
		RuleID:   "TEST001",
		Title:    "Test Issue",
		Severity: "WARNING",
		File:     "test.go",
		Line:     10,
		Column:   5,
		Message:  "Test message",
	}

	fix := &FixSuggestion{
		Before:      "hardcoded_password = \"secret\"",
		After:       "password = os.Getenv(\"PASSWORD\")",
		Explanation: "Use environment variable instead of hardcoding secrets",
		Confidence:  0.95,
	}

	ctx := &CodeContext{
		FunctionName: "authenticate",
	}

	issueWithCtx := &IssueWithContext{
		Issue:   originalIssue,
		Context: ctx,
		Fix:     fix,
	}

	result := issueWithCtx.ToIssue()

	// Check that the suggestion is populated from fix
	if result.Suggestion != fix.Explanation {
		t.Errorf("expected suggestion %q, got %q", fix.Explanation, result.Suggestion)
	}

	// Check that other fields are preserved
	if result.RuleID != originalIssue.RuleID {
		t.Errorf("expected RuleID %q, got %q", originalIssue.RuleID, result.RuleID)
	}
	if result.Title != originalIssue.Title {
		t.Errorf("expected Title %q, got %q", originalIssue.Title, result.Title)
	}
	if result.Severity != originalIssue.Severity {
		t.Errorf("expected Severity %q, got %q", originalIssue.Severity, result.Severity)
	}
}

func TestIssueWithContext_ToIssue_NilFix(t *testing.T) {
	originalIssue := types.Issue{
		RuleID:   "TEST002",
		Title:    "Another Issue",
		Severity: "INFO",
		File:     "main.go",
		Line:     1,
		Column:   1,
		Message:  "Another message",
	}

	issueWithCtx := &IssueWithContext{
		Issue: originalIssue,
		Fix:   nil, // No fix suggestion
	}

	result := issueWithCtx.ToIssue()

	// Check that suggestion is empty when fix is nil
	if result.Suggestion != "" {
		t.Errorf("expected empty suggestion when fix is nil, got %q", result.Suggestion)
	}

	// Check that other fields are preserved
	if result.RuleID != originalIssue.RuleID {
		t.Errorf("expected RuleID %q, got %q", originalIssue.RuleID, result.RuleID)
	}
}

func TestIssueWithContext_ToIssue_NoFix(t *testing.T) {
	originalIssue := types.Issue{
		RuleID:   "TEST003",
		Title:    "Yet Another Issue",
		Severity: "SEVERE",
		File:     "util.go",
		Line:     50,
		Column:   10,
		Message:  "Yet another message",
	}

	issueWithCtx := &IssueWithContext{
		Issue: originalIssue,
		Fix:   nil, // No fix suggestion
	}

	result := issueWithCtx.ToIssue()

	// Check that suggestion is empty when fix is nil
	if result.Suggestion != "" {
		t.Errorf("expected empty suggestion when fix is nil, got %q", result.Suggestion)
	}

	// Check that other fields are preserved
	if result.RuleID != originalIssue.RuleID {
		t.Errorf("expected RuleID %q, got %q", originalIssue.RuleID, result.RuleID)
	}
}
