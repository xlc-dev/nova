package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetDBAdapter remains the same, no file I/O.
func TestGetDBAdapter(t *testing.T) {
	cases := []struct {
		dbImport string
		want     string
	}{
		{"sqlite", "modernc.org/sqlite"},
		{"postgres", "github.com/lib/pq"},
		{"mysql", "github.com/go-sql-driver/mysql"},
		{"unknown", ""},
	}
	for _, c := range cases {
		got := getDBAdapter(c.dbImport)
		if got != c.want {
			t.Errorf("getDBAdapter(%q) = %q; want %q",
				c.dbImport, got, c.want)
		}
	}
}

// TestProcessPathTemplate remains the same, no file I/O.
func TestProcessPathTemplate(t *testing.T) {
	data := templateData{
		Name:      "proj",
		DBImport:  "sqlite",
		DBAdapter: getDBAdapter("sqlite"),
	}
	in := "dir/{{.Name}}.tmpl"
	out, err := processPathTemplate(in, "my-project", data)
	if err != nil {
		t.Fatalf("processPathTemplate error: %v", err)
	}
	want := filepath.Join("my-project", "dir", "proj")
	if out != want {
		t.Errorf("got path %q; want %q", out, want)
	}
}

// TestProcessContentTemplate remains the same, no file I/O.
func TestProcessContentTemplate(t *testing.T) {
	name := "file.tmpl"
	content := `package {{.Name}}

import "{{.DBAdapter}}"
`
	data := templateData{
		Name:      "mypkg",
		DBImport:  "sqlite",
		DBAdapter: getDBAdapter("sqlite"),
	}
	got, err := processContentTemplate(name, content, data)
	if err != nil {
		t.Fatalf("processContentTemplate error: %v", err)
	}
	if !strings.Contains(got, "package mypkg") {
		t.Error("output missing package declaration")
	}
	if !strings.Contains(got, `import "modernc.org/sqlite"`) {
		t.Error("output missing correct import path")
	}
}

// TestHandleDirectory now uses t.TempDir() for safety.
func TestHandleDirectory(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "a", "b", "c")

	if err := handleDirectory(target, false); err != nil {
		t.Fatalf("handleDirectory failed: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("%q is not a directory", target)
	}
}

// TestCreateMinimalIntegration now uses t.TempDir() to prevent deleting source.
func TestCreateMinimalIntegration(t *testing.T) {
	projectDir := t.TempDir()

	if err := CreateMinimal(projectDir, false, "sqlite"); err != nil {
		t.Fatalf("CreateMinimal failed: %v", err)
	}
	info, err := os.Stat(projectDir)
	if err != nil || !info.IsDir() {
		t.Fatalf("project dir %q not created", projectDir)
	}

	mainPath := filepath.Join(projectDir, "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("generated file %q missing: %v", mainPath, err)
	}
	out, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("reading %q: %v", mainPath, err)
	}
	adapter := getDBAdapter("sqlite")
	if !strings.Contains(string(out), adapter) {
		t.Errorf("generated main.go missing import %q", adapter)
	}
}

// TestCreateStructuredAndTODO is updated to use t.TempDir() for best practice.
func TestCreateStructuredAndTODO(t *testing.T) {
	cases := []struct {
		name string
		call func(string, bool, string) error
	}{
		{"structured", CreateStructured},
		{"todo", CreateTODO},
	}
	for _, tc := range cases {
		projectDir := t.TempDir()

		if err := tc.call(projectDir, false, "postgres"); err != nil {
			t.Errorf("Create%s failed: %v", strings.Title(tc.name), err)
			continue
		}
		if info, err := os.Stat(projectDir); err != nil || !info.IsDir() {
			t.Errorf("%s project dir not created", tc.name)
		}
	}
}
