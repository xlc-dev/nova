package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// templateConfig holds configuration for creating a project from templates.
// TemplateFS is the embedded filesystem containing the templates.
// TemplateDir is the root directory within TemplateFS to walk.
// Verbose controls whether creation steps are logged to stdout.
type templateConfig struct {
	TemplateFS  embed.FS
	TemplateDir string
	Verbose     bool
}

// templateData holds the data passed to path and content templates.
// ProjectName is the project name, DBImport is the database identifier,
// and DBAdapter is the Go import path for the chosen database adapter.
type templateData struct {
	ProjectName string
	DBImport    string
	DBAdapter   string
}

// getDBAdapter returns the Go import path for a given database identifier.
// Supported identifiers are "sqlite", "postgres", and "mysql".
// Returns the empty string for unknown identifiers.
func getDBAdapter(dbImport string) string {
	switch dbImport {
	case "sqlite":
		return "modernc.org/sqlite"
	case "postgres":
		return "github.com/lib/pq"
	case "mysql":
		return "github.com/go-sql-driver/mysql"
	default:
		return ""
	}
}

// createFromTemplate walks the embedded TemplateFS under TemplateDir,
// processes each path and file as a Go text/template applied to
// templateData, and writes the results under the new project directory
// named by 'ProjectName'. Returns any error encountered.
func createFromTemplate(
	projectDir string,
	projectName string,
	config templateConfig,
	dbImport string,
) error {
	// Ensure the root project directory exists before we start
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Prepare the data for templating
	data := templateData{
		ProjectName: projectName,
		DBImport:    dbImport,
		DBAdapter:   getDBAdapter(dbImport),
	}

	// Walk the embedded filesystem and handle each entry
	return fs.WalkDir(config.TemplateFS, config.TemplateDir, func(
		path string, d fs.DirEntry, err error,
	) error {
		if err != nil {
			return fmt.Errorf("filesystem walk error: %w", err)
		}

		// Compute the path relative to the template root
		relPath, err := filepath.Rel(config.TemplateDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Process the path as a template, stripping .tmpl extension
		targetPath, err := processPathTemplate(relPath, projectDir, data)
		if err != nil {
			return fmt.Errorf("failed to process path template: %w", err)
		}

		// If it's a directory, create it; otherwise process a file
		if d.IsDir() {
			if path == config.TemplateDir {
				return nil // Skip the root template dir itself
			}
			return handleDirectory(targetPath, config.Verbose)
		}
		return handleFile(
			config.TemplateFS,
			path,
			targetPath,
			data,
			config.Verbose,
		)
	})
}

// processPathTemplate applies Go text/template to a file or directory path,
// stripping any ".tmpl" suffix from the base name. The result is joined under
// the project root 'ProjectName'. Returns the final filesystem path.
func processPathTemplate(
	path, projectDir string, data templateData,
) (string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// Remove .tmpl extension if present
	if strings.HasSuffix(base, ".tmpl") {
		base = strings.TrimSuffix(base, ".tmpl")
	}

	// Reassemble the path without the suffix
	newPath := filepath.Join(dir, base)
	tmpl, err := template.New("path").Parse(newPath)
	if err != nil {
		return "", err
	}

	var processed bytes.Buffer
	if err := tmpl.Execute(&processed, data); err != nil {
		return "", err
	}

	// Prepend the project root directory
	return filepath.Join(projectDir, processed.String()), nil
}

// handleDirectory creates the directory (and parents) at path.
// If verbose is true, prints the action to stdout.
func handleDirectory(path string, verbose bool) error {
	if verbose {
		fmt.Printf("Creating directory: %s\n", path)
	}
	return os.MkdirAll(path, 0755)
}

// handleFile reads a template file from templateFS at srcPath, applies
// Go text/template with templateData to its content, and writes the
// result to destPath. Logs the action if verbose is true.
func handleFile(
	templateFS embed.FS,
	srcPath, destPath string,
	data templateData,
	verbose bool,
) error {
	if verbose {
		fmt.Printf("Creating file: %s\n", destPath)
	}

	// Read the raw template content
	content, err := fs.ReadFile(templateFS, srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Process the content as a template
	processed, err := processContentTemplate(srcPath, string(content), data)
	if err != nil {
		return fmt.Errorf("failed to process content template: %w", err)
	}

	// Write the final content to disk
	return os.WriteFile(destPath, []byte(processed), 0644)
}

// processContentTemplate parses and executes a text/template with name,
// content, and data, returning the processed string or an error.
func processContentTemplate(
	name, content string, data templateData,
) (string, error) {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return "", err
	}
	var processed bytes.Buffer
	if err := tmpl.Execute(&processed, data); err != nil {
		return "", err
	}
	return processed.String(), nil
}

//go:embed minimal/*
var minimalTemplate embed.FS

//go:embed structured/*
var structuredTemplate embed.FS

// CreateMinimal generates a new project using the minimal template layout.
// name is the project directory, verbose enables logging, and dbImport
// selects the database adapter ("sqlite", "postgres", "mysql").
func CreateMinimal(projectDir string, verbose bool, dbImport string) error {
	projectName := filepath.Base(projectDir)
	return createFromTemplate(projectDir, projectName, templateConfig{
		TemplateFS:  minimalTemplate,
		TemplateDir: "minimal",
		Verbose:     verbose,
	}, dbImport)
}

// CreateStructured generates a new project using the structured template layout.
// name is the project directory, verbose enables logging, and dbImport
// selects the database adapter ("sqlite", "postgres", "mysql").
func CreateStructured(projectDir string, verbose bool, dbImport string) error {
	projectName := filepath.Base(projectDir)
	return createFromTemplate(projectDir, projectName, templateConfig{
		TemplateFS:  structuredTemplate,
		TemplateDir: "structured",
		Verbose:     verbose,
	}, dbImport)
}
