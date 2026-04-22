# Rule Development Guide

This guide covers everything needed to write effective detection rules for CodeSentry.

## Rule YAML Schema

Every rule is a YAML file. Fields marked **required** must be present.

```yaml
id: UNIQUE_ID              # Required: Unique identifier (SCREAMING_SNAKE_CASE)
name: Human Name           # Required: Display name shown in output
description: What it does  # Required: Detailed explanation of the vulnerability
severity: SEVERE|WARNING|INFO   # Required: Impact level
category: security|performance  # Required: Determines --security / --performance flag
languages:                 # Required: List of language identifiers
  - go
  - python
suggestion: How to fix it # Recommended: Fix guidance shown in output
patterns:                 # Required: At least one pattern
  - type: regex           # Required: Pattern type (regex only fully implemented)
    pattern: 'regex'      # Required: Regular expression
    comment: Description  # Required: Human-readable match description
```

## Severity Levels

| Level   | When to use                                                            |
|---------|------------------------------------------------------------------------|
| SEVERE  | Critical security issues: injection, RCE, hardcoded secrets, auth bypass |
| WARNING | Potential bugs, code smells, resource leaks, deprecated API usage        |
| INFO    | Informational: style suggestions, minor improvements                      |

## Pattern Types

### `regex` (Fully Implemented)

The primary detection mechanism. Regex patterns are matched **per line** against the source file.

**How it works**: The parser splits the file content into lines, then runs each line through every regex pattern. When a pattern matches, it produces a Finding at that line number.

**Important**: Patterns should be written to match a **single line**. Multi-line patterns are not supported.

### `ast` (Declared but Not Implemented)

The YAML schema supports `type: ast`, but no parser currently implements AST-based pattern matching except Go's hardcoded AST checks in `langs/golang/parser.go`. Do not rely on `type: ast` for cross-language rules.

## Rule ID Naming Convention

All rule IDs must follow **SCREAMING_SNAKE_CASE** consistently.

| Correct          | Incorrect            | Issue                          |
|-----------------|---------------------|--------------------------------|
| `HARDCODED_SECRET` | `TSPrototype_POLLUTION` | Mixed case (use all caps)      |
| `SQL_INJECTION` | `SQLInjection`       | Not screaming snake            |
| `PYTHON_PICKLE` | `Python_Pickle`      | Python is conventionally all caps in IDs |

## Language Identifiers

Use lowercase language names matching the parser registration:

| Language    | Identifier   | Parser language() return value |
|------------|--------------|-------------------------------|
| Go         | `go`         | `"go"`                        |
| Python     | `python`      | `"python"`                    |
| TypeScript | `typescript` | `"typescript"`               |
| JavaScript | `javascript` | `"javascript"`                |
| Java       | `java`       | `"java"`                      |
| Ruby       | `ruby`       | `"ruby"`                      |
| Rust       | `rust`       | `"rust"`                      |
| C++        | `cpp`        | `"cpp"`                       |
| PHP        | `php`        | `"php"`                       |
| Swift      | `swift`      | `"swift"`                     |
| Kotlin     | `kotlin`     | `"kotlin"`                    |

## Writing Effective Regex Patterns

### Basic Structure

```yaml
patterns:
  - type: regex
    pattern: '\beval\s*\('
    comment: Dangerous use of eval()
```

### Anchor to Reduce False Positives

**Word boundaries** prevent partial matches within identifiers:

```yaml
# Bad: matches "ieval(..." or "myeval()"
pattern: 'eval\s*\('

# Good: only matches standalone "eval("
pattern: '\beval\s*\('
```

**Start-of-line anchors** skip comments and strings (use carefully):

```yaml
# Only matches at beginning of actual code line
pattern: '(?i)^[^/]*(password|token)\s*[=:]'
```

> The `^[^/]*` skips lines containing `//` comments at the start, but it also skips lines with inline comments. Test thoroughly.

### Case-Insensitive Matching

Use `(?i)` prefix for case-insensitive matching, but apply it only to the necessary part:

```yaml
# Bad: makes entire pattern case-insensitive
pattern: '(?i)password\s*=\s*".*"'

# Good: only "password" and common variants are case-insensitive
pattern: '(?i)(password|token|api_key)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
```

### Escaping in YAML

Special YAML characters must be quoted: `:`, `#`, `|`, `>`, `[`, `]`, `{`, `}`, `!`, `*`

```yaml
# Escape dot and parentheses in regex
pattern: '(?i)\.innerHTML\s*='

# Escape special characters in character classes
pattern: '["\x27]alice["\x27]'
```

### Common Vulnerability Patterns

#### Hardcoded Secrets

Detect API keys, passwords, tokens hardcoded in source:

```yaml
patterns:
  # Standard key=value patterns
  - type: regex
    pattern: '(?i)(password|token|api_key|secret|apiKey|secretKey)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
    comment: Possible hardcoded secret

  # Long base64-like strings (possible keys/tokens)
  - type: regex
    pattern: '(?i)const\s+\w+\s*=\s*["\x27][a-zA-Z0-9_\-]{20,}["\x27]'
    comment: Possible hardcoded secret (long string constant)
```

#### SQL Injection

Detect string concatenation in database queries:

```yaml
patterns:
  # String concatenation in execute/query
  - type: regex
    pattern: '(?i)(execute|query|exec)\s*\([^)]*\+[^)]*\)'
    comment: String concatenation in SQL query — use parameterized queries

  # f-string or format string with + in SQL context
  - type: regex
    pattern: '(?i)f["\x27].*\$\{?[^}]+\}?.*".*\%s'
    comment: Potential SQL injection via formatted string
```

#### Command Injection

Detect `exec.Command` with shell string concatenation:

```yaml
patterns:
  - type: regex
    pattern: 'exec\.Command\s*\(\s*["\x27][^"\x27]+["\x27]\s*\+\s*\w+'
    comment: Shell command constructed from string concatenation — use separate arguments
```

#### Dangerous Deserialization

Detect unsafe deserialization patterns:

```yaml
# Python pickle
patterns:
  - type: regex
    pattern: 'pickle\.loads?\s*\('
    comment: Unsafe pickle deserialization — use json.loads() for untrusted data

# Python yaml.load without SafeLoader
patterns:
  - type: regex
    pattern: 'yaml\.load\s*\([^)]*(?<!SafeLoader|Loader|Safe\)'
    comment: Unsafe YAML load without SafeLoader — may execute arbitrary code

# Java deserialization
patterns:
  - type: regex
    pattern: 'ObjectInputStream|readObject\s*\('
    comment: Unsafe Java deserialization — validate and sanitize input
```

#### Path Traversal

Detect file operations with unsanitized user input:

```yaml
patterns:
  - type: regex
    pattern: 'os\.Open\s*\([^)]*\+[^)]*\)'
    comment: Path concatenation in file open — validate and sanitize path components

  - type: regex
    pattern: 'open\s*\([^)]*\+[^)]*\)'
    comment: Potential path traversal — use secure path joining
```

#### Sensitive Data Logging

Detect logging of sensitive values:

```yaml
patterns:
  - type: regex
    pattern: '(?i)(log|fmt\.Print|echo|console\.log)\s*\([^)]*(password|token|secret|key|credential)[^)]*\)'
    comment: Sensitive data being logged — remove or mask before logging
```

#### Dangerous JavaScript Patterns

```yaml
patterns:
  # eval()
  - type: regex
    pattern: '(?<![.\w])eval\s*\('
    comment: Dangerous eval() — may execute arbitrary code

  # innerHTML XSS
  - type: regex
    pattern: '(?i)\.innerHTML\s*='
    comment: innerHTML assignment — use textContent or sanitize input

  # new Function()
  - type: regex
    pattern: 'new\s+Function\s*\('
    comment: Dynamic function creation — equivalent to eval()

  # Prototype pollution
  - type: regex
    pattern: '(?i)(Object\.assign|merge|extend)\s*\([^,]*req'
    comment: Potential prototype pollution attack
```

## Language-Specific Considerations

### Go

Go rules run on both regex and AST analysis.

**AST-based checks** (hardcoded in `langs/golang/parser.go`):
- `GOROUTINE_LEAK`: Detects `go` statements without `errgroup`
- `CONTEXT_LEAK`: Detects context passed to goroutines (requires JWT import)
- `RESOURCE_LEAK`: Detects opened resources not closed

**Regex patterns** cover:
- String concatenation in SQL (`database/sql`)
- Hardcoded secrets
- Logging sensitive data

### Python

- Regex patterns don't see indentation
- `pickle.load()` and `yaml.load()` are high-risk
- f-strings with variables in SQL: `f"SELECT * FROM {table}"`
- Subprocess calls: `subprocess.call(cmd, shell=True)`

### JavaScript/TypeScript

- Template literals: backtick strings with `${...}`
- `dangerouslySetInnerHTML` in React
- `eval()` and `new Function()`
- Prototype pollution via `Object.assign` with user input

### PHP

- `$_GET`, `$_POST`, `$_REQUEST` in SQL queries
- `unserialize()` on user input
- `eval()` and `assert()` on user input
- `exec()` with string concatenation

### Ruby

- String interpolation in SQL: `"SELECT * FROM #{table}"`
- `YAML.load()` instead of `YAML.safe_load()`
- `send()` with user input

## Rule Quality Checklist

Before submitting a rule:

- [ ] Rule ID is SCREAMING_SNAKE_CASE
- [ ] At least one test case exists (manual or golden file)
- [ ] Patterns are anchored where possible
- [ ] Special YAML characters are quoted
- [ ] Case-insensitive matching is applied only where needed
- [ ] `languages` list matches actual parser registrations
- [ ] `category` is `security` or `performance` (no typos)
- [ ] `suggestion` provides actionable fix guidance
- [ ] No false positive test cases found

## Rule File Organization

```
rules/
├── security/               # Cross-language security rules
│   ├── hardcoded_secret.yaml
│   ├── sql_injection.yaml
│   └── sensitive_log.yaml
├── go/                    # Go-specific rules
│   ├── context_leak.yaml   # (AST-based, patterns: [])
│   ├── goroutine_leak.yaml
│   └── unsafe_deserialization.yaml
├── python/                 # Python-specific rules
│   ├── pickle.yaml
│   └── yaml_load.yaml
└── <lang>/
    └── <rule_id>.yaml
```

Rules in `security/` should be cross-language and have the broadest `languages` list possible. Language-specific rules go in their respective directories.
