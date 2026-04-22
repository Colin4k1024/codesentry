package rust

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&RustParser{})
}

type RustParser struct {
	parserpkg.BaseRegexParser
}

func (p *RustParser) Language() string { return "rust" }

func (p *RustParser) Extensions() []string {
	return []string{".rs"}
}

func (p *RustParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
