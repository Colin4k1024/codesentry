package swift

import (
	"regexp"
	"strings"

	parserpkg "github.com/goreview/goreview/internal/parser"
	"github.com/goreview/goreview/internal/rules"
)

func init() {
	parserpkg.Register(&SwiftParser{})
}

type SwiftParser struct{}

func (p *SwiftParser) Language() string { return "swift" }

func (p *SwiftParser) Extensions() []string {
	return []string{".swift"}
}

func (p *SwiftParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	var findings []parserpkg.Finding
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
					findings = append(findings, parserpkg.Finding{
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

	return findings, nil
}
