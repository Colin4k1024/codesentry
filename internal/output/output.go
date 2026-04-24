package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/tabwriter"
	"unsafe"

	"github.com/Colin4k1024/codesentry/internal/types"
)

// Format output format type
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatSarif Format = "sarif"
	FormatGHSL  Format = "ghsl"
	FormatCLang Format = "clang"
)

// ANSI color codes
var colorEnabled = false // 默认关闭，等待自动检测

func init() {
	// 自动检测 isTTY，只有在终端时才启用颜色
	colorEnabled = isTTY()
}

func SetColorEnabled(enabled bool) {
	colorEnabled = enabled
}

// isTTY 检测输出是否连接到终端
func isTTY() bool {
	return isTerminal(os.Stdout.Fd())
}

// isTerminal 使用系统调用检测文件描述符是否为终端
func isTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

// ioctlReadTermios 是 TIOCGETA 的值
const ioctlReadTermios = 0x402C7413

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorBold   = "\033[1m"
	colorReset  = "\033[0m"
)

const toolVersion = "1.0.0"

// FormatForPath infers the output format from a file path extension.
func FormatForPath(outputPath string) Format {
	switch strings.ToLower(filepath.Ext(outputPath)) {
	case ".json":
		return FormatJSON
	case ".sarif":
		return FormatSarif
	default:
		return FormatText
	}
}

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
	case FormatGHSL:
		output, err = formatGHSL(result)
		if err != nil {
			return fmt.Errorf("failed to format GHSL: %w", err)
		}
	case FormatCLang:
		output, err = formatCLang(result)
		if err != nil {
			return fmt.Errorf("failed to format CLang: %w", err)
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
	sb.WriteString(colorize("\n=== CodeSentry Scan Results ===\n", colorBold))
	sb.WriteString(fmt.Sprintf("Files scanned: %d\n", result.FilesScanned))
	sb.WriteString(fmt.Sprintf("Total issues: %d\n", result.TotalIssues))
	sb.WriteString(fmt.Sprintf("  %s: %d\n", colorize("SEVERE", colorRed), result.Severe))
	sb.WriteString(fmt.Sprintf("  %s: %d\n", colorize("WARNING", colorYellow), result.Warning))
	sb.WriteString(fmt.Sprintf("  %s: %d\n", colorize("INFO", colorBlue), result.Info))
	sb.WriteString(fmt.Sprintf("Duration: %v\n", result.Duration))
	sb.WriteString("\n")

	if len(result.Issues) > 0 {
		sb.WriteString(colorize("=== Issues Found ===\n", colorBold))
		for _, issue := range result.Issues {
			severity := issue.Severity
			if severity == "" {
				severity = "WARNING"
			}
			sevColor := colorYellow
			if severity == "SEVERE" {
				sevColor = colorRed
			} else if severity == "INFO" {
				sevColor = colorBlue
			}
			sb.WriteString(fmt.Sprintf("[%s] %s\n", colorize(severity, sevColor), issue.Title))
			sb.WriteString(fmt.Sprintf("  %s: %s:%d:%d\n", colorize("File", colorBlue), colorize(issue.File, colorGreen), issue.Line, issue.Column))
			sb.WriteString(fmt.Sprintf("  %s\n", issue.Message))
			if issue.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("  %s: %s\n", colorize("Suggestion", colorGreen), issue.Suggestion))
			}
			sb.WriteString("\n")
		}
	}

	return []byte(sb.String())
}

func colorize(s string, color string) string {
	if !colorEnabled {
		return s
	}
	return color + s + colorReset
}

func formatSarif(result *types.Result) ([]byte, error) {
	sarif := map[string]interface{}{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":           "CodeSentry",
						"version":        toolVersion,
						"informationUri": "https://github.com/Colin4k1024/codesentry",
						"rules":          []map[string]interface{}{},
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

// formatGHSL formats results in GitHub Security Lab alert-like JSON structure
func formatGHSL(result *types.Result) ([]byte, error) {
	type alert struct {
		Rule        string `json:"rule"`
		Title       string `json:"title"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
		File        string `json:"file"`
		Line        int    `json:"line"`
		Column      int    `json:"column"`
		Suggestion  string `json:"suggestion,omitempty"`
	}
	alerts := []alert{}
	for _, issue := range result.Issues {
		a := alert{
			Rule:        issue.RuleID,
			Title:       issue.Title,
			Severity:    issue.Severity,
			Description: issue.Message,
			File:        issue.File,
			Line:        issue.Line,
			Column:      issue.Column,
		}
		if issue.Suggestion != "" {
			a.Suggestion = issue.Suggestion
		}
		alerts = append(alerts, a)
	}
	return json.MarshalIndent(map[string]interface{}{
		"tool":       "CodeSentry goreview",
		"version":    "1.0.0",
		"filesScanned": result.FilesScanned,
		"totalIssues": result.TotalIssues,
		"summary": map[string]int{
			"severe":  result.Severe,
			"warning": result.Warning,
			"info":    result.Info,
		},
		"alerts": alerts,
	}, "", "  ")
}

// formatCLang formats results in a CLang-Tidy compatible line-by-line format
func formatCLang(result *types.Result) ([]byte, error) {
	var sb strings.Builder
	for _, issue := range result.Issues {
		severity := "warning"
		if issue.Severity == types.SEVERE {
			severity = "error"
		} else if issue.Severity == types.INFO {
			severity = "note"
		}
		sb.WriteString(fmt.Sprintf("%s:%d:%d: %s: %s [%s]\n",
			issue.File, issue.Line, issue.Column,
			severity, issue.Message, issue.RuleID))
		if issue.Suggestion != "" {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", "suggestion", issue.Suggestion))
		}
	}
	return []byte(sb.String()), nil
}

// NewTabWriter creates a new tabwriter
func NewTabWriter(output io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(output, 0, 8, 2, ' ', 0)
}
