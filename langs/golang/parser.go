package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	parserpkg "github.com/goreview/goreview/internal/parser"
	"github.com/goreview/goreview/internal/rules"
)

func init() {
	parserpkg.Register(&GoParser{})
}

type GoParser struct{}

func (p *GoParser) Language() string { return "go" }

func (p *GoParser) Extensions() []string { return []string{".go"} }

func (p *GoParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
	var findings []parserpkg.Finding
	text := string(content)
	lines := strings.Split(text, "\n")

	for _, rule := range langRules {
		for _, pattern := range rule.Patterns {
			if pattern.Type != "regex" {
				continue
			}
			re, err := regexp.Compile(pattern.Pattern)
			if err != nil {
				continue
			}
			for lineNum, line := range lines {
				if re.MatchString(line) {
					findings = append(findings, parserpkg.Finding{
						RuleID:   rule.ID,
						Line:     lineNum + 1,
						Column:   1,
						EndLine:  lineNum + 1,
						Severity: rule.Severity,
						Message:  pattern.Comment,
					})
				}
			}
		}
	}

	// Also do AST-based checks for Go-specific patterns that regex can't catch
	astFindings := p.checkAST(filePath, text, langRules)
	findings = append(findings, astFindings...)

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
	hasContextLeakRule := false
	hasResourceLeakRule := false
	for _, r := range langRules {
		if r.ID == "GOROUTINE_LEAK" {
			hasGoroutineLeakRule = true
		}
		if r.ID == "CONTEXT_LEAK" {
			hasContextLeakRule = true
		}
		if r.ID == "RESOURCE_LEAK" {
			hasResourceLeakRule = true
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
				p := fset.Position(pos)
				findings = append(findings, parserpkg.Finding{
					RuleID:   "RESOURCE_LEAK",
					Line:     p.Line,
					Column:   p.Column,
					EndLine:  p.Line,
					Severity: "WARNING",
					Message:  "Resource '" + name + "' opened but never closed",
				})
			}
		}
	}

	if hasContextLeakRule && usesJWT {
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

	return findings
}

func isNilExpr(n ast.Expr) bool {
	if ident, ok := n.(*ast.Ident); ok {
		return ident.Name == "nil"
	}
	return false
}
