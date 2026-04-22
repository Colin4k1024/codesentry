package ruby

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&RubyParser{})
}

type RubyParser struct {
	parserpkg.BaseRegexParser
}

func (p *RubyParser) Language() string { return "ruby" }

func (p *RubyParser) Extensions() []string {
	return []string{".rb"}
}

func (p *RubyParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
