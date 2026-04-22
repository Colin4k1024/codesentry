package typescript

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&TSParser{})
}

type TSParser struct {
	parserpkg.BaseRegexParser
}

func (p *TSParser) Language() string { return "typescript" }

func (p *TSParser) Extensions() []string {
	return []string{".ts", ".tsx", ".mts", ".cts"}
}

func (p *TSParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
