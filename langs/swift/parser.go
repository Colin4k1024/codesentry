package swift

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&SwiftParser{})
}

type SwiftParser struct {
	parserpkg.BaseRegexParser
}

func (p *SwiftParser) Language() string { return "swift" }

func (p *SwiftParser) Extensions() []string {
	return []string{".swift"}
}

func (p *SwiftParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
