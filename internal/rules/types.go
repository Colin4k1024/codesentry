package rules

// Rule represents a code review rule
type Rule struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Severity    string   `yaml:"severity"`
	Category    string   `yaml:"category"`
	Languages   []string `yaml:"languages"`
	Suggestion  string   `yaml:"suggestion"`
	Patterns    []Pattern `yaml:"patterns"`
}

// Pattern represents a pattern to match within code
type Pattern struct {
	Type    string `yaml:"type"` // regex, ast
	Pattern string `yaml:"pattern"`
	Comment string `yaml:"comment"`
}
