package rules

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadRules loads all rule YAML files from the specified directory
func LoadRules(dir string) []Rule {
	var rules []Rule

	// Walk the rules directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process YAML files
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		// Read the file
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not read rule file %s: %v\n", path, err)
			return nil
		}

		// Parse YAML
		var rule Rule
		if err := yaml.Unmarshal(data, &rule); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not parse rule file %s: %v\n", path, err)
			return nil
		}

		// Skip rules without an ID
		if rule.ID == "" {
			fmt.Fprintf(os.Stderr, "Warning: Rule file %s has no ID, skipping\n", path)
			return nil
		}

		rules = append(rules, rule)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking rules directory: %v\n", err)
		return nil
	}

	return rules
}

// FilterByLanguage returns rules that support the given language
func FilterByLanguage(lang string, rules []Rule) []Rule {
	var filtered []Rule
	for _, rule := range rules {
		for _, supportedLang := range rule.Languages {
			if supportedLang == lang {
				filtered = append(filtered, rule)
				break
			}
		}
	}
	return filtered
}
