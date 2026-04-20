package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/goreview/goreview/internal/types"
)

// Format output format type
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatSarif Format = "sarif"
)

// Write writes the result in the specified format
func Write(result *types.Result, format Format, outputPath string) error {
	var err error
	var output []byte

	switch format {
	case FormatJSON:
		output, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	case FormatText:
		output = formatText(result)
	case FormatSarif:
		output, err = formatSarif(result)
		if err != nil {
			return fmt.Errorf("failed to format SARIF: %w", err)
		}
	default:
		output = formatText(result)
	}

	// Write to file or stdout
	if outputPath != "" {
		err = os.WriteFile(outputPath, output, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		fmt.Print(string(output))
	}

	return nil
}

func formatText(result *types.Result) []byte {
	var sb strings.Builder
	sb.WriteString("\n=== GoReview Scan Results ===\n")
	sb.WriteString(fmt.Sprintf("Files scanned: %d\n", result.FilesScanned))
	sb.WriteString(fmt.Sprintf("Total issues: %d\n", result.TotalIssues))
	sb.WriteString(fmt.Sprintf("  SEVERE:   %d\n", result.Severe))
	sb.WriteString(fmt.Sprintf("  WARNING:  %d\n", result.Warning))
	sb.WriteString(fmt.Sprintf("  INFO:     %d\n", result.Info))
	sb.WriteString(fmt.Sprintf("Duration: %v\n", result.Duration))
	sb.WriteString("\n")

	if len(result.Issues) > 0 {
		sb.WriteString("=== Issues Found ===\n")
		for _, issue := range result.Issues {
			severity := issue.Severity
			if severity == "" {
				severity = "WARNING"
			}
			sb.WriteString(fmt.Sprintf("[%s] %s\n", severity, issue.Title))
			sb.WriteString(fmt.Sprintf("  File: %s:%d:%d\n", issue.File, issue.Line, issue.Column))
			sb.WriteString(fmt.Sprintf("  %s\n", issue.Message))
			if issue.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("  Suggestion: %s\n", issue.Suggestion))
			}
			sb.WriteString("\n")
		}
	}

	return []byte(sb.String())
}

func formatSarif(result *types.Result) ([]byte, error) {
	sarif := map[string]interface{}{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":            "GoReview",
						"version":         "0.3.0",
						"informationUri":  "https://github.com/goreview/goreview",
						"rules":           []map[string]interface{}{},
					},
				},
				"results": []map[string]interface{}{},
			},
		},
	}

	// Convert issues to SARIF format
	results := make([]map[string]interface{}, 0, len(result.Issues))
	rules := make(map[string]map[string]interface{})

	for _, issue := range result.Issues {
		ruleID := issue.RuleID
		if ruleID == "" {
			ruleID = "UNKNOWN"
		}

		// Add rule if not already added
		if _, exists := rules[ruleID]; !exists {
			rules[ruleID] = map[string]interface{}{
				"id":               ruleID,
				"name":             issue.Title,
				"shortDescription": map[string]interface{}{"text": issue.Title},
				"fullDescription":  map[string]interface{}{"text": issue.Message},
				"defaultLevel":     mapSeverityToSarifLevel(issue.Severity),
			}
		}

		resultEntry := map[string]interface{}{
			"ruleId": ruleID,
			"level":  mapSeverityToSarifLevel(issue.Severity),
			"message": map[string]interface{}{
				"text": issue.Message,
			},
			"locations": []map[string]interface{}{
				{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]interface{}{
							"uri": issue.File,
						},
						"region": map[string]interface{}{
							"startLine":   issue.Line,
							"startColumn": issue.Column,
						},
					},
				},
			},
		}
		results = append(results, resultEntry)
	}

	// Update runs with results and rules
	if len(sarif["runs"].([]map[string]interface{})) > 0 {
		sarif["runs"].([]map[string]interface{})[0]["results"] = results

		driver := sarif["runs"].([]map[string]interface{})[0]["tool"].(map[string]interface{})["driver"].(map[string]interface{})
		rulesList := make([]map[string]interface{}, 0, len(rules))
		for _, r := range rules {
			rulesList = append(rulesList, r)
		}
		driver["rules"] = rulesList
	}

	return json.MarshalIndent(sarif, "", "  ")
}

func mapSeverityToSarifLevel(severity string) string {
	switch severity {
	case types.SEVERE:
		return "error"
	case types.WARNING:
		return "warning"
	case types.INFO:
		return "note"
	default:
		return "warning"
	}
}

// NewTabWriter creates a new tabwriter
func NewTabWriter(output io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(output, 0, 8, 2, ' ', 0)
}
