package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	parserpkg "github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
)

func init() {
	parserpkg.Register(&GoParser{})
}

type GoParser struct {
	parserpkg.BaseRegexParser
}

func (p *GoParser) Language() string { return "go" }

func (p *GoParser) Extensions() []string { return []string{".go"} }

func (p *GoParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	var findings []parserpkg.Finding

	// AST-based findings first (higher priority - more accurate than regex)
	astFindings := p.checkAST(filePath, string(content), langRules)
	findings = append(findings, astFindings...)

	// Regex-based findings second (may be deduplicated if AST found same issue)
	regexFindings := p.BaseRegexParser.ParseRegex(content, langRules)
	findings = append(findings, regexFindings...)

	return findings, nil
}

func (p *GoParser) checkAST(filePath string, content string, langRules []rules.Rule) []parserpkg.Finding {
	var findings []parserpkg.Finding
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return findings
	}

	usesErrgroup := false
	usesJWT := false
	for _, imp := range file.Imports {
		if strings.Contains(imp.Path.Value, "golang.org/x/sync/errgroup") {
			usesErrgroup = true
		}
		if strings.Contains(imp.Path.Value, "jwt") {
			usesJWT = true
		}
	}

	hasGoroutineLeakRule := false
	hasResourceLeakRule := false
	hasJWTErrorRule := false
	hasSQLInjectionRule := false
	for _, r := range langRules {
		if r.ID == "GOROUTINE_LEAK" {
			hasGoroutineLeakRule = true
		}
		if r.ID == "RESOURCE_LEAK" {
			hasResourceLeakRule = true
		}
		if r.ID == "JWT_ERROR" {
			hasJWTErrorRule = true
		}
		if r.ID == "SQL_INJECTION" {
			hasSQLInjectionRule = true
		}
	}

	openedResources := make(map[string]token.Pos)
	closedResources := make(map[string]bool)

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GoStmt:
			if hasGoroutineLeakRule && !usesErrgroup {
				pos := fset.Position(node.Go)
				findings = append(findings, parserpkg.Finding{
					RuleID:   "GOROUTINE_LEAK",
					Line:     pos.Line,
					Column:   pos.Column,
					EndLine:  pos.Line,
					Severity: "WARNING",
					Message:  "Goroutine started without errgroup - no guarantee of graceful shutdown",
				})
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				if i >= len(node.Rhs) {
					continue
				}
				if ident, ok := lhs.(*ast.Ident); ok {
					if call, ok := node.Rhs[i].(*ast.CallExpr); ok {
						if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
							if base, ok := sel.X.(*ast.Ident); ok && base.Name == "sql" {
								if sel.Sel.Name == "Open" || sel.Sel.Name == "Query" || sel.Sel.Name == "QueryRow" {
									openedResources[ident.Name] = call.Pos()
								}
							}
						}
					}
				}
			}
		case *ast.ExprStmt:
			if call, ok := node.X.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "Close" {
						if ident, ok := sel.X.(*ast.Ident); ok {
							closedResources[ident.Name] = true
						}
					}
				}
			}
		case *ast.DeferStmt:
			if call, ok := node.Call.Fun.(*ast.SelectorExpr); ok {
				if call.Sel.Name == "Close" {
					if ident, ok := call.X.(*ast.Ident); ok {
						closedResources[ident.Name] = true
					}
				}
			}
		}
		return true
	})

	if hasResourceLeakRule {
		for name, pos := range openedResources {
			if !closedResources[name] {
				pos2 := fset.Position(pos)
				findings = append(findings, parserpkg.Finding{
					RuleID:   "RESOURCE_LEAK",
					Line:     pos2.Line,
					Column:   pos2.Column,
					EndLine:  pos2.Line,
					Severity: "WARNING",
					Message:  "Resource '" + name + "' opened but never closed",
				})
			}
		}
	}

	if hasJWTErrorRule && usesJWT {
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Parse" {
					if len(call.Args) >= 2 {
						if isNilExpr(call.Args[1]) {
							pos := fset.Position(call.Pos())
							findings = append(findings, parserpkg.Finding{
								RuleID:   "JWT_ERROR",
								Line:     pos.Line,
								Column:   pos.Column,
								EndLine:  pos.Line,
								Severity: "SEVERE",
								Message:  "JWT parsed with nil key function - signature not verified",
							})
						}
					}
				}
			}
			return true
		})
	}

	if hasSQLInjectionRule {
		sqlFindings := checkSQLInjection(file, fset)
		findings = append(findings, sqlFindings...)
	}

	return findings
}

// isNilExpr checks if an expression is a nil literal
func isNilExpr(n ast.Expr) bool {
	if ident, ok := n.(*ast.Ident); ok {
		return ident.Name == "nil"
	}
	return false
}

// checkSQLInjection performs AST-based SQL injection detection
func checkSQLInjection(file *ast.File, fset *token.FileSet) []parserpkg.Finding {
	var findings []parserpkg.Finding

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a SQL method call
		if !isSQLCall(call) {
			return true
		}

		// Need at least one argument (the SQL query string)
		if len(call.Args) < 1 {
			return true
		}

		firstArg := call.Args[0]

		// SAFE: First arg is a simple string literal (parameterized query)
		if isSimpleStringLit(firstArg) {
			return true // No finding - parameterized query
		}

		// UNSAFE: First arg involves string concatenation or formatting
		if involvesStringConcat(firstArg) || isFormattedString(firstArg) {
			pos := fset.Position(call.Pos())
			findings = append(findings, parserpkg.Finding{
				RuleID:   "SQL_INJECTION",
				Line:     pos.Line,
				Column:   pos.Column,
				EndLine:  pos.Line,
				Severity: "SEVERE",
				Message:  "SQL query constructed via string concatenation or formatting - use parameterized queries instead",
			})
		}

		return true
	})

	return findings
}

// isSQLCall checks if a call expression is a SQL database method call
func isSQLCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Collect all identifiers in the selector chain
	// e.g., for s.pool.QueryRow, we get ["s", "pool"]
	idents := collectSelectorIdents(sel.X)
	if len(idents) == 0 {
		return false
	}

	// Check if any identifier in the chain is a database variable
	for _, ident := range idents {
		if strings.HasPrefix(ident.Name, "db") ||
			strings.HasPrefix(ident.Name, "sql") ||
			ident.Name == "stmt" || ident.Name == "dba" || ident.Name == "database" ||
			ident.Name == "pool" {
			// Check method name for SQL operations
			switch sel.Sel.Name {
			case "Query", "QueryRow", "Exec", "QueryContext", "ExecContext", "QueryRowContext", "Prepare", "PrepareContext":
				return true
			}
		}
	}

	return false
}

// collectSelectorIdents recursively collects all identifiers from nested selector expressions
// e.g., for s.pool.QueryRow, sel.X is a SelectorExpr(s.pool), and we collect "s" and "pool"
func collectSelectorIdents(expr ast.Expr) []*ast.Ident {
	if ident, ok := expr.(*ast.Ident); ok {
		return []*ast.Ident{ident}
	}
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		return append(collectSelectorIdents(sel.X), collectSelectorIdents(sel.Sel)...)
	}
	return nil
}

// isSimpleStringLit checks if an expression is a simple string literal
func isSimpleStringLit(expr ast.Expr) bool {
	if basicLit, ok := expr.(*ast.BasicLit); ok {
		return basicLit.Kind == token.STRING
	}
	return false
}

// involvesStringConcat checks if an expression involves string concatenation
func involvesStringConcat(expr ast.Expr) bool {
	if bin, ok := expr.(*ast.BinaryExpr); ok && bin.Op == token.ADD {
		return true
	}
	return false
}

// isFormattedString checks if an expression involves string formatting (fmt.Sprintf, Queryf, etc.)
func isFormattedString(expr ast.Expr) bool {
	// Check for CallExpr (fmt.Sprintf, Queryf, etc.)
	if call, ok := expr.(*ast.CallExpr); ok {
		// Check for fmt.Sprintf
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				if ident.Name == "fmt" && (sel.Sel.Name == "Sprintf" || sel.Sel.Name == "Errorf") {
					return true
				}
			}
		}
		// Check for direct identifiers like Queryf, Execf
		if ident, ok := call.Fun.(*ast.Ident); ok {
			if strings.HasSuffix(ident.Name, "f") && strings.HasSuffix(ident.Name, "Query") == false {
				// Likely a fmt-style function
				return true
			}
		}
	}
	return false
}
