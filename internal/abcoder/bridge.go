// Package abcoder provides integration with cloudwego/abcoder for code context understanding
package abcoder

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudwego/abcoder/lang"
	"github.com/cloudwego/abcoder/lang/collect"
	"github.com/cloudwego/abcoder/lang/uniast"
)

// Bridge wraps abcoder UniAST parser for CodeSentry
type Bridge struct {
	repo     *uniast.Repository
	repoPath string
	mu       sync.RWMutex
}

// CodeContext holds the context around a code location
type CodeContext struct {
	File            string
	Line            int
	Column          int
	FunctionName    string
	FunctionContent string
	FunctionCalls   []CallInfo
	MethodCalls     []CallInfo
	LocalVariables  []VarInfo
	GlobalVariables []VarInfo
	UsedTypes       []TypeInfo
}

// CallInfo describes a function/method call
type CallInfo struct {
	Name    string
	PkgPath string
	File    string
	Line    int
}

// VarInfo describes a variable
type VarInfo struct {
	Name    string
	Type    string
	Content string
	Line    int
}

// TypeInfo describes a type usage
type TypeInfo struct {
	Name    string
	PkgPath string
	Kind    string
}

// NewBridge creates a new abcoder Bridge
func NewBridge(repoPath string) (*Bridge, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	return &Bridge{
		repoPath: absPath,
	}, nil
}

// Parse parses the repository using abcoder
func (b *Bridge) Parse(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	opts := lang.ParseOptions{
		CollectOption: collect.CollectOption{
			Language:           uniast.Golang,
			LoadExternalSymbol: false,
			NotNeedTest:        false,
			NoNeedComment:      true,
		},
	}

	result, err := lang.Parse(ctx, b.repoPath, opts)
	if err != nil {
		return fmt.Errorf("abcoder parse failed: %w", err)
	}

	var repo uniast.Repository
	if err := json.Unmarshal(result, &repo); err != nil {
		return fmt.Errorf("failed to unmarshal UniAST: %w", err)
	}

	b.repo = &repo
	return nil
}

// GetContext returns the code context at the given file and line
func (b *Bridge) GetContext(file string, line int) (*CodeContext, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.repo == nil {
		return nil, fmt.Errorf("repository not parsed, call Parse first")
	}

	// Find the file in repository
	f, mod := b.repo.GetFile(file)
	if f == nil {
		return nil, fmt.Errorf("file not found in repository: %s", file)
	}

	ctx := &CodeContext{
		File:            file,
		Line:            line,
		FunctionCalls:   []CallInfo{},
		MethodCalls:     []CallInfo{},
		LocalVariables:  []VarInfo{},
		GlobalVariables: []VarInfo{},
		UsedTypes:       []TypeInfo{},
	}

	if f.Package == "" {
		return ctx, nil
	}

	pkg := mod.Packages[f.Package]
	if pkg == nil {
		return ctx, nil
	}

	// Find function containing the line
	for _, fn := range pkg.Functions {
		if fn.File == file && fn.Line <= line {
			// Check if line is within function
			if fn.Line <= line && (fn.Line+strings.Count(fn.Content, "\n")) >= (line-fn.Line) {
				ctx.FunctionName = fn.Name
				ctx.FunctionContent = fn.Content

				// Collect function calls
				for _, call := range fn.FunctionCalls {
					ctx.FunctionCalls = append(ctx.FunctionCalls, CallInfo{
						Name:    call.Name,
						PkgPath: string(call.ModPath),
						File:    call.File,
						Line:    call.Line,
					})
				}

				// Collect method calls
				for _, call := range fn.MethodCalls {
					ctx.MethodCalls = append(ctx.MethodCalls, CallInfo{
						Name:    call.Name,
						PkgPath: string(call.ModPath),
						File:    call.File,
						Line:    call.Line,
					})
				}

				// Collect local variables (from Params and Results)
				for _, dep := range fn.Params {
					ctx.LocalVariables = append(ctx.LocalVariables, VarInfo{
						Name:    dep.Name,
						Type:    dep.PkgPath,
						Content: "",
						Line:    dep.Line,
					})
				}

				// Collect global variables used
				for _, dep := range fn.GlobalVars {
					if varNode := b.repo.GetVar(dep.Identity); varNode != nil {
						ctx.GlobalVariables = append(ctx.GlobalVariables, VarInfo{
							Name:    varNode.Name,
							Type:    b.getTypeName(varNode.Type),
							Content: varNode.Content,
							Line:    varNode.Line,
						})
					}
				}

				// Collect used types
				for _, dep := range fn.Types {
					if typeNode := b.repo.GetType(dep.Identity); typeNode != nil {
						ctx.UsedTypes = append(ctx.UsedTypes, TypeInfo{
							Name:    typeNode.Name,
							PkgPath: string(typeNode.ModPath),
							Kind:    string(typeNode.TypeKind),
						})
					}
				}

				break
			}
		}
	}

	return ctx, nil
}

// GetFunction returns the function with the given identity
func (b *Bridge) GetFunction(nodeID string) (*uniast.Function, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.repo == nil {
		return nil, fmt.Errorf("repository not parsed")
	}

	identity := parseIdentity(nodeID)
	fn := b.repo.GetFunction(identity)
	if fn == nil {
		return nil, fmt.Errorf("function not found: %s", nodeID)
	}
	return fn, nil
}

// GetVariable returns the variable with the given identity
func (b *Bridge) GetVariable(nodeID string) (*uniast.Var, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.repo == nil {
		return nil, fmt.Errorf("repository not parsed")
	}

	identity := parseIdentity(nodeID)
	v := b.repo.GetVar(identity)
	if v == nil {
		return nil, fmt.Errorf("variable not found: %s", nodeID)
	}
	return v, nil
}

// GetCallChain returns the call chain for a function
func (b *Bridge) GetCallChain(nodeID string) ([]string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.repo == nil {
		return nil, fmt.Errorf("repository not parsed")
	}

	identity := parseIdentity(nodeID)
	fn := b.repo.GetFunction(identity)
	if fn == nil {
		return nil, fmt.Errorf("function not found: %s", nodeID)
	}

	var chain []string
	for _, call := range fn.FunctionCalls {
		chain = append(chain, call.Identity.Full())
	}
	return chain, nil
}

// parseIdentity parses a node ID string into Identity
func parseIdentity(nodeID string) uniast.Identity {
	// Format: ModPath?PkgPath#Name
	var modPath, pkgPath, name string

	if idx := strings.Index(nodeID, "?"); idx != -1 {
		modPath = nodeID[:idx]
		nodeID = nodeID[idx+1:]
	}

	if idx := strings.Index(nodeID, "#"); idx != -1 {
		pkgPath = nodeID[:idx]
		name = nodeID[idx+1:]
	} else {
		name = nodeID
	}

	return uniast.Identity{
		ModPath: uniast.ModPath(modPath),
		PkgPath: uniast.PkgPath(pkgPath),
		Name:    name,
	}
}

func (b *Bridge) getTypeName(t *uniast.Identity) string {
	if t == nil {
		return ""
	}
	return t.Name
}

// IsAvailable checks if abcoder is available for the given file
func IsAvailable(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".go":
		return true
	default:
		return false
	}
}
