package javascript

import (
	"regexp"
	"strings"

	parserpkg "github.com/goreview/goreview/internal/parser"
	"github.com/goreview/goreview/internal/rules"
)

func init() {
	parserpkg.Register(&JSParser{})
}

type JSParser struct{}

func (p *JSParser) Language() string { return "javascript" }

func (p *JSParser) Extensions() []string {
	return []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs"}
}

func (p *JSParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
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
