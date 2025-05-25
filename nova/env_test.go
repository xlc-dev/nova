package nova

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadDotenv_NonExistentFile ensures that a non-existent .env file
// is silently ignored and no error is returned.
func TestLoadDotenv_NonExistentFile(t *testing.T) {
	if err := LoadDotenv("nonexistent.env"); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

// TestLoadDotenv_SimpleParsing verifies basic parsing of .env content,
// including trimming, quoted values, escapes, and literal hashes.
func TestLoadDotenv_SimpleParsing(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := `
# Sample .env
FOO=bar
SINGLE='single quoted'
DOUBLE="double quoted"
ESCAPED="a\nb\tc\\d"
LITERAL_HASH=foo#bar
TRAILING=foo\\
`
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("writing .env file: %v", err)
	}
	// Clean up any pre-existing variables
	for _, k := range []string{
		"FOO", "SINGLE", "DOUBLE", "ESCAPED", "LITERAL_HASH", "TRAILING",
	} {
		os.Unsetenv(k)
	}
	if err := LoadDotenv(file); err != nil {
		t.Fatalf("LoadDotenv returned error: %v", err)
	}
	expected := map[string]string{
		"FOO":          "bar",
		"SINGLE":       "single quoted",
		"DOUBLE":       "double quoted",
		"ESCAPED":      "a\nb\tc\\d",
		"LITERAL_HASH": "foo#bar",
		"TRAILING":     "foo\\",
	}
	for key, want := range expected {
		val, ok := os.LookupEnv(key)
		if !ok {
			t.Errorf("expected %q to be set", key)
			continue
		}
		if val != want {
			t.Errorf("for %q: got %q, want %q", key, val, want)
		}
	}
}

// TestLoadDotenv_ExpansionAndCycle verifies variable expansion from other
// entries and the host environment, and ensures cycles yield empty values.
func TestLoadDotenv_ExpansionAndCycle(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "dotenv")
	content := `
VAR1=hello
VAR2=$VAR1 world
VAR3=${VAR1}_suffix
REFOS=$OSVAR
A=$B
B=$A
`
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("writing dotenv file: %v", err)
	}
	// Prepare environment: OSVAR exists before loading
	t.Setenv("OSVAR", "osvalue")
	// Clean up any old values
	for _, k := range []string{"VAR1", "VAR2", "VAR3", "REFOS", "A", "B"} {
		os.Unsetenv(k)
	}
	if err := LoadDotenv(file); err != nil {
		t.Fatalf("LoadDotenv returned error: %v", err)
	}
	// Check proper expansions
	expTests := []struct{ key, want string }{
		{"VAR1", "hello"},
		{"VAR2", "hello world"},
		{"VAR3", "hello_suffix"},
		{"REFOS", "osvalue"},
	}
	for _, tc := range expTests {
		val, ok := os.LookupEnv(tc.key)
		if !ok || val != tc.want {
			t.Errorf("for %q: got (%q, %v), want (%q, true)",
				tc.key, val, ok, tc.want)
		}
	}
	// Cycle detection: A and B should be set to empty
	for _, k := range []string{"A", "B"} {
		val, ok := os.LookupEnv(k)
		if !ok {
			t.Errorf("expected %q to exist (even if empty)", k)
		}
		if val != "" {
			t.Errorf("expected %q to be empty, got %q", k, val)
		}
	}
}

// TestLoadDotenv_SkipExisting ensures that already set environment
// variables are not overwritten by LoadDotenv.
func TestLoadDotenv_SkipExisting(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := `
EXISTING=changed
NEW=added
`
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("writing .env file: %v", err)
	}
	t.Setenv("EXISTING", "original")
	os.Unsetenv("NEW")
	if err := LoadDotenv(file); err != nil {
		t.Fatalf("LoadDotenv returned error: %v", err)
	}
	if val, _ := os.LookupEnv("EXISTING"); val != "original" {
		t.Errorf("EXISTING was overwritten: got %q, want %q", val, "original")
	}
	if val, ok := os.LookupEnv("NEW"); !ok || val != "added" {
		t.Errorf("NEW: got (%q, %v), want (%q, true)", val, ok, "added")
	}
}

// TestLoadDotenv_UnclosedQuotesAndInvalid verifies that lines with unclosed
// quotes or invalid keys are skipped, while valid lines are loaded.
func TestLoadDotenv_UnclosedQuotesAndInvalid(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := `
UNCLOSED="missing end
BAD KEY=foo
GOOD=bar
`
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("writing .env file: %v", err)
	}
	// Clean up
	os.Unsetenv("UNCLOSED")
	os.Unsetenv("BAD KEY")
	os.Unsetenv("GOOD")
	if err := LoadDotenv(file); err != nil {
		t.Fatalf("LoadDotenv returned error: %v", err)
	}
	if _, ok := os.LookupEnv("UNCLOSED"); ok {
		t.Errorf("UNCLOSED should be skipped due to unclosed quotes")
	}
	if _, ok := os.LookupEnv("BAD KEY"); ok {
		t.Errorf("BAD KEY should be skipped due to invalid key")
	}
	if val, ok := os.LookupEnv("GOOD"); !ok || val != "bar" {
		t.Errorf("GOOD: got (%q, %v), want (%q, true)", val, ok, "bar")
	}
}
