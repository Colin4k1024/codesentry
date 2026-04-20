# CodeSentry

A fast, extensible static analysis and code review tool with support for **11 programming languages**.

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8)
![Languages](https://img.shields.io/badge/languages-11-blue)

## Features

- **Multi-language support**: Go, Python, TypeScript, JavaScript, Java, Ruby, Rust, C++, PHP, Swift, Kotlin
- **YAML-based rules**: Easy to add custom security and performance rules
- **Multiple output formats**: Text, JSON, SARIF
- **No dependencies**: Single binary, works out of the box
- **Extensible architecture**: Add parsers for any language

## Installation

### From Source

```bash
git clone https://github.com/Colin4k1024/codesentry_refactor.git
cd codesentry_refactor
go build -o codesentry ./cmd/codesentry
```

### Pre-built Binary

Download from the [Releases](https://github.com/Colin4k1024/codesentry_refactor/releases) page.

## Quick Start

```bash
# Scan a directory for security issues
./codesentry scan ./src --security

# Scan with all rules
./codesentry scan ./src

# Scan with performance rules
./codesentry scan ./src --performance

# Output to JSON
./codesentry scan ./src --security -o results.json

# List supported languages
./codesentry languages
```

## Supported Languages

| Language | Extensions | Status |
|----------|------------|--------|
| Go | `.go` | ✅ Full AST + Regex |
| Python | `.py`, `.pyw`, `.pyi` | ✅ Regex |
| TypeScript | `.ts`, `.tsx`, `.mts`, `.cts` | ✅ Regex |
| JavaScript | `.js`, `.jsx`, `.mjs`, `.cjs` | ✅ Regex |
| Java | `.java` | ✅ Regex |
| Ruby | `.rb` | ✅ Regex |
| Rust | `.rs` | ✅ Regex |
| C++ | `.cpp`, `.cc`, `.cxx`, `.c++`, `.h`, `.hpp` | ✅ Regex |
| PHP | `.php` | ✅ Regex |
| Swift | `.swift` | ✅ Regex |
| Kotlin | `.kt`, `.kts` | ✅ Regex |

## Security Rules

### Cross-language Rules
- **Hardcoded Secret**: Detects API keys, passwords, tokens hardcoded in source
- **SQL Injection**: Detects string concatenation in SQL queries
- **Sensitive Data Logging**: Detects passwords/tokens logged to output

### Language-Specific Rules
| Language | Rules |
|----------|-------|
| Go | Goroutine Leak, Context Leak |
| Python | Pickle Deserialization |
| TypeScript/JavaScript | Dangerous eval(), XSS (innerHTML) |
| Java | Unsafe Deserialization |
| Ruby | Dangerous YAML.load |
| Rust | Unsafe Code Blocks |
| C++ | Buffer Overflow (strcpy, sprintf, gets) |
| PHP | unserialize(), eval(), assert() |
| Swift | Deprecated UIWebView |

## Configuration

### Rule Categories

Rules are organized by category:
- `security` - Security vulnerabilities
- `performance` - Performance issues

### Custom Rules

Rules are defined in YAML files under `rules/`:
- `rules/security/` - Security rules
- `rules/<lang>/` - Language-specific rules

Example rule (`rules/security/sql_injection.yaml`):
```yaml
id: SQL_INJECTION
name: SQL Injection
description: Detects string concatenation in SQL queries
severity: SEVERE
category: security
languages:
  - go
  - python
  - java
suggestion: Use parameterized queries
patterns:
  - type: regex
    pattern: '(?i)execute.*\+'
    comment: String concatenation in SQL query
```

## Architecture

```
codesentry/
├── cmd/codesentry/     # CLI entry point
├── internal/
│   ├── engine/       # Scanning engine
│   ├── parser/       # Parser registry
│   ├── rules/       # Rule loading
│   └── types/        # Shared types
├── langs/            # Language parsers
│   ├── golang/
│   ├── python/
│   ├── typescript/
│   └── ...
└── rules/            # YAML rule definitions
```

### Adding a New Language

1. Create a new parser in `langs/<lang>/parser.go`:
```go
package <lang>

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

func (p *Parser) Language() string { return "<lang>" }
func (p *Parser) Extensions() []string { return []string{".ext"} }

func (p *Parser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
    // Implement parsing logic
}
```

2. Register in `cmd/codesentry/langs.go`:
```go
import _ "github.com/Colin4k1024/codesentry/langs/<lang>"
```

3. Add rules in `rules/<lang>/`

4. Build and test:
```bash
go build -o codesentry ./cmd/codesentry
./codesentry languages  # Verify language appears
./codesentry scan /tmp/test.<ext> --security  # Test
```

## Output Formats

### Text (Default)
```
=== CodeSentry Scan Results ===
Files scanned: 3
Total issues: 5
  SEVERE:   3
  WARNING:  2
Duration: 1.5ms

=== Issues Found ===
[SEVERE] Hardcoded Secret
  File: src/auth.py:10:1
  Possible hardcoded secret
  Suggestion: Use environment variables or a secrets manager
```

### JSON (`-o results.json`)
```json
{
  "timestamp": "2026-04-20T10:00:00Z",
  "files_scanned": 3,
  "total_issues": 5,
  "issues": [...]
}
```

### SARIF (for CI integration)
```json
{
  "$schema": "https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-schema.json",
  "version": "2.1.0",
  "runs": [{
    "results": [...]
  }]
}
```

## CI/CD Integration

### GitHub Actions
```yaml
- name: Run CodeSentry
  run: |
    curl -sL https://github.com/Colin4k1024/codesentry_refactor/releases/latest/download/codesentry_linux_amd64 -o codesentry
    chmod +x codesentry
    ./codesentry scan ./src --security -o codesentry-results.json
```

### GitLab CI
```yaml
codesentry:
  script:
    - ./codesentry scan ./src --security -o codesentry-results.sarif
  artifacts:
    reports:
      sast: codesentry-results.sarif
```

## License

MIT License
