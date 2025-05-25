package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateReferenceMarkdown_Nonexistent ensures error when input directory missing.
func TestGenerateReferenceMarkdown_Nonexistent(t *testing.T) {
	outFile := filepath.Join(t.TempDir(), "out.md")
	err := generateReferenceMarkdown("does_not_exist", outFile)
	if err == nil || !strings.Contains(err.Error(), "input directory does not exist") {
		t.Fatalf("expected missing input dir error, got %v", err)
	}
}

// TestGenerateReferenceMarkdown_Basic creates a minimal Go package with
// one const, var, func, and type+method, then verifies the generated
// Markdown contains appropriate sections, TOC entries, and code blocks.
func TestGenerateReferenceMarkdown_Basic(t *testing.T) {
	// Set up temporary input directory
	inDir := t.TempDir()
	src := `
	// Package sample provides example functionality.
	package sample

	// Pi is the mathematical constant.
	const Pi = 3.14

	// Name holds the default name.
	var Name = "test"

	// Hello returns a greeting.
	func Hello(who string) string {
		return "Hello, " + who
	}

	// Counter counts things.
	type Counter int

	// Increment increases the counter.
	func (c *Counter) Increment() {
		*c++
	}
	`
	if err := os.WriteFile(filepath.Join(inDir, "sample.go"), []byte(src), 0644); err != nil {
		t.Fatalf("failed to write sample.go: %v", err)
	}

	// Generate the docs
	outFile := filepath.Join(t.TempDir(), "REF.md")
	if err := generateReferenceMarkdown(inDir, outFile); err != nil {
		t.Fatalf("generateReferenceMarkdown failed: %v", err)
	}

	// Read the output
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	md := string(data)

	// Check title
	if !strings.HasPrefix(md, "# Reference") {
		t.Error("missing title '# Reference'")
	}
	// Check TOC contains sections
	for _, sec := range []string{"Overview", "Constants", "Variables", "Functions", "Types"} {
		if !strings.Contains(md, "- ["+sec+"]") {
			t.Errorf("TOC missing section %q", sec)
		}
	}
	// Check Const Pi entry
	if !strings.Contains(md, "### `Pi`") || !strings.Contains(md, "const Pi = 3.14") {
		t.Error("missing or malformed const Pi documentation")
	}
	// Check Var Name entry
	if !strings.Contains(md, "### `Name`") || !strings.Contains(md, "var Name = \"test\"") {
		t.Error("missing or malformed var Name documentation")
	}
	// Check Func Hello entry
	if !strings.Contains(md, "### `Hello`") || !strings.Contains(md, "func Hello(who string) string") {
		t.Error("missing or malformed func Hello documentation")
	}
	// Check Type Counter entry and its method
	if !strings.Contains(md, "### `Counter`") || !strings.Contains(md, "type Counter int") {
		t.Error("missing or malformed type Counter documentation")
	}
	if !strings.Contains(md, "### `Increment`") || !strings.Contains(md, "func (c *Counter) Increment()") {
		t.Error("missing or malformed method Increment documentation")
	}
}

// TestGenerateAnchor verifies that generateAnchor produces expected Markdown anchors.
func TestGenerateAnchor(t *testing.T) {
	cases := map[string]string{
		"Simple Text":     "simple-text",
		"`Code` Example":  "code-example",
		"Mixed_Case 123":  "mixedcase-123",
		"Trailing - dash": "trailing---dash",
		"***Stars***":     "stars",
	}
	for input, want := range cases {
		if got := generateAnchor(input); got != want {
			t.Errorf("generateAnchor(%q) = %q; want %q", input, got, want)
		}
	}
}

// TestFormatDocText ensures that formatDocText converts paragraphs and preserves code blocks.
func TestFormatDocText(t *testing.T) {
	raw := `This is a paragraph.
It has two lines.

    indented code
More text after code.

- list item one
- list item two
`
	out := formatDocText(raw)
	if !strings.Contains(out, "```text") {
		t.Error("expected code fence for indented block")
	}
	// List items should be present
	if !strings.Contains(out, "- list item one") || !strings.Contains(out, "- list item two") {
		t.Error("list items not formatted correctly")
	}
}
