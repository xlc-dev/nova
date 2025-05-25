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

	"github.com/xlc-dev/nova/nova"
	"github.com/xlc-dev/nova/templates"
)

var logo = " _   _\n" +
	"| \\ | |\n" +
	"|  \\| |  ___ __   __ __ _\n" +
	"| . ` | / _ \\\\ \\ / // _` |\n" +
	"| |\\  || (_) |\\ V /| (_| |\n" +
	"|_| \\_| \\___/  \\_/  \\____|\n"

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
				Action: func(ctx *nova.Context) error {
					// Ensure exactly one argument for the project name
					if len(ctx.Args()) != 1 {
						return fmt.Errorf("expected exactly one argument: <project-name>")
					}
					projectName := ctx.Args()[0]

					// Get the absolute path for the project directory
					projectDir, err := filepath.Abs(projectName)
					if err != nil {
						return fmt.Errorf("failed to get absolute path: %w", err)
					}
					if _, err := os.Stat(projectDir); err == nil {
						return fmt.Errorf("directory '%s' already exists", projectName)
					}

					reader := bufio.NewReader(os.Stdin)

					// Helper function for getting validated input.
					// If validOptions is nil or empty, valid responses are expected to be "y" or "n".
					getInput := func(prompt string, validOptions []string) (string, error) {
						for {
							fmt.Print(prompt)
							input, err := reader.ReadString('\n')
							if err != nil {
								return "", err
							}
							input = strings.TrimSpace(strings.ToLower(input))
							if len(validOptions) == 0 {
								if input == "y" || input == "n" {
									return input, nil
								}
							} else {
								for _, option := range validOptions {
									if input == option {
										return input, nil
									}
								}
							}
							fmt.Printf("Invalid input. Valid options: %v\n", validOptions)
						}
					}

					// Template selection
					templateOptions := []string{"minimal", "todo", "structured"}
					templateChoice, err := getInput(
						"Select template (minimal, todo, structured): ",
						templateOptions,
					)
					if err != nil {
						return err
					}

					// Map the chosen template to the corresponding creation function
					createFns := map[string]func(string, bool, string) error{
						"minimal":    templates.CreateMinimal,
						"todo":       templates.CreateTODO,
						"structured": templates.CreateStructured,
					}
					createFn, exists := createFns[templateChoice]
					if !exists {
						return fmt.Errorf("unknown template: %s", templateChoice)
					}

					// Database setup
					chosenAdapter := ""
					dbAnswer, err := getInput("Would you like to add a database to your project? (y/n): ", nil)
					if err != nil {
						return err
					}

					// dbConfig contains your adapter configuration
					dbConfig := map[string]struct {
						pkg    string
						envURL string
					}{
						"sqlite":   {"modernc.org/sqlite", "DATABASE_URL=file:database.db?cache=shared&mode=rwc\n"},
						"postgres": {"github.com/lib/pq", "DATABASE_URL=postgres://user:password@localhost/dbname?sslmode=disable\n"},
						"mysql":    {"github.com/go-sql-driver/mysql", "DATABASE_URL=mysql://user:password@tcp(127.0.0.1:3306)/dbname\n"},
					}

					if dbAnswer == "y" {
						adapterOptions := []string{"sqlite", "postgres", "mysql"}
						chosenAdapter, err = getInput(
							"Select your database adapter (sqlite, postgres, mysql): ",
							adapterOptions,
						)
						if err != nil {
							return err
						}
					}

					// Create project from template
					if err := createFn(projectName, ctx.Bool("verbose"), chosenAdapter); err != nil {
						return fmt.Errorf("failed to create project: %w", err)
					}

					// Initialize module and add replace directive
					initCmd := exec.Command("go", "mod", "init", projectName)
					initCmd.Dir = projectDir
					if err := initCmd.Run(); err != nil {
						return fmt.Errorf("failed to initialize go module: %w", err)
					}

					// Setup database (if chosen)
					if config, exists := dbConfig[chosenAdapter]; exists {
						if err := os.WriteFile(filepath.Join(projectDir, ".env"), []byte(config.envURL), 0644); err != nil {
							return fmt.Errorf("failed to create .env file: %w", err)
						}
					}

					// Run remaining commands
					for _, cmd := range [][]string{
						{"go", "mod", "tidy"},
						{"go", "fmt", "./..."},
					} {
						command := exec.Command(cmd[0], cmd[1:]...)
						command.Dir = projectDir
						if err := command.Run(); err != nil {
							return fmt.Errorf("failed to run %s: %w", strings.Join(cmd, " "), err)
						}
					}

					// Initialize git if requested
					gitAnswer, err := getInput("Initialize a git repository? (y/n): ", nil)
					if err != nil {
						return err
					}
					if gitAnswer == "y" {
						gitCmd := exec.Command("git", "init")
						gitCmd.Dir = projectDir
						if err := gitCmd.Run(); err != nil {
							return fmt.Errorf("failed to initialize git: %w", err)
						}

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
*.env
*.db
%s
`, projectName)
						if err := os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitignore), 0644); err != nil {
							return fmt.Errorf("failed to create .gitignore: %w", err)
						}
						fmt.Println("Git repository and .gitignore initialized successfully.")
					}

					// Ask about adding a Makefile
					makefileAnswer, err := getInput("Would you like to add a Makefile? (y/n): ", nil)
					if err != nil {
						return err
					}
					if makefileAnswer == "y" {
						makefile := fmt.Sprintf(`BINARY_NAME=%s

.PHONY: build clean fmt test help
default: build

build:
	@go build -o $(BINARY_NAME)

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
`, projectName)
						if err := os.WriteFile(filepath.Join(projectDir, "Makefile"), []byte(makefile), 0644); err != nil {
							return fmt.Errorf("failed to create Makefile: %w", err)
						}
						fmt.Println("Makefile created successfully.")
					}

					fmt.Printf("Project '%s' created successfully.\n", projectName)
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
