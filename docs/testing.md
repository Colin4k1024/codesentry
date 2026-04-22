# Testing Guide

> **Version:** 0.3.0
> **Last Updated:** 2026-04-22

## Overview

CodeSentry uses a multi-layered testing strategy:

| Test Type          | Package                     | Purpose                                       |
|-------------------|-----------------------------|-----------------------------------------------|
| Unit Tests        | `internal/rules/`, `internal/parser/` | Rule loading, regex parsing, registry         |
| Engine Tests      | `internal/engine/`           | Scan logic, filtering, deduplication, skips   |
| Integration Tests | `cmd/goreview/`              | Parser registration, golden file end-to-end  |
| E2E Tests         | `internal/abcoder/`          | abcoder integration, context retrieval         |

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run tests for a specific package
go test -v ./internal/engine/...
go test -v ./internal/parser/...
go test -v ./cmd/goreview/...

# Run tests matching a pattern
go test -v -run TestEngine ./...
go test -v -run TestParser ./...
```

## Test Packages

### `internal/rules/` — Rule Loading

| Test File           | What It Tests                                         |
|--------------------|------------------------------------------------------|
| `loader_test.go`   | `LoadRules()`, `FilterByLanguage()`, invalid YAML handling, rules without ID |

**Key tests:**
- `TestLoadRules`: Verifies all YAML files in `rules/` load successfully with required fields
- `TestLoadRules_FileByFile`: Verifies specific rule files (e.g., `HARDCODED_SECRET`, `SQL_INJECTION`) are found
- `TestFilterByLanguage`: Verifies rules are correctly filtered by language
- `TestLoadRules_InvalidYAML`: Verifies invalid YAML files are skipped gracefully
- `TestLoadRules_NoID`: Verifies rules without `id` field are skipped

### `internal/parser/` — Parser Infrastructure

| Test File           | What It Tests                                         |
|--------------------|------------------------------------------------------|
| `base_test.go`     | `BaseRegexParser.ParseRegex()` logic                 |
| `golden_test.go`   | Golden file loading and comparison utilities          |

**Key tests:**
- `TestBaseRegexParser_ParseRegex`: Tests regex matching across single/multiple matches, invalid regex handling, non-regex pattern types
- `TestBaseRegexParser_ParseRegex_Severity`: Verifies severity and rule ID propagate to findings

### `internal/engine/` — Scanning Engine

| Test File           | What It Tests                                         |
|--------------------|------------------------------------------------------|
| `engine_test.go`   | Full scan pipeline: file discovery, rule matching, filtering, deduplication |

**Key tests:**
- `TestEngine_Scan`: Basic scan of Go + Python files
- `TestEngine_Scan_SecurityFilter`: `--security` flag filters non-security rules
- `TestEngine_Scan_PerformanceFilter`: `--performance` flag filters non-performance rules
- `TestEngine_Scan_NoFilter`: No flag returns all rules
- `TestEngine_Scan_SkipDirectories`: `node_modules`, `.git`, `vendor` are skipped
- `TestEngine_Scan_ExcludePattern`: `--exclude` flag filters directories
- `TestEngine_Scan_Dedup`: Same line+rule deduped
- `TestEngine_Scan_UnknownExtension`: Unknown extensions return no findings
- `TestEngine_Scan_ReadError`: Unreadable files handled gracefully

### `cmd/goreview/` — CLI Integration

| Test File           | What It Tests                                         |
|--------------------|------------------------------------------------------|
| `registry_test.go` | All 11 language parsers register correctly            |
| `golden_test.go`   | End-to-end scan matching against golden files        |

**Key tests:**
- `TestParserRegistry`: Confirms all 11 languages are registered
- `TestDetectFromPath`: Extension → language detection mapping
- `TestDetectFromPathUnknown`: Unknown extensions return nil parser
- `TestGoldenFile_HARDCODED_SECRET`: Full scan pipeline with golden comparison
- `TestGoldenFile_SQL_INJECTION`: SQL injection detection
- `TestGoldenFile_GOROUTINE_LEAK`: Go AST-based goroutine leak detection

### `internal/abcoder/` — AI Context

| Test File           | What It Tests                                         |
|--------------------|------------------------------------------------------|
| `bridge_test.go`   | `Bridge` creation, `IsAvailable()`, identity parsing |
| `skill_test.go`    | `SkillAgent` fix generation, structured output        |
| `fallback_test.go` | Fallback suggestion lookup and formatting             |
| `e2e_test.go`      | Full abcoder pipeline (requires abcoder binary)       |

## Golden File Testing

Golden file tests provide a reproducible reference for rule behavior.

### Test Data Location

```
testdata/
└── rules/
    ├── HARDCODED_SECRET.input.go    # Source code with known issues
    └── HARDCODED_SECRET.golden.json  # Expected findings (JSON)
```

### Golden File Format

```json
{
  "rule_id": "HARDCODED_SECRET",
  "findings": [
    {
      "line": 1,
      "column": 1,
      "end_line": 1,
      "severity": "SEVERE",
      "message": "Possible hardcoded secret"
    }
  ]
}
```

### Creating a Golden Test

1. Create an input file with known issues in `testdata/rules/`
2. Run the scanner manually to get actual findings
3. Create the `.golden.json` file with expected findings
4. Write the test:

```go
func TestGoldenFile_MY_RULE(t *testing.T) {
    tmpDir := t.TempDir()

    inputFile := filepath.Join(tmpDir, "test.py")
    inputContent := `password = "hardcoded_secret"`
    if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
        t.Fatal(err)
    }

    // Load rules
    rulesDir := filepath.Join("..", "..", "rules")
    allRules := rules.LoadRules(rulesDir)
    var myRules []rules.Rule
    for _, r := range allRules {
        if r.ID == "MY_RULE" {
            myRules = append(myRules, r)
        }
    }

    if len(myRules) == 0 {
        t.Fatal("MY_RULE rule not found")
    }

    e := engine.New(myRules)
    cfg := &engine.Config{}

    result, err := e.Scan([]string{inputFile}, cfg)
    if err != nil {
        t.Fatalf("Scan() returned error: %v", err)
    }

    // Compare with golden
    golden := loadGoldenFile(t, "MY_RULE")
    CompareFindings(t, convertToFindings(result.Issues), golden)
}
```

## Coverage

Run coverage analysis:

```bash
# Coverage for all packages
go test -cover ./...

# Per-package coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage Targets

| Package | Target |
|---------|--------|
| `internal/engine/` | ≥ 80% |
| `internal/parser/` | ≥ 80% |
| `internal/rules/` | ≥ 80% |
| `internal/abcoder/` | ≥ 70% |

## Mocking

The abcoder package uses real abcoder parsing in E2E tests (`e2e_test.go`). For unit tests that don't need actual abcoder behavior, mock the `Bridge` interface or use the `FallbackHandler`.

## Test Data

Test files are created in `t.TempDir()` for isolation. Do not use hardcoded paths like `/tmp/test.go`.

```go
// ✅ Correct: uses temp directory
tmpDir := t.TempDir()
testFile := filepath.Join(tmpDir, "test.py")

// ❌ Incorrect: uses fixed path
testFile := "/tmp/test.py"
```

## CI Integration

Add to GitHub Actions:

```yaml
- name: Run tests
  run: |
    go test -cover ./...

- name: Check formatting
  run: |
    go fmt ./...
    git diff --exit-code
```

## Debugging Failed Tests

```bash
# Run with full output
go test -v ./internal/engine/...

# Run specific test
go test -v -run TestEngine_Scan_SecurityFilter ./internal/engine/...

# Attach debugger
go test -debug ./internal/engine/...
```
