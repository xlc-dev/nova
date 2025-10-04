package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// generateReferenceMarkdown creates a Markdown documentation file from Go source code comments.
// It parses the package found in inputDir and writes the formatted Markdown to outputFile.
func generateReferenceMarkdown(inputDir, outputFile string) error {
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory does not exist: %s", inputDir)
	}

	outputDirPath := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		// Log warning but continue, os.Create might still work
		log.Printf("Warning: Could not ensure output directory %s exists: %v", outputDirPath, err)
	}

	log.Printf("Generating reference docs from '%s' to '%s'", inputDir, outputFile)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, inputDir, func(fi os.FileInfo) bool {
		// Basic filter to ignore test files
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)

	if err != nil {
		return fmt.Errorf("failed to parse directory %s: %w", inputDir, err)
	}

	// Find the primary package in the directory
	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
		break // Use the first package found
	}
	if pkg == nil {
		return fmt.Errorf("no non-test Go package found in directory: %s", inputDir)
	}

	// Extract documentation (exported symbols only by default)
	docPkg := doc.New(pkg, pkg.Name, doc.AllDecls)

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer out.Close()

	// Add header
	header := "{{ title: Nova - Reference }}\n\n{{ include-block: doc.html markdown=\"true\" }}\n\n"
	_, err = out.WriteString(header)
	if err != nil {
		return fmt.Errorf("failed to write header to output file: %w", err)
	}

	var contentBuf bytes.Buffer

	title := "# Reference\n\n"

	// Package Documentation (if any)
	if docPkg.Doc != "" {
		contentBuf.WriteString(fmt.Sprintf("## %s\n\n", "Overview"))
		contentBuf.WriteString(formatDocText(docPkg.Doc))
		contentBuf.WriteString("\n\n")
	}

	// Write Title
	_, err = out.WriteString(title)
	if err != nil {
		return fmt.Errorf("failed to write title to output file: %w", err)
	}

	// Generate and Write Table of Contents
	toc := generateTOC(docPkg)
	if toc != "" {
		_, err = out.WriteString(toc)
		if err != nil {
			return fmt.Errorf("failed to write TOC to output file: %w", err)
		}
		_, err = out.WriteString("\n\n") // Add separation after TOC
		if err != nil {
			return fmt.Errorf("failed to write separator after TOC: %w", err)
		}
	}

	// Constants
	if len(docPkg.Consts) > 0 {
		contentBuf.WriteString("## Constants\n\n")
		for _, c := range docPkg.Consts {
			writeDocItem(&contentBuf, fset, c.Doc, c.Names, c.Decl, 3)
		}
	}

	// Variables
	if len(docPkg.Vars) > 0 {
		contentBuf.WriteString("## Variables\n\n")
		for _, v := range docPkg.Vars {
			writeDocItem(&contentBuf, fset, v.Doc, v.Names, v.Decl, 3)
		}
	}

	// Functions
	if len(docPkg.Funcs) > 0 {
		contentBuf.WriteString("## Functions\n\n")
		for _, f := range docPkg.Funcs {
			writeDocItem(&contentBuf, fset, f.Doc, []string{f.Name}, f.Decl, 3)
		}
	}

	// Types
	if len(docPkg.Types) > 0 {
		contentBuf.WriteString("## Types\n\n")
		for _, t := range docPkg.Types {
			writeDocItem(&contentBuf, fset, t.Doc, []string{t.Name}, t.Decl, 3)

			if len(t.Consts) > 0 {
				contentBuf.WriteString("#### Associated Constants\n\n")
				for _, c := range t.Consts {
					writeDocItem(&contentBuf, fset, c.Doc, c.Names, c.Decl, 4)
				}
			}
			if len(t.Vars) > 0 {
				contentBuf.WriteString("#### Associated Variables\n\n")
				for _, v := range t.Vars {
					writeDocItem(&contentBuf, fset, v.Doc, v.Names, v.Decl, 4)
				}
			}
			if len(t.Funcs) > 0 {
				contentBuf.WriteString("#### Associated Functions\n\n")
				for _, f := range t.Funcs {
					writeDocItem(&contentBuf, fset, f.Doc, []string{f.Name}, f.Decl, 4)
				}
			}
			if len(t.Methods) > 0 {
				contentBuf.WriteString("#### Methods\n\n")
				for _, m := range t.Methods {
					writeDocItem(&contentBuf, fset, m.Doc, []string{m.Name}, m.Decl, 4)
				}
			}
		}
	}

	// Write Main Content
	_, err = contentBuf.WriteTo(out)
	if err != nil {
		return fmt.Errorf("failed to write content buffer to output file: %w", err)
	}

	// Add footer
	footer := "{{ endinclude }}"
	_, err = out.WriteString(footer)
	if err != nil {
		return fmt.Errorf("failed to write footer to output file: %w", err)
	}

	log.Printf("Successfully generated reference docs with TOC to %s", outputFile)
	return nil
}

func generateTOC(docPkg *doc.Package) string {
	var tocBuf bytes.Buffer
	tocBuf.WriteString("## Table of Contents\n\n")

	hasContent := false

	if docPkg.Doc != "" {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Overview", "overview"))
		hasContent = true
	}

	if len(docPkg.Consts) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Constants", "constants"))
		hasContent = true
	}

	if len(docPkg.Vars) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Variables", "variables"))
		hasContent = true
	}

	if len(docPkg.Funcs) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Functions", "functions"))
		hasContent = true
	}

	if len(docPkg.Types) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Types", "types"))
		hasContent = true
		for _, t := range docPkg.Types {
			typeTitle := fmt.Sprintf("`%s`", t.Name)
			typeAnchor := generateAnchor(t.Name)
			tocBuf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", typeTitle, typeAnchor))

			if len(t.Consts) > 0 || len(t.Vars) > 0 || len(t.Funcs) > 0 || len(t.Methods) > 0 {
				if len(t.Consts) > 0 {
					tocBuf.WriteString(fmt.Sprintf("    - [Associated Constants](#%s-constants)\n", typeAnchor))
				}
				if len(t.Vars) > 0 {
					tocBuf.WriteString(fmt.Sprintf("    - [Associated Variables](#%s-variables)\n", typeAnchor))
				}
				if len(t.Funcs) > 0 {
					tocBuf.WriteString(fmt.Sprintf("    - [Associated Functions](#%s-functions)\n", typeAnchor))
				}
				if len(t.Methods) > 0 {
					tocBuf.WriteString(fmt.Sprintf("    - [Methods](#%s-methods)\n", typeAnchor))
				}
			}
		}
	}

	if !hasContent {
		return ""
	}

	return tocBuf.String()
}

// writeDocItem formats a single documentation item (const, var, func, type, method)
// and writes it to the content buffer.
func writeDocItem(contentBuf *bytes.Buffer, fset *token.FileSet, docComment string, names []string, decl ast.Node, level int) {
	displayNames := strings.Join(names, ", ")
	itemTitle := fmt.Sprintf("`%s`", displayNames)
	itemAnchor := generateAnchor(displayNames)

	// Create a clean anchor tag that is compatible with most Markdown renderers
	fmt.Fprintf(contentBuf, "<a id=\"%s\"></a>\n", itemAnchor)
	fmt.Fprintf(contentBuf, "%s %s\n\n", strings.Repeat("#", level), itemTitle)

	// Print the declaration (signature) using go/printer
	var declBuf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
	err := cfg.Fprint(&declBuf, fset, decl)
	if err != nil {
		log.Printf("Warning: Failed to print declaration for %s: %v", displayNames, err)
		contentBuf.WriteString("```go\n// Error printing declaration\n```\n\n")
	} else {
		contentBuf.WriteString("```go\n")
		contentBuf.Write(declBuf.Bytes())
		contentBuf.WriteString("\n```\n\n")
	}

	// Write the documentation comment
	if docComment != "" {
		contentBuf.WriteString(formatDocText(docComment))
		contentBuf.WriteString("\n\n")
	}
}

func formatDocText(text string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(text)

	// Use doc.ToHTML to handle formatting, then convert to Markdown
	var buf bytes.Buffer
	doc.ToHTML(&buf, trimmed, nil)
	html := buf.String()

	// Basic conversion from HTML to Markdown
	md := strings.ReplaceAll(html, "<p>", "")
	md = strings.ReplaceAll(md, "</p>", "\n\n")
	md = strings.ReplaceAll(md, "<pre>", "```go\n")
	md = strings.ReplaceAll(md, "</pre>", "\n```\n\n")
	md = strings.ReplaceAll(md, "<code>", "`")
	md = strings.ReplaceAll(md, "</code>", "`")
	md = strings.ReplaceAll(md, "<ul>", "")
	md = strings.ReplaceAll(md, "</ul>", "")
	md = strings.ReplaceAll(md, "<li>", "- ")
	md = strings.ReplaceAll(md, "</li>", "\n")

	// Clean up extra newlines
	md = strings.TrimSpace(md)
	md = regexp.MustCompile(`\n{3,}`).ReplaceAllString(md, "\n\n")

	return md
}

func generateAnchor(text string) string {
	// Remove backticks first
	text = strings.ReplaceAll(text, "`", "")
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")

	// Remove any characters that are not letters, numbers, or hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	text = reg.ReplaceAllString(text, "")

	// Replace multiple hyphens with a single hyphen
	text = regexp.MustCompile("-+").ReplaceAllString(text, "-")

	// Ensure it doesn't start or end with a hyphen
	text = strings.Trim(text, "-")
	return text
}
