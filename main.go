package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"slices"

	"github.com/xlc-dev/nova/nova"
	"github.com/xlc-dev/nova/templates"
)

var logo = " _   _\n" +
	"| \\ | |\n" +
	"|  \\| |  ___ __   __ __ _\n" +
	"| . ` | / _ \\\\ \\ / // _` |\n" +
	"| |\\  || (_) |\\ V /| (_| |\n" +
	"|_| \\_| \\___/  \\_/  \\____|\n"

// projectGenerator holds the state and logic for creating a new project.
type projectGenerator struct {
	projectName string
	projectDir  string
	isVerbose   bool
	reader      *bufio.Reader

	// User choices
	templateChoice string
	dbAdapter      string
	useGit         bool
	useMakefile    bool
}

// newProjectGenerator is the constructor for our generator.
func newProjectGenerator(
	projectName string,
	isVerbose bool,
) (*projectGenerator, error) {
	projectDir, err := filepath.Abs(projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(projectDir); err == nil {
		return nil, fmt.Errorf("directory '%s' already exists", projectName)
	}

	return &projectGenerator{
		projectName: projectName,
		projectDir:  projectDir,
		isVerbose:   isVerbose,
		reader:      bufio.NewReader(os.Stdin),
	}, nil
}

// run executes the full project generation workflow.
func (g *projectGenerator) run() error {
	if err := os.MkdirAll(g.projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}
	fmt.Printf("Created directory: %s\n", g.projectDir)

	if err := g.promptUserForChoices(); err != nil {
		return err
	}

	if err := g.createFromTemplate(); err != nil {
		return err
	}

	if err := g.initializeGoModule(); err != nil {
		return err
	}

	if g.dbAdapter != "" {
		if err := g.setupDatabase(); err != nil {
			return err
		}
	}

	if g.useGit {
		if err := g.initializeGit(); err != nil {
			return err
		}
	}

	if g.useMakefile {
		if err := g.createMakefile(); err != nil {
			return err
		}
	}

	return nil
}

// promptUserForChoices handles all interactive questions.
func (g *projectGenerator) promptUserForChoices() error {
	var err error

	// Template selection
	templateOptions := []string{"minimal", "todo", "structured"}
	g.templateChoice, err = g.ask(
		"Select template (minimal, todo, structured): ",
		templateOptions,
	)
	if err != nil {
		return err
	}

	dbAnswer, err := g.ask(
		"Would you like to add a database? (y/n): ",
		[]string{"y", "n"},
	)
	if err != nil {
		return err
	}
	if dbAnswer == "y" {
		adapterOptions := []string{"sqlite", "postgres", "mysql"}
		g.dbAdapter, err = g.ask(
			"Select database adapter (sqlite, postgres, mysql): ",
			adapterOptions,
		)
		if err != nil {
			return err
		}
	}

	gitAnswer, err := g.ask(
		"Initialize a git repository? (y/n): ",
		[]string{"y", "n"},
	)
	if err != nil {
		return err
	}
	g.useGit = (gitAnswer == "y")

	makefileAnswer, err := g.ask(
		"Would you like to add a Makefile? (y/n): ",
		[]string{"y", "n"},
	)
	if err != nil {
		return err
	}
	g.useMakefile = (makefileAnswer == "y")

	return nil
}

// createFromTemplate executes the chosen template creation function.
func (g *projectGenerator) createFromTemplate() error {
	createFns := map[string]func(projectDir string, isVerbose bool, dbAdapter string) error{
		"minimal":    templates.CreateMinimal,
		"todo":       templates.CreateTODO,
		"structured": templates.CreateStructured,
	}
	createFn, exists := createFns[g.templateChoice]
	if !exists {
		return fmt.Errorf("internal error: unknown template '%s'", g.templateChoice)
	}

	fmt.Println("Creating project from template...")
	err := createFn(g.projectDir, g.isVerbose, g.dbAdapter)
	if err != nil {
		return fmt.Errorf("failed to create from template: %w", err)
	}
	return nil
}

// initializeGoModule runs `go mod init`, `go mod tidy`, and `go fmt`.
func (g *projectGenerator) initializeGoModule() error {
	fmt.Println("Initializing Go module...")
	if err := g.runCommand("go", "mod", "init", g.projectName); err != nil {
		return fmt.Errorf("failed to run 'go mod init': %w", err)
	}
	if err := g.runCommand("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %w", err)
	}
	if err := g.runCommand("go", "fmt", "./..."); err != nil {
		return fmt.Errorf("failed to run 'go fmt': %w", err)
	}
	return nil
}

// setupDatabase creates the .env file with the correct DATABASE_URL.
func (g *projectGenerator) setupDatabase() error {
	dbConfig := map[string]string{
		"sqlite":   "DATABASE_URL=file:database.db?cache=shared&mode=rwc\n",
		"postgres": "DATABASE_URL=postgres://user:password@localhost/dbname?sslmode=disable\n",
		"mysql":    "DATABASE_URL=mysql://user:password@tcp(127.0.0.1:3306)/dbname\n",
	}
	envURL, exists := dbConfig[g.dbAdapter]
	if !exists {
		return fmt.Errorf("internal error: unknown db adapter '%s'", g.dbAdapter)
	}

	fmt.Println("Creating .env file...")
	envPath := filepath.Join(g.projectDir, ".env")
	if err := os.WriteFile(envPath, []byte(envURL), 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}
	return nil
}

// initializeGit creates a .gitignore file and runs `git init`.
func (g *projectGenerator) initializeGit() error {
	fmt.Println("Initializing git repository...")
	gitignore := fmt.Sprintf(`.DS_Store
Thumbs.db
*.exe
*.exe~
*.dll
*.so
*.dylib
.idea/
.vscode/
*~
*.swp
.env
*.db
%s
`, g.projectName)
	gitignorePath := filepath.Join(g.projectDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignore), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	if err := g.runCommand("git", "init"); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}
	fmt.Println("Git repository and .gitignore initialized successfully.")
	return nil
}

// createMakefile generates and writes the Makefile.
func (g *projectGenerator) createMakefile() error {
	fmt.Println("Creating Makefile...")
	var buildTarget string
	if g.templateChoice == "structured" {
		buildTarget = fmt.Sprintf("./cmd/%s", g.projectName)
	} else {
		buildTarget = "."
	}
	makefileContent := fmt.Sprintf(`BINARY_NAME=%s

.PHONY: build clean fmt test help
default: build

build:
	@go build -o $(BINARY_NAME) %s

clean:
	@rm -f $(BINARY_NAME)

fmt:
	@goimports -w .
	@go fmt ./...

test:
	@go test ./... -v

help:
	@echo "Available Make targets:"
	@echo "  build : Build the Go application (default)"
	@echo "  clean : Remove the built binary ($(BINARY_NAME))"
	@echo "  fmt   : Format Go source code (using goimports)"
	@echo "  test  : Run Go tests"
	@echo "  help  : Show this help message"
`, g.projectName, buildTarget)
	makefile_path := filepath.Join(g.projectDir, "Makefile")
	if err := os.WriteFile(makefile_path, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Makefile: %w", err)
	}
	fmt.Println("Makefile created successfully.")
	return nil
}

// runCommand is a helper to execute external commands in the project directory.
func (g *projectGenerator) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = g.projectDir
	if g.isVerbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command '%s %s' failed: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

// ask is a generic helper for prompting the user and validating input.
func (g *projectGenerator) ask(
	prompt string,
	validOptions []string,
) (string, error) {
	for {
		fmt.Print(prompt)
		input, err := g.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read user input: %w", err)
		}
		input = strings.TrimSpace(strings.ToLower(input))

		if slices.Contains(validOptions, input) {
			return input, nil
		}
		fmt.Printf(
			"Invalid input. Please choose one of: %v\n",
			validOptions,
		)
	}
}

func main() {
	config := &nova.CLI{
		Name:        "Nova",
		Version:     "0.0.1",
		Description: "A modern do-it-all Golang framework to create REST APIs with ease.",
		GlobalFlags: []nova.Flag{
			&nova.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "Enable verbose output",
			},
		},

		Commands: []*nova.Command{
			{
				Name:        "gendoc",
				Usage:       "Generate markdown reference docs from Go code",
				Description: "Parses Go source comments in a specified directory and outputs Markdown documentation. This is made for Nova's own use, but you can try it out if you like :)",
				Flags: []nova.Flag{
					&nova.StringFlag{
						Name:    "input",
						Default: "./nova",
						Usage:   "Directory containing the Go package source files",
					},
					&nova.StringFlag{
						Name:    "output",
						Default: "./docs/src/reference.md",
						Usage:   "Path to the output Markdown file",
					},
				},
				Action: func(ctx *nova.Context) error {
					inputDir := ctx.String("input")
					outputFile := ctx.String("output")
					err := generateReferenceMarkdown(inputDir, outputFile)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
						return err
					}
					return nil
				},
			},
			{
				Name:        "migrate",
				Usage:       "Manages database migrations (up, down, new)",
				Description: "Applies pending migrations (up), rolls back migrations (down), or creates a new migration file (new).",
				ArgsUsage:   "<up|down|new> [steps|migration_name]",
				Action: func(ctx *nova.Context) error {
					// Load environment variables from a .env file if it exists
					if err := nova.LoadDotenv(); err != nil {
						log.Fatalf("Error loading .env: %v", err)
					}

					// Read the full DSN from DATABASE_URL
					dsn := os.Getenv("DATABASE_URL")
					if dsn == "" {
						return fmt.Errorf("DATABASE_URL environment variable is not set")
					}

					args := ctx.Args()
					if len(args) < 1 {
						return fmt.Errorf("expected a migration action: up, down, or new")
					}
					action := args[0]

					// Determine the driver based on the DSN prefix
					var driver string
					if strings.HasPrefix(dsn, "postgres://") {
						driver = "pq"
					} else if strings.HasPrefix(dsn, "mysql://") {
						driver = "mysql"
					} else if strings.HasPrefix(dsn, "file:") || strings.HasSuffix(dsn, ".db") {
						driver = "sqlite"
					} else {
						return fmt.Errorf("unsupported DSN: %s", dsn)
					}

					db, err := sql.Open(driver, dsn)
					if err != nil {
						return err
					}
					defer db.Close()

					switch action {
					case "up":
						steps := 0
						if len(args) > 1 {
							steps, err = strconv.Atoi(args[1])
							if err != nil {
								return fmt.Errorf("invalid steps: %w", err)
							}
						}
						return nova.MigrateUp(db, steps)
					case "down":
						steps := 1
						if len(args) > 1 {
							steps, err = strconv.Atoi(args[1])
							if err != nil {
								return fmt.Errorf("invalid steps: %w", err)
							}
						}
						return nova.MigrateDown(db, steps)
					case "new":
						if len(args) < 2 {
							return fmt.Errorf("migration name required for 'new' action")
						}
						migrationName := args[1]
						return nova.CreateNewMigration(migrationName)
					default:
						return fmt.Errorf("unknown migration action: %s", action)
					}
				},
			},
			{
				Name:        "new",
				Aliases:     []string{"n"},
				Usage:       "Creates a new project",
				Description: "Creates a new project directory with the basic structure.",
				ArgsUsage:   "<project-name>",
				Action: func(ctx *nova.Context) (err error) {
					if len(ctx.Args()) != 1 {
						return fmt.Errorf("expected exactly one argument: <project-name>")
					}
					projectName := ctx.Args()[0]

					g, err := newProjectGenerator(projectName, ctx.Bool("verbose"))
					if err != nil {
						return err
					}

					defer func() {
						if err != nil {
							// g.cleanupOnFailure()
						}
					}()

					err = g.run()
					if err != nil {
						return fmt.Errorf("failed to create project: %w", err)
					}

					fmt.Printf("\nProject '%s' created successfully.\n", projectName)
					return nil
				},
			},
		},
	}

	cli, err := nova.NewCLI(config)
	if err != nil {
		log.Fatalf("Failed to initialize CLI: %v", err)
	}

	fmt.Println(logo)
	err = cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
