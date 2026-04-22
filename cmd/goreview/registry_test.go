package main

import (
	"testing"

	"github.com/Colin4k1024/codesentry/internal/parser"
)

func TestParserRegistry(t *testing.T) {
	// Test that all 11 language parsers are registered
	languages := parser.List()

	expectedLangs := []string{
		"cpp", "go", "java", "javascript",
		"kotlin", "php", "python", "ruby",
		"rust", "swift", "typescript",
	}

	if len(languages) != len(expectedLangs) {
		t.Errorf("expected %d languages, got %d: %v", len(expectedLangs), len(languages), languages)
	}

	for _, want := range expectedLangs {
		found := false
		for _, got := range languages {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected language %q not found in %v", want, languages)
		}
	}
}

func TestDetectFromPath(t *testing.T) {
	tests := []struct {
		path     string
		wantLang string
	}{
		{"file.go", "go"},
		{"file.py", "python"},
		{"file.js", "javascript"},
		{"file.ts", "typescript"},
		{"file.java", "java"},
		{"file.rb", "ruby"},
		{"file.rs", "rust"},
		{"file.cpp", "cpp"},
		{"file.php", "php"},
		{"file.swift", "swift"},
		{"file.kt", "kotlin"},
		{"file.tsx", "typescript"},
		{"file.jsx", "javascript"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			p := parser.DetectFromPath(tt.path)
			if p == nil {
				t.Fatalf("DetectFromPath(%q) returned nil", tt.path)
			}
			if got := p.Language(); got != tt.wantLang {
				t.Errorf("DetectFromPath(%q).Language() = %q, want %q", tt.path, got, tt.wantLang)
			}
		})
	}
}

func TestDetectFromPathUnknown(t *testing.T) {
	tests := []string{
		"file.unknown",
		"file.txt",
		"file.csv",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			p := parser.DetectFromPath(path)
			if p != nil {
				t.Errorf("DetectFromPath(%q) = %v, want nil", path, p)
			}
		})
	}
}
