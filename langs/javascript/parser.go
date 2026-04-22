package javascript

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&JSParser{})
}

type JSParser struct {
	parserpkg.BaseRegexParser
}

func (p *JSParser) Language() string { return "javascript" }

func (p *JSParser) Extensions() []string {
	return []string{".js", ".jsx", ".mjs", ".cjs"}
}

func (p *JSParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
