package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vtbarreto/CLinicius/internal/analyzer"
)

const sampleGoFile = `package sample

import (
	"fmt"
	"os"
	"strings"
)

func hello() { fmt.Println(strings.TrimSpace(os.Getenv("HOME"))) }
`

func writeTempGoFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestParseFileImports(t *testing.T) {
	path := writeTempGoFile(t, sampleGoFile)

	imports, err := analyzer.ParseFileImports(path)
	if err != nil {
		t.Fatalf("ParseFileImports() error = %v", err)
	}

	want := map[string]bool{"fmt": true, "os": true, "strings": true}
	if len(imports) != len(want) {
		t.Fatalf("got %d imports, want %d: %v", len(imports), len(want), imports)
	}
	for _, imp := range imports {
		if !want[imp] {
			t.Errorf("unexpected import %q", imp)
		}
	}
}

func TestParseFileImports_InvalidSyntax(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.go")
	if err := os.WriteFile(path, []byte("package bad\nimport (bad syntax"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := analyzer.ParseFileImports(path)
	if err == nil {
		t.Fatal("expected error for invalid Go syntax, got nil")
	}
}

func TestParseFileImports_NoImports(t *testing.T) {
	src := "package noImports\n\nfunc Foo() {}\n"
	path := writeTempGoFile(t, src)

	imports, err := analyzer.ParseFileImports(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(imports) != 0 {
		t.Errorf("expected 0 imports, got %v", imports)
	}
}

func TestParseDirImports(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"a.go": "package pkg\nimport \"fmt\"\nfunc A() { fmt.Println() }",
		"b.go": "package pkg\nimport \"os\"\nfunc B() { _ = os.Stderr }",
	}
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result, err := analyzer.ParseDirImports(dir)
	if err != nil {
		t.Fatalf("ParseDirImports() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 FileImports entries, got %d", len(result))
	}

	found := make(map[string]bool)
	for _, fi := range result {
		for _, imp := range fi.Imports {
			found[imp] = true
		}
	}
	if !found["fmt"] || !found["os"] {
		t.Errorf("expected both fmt and os imports, got %v", found)
	}
}
