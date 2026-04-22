package cpp

import (
	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&CppParser{})
}

type CppParser struct {
	parserpkg.BaseRegexParser
}

func (p *CppParser) Language() string { return "cpp" }

func (p *CppParser) Extensions() []string {
	return []string{".cpp", ".cc", ".cxx", ".c++", ".h", ".hpp"}
}

func (p *CppParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
