# Contributing to CodeSentry

Thank you for your interest in contributing!

## Ways to Contribute

- **Add new language parsers** — Implement a parser for a language not yet supported
- **Improve existing rules** — Add more accurate regex patterns or AST-based checks
- **Add new rules** — Create new security or performance rules
- **Bug fixes** — Fix incorrect detections or false positives
- **Documentation** — Improve docs, add examples, translate
- **Test coverage** — Add unit tests and golden file tests

---

## Development Setup

```bash
# Clone the repository
git clone https://github.com/Colin4k1024/codesentry_refactor.git
cd codesentry_refactor

# Install dependencies
go mod download

# Build
go build -o codesentry ./cmd/goreview

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

---

## Project Structure Overview

```
codesentry/
├── cmd/goreview/         # CLI application entry point
├── internal/
│   ├── engine/           # Core scanning engine
│   ├── parser/           # Parser registry and base types
│   ├── rules/            # Rule loading and types
│   ├── abcoder/          # AI context enrichment (Go only)
│   ├── output/           # Output formatting (text, JSON, SARIF)
│   └── types/            # Core data structures
├── langs/                # Language parsers (plugin pattern)
└── rules/                # YAML rule definitions
```

---

## Adding a New Language Parser

### 1. Create the Parser Package

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

type Parser struct {
    parserpkg.BaseRegexParser  // Embed for standard regex matching
}

func (p *Parser) Language() string { return "<language>" }
func (p *Parser) Extensions() []string { return []string{".ext"} }

func (p *Parser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
    // For most languages: delegate to BaseRegexParser
    return p.BaseRegexParser.ParseRegex(content, langRules), nil
}
```

> **Important**: Use `parserpkg.BaseRegexParser` embedding (composition, not inheritance). Do not shadow the receiver variable `p` with local declarations like `p := fset.Position(...)` — this causes hard-to-debug issues in Go.

### 2. Register the Parser

In `cmd/goreview/langs.go`, add a blank import:

```go
import (
    _ "github.com/Colin4k1024/codesentry/langs/<language>"
)
```

### 3. Add Rules

Create a `rules/<language>/` directory with YAML rule files. See [docs/RULES.md](docs/RULES.md) for the rule schema.

### 4. Test Your Parser

```bash
go build -o codesentry ./cmd/goreview
./codesentry languages      # Should list your new language
./codesentry scan ./testdir # Test on a sample file
```

---

## Writing Rules

### Rule YAML Schema

```yaml
id: UNIQUE_ID              # Required: Unique identifier (SCREAMING_SNAKE_CASE)
name: Human Name           # Required: Display name
description: Description  # Required: What the rule detects
severity: SEVERE|WARNING|INFO
category: security|performance
languages:                 # List of supported languages
  - go
  - python
suggestion: How to fix    # Recommended fix suggestion
patterns:                  # Detection patterns
  - type: regex           # Currently only regex is fully implemented
    pattern: 'regex'      # Regular expression
    comment: Description  # Human-readable match description
```

### Severity Levels

| Level   | Description                                                    |
|---------|----------------------------------------------------------------|
| SEVERE  | Critical security vulnerability (SQL injection, RCE, etc.)       |
| WARNING | Potential issue or code smell                                   |
| INFO    | Informational notice                                            |

### Regex Best Practices

#### 1. Use Case-Insensitive Flags Sparingly

```yaml
# Bad: Makes everything case-insensitive
pattern: '(?i)password.*=.*".*"'

# Good: Only match case-insensitive parts
pattern: '(?i)(password|token)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
```

#### 2. Escape Special Characters

In YAML, these characters have special meaning: `:`, `#`, `|`, `>`

Always quote patterns containing these:

```yaml
pattern: '(?i)\.innerHTML\s*='
```

#### 3. Anchor When Possible

```yaml
# Better: Match at word boundary
pattern: '\beval\s*\('

# Match only at start of assignment (skip comments)
pattern: '(?i)^[^/]*(password|token)\s*[=:]'
```

#### 4. Common Patterns

**Hardcoded Secrets**
```yaml
pattern: '(?i)(password|token|api_key|secret)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
```

**SQL Injection**
```yaml
pattern: '(?i)(execute|query)\s*\([^)]*\+[^)]*\)'
```

**Sensitive Logging**
```yaml
pattern: '(?i)(log|print|echo)\([^)]*(password|token|secret)[^)]*\)'
```

### Testing Your Rule

1. Create a test file with known issues:

```python
# /tmp/test_rule.py
password = "hardcoded_secret_123"
```

2. Run the scanner:

```bash
./codesentry scan /tmp/test_rule.py --security
```

3. Verify your rule triggers on the test file.

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test -v ./internal/parser/...
```

### Test Types

#### 1. Unit Tests

Standard Go tests in `*_test.go` files.

```bash
go test -v ./internal/rules/...
```

#### 2. Golden File Tests

For parser and rule testing, golden file tests compare output against reference files.

```
testdata/
└── rules/
    ├── HARDCODED_SECRET.input.go   # Input file with known issues
    └── HARDCODED_SECRET.golden.json # Expected JSON output
```

Run golden tests:
```bash
go test -v ./internal/parser/...   # Parser registry tests
go test -v ./cmd/goreview/...      # Language-specific tests
```

#### 3. End-to-End Tests

End-to-end tests verify the full scanning pipeline:

```bash
go test -v ./internal/abcoder/...
```

### Writing Tests

#### Parser Tests

```go
func TestPythonParser(t *testing.T) {
    content := []byte(`password = "hardcoded_secret_123"`)
    findings, err := parser.Parse("/test/test.py", content, testRules)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(findings) != 1 {
        t.Errorf("expected 1 finding, got %d", len(findings))
    }
}
```

#### Golden File Tests

```go
func TestGolden(t *testing.T) {
    inputPath := "testdata/rules/HARDCODED_SECRET.input.go"
    goldenPath := "testdata/rules/HARDCODED_SECRET.golden.json"

    // ... run scanner on input ...

    // Compare with golden file
    if diff := cmp.Diff(golden, actual); diff != "" {
        t.Errorf("mismatch (-golden +actual):\n%s", diff)
    }
}
```

---

## abcoder Integration (Go AI Context)

CodeSentry uses [cloudwego/abcoder](https://github.com/cloudwego/abcoder) to provide code context understanding for Go files.

### How It Works

1. When scanning Go files without `--no-ai`, the engine invokes abcoder to parse the repository
2. For each finding, abcoder retrieves context: function name, variables, call chains
3. The context enriches the suggestion, making fix recommendations more precise
4. For non-Go languages, fallback suggestions from the YAML rule are used

### abcoder Architecture

```
abcoder Bridge (internal/abcoder/bridge.go)
├── UniAST Repository  ← Parsed via cloudwego/abcoder
├── CodeContext       ← Function, variables, calls at finding location
└── SkillAgent        ← Fix suggestion generation
```

### Key Constraints

- **Go only**: abcoder UniAST currently only supports Go AST parsing
- **Repository-level**: Context requires parsing the entire repository (adds latency)
- **Fallback**: Non-Go languages use static suggestions from YAML rules

### Configuration

```bash
# Skip AI context enrichment (faster scanning)
./codesentry scan ./src --no-ai

# With AI context (default for Go files)
./codesentry scan ./src --security
```

---

## Code Style

- Run `go fmt ./...` before committing
- Follow standard Go conventions (Effective Go, Go Code Review Comments)
- Add tests for new functionality
- Keep PRs focused and small

---

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/language-support`)
3. Make your changes
4. Ensure tests pass (`go test ./...`)
5. Ensure code is formatted (`go fmt ./...`)
6. Commit with clear, descriptive messages
7. Push and open a Pull Request

### PR Checklist

- [ ] Code builds successfully (`go build ./...`)
- [ ] All tests pass (`go test ./...`)
- [ ] New rules have corresponding test cases
- [ ] New language parser has been tested manually
- [ ] Documentation updated (README, docs/, etc.)
- [ ] No unintended debug or dead code left behind

---

## Architecture Decision Records

Major architectural decisions are documented in `docs/memory/decisions.md`. Key decisions include:

| ADR | Decision | Rationale |
|-----|----------|-----------|
| D1  | BaseRegexParser via composition (embedding) | Go doesn't support multiple inheritance; composition allows AST + regex |
| D2  | GoParser keeps AST logic separate | AST checks are Go-specific, orthogonal to regex |
| D3  | Standard Go testing + go-cmp | Avoid heavyweight deps (testify) while solving struct comparison |
| D4  | Registry tests in cmd/goreview | internal/parser can't import langs (circular dependency) |
| D5  | JS-only handles .js/.jsx/.mjs/.cjs, TS handles .ts/.tsx | Fix extension detection bug where TS files were misidentified |
