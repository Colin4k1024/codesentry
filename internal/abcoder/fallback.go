package abcoder

import (
	"strings"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

// FallbackHandler provides fallback suggestions when abcoder is unavailable
type FallbackHandler struct {
	// ruleFixes maps rule IDs to fallback fix templates
	ruleFixes map[string]*FallbackFix
}

// FallbackFix represents a fallback fix template
type FallbackFix struct {
	Before   string
	After    string
	Pattern  string // Description of the pattern to fix
	Template string // Suggestion template
}

// NewFallbackHandler creates a new FallbackHandler with default mappings
func NewFallbackHandler() *FallbackHandler {
	return &FallbackHandler{
		ruleFixes: map[string]*FallbackFix{
			"SQL_INJECTION": {
				Pattern:  "String concatenation in SQL query",
				Template: "Use parameterized queries: db.Query(\"SELECT ... WHERE id=?\", id)",
				Before:   `"SELECT ... WHERE id=" + id`,
				After:    `"SELECT ... WHERE id=?", id`,
			},
			"HARDCODED_SECRET": {
				Pattern:  "Hardcoded sensitive value",
				Template: "Use environment variables: os.Getenv(\"KEY\") or secrets manager",
				Before:   `apiKey := "YOUR_API_KEY_HERE"`,
				After:    `apiKey := os.Getenv("API_KEY")`,
			},
			"EXECUTION": {
				Pattern:  "Shell command with string concatenation",
				Template: "Use exec.Command with separate arguments: exec.Command(\"cmd\", arg1, arg2)",
				Before:   `exec.Command("ls " + path)`,
				After:    `exec.Command("ls", path)`,
			},
			"XSS": {
				Pattern:  "Direct HTML insertion",
				Template: "Sanitize output or use safe DOM APIs",
				Before:   `element.innerHTML = userInput`,
				After:    `element.textContent = userInput`,
			},
			"DANGEROUS_EVAL": {
				Pattern:  "Use of eval() with external input",
				Template: "Avoid eval(); use safer alternatives like JSON.parse for data",
				Before:   `eval(userInput)`,
				After:    `JSON.parse(userInput)`,
			},
			"PATH_TRAVERSAL": {
				Pattern:  "Direct path concatenation with user input",
				Template: "Use path.Clean() and validate path components",
				Before:   `os.Open(userPath)`,
				After:    `os.Open(path.Clean(userPath))`,
			},
			"DESERIALIZATION": {
				Pattern:  "Unsafe deserialization",
				Template: "Use safe deserialization; avoid pickle/ marshal with untrusted data",
				Before:   `pickle.loads(data)`,
				After:    `json.loads(data)`,
			},
		},
	}
}

// GetFix returns a fallback fix for the given rule
func (h *FallbackHandler) GetFix(ruleID string) *FallbackFix {
	if fix, ok := h.ruleFixes[ruleID]; ok {
		return fix
	}
	return nil
}

// BuildSuggestion builds a fallback suggestion from a rule
func (h *FallbackHandler) BuildSuggestion(rule *rules.Rule) string {
	if rule == nil {
		return "Review and fix this code issue"
	}

	// If rule has its own suggestion, use it
	if rule.Suggestion != "" {
		return rule.Suggestion
	}

	// Try to find a matching fallback fix
	if fix := h.GetFix(rule.ID); fix != nil {
		return fix.Template
	}

	// Default suggestion based on category
	switch rule.Category {
	case "security":
		return "Security issue detected: review and fix according to security best practices"
	case "performance":
		return "Performance issue detected: consider optimizing this code"
	default:
		return "Code issue detected: review and fix"
	}
}

// IsFallbackNeeded determines if fallback should be used
func IsFallbackNeeded(file string, err error) bool {
	// If file is not a Go file, we need fallback (abcoder only supports Go AST)
	if !IsAvailable(file) {
		return true
	}

	// If there was an error parsing, fallback to suggestion
	if err != nil {
		return true
	}

	return false
}

// FormatFallback formats a fallback suggestion as a readable string
func (h *FallbackHandler) FormatFallback(ruleID, suggestion string) string {
	fix := h.GetFix(ruleID)
	if fix == nil {
		return suggestion
	}

	var sb strings.Builder
	sb.WriteString("【修复建议】\n\n")
	sb.WriteString("问题: " + fix.Pattern + "\n")
	sb.WriteString("建议: " + fix.Template + "\n\n")

	if suggestion != "" && suggestion != fix.Template {
		sb.WriteString("原始建议: " + suggestion + "\n")
	}

	return sb.String()
}
