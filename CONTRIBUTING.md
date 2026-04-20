# Contributing to CodeSentry

Thank you for your interest in contributing!

## Ways to Contribute

- **Add new language parsers** - Implement a parser for a language not yet supported
- **Improve existing rules** - Add more accurate regex patterns or AST-based checks
- **Add new rules** - Create new security or performance rules
- **Bug fixes** - Fix incorrect detections or false positives
- **Documentation** - Improve docs, add examples, translate

## Development Setup

```bash
# Clone the repository
git clone https://github.com/Colin4k1024/codesentry_refactor.git
cd codesentry_refactor

# Build
go build -o codesentry ./cmd/codesentry

# Run tests
go test ./...
```

## Adding a New Language

### 1. Create the Parser

Create `langs/<language>/parser.go`:

```go
package <language>

import (
    "regexp"
    "strings"

    parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
    "github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
    parserpkg.Register(&Parser{})
}

type Parser struct{}

func (p *Parser) Language() string { return "<language>" }
func (p *Parser) Extensions() []string { return []string{".ext"} }

func (p *Parser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
    var findings []parserpkg.Finding
    text := string(content)
    lines := strings.Split(text, "\n")

    for _, rule := range langRules {
        for _, pattern := range rule.Patterns {
            if pattern.Type != "regex" {
                continue
            }
            re, err := regexp.Compile(pattern.Pattern)
            if err != nil {
                continue
            }

            for lineNum, line := range lines {
                if re.MatchString(line) {
                    findings = append(findings, parserpkg.Finding{
                        RuleID:   rule.ID,
                        Line:     lineNum + 1,
                        Column:   1,
                        EndLine:  lineNum + 1,
                        Severity: rule.Severity,
                        Message:  pattern.Comment,
                    })
                }
            }
        }
    }
    return findings, nil
}
```

### 2. Register the Parser

Add to `cmd/codesentry/langs.go`:
```go
import _ "github.com/Colin4k1024/codesentry/langs/<language>"
```

### 3. Add Rules

Create `rules/<language>/` directory with YAML rule files.

Example `rules/<language>/hardcoded_secret.yaml`:
```yaml
id: LANG_HARDCODED_SECRET
name: Hardcoded Secret
description: Detects hardcoded secrets in <language>
severity: SEVERE
category: security
languages:
  - <language>
suggestion: Use environment variables
patterns:
  - type: regex
    pattern: '(?i)(password|token|secret)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
    comment: Possible hardcoded secret
```

### 4. Test

```bash
go build -o codesentry ./cmd/codesentry
./codesentry languages  # Should list your new language
```

## Rule YAML Schema

```yaml
id: RULE_ID              # Unique identifier (SCREAMING_SNAKE_CASE)
name: Human Name         # Display name
description: Description # What the rule detects
severity: SEVERE|WARNING|INFO
category: security|performance
languages:               # List of supported languages
  - go
  - python
suggestion: How to fix  # Optional fix suggestion
patterns:                # Detection patterns
  - type: regex         # Currently only regex is supported
    pattern: 'regex'    # Regular expression
    comment: Description of the match
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/language-support`)
3. Make your changes
4. Build and test (`go build && ./codesentry scan ./...`)
5. Commit with clear messages
6. Push and open a PR

## Code Style

- Run `go fmt ./...` before committing
- Follow standard Go conventions
- Add tests for new functionality
