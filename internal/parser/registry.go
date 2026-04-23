package parser

import "github.com/Colin4k1024/codesentry/internal/rules"

// Parser is implemented by language-specific parsers
type Parser interface {
	Language() string
	Extensions() []string
	Parse(filePath string, content []byte, langRules []rules.Rule) ([]Finding, error)
}

// Finding represents a single issue found during parsing
type Finding struct {
	RuleID   string
	Line     int
	Column   int
	EndLine  int
	Message  string
	Severity string
}

var registry = make(map[string]Parser)

// Register adds a parser to the registry
func Register(p Parser) { registry[p.Language()] = p }

// Get returns a parser for the given language
func Get(lang string) (Parser, bool) {
	p, ok := registry[lang]
	return p, ok
}

// List returns all registered languages
func List() []string {
	var langs []string
	for lang := range registry {
		langs = append(langs, lang)
	}
	return langs
}

// DetectFromPath returns a parser based on file extension
func DetectFromPath(path string) Parser {
	for _, p := range registry {
		for _, ext := range p.Extensions() {
			if hasExt(path, ext) {
				return p
			}
		}
	}
	return nil
}

func hasExt(path, ext string) bool {
	if len(path) < len(ext) {
		return false
	}
	return path[len(path)-len(ext):] == ext
}
