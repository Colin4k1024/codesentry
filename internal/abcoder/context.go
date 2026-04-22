// Package abcoder provides integration with cloudwego/abcoder for code context understanding
package abcoder

import (
	"github.com/Colin4k1024/codesentry/internal/types"
)

// FixSuggestion represents a suggested fix for a code issue
type FixSuggestion struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Explanation string `json:"explanation"`
	Confidence   float64 `json:"confidence"`
}

// IssueWithContext represents an issue with additional code context
type IssueWithContext struct {
	Issue   types.Issue
	Context *CodeContext
	Fix     *FixSuggestion
}

// ToIssue converts IssueWithContext back to types.Issue with fix info
func (i *IssueWithContext) ToIssue() types.Issue {
	issue := i.Issue
	if i.Fix != nil {
		issue.Suggestion = i.Fix.Explanation
	}
	return issue
}
