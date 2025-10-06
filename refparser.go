package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"html"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// generateReferenceMarkdown generates a Markdown reference file for the package
// found in inputDir and writes it to outputFile.
func generateReferenceMarkdown(inputDir, outputFile string) error {
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory does not exist: %s", inputDir)
	}

	outputDirPath := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		log.Printf("Warning: Could not ensure output directory %s exists: %v", outputDirPath, err)
	}

	log.Printf("Generating reference docs from '%s' to '%s'", inputDir, outputFile)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, inputDir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse directory %s: %w", inputDir, err)
	}

	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
		break
	}
	if pkg == nil {
		return fmt.Errorf("no non-test Go package found in directory: %s", inputDir)
	}

	docPkg := doc.New(pkg, pkg.Name, doc.AllDecls)

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer out.Close()

	header := "{{ title: Nova - Reference }}\n\n{{ include-block: doc.html markdown=\"true\" }}\n\n"
	if _, err := out.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header to output file: %w", err)
	}

	var content bytes.Buffer
	content.WriteString("# Reference\n\n")

	toc := generateTOC(docPkg)
	if toc != "" {
		content.WriteString(toc)
		content.WriteString("\n\n")
	}

	if docPkg.Doc != "" {
		content.WriteString("## Overview\n\n")
		content.WriteString(formatDocText(docPkg.Doc))
		content.WriteString("\n\n")
	}

	if len(docPkg.Consts) > 0 {
		content.WriteString("## Constants\n\n")
		for _, c := range docPkg.Consts {
			writeDocItem(&content, fset, c.Doc, c.Names, c.Decl, 3)
		}
	}

	if len(docPkg.Vars) > 0 {
		content.WriteString("## Variables\n\n")
		for _, v := range docPkg.Vars {
			writeDocItem(&content, fset, v.Doc, v.Names, v.Decl, 3)
		}
	}

	if len(docPkg.Funcs) > 0 {
		content.WriteString("## Functions\n\n")
		for _, f := range docPkg.Funcs {
			writeDocItem(&content, fset, f.Doc, []string{f.Name}, f.Decl, 3)
		}
	}

	if len(docPkg.Types) > 0 {
		content.WriteString("## Types\n\n")
		for _, t := range docPkg.Types {
			fmt.Fprintf(&content, "### `%s`\n\n", t.Name)
			printDeclaration(&content, fset, t.Decl, t.Name)
			if t.Doc != "" {
				content.WriteString(formatDocText(t.Doc))
				content.WriteString("\n\n")
			}

			if len(t.Consts) > 0 {
				content.WriteString("#### Associated Constants\n\n")
				for _, c := range t.Consts {
					writeDocItem(&content, fset, c.Doc, c.Names, c.Decl, 4)
				}
			}
			if len(t.Vars) > 0 {
				content.WriteString("#### Associated Variables\n\n")
				for _, v := range t.Vars {
					writeDocItem(&content, fset, v.Doc, v.Names, v.Decl, 4)
				}
			}
			if len(t.Funcs) > 0 {
				content.WriteString("#### Associated Functions\n\n")
				for _, f := range t.Funcs {
					writeDocItem(&content, fset, f.Doc, []string{f.Name}, f.Decl, 4)
				}
			}
			if len(t.Methods) > 0 {
				content.WriteString("#### Methods\n\n")
				for _, m := range t.Methods {
					writeDocItem(&content, fset, m.Doc, []string{m.Name}, m.Decl, 4)
				}
			}
		}
	}

	if _, err := content.WriteTo(out); err != nil {
		return fmt.Errorf("failed to write content buffer to output file: %w", err)
	}

	footer := "{{ endinclude }}"
	if _, err := out.WriteString(footer); err != nil {
		return fmt.Errorf("failed to write footer to output file: %w", err)
	}

	log.Printf("Successfully generated reference docs to %s", outputFile)
	return nil
}

// generateTOC builds a Table of Contents using GitHub-style heading anchors.
func generateTOC(docPkg *doc.Package) string {
	var tocBuf bytes.Buffer
	tocBuf.WriteString("## Table of Contents\n\n")

	hasContent := false

	if docPkg.Doc != "" {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Overview", generateAnchor("Overview")))
		hasContent = true
	}
	if len(docPkg.Consts) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Constants", generateAnchor("Constants")))
		hasContent = true
	}
	if len(docPkg.Vars) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Variables", generateAnchor("Variables")))
		hasContent = true
	}
	if len(docPkg.Funcs) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Functions", generateAnchor("Functions")))
		hasContent = true
	}
	if len(docPkg.Types) > 0 {
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", "Types", generateAnchor("Types")))
		hasContent = true
		for _, t := range docPkg.Types {
			typeTitle := fmt.Sprintf("`%s`", t.Name)
			typeAnchor := generateAnchor(t.Name)
			tocBuf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", typeTitle, typeAnchor))
		}
	}

	if !hasContent {
		return ""
	}
	return tocBuf.String()
}

// writeDocItem writes a documentation item heading, its declaration, and doc text.
func writeDocItem(contentBuf *bytes.Buffer, fset *token.FileSet, docComment string, names []string, decl ast.Node, level int) {
	displayName := strings.Join(names, ", ")
	fmt.Fprintf(contentBuf, "%s `%s`\n\n", strings.Repeat("#", level), displayName)

	printDeclaration(contentBuf, fset, decl, displayName)

	if docComment != "" {
		contentBuf.WriteString(formatDocText(docComment))
		contentBuf.WriteString("\n\n")
	}
}

// printDeclaration writes an AST declaration into a Go code fence.
func printDeclaration(buf *bytes.Buffer, fset *token.FileSet, decl ast.Node, name string) {
	var declBuf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
	if err := cfg.Fprint(&declBuf, fset, decl); err != nil {
		log.Printf("Warning: Failed to print declaration for %s: %v", name, err)
		buf.WriteString("```go\n// Error printing declaration\n```\n\n")
		return
	}
	buf.WriteString("```go\n")
	buf.Write(declBuf.Bytes())
	buf.WriteString("\n```\n\n")
}

// formatDocText converts Go doc comments to Markdown, unescaping HTML,
// turning "Parameters:" into a small header and cleaning code fences.
func formatDocText(text string) string {
	trimmed := strings.TrimSpace(text)
	var buf bytes.Buffer
	doc.ToHTML(&buf, trimmed, nil)
	htmlStr := buf.String()

	md := strings.ReplaceAll(htmlStr, "<p>", "")
	md = strings.ReplaceAll(md, "</p>", "\n\n")
	md = strings.ReplaceAll(md, "<pre>", "```go\n")
	md = strings.ReplaceAll(md, "</pre>", "\n```\n\n")
	md = strings.ReplaceAll(md, "<code>", "`")
	md = strings.ReplaceAll(md, "</code>", "`")
	md = strings.ReplaceAll(md, "<ul>", "")
	md = strings.ReplaceAll(md, "</ul>", "")
	md = strings.ReplaceAll(md, "<li>", "- ")
	md = strings.ReplaceAll(md, "</li>", "\n")

	md = strings.TrimSpace(md)

	// Unescape HTML entities.
	md = html.UnescapeString(md)

	// Convert a "Parameters" paragraph into a small header.
	md = regexp.MustCompile(`(?m)^\s*Parameters:\s*$`).ReplaceAllString(md, "#### Parameters")

	// Ensure code fences don't end up with an extra blank line before the closing fence.
	reFence := regexp.MustCompile("(?s)```go\\n(.*?)\\n+```")
	md = reFence.ReplaceAllString(md, "```go\n$1\n```")

	// Collapse excessive blank lines.
	md = regexp.MustCompile(`\n{3,}`).ReplaceAllString(md, "\n\n")

	return strings.TrimSpace(md)
}

// generateAnchor returns a GitHub-style slug for text.
func generateAnchor(text string) string {
	s := strings.TrimSpace(text)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "`", "")
	s = regexp.MustCompile(`[^a-z0-9 -]+`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
