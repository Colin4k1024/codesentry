# CodeSentry Rules Reference

Detailed reference for all security and performance rules supported by CodeSentry.

## Security Rules

### Cross-Language Rules

These rules apply to multiple languages:

| Rule ID | Languages | Severity | Description |
|---------|-----------|----------|-------------|
| HARDCODED_SECRET | go, javascript, python, typescript, java, ruby, rust, php, kotlin | SEVERE | Detects hardcoded API keys, passwords, tokens, and secrets |
| SQL_INJECTION | go, python, java, typescript, rust, php, cpp | SEVERE | Detects string concatenation in SQL queries |
| SENSITIVE_LOG | go, python, java, typescript, rust, php, cpp | WARNING | Detects sensitive data in log statements |
| COMMAND_INJECTION | go, python, java, typescript | SEVERE | Detects unsanitized system command execution |
| PATH_TRAVERSAL | go, python, java, typescript | WARNING | Detects unsanitized file path operations |
| INSECURE_DESERIALIZATION | go, python, java, typescript | SEVERE | Detects unsafe deserialization of untrusted data |

### Language-Specific Rules

#### Go Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| GOROUTINE_LEAK | WARNING | Goroutines started without proper lifecycle management |
| CONTEXT_LEAK | WARNING | Context not properly cancelled or passed |
| JWT_ERROR | WARNING | JWT parsing error handling issues |
| RESOURCE_LEAK | WARNING | Unreleased resources (files, connections) |
| UNSAFE_DESERIALIZATION | SEVERE | Unsafe deserialization using encoding.BinaryUnmarshaler |
| PATH_TRAVERSAL | WARNING | Unsanitized file path operations |
| EXECUTION | SEVERE | Execution of external commands with user input |

#### Python Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| PYTHON_PICKLE | SEVERE | Pickle deserialization of untrusted data |
| YAML_LOAD | SEVERE | Unsafe YAML loading |
| PYTHON_SUBPROCESS | WARNING | Shell injection in subprocess calls |

#### TypeScript/JavaScript Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| TS_EVAL | SEVERE | Dangerous eval() usage |
| TS_PROTOTYPE_POLLUTION | SEVERE | Prototype pollution vulnerability |
| TS_XSS | SEVERE | innerHTML XSS vulnerability |

#### Java Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| JAVA_DESERIALIZATION | SEVERE | Unsafe Java deserialization |

#### Rust Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| RUST_UNSAFE_BLOCK | WARNING | Unchecked unsafe code blocks |

#### PHP Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| PHP_UNSAFE_FUNC | WARNING | Use of dangerous PHP functions |

## Performance Rules

### Cross-Language Performance Rules

| Rule ID | Languages | Severity | Description |
|---------|-----------|----------|-------------|
| GOROUTINE_LEAK | go | WARNING | Goroutines not properly terminated with errgroup or context |
| CONTEXT_LEAK | go | WARNING | Context not properly passed or cancelled |
| MEMORY_LEAK | go | WARNING | Unreleased memory allocations |
| INEFFICIENT_STRING | go | WARNING | Inefficient string concatenation in loops |

### Go Performance Rules

| Rule ID | Severity | Description |
|---------|----------|-------------|
| GOROUTINE_LEAK | WARNING | Goroutines without lifecycle management (no errgroup, context, or sync.WaitGroup) |
| CONTEXT_LEAK | WARNING | Context values not properly cancelled or leaked |
| RESOURCE_LEAK | WARNING | Files, connections, or other resources not properly closed |

## Rule Pattern Examples

### HARDCODED_SECRET Pattern

```yaml
id: HARDCODED_SECRET
name: Hardcoded Secret
description: Detects hardcoded API keys, passwords, tokens
severity: SEVERE
category: security
languages: [go, python, javascript, typescript, java]
suggestion: Use environment variables or a secrets manager
patterns:
  - type: regex
    pattern: '(?i)(api[_-]?key|secret|password|token|auth)\s*[:=]\s*["\'][^"\']{8,}["\']'
    comment: Hardcoded credential detected
```

### SQL_INJECTION Pattern

```yaml
id: SQL_INJECTION
name: SQL Injection
description: Detects string concatenation in SQL queries
severity: SEVERE
category: security
languages: [go, python, java]
suggestion: Use parameterized queries
patterns:
  - type: regex
    pattern: '(?i)(execute|query|exec)\s*\([^)]*\+[^)]*\)'
    comment: String concatenation in SQL query
```

### GOROUTINE_LEAK Pattern

```yaml
id: GOROUTINE_LEAK
name: Goroutine Leak
description: Goroutines without proper lifecycle management
severity: WARNING
category: performance
languages: [go]
suggestion: Use errgroup or context to manage goroutine lifecycle
patterns:
  - type: regex
    pattern: 'go\s+func\s*\('
    comment: Goroutine started without lifecycle management
```

## Severity Levels

| Level | Meaning |
|-------|---------|
| **SEVERE** | Critical security vulnerability - must fix immediately |
| **WARNING** | Performance issue or moderate security concern - should fix |
| **INFO** | Best practice suggestion - consider fixing |

## Rule File Format

Rules are defined in YAML files under the `rules/` directory:

```yaml
id: UNIQUE_RULE_ID
name: Human Readable Name
description: What this rule detects
severity: SEVERE|WARNING|INFO
category: security|performance
languages: [go, python, javascript]
suggestion: How to fix the issue
patterns:
  - type: regex
    pattern: 'regex pattern here'
    comment: Description of what this pattern matches
```

## Adding Custom Rules

1. Create a YAML file in the appropriate `rules/` subdirectory:
   - `rules/security/` for cross-language security rules
   - `rules/<language>/` for language-specific rules

2. Follow the rule schema above

3. Test the rule:
   ```bash
   goreview scan ./src --security
   ```

## Excluding Rules

Use `--exclude` flag to skip specific paths:

```bash
goreview scan ./src --exclude "vendor" --exclude "testdata" --exclude "node_modules"
```
