package nova

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// TestNewCLI_Validation verifies that NewCLI correctly validates the CLI
// configuration, returning errors for nil input, missing fields, and
// reserved flag names.
func TestNewCLI_Validation(t *testing.T) {
	cases := []struct {
		name     string
		cli      *CLI
		wantErr  bool
		contains string
	}{
		{
			name:     "nil CLI",
			cli:      nil,
			wantErr:  true,
			contains: "configuration cannot be nil",
		},
		{
			name: "empty Name",
			cli: &CLI{
				Name:    " ",
				Version: "1.0",
			},
			wantErr:  true,
			contains: "Name cannot be empty",
		},
		{
			name: "empty Version",
			cli: &CLI{
				Name:    "app",
				Version: " ",
			},
			wantErr:  true,
			contains: "Version cannot be empty",
		},
		{
			name: "reserved global flag",
			cli: &CLI{
				Name:    "app",
				Version: "1",
				GlobalFlags: []Flag{
					&BoolFlag{Name: "version", Usage: "bad"},
				},
			},
			wantErr:  true,
			contains: "flag name 'version' is reserved",
		},
		{
			name: "good CLI",
			cli: &CLI{
				Name:    "app",
				Version: "1.0",
			},
			wantErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewCLI(c.cli)
			if c.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), c.contains) {
					t.Fatalf("error %q does not contain %q", err.Error(), c.contains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// runCLI is a helper that runs the CLI with the provided arguments,
// captures stdout, and returns the output and any error.
func runCLI(cli *CLI, args []string) (string, error) {
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.Run(append([]string{"cmd"}, args...))

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldOut

	return buf.String(), err
}

// TestRun_GlobalVersionAndHelp ensures that the global --version flag
// prints the version and that running with no args prints the help.
func TestRun_GlobalVersionAndHelp(t *testing.T) {
	cli, _ := NewCLI(&CLI{
		Name:    "myapp",
		Version: "9.9",
	})

	out, err := runCLI(cli, []string{"--version"})
	if err != nil {
		t.Fatalf("version returned error: %v", err)
	}
	want := "myapp version 9.9"
	if !strings.Contains(out, want) {
		t.Errorf("got %q, want to contain %q", out, want)
	}

	out, err = runCLI(cli, []string{})
	if err != nil {
		t.Fatalf("help returned error: %v", err)
	}
	if !strings.Contains(out, "Usage: myapp") {
		t.Errorf("help output missing Usage: got %q", out)
	}
}

// TestRun_UnknownCommand checks that invoking an undefined command
// returns an appropriate error.
func TestRun_UnknownCommand(t *testing.T) {
	cli, _ := NewCLI(&CLI{
		Name:    "app",
		Version: "1",
	})
	_, err := runCLI(cli, []string{"foo"})
	if err == nil || !strings.Contains(err.Error(), `unknown command "foo"`) {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

// TestRun_DefaultAction verifies that the default Action is invoked
// with the correct arguments when no subcommand is specified.
func TestRun_DefaultAction(t *testing.T) {
	called := false
	cli, _ := NewCLI(&CLI{
		Name:    "app",
		Version: "1",
		Action: func(ctx *Context) error {
			called = true
			if len(ctx.Args()) != 2 || ctx.Args()[0] != "a" || ctx.Args()[1] != "b" {
				return fmt.Errorf("unexpected args: %v", ctx.Args())
			}
			return nil
		},
	})
	out, err := runCLI(cli, []string{"a", "b"})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !called {
		t.Error("default Action was not called")
	}
	if out != "" {
		t.Errorf("expected no stdout, got %q", out)
	}
}

// TestStringSliceFlag ensures that StringSliceFlag collects multiple
// occurrences of a flag into a slice, preserving order.
func TestStringSliceFlag(t *testing.T) {
	cmd := &Command{
		Name:  "tag",
		Usage: "collect tags",
		Action: func(ctx *Context) error {
			tags := ctx.StringSlice("tag")
			fmt.Fprint(os.Stdout, strings.Join(tags, ","))
			return nil
		},
		Flags: []Flag{
			&StringSliceFlag{Name: "tag", Aliases: []string{"t"}, Usage: "add tag"},
		},
	}
	cli, _ := NewCLI(&CLI{
		Name:     "app",
		Version:  "1",
		Commands: []*Command{cmd},
	})

	out, err := runCLI(cli, []string{"tag", "--tag", "one", "-t", "two"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// order is preserved
	if out != "one,two" {
		t.Errorf("got %q, want %q", out, "one,two")
	}
}
