package kotlin

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&KotlinParser{})
}

type KotlinParser struct {
	parserpkg.BaseRegexParser
}

func (p *KotlinParser) Language() string { return "kotlin" }

func (p *KotlinParser) Extensions() []string {
	return []string{".kt", ".kts"}
}

func (p *KotlinParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
