package php

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&PHPParser{})
}

type PHPParser struct {
	parserpkg.BaseRegexParser
}

func (p *PHPParser) Language() string { return "php" }

func (p *PHPParser) Extensions() []string {
	return []string{".php"}
}

func (p *PHPParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
