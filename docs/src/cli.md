# CLI

Nova has a built in CLI feature. It simplifies common CLI tasks by providing:

-   **Structured Application:** Define your application commands clearly.
-   **Type-Safe Flag Parsing:** Supports `string`, `int`, `bool`, `float64`, and `[]string` flags.
-   **Automatic Help Generation:** Generates help text for the application and individual commands, including usage, descriptions, flags, and aliases. Features a built-in `help` command.
-   **Built-in Version Flag:** Automatically handles `--version` and `-v` global flags.
-   **Global and Command-Specific Flags:** Define flags that apply to the entire application or only to specific commands.
-   **Required Flag Validation:** Mark flags as mandatory, and Nova will enforce their presence.
-   **Default Values:** Specify default values for flags.
-   **Aliases:** Define short or alternative names for commands and flags.
-   **Context-Aware Actions:** Command actions receive a `Context` object providing access to parsed flags, arguments, and application metadata.
-   **Initialization Validation:** `NewCLI` validates your configuration for conflicts (e.g., reserved names) before running.

## Table of Contents

1.  [Getting Started](#getting-started)
2.  [Core Concepts](#core-concepts)
    -   [The `CLI` Struct](#the-cli-struct)
    -   [The `Command` Struct](#the-command-struct)
    -   [The `Flag` Interface](#the-flag-interface)
    -   [The `Context` Object](#the-context-object)
    -   [Execution Flow (`Run`)](#execution-flow-run)
3.  [Defining Commands](#defining-commands)
4.  [Working with Flags](#working-with-flags)
    -   [Flag Types](#flag-types)
    -   [Flag Definition](#flag-definition)
    -   [Global vs. Command Flags](#global-vs-command-flags)
    -   [Required Flags](#required-flags)
    -   [Default Values](#default-values)
    -   [Aliases](#aliases)
    -   [String Slice Flags](#string-slice-flags)
5.  [Context API](#context-api)
    -   [Accessing Arguments](#accessing-arguments)
    -   [Accessing Flag Values](#accessing-flag-values)
    -   [Accessing Metadata](#accessing-metadata)
6.  [Help and Version](#help-and-version)
7.  [Error Handling](#error-handling)
8.  [Full Example Application](#full-example-application)

## Getting Started

Here's a minimal example to get you started:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xlc-dev/nova/nova"
)

func main() {
	// Define the application structure
	app := &nova.CLI{
		Name:        "myapp",
		Version:     "1.0.0",
		Description: "A simple example CLI application.",
		Commands: []*nova.Command{
			{
				Name:   "greet",
				Usage:  "Prints a friendly greeting",
				Action: greetAction, // Function to execute
			},
		},
		// Optional: Define a default action if no command is given
		// Action: defaultAction,
	}

	// Validate the CLI configuration and initialize internal flags.
	// It's crucial to use NewCLI before Run.
	cli, err := nova.NewCLI(app)
	if err != nil {
		log.Fatalf("Failed to initialize CLI: %v", err)
	}

	// Run the application, passing command-line arguments
	if err := cli.Run(os.Args); err != nil {
		// Errors from parsing, validation, or the action itself are returned here.
		log.Fatalf("Error running CLI: %v", err)
	}
}

// greetAction is the function executed by the 'greet' command.
// It receives a Context object.
func greetAction(ctx *nova.Context) error {
	fmt.Println("Hello from Nova!")
	// You can access flags and arguments via ctx here if needed
	return nil // Return nil on success, or an error
}

/*
// Optional default action
func defaultAction(ctx *nova.Context) error {
    fmt.Println("Running default action. Use 'myapp help' for commands.")
    return nil
}
*/
```

Build and run:

```bash
go build -o myapp
./myapp greet
# Output: Hello from Nova!

./myapp --version
# Output: myapp version 1.0.0

./myapp help
# Output: Shows application help

./myapp help greet
# Output: Shows help specific to the greet command
```

## Core Concepts

### The `CLI` Struct

The `nova.CLI` struct is the root of your application. It holds metadata and configuration:

-   `Name` (string, **Required**): The name of your application, used in help messages.
-   `Version` (string, **Required**): The application version, displayed by the `--version` flag.
-   `Description` (string): A short description of the application, shown in the main help.
-   `Commands` ([]*`Command`): A slice of commands the application supports.
-   `Action` (`ActionFunc`): A function to run if no command is specified on the command line. If `nil` and no command is given, the main help is shown.
-   `GlobalFlags` ([]`Flag`): Flags that apply to the application globally, regardless of the command being run. Parsed before command flags.
-   `Authors` (string): Optional author information, displayed in the main help.

**Important:** Always initialize your application using `nova.NewCLI(app)`. This function validates your configuration (checking for required fields, reserved names like `help` for commands or `version` for global flags) and sets up internal flags before you call `Run`.

### The `Command` Struct

The `nova.Command` struct defines a specific action your application can perform.

-   `Name` (string, **Required**): The primary name used to invoke the command. Cannot be `help`.
-   `Aliases` ([]string): Alternative names for the command. Cannot include `h`.
-   `Usage` (string): A short, one-line description shown in the main command list.
-   `Description` (string): A more detailed description shown in the command's specific help (`myapp help <command>`).
-   `ArgsUsage` (string): Describes the expected arguments (e.g., `<input> [output]`), shown in the command's help.
-   `Action` (`ActionFunc`, **Required**): The function to execute when the command is invoked. It receives a `*nova.Context`.
-   `Flags` ([]`Flag`): Flags specific to this command. Cannot be named `help` or have an alias `h`.

### The `Flag` Interface

`nova.Flag` is an interface implemented by concrete flag types. You don't use the interface directly but rather the specific types:

-   `nova.StringFlag`
-   `nova.IntFlag`
-   `nova.BoolFlag`
-   `nova.Float64Flag`
-   `nova.StringSliceFlag`

Each flag type defines how a command-line option is parsed and stored. Key properties include `Name`, `Aliases`, `Usage`, `Default`, and `Required`.

### The `Context` Object

An instance of `nova.Context` is passed to every `ActionFunc`. It provides access to runtime information:

-   Parsed flag values (both command-specific and global).
-   Positional arguments remaining after flag parsing.
-   Metadata about the `CLI` and the currently executing `Command`.

### Execution Flow (`Run`)

When you call `cli.Run(os.Args)`:

1.  **Global Flags Parsing:** Nova parses global flags (including the built-in `--version`/`-v`). If `--version` is found, it prints the version and exits.
2.  **Global Flag Validation:** Required global flags are checked.
3.  **Command Identification:** Nova looks at the next argument to see if it matches a command name or alias (including the built-in `help` command).
4.  **Action Determination:**
    *   **Command Found:** If it's the `help` command, its action is run. Otherwise, the command's specific flags are parsed and validated. The command's `ActionFunc` is then executed with a `Context` containing command flags, global flags, and remaining arguments.
    *   **No Command Found & `CLI.Action` Defined:** The global `ActionFunc` is executed with a `Context` containing global flags and all non-flag arguments.
    *   **No Command Found & No `CLI.Action`:** If arguments were provided but didn't match a command, an "unknown command" error is returned. If no arguments were provided, the main application help is displayed.
5.  **Error Handling:** Any errors during parsing, validation, or action execution are returned by `Run`.

## Defining Commands

Define commands by creating instances of `nova.Command` and assigning them to the `CLI.Commands` slice.

```go
// Action function for the 'create' command
func createAction(ctx *nova.Context) error {
	// Access flags and args via ctx
	resourceName := ctx.Args()[0] // Example: Get first argument
	fmt.Printf("Creating resource: %s\n", resourceName)
	return nil
}

// Command definition
var createCmd = &nova.Command{
	Name:        "create",          // Primary name (e.g., `./myapp create`)
	Aliases:     []string{"c"},     // Aliases (e.g., `./myapp c`)
	Usage:       "Create a new resource", // Short help text
	Description: `Creates a new resource based on the provided name.
This command demonstrates basic command structure.`, // Long help text
	ArgsUsage:   "<resource-name>", // Argument syntax help
	Action:      createAction,      // Function to run
	Flags:       []nova.Flag{ /* Command-specific flags go here */ },
}

// In main():
app := &nova.CLI{
	// ... other CLI fields
	Commands: []*nova.Command{
		createCmd,
		// ... other commands
	},
}
```

**Reserved Names:** Remember, command `Name` cannot be `help`, and `Aliases` cannot include `h`.

## Working with Flags

Flags provide options and configuration for your application and commands.

### Flag Types

Nova provides the following concrete flag types, all implementing the `nova.Flag` interface:

| Type              | Go Type      | Example Usage           | Description                                    |
| ----------------- | ------------ | ----------------------- | ---------------------------------------------- |
| `StringFlag`      | `string`     | `--name="John Doe"`     | Accepts a text value.                          |
| `IntFlag`         | `int`        | `--port=8080`           | Accepts an integer value.                      |
| `BoolFlag`        | `bool`       | `--verbose`             | Acts as a switch (true if present).            |
| `Float64Flag`     | `float64`    | `--ratio=1.5`           | Accepts a floating-point number.               |
| `StringSliceFlag` | `[]string`   | `--tag foo --tag bar`   | Accepts multiple string values (repeatable).   |

### Flag Definition

Define flags by creating instances of the flag types and assigning them to `CLI.GlobalFlags` or `Command.Flags`.

```go
// StringFlag Example
&nova.StringFlag{
	Name:        "output",                       // Long name: --output
	Aliases:     []string{"o"},                  // Short name: -o
	Usage:       "Specify the output file path", // Help text
	Default:     "stdout",                       // Default value if flag not provided
	Required:    false,                          // The flag is optional
},

// BoolFlag Example
&nova.BoolFlag{
	Name:        "verbose",
	Aliases:     []string{"V"},              // Note: -v is reserved globally for version
	Usage:       "Enable verbose logging",
	Default:     false,
},

// IntFlag Example
&nova.IntFlag{
	Name:        "retries",
	Usage:       "Number of times to retry",
	Default:     3,
	Required:    true,                       // This flag must be provided
},

// Float64Flag Example
&nova.Float64Flag{
	Name:        "ratio",
	Aliases:     []string{"r"},
	Usage:       "A floating-point value representing a ratio",
	Default:     1.0,
	Required:    false,
},

// StringSliceFlag Example
&nova.StringSliceFlag{
	Name:        "tag",
	Aliases:     []string{"t"},
	Usage:       "Add one or more tags (can be specified multiple times)",
	Default:     []string{},               // Defaults to an empty slice if flag not provided
	Required:    false,
},
```

### Global vs. Command Flags

-   **Global Flags:** Defined in `CLI.GlobalFlags`. They are available and parsed *before* any command is run. Useful for options like `--config`, `--verbose`, or `--region`.
    -   **Reserved:** Cannot use `Name: "version"` or `Aliases: []string{"v"}`.
-   **Command Flags:** Defined in `Command.Flags`. They are only available and parsed when that specific command is invoked.
    -   **Reserved:** Cannot use `Name: "help"` or `Aliases: []string{"h"}`.

```go
app := &nova.CLI{
	// ...
	GlobalFlags: []nova.Flag{
		&nova.BoolFlag{
			Name:        "verbose",
			Usage:       "Enable verbose output globally",
		},
	},
	Commands: []*nova.Command{
		{
			Name:   "serve",
			Usage:  "Start a server",
			Action: serveAction,
			Flags: []nova.Flag{
				&nova.IntFlag{
					Name:        "port",
					Aliases:     []string{"p"},
					Usage:       "Port to listen on",
					Default:     8080,
				},
			},
		},
	},
}
```

### Required Flags

Set `Required: true` in the flag definition. Nova automatically checks if required flags were provided during the `Run` process *after* parsing. If a required flag is missing, `Run` returns an error. This applies to all flag types except `BoolFlag`. For `StringSliceFlag`, it ensures the flag was provided at least once.

```go
&nova.StringFlag{
	Name:        "api-key",
	Usage:       "Your API key for authentication",
	Required:    true, // ./myapp command --api-key=XYZ (Required)
},
```

### Default Values

Set the `Default` field in the flag definition. If the user provides the flag, the user's value overrides the default.

```go
&nova.IntFlag{
	Name:        "timeout",
	Usage:       "Request timeout in seconds",
	Default:     30, // Defaults to 30 if --timeout is not specified
},
```

### Aliases

The `Aliases` field (`[]string`) provides alternative names for flags.

-   Single-character aliases are typically prefixed with a single hyphen (`-a`).
-   Multi-character aliases are prefixed with double hyphens (`--alias`).

Nova handles the prefixing automatically based on the alias length when generating help text.

```go
&nova.StringFlag{
	Name:        "file",          // --file
	Aliases:     []string{"f"},   // -f
	Usage:       "Input filename",
},
```

### String Slice Flags

`StringSliceFlag` allows a flag to be specified multiple times, collecting all values into a `[]string`.

```go
&nova.StringSliceFlag{
	Name:        "tag",           // --tag
	Aliases:     []string{"t"},   // -t
	Usage:       "Add a tag (can be specified multiple times)",
	Default:     []string{"default-tag"}, // Optional default slice
	Required:    false,
},
```

Usage: `./myapp process --tag production -t us-east-1 --tag webserver`
In the action, `ctx.StringSlice("tag")` would return `[]string{"production", "us-east-1", "webserver"}`.

## Context API

The `ActionFunc` receives a `*nova.Context` pointer, which is your gateway to runtime information.

```go
func myAction(ctx *nova.Context) error {
	// Accessing Arguments
	// Args() returns positional arguments AFTER flags have been parsed.
	args := ctx.Args() // Type: []string
	if len(args) > 0 {
		fmt.Printf("First argument: %s\n", args[0])
	}

	// Accessing Flag Values
	// Use type-specific methods. Nova checks command flags first, then global flags.
	// Returns the zero value for the type if the flag wasn't found or type mismatches.
	configFile := ctx.String("config") // Checks command's --config, then global --config
	port := ctx.Int("port")            // Checks command's --port, then global --port
	verbose := ctx.Bool("verbose")     // Checks command's --verbose, then global --verbose
	tags := ctx.StringSlice("tag")     // Checks command's --tag, then global --tag

	fmt.Printf("Config: %s, Port: %d, Verbose: %t, Tags: %v\n",
		configFile, port, verbose, tags)

	// Accessing Metadata
	appName := ctx.CLI.Name           // Get the application name
	appVersion := ctx.CLI.Version     // Get the application version
	commandName := ""
	if ctx.Command != nil { // Command is nil if running the global Action
		commandName = ctx.Command.Name // Get the name of the running command
	}

	fmt.Printf("Running command '%s' in app '%s' v%s\n",
		commandName, appName, appVersion)

	return nil
}
```

## Help and Version

-   **Version:** Nova automatically adds a global `--version` flag (and `-v` alias). When used, it prints `AppName version AppVersion` and exits. You don't need to define this flag yourself.
-   **Help:** Nova provides a built-in `help` command.
    -   Running `myapp help` shows the main application help (description, usage, commands, global options).
    -   Running `myapp help <command>` shows detailed help for that specific command (description, usage, arguments, command-specific options, aliases).
-   **Command Help Flags (`-h`/`--help`):** Nova *reserves* the names `help` and `h` for flags within a command's definition (`Command.Flags`). You cannot define flags with these names. Users should use the `help` command (`myapp help <command>`) to get help for a specific command.

## Error Handling

-   `nova.NewCLI(app)`: Returns an error if the `CLI` configuration is invalid (e.g., missing `Name`/`Version`, conflicting reserved names). Check this error before proceeding.
-   `cli.Run(args)`: Returns errors encountered during execution:
    -   Flag parsing errors (e.g., invalid value type).
    -   Validation errors (e.g., missing required flags).
    -   "Unknown command" errors.
    -   Any error returned by your `ActionFunc`.

Always check the error returned by `Run` and handle it appropriately (e.g., log it, exit with a non-zero status).

```go
// In main()
cli, err := nova.NewCLI(app)
if err != nil {
    log.Fatalf("Initialization error: %v", err) // Handle NewCLI error
}

if err := cli.Run(os.Args); err != nil {
    log.Fatalf("Runtime error: %v", err) // Handle Run error
}
```

## Full Example Application

This example demonstrates global flags, command flags, required flags, default values, arguments, and context usage.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xlc-dev/nova/nova"
)

func main() {
	app := &nova.CLI{
		Name:        "greeter-server",
		Version:     "1.1.0",
		Description: "A versatile app to greet users or start a server.",
		Authors:     "Nova Developer",
		GlobalFlags: []nova.Flag{
			&nova.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug output globally",
			},
			&nova.StringFlag{
				Name:    "log-file",
				Usage:   "Path to write logs",
				Default: "",
			},
		},
		Commands: []*nova.Command{
			{
				Name:      "hello",
				Aliases:   []string{"hi"},
				Usage:     "Prints a customizable greeting",
				ArgsUsage: "<name>", // Expects one argument
				Action:    helloAction,
				Flags: []nova.Flag{
					&nova.StringFlag{
						Name:    "greeting",
						Aliases: []string{"g"},
						Usage:   "Greeting phrase",
						Default: "Hello",
					},
					&nova.IntFlag{
						Name:    "repeat",
						Aliases: []string{"r"},
						Usage:   "Number of times to repeat the greeting",
						Default: 1,
					},
				},
			},
			{
				Name:        "serve",
				Usage:       "Starts an HTTP(S) server",
				Description: "Starts a simple web server on the specified host and port.",
				Action:      serveAction,
				Flags: []nova.Flag{
					&nova.IntFlag{
						Name:     "port",
						Aliases:  []string{"p"},
						Usage:    "Server port number",
						Required: true, // Port is mandatory for serve
					},
					&nova.StringFlag{
						Name:    "host",
						Usage:   "Server host address",
						Default: "127.0.0.1",
					},
					&nova.StringFlag{
						Name:  "tls-cert",
						Usage: "Path to TLS certificate file (enables HTTPS)",
					},
				},
			},
		},
	}

	cli, err := nova.NewCLI(app)
	if err != nil {
		log.Fatalf("Failed to initialize CLI: %v", err)
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

// logMessage now takes nova.Context to access global flags
func logMessage(ctx *nova.Context, format string, a ...any) {
	// Get global flag values from context
	debugMode := ctx.Bool("debug")
	logFileDest := ctx.String("log-file")

	msg := fmt.Sprintf(format, a...)

	// Log to file if specified
	if logFileDest != "" {
		f, err := os.OpenFile(
			logFileDest,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err == nil {
			defer f.Close()
			if _, writeErr := f.WriteString(msg + "\n"); writeErr != nil {
				// Log the error of writing to the file to stderr
				fmt.Fprintf(
					os.Stderr,
					"Error writing to log file %s: %v\n",
					logFileDest,
					writeErr,
				)
			}
		} else {
			// Log the error of opening the file to stderr
			fmt.Fprintf(
				os.Stderr,
				"Error opening log file %s: %v\n",
				logFileDest,
				err,
			)
		}
	}

	// Conditional console output:
	// Print if debug mode is on, or if no log file is specified.
	if debugMode || logFileDest == "" {
		fmt.Println(msg)
	}
}

func helloAction(ctx *nova.Context) error {
	// Access global flag "debug" via context
	if ctx.Bool("debug") {
		logMessage(ctx, "Debug: Running 'hello' command")
	}

	args := ctx.Args()
	if len(args) == 0 {
		// ArgsUsage suggests <name>, ActionFunc should still validate
		return fmt.Errorf("missing required argument: <name>")
	}
	name := args[0]

	// Access command-specific flags via context
	greeting := ctx.String("greeting")
	repeat := ctx.Int("repeat")

	message := fmt.Sprintf("%s, %s!", greeting, name)

	for range repeat {
		logMessage(ctx, message)
	}

	if ctx.Bool("debug") {
		logMessage(ctx, "Debug: 'hello' command finished")
	}
	return nil
}

func serveAction(ctx *nova.Context) error {
	if ctx.Bool("debug") {
		logMessage(ctx, "Debug: Running 'serve' command")
		// Access global flag "log-file" via context for logging
		logMessage(ctx, "Debug: Global log file: %q", ctx.String("log-file"))
	}

	// Access command-specific flags via context
	host := ctx.String("host")
	port := ctx.Int("port") // Required flag, nova should ensure it's present
	tlsCert := ctx.String("tls-cert")

	protocol := "HTTP"
	if tlsCert != "" {
		protocol = "HTTPS"
	}

	logMessage(ctx, "Starting %s server on %s:%d...", protocol, host, port)

	if tlsCert != "" {
		logMessage(ctx, "Using TLS certificate: %s", tlsCert)
		fmt.Printf(
			"Simulating HTTPS server start on %s:%d with cert %s\n",
			host,
			port,
			tlsCert,
		)
	} else {
		// Add actual HTTP server start logic here
		fmt.Printf("Simulating HTTP server start on %s:%d\n", host, port)
	}

	logMessage(ctx, "Server running (simulation). Press Ctrl+C to stop.")
	// Block or wait for server goroutine in a real app
	select {} // Simulate running server
}
```
