package engine

import (
	"context"
	"path/filepath"

	"github.com/Colin4k1024/codesentry/internal/abcoder"
	"github.com/Colin4k1024/codesentry/internal/types"
)

// ScanWithContext scans and enriches issues with abcoder code context
func (e *Engine) ScanWithContext(paths []string, cfg *Config) (*types.Result, error) {
	// First, do regular scan
	result, err := e.Scan(paths, cfg)
	if err != nil {
		return nil, err
	}

	// If NoAI is set, skip context enrichment
	if cfg.NoAI {
		return result, nil
	}

	// Find the repo root (common ancestor of all paths)
	repoRoot := findRepoRoot(paths)

	// Create abcoder bridge
	bridge, err := abcoder.NewBridge(repoRoot)
	if err != nil {
		// If bridge creation fails, return original result
		return result, nil
	}

	// Parse the repository
	ctx := context.Background()
	if err := bridge.Parse(ctx); err != nil {
		// If parse fails, return original result
		return result, nil
	}

	// Enrich each issue with context (only for Go files)
	for i := range result.Issues {
		issue := &result.Issues[i]
		if !abcoder.IsAvailable(issue.File) {
			continue
		}

		// Get relative path from repo root
		relPath, err := filepath.Rel(repoRoot, issue.File)
		if err != nil {
			continue
		}

		// Get code context
		codeCtx, err := bridge.GetContext(relPath, issue.Line)
		if err != nil {
			continue
		}

		// If we have context, enhance the suggestion
		if codeCtx.FunctionName != "" {
			issue.Suggestion = buildContextualSuggestion(issue, codeCtx)
		}
	}

	return result, nil
}

// findRepoRoot finds the common root directory for paths
func findRepoRoot(paths []string) string {
	if len(paths) == 0 {
		return "."
	}

	repoRoot := paths[0]
	for _, path := range paths[1:] {
		repoRoot = commonPrefix(repoRoot, path)
	}
	return repoRoot
}

// commonPrefix returns the common prefix of two paths
func commonPrefix(a, b string) string {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	i := 0
	for i < minLen && a[i] == b[i] {
		i++
	}

	return a[:i]
}

// buildContextualSuggestion builds an enhanced suggestion with code context
func buildContextualSuggestion(issue *types.Issue, ctx *abcoder.CodeContext) string {
	suggestion := issue.Suggestion

	if ctx.FunctionName != "" {
		suggestion = suggestion + "\n\n[Context]\n"
		suggestion = suggestion + "Function: " + ctx.FunctionName + "\n"

		if len(ctx.FunctionCalls) > 0 {
			suggestion = suggestion + "Calls: "
			for j, call := range ctx.FunctionCalls {
				if j > 0 {
					suggestion = suggestion + ", "
				}
				suggestion = suggestion + call.Name
			}
			suggestion = suggestion + "\n"
		}

		if len(ctx.LocalVariables) > 0 {
			suggestion = suggestion + "Variables: "
			for j, v := range ctx.LocalVariables {
				if j > 0 {
					suggestion = suggestion + ", "
				}
				suggestion = suggestion + v.Name + " (" + v.Type + ")"
			}
			suggestion = suggestion + "\n"
		}
	}

	return suggestion
}
