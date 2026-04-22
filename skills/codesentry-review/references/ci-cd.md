# CI/CD Integration Guide

Guide for integrating CodeSentry goreview scans into CI/CD pipelines.

## GitHub Actions

### Basic Security Scan

```yaml
name: Security Scan

on: [push, pull_request]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install goreview
        run: go install github.com/Colin4k1024/codesentry/cmd/goreview@latest

      - name: Run security scan
        run: goreview scan ./src --security --no-ai

      - name: Upload results
        if: failure()
        run: |
          goreview scan ./src --security --no-ai > security-results.txt
          echo "Security issues found:"
          cat security-results.txt
```

### GitHub Actions with SARIF (GitHub Code Scanning)

Note: goreview currently outputs text format. For GitHub Code Scanning integration, parse the text output.

```yaml
name: Security Scan with SARIF

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install goreview
        run: go install github.com/Colin4k1024/codesentry/cmd/goreview@latest

      - name: Run security scan
        run: |
          goreview scan ./src --security --no-ai > scan-results.txt
          # Parse and format for SARIF if needed

      - name: Upload results to GitHub Security
        if: always()
        run: |
          # GitHub Code Scanning accepts SARIF format
          # For now, just archive the text results
          cat scan-results.txt || echo "No issues found"
```

## GitLab CI

### Basic Security Scan

```yaml
security_scan:
  image: golang:latest
  before_script:
    - go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
  script:
    - goreview scan ./src --security --no-ai
  artifacts:
    when: always
    paths:
      - security-results.txt
    reports:
      # Note: GitLab SAST expects SARIF format
      # goreview text output needs conversion for full integration
```

## Local CI Script

### Bash Script with Exit Code

```bash
#!/bin/bash
# save as: scripts/security-scan.sh

set -e

echo "Running CodeSentry security scan..."
goreview scan ./src --security --no-ai > scan-results.txt

# Check for SEVERE issues
if grep -q "SEVERE:" scan-results.txt; then
  echo "=========================================="
  echo "SECURITY ISSUES FOUND - BUILD FAILED"
  echo "=========================================="
  cat scan-results.txt
  exit 1
fi

# Check for WARNING issues
if grep -q "WARNING:" scan-results.txt; then
  echo "=========================================="
  echo "PERFORMANCE ISSUES FOUND"
  echo "=========================================="
  cat scan-results.txt
  # Uncomment to fail on warnings:
  # exit 1
fi

echo "Scan completed successfully - no SEVERE issues found"
exit 0
```

### Makefile Integration

```makefile
.PHONY: security-scan

security-scan:
	@go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
	@goreview scan ./src --security --no-ai || (echo "Security scan failed" && exit 1)

security-scan-all:
	@go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
	@goreview scan ./src || (echo "Code review failed" && exit 1)

performance-scan:
	@go install github.com/Colin4k1024/codesentry/cmd/goreview@latest
	@goreview scan ./src --performance || (echo "Performance scan failed" && exit 1)
```

## Pre-commit Hook

### Shell Pre-commit Hook

```bash
#!/bin/bash
# save as: .git/hooks/pre-commit

echo "Running pre-commit security scan..."

goreview scan ./src --security --no-ai > /tmp/pre-commit-scan.txt

if grep -q "SEVERE:" /tmp/pre-commit-scan.txt; then
  echo "=========================================="
  echo "SECURITY ISSUES FOUND - COMMIT REJECTED"
  echo "=========================================="
  cat /tmp/pre-commit-scan.txt
  rm /tmp/pre-commit-scan.txt
  exit 1
fi

echo "Pre-commit scan passed"
rm /tmp/pre-commit-scan.txt
exit 0
```

### Python Pre-commit Hook

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: codesentry-security
        name: CodeSentry Security Scan
        entry: bash -c 'goreview scan ./src --security --no-ai'
        language: system
        pass_filenames: false
```

## Docker Integration

### Build-time Security Scan

```dockerfile
# Build your application
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

# Scan the built application
FROM golang:1.23
WORKDIR /app
COPY --from=builder /app/myapp .
RUN go install github.com/Colin4k1024/codesentry/cmd/goreview@latest && \
    goreview scan . --security --no-ai || true
```

## Integration with Other Tools

### Slack Notification on Failure

```yaml
- name: Run security scan
  run: |
    goreview scan ./src --security --no-ai > scan-results.txt || {
      echo "Security issues found, notifying Slack..."
      curl -X POST $SLACK_WEBHOOK \
        -H 'Content-Type: application/json' \
        --data "{\"text\":\"CodeSentry found security issues: $(cat scan-results.txt)\"}"
      exit 1
    }
```

### JIRA Integration

```bash
#!/bin/bash
goreview scan ./src --security --no-ai > scan-results.txt
if grep -q "SEVERE:" scan-results.txt; then
  # Create JIRA ticket
  curl -X POST $JIRA_API/issues \
    -H "Authorization: Bearer $JIRA_TOKEN" \
    -d '{
      "fields": {
        "project": {"key": "SEC"},
        "summary": "CodeSentry: Security issues found",
        "description": "'"$(cat scan-results.txt)"'",
        "issuetype": {"name": "Bug"},
        "priority": {"name": "High"}
      }
    }'
fi
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Scan completed successfully (issues may still be found) |
| 1 | Error during scan (file not found, etc.) |

Note: goreview does not currently return non-zero exit codes for findings. Use grep to check output for SEVERE/WARNING indicators.

## Best Practices

1. **Run with `--no-ai`** in CI for faster scans
2. **Exclude vendor/test/generated directories** to reduce noise
3. **Archive scan results** for audit trail
4. **Fail on SEVERE** but review WARNINGs manually
5. **Run regularly** (daily or on every push to main)
6. **Track trends** - monitor if issues are increasing or decreasing
