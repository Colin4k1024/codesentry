# AST Analysis Integration Proposal for CodeSentry

> **Date:** 2026-04-25  
> **Author:** CodeSentry Research  
> **Status:** Draft for Review

## Executive Summary

This document proposes integrating AST-level analysis into the CodeSentry rules engine to enable contextual security checks like distinguishing parameterized SQL queries from string-concatenated queries. The current architecture uses regex-based pattern matching which cannot understand code structure, leading to false positives and limited detection capabilities.

---

## 1. Current Architecture & Limitations

### 1.1 How Scanning Currently Works

```
engine.Scan(paths, cfg)
    ├── filepath.Walk(paths) → collect files
    ├── Group rules by language
    └── For each file:
        ├── parser.DetectFromPath(file) → Parser
        ├── parser.Parse(file, content, langRules) → []Finding
        └── Finding → Issue (lookup rule metadata)
```

### 1.2 Parser Architecture

| Parser | Detection Type | AST Capability |
|--------|---------------|-----------------|
| `golang` | AST + Regex | `go/ast` standard library, hardcoded checks for GOROUTINE_LEAK, RESOURCE_LEAK, JWT_ERROR |
| `python` | Regex only | None |
| `typescript` | Regex only | None |
| `javascript` | Regex only | None |
| `java` | Regex only | None |
| `ruby` | Regex only | None |
| `rust` | Regex only | None |
| `cpp` | Regex only | None |
| `php` | Regex only | None |
| `swift` | Regex only | None |
| `kotlin` | Regex only | None |

### 1.3 Current Limitations

1. **Regex-Based Detection**: SQL_INJECTION rule uses patterns like:
   ```yaml
   pattern: "(SELECT|INSERT|...).*(\\+|\\&|\\|f\"\\{)"
   ```
   This cannot distinguish:
   - `db.Query("SELECT * FROM users WHERE id = " + userID)` (vulnerable)
   - `db.Query("SELECT * FROM users WHERE id = $1", userID)` (safe, parameterized)

2. **No Context Awareness**: Regex sees a line as a string, not as code structure. A SQL query inside a parameterized API call looks identical to a concatenated one.

3. **Hardcoded AST Logic**: Go parser has AST checks written directly in `langs/golang/parser.go:checkAST()`. Adding AST checks for new rules requires modifying parser code.

4. **Pattern Type Schema Already Declared but Unused**: The `Pattern` struct in `internal/rules/types.go` has a `Type` field:
   ```go
   type Pattern struct {
       Type    string `yaml:"type"` // regex, ast
       Pattern string `yaml:"pattern"`
       Comment string `yaml:"comment"`
   }
   ```
   But `type: ast` is not implemented in `BaseRegexParser.ParseRegex()`.

5. **abcoder is Contextual, Not Structural**: The `cloudwego/abcoder` integration provides function/variable context around a finding but doesn't enable structural AST pattern matching.

---

## 2. Available AST Capabilities

### 2.1 Go (`go/ast` Standard Library)

The Go parser already uses `go/ast` and `go/parser`. Example from `langs/golang/parser.go`:

```go
fset := token.NewFileSet()
file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)

// AST traversal via ast.Inspect
ast.Inspect(file, func(n ast.Node) bool {
    switch node := n.(type) {
    case *ast.GoStmt:
        // Check goroutine leak
    case *ast.CallExpr:
        // Check function calls
    }
    return true
})
```

**What we can detect:**
- Call expressions (`*ast.CallExpr`) - function calls, method calls
- Call arguments - distinguish `Query(sql, args)` from `Query(sql + args)`
- Identifiers and selector expressions - `db.Query` vs `sql.Open`
- Import paths - detect `database/sql` usage

### 2.2 Other Languages

| Language | Parser | External AST Library Available |
|----------|--------|-------------------------------|
| Python | Regex only | `tree-sitter-python` via `github.com/smacker/go-tree-sitter` (already in go.mod) |
| JavaScript/TypeScript | Regex only | `tree-sitter-javascript` (already in go.mod via go-tree-sitter) |
| Java | Regex only | `tree-sitter-java` |
| Ruby | Regex only | `tree-sitter-ruby` |
| Rust | Regex only | `tree-sitter-rust` |
| C/C++ | Regex only | `tree-sitter-cpp` |

**Note**: `github.com/smacker/go-tree-sitter` is already in `go.mod` as an indirect dependency.

---

## 3. Proposed Integration Approaches

### Option A: Add `ast_callback` Field to Rule Schema

**Concept**: Extend the YAML rule schema to include an AST callback function reference.

```yaml
id: SQL_INJECTION
name: SQL Injection
severity: SEVERE
category: security
languages: [go]
suggestion: Use parameterized queries
patterns:
  - type: regex
    pattern: "SELECT|INSERT|UPDATE|DELETE"
    comment: SQL keyword found
ast_analyzer: SQLInjectionAnalyzer  # New field - references a registered analyzer
```

**Pros:**
- Rules stay declarative
- Separates AST logic from rule definitions

**Cons:**
- Requires maintaining a registry of AST analyzers
- YAML cannot express complex AST logic - would still need Go code

### Option B: Add `ast_query` Field for Pattern-Based AST Matching

**Concept**: Extend the `Pattern` struct with an AST query format.

```yaml
id: SQL_INJECTION
name: SQL Injection
patterns:
  - type: ast
    language: go
    query: |
      (call_expr
        function: (selector_expr
          object: (identifier "db")
          method: (identifier "Query"))
        arguments: (concatenated_string))
```

**Pros:**
- Declarative pattern-based approach
- Could use tree-sitter queries for non-Go languages

**Cons:**
- Complex to implement
- Limited expressiveness for complex contextual analysis

### Option C: Language-Specific AST Analyzers (Recommended)

**Concept**: Each language parser implements an `ASTAnalyzer` interface. The engine passes AST nodes to registered analyzers.

```go
// internal/parser/ast.go
type ASTCallback func(node interface{}, ctx *ASTContext) []Finding

type ASTContext struct {
    FilePath    string
    Content     []byte
    Fset        *token.FileSet  // for Go
    File        interface{}     // *ast.File for Go
    RuleID      string
}

type ASTAnalyzer interface {
    Analyze(ctx *ASTContext) []Finding
    RuleIDs() []string
}
```

**Pros:**
- Follows existing parser plugin pattern
- Each language controls its own AST logic
- Can use native AST libraries per language
- Extensible - new analyzers register themselves

**Cons:**
- Requires implementing Go interfaces for each language
- More code than declarative patterns

### Option D: Hybrid - Register AST Callbacks Per Rule ID

**Concept**: The Go parser has a map of `ruleID → AST analysis function`. Rules declare which AST checks they need.

```go
// In golang parser
var astCheckers = map[string]func(*GoParser, *ast.File, *token.FileSet) []Finding{
    "SQL_INJECTION": checkSQLInjection,
    "PATH_TRAVERSAL": checkPathTraversal,
}
```

```yaml
# In rules/go/sql_injection.yaml
id: SQL_INJECTION
ast_checks: [SQL_INJECTION]  # Request specific AST checks
```

**Pros:**
- Clean separation - rule declares intent, parser implements
- Follows existing `init()` registration pattern

**Cons:**
- Still requires per-language implementation

---

## 4. Recommended Approach: Option C with Go AST as Prototype

**Rationale:**
1. The Go parser already has AST infrastructure (`go/ast`, `token.FileSet`, `ast.Inspect`)
2. The parser plugin pattern (`Parser` interface + `init()` registration) scales to other languages
3. `tree-sitter` bindings are already in `go.mod` as indirect dependency
4. This approach doesn't break existing rules - regex and AST can coexist

### 4.1 Architecture

```
internal/parser/ast/
    ├── ast.go           # ASTCallback, ASTContext interfaces
    ├── registry.go      # Global AST analyzer registry
    └── findings.go      # Finding struct extension for AST

langs/golang/
    ├── parser.go        # Existing, add ASTAnalyzer registration
    └── analyzers/
        ├── sql_injection.go
        ├── path_traversal.go
        └── resource_leak.go    # Already exists as hardcoded, extract to here
```

### 4.2 Interface Design

```go
// internal/parser/ast/ast.go

package ast

// Context passed to AST analyzers
type Context struct {
    FilePath string
    Content  []byte
    Fset     *token.FileSet
    File     interface{}  // *ast.File for Go, *tree_sitter.Tree for others
}

// Finding from AST analysis with optional node reference
type Finding struct {
    RuleID   string
    Line     int
    Column   int
    EndLine  int
    Message  string
    Severity string
    // Optional: reference to AST node for deduplication
    NodeRef  interface{}
}

// Analyzer is implemented by language-specific AST analyzers
type Analyzer interface {
    // Analyze runs AST analysis and returns findings
    Analyze(ctx *Context, ruleIDs []string) []Finding
    
    // Language returns the supported language
    Language() string
}
```

### 4.3 Engine Integration

In `engine.Scan()`, after `parser.Parse()` returns findings, call AST analyzers:

```go
// In engine.go Scan()
for i, filePath := range filesToScan {
    // ... existing file reading and parser detection ...
    
    langRules := rulesByLang[p.Language()]
    findings, err := p.Parse(filePath, content, langRules)
    
    // NEW: Call AST analyzers
    if analyzer := parser.GetASTAnalyzer(p.Language()); analyzer != nil {
        astFindings := analyzer.Analyze(&parser.ASTContext{
            FilePath: filePath,
            Content:  content,
            Fset:     /* from parser */,
            File:     /* AST from parser */,
        }, extractRuleIDs(langRules))
        findings = append(findings, astFindings...)
    }
}
```

### 4.4 Parser Changes for Go

Modify `GoParser` to implement `ASTAnalyzer` and expose its AST:

```go
// langs/golang/parser.go

func (p *GoParser) Language() string { return "go" }

func (p *GoParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
    // Parse regex-based rules
    regexFindings := p.BaseRegexParser.ParseRegex(content, langRules)
    
    // Parse AST-based rules
    astFindings := p.checkAST(filePath, string(content), langRules)
    
    return append(regexFindings, astFindings...), nil
}

// New: Expose AST for engine-level analysis
func (p *GoParser) ParseWithAST(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, *ast.File, *token.FileSet, error) {
    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
    if err != nil {
        return nil, nil, nil, err
    }
    
    // ... analysis ...
    return findings, file, fset, nil
}
```

---

## 5. Example: AST-Based SQL_INJECTION Rule for Go

### 5.1 Current Regex-Based Rule (False Positives)

```yaml
# rules/security/sql_injection.yaml
id: SQL_INJECTION
name: SQL Injection
severity: SEVERE
category: security
languages: [go]
suggestion: Use parameterized queries instead of string formatting
patterns:
  - type: regex
    pattern: "(SELECT|INSERT|...).*(\\+|\\&|\\.format\\(|f\"\\{)"
    comment: SQL with dynamic string construction
```

**Problem:** Flags `db.Query("SELECT * FROM users", userID)` as vulnerable (false positive).

### 5.2 AST-Based Analysis

**Safe patterns (should NOT flag):**
```go
// Parameterized - arguments separate from SQL string
db.Query("SELECT * FROM users WHERE id = $1", userID)
db.Query("SELECT * FROM users WHERE id = ?", userID)
db.Query("SELECT * FROM users WHERE id = $1", id)

// Prepared statements
stmt, _ := db.Prepare("SELECT * FROM users WHERE id = ?")
stmt.QueryRow(id)
```

**Unsafe patterns (SHOULD flag):**
```go
// String concatenation - vulnerable
db.Query("SELECT * FROM users WHERE id = " + userID)

// Template with direct interpolation - vulnerable  
db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID))

// Template string - vulnerable
db.Queryf("SELECT * FROM users WHERE id = %s", userID)
```

### 5.3 AST Analysis Logic

```go
// langs/golang/analyzers/sql_injection.go

package golang

import (
    "go/ast"
    "go/parser"
    "go/token"
    "strings"
    
    "github.com/Colin4k1024/codesentry/internal/parser"
)

// sqlInjectionChecker implements AST-based SQL injection detection
func checkSQLInjection(file *ast.File, fset *token.FileSet, ruleIDs []string) []parser.Finding {
    var findings []parser.Finding
    
    hasRule := false
    for _, id := range ruleIDs {
        if id == "SQL_INJECTION" {
            hasRule = true
            break
        }
    }
    if !hasRule {
        return findings
    }
    
    ast.Inspect(file, func(n ast.Node) bool {
        call, ok := n.(*ast.CallExpr)
        if !ok {
            return true
        }
        
        // Check if it's a SQL method call
        if !isSQLCall(call) {
            return true
        }
        
        // Check argument structure
        if len(call.Args) < 1 {
            return true
        }
        
        firstArg := call.Args[0]
        
        // SAFE: First arg is a simple string literal (parameterized query)
        if isSimpleStringLit(firstArg) {
            return true  // No finding - parameterized query
        }
        
        // UNSAFE: First arg involves string concatenation or formatting
        if involvesStringConcat(firstArg) || isFormattedString(firstArg) {
            pos := fset.Position(call.Pos())
            findings = append(findings, parser.Finding{
                RuleID:   "SQL_INJECTION",
                Line:     pos.Line,
                Column:   pos.Column,
                EndLine:  pos.Line,
                Severity: "SEVERE",
                Message:  "SQL query constructed via string concatenation - use parameterized queries",
            })
        }
        
        return true
    })
    
    return findings
}

func isSQLCall(call *ast.CallExpr) bool {
    // Check for db.Query, db.QueryRow, db.Exec, etc.
    sel, ok := call.Fun.(*ast.SelectorExpr)
    if !ok {
        return false
    }
    
    // Check if object is a database variable
    if ident, ok := sel.X.(*ast.Ident); ok {
        // Common db variable names
        if strings.HasPrefix(ident.Name, "db") || 
           strings.HasPrefix(ident.Name, "sql") ||
           ident.Name == "stmt" {
            // Check method name
            switch sel.Sel.Name {
            case "Query", "QueryRow", "Exec", "QueryContext", "ExecContext":
                return true
            }
        }
    }
    
    return false
}

func isSimpleStringLit(expr ast.Expr) bool {
    _, ok := expr.(*ast.BasicLit)
    return ok && expr.(*ast.BasicLit).Kind == token.STRING
}

func involvesStringConcat(expr ast.Expr) bool {
    // Check for binary expressions with +
    if bin, ok := expr.(*ast.BinaryExpr); ok && bin.Op == token.ADD {
        return true
    }
    return false
}

func isFormattedString(expr ast.Expr) bool {
    // Check for fmt.Sprintf, fmt.Errorf, template.Execute, etc.
    if call, ok := expr.(*ast.CallExpr); ok {
        if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
            if ident, ok := sel.X.(*ast.Ident); ok {
                if (ident.Name == "fmt" && strings.HasPrefix(sel.Sel.Name, "Sprintf")) ||
                   (ident.Name == "fmt" && strings.HasPrefix(sel.Sel.Name, "Errorf")) {
                    return true
                }
            }
        }
    }
    return false
}
```

### 5.4 Updated Rule Definition

```yaml
# rules/go/sql_injection.yaml (new file, language-specific)
id: SQL_INJECTION
name: SQL Injection
description: Detects SQL queries constructed via string concatenation
severity: SEVERE
category: security
languages: [go]
suggestion: Use db.Query("SELECT ... WHERE id = $1", param) instead of string concatenation
patterns: []
ast_analysis: go_sql_injection  # References golang AST analyzer
```

---

## 6. Implementation Phases

### Phase 1: Extract Go AST Checks to Analyzer Package
- Move hardcoded AST checks from `langs/golang/parser.go:checkAST()` to `langs/golang/analyzers/`
- Create `Analyzer` interface in `internal/parser/ast/`
- Register analyzers via `init()` in each language package

### Phase 2: Add SQL Injection AST Analyzer
- Implement `checkSQLInjection()` in `langs/golang/analyzers/sql_injection.go`
- Update rule YAML to use `patterns: []` and `ast_analysis: sql_injection`
- Test with true positives (concatenated queries) and false positives (parameterized)

### Phase 3: Generalize Engine Integration
- Modify `engine.Scan()` to call AST analyzers after regex matching
- Add `parser.ASTContext` interface to expose AST from parsers
- Update `Parser` interface to optionally return AST context

### Phase 4: Tree-sitter for Non-Go Languages
- Implement `TSAnalyzer` using `tree-sitter` (already in go.mod)
- Add analyzers for JavaScript/TypeScript (SQL injection, XSS)
- Add analyzers for Python (SQL injection, command injection)

---

## 7. Summary

| Aspect | Current State | Proposed State |
|--------|--------------|----------------|
| SQL_INJECTION | Regex, high false positives | AST-aware, contextual |
| Go Parser | Hardcoded AST checks | Extensible analyzer registry |
| Other Languages | Regex only | Tree-sitter based AST |
| Engine | Regex findings only | AST findings after regex |
| Rule Schema | `type: ast` unused | Implemented via `ast_analysis` field |

The proposed architecture builds on existing patterns (parser registry, `init()` registration) and doesn't break existing rules. It adds AST capability incrementally: first extract Go AST checks, then add SQL injection, then extend to other languages.

---

## Appendix: Files to Modify

| File | Change |
|------|--------|
| `internal/rules/types.go` | Add `ASTAnalysis` field to `Rule` struct |
| `internal/parser/ast/ast.go` | New file: `ASTContext`, `Analyzer` interface |
| `internal/parser/ast/registry.go` | New file: analyzer registration |
| `internal/parser/registry.go` | Add `GetASTAnalyzer(lang)` function |
| `internal/engine/engine.go` | Call AST analyzers after regex matching |
| `langs/golang/parser.go` | Implement `Analyzer` interface, expose AST |
| `langs/golang/analyzers/sql_injection.go` | New file: AST-based SQL injection check |
| `rules/go/sql_injection.yaml` | New file: language-specific SQL injection rule |
| `rules/security/sql_injection.yaml` | Keep for non-Go languages |
