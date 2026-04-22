---
name: codesentry-review
description: >
  Use this skill when reviewing code for security vulnerabilities,
  performance issues, or quality problems. Activates for security scans,
  performance audits, code quality checks, and CI/CD integration tasks.
  Supports Go, Python, TypeScript, JavaScript, Java, Ruby, Rust, C++, PHP, Swift, and Kotlin.
origin: ECC
---

# CodeSentry Review Skill

Automated code review using [CodeSentry goreview](https://github.com/Colin4k1024/codesentry) CLI tool for security and performance analysis.

## When to Activate

- Scanning code for security vulnerabilities (secrets, injection, etc.)
- Auditing code for performance issues (goroutine leaks, context misuse)
- Pre-commit code quality checks
- CI/CD pipeline security scanning
- Code review requests involving supported languages
- Request mentions "goreview", "codesentry", "security scan", or "code review"

## Prerequisites

### Installation

```bash
go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
```

### Verification

```bash
# Check version
goreview version

# List supported languages
goreview languages
```

## Workflow

### 1. Determine Scan Type

Choose the appropriate scan based on review goals:

| Goal | Flag | Rules |
|------|------|-------|
| Full review | (default) | All security + performance rules |
| Security only | `--security` | HARDCODED_SECRET, SQL_INJECTION, SENSITIVE_LOG, etc. |
| Performance only | `--performance` | GOROUTINE_LEAK, CONTEXT_LEAK, etc. |

### 2. Execute Scan

```bash
# Basic syntax
goreview scan <path> [flags]

# Common flags:
#   --security          Scan security rules only
#   --performance       Scan performance rules only
#   --exclude <path>    Exclude specific paths (can repeat)
#   --no-ai            Disable AI enhancement
#   -o, --output <file> Output to file
```

### 3. Execute Commands

```bash
# Full scan (all rules)
goreview scan ./src

# Security scan
goreview scan ./src --security

# Performance scan excluding vendor and testdata
goreview scan ./src --performance --exclude "vendor" --exclude "testdata"

# Quick scan
goreview scan ./src --security --no-ai
```

### 4. Interpret Results

Output is text format with severity indicators:

```
=== CodeSentry Scan Results ===
Files scanned: 8
Total issues: 7
  SEVERE:   7
  WARNING:  0
  INFO:     0
Duration: 4.356125ms

=== Issues Found ===
[SEVERE] Hardcoded Secret
  File: path/to/file.go:31:1
  Possible hardcoded secret
  Suggestion: Use environment variables or a secrets manager
```

## Supported Languages

| Language      | Extensions                              | Security Rules | Performance Rules |
|---------------|----------------------------------------|---------------|------------------|
| Go            | `.go`                                  | Yes           | Yes              |
| Python        | `.py`, `.pyw`, `.pyi`                  | Yes           | Yes              |
| TypeScript    | `.ts`, `.tsx`, `.mts`, `.cts`           | Yes           | Yes              |
| JavaScript    | `.js`, `.jsx`, `.mjs`, `.cjs`           | Yes           | Yes              |
| Java          | `.java`                                | Yes           | Yes              |
| Ruby          | `.rb`                                  | Yes           | Yes              |
| Rust          | `.rs`                                  | Yes           | Yes              |
| C++           | `.cpp`, `.cc`, `.cxx`, `.c++`, `.h`, `.hpp` | Yes       | Yes              |
| PHP           | `.php`                                 | Yes           | Yes              |
| Swift         | `.swift`                               | Yes           | Yes              |
| Kotlin        | `.kt`, `.kts`                          | Yes           | Yes              |

## Security Rules Reference

| Rule ID                  | Severity | Description |
|--------------------------|----------|-------------|
| HARDCODED_SECRET         | SEVERE   | API keys, passwords, tokens in source |
| SQL_INJECTION            | SEVERE   | Unsanitized SQL query construction |
| SENSITIVE_LOG            | WARNING  | Sensitive data in log statements |
| COMMAND_INJECTION        | SEVERE   | Unsanitized system commands |
| PATH_TRAVERSAL           | WARNING  | Unsanitized file path operations |
| INSECURE_DESERIALIZATION | SEVERE   | Unsafe deserialization |
| GOROUTINE_LEAK           | WARNING  | Goroutines not properly terminated |
| CONTEXT_LEAK            | WARNING  | Context not properly cancelled |
| JWT_ERROR                | WARNING  | JWT error handling issues |
| RESOURCE_LEAK            | WARNING  | Unreleased resources |

## Performance Rules Reference

| Rule ID              | Severity | Description |
|----------------------|----------|-------------|
| GOROUTINE_LEAK       | WARNING  | Goroutines not properly terminated |
| CONTEXT_LEAK         | WARNING  | Context not properly cancelled |
| MEMORY_LEAK          | WARNING  | Unreleased memory allocations |
| INEFFICIENT_STRING   | WARNING  | Inefficient string operations |

## Checklist

### Pre-Scan
- [ ] goreview is installed (`goreview version` works)
- [ ] Target directory contains supported language files
- [ ] Exclusions configured for vendor, testdata, node_modules
- [ ] Scan type selected (security / performance / all)

### Post-Scan
- [ ] Review all SEVERE severity findings first
- [ ] Address findings before proceeding with changes
- [ ] Consider adding exclusions for false positives
- [ ] Document findings for code review report

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run CodeSentry Security Scan
  run: |
    go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
    goreview scan ./src --security -o security-results.txt
  continue-on-error: true
```

### Local CI Script

```bash
#!/bin/bash
# Fail build on SEVERE findings
goreview scan ./src --security --no-ai > scan-results.txt
if grep -q "SEVERE:" scan-results.txt; then
  echo "Security issues found!"
  cat scan-results.txt
  exit 1
fi
echo "Security scan passed"
```

## Examples

### Security Scan
```
User: Scan the src directory for security issues
Agent: Executes `goreview scan ./src --security` and presents findings
```

### Performance Audit
```
User: Check for performance issues in the codebase
Agent: Executes `goreview scan ./src --performance` and summarizes findings
```

### Full Review
```
User: Run a full code review
Agent: Executes `goreview scan ./src` and provides comprehensive report
```

### Excluding Problematic Paths
```
User: Scan everything except vendor and test files
Agent: Executes `goreview scan ./src --exclude "vendor" --exclude "testdata" --exclude "*_test.go"`
```

## Output Severity Levels

- **SEVERE**: Critical security vulnerabilities requiring immediate attention
- **WARNING**: Performance issues or moderate security concerns
- **INFO**: Informational findings, best practices suggestions

## Tips

1. **Prioritize SEVERE findings** - These are critical security issues
2. **Use `--no-ai`** in CI to speed up scanning
3. **Configure exclusions** for generated code, vendor, and test files
4. **Review WARNINGs** for performance improvements
5. **AI enhancement** (abcoder) currently only works with Go files
