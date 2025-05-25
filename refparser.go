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

	// Convert file map to slice for doc.NewFromFiles
	astFiles := make([]*ast.File, 0, len(pkg.Files))
	for _, file := range pkg.Files {
		astFiles = append(astFiles, file)
	}

	// Extract documentation (exported symbols only by default)
	docPkg, err := doc.NewFromFiles(fset, astFiles, pkg.Name)
	if err != nil {
		return fmt.Errorf("failed to create doc package for '%s': %w", pkg.Name, err)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer out.Close()

	var tocBuf bytes.Buffer
	var contentBuf bytes.Buffer

	title := "# Reference\n\n"
	tocBuf.WriteString("## Table of Contents\n\n")

	// Package Documentation (if any)
	if docPkg.Doc != "" {
		sectionTitle := "Overview"
		anchor := generateAnchor(sectionTitle)
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", sectionTitle, anchor))
		contentBuf.WriteString(fmt.Sprintf("## %s\n\n", sectionTitle))
		contentBuf.WriteString(formatDocText(docPkg.Doc))
		contentBuf.WriteString("\n\n")
	}

	// Constants
	if len(docPkg.Consts) > 0 {
		sectionTitle := "Constants"
		anchor := generateAnchor(sectionTitle)
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", sectionTitle, anchor))
		contentBuf.WriteString(fmt.Sprintf("## %s\n\n", sectionTitle))
		for _, c := range docPkg.Consts {
			displayNames := strings.Join(c.Names, ", ")
			itemTitle := fmt.Sprintf("`%s`", displayNames)
			itemAnchor := generateAnchor(itemTitle)
			tocBuf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", itemTitle, itemAnchor))
			writeDocItem(&contentBuf, fset, c.Doc, c.Names, c.Decl)
		}
	}

	// Variables
	if len(docPkg.Vars) > 0 {
		sectionTitle := "Variables"
		anchor := generateAnchor(sectionTitle)
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", sectionTitle, anchor))
		contentBuf.WriteString(fmt.Sprintf("## %s\n\n", sectionTitle))
		for _, v := range docPkg.Vars {
			displayNames := strings.Join(v.Names, ", ")
			itemTitle := fmt.Sprintf("`%s`", displayNames)
			itemAnchor := generateAnchor(itemTitle)
			tocBuf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", itemTitle, itemAnchor))
			writeDocItem(&contentBuf, fset, v.Doc, v.Names, v.Decl)
		}
	}

	// Functions
	if len(docPkg.Funcs) > 0 {
		sectionTitle := "Functions"
		anchor := generateAnchor(sectionTitle)
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", sectionTitle, anchor))
		contentBuf.WriteString(fmt.Sprintf("## %s\n\n", sectionTitle))
		for _, f := range docPkg.Funcs {
			itemTitle := fmt.Sprintf("`%s`", f.Name)
			itemAnchor := generateAnchor(itemTitle)
			tocBuf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", itemTitle, itemAnchor))
			writeDocItem(&contentBuf, fset, f.Doc, []string{f.Name}, f.Decl)
		}
	}

	// Types
	if len(docPkg.Types) > 0 {
		sectionTitle := "Types"
		anchor := generateAnchor(sectionTitle)
		tocBuf.WriteString(fmt.Sprintf("- [%s](#%s)\n", sectionTitle, anchor))
		contentBuf.WriteString(fmt.Sprintf("## %s\n\n", sectionTitle))
		for _, t := range docPkg.Types {
			typeTitle := fmt.Sprintf("`%s`", t.Name)
			typeAnchor := generateAnchor(typeTitle)
			tocBuf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", typeTitle, typeAnchor))
			// Document the type itself
			writeDocItem(&contentBuf, fset, t.Doc, []string{t.Name}, t.Decl)

			if len(t.Consts) > 0 {
				subSectionTitle := "Associated Constants"
				contentBuf.WriteString(fmt.Sprintf("#### %s\n\n", subSectionTitle))
				for _, c := range t.Consts {
					displayNames := strings.Join(c.Names, ", ")
					itemTitle := fmt.Sprintf("`%s`", displayNames)
					itemAnchor := generateAnchor(itemTitle)
					tocBuf.WriteString(fmt.Sprintf("    - [%s](#%s)\n", itemTitle, itemAnchor))
					writeDocItem(&contentBuf, fset, c.Doc, c.Names, c.Decl)
				}
			}
			if len(t.Vars) > 0 {
				subSectionTitle := "Associated Variables"
				contentBuf.WriteString(fmt.Sprintf("#### %s\n\n", subSectionTitle))
				for _, v := range t.Vars {
					displayNames := strings.Join(v.Names, ", ")
					itemTitle := fmt.Sprintf("`%s`", displayNames)
					itemAnchor := generateAnchor(itemTitle)
					tocBuf.WriteString(fmt.Sprintf("    - [%s](#%s)\n", itemTitle, itemAnchor))
					writeDocItem(&contentBuf, fset, v.Doc, v.Names, v.Decl)
				}
			}
			if len(t.Funcs) > 0 {
				subSectionTitle := "Associated Functions"
				contentBuf.WriteString(fmt.Sprintf("#### %s\n\n", subSectionTitle))
				for _, f := range t.Funcs {
					itemTitle := fmt.Sprintf("`%s`", f.Name)
					itemAnchor := generateAnchor(itemTitle)
					tocBuf.WriteString(fmt.Sprintf("    - [%s](#%s)\n", itemTitle, itemAnchor))
					writeDocItem(&contentBuf, fset, f.Doc, []string{f.Name}, f.Decl)
				}
			}
			if len(t.Methods) > 0 {
				subSectionTitle := "Methods"
				contentBuf.WriteString(fmt.Sprintf("#### %s\n\n", subSectionTitle))
				for _, m := range t.Methods {
					itemTitle := fmt.Sprintf("`%s.%s`", t.Name, m.Name)
					itemAnchor := generateAnchor(itemTitle)
					tocBuf.WriteString(fmt.Sprintf("    - [%s](#%s)\n", itemTitle, itemAnchor))
					writeDocItem(&contentBuf, fset, m.Doc, []string{m.Name}, m.Decl)
				}
			}
		}
	}

	// Write Title
	_, err = out.WriteString(title)
	if err != nil {
		return fmt.Errorf("failed to write title to output file: %w", err)
	}

	// Write Table of Contents
	if tocBuf.Len() > len("## Table of Contents\n\n") { // Only write TOC if it has entries
		_, err = tocBuf.WriteTo(out)
		if err != nil {
			return fmt.Errorf("failed to write TOC buffer to output file: %w", err)
		}
		_, err = out.WriteString("\n\n") // Add separation after TOC
		if err != nil {
			return fmt.Errorf("failed to write separator after TOC: %w", err)
		}
	}

	// Write Main Content
	_, err = contentBuf.WriteTo(out)
	if err != nil {
		return fmt.Errorf("failed to write content buffer to output file: %w", err)
	}

	log.Printf("Successfully generated reference docs with TOC to %s", outputFile)
	return nil
}

// writeDocItem formats a single documentation item (const, var, func, type, method)
// and writes it to the content buffer.
func writeDocItem(contentBuf *bytes.Buffer, fset *token.FileSet, docComment string, names []string, decl ast.Node) {
	// Use the first name for the header, join for display if multiple (const/var blocks)
	// For methods, names only contains the method name, not the receiver type.
	// The anchor generation in the main function handles creating unique anchors like Type.Method.
	displayNames := strings.Join(names, ", ")
	itemTitle := fmt.Sprintf("`%s`", displayNames) // Title used for the H3 header

	// Generate an anchor based *only* on the display name(s).
	// For methods, this will be `#methodname`. For associated funcs, `#funcname`.
	// For types, `#typename`. For top-level funcs/consts/vars, `#name`.
	// The TOC links generated earlier might be more specific (e.g., `#type.methodname`)
	// but the header ID generated by Markdown processors will likely match this simpler anchor.
	itemAnchor := generateAnchor(itemTitle)

	contentBuf.WriteString(fmt.Sprintf("### %s {#%s}\n\n", itemTitle, itemAnchor)) // Header for the item with explicit anchor

	// Print the declaration (signature) using go/printer
	var declBuf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
	err := cfg.Fprint(&declBuf, fset, decl)
	if err != nil {
		log.Printf("Warning: Failed to print declaration for %s: %v", displayNames, err)
		contentBuf.WriteString("```go\n// Error printing declaration\n```\n\n")
	} else {
		contentBuf.WriteString("```go\n")
		contentBuf.Write(declBuf.Bytes()) // Write the formatted declaration
		contentBuf.WriteString("\n```\n\n")
	}

	// Write the documentation comment
	if docComment != "" {
		contentBuf.WriteString(formatDocText(docComment))
		contentBuf.WriteString("\n\n")
	}
	contentBuf.WriteString("---\n\n") // Separator
}

// formatDocText performs basic formatting on Go doc strings for Markdown.
// (No changes needed in this function)
func formatDocText(text string) string {
	trimmed := strings.TrimSpace(text)
	// Split into paragraphs based on double newline
	paragraphs := strings.Split(trimmed, "\n\n")
	for i, p := range paragraphs {
		if strings.HasPrefix(p, "    ") || strings.HasPrefix(p, "\t") {
			lines := strings.Split(p, "\n")
			cleanedLines := make([]string, len(lines))
			for j, line := range lines {
				cleanedLines[j] = strings.TrimPrefix(strings.TrimPrefix(line, "    "), "\t")
			}
			paragraphs[i] = "```text\n" + strings.Join(cleanedLines, "\n") + "\n```"
		} else {
			// Replace single newlines within a paragraph with spaces for better flow
			// Keep code examples formatted correctly by not replacing newlines inside ``` blocks
			if !strings.HasPrefix(paragraphs[i], "```") {
				paragraphs[i] = strings.ReplaceAll(p, "\n", " ")
			}
			// Basic fix for lists potentially getting joined: re-insert newline before list items
			paragraphs[i] = strings.ReplaceAll(paragraphs[i], " - ", "\n - ")
			paragraphs[i] = strings.ReplaceAll(paragraphs[i], " * ", "\n * ")
			paragraphs[i] = regexp.MustCompile(` (\d+)\. `).ReplaceAllString(paragraphs[i], "\n$1. ")

		}
	}
	// Join paragraphs with double newline for Markdown paragraph separation
	return strings.Join(paragraphs, "\n\n")
}

// generateAnchor creates a Markdown-friendly anchor link from a string.
// It converts to lowercase, replaces spaces with hyphens, and removes non-alphanumeric characters except hyphens.
func generateAnchor(text string) string {
	// Remove backticks first
	text = strings.ReplaceAll(text, "`", "")
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")

	// Remove any characters that are not letters, numbers, or hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	text = reg.ReplaceAllString(text, "")
	// Ensure it doesn't start or end with a hyphen
	text = strings.Trim(text, "-")
	return text
}
