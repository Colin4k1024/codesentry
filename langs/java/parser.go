package java

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&JavaParser{})
}

type JavaParser struct {
	parserpkg.BaseRegexParser
}

func (p *JavaParser) Language() string { return "java" }

func (p *JavaParser) Extensions() []string {
	return []string{".java"}
}

func (p *JavaParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
