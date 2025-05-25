package nova

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestHasAllowedExtension verifies that hasAllowedExtension correctly
// identifies filenames with allowed extensions.
func TestHasAllowedExtension(t *testing.T) {
	cases := []struct {
		filename string
		exts     []string
		want     bool
	}{
		{"foo.go", []string{".go"}, true},
		{"foo.txt", []string{".go"}, false},
		{"foo.txt", []string{" .txt", ".md"}, true},
		{"foo", []string{".go"}, false},
		{"foo.go", []string{}, false},
		{"", []string{".go"}, false},
	}
	for _, c := range cases {
		got := hasAllowedExtension(c.filename, c.exts)
		if got != c.want {
			t.Errorf("hasAllowedExtension(%q, %v) = %v; want %v",
				c.filename, c.exts, got, c.want)
		}
	}
}

// TestRecompileSuccess creates a temporary Go module with a minimal main.go
// and verifies that recompile succeeds and produces a non-empty binary.
func TestRecompileSuccess(t *testing.T) {
	tempDir := t.TempDir()
	prevWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(prevWd); err != nil {
			t.Fatalf("failed to restore working dir: %v", err)
		}
	}()

	// Write go.mod
	modFile := filepath.Join(tempDir, "go.mod")
	modContent := "module example.com/testmod\n"
	if err := os.WriteFile(modFile, []byte(modContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	// Write main.go
	mainFile := filepath.Join(tempDir, "main.go")
	mainContent := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	// Run recompile
	if err := recompile(false); err != nil {
		t.Fatalf("recompile failed: %v", err)
	}

	// Check for binary in CWD with possible names
	possibleBinaryNames := []string{"main", filepath.Base(tempDir), "testmod"}
	if runtime.GOOS == "windows" {
		for i := range possibleBinaryNames {
			possibleBinaryNames[i] += ".exe"
		}
	}

	cwdAfter, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get CWD after recompile: %v", err)
	}

	for _, binaryName := range possibleBinaryNames {
		binPath := filepath.Join(cwdAfter, binaryName)
		info, err := os.Stat(binPath)
		if err == nil && info.Size() > 0 {
			return // Success, exit early
		}
	}

	t.Fatalf("No valid binary found in %s", cwdAfter)
}

// TestRecompileError verifies that recompile returns an error when there
// are no Go files to build.
func TestRecompileError(t *testing.T) {
	tempDir := t.TempDir()
	prevWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(prevWd); err != nil {
			t.Fatalf("failed to restore working dir: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	if err := recompile(false); err == nil {
		t.Errorf("expected error for recompile in empty dir, got nil")
	}
}

// TestWatchAndCompileInvalidDir ensures watchAndCompile returns an error
// when given a non-existent directory.
func TestWatchAndCompileInvalidDir(t *testing.T) {
	err := watchAndCompile("nonexistent_dir", false, []string{".go"})
	if err == nil {
		t.Errorf("expected error for invalid directory, got nil")
	}
}

// TestWatchAndReloadInvalidDir ensures watchAndReload returns an error
// when given a non-existent directory.
func TestWatchAndReloadInvalidDir(t *testing.T) {
	ch := make(chan struct{}, 1)
	err := watchAndReload("nonexistent_dir", "exe", false, []string{".go"}, ch)
	if err == nil {
		t.Errorf("expected error for invalid directory, got nil")
	}
}
