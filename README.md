# CodeSentry

A fast, extensible static analysis and AI-powered code review tool with support for **11 programming languages**.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

## Features

- **Multi-language support**: Go, Python, TypeScript, JavaScript, Java, Ruby, Rust, C++, PHP, Swift, Kotlin
- **YAML-based rules**: Easy to add custom security and performance rules
- **AI-powered context**: Uses [cloudwego/abcoder](https://github.com/cloudwego/abcoder) UniAST for Go code context understanding
- **Multiple output formats**: Text, JSON, SARIF
- **Single binary**: No external dependencies required for scanning
- **Extensible architecture**: Add parsers for any language

## Supported Languages

| Language      | Extensions                              | Parser Type        |
|---------------|----------------------------------------|--------------------|
| Go            | `.go`                                  | AST + Regex        |
| Python        | `.py`, `.pyw`, `.pyi`                  | Regex              |
| TypeScript    | `.ts`, `.tsx`, `.mts`, `.cts`           | Regex              |
| JavaScript    | `.js`, `.jsx`, `.mjs`, `.cjs`           | Regex              |
| Java          | `.java`                                | Regex              |
| Ruby          | `.rb`                                  | Regex              |
| Rust          | `.rs`                                  | Regex              |
| C++           | `.cpp`, `.cc`, `.cxx`, `.c++`, `.h`, `.hpp` | Regex          |
| PHP           | `.php`                                 | Regex              |
| Swift         | `.swift`                               | Regex              |
| Kotlin        | `.kt`, `.kts`                          | Regex              |

> **Note**: AI-powered code context (via abcoder) is currently available for **Go only**. Other languages fall back to static regex-based suggestions.

## Installation

### Using go install (Recommended)

```bash
go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
```

This installs the `goreview` binary to your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your `$PATH`.

### From Source

```bash
git clone https://github.com/Colin4k1024/codesentry.git
cd codesentry
go build -o goreview ./cmd/goreview
```

### Pre-built Binary

Download from the [Releases](https://github.com/Colin4k1024/codesentry/releases) page.

## Quick Start

```bash
# Scan a directory for security issues
goreview scan ./src --security

# Scan with all rules (default)
goreview scan ./src

# Scan with performance rules only
goreview scan ./src --performance

# Disable AI-powered suggestions (faster, no context enrichment)
goreview scan ./src --security --no-ai

# Output to JSON file
goreview scan ./src --security -o results.json

# Output in SARIF format (CI/CD integration)
goreview scan ./src --security -o results.sarif

# List supported languages
goreview languages

# Show version
goreview version
```

## Security Rules

### Cross-language Rules

| Rule ID           | Name                   | Severity | Languages                                    |
|-------------------|------------------------|----------|----------------------------------------------|
| `HARDCODED_SECRET`| Hardcoded Secret       | SEVERE   | go, javascript, python, typescript, java, ruby, rust, php, kotlin |
| `SQL_INJECTION`   | SQL Injection          | SEVERE   | go, python, java, typescript, rust, php, cpp |
| `SENSITIVE_LOG`   | Sensitive Data Logging | WARNING  | go, python, java, typescript, rust, php, cpp |

### Language-Specific Rules

| Language    | Rules                                                    |
|-------------|----------------------------------------------------------|
| Go          | Goroutine Leak (`GOROUTINE_LEAK`), Context Leak (`CONTEXT_LEAK`), Unsafe Deserialization (`UNSAFE_DESERIALIZATION`), Path Traversal (`PATH_TRAVERSAL`), Exec with Input (`EXECUTION`) |
| Python      | Pickle Deserialization (`PYTHON_PICKLE`), YAML Load (`YAML_LOAD`), Subprocess Shell (`PYTHON_SUBPROCESS`) |
| TypeScript  | Dangerous eval() (`TS_EVAL`), Prototype Pollution (`TSPrototype_POLLUTION`), innerHTML XSS (`TS_XSS`) |
| Java        | Unsafe Deserialization (`JAVA_DESERIALIZATION`)          |
| Rust        | Unsafe Block (`RUST_UNSAFE_BLOCK`)                      |
| PHP         | Dangerous Functions (`PHP_UNSAFE_FUNC`)                  |

## Configuration

### Rule Categories

Rules are organized by category:
- `security` — Security vulnerabilities (enabled with `--security`)
- `performance` — Performance issues (enabled with `--performance`)

### Custom Rules

Rules are defined in YAML files under `rules/`:

```
rules/
├── security/           # Cross-language security rules
├── go/                 # Go-specific rules
├── python/             # Python-specific rules
├── typescript/         # TypeScript/JavaScript-specific rules
└── ...
```

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
    pattern: '(?i)(execute|query)\s*\([^)]*\+[^)]*\)'
    comment: String concatenation in SQL query
```

### Rule YAML Schema

```yaml
id: UNIQUE_ID              # Required: Unique identifier (SCREAMING_SNAKE_CASE)
name: Human Name           # Required: Display name
description: What it does  # Required: Detailed description
severity: SEVERE|WARNING|INFO
category: security|performance
languages:                 # List of supported languages
  - go
  - python
suggestion: How to fix it # Recommended: Fix suggestion
patterns:                  # Detection patterns
  - type: regex           # Currently only "regex" is fully implemented
    pattern: 'regex'      # Regular expression pattern
    comment: Description  # Human-readable description of the match
```

## Architecture

```
codesentry/
├── cmd/goreview/           # CLI entry point (Cobra commands)
│   ├── main.go              # Application entry point
│   ├── root.go              # Root command and version
│   ├── scan.go              # scan command implementation
│   ├── languages.go         # languages command implementation
│   └── langs.go             # Language parser registration (blank imports)
│
├── internal/
│   ├── engine/              # Core scanning engine
│   │   ├── engine.go        # File traversal, rule matching, deduplication
│   │   └── engine_abcoder.go # AI context enrichment (abcoder integration)
│   ├── parser/              # Parser registry and base types
│   │   ├── registry.go      # Language parser registration (Plugin pattern)
│   │   └── base.go          # BaseRegexParser - shared regex logic
│   ├── rules/               # Rule loading and types
│   │   ├── loader.go        # YAML rule file loader
│   │   └── types.go         # Rule and Pattern struct definitions
│   ├── abcoder/             # abcoder integration (cloudwego/abcoder)
│   │   ├── bridge.go        # UniAST Repository wrapper, context retrieval
│   │   ├── skill.go         # Skill agent for fix suggestion generation
│   │   ├── fallback.go       # Fallback suggestions for non-Go languages
│   │   └── context.go       # CodeContext data structures
│   ├── output/              # Output formatting
│   │   └── output.go        # Text, JSON, SARIF formatters
│   └── types/
│       └── types.go         # Issue, Result data structures
│
├── langs/                   # Language parsers (plugin pattern)
│   ├── golang/parser.go     # Go parser (AST + Regex)
│   ├── python/parser.go     # Python parser (Regex only)
│   ├── typescript/parser.go # TypeScript parser (Regex only)
│   └── ...                  # Other language parsers
│
└── rules/                   # YAML rule definitions
    ├── security/            # Cross-language security rules
    └── <lang>/              # Per-language rules
```

### Key Design Decisions

1. **Parser Plugin Pattern**: Each language parser registers itself via `init()` and `parser.Register()`. The engine does not import parsers directly; instead `cmd/goreview/langs.go` contains blank imports that trigger registration.

2. **BaseRegexParser**: Ten of eleven parsers embed `BaseRegexParser` which provides standard regex matching. Only the Go parser adds AST-based analysis.

3. **abcoder Integration**: The `abcoder` package wraps cloudwego/abcoder for Go code context understanding. It provides context enrichment (function name, variables, call chain) for Go files, with fallback suggestions for other languages.

4. **Rule-driven Detection**: Detection logic is defined in YAML rules, not hardcoded. The `type: regex` pattern is fully functional; `type: ast` is declared in the schema but not yet implemented by parsers (only Go's hardcoded AST checks work).

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run CodeSentry
  run: |
    curl -sL https://github.com/Colin4k1024/codesentry/releases/latest/download/goreview_darwin_amd64.tar.gz -o goreview.tar.gz
    tar -xzf goreview.tar.gz
    ./goreview scan ./src --security -o codesentry-results.json
```

### GitLab CI

```yaml
goreview:
  script:
    - ./goreview scan ./src --security -o codesentry-results.sarif
  artifacts:
    reports:
      sast: codesentry-results.sarif
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
  "severe": 3,
  "warning": 2,
  "info": 0,
  "issues": [...]
}
```

### SARIF (`-o results.sarif`)

Full [SARIF 2.1.0](https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-schema-v2.1.0.html) compliant output for integration with GitHub Code Scanning, GitLab SAST, and other security tools.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development guide including:
- Project setup
- Adding new language parsers
- Writing rules
- Testing
- abcoder integration details

## License

MIT License
