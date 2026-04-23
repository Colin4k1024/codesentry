package parser

import (
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

func TestBaseRegexParser_ParseRegex(t *testing.T) {
	p := &BaseRegexParser{}

	tests := []struct {
		name      string
		content   string
		rules     []rules.Rule
		wantCount int
		wantLines []int
	}{
		{
			name:      "no rules",
			content:   "password = \"secret123\"\n",
			rules:     []rules.Rule{},
			wantCount: 0,
		},
		{
			name:    "no matches",
			content: "x = 123\n",
			rules: []rules.Rule{
				{
					ID:       "TEST",
					Severity: "WARNING",
					Patterns: []rules.Pattern{
						{Type: "regex", Pattern: "(?i)password", Comment: "test"},
					},
				},
			},
			wantCount: 0,
		},
		{
			name:    "single match",
			content: "password = \"secret123\"\n",
			rules: []rules.Rule{
				{
					ID:       "TEST",
					Severity: "WARNING",
					Patterns: []rules.Pattern{
						{Type: "regex", Pattern: "(?i)password", Comment: "test"},
					},
				},
			},
			wantCount: 1,
			wantLines: []int{1},
		},
		{
			name:    "multiple matches",
			content: "password = \"secret1\"\ntoken = \"secret2\"\n",
			rules: []rules.Rule{
				{
					ID:       "TEST",
					Severity: "WARNING",
					Patterns: []rules.Pattern{
						{Type: "regex", Pattern: "(?i)(password|token)", Comment: "test"},
					},
				},
			},
			wantCount: 2,
			wantLines: []int{1, 2},
		},
		{
			name:    "skips non-regex patterns",
			content: "password = \"secret\"\n",
			rules: []rules.Rule{
				{
					ID:       "TEST",
					Severity: "WARNING",
					Patterns: []rules.Pattern{
						{Type: "ast", Pattern: "some.ast.query", Comment: "ast not supported"},
					},
				},
			},
			wantCount: 0,
		},
		{
			name:    "invalid regex skipped",
			content: "password = \"secret\"\n",
			rules: []rules.Rule{
				{
					ID:       "TEST",
					Severity: "WARNING",
					Patterns: []rules.Pattern{
						{Type: "regex", Pattern: "[invalid", Comment: "bad regex"},
					},
				},
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := p.ParseRegex([]byte(tt.content), tt.rules)

			if len(findings) != tt.wantCount {
				t.Errorf("ParseRegex() returned %d findings, want %d", len(findings), tt.wantCount)
			}

			if len(tt.wantLines) > 0 {
				for i, wantLine := range tt.wantLines {
					if i >= len(findings) {
						t.Errorf("missing finding for line %d", wantLine)
						continue
					}
					if findings[i].Line != wantLine {
						t.Errorf("finding[%d].Line = %d, want %d", i, findings[i].Line, wantLine)
					}
				}
			}
		})
	}
}

func TestBaseRegexParser_ParseRegex_Severity(t *testing.T) {
	p := &BaseRegexParser{}

	content := "password = \"secret\"\n"
	rules := []rules.Rule{
		{
			ID:       "SEVERE_TEST",
			Severity: "SEVERE",
			Patterns: []rules.Pattern{
				{Type: "regex", Pattern: "(?i)password", Comment: "hardcoded secret"},
			},
		},
	}

	findings := p.ParseRegex([]byte(content), rules)

	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}

	if findings[0].Severity != "SEVERE" {
		t.Errorf("finding.Severity = %q, want %q", findings[0].Severity, "SEVERE")
	}

	if findings[0].RuleID != "SEVERE_TEST" {
		t.Errorf("finding.RuleID = %q, want %q", findings[0].RuleID, "SEVERE_TEST")
	}

	if findings[0].Message != "hardcoded secret" {
		t.Errorf("finding.Message = %q, want %q", findings[0].Message, "hardcoded secret")
	}
}
