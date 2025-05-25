# Database Migrations

Nova provides a simple, file-based database migration tool accessible directly through the `nova` command-line binary. It allows you to manage database schema changes using plain SQL files stored in a dedicated folder, driven by simple commands.

- **CLI Driven:** Manage migrations using `nova migrate <action>`.
- **Plain SQL:** Write standard SQL for your target database.
- **Version Tracking:** Automatically tracks the applied migration version in a dedicated schema table.
- **Up/Down Migrations:** Each migration file contains SQL for both applying (`up`) and reverting (`down`) the change.
- **Timestamp-Based Ordering:** Migration files are ordered based on a timestamp prefix in their filename.
- **Step Control:** Apply or roll back all pending migrations or a specific number of steps via command arguments.
- **Environment Configuration:** Reads database connection details from the `DATABASE_URL` environment variable (supports `.env` files).
- **Driver Detection:** Automatically detects the database driver (`postgres`, `mysql`, `sqlite`) based on the `DATABASE_URL` prefix.
- **Low dependencies:** Relies only on the Go standard library and database drivers for Postgres, MySQL/MariaDB and SQLite.

## Table of Contents

1.  [Getting Started](#getting-started)
2.  [Configuration (`DATABASE_URL`)](#configuration-database_url)
3.  [Core Concepts](#core-concepts)
    - [The `migrations` Folder](#the-migrations-folder)
    - [Migration File Naming](#migration-file-naming)
    - [Migration File Structure](#migration-file-structure)
    - [Version Tracking (`schema_version` table)](#version-tracking-schema_version-table)
4.  [CLI Commands](#cli-commands)
    - [`nova migrate new`](#nova-migrate-new)
    - [`nova migrate up`](#nova-migrate-up)
    - [`nova migrate down`](#nova-migrate-down)
5.  [Programmatic Usage (Advanced)](#programmatic-usage-advanced)
    - [`CreateNewMigration`](#createnewmigration-programmatic)
    - [`MigrateUp`](#migrateup-programmatic)
    - [`MigrateDown`](#migratedown-programmatic)

## Getting Started

1.  **Set `DATABASE_URL`:** Ensure the `DATABASE_URL` environment variable is set correctly for your target database. You can also place it in a `.env` file in the directory where you run `nova`.

    ```bash
    # Example for PostgreSQL
    export DATABASE_URL="postgres://user:password@host:port/dbname?sslmode=disable"

    # Example for SQLite
    export DATABASE_URL="file:./my_app_data.db"
    ```

2.  **Create the `migrations` Folder:** In the root of your project (or where you run `nova`), create the folder:

    ```bash
    mkdir migrations
    ```

3.  **Create Your First Migration:** Use the `new` command:

    ```bash
    nova migrate new create_users_table
    # Output: Created new migration: migrations/1678886400_create_users_table.sql (timestamp will vary)
    ```

4.  **Edit the SQL File:** Open the generated `.sql` file and add your SQL statements under the appropriate delimiters:

    ```sql
    -- migrations/1678886400_create_users_table.sql

    -- migrate:up
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    -- migrate:down
    DROP TABLE IF EXISTS users;
    ```

5.  **Apply the Migration:** Use the `up` command:
    ```bash
    nova migrate up
    # Output: Applied migration migrations/1678886400_create_users_table.sql
    # Output: Successfully applied 1 migration(s).
    ```
    Your database schema is now updated!

## Configuration (`DATABASE_URL`)

The `nova migrate` command requires the `DATABASE_URL` environment variable to connect to your database.

- **Format:** Use a standard DSN (Data Source Name) URL format.
- **Supported Drivers (Auto-Detected):**
  - PostgreSQL (`postgres://...`)
  - MySQL (`mysql://...`)
  - SQLite (`file:...` or path ending in `.db`)
- **`.env` File:** If a `.env` file exists in the current directory, `nova migrate` will attempt to load environment variables from it. Variables already present in the environment take precedence.

## Core Concepts

These concepts explain how the migration files and tracking work.

### The `migrations` Folder

- Must be named `migrations`.
- Must exist in the directory where you execute the `nova migrate` command or call the programmatic functions.
- Contains all your `.sql` migration files.

### Migration File Naming

- **Format:** `<version>_<descriptive_name>.sql`
- `<version>`: An integer, typically a Unix timestamp (e.g., `1678886400`), used for ordering. Generated automatically by `nova migrate new`.
- `<descriptive_name>`: Briefly explains the migration's purpose (e.g., `create_users_table`).

### Migration File Structure

- Plain SQL file (`.sql`).
- Must contain `-- migrate:up` and `-- migrate:down` delimiters.
- SQL statements under `-- migrate:up` are executed by `nova migrate up`.
- SQL statements under `-- migrate:down` are executed by `nova migrate down`.

```sql
-- migrations/TIMESTAMP_add_email_to_users.sql

-- migrate:up
ALTER TABLE users ADD COLUMN email VARCHAR(255) UNIQUE;
CREATE INDEX idx_users_email ON users(email);

-- migrate:down
DROP INDEX IF EXISTS idx_users_email;
ALTER TABLE users DROP COLUMN IF EXISTS email;
```

### Version Tracking (`schema_version` table)

- The first time `nova migrate up` or `nova migrate down` runs (either via CLI or programmatically), it automatically creates a table named `schema_version` in your database.
- This table stores the `<version>` number of the most recently applied migration.
- The migration commands use this table to determine which migrations need to be applied or rolled back.

## CLI Commands

Manage your database schema using these subcommands of `nova migrate`.

### `nova migrate new`

Creates a new migration file skeleton.

- **Syntax:** `nova migrate new <migration_name>`
- **Arguments:**
  - `<migration_name>` (**Required**): A descriptive name for the migration (e.g., `add_indexes_to_orders`). Use underscores for spaces.
- **Behavior:**
  - Creates the `migrations` folder if needed.
  - Generates a filename: `<timestamp>_<migration_name>.sql`.
  - Writes the basic `-- migrate:up` and `-- migrate:down` template into the file.
- **Example:**
  ```bash
  nova migrate new add_last_login_to_users
  # Creates migrations/1678886402_add_last_login_to_users.sql (example)
  ```

### `nova migrate up`

Applies pending migrations to the database.

- **Syntax:** `nova migrate up [steps]`
- **Arguments:**
  - `[steps]` (Optional): An integer specifying the maximum number of migrations to apply. If omitted or `0`, applies _all_ pending migrations.
- **Behavior:**
  - Connects to the database using `DATABASE_URL`.
  - Reads the current version from the `schema_version` table.
  - Finds all `.sql` files in the `migrations` folder with a version greater than the current version.
  - Sorts these pending migrations chronologically.
  - Executes the SQL statements under `-- migrate:up` for each pending migration, up to the specified `[steps]` limit.
  - Updates the `schema_version` table after each successful migration file application.
  - Prints status messages.
- **Examples:**

  ```sh
  # Apply all pending migrations
  nova migrate up

  # Apply only the next 2 pending migrations
  nova migrate up 2
  ```

### `nova migrate down`

Rolls back previously applied migrations.

- **Syntax:** `nova migrate down [steps]`
- **Arguments:**
  - `[steps]` (Optional): An integer specifying the exact number of migrations to roll back. If omitted or `0`, defaults to rolling back **one** migration.
- **Behavior:**
  - Connects to the database using `DATABASE_URL`.
  - Reads the current version from the `schema_version` table.
  - Finds all `.sql` files in the `migrations` folder with a version less than or equal to the current version.
  - Sorts these applied migrations in _reverse_ chronological order.
  - Executes the SQL statements under `-- migrate:down` for the most recent applied migrations, up to the specified number of `[steps]`.
  - Updates the `schema_version` table after each successful rollback to reflect the new latest version.
  - Prints status messages.
- **Examples:**

  ```sh
  # Roll back the single most recent migration
  nova migrate down
  # OR
  nova migrate down 1

  # Roll back the last 3 applied migrations
  nova migrate down 3
  ```

## Programmatic Usage (Advanced)

While the `nova migrate` CLI command is the recommended way to manage migrations, the underlying functions are exported and can be used directly within your Go code if you need more control or want to embed migration logic into your application's startup sequence or custom tooling.

**Note:** When using these functions programmatically, you are responsible for obtaining the `*sql.DB` database connection handle yourself. The functions do **not** read the `DATABASE_URL` environment variable out of the box.

### `CreateNewMigration`

```go
func CreateNewMigration(name string) error
```

Identical in behavior to the CLI's `new` action, but called from Go code.

- **Parameters:**
  - `name` (string): Descriptive name for the migration.
- **Returns:** An error if file creation fails, otherwise `nil`.

**Usage:**

```go
import "github.com/xlc-dev/nova/nova"

// ... inside your Go code ...
err := nova.CreateNewMigration("seed_initial_data")
if err != nil {
    log.Printf("Failed to create migration file: %v", err)
    // Handle error appropriately
}
```

### `MigrateUp`

```go
func MigrateUp(db *sql.DB, steps int) error
```

Identical in behavior to the CLI's `up` action, but called from Go code.

- **Parameters:**
  - `db` (\*sql.DB): An active database connection handle obtained via `sql.Open`.
  - `steps` (int): Max number of migrations to apply (0 = all).
- **Returns:** An error if any migration step fails, otherwise `nil`.

**Usage:**

```go
import (
    "database/sql"

	_ "github.com/lib/pq" // Import your DB driver
    "github.com/xlc-dev/nova/nova"
)

// ... inside your Go code ...
db, err := sql.Open("postgres", "your_connection_string")
if err != nil { /* handle error */ }
defer db.Close()

log.Println("Applying database migrations...")
// Apply all pending migrations during application startup
err = nova.MigrateUp(db, 0)
if err != nil {
    log.Fatalf("Database migration failed: %v", err) // Critical error on startup
}
log.Println("Database migrations applied successfully.")
```

### `MigrateDown`

```go
func MigrateDown(db *sql.DB, steps int) error
```

Identical in behavior to the CLI's `down` action, but called from Go code.

- **Parameters:**
  - `db` (\*sql.DB): An active database connection handle.
  - `steps` (int): Number of migrations to roll back (0 or negative = 1).
- **Returns:** An error if any rollback step fails, otherwise `nil`.

**Usage:**

```go
import (
    "database/sql"

	_ "github.com/lib/pq" // Import your DB driver
    "github.com/xlc-dev/nova/nova"
)

// ... inside custom tooling or specific Go code ...
db, err := sql.Open("postgres", "your_connection_string")
if err != nil { /* handle error */ }
defer db.Close()

log.Println("Rolling back the last migration...")
// Roll back one migration
err = nova.MigrateDown(db, 1)
if err != nil {
    log.Printf("Rollback failed: %v", err)
    // Handle error appropriately
} else {
    log.Println("Rollback successful.")
}
```

**Important:** When using these functions programmatically, the status messages (`Applied migration...`, `Rolled back migration...`) are still printed to standard output. You might want to capture/redirect stdout if integrating into a larger system where this output is undesirable.
