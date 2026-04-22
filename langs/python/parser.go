package python

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&PythonParser{})
}

type PythonParser struct {
	parserpkg.BaseRegexParser
}

func (p *PythonParser) Language() string { return "python" }

func (p *PythonParser) Extensions() []string {
	return []string{".py", ".pyw", ".pyi"}
}

func (p *PythonParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
