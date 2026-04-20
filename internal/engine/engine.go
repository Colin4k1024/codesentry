package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Colin4k1024/codesentry/internal/parser"
	"github.com/Colin4k1024/codesentry/internal/rules"
	"github.com/Colin4k1024/codesentry/internal/types"
)

// Config holds configuration for scanning
type Config struct {
	Security    bool
	Performance bool
	NoAI        bool
	Output      string
	Exclude     []string
}

// Engine is the core scanning engine
type Engine struct {
	rules []rules.Rule
}

// New creates a new Engine with the given rules
func New(rules []rules.Rule) *Engine {
	return &Engine{rules: rules}
}

// Scan scans the given paths for issues
func (e *Engine) Scan(paths []string, cfg *Config) (*types.Result, error) {
	result := &types.Result{
		Timestamp: time.Now(),
		Issues:    []types.Issue{},
	}

	// Collect all files to scan
	var filesToScan []string
	for _, path := range paths {
		if err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == "node_modules" || info.Name() == ".git" || info.Name() == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}
			for _, excl := range cfg.Exclude {
				if strings.Contains(filePath, excl) {
					return nil
				}
			}
			ext := strings.ToLower(filepath.Ext(filePath))
			if !isCodeFile(ext) {
				return nil
			}
			filesToScan = append(filesToScan, filePath)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	result.FilesScanned = len(filesToScan)

	// Group rules by language
	rulesByLang := make(map[string][]rules.Rule)
	for _, rule := range e.rules {
		for _, lang := range rule.Languages {
			rulesByLang[lang] = append(rulesByLang[lang], rule)
		}
	}
	// Process each file
	seen := make(map[string]bool)
	for _, filePath := range filesToScan {
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		p := parser.DetectFromPath(filePath)
		if p == nil {
			continue
		}

		langRules := rulesByLang[p.Language()]
		findings, err := p.Parse(filePath, content, langRules)
		if err != nil {
			continue
		}

		for _, finding := range findings {
			dedupKey := fmt.Sprintf("%s:%d:%s", filePath, finding.Line, finding.RuleID)
			if seen[dedupKey] {
				continue
			}
			seen[dedupKey] = true

			// Find matching rule
			var matchedRule *rules.Rule
			for i := range e.rules {
				if e.rules[i].ID == finding.RuleID {
					matchedRule = &e.rules[i]
					break
				}
			}

			severity := finding.Severity
			if severity == "" && matchedRule != nil {
				severity = matchedRule.Severity
			}
			if severity == "" {
				severity = "WARNING"
			}

			issue := types.Issue{
				File:     filePath,
				Line:     finding.Line,
				Column:   finding.Column,
				EndLine:  finding.EndLine,
				RuleID:   finding.RuleID,
				Severity: severity,
				Message:  finding.Message,
				Source:   "static",
			}

			if matchedRule != nil {
				issue.Title = matchedRule.Name
				issue.Category = matchedRule.Category
				issue.Suggestion = matchedRule.Suggestion

				// Apply category filter
				if cfg.Security && matchedRule.Category != "security" {
					continue
				}
				if cfg.Performance && matchedRule.Category != "performance" {
					continue
				}
			}

			result.Issues = append(result.Issues, issue)
		}
	}

	result.TotalIssues = len(result.Issues)
	for _, issue := range result.Issues {
		switch issue.Severity {
		case "SEVERE":
			result.Severe++
		case "WARNING":
			result.Warning++
		case "INFO":
			result.Info++
		default:
			result.Warning++
		}
	}

	return result, nil
}

func isCodeFile(ext string) bool {
	codeExtensions := map[string]bool{
		".go": true, ".js": true, ".jsx": true, ".ts": true, ".tsx": true,
		".py": true, ".java": true, ".rb": true, ".rs": true, ".c": true,
		".cpp": true, ".h": true, ".php": true, ".swift": true, ".kt": true,
	}
	return codeExtensions[ext]
}
