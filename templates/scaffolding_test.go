package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetDBAdapter verifies that getDBAdapter returns the correct import
// path for known databases and empty string for unknown.
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

// TestProcessPathTemplate ensures that processPathTemplate strips the
// .tmpl extension and applies data to path templates.
func TestProcessPathTemplate(t *testing.T) {
	data := templateData{
		Name:      "proj",
		DBImport:  "sqlite",
		DBAdapter: getDBAdapter("sqlite"),
	}
	// path template uses Name in directory and filename
	in := "dir/{{.Name}}.tmpl"
	out, err := processPathTemplate(in, "proj", data)
	if err != nil {
		t.Fatalf("processPathTemplate error: %v", err)
	}
	want := filepath.Join("proj", "dir", "proj")
	if out != want {
		t.Errorf("got path %q; want %q", out, want)
	}
}

// TestProcessContentTemplate verifies that templates in file content
// are executed correctly.
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

// TestHandleDirectory checks that handleDirectory creates nested dirs
// and prints verbose messages when requested.
func TestHandleDirectory(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "a", "b", "c")
	// verbose = false: should create without error
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

// TestCreateMinimalIntegration runs CreateMinimal on the embedded
// minimalTemplate and checks that expected files are generated.
func TestCreateMinimalIntegration(t *testing.T) {
	name := "testminimal"
	defer os.RemoveAll(name)

	if err := CreateMinimal(name, false, "sqlite"); err != nil {
		t.Fatalf("CreateMinimal failed: %v", err)
	}
	// root dir exists
	info, err := os.Stat(name)
	if err != nil || !info.IsDir() {
		t.Fatalf("project dir %q not created", name)
	}

	// check that main.go exists
	mainPath := filepath.Join(name, "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("generated file %q missing: %v", mainPath, err)
	}
	// read and verify it contains the sqlite adapter import
	out, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("reading %q: %v", mainPath, err)
	}
	adapter := getDBAdapter("sqlite")
	if !strings.Contains(string(out), adapter) {
		t.Errorf("generated main.go missing import %q", adapter)
	}
}

// TestCreateStructuredAndTODO ensures CreateStructured and CreateTODO
// create their project roots without error.
func TestCreateStructuredAndTODO(t *testing.T) {
	cases := []struct {
		name string
		call func(string, bool, string) error
	}{
		{"structured", CreateStructured},
		{"todo", CreateTODO},
	}
	for _, tc := range cases {
		dir := "proj_" + tc.name
		os.RemoveAll(dir)
		if err := tc.call(dir, false, "postgres"); err != nil {
			t.Errorf("Create%s failed: %v", strings.Title(tc.name), err)
			continue
		}
		// root dir exists
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Errorf("%s project dir not created", tc.name)
		}
		os.RemoveAll(dir)
	}
}
