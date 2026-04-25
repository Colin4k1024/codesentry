package parser

import (
	"testing"

	"github.com/Colin4k1024/codesentry/internal/rules"
)

type mockParser struct {
	lang string
	exts []string
}

func (m *mockParser) Language() string     { return m.lang }
func (m *mockParser) Extensions() []string { return m.exts }
func (m *mockParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]Finding, error) {
	return nil, nil
}

func TestRegisterAndGet(t *testing.T) {
	// Register a mock parser
	p := &mockParser{lang: "testlang", exts: []string{".test"}}
	Register(p)

	got, ok := Get("testlang")
	if !ok {
		t.Fatal("Get(testlang) returned false, want true")
	}
	if got.Language() != "testlang" {
		t.Errorf("got.Language() = %q, want %q", got.Language(), "testlang")
	}

	// Clean up
	delete(registry, "testlang")
}

func TestGetNotFound(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Error("Get(nonexistent) returned true, want false")
	}
}

func TestList(t *testing.T) {
	before := List()

	// Register a unique mock
	p := &mockParser{lang: "testlang2", exts: []string{".tst"}}
	Register(p)
	defer delete(registry, "testlang2")

	after := List()
	if len(after) <= len(before) {
		t.Error("List() did not grow after registration")
	}

	found := false
	for _, lang := range after {
		if lang == "testlang2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("testlang2 not found in List()")
	}
}

func TestDetectFromPath(t *testing.T) {
	// Register mock parsers
	Register(&mockParser{lang: "foo", exts: []string{".foo"}})
	Register(&mockParser{lang: "bar", exts: []string{".bar"}})
	defer func() {
		delete(registry, "foo")
		delete(registry, "bar")
	}()

	if p := DetectFromPath("/path/to/script.foo"); p == nil || p.Language() != "foo" {
		t.Errorf("DetectFromPath script.foo = %v (lang=%q), want lang=foo", p, p.Language())
	}
	if p := DetectFromPath("/path/to/script.bar"); p == nil || p.Language() != "bar" {
		t.Errorf("DetectFromPath script.bar = %v (lang=%q), want lang=bar", p, p.Language())
	}
	if p := DetectFromPath("/path/to/script.unknown"); p != nil {
		t.Errorf("DetectFromPath script.unknown = %v, want nil", p)
	}
}

func TestHasExt(t *testing.T) {
	if !hasExt("test.js", ".js") {
		t.Error("hasExt(test.js, .js) = false, want true")
	}
	if hasExt("test.js", ".ts") {
		t.Error("hasExt(test.js, .ts) = true, want false")
	}
	if !hasExt(".js", ".js") {
		t.Error("hasExt(.js, .js) = false, want true")
	}
	if hasExt(".js", ".jsx") {
		t.Error("hasExt(.js, .jsx) = true, want false")
	}
	if hasExt("test", ".test") {
		t.Error("hasExt(test, .test) = true, want false (too short)")
	}
}
