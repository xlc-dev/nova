package nova

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// migrationsFolder is the folder where migration SQL files are stored.
	migrationsFolder = "migrations"
	// versionTable is the name of the table used to store the current migration version.
	versionTable = "schema_version"
)

// updateVersionTable updates the current migration version in the version table.
func updateVersionTable(db *sql.DB, version int64) error {
	_, err := db.Exec("UPDATE "+versionTable+" SET version = $1", version)
	return err
}

// getMigrationFiles returns a list of SQL file names found in the migrations folder.
func getMigrationFiles() ([]string, error) {
	entries, err := os.ReadDir(migrationsFolder)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// parseMigrationVersion attempts to parse a numeric migration version from the given file name.
// The file name is expected to start with a numeric version (typically a timestamp),
// followed by an underscore.
func parseMigrationVersion(fileName string) (int64, error) {
	parts := strings.Split(fileName, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid migration file name: %s", fileName)
	}
	versionStr := parts[0]
	v, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing migration version in %s: %v", fileName, err)
	}
	return v, nil
}

// splitSQLStatements splits the content of a migration file into separate SQL statements
// based on the provided action ("up" or "down").
// The migration file must contain delimiters:
//
//	"-- migrate:up"   marks the start of the statements for applying (up) migration.
//	"-- migrate:down" marks the start of the statements for rolling back (down) migration.
func splitSQLStatements(content, action string) ([]string, error) {
	upDelimiter := "-- migrate:up"
	downDelimiter := "-- migrate:down"
	var statements []string
	var capture bool
	var currentStmt string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == upDelimiter {
			capture = action == "up"
			continue
		}
		if trimmed == downDelimiter {
			capture = action == "down"
			continue
		}
		if !capture {
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		currentStmt += " " + line
		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, strings.TrimSpace(currentStmt))
			currentStmt = ""
		}
	}
	if currentStmt != "" {
		statements = append(statements, strings.TrimSpace(currentStmt))
	}
	return statements, nil
}

// getCurrentVersion retrieves the current migration version from the version table.
// If the version table does not exist, it creates it and sets the version to 0.
func getCurrentVersion(db *sql.DB) (int64, error) {
	var version int64
	err := db.QueryRow("SELECT version FROM " + versionTable).Scan(&version)
	if err != nil {
		// If the table does not exist, create it and insert a version of 0.
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + versionTable + " (version INTEGER)")
		if err != nil {
			return 0, err
		}
		_, err = db.Exec("INSERT INTO " + versionTable + " (version) VALUES (0)")
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	return version, nil
}

// MigrateUp applies pending migrations to the provided database.
// It reads migration files from the migrations folder, sorts them in ascending order,
// and applies each migration with a version greater than the current database version.
// The parameter steps indicates how many migrations to apply:
// if steps is 0 the function applies all pending migrations.
//
// Parameters:
//
//	db    - The database handle (from database/sql).
//	steps - The maximum number of migrations to apply (0 means apply all).
//
// Returns an error if migration fails, otherwise nil.
func MigrateUp(db *sql.DB, steps int) error {
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return err
	}

	files, err := getMigrationFiles()
	if err != nil {
		return err
	}
	// Sort migration files in ascending order (oldest first).
	sort.Strings(files)

	appliedCount := 0
	for _, file := range files {
		// If steps > 0 and we've applied enough migrations, exit loop.
		if steps > 0 && appliedCount >= steps {
			break
		}
		migrationVersion, err := parseMigrationVersion(file)
		if err != nil {
			return err
		}
		if migrationVersion > currentVersion {
			migrationPath := filepath.Join(migrationsFolder, file)
			content, err := os.ReadFile(migrationPath)
			if err != nil {
				return err
			}

			statements, err := splitSQLStatements(string(content), "up")
			if err != nil {
				return err
			}

			for _, stmt := range statements {
				if _, err := db.Exec(stmt); err != nil {
					return fmt.Errorf("error executing migration %s: %v", file, err)
				}
			}

			currentVersion = migrationVersion
			if err := updateVersionTable(db, currentVersion); err != nil {
				return err
			}

			fmt.Printf("Applied migration %s\n", file)
			appliedCount++
		}
	}

	if appliedCount == 0 {
		fmt.Println("No pending migrations to apply.")
	} else {
		fmt.Printf("Successfully applied %d migration(s).\n", appliedCount)
	}
	return nil
}

// MigrateDown rolls back migrations on the provided database.
// It reads migration files from the migrations folder, sorts them in descending order,
// and applies the rollback (down) statements for each migration file where the migration version
// is less than or equal to the current version. The parameter steps indicates how many migrations to roll back:
// if steps is 0 the function rolls back one migration by default.
//
// Parameters:
//
//	db    - The database handle (from database/sql).
//	steps - The number of migrations to roll back (0 means 1 migration).
//
// Returns an error if rollback fails, otherwise nil.
func MigrateDown(db *sql.DB, steps int) error {
	if steps <= 0 {
		steps = 1
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return err
	}

	files, err := getMigrationFiles()
	if err != nil {
		return err
	}
	// Reverse sort to start from the latest migration.
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	rollbackCount := 0
	for _, file := range files {
		// Stop if we've rolled back the requested number of migrations.
		if rollbackCount >= steps {
			break
		}
		migrationVersion, err := parseMigrationVersion(file)
		if err != nil {
			return err
		}
		if migrationVersion <= currentVersion && currentVersion > 0 {
			migrationPath := filepath.Join(migrationsFolder, file)
			content, err := os.ReadFile(migrationPath)
			if err != nil {
				return err
			}
			statements, err := splitSQLStatements(string(content), "down")
			if err != nil {
				return err
			}
			for _, stmt := range statements {
				if _, err := db.Exec(stmt); err != nil {
					return fmt.Errorf("error executing rollback for %s: %v", file, err)
				}
			}
			// Decrement current version.
			currentVersion = migrationVersion - 1
			if err := updateVersionTable(db, currentVersion); err != nil {
				return err
			}
			fmt.Printf("Rolled back migration %s\n", file)
			rollbackCount++
		}
	}

	if rollbackCount == 0 {
		fmt.Println("No migrations were rolled back.")
	} else {
		fmt.Printf("Successfully rolled back %d migration(s).\n", rollbackCount)
	}
	return nil
}

// CreateNewMigration creates a new migration file with a basic "up" and "down" template.
// The new file is named using the current Unix timestamp followed by an underscore and the provided name.
// The migration file is saved in the migrationsFolder and contains two sections separated by
// "-- migrate:up" (for the migration) and "-- migrate:down" (for the rollback).
//
// Parameters:
//
//	name - A descriptive name for the migration (e.g., "create_users_table").
//
// Returns an error if the file creation fails, otherwise nil.
func CreateNewMigration(name string) error {
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s.sql", timestamp, name)
	filePath := filepath.Join(migrationsFolder, fileName)

	// Ensure the migrations folder exists.
	if _, err := os.Stat(migrationsFolder); os.IsNotExist(err) {
		os.Mkdir(migrationsFolder, os.ModePerm)
	}

	content := "-- migrate:up\n\n-- migrate:down\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Printf("Created new migration: %s\n", fileName)
	return nil
}
