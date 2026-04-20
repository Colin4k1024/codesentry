# Rule Development Guide

This guide explains how to write effective detection rules for GoReview.

## Rule YAML Structure

Every rule is a YAML file with the following fields:

```yaml
id: UNIQUE_ID              # Required: Unique identifier
name: Human Name           # Required: Display name
description: What it does  # Required: Detailed description
severity: SEVERE|WARNING|INFO
category: security|performance
languages:
  - go
  - python
suggestion: How to fix it
patterns:
  - type: regex
    pattern: 'regex pattern'
    comment: What this pattern matches
```

## Severity Levels

| Level | Description |
|-------|-------------|
| SEVERE | Critical security vulnerability (SQL injection, RCE, etc.) |
| WARNING | Potential issue or code smell |
| INFO | Informational notice |

## Regex Best Practices

### 1. Use Case-Insensitive Flags Sparingly

```yaml
# Bad: Makes everything case-insensitive
pattern: '(?i)password.*=.*".*"'

# Good: Only match case-insensitive parts
pattern: '(?i)(password|token)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
```

### 2. Escape Special Characters

In YAML, these characters have special meaning: `:`, `#`, `|`, `>`

Always quote patterns containing these:
```yaml
pattern: '(?i)\.innerHTML\s*='
```

### 3. Anchor When Possible

```yaml
# Better: Match at word boundary
pattern: '\beval\s*\('

# Match only at start of assignment
pattern: '(?i)^[^/]*(password|token)\s*[=:]'
```

### 4. Common Patterns

#### Hardcoded Secrets
```yaml
pattern: '(?i)(password|token|api_key|secret)\s*[=:]\s*["\x27][^"\x27]{8,}["\x27]'
```

#### SQL Injection
```yaml
pattern: '(?i)(execute|query)\s*\([^)]*\+[^)]*\)'
```

#### Sensitive Logging
```yaml
pattern: '(?i)(log|print|echo)\([^)]*(password|token|secret)[^)]*\)'
```

## Testing Rules

### 1. Create a Test File

Create a file with known issues:
```python
# /tmp/test_rule.py
password = "hardcoded_secret_123"
```

### 2. Run the Scanner

```bash
./goreview scan /tmp/test_rule.py --security
```

### 3. Verify Detection

You should see your rule trigger on the test file.

## Language-Specific Considerations

### Python
- Indentation matters but regex doesn't see it
- `pickle.load()` and `yaml.load()` are common risks
- f-strings with variables in SQL queries

### JavaScript/TypeScript
- Template literals: `` `SELECT * FROM ${table}` ``
- React's `dangerouslySetInnerHTML`
- `eval()` and `new Function()`

### Go
- Goroutines without errgroup
- Context passed to goroutines
- JWT parsed with nil key function

### PHP
- `$_GET`, `$_POST`, `$_REQUEST` in SQL
- `unserialize()` on user input
- `eval()` on user input

### Ruby
- String interpolation in SQL: `"SELECT * FROM #{table}"`
- `YAML.load()` instead of `YAML.safe_load()`
- `send()` with user input
