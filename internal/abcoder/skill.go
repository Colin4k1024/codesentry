package abcoder

import (
	"context"
	"encoding/json"
	"fmt"
)

// SkillInput represents the input for the code-reviewer skill
type SkillInput struct {
	Task          string       `json:"task"`
	Vulnerability string       `json:"vulnerability"`
	File          string       `json:"file"`
	Line          int          `json:"line"`
	RuleID        string       `json:"rule_id"`
	Suggestion    string       `json:"suggestion"`
	Context       *CodeContext `json:"context,omitempty"`
}

// SkillOutput represents the output from the code-reviewer skill
type SkillOutput struct {
	Before      string   `json:"before"`
	After       string   `json:"after"`
	Explanation string   `json:"explanation"`
	Confidence  float64  `json:"confidence"`
	Warnings    []string `json:"warnings,omitempty"`
}

// SkillAgent provides integration with Claude Code's code-reviewer skill
type SkillAgent struct {
	bridge *Bridge
}

// NewSkillAgent creates a new SkillAgent
func NewSkillAgent(bridge *Bridge) *SkillAgent {
	return &SkillAgent{bridge: bridge}
}

// GenerateFix generates a fix suggestion using the code context
// This is a simplified implementation that uses structured output
// In production, this would invoke Claude Code's skill via MCP
func (s *SkillAgent) GenerateFix(ctx context.Context, ruleID, suggestion, file string, line int) (*SkillOutput, error) {
	// Get code context if available
	var codeCtx *CodeContext
	if s.bridge != nil && IsAvailable(file) {
		var err error
		codeCtx, err = s.bridge.GetContext(file, line)
		if err != nil {
			// If context retrieval fails, continue without context
			codeCtx = nil
		}
	}

	// Build skill input
	input := SkillInput{
		Task:          "generate_fix",
		Vulnerability: ruleID,
		File:          file,
		Line:          line,
		RuleID:        ruleID,
		Suggestion:    suggestion,
		Context:       codeCtx,
	}

	// In a full implementation, this would:
	// 1. Serialize input to JSON
	// 2. Invoke Claude Code skill via MCP
	// 3. Parse the skill output

	// For now, we return a structured output that can be used by external tools
	return s.buildStructuredOutput(input)
}

// buildStructuredOutput creates a structured fix suggestion
// This will be replaced by actual skill invocation
func (s *SkillAgent) buildStructuredOutput(input SkillInput) (*SkillOutput, error) {
	output := &SkillOutput{
		Confidence: 0.7, // Default confidence
		Warnings:   []string{},
	}

	// Build explanation based on rule type
	switch input.RuleID {
	case "SQL_INJECTION":
		output.Before = "query := \"SELECT * FROM users WHERE id=\" + id"
		output.After = "query := \"SELECT * FROM users WHERE id=?\", id"
		output.Explanation = "使用参数化查询避免 SQL 注入风险"
	case "HARDCODED_SECRET":
		output.Before = "apiKey := \"YOUR_API_KEY_HERE\""
		output.After = "apiKey := os.Getenv(\"API_KEY\")"
		output.Explanation = "使用环境变量存储敏感信息，避免硬编码"
	case "EXECUTION":
		output.Before = "exec.Command(\"ls \" + userInput)"
		output.After = "exec.Command(\"ls\", userInput)"
		output.Explanation = "使用参数化命令执行，避免 shell 注入"
	default:
		// Use suggestion as explanation
		output.Before = "<需要修复的代码>"
		output.After = "<修复后的代码>"
		output.Explanation = input.Suggestion
		if output.Explanation == "" {
			output.Explanation = "建议审查并修复此代码问题"
		}
	}

	// Add context warnings if available
	if input.Context != nil {
		if len(input.Context.FunctionCalls) > 3 {
			output.Warnings = append(output.Warnings, "函数有多个外部调用，可能需要更全面的修复")
		}
		if input.Context.FunctionName == "main" {
			output.Warnings = append(output.Warnings, "问题出现在入口函数，需要特别谨慎")
		}
	}

	return output, nil
}

// ToJSON converts the skill output to JSON for external consumption
func (s *SkillOutput) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// FormatFix formats a fix suggestion as a readable string
func (s *SkillOutput) FormatFix() string {
	format := `
=== 修复建议 ===

修复前:
%s

修复后:
%s

原因: %s
置信度: %.0f%%
`
	if len(s.Warnings) > 0 {
		format += "\n警告:\n"
		for _, w := range s.Warnings {
			format += fmt.Sprintf("  - %s\n", w)
		}
	}
	return fmt.Sprintf(format, s.Before, s.After, s.Explanation, s.Confidence*100)
}
