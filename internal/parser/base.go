package parser

import (
	"regexp"
	"strings"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

// BaseRegexParser provides common regex-based parsing logic for language parsers.
// Embed this struct in language parsers to reuse the standard regex matching logic.
// For AST-based checks, embed BaseRegexParser and implement additional methods.
type BaseRegexParser struct{}

// ParseRegex implements the standard regex-based rule matching.
// Call this from your Parse method after initializing BaseRegexParser.
func (p *BaseRegexParser) ParseRegex(content []byte, langRules []rules.Rule) []Finding {
	var findings []Finding
	text := string(content)
	lines := strings.Split(text, "\n")

	for _, rule := range langRules {
		for _, pattern := range rule.Patterns {
			if pattern.Type != "regex" {
				continue
			}
			re, err := regexp.Compile(pattern.Pattern)
			if err != nil {
				continue
			}

			for lineNum, line := range lines {
				if re.MatchString(line) {
					findings = append(findings, Finding{
						RuleID:   rule.ID,
						Line:     lineNum + 1,
						Column:   1,
						EndLine:  lineNum + 1,
						Severity: rule.Severity,
						Message:  pattern.Comment,
					})
				}
			}
		}
	}

	return findings
}
