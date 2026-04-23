# CodeSentry Architecture

> **Version:** 1.0.0
> **Last Updated:** 2026-04-23

## Overview

CodeSentry is a CLI static analysis tool that scans source code for security vulnerabilities and code quality issues. It uses a rule-driven detection engine with plugin-based language parsers, and optionally enriches findings with AI-powered code context for Go files via [cloudwego/abcoder](https://github.com/cloudwego/abcoder).

## High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                    CLI (Cobra)                       │
│  scan  languages  version                           │
└──────────┬──────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────┐
│               Scanning Engine (engine.go)            │
│  1. File discovery (filepath.Walk)                   │
│  2. Language detection (extension → parser)         │
│  3. Rule matching (per-file, per-rule, per-pattern) │
│  4. Deduplication (file:line:ruleID)               │
│  5. [AI context enrichment for Go] (abcoder)       │
└──────────┬─────────────────────────────┬────────────┘
           │                             │
           ▼                             ▼
┌─────────────────────┐       ┌──────────────────────┐
│   Parser Registry    │       │  abcoder Integration  │
│ (plugin pattern)     │       │  (Go only)            │
│  11 language parsers │       │  UniAST → Context    │
└─────────────────────┘       └──────────────────────┘
           │                             │
           ▼                             ▼
┌─────────────────────┐       ┌──────────────────────┐
│   Rule Loader       │       │  SkillAgent / Fallback│
│  YAML files → Rules │       │  Fix suggestions      │
└─────────────────────┘       └──────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────┐
│              Output Formatters (output.go)           │
│  Text  /  JSON  /  SARIF 2.1.0                      │
└─────────────────────────────────────────────────────┘
```

## Core Components

### 1. CLI Layer (`cmd/goreview/`)

| File         | Responsibility                                      |
|-------------|-----------------------------------------------------|
| `main.go`   | Entry point, calls `Execute()`                      |
| `root.go`    | Root Cobra command, version constant, global flags  |
| `scan.go`    | `scan` subcommand: loads rules, runs engine, outputs results |
| `languages.go` | `languages` subcommand: lists registered parsers   |
| `langs.go`   | Blank imports for all language parsers (triggers `init()` registration) |

**Build entry point is `cmd/goreview/`, not `cmd/codesentry/`.**

### 2. Engine (`internal/engine/`)

The engine orchestrates the scan pipeline.

#### `engine.go` — Core Scan Loop

```
Scan(paths, cfg) → *types.Result
```

Key steps:

1. **File Discovery**: `filepath.Walk` traverses paths, skips `node_modules`, `.git`, `vendor`, and files without recognized extensions.

2. **Language Detection**: `parser.DetectFromPath(filePath)` looks up the file extension in the parser registry.

3. **Rule Matching**: For each file:
   - Group applicable rules by language (`rulesByLang`)
   - Call `parser.Parse(filePath, content, langRules)` → `[]Finding`
   - Convert `Finding` → `Issue` using matched rule metadata

4. **Deduplication**: Uses `file:line:ruleID` key to avoid reporting the same issue twice.

5. **Category Filtering**: After matching, issues are filtered by `--security` / `--performance` flags.

#### `engine_abcoder.go` — AI Context Enrichment

```
ScanWithContext(paths, cfg) → *types.Result
```

Wraps the regular scan with abcoder-based context enrichment:

1. Parses entire repository via abcoder UniAST
2. For each Go finding, retrieves `CodeContext`: function name, variables, call chain
3. Builds contextual suggestion by appending context to the rule's `suggestion` field
4. Non-Go files skip enrichment (fallback to static suggestions)

### 3. Parser Registry (`internal/parser/`)

#### `registry.go` — Plugin Pattern

```go
type Parser interface {
    Language() string                    // e.g., "go", "python"
    Extensions() []string                // e.g., [".go"], [".py", ".pyw"]
    Parse(filePath string, content []byte, langRules []rules.Rule) ([]Finding, error)
}
```

The registry is a `map[string]Parser`. Language parsers register themselves via `init()`:

```go
// In langs/golang/parser.go
func init() {
    parserpkg.Register(&GoParser{})
}
```

This allows the engine to be decoupled from individual language parsers.

#### `base.go` — BaseRegexParser

Ten of eleven parsers embed `BaseRegexParser`, which provides the standard regex matching loop:

```go
type BaseRegexParser struct{}

func (p *BaseRegexParser) ParseRegex(content []byte, langRules []rules.Rule) []Finding {
    // For each rule, for each regex pattern:
    //   - Compile regex
    //   - For each line: if matches → append Finding
}
```

**Why composition (embedding) over inheritance**: Go doesn't support multiple inheritance, but struct embedding allows `GoParser` to embed `BaseRegexParser` while also implementing its own AST checks.

### 4. Language Parsers (`langs/`)

| Parser     | Detection Type          | Notes                                          |
|-----------|------------------------|------------------------------------------------|
| `golang`  | AST + Regex            | Uses `go/ast` for goroutine/context analysis   |
| `python`  | Regex only             |                                                |
| `typescript` | Regex only          | Handles `.ts`, `.tsx`, `.mts`, `.cts`          |
| `javascript` | Regex only          | Handles `.js`, `.jsx`, `.mjs`, `.cjs`          |
| `java`    | Regex only             |                                                |
| `ruby`    | Regex only             |                                                |
| `rust`    | Regex only             |                                                |
| `cpp`     | Regex only             | Handles `.cpp`, `.cc`, `.cxx`, `.c++`, `.h`, `.hpp` |
| `php`     | Regex only             |                                                |
| `swift`   | Regex only             |                                                |
| `kotlin`  | Regex only             |                                                |

### 5. Rules (`internal/rules/`)

#### `types.go`

```go
type Rule struct {
    ID          string
    Name        string
    Description string
    Severity    string   // "SEVERE" | "WARNING" | "INFO"
    Category    string   // "security" | "performance"
    Languages   []string
    Suggestion  string
    Patterns    []Pattern
}

type Pattern struct {
    Type    string  // "regex" (only fully implemented type)
    Pattern string  // The regex string
    Comment string  // Human-readable description
}
```

#### `loader.go`

`LoadRules(dir string) []Rule` walks the `rules/` directory and loads all `.yaml` / `.yml` files. Each file maps to exactly one `Rule` struct.

### 6. abcoder Integration (`internal/abcoder/`)

#### `bridge.go` — UniAST Repository Wrapper

```go
type Bridge struct {
    repo     *uniast.Repository  // cloudwego/abcoder's UniAST representation
    repoPath string
    mu       sync.RWMutex
}
```

Key methods:
- `Parse(ctx)`: Parses the repository using abcoder, populates `repo`
- `GetContext(file, line)`: Retrieves `CodeContext` (function, variables, calls) at a location
- `IsAvailable(file)`: Returns `true` only for `.go` files

`CodeContext` structure:
```go
type CodeContext struct {
    File, Line, Column int
    FunctionName, FunctionContent string
    FunctionCalls, MethodCalls   []CallInfo
    LocalVariables, GlobalVariables []VarInfo
    UsedTypes                   []TypeInfo
}
```

#### `skill.go` — Fix Suggestion Generation

`SkillAgent.GenerateFix()` takes a rule ID and code context, and returns a structured fix suggestion (`Before`, `After`, `Explanation`, `Confidence`).

> **Current status**: The actual LLM invocation via MCP is stubbed out. The implementation returns hardcoded suggestions based on rule ID matching. Full skill integration requires connecting to Claude Code's code-reviewer skill via MCP.

#### `fallback.go` — Non-Go Fallback

When abcoder is unavailable (non-Go files, parse errors), `FallbackHandler` provides static suggestions from a built-in mapping keyed by rule ID.

### 7. Output (`internal/output/`)

Three output formats are supported:

| Format   | File Extension | Use Case                                    |
|----------|---------------|---------------------------------------------|
| Text     | stdout, other | Human-readable CLI output                   |
| JSON     | `.json`       | Scripted processing, CI integration          |
| SARIF    | `.sarif`      | GitHub Code Scanning, GitLab SAST, tools    |

The `scan -o <path>` command infers the format from the output file extension. Unknown extensions fall back to text.

SARIF output follows the [OASIS SARIF 2.1.0](https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-schema-v2.1.0.html) specification.

### 8. Types (`internal/types/`)

Core data structures shared across packages:

```go
type Issue struct {
    ID, Severity, Title, Message string
    File string
    Line, Column, EndLine int
    RuleID, Category, Suggestion string
    Source string  // "static" or "ai"
}

type Result struct {
    TotalFiles, TotalIssues int
    Severe, Warning, Info int
    Duration time.Duration
    Timestamp time.Time
    Issues []Issue
    FilesScanned int
}
```

## Data Flow

```
User runs: ./goreview scan ./src --security
          │
          ▼
RootCmd.Execute() → scanCmd.Run()
          │
          ▼
rules.LoadRules("rules") → []rules.Rule
          │
          ▼
engine.New(rules) → *Engine
          │
          ▼
eng.Scan(paths, cfg) → *types.Result
          │
          ├─ filepath.Walk(paths) → files[]
          │
          ├─ For each file:
          │   ├─ parser.DetectFromPath(file) → Parser
          │   ├─ parser.Parse(file, content, rules) → []Finding
          │   ├─ Finding → Issue (lookup rule metadata)
          │   └─ Deduplicate (file:line:ruleID)
          │
          ├─ [If not --no-ai and file is .go:]
          │   └─ abcoder.GetContext(file, line) → enrich Issue
          │
          └─ Count by severity
          │
          ▼
outputResults(result) → Text/JSON/SARIF
```

## Key Design Decisions

### ADR-1: Parser Plugin Registration via `init()`

Each language parser registers itself at package init time via `parser.Register()`. The `cmd/goreview/langs.go` blank-imports all parser packages to trigger registration. This keeps the engine agnostic to available languages.

### ADR-2: BaseRegexParser via Struct Embedding

Ten parsers embed `BaseRegexParser` for the standard per-line regex matching loop. GoParser additionally embeds `BaseRegexParser` but overrides `Parse()` to add AST-based checks. Composition over inheritance allows this mixing without requiring multiple inheritance.

### ADR-3: abcoder for Go Only

`abcoder.IsAvailable()` returns `true` only for `.go` files, as abcoder's UniAST currently only supports Go AST parsing. Non-Go languages use the `FallbackHandler` for static suggestions.

### ADR-4: Rule IDs in SCREAMING_SNAKE_CASE

All rule IDs use `SCREAMING_SNAKE_CASE` for consistency. Mixed-case IDs (e.g., `TSPrototype_POLLUTION`) should be corrected.

### ADR-5: `patterns: []` in AST-only Rules

Rules like `CONTEXT_LEAK` that rely entirely on AST checks (no regex) declare `patterns: []` in YAML. The parser checks `hasContextLeakRule` flag and runs AST logic directly, bypassing regex.

## Performance Characteristics

| Operation              | Complexity           | Notes                                      |
|------------------------|----------------------|--------------------------------------------|
| File discovery         | O(n) files           | Single filepath.Walk pass                  |
| Regex matching         | O(n × r × p × l)    | Files × Rules × Patterns × Lines           |
| abcoder parse          | O(repo size)         | Runs once per scan (if not --no-ai)       |
| abcoder GetContext     | O(functions in pkg)  | Linear scan per finding                   |
| Deduplication          | O(issues)            | HashMap lookup per finding                 |

The `--no-ai` flag skips abcoder parsing entirely, significantly reducing latency for large repositories.

## Dependencies

| Dependency               | Version | Purpose                                      |
|-------------------------|---------|---------------------------------------------|
| spf13/cobra             | v1.8.1  | CLI framework                               |
| gopkg.in/yaml.v3        | v3.0.1  | YAML rule parsing                           |
| cloudwego/abcoder        | v0.3.1  | Go AST → UniAST + code context              |
| google/go-cmp            | v0.7.0  | Struct comparison in tests                  |

## Future Considerations

- **AST YAML pattern engine**: Implement `type: ast` in the YAML schema as an actual AST pattern matcher
- **MCP integration**: Connect `SkillAgent` to Claude Code's code-reviewer skill via MCP for actual LLM-powered fix generation
- **Multi-language abcoder**: As abcoder adds support for more languages, extend `IsAvailable()` accordingly
- **Rule conflict detection**: Detect overlapping patterns across rules that could cause duplicate findings
