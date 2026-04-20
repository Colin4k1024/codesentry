package types

import "time"

// Severity constants
const (
	SEVERE   = "SEVERE"
	WARNING  = "WARNING"
	INFO     = "INFO"
)

// Issue represents a code issue found during scanning
type Issue struct {
	ID        string `json:"id"`
	Severity  string `json:"severity"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	EndLine   int    `json:"end_line"`
	RuleID    string `json:"rule_id"`
	Category  string `json:"category"`
	Suggestion string `json:"suggestion"`
	Source    string `json:"source"`
}

// Result represents the overall scan result
type Result struct {
	TotalFiles    int      `json:"total_files"`
	TotalIssues   int      `json:"total_issues"`
	Severe        int      `json:"severe"`
	Warning       int      `json:"warning"`
	Info          int      `json:"info"`
	Duration      time.Duration `json:"duration"`
	Timestamp     time.Time `json:"timestamp"`
	Issues        []Issue  `json:"issues"`
	FilesScanned  int       `json:"files_scanned"`
}
