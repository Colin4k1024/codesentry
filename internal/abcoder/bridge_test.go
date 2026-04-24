package abcoder

import (
	"context"
	"testing"
)

func TestIsAvailable(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.GO", true},
		{"test.js", false},
		{"test.py", false},
		{"test.java", false},
		{"test.rs", false},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			result := IsAvailable(tt.file)
			if result != tt.expected {
				t.Errorf("IsAvailable(%q) = %v, want %v", tt.file, result, tt.expected)
			}
		})
	}
}

func TestNewBridge(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Errorf("NewBridge with valid path returned error: %v", err)
	}
	if bridge == nil {
		t.Error("NewBridge returned nil bridge")
	}
}

func TestBridgeParse(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	ctx := context.Background()
	err = bridge.Parse(ctx)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
}

func TestBridgeGetContext(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	ctx := context.Background()
	if err := bridge.Parse(ctx); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Try to get context for a non-existent file
	_, err = bridge.GetContext("nonexistent.go", 10)
	if err == nil {
		t.Error("GetContext for non-existent file should return error")
	}
}

func TestBridgeGetContextNotParsed(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	_, err = bridge.GetContext("test.go", 10)
	if err == nil {
		t.Error("GetContext without Parse should return error")
	}
}

func TestBridgeGetFunction(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	ctx := context.Background()
	if err := bridge.Parse(ctx); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// GetFunction with non-existent identity should fail
	_, err = bridge.GetFunction("nonexistent")
	if err == nil {
		t.Error("GetFunction with non-existent identity should return error")
	}
}

func TestBridgeGetVariable(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	ctx := context.Background()
	if err := bridge.Parse(ctx); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// GetVariable with non-existent identity should fail
	_, err = bridge.GetVariable("nonexistent")
	if err == nil {
		t.Error("GetVariable with non-existent identity should return error")
	}
}

func TestBridgeGetCallChain(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	ctx := context.Background()
	if err := bridge.Parse(ctx); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// GetCallChain with non-existent identity should fail
	_, err = bridge.GetCallChain("nonexistent")
	if err == nil {
		t.Error("GetCallChain with non-existent identity should return error")
	}
}

func TestParseIdentity(t *testing.T) {
	// Test parseIdentity with different formats
	tests := []struct {
		nodeID string
		name   string
	}{
		{"Name", "Name"},
		{"PkgPath#Name", "Name"},
		{"ModPath?PkgPath#Name", "Name"},
	}

	for _, tt := range tests {
		t.Run(tt.nodeID, func(t *testing.T) {
			identity := parseIdentity(tt.nodeID)
			if identity.Name != tt.name {
				t.Errorf("parseIdentity(%q).Name = %q, want %q", tt.nodeID, identity.Name, tt.name)
			}
		})
	}
}

func TestGetTypeName(t *testing.T) {
	bridge, err := NewBridge(".")
	if err != nil {
		t.Fatalf("NewBridge failed: %v", err)
	}

	// Test getTypeName with nil
	result := bridge.getTypeName(nil)
	if result != "" {
		t.Errorf("getTypeName(nil) = %q, want empty string", result)
	}
}
