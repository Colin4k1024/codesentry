package abcoder

import (
	"context"
	"testing"
)

func TestSkillAgent_GenerateFix(t *testing.T) {
	agent := NewSkillAgent(nil)

	tests := []struct {
		name       string
		ruleID     string
		suggestion string
		file       string
		line       int
	}{
		{
			name:       "SQL Injection",
			ruleID:     "SQL_INJECTION",
			suggestion: "Use parameterized queries",
			file:       "test.go",
			line:       42,
		},
		{
			name:       "Hardcoded Secret",
			ruleID:     "HARDCODED_SECRET",
			suggestion: "Use environment variables",
			file:       "config.go",
			line:       10,
		},
		{
			name:       "Command Execution",
			ruleID:     "EXECUTION",
			suggestion: "Avoid shell injection",
			file:       "run.go",
			line:       20,
		},
		{
			name:       "Unknown Rule",
			ruleID:     "UNKNOWN_RULE",
			suggestion: "Review this code",
			file:       "test.go",
			line:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := agent.GenerateFix(context.Background(), tt.ruleID, tt.suggestion, tt.file, tt.line)
			if err != nil {
				t.Fatalf("GenerateFix failed: %v", err)
			}
			if output == nil {
				t.Fatal("GenerateFix returned nil")
			}
			if output.Before == "" {
				t.Error("Before code is empty")
			}
			if output.After == "" {
				t.Error("After code is empty")
			}
			if output.Explanation == "" {
				t.Error("Explanation is empty")
			}
			if output.Confidence < 0 || output.Confidence > 1 {
				t.Errorf("Confidence %f is out of range [0, 1]", output.Confidence)
			}
		})
	}
}

func TestSkillOutput_ToJSON(t *testing.T) {
	output := &SkillOutput{
		Before:      "old code",
		After:       "new code",
		Explanation: "because it's better",
		Confidence:  0.85,
	}

	data, err := output.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("ToJSON returned empty data")
	}
}

func TestSkillOutput_FormatFix(t *testing.T) {
	output := &SkillOutput{
		Before:      "old code",
		After:       "new code",
		Explanation: "because it's better",
		Confidence:  0.85,
		Warnings:    []string{"warning 1", "warning 2"},
	}

	formatted := output.FormatFix()
	if formatted == "" {
		t.Error("FormatFix returned empty string")
	}

	// Should contain key elements
	if !contains(formatted, "修复前") {
		t.Error("FormatFix missing '修复前'")
	}
	if !contains(formatted, "修复后") {
		t.Error("FormatFix missing '修复后'")
	}
	if !contains(formatted, "原因") {
		t.Error("FormatFix missing '原因'")
	}
	if !contains(formatted, "85%") {
		t.Error("FormatFix missing confidence")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
