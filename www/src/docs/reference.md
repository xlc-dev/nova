{{ title: Nova - Reference }}

{{ include-block: doc.html markdown="true" }}

# Reference

## Table of Contents

- [Constants](#constants)
- [Variables](#variables)
- [Functions](#functions)
- [Types](#types)
  - [`ActionFunc`](#actionfunc)
  - [`AuthValidator`](#authvalidator)
  - [`BasicAuthConfig`](#basicauthconfig)
  - [`BoolFlag`](#boolflag)
  - [`CLI`](#cli)
  - [`CORSConfig`](#corsconfig)
  - [`CSRFConfig`](#csrfconfig)
  - [`CacheControlConfig`](#cachecontrolconfig)
  - [`Command`](#command)
  - [`Components`](#components)
  - [`ConcurrencyLimiterConfig`](#concurrencylimiterconfig)
  - [`Context`](#context)
  - [`DocumentConfig`](#documentconfig)
  - [`ETagConfig`](#etagconfig)
  - [`Element`](#element)
  - [`EnforceContentTypeConfig`](#enforcecontenttypeconfig)
  - [`Flag`](#flag)
  - [`Float64Flag`](#float64flag)
  - [`ForceHTTPSConfig`](#forcehttpsconfig)
  - [`Group`](#group)
  - [`GzipConfig`](#gzipconfig)
  - [`HTMLDocument`](#htmldocument)
  - [`HTMLElement`](#htmlelement)
  - [`HandlerFunc`](#handlerfunc)
  - [`HeaderObject`](#headerobject)
  - [`HealthCheckConfig`](#healthcheckconfig)
  - [`IPFilterConfig`](#ipfilterconfig)
  - [`Info`](#info)
  - [`IntFlag`](#intflag)
  - [`LoggingConfig`](#loggingconfig)
  - [`MaintenanceModeConfig`](#maintenancemodeconfig)
  - [`MaxRequestBodySizeConfig`](#maxrequestbodysizeconfig)
  - [`MediaTypeObject`](#mediatypeobject)
  - [`MethodOverrideConfig`](#methodoverrideconfig)
  - [`Middleware`](#middleware)
  - [`OpenAPI`](#openapi)
  - [`OpenAPIConfig`](#openapiconfig)
  - [`Operation`](#operation)
  - [`ParameterObject`](#parameterobject)
  - [`ParameterOption`](#parameteroption)
  - [`PathItem`](#pathitem)
  - [`RateLimiterConfig`](#ratelimiterconfig)
  - [`RealIPConfig`](#realipconfig)
  - [`RecoveryConfig`](#recoveryconfig)
  - [`RequestBodyObject`](#requestbodyobject)
  - [`RequestIDConfig`](#requestidconfig)
  - [`ResponseContext`](#responsecontext)
  - [`ResponseObject`](#responseobject)
  - [`ResponseOption`](#responseoption)
  - [`RouteOptions`](#routeoptions)
  - [`Router`](#router)
  - [`SchemaObject`](#schemaobject)
  - [`SecurityHeadersConfig`](#securityheadersconfig)
  - [`StringFlag`](#stringflag)
  - [`StringSliceFlag`](#stringsliceflag)
  - [`TimeoutConfig`](#timeoutconfig)
  - [`TrailingSlashRedirectConfig`](#trailingslashredirectconfig)
  - [`ValidationErrors`](#validationerrors)
  - [`bufferingResponseWriterInterceptor`](#bufferingresponsewriterinterceptor)
  - [`contextKey`](#contextkey)
  - [`gzipResponseWriter`](#gzipresponsewriter)
  - [`responseWriterInterceptor`](#responsewriterinterceptor)
  - [`route`](#route)
  - [`schemaGenCtx`](#schemagenctx)
  - [`segment`](#segment)
  - [`stringSliceValue`](#stringslicevalue)
  - [`textNode`](#textnode)
  - [`timeoutResponseWriter`](#timeoutresponsewriter)
  - [`visitor`](#visitor)


## Constants

### `migrationsFolder, versionTable`

```go
const (
	// migrationsFolder is the folder where migration SQL files are stored.
	migrationsFolder = "migrations"
	// versionTable is the name of the table used to store the current migration version.
	versionTable = "schema_version"
)
```

## Variables

### `ErrMissingContentType`

```go
var ErrMissingContentType = fmt.Errorf("missing Content-Type header")
```

ErrMissingContentType indicates the Content-Type header was absent.

### `ErrUnsupportedContentType`

```go
var ErrUnsupportedContentType = fmt.Errorf("unsupported Content-Type")
```

ErrUnsupportedContentType indicates the Content-Type was not in the allowed list.

### `swaggerUIFS`

```go
var swaggerUIFS embed.FS
```

### `validationMessages`

```go
var validationMessages = map[string]map[string]string{
	"en": {
		"required":     "Field '%s' is required",
		"minlength":    "Field '%s' must be at least %d characters long",
		"maxlength":    "Field '%s' must be at most %d characters long",
		"min":          "Field '%s' must be at least %v",
		"max":          "Field '%s' must be at most %v",
		"pattern":      "Field '%s' does not match required pattern",
		"enum":         "Field '%s' must be one of: %s",
		"email":        "Field '%s' must be a valid email address",
		"url":          "Field '%s' must be a valid URL",
		"uuid":         "Field '%s' must be a valid UUID",
		"date-time":    "Field '%s' must be a valid RFC3339 date-time",
		"date":         "Field '%s' must be a valid date (YYYY-MM-DD)",
		"time":         "Field '%s' must be a valid time (HH:MM:SS)",
		"phone":        "Field '%s' must be a valid phone number",
		"password":     "Field '%s' password must be at least 8 characters long",
		"alphanumeric": "Field '%s' must contain only alphanumeric characters",
		"alpha":        "Field '%s' must contain only alphabetic characters",
		"numeric":      "Field '%s' must contain only numeric characters",
		"minItems":     "Field '%s' must have at least %d items",
		"maxItems":     "Field '%s' must have at most %d items",
		"uniqueItems":  "Field '%s' must have unique items",
		"multipleOf":   "Field '%s' must be a multiple of %v",
	},
	"es": {
		"required":     "El campo '%s' es obligatorio",
		"minlength":    "El campo '%s' debe tener al menos %d caracteres",
		"maxlength":    "El campo '%s' debe tener como máximo %d caracteres",
		"min":          "El campo '%s' debe ser al menos %v",
		"max":          "El campo '%s' debe ser como máximo %v",
		"pattern":      "El campo '%s' no coincide con el patrón requerido",
		"enum":         "El campo '%s' debe ser uno de: %s",
		"email":        "El campo '%s' debe ser una dirección de correo válida",
		"url":          "El campo '%s' debe ser una URL válida",
		"uuid":         "El campo '%s' debe ser un UUID válido",
		"date-time":    "El campo '%s' debe ser una fecha-hora RFC3339 válida",
		"date":         "El campo '%s' debe ser una fecha válida (AAAA-MM-DD)",
		"time":         "El campo '%s' debe ser una hora válida (HH:MM:SS)",
		"phone":        "El campo '%s' debe ser un número de teléfono válido",
		"password":     "El campo '%s' debe tener al menos 8 caracteres",
		"alphanumeric": "El campo '%s' debe contener solo caracteres alfanuméricos",
		"alpha":        "El campo '%s' debe contener solo caracteres alfabéticos",
		"numeric":      "El campo '%s' debe contener solo caracteres numéricos",
		"minItems":     "El campo '%s' debe tener al menos %d elementos",
		"maxItems":     "El campo '%s' debe tener como máximo %d elementos",
		"uniqueItems":  "El campo '%s' debe tener elementos únicos",
		"multipleOf":   "El campo '%s' debe ser múltiplo de %v",
	},
	"fr": {
		"required":     "Le champ '%s' est obligatoire",
		"minlength":    "Le champ '%s' doit contenir au moins %d caractères",
		"maxlength":    "Le champ '%s' doit contenir au plus %d caractères",
		"min":          "Le champ '%s' doit être d'au moins %v",
		"max":          "Le champ '%s' doit être d'au plus %v",
		"pattern":      "Le champ '%s' ne correspond pas au motif requis",
		"enum":         "Le champ '%s' doit être l'un de: %s",
		"email":        "Le champ '%s' doit être une adresse email valide",
		"url":          "Le champ '%s' doit être une URL valide",
		"uuid":         "Le champ '%s' doit être un UUID valide",
		"date-time":    "Le champ '%s' doit être une date-heure RFC3339 valide",
		"date":         "Le champ '%s' doit être une date valide (AAAA-MM-JJ)",
		"time":         "Le champ '%s' doit être une heure valide (HH:MM:SS)",
		"phone":        "Le champ '%s' doit être un numéro de téléphone valide",
		"password":     "Le champ '%s' doit contenir au moins 8 caractères",
		"alphanumeric": "Le champ '%s' ne doit contenir que des caractères alphanumériques",
		"alpha":        "Le champ '%s' ne doit contenir que des caractères alphabétiques",
		"numeric":      "Le champ '%s' ne doit contenir que des caractères numériques",
		"minItems":     "Le champ '%s' doit avoir au moins %d éléments",
		"maxItems":     "Le champ '%s' doit avoir au plus %d éléments",
		"uniqueItems":  "Le champ '%s' doit avoir des éléments uniques",
		"multipleOf":   "Le champ '%s' doit être un multiple de %v",
	},
	"de": {
		"required":     "Das Feld '%s' ist erforderlich",
		"minlength":    "Das Feld '%s' muss mindestens %d Zeichen lang sein",
		"maxlength":    "Das Feld '%s' darf höchstens %d Zeichen lang sein",
		"min":          "Das Feld '%s' muss mindestens %v sein",
		"max":          "Das Feld '%s' darf höchstens %v sein",
		"pattern":      "Das Feld '%s' entspricht nicht dem erforderlichen Muster",
		"enum":         "Das Feld '%s' muss eines von: %s sein",
		"email":        "Das Feld '%s' muss eine gültige E-Mail-Adresse sein",
		"url":          "Das Feld '%s' muss eine gültige URL sein",
		"uuid":         "Das Feld '%s' muss eine gültige UUID sein",
		"date-time":    "Das Feld '%s' muss eine gültige RFC3339 Datum-Zeit sein",
		"date":         "Das Feld '%s' muss ein gültiges Datum sein (JJJJ-MM-TT)",
		"time":         "Das Feld '%s' muss eine gültige Zeit sein (HH:MM:SS)",
		"phone":        "Das Feld '%s' muss eine gültige Telefonnummer sein",
		"password":     "Das Feld '%s' muss mindestens 8 Zeichen lang sein",
		"alphanumeric": "Das Feld '%s' darf nur alphanumerische Zeichen enthalten",
		"alpha":        "Das Feld '%s' darf nur alphabetische Zeichen enthalten",
		"numeric":      "Das Feld '%s' darf nur numerische Zeichen enthalten",
		"minItems":     "Das Feld '%s' muss mindestens %d Elemente haben",
		"maxItems":     "Das Feld '%s' darf höchstens %d Elemente haben",
		"uniqueItems":  "Das Feld '%s' muss eindeutige Elemente haben",
		"multipleOf":   "Das Feld '%s' muss ein Vielfaches von %v sein",
	},
	"nl": {
		"required":     "Veld '%s' is verplicht",
		"minlength":    "Veld '%s' moet minimaal %d karakters lang zijn",
		"maxlength":    "Veld '%s' mag maximaal %d karakters lang zijn",
		"min":          "Veld '%s' moet minimaal %v zijn",
		"max":          "Veld '%s' mag maximaal %v zijn",
		"pattern":      "Veld '%s' voldoet niet aan het vereiste patroon",
		"enum":         "Veld '%s' moet een van de volgende zijn: %s",
		"email":        "Veld '%s' moet een geldig e-mailadres zijn",
		"url":          "Veld '%s' moet een geldige URL zijn",
		"uuid":         "Veld '%s' moet een geldige UUID zijn",
		"date-time":    "Veld '%s' moet een geldige RFC3339 datum-tijd zijn",
		"date":         "Veld '%s' moet een geldige datum zijn (JJJJ-MM-DD)",
		"time":         "Veld '%s' moet een geldige tijd zijn (UU:MM:SS)",
		"phone":        "Veld '%s' moet een geldig telefoonnummer zijn",
		"password":     "Veld '%s' wachtwoord moet minimaal 8 karakters lang zijn",
		"alphanumeric": "Veld '%s' mag alleen alfanumerieke karakters bevatten",
		"alpha":        "Veld '%s' mag alleen alfabetische karakters bevatten",
		"numeric":      "Veld '%s' mag alleen numerieke karakters bevatten",
		"minItems":     "Veld '%s' moet minimaal %d items bevatten",
		"maxItems":     "Veld '%s' mag maximaal %d items bevatten",
		"uniqueItems":  "Veld '%s' moet unieke items bevatten",
		"multipleOf":   "Veld '%s' moet een veelvoud zijn van %v",
	},
}
```

Translation messages for validation errors.

## Functions

### `CreateNewMigration`

```go
func CreateNewMigration(name string) error
```

CreateNewMigration creates a new migration file with a basic "up" and "down" template.
The new file is named using the current Unix timestamp followed by an underscore and the provided name.
The migration file is saved in the migrationsFolder and contains two sections separated by
"-- migrate:up" (for the migration) and "-- migrate:down" (for the rollback).
#### Parameters
```go
name - A descriptive name for the migration (e.g., "create_users_table").
```

Returns an error if the file creation fails, otherwise nil.

### `GetBasicAuthUser`

```go
func GetBasicAuthUser(ctx context.Context) string
```

GetBasicAuthUser retrieves the authenticated username from the context, if available
via BasicAuthMiddleware using the default key and if configured to store it.

### `GetBasicAuthUserWithKey`

```go
func GetBasicAuthUserWithKey(ctx context.Context, key contextKey) string
```

GetBasicAuthUserWithKey retrieves the authenticated username from the context
using a specific key.

### `GetCSRFToken`

```go
func GetCSRFToken(ctx context.Context) string
```

GetCSRFToken retrieves the CSRF token from the context, if available via
CSRFMiddleware using the default key. This is the token expected in subsequent
unsafe requests.

### `GetCSRFTokenWithKey`

```go
func GetCSRFTokenWithKey(ctx context.Context, key contextKey) string
```

GetCSRFTokenWithKey retrieves the CSRF token from the context using a specific key.

### `GetRealIP`

```go
func GetRealIP(ctx context.Context) string
```

GetRealIP retrieves the client's real IP from the context, if available via
RealIPMiddleware using the default key.

### `GetRealIPWithKey`

```go
func GetRealIPWithKey(ctx context.Context, key contextKey) string
```

GetRealIPWithKey retrieves the client's real IP from the context using a specific key.

### `GetRequestID`

```go
func GetRequestID(ctx context.Context) string
```

GetRequestID retrieves the request ID from the context, if available via
RequestIDMiddleware using the default key.

### `GetRequestIDWithKey`

```go
func GetRequestIDWithKey(ctx context.Context, key contextKey) string
```

GetRequestIDWithKey retrieves the request ID from the context using a specific key.

### `LoadDotenv`

```go
func LoadDotenv(paths ...string) error
```

LoadDotenv loads variables from a .env file, expands them, and sets them
as environment variables. If no path is provided, ".env" is used.
If the specified file doesn't exist, it is silently ignored.

### `MigrateDown`

```go
func MigrateDown(db *sql.DB, steps int) error
```

MigrateDown rolls back migrations on the provided database.
It reads migration files from the migrations folder, sorts them in descending order,
and applies the rollback (down) statements for each migration file where the migration version
is less than or equal to the current version. The parameter steps indicates how many migrations to roll back:
if steps is 0 the function rolls back one migration by default.
#### Parameters
```go
db    - The database handle (from database/sql).
steps - The number of migrations to roll back (0 means 1 migration).
```

Returns an error if rollback fails, otherwise nil.

### `MigrateUp`

```go
func MigrateUp(db *sql.DB, steps int) error
```

MigrateUp applies pending migrations to the provided database.
It reads migration files from the migrations folder, sorts them in ascending order,
and applies each migration with a version greater than the current database version.
The parameter steps indicates how many migrations to apply:
if steps is 0 the function applies all pending migrations.
#### Parameters
```go
db    - The database handle (from database/sql).
steps - The maximum number of migrations to apply (0 means apply all).
```

Returns an error if migration fails, otherwise nil.

### `Serve`

```go
func Serve(ctx *Context, router http.Handler) error
```

Serve launches the web server (concurrently) with graceful shutdown and live reloading (if enabled).
It wraps key goroutines with recovery blocks to avoid crashes due to unexpected errors.
It also allows for logging customization via context options.

### `bindFormToStruct`

```go
func bindFormToStruct(form url.Values, v any) error
```

bindFormToStruct uses reflection to bind form values to struct fields.
Field names are matched using JSON tags (without omitempty) or struct field names.
Supports automatic type conversion for string, bool, int, and float types.

### `buildPathString`

```go
func buildPathString(segments []segment, prefix string) string
```

buildPathString constructs a URI path string from segments and an optional
parent prefix. Parameters are wrapped in curly braces, e.g., "{id}".

### `checkCommandConflicts`

```go
func checkCommandConflicts(cmd *Command) error
```

checkCommandConflicts validates a command against reserved names/aliases and basic requirements.

### `checkFlagConflicts`

```go
func checkFlagConflicts(f Flag, isCommandFlag bool) error
```

checkFlagConflicts validates a flag against reserved names/aliases and basic requirements.
`isCommandFlag` determines if command-level restrictions apply (reserving help/h).

### `collectRoutes`

```go
func collectRoutes(r *Router, spec *OpenAPI, schemaCtx *schemaGenCtx, parentPath string)
```

collectRoutes traverses the router hierarchy starting at r, and populates
the OpenAPI spec with PathItems and Operations based on registered routes.

### `detectLanguage`

```go
func detectLanguage(acceptLanguage string) string
```

detectLanguage extracts the preferred language from Accept-Language header

### `formatFlagNames`

```go
func formatFlagNames(name string, aliases []string) string
```

formatFlagNames combines the primary name and aliases into a comma-separated string
suitable for help text (e.g., "-n, --name"). It sorts them for consistency.

### `getCurrentVersion`

```go
func getCurrentVersion(db *sql.DB) (int64, error)
```

getCurrentVersion retrieves the current migration version from the version table.
If the version table does not exist, it creates it and sets the version to 0.

### `getMessage`

```go
func getMessage(lang, key string, args ...any) string
```

getMessage gets a localized validation message

### `getMigrationFiles`

```go
func getMigrationFiles() ([]string, error)
```

getMigrationFiles returns a list of SQL file names found in the migrations folder.

### `hasAllowedExtension`

```go
func hasAllowedExtension(filename string, exts []string) bool
```

hasAllowedExtension returns true if filename ends with any of the allowed extensions.

### `joinPaths`

```go
func joinPaths(a, b string) string
```

joinPaths concatenates two URL path segments ensuring exactly one "/" between them.
It handles edge cases where either path is empty and normalizes trailing/leading slashes.

### `matchSegments`

```go
func matchSegments(path string, segments []segment) (bool, map[string]string)
```

matchSegments checks if the given URL path matches the compiled segments.
If the path matches, it returns true and a map of parameter names to their values.
Parameter validation is performed using regex patterns if specified.

### `parseFlagSet`

```go
func parseFlagSet(cli *CLI, flags []Flag, args []string, name string, output io.Writer) (*flag.FlagSet, error)
```

parseFlagSet handles the common logic of creating a flag.FlagSet, applying flags to it,
and parsing arguments. It returns the parsed set and any parsing error.

### `parseMigrationVersion`

```go
func parseMigrationVersion(fileName string) (int64, error)
```

parseMigrationVersion attempts to parse a numeric migration version from the given file name.
The file name is expected to start with a numeric version (typically a timestamp),
followed by an underscore.

### `printFlagsHelp`

```go
func printFlagsHelp(w io.Writer, title string, flags []Flag, isCommandHelp bool)
```

printFlagsHelp formats and prints a list of flags to the writer.

### `recompile`

```go
func recompile(verbose bool) error
```

recompile triggers the recompilation of the application.

### `setFieldValue`

```go
func setFieldValue(field reflect.Value, value string) error
```

setFieldValue sets a reflect.Value based on its type and the string value from form data.
Handles automatic type conversion for common types including special checkbox handling.

### `splitSQLStatements`

```go
func splitSQLStatements(content, action string) ([]string, error)
```

splitSQLStatements splits the content of a migration file into separate SQL statements
based on the provided action ("up" or "down").
The migration file must contain delimiters:
```go
"-- migrate:up"   marks the start of the statements for applying (up) migration.
"-- migrate:down" marks the start of the statements for rolling back (down) migration.
```

### `updateVersionTable`

```go
func updateVersionTable(db *sql.DB, version int64) error
```

updateVersionTable updates the current migration version in the version table.

### `validateFlags`

```go
func validateFlags(flags []Flag, set *flag.FlagSet) error
```

validateFlags iterates through a list of user-defined flags and calls their Validate method
against the parsed flag set. It collects and returns any validation errors.

### `validateRequiredFlag`

```go
func validateRequiredFlag(set *flag.FlagSet, name string, required, checkEmpty bool) error
```

validateRequiredFlag provides common validation logic for required flags (String, Int, Float64).
It checks if the flag was set and, optionally, if its value is non-empty (for strings).

### `validateRequiredSliceFlag`

```go
func validateRequiredSliceFlag(set *flag.FlagSet, name string, required bool) error
```

validateRequiredSliceFlag handles validation for required slice flags.

### `validateStruct`

```go
func validateStruct(val any, lang string) error
```

validateStruct walks a struct’s exported fields, applies all tag‐based
validations, and returns every violation found rather than failing fast.
val may be a struct or a pointer to struct. lang selects the locale for
error messages.

### `vlog`

```go
func vlog(verbose bool, format string, a ...any)
```

vlog prints messages only if verbose is enabled.

### `watchAndCompile`

```go
func watchAndCompile(dir string, verbose bool, exts []string) error
```

watchAndCompile watches the given directory for file changes that match allowed extensions
and triggers a recompilation when a change is detected.

### `watchAndReload`

```go
func watchAndReload(
	dir, exe string,
	verbose bool,
	exts []string,
	reloadCh chan<- struct{},
) error
```

watchAndReload watches dir for changes to allowed extensions, rebuilds this binary
in‐place and then signals reloadCh once so Serve can exec the new version.

## Types

### `ActionFunc`

```go
type ActionFunc func(ctx *Context) error
```

ActionFunc defines the function signature for CLI actions (both global and command-specific).
It receives a Context containing parsed flags and arguments.

### `AuthValidator`

```go
type AuthValidator func(username, password string) bool
```

AuthValidator is a function type that validates the provided username and password.
It returns true if the credentials are valid.

### `BasicAuthConfig`

```go
type BasicAuthConfig struct {
	// Realm is the authentication realm presented to the client. Defaults to "Restricted".
	Realm string
	// Validator is the function used to validate username/password combinations. Required.
	Validator AuthValidator
	// StoreUserInContext determines whether to store the validated username in the request context.
	// Defaults to false.
	StoreUserInContext bool
	// ContextKey is the key used to store the username if StoreUserInContext is true.
	// Defaults to the package's internal basicAuthUserKey.
	ContextKey contextKey
}
```

BasicAuthConfig holds configuration for BasicAuthMiddleware.

### `BoolFlag`

```go
type BoolFlag struct {
	// Name is the primary identifier for the flag (e.g., "verbose"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "V"). Single letters are used as -V.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value. Note: Presence of the flag usually implies true.
	Default bool
	// Required is ignored for BoolFlag as presence implies the value.
	Required bool
}
```

BoolFlag defines a flag that acts as a boolean switch (true if present, false otherwise).

#### Methods

#### `Apply`

```go
func (f *BoolFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the bool flag with the flag.FlagSet.

#### `GetAliases`

```go
func (f *BoolFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

#### `GetName`

```go
func (f *BoolFlag) GetName() string
```

GetName returns the primary name of the flag.

#### `IsRequired`

```go
func (f *BoolFlag) IsRequired() bool
```

IsRequired always returns false for boolean flags.

#### `String`

```go
func (f *BoolFlag) String() string
```

String returns the help text representation of the flag.

#### `Validate`

```go
func (f *BoolFlag) Validate(set *flag.FlagSet) error
```

Validate always returns nil for boolean flags.

### `CLI`

```go
type CLI struct {
	// Name of the application (Required).
	Name string
	// Version string for the application (Required). Displayed with the --version flag.
	Version string
	// Description provides a brief explanation of the application's purpose, shown in help text.
	Description string
	// Commands is a list of commands the application supports.
	Commands []*Command
	// Action is the default function to run when no command is specified.
	// If nil and no command is given, help is shown.
	Action ActionFunc
	// GlobalFlags are flags applicable to the application as a whole, before any command.
	// Flags named "version" or aliased "v" are reserved and cannot be used.
	// The "help" flag/alias is handled solely by the built-in 'help' command.
	GlobalFlags []Flag
	// Authors lists the application's authors, shown in the main help text.
	Authors string

	// internal storage for the parsed global flag set
	globalSet *flag.FlagSet
	// internal storage for the built-in nova flags
	internalGlobalFlags []Flag
	// Internal storage for flag values (pointers needed by std/flag)
	internalFlagPtrs map[string]any // Map flag name to pointer (*string, *int, etc.)
}
```

CLI represents the command-line application structure.
It holds the application's metadata, commands, global flags, and the main action.
Name and Version are required fields when creating a CLI instance via NewCLI.

#### Associated Functions

#### `NewCLI`

```go
func NewCLI(cli *CLI) (*CLI, error)
```

NewCLI creates and validates a new CLI application instance based on the provided configuration.
The Name and Version fields in the input CLI struct are required.
It checks for conflicts with reserved names/aliases (help command, h alias, version flag, v alias)
and basic flag/command requirements.
Returns the validated CLI instance or an error if validation fails.

#### Methods

#### `AddCommand`

```go
func (c *CLI) AddCommand(cmd *Command) error
```

AddCommand registers a new command with the application *after* initial NewCLI() validation.
It performs the same conflict checks as NewCLI(). It is generally recommended to define
all commands and flags within the struct passed to NewCLI().

#### `Run`

```go
func (c *CLI) Run(arguments []string) error
```

Run executes the CLI application based on the provided command-line arguments.
Call NewCLI() to create and validate the CLI instance before calling Run.
Run parses flags, handles the built-in version flag and help command, validates required flags,
and executes the appropriate action (global or command-specific).

#### `ShowHelp`

```go
func (c *CLI) ShowHelp(w io.Writer)
```

ShowHelp prints the main help message for the application to the specified writer.

#### `createHelpCommand`

```go
func (c *CLI) createHelpCommand() *Command
```

createHelpCommand generates the internal help command definition.

#### `findCommand`

```go
func (c *CLI) findCommand(name string) *Command
```

findCommand searches the CLI's commands list for a command matching the given name or alias.
Returns the *Command if found, otherwise nil. It checks commands added by the user
AND the internal help command.

#### `setupInternalGlobalFlags`

```go
func (c *CLI) setupInternalGlobalFlags()
```

setupInternalGlobalFlags defines and stores the built-in global flag.

### `CORSConfig`

```go
type CORSConfig struct {
	// AllowedOrigins is a list of origins that are allowed to make cross-site requests.
	// An origin of "*" allows any origin. Defaults to allowing no origins (empty list).
	AllowedOrigins []string
	// AllowedMethods is a list of HTTP methods that are allowed.
	// Defaults to "GET, POST, PUT, DELETE, PATCH, OPTIONS".
	AllowedMethods []string
	// AllowedHeaders is a list of request headers that clients are allowed to send.
	// Defaults to "Content-Type, Authorization, X-Request-ID". Use "*" to allow any header.
	AllowedHeaders []string
	// ExposedHeaders is a list of response headers that clients can access.
	// Defaults to an empty list.
	ExposedHeaders []string
	// AllowCredentials indicates whether the browser should include credentials (like cookies)
	// with the request. Cannot be used with AllowedOrigins = ["*"]. Defaults to false.
	AllowCredentials bool
	// MaxAgeSeconds specifies how long the result of a preflight request can be cached
	// in seconds. Defaults to 86400 (24 hours). A value of -1 disables caching.
	MaxAgeSeconds int
}
```

CORSConfig holds configuration for CORSMiddleware.

### `CSRFConfig`

```go
type CSRFConfig struct {
	// Logger specifies the logger instance for errors. Defaults to log.Default().
	Logger *log.Logger
	// FieldName is the name of the form field to check for the CSRF token.
	// Defaults to "csrf_token".
	FieldName string
	// HeaderName is the name of the HTTP header to check for the CSRF token.
	// Defaults to "X-CSRF-Token".
	HeaderName string
	// CookieName is the name of the cookie used to store the CSRF secret.
	// Defaults to "_csrf". It should be HttpOnly.
	CookieName string
	// ContextKey is the key used to store the generated token in the request context.
	// Defaults to the package's internal csrfTokenKey.
	ContextKey contextKey
	// ErrorHandler is called when CSRF validation fails.
	// Defaults to sending a 403 Forbidden response.
	ErrorHandler http.HandlerFunc
	// CookiePath sets the path attribute of the CSRF cookie. Defaults to "/".
	CookiePath string
	// CookieDomain sets the domain attribute of the CSRF cookie. Defaults to "".
	CookieDomain string
	// CookieMaxAge sets the max-age attribute of the CSRF cookie.
	// Defaults to 12 hours.
	CookieMaxAge time.Duration
	// CookieSecure sets the secure attribute of the CSRF cookie.
	// Defaults to false (for HTTP testing). Set to true in production with HTTPS.
	CookieSecure bool
	// CookieSameSite sets the SameSite attribute of the CSRF cookie.
	// Defaults to http.SameSiteLaxMode.
	CookieSameSite http.SameSite
	// TokenLength is the byte length of the generated token. Defaults to 32.
	TokenLength int
	// SkipMethods is a list of HTTP methods to skip CSRF checks for.
	// Defaults to ["GET", "HEAD", "OPTIONS", "TRACE"].
	SkipMethods []string
}
```

CSRFConfig holds configuration for CSRFMiddleware.

### `CacheControlConfig`

```go
type CacheControlConfig struct {
	// CacheControlValue is the string to set for the Cache-Control header. Required.
	// Example values: "no-store", "no-cache", "public, max-age=3600"
	CacheControlValue string
}
```

CacheControlConfig holds configuration for CacheControlMiddleware.

### `Command`

```go
type Command struct {
	// Name is the primary identifier for the command (Required).
	Name string
	// Aliases provide alternative names for invoking the command. "h" is reserved.
	Aliases []string
	// Usage provides a short description of the command's purpose, shown in the main help list.
	Usage string
	// Description gives a more detailed explanation of the command, shown in the command's help.
	Description string
	// ArgsUsage describes the expected arguments for the command, shown in the command's help.
	// Example: "<input-file> [output-file]"
	ArgsUsage string
	// Action is the function executed when this command is invoked (Required).
	Action ActionFunc
	// Flags are the options specific to this command.
	// Flags named "help" or aliased "h" are reserved and cannot be used.
	Flags []Flag
}
```

Command defines a specific action the CLI application can perform.
It includes metadata, flags specific to the command, and the action function.
Commands named "help" or aliased "h" are reserved.

#### Methods

#### `ShowHelp`

```go
func (cmd *Command) ShowHelp(w io.Writer, appName string)
```

ShowHelp prints the help message for a specific command to the specified writer.

### `Components`

```go
type Components struct {
	Schemas map[string]*SchemaObject `json:"schemas,omitempty"`
}
```

Components holds reusable schema definitions for the OpenAPI spec.

### `ConcurrencyLimiterConfig`

```go
type ConcurrencyLimiterConfig struct {
	// MaxConcurrent is the maximum number of requests allowed to be processed concurrently. Required.
	MaxConcurrent int
	// WaitTimeout is the maximum duration a request will wait for a slot before failing.
	// If zero or negative, requests wait indefinitely.
	WaitTimeout time.Duration
	// OnLimitExceeded allows custom handling when the concurrency limit is hit and timeout occurs (if set).
	// If nil, sends 503 Service Unavailable.
	OnLimitExceeded func(w http.ResponseWriter, r *http.Request)
}
```

ConcurrencyLimiterConfig holds configuration for ConcurrencyLimiterMiddleware.

### `Context`

```go
type Context struct {
	// CLI points to the parent application instance.
	CLI *CLI
	// Command points to the specific command being executed (nil for the global action).
	Command *Command
	// flagSet holds the parsed flag set relevant to the current context
	// (either command-specific or global if no command is running).
	flagSet *flag.FlagSet
	// args holds the non-flag arguments remaining after flag parsing.
	args []string
	// globalSet always holds the parsed global flag set for accessing global flags
	// even within a command context.
	globalSet *flag.FlagSet
}
```

Context provides access to parsed flags, arguments, and application/command metadata
within an ActionFunc.

#### Methods

#### `Args`

```go
func (c *Context) Args() []string
```

Args returns the non-flag arguments remaining after parsing for the current context.

#### `Bool`

```go
func (c *Context) Bool(name string) bool
```

Bool returns the boolean value of a flag specified by name.
It checks command flags first, then global flags. Returns false if not found or type mismatch.

#### `Float64`

```go
func (c *Context) Float64(name string) float64
```

Float64 returns the float64 value of a flag specified by name.
It checks command flags first, then global flags. Returns 0.0 if not found or type mismatch.

#### `Int`

```go
func (c *Context) Int(name string) int
```

Int returns the integer value of a flag specified by name.
It checks command flags first, then global flags. Returns 0 if not found or type mismatch.

#### `String`

```go
func (c *Context) String(name string) string
```

String returns the string value of a flag specified by name.
It checks command flags first, then global flags. Returns "" if not found or type mismatch.

#### `StringSlice`

```go
func (c *Context) StringSlice(name string) []string
```

StringSlice returns the []string value of a flag specified by name.
It checks command flags first, then global flags. Returns nil if not found or type mismatch.

### `DocumentConfig`

```go
type DocumentConfig struct {
	Lang        string        // Lang attribute for `<html>` tag, defaults to "en".
	Title       string        // Content for `<title>` tag, defaults to "Document".
	Charset     string        // Charset for `<meta charset>`, defaults to "utf-8".
	Viewport    string        // Content for `<meta name="viewport">`, defaults to "width=device-width, initial-scale=1".
	Description string        // Content for `<meta name="description">`. If empty, tag is omitted.
	Keywords    string        // Content for `<meta name="keywords">`. If empty, tag is omitted.
	Author      string        // Content for `<meta name="author">`. If empty, tag is omitted.
	HeadExtras  []HTMLElement // Additional HTMLElements to be included in the `<head>` section.
}
```

DocumentConfig provides configuration options for creating a new HTML document.
Fields left as zero-values (e.g., empty strings) will use sensible defaults
or be omitted if optional (like Description, Keywords, Author).

### `ETagConfig`

```go
type ETagConfig struct {
	// Weak determines if weak ETags (prefixed with W/) should be generated.
	// Weak ETags indicate semantic equivalence, not byte-for-byte identity.
	// Defaults to false (strong ETags).
	Weak bool
	// SkipNoContent determines whether to skip ETag generation/checking for
	// responses with status 204 No Content. Defaults to true.
	SkipNoContent bool
}
```

ETagConfig holds configuration for ETagMiddleware.

### `Element`

```go
type Element struct {
	tag        string
	content    string
	attributes map[string]string
	children   []HTMLElement
	selfClose  bool
}
```

Element represents an HTML element with attributes, content, and child elements.
It supports both self-closing elements (like `<img>`) and container elements (like `<div>`).
Elements can be chained using fluent API methods for convenient construction.

#### Associated Functions

#### `A`

```go
func A(href string, content ...HTMLElement) *Element
```

A creates an `<a>` anchor element.

#### `Abbr`

```go
func Abbr(content ...HTMLElement) *Element
```

Abbr creates an `<abbr>` abbreviation element.

#### `Address`

```go
func Address(content ...HTMLElement) *Element
```

Address creates an `<address>` semantic element.

#### `Article`

```go
func Article(content ...HTMLElement) *Element
```

Article creates an `<article>` semantic element.

#### `Aside`

```go
func Aside(content ...HTMLElement) *Element
```

Aside creates an `<aside>` semantic element.

#### `Audio`

```go
func Audio(content ...HTMLElement) *Element
```

Audio creates an `<audio>` element.

#### `B`

```go
func B(content ...HTMLElement) *Element
```

B creates a `<b>` element for stylistically offset text.

#### `Base`

```go
func Base(href string) *Element
```

Base creates a `<base>` element.

#### `Blockquote`

```go
func Blockquote(content ...HTMLElement) *Element
```

Blockquote creates a `<blockquote>` element.

#### `Body`

```go
func Body(content ...HTMLElement) *Element
```

Body creates a `<body>` element.

#### `Br`

```go
func Br() *Element
```

Br creates a self-closing `<br>` line break element.

#### `Button`

```go
func Button(content ...HTMLElement) *Element
```

Button creates a `<button>` element.

#### `ButtonInput`

```go
func ButtonInput(valueText string) *Element
```

ButtonInput creates an `<input type="button">`.

#### `Caption`

```go
func Caption(content ...HTMLElement) *Element
```

Caption creates a `<caption>` element for a table.

#### `CheckboxInput`

```go
func CheckboxInput(name string) *Element
```

CheckboxInput creates an `<input type="checkbox">`.

#### `Cite`

```go
func Cite(content ...HTMLElement) *Element
```

Cite creates a `<cite>` element.

#### `Code`

```go
func Code(content ...HTMLElement) *Element
```

Code creates a `<code>` element.

#### `Col`

```go
func Col() *Element
```

Col creates a `<col>` element. It's self-closing.

#### `Colgroup`

```go
func Colgroup(content ...HTMLElement) *Element
```

Colgroup creates a `<colgroup>` element.

#### `ColorInput`

```go
func ColorInput(name string) *Element
```

ColorInput creates an `<input type="color">` field.

#### `Datalist`

```go
func Datalist(id string, content ...HTMLElement) *Element
```

Datalist creates a `<datalist>` element.

#### `DateInput`

```go
func DateInput(name string) *Element
```

DateInput creates an `<input type="date">` field.

#### `DateTimeLocalInput`

```go
func DateTimeLocalInput(name string) *Element
```

DateTimeLocalInput creates an `<input type="datetime-local">` field.

#### `Details`

```go
func Details(content ...HTMLElement) *Element
```

Details creates a `<details>` element.

#### `Dfn`

```go
func Dfn(content ...HTMLElement) *Element
```

Dfn creates a `<dfn>` definition element.

#### `DialogEl`

```go
func DialogEl(content ...HTMLElement) *Element
```

DialogEl creates a `<dialog>` element. Renamed to avoid potential conflicts.

#### `Div`

```go
func Div(content ...HTMLElement) *Element
```

Div creates a `<div>` element.

#### `Em`

```go
func Em(content ...HTMLElement) *Element
```

Em creates an `<em>` emphasis element.

#### `EmailInput`

```go
func EmailInput(name string) *Element
```

EmailInput creates an `<input type="email">` field.

#### `EmbedEl`

```go
func EmbedEl(src string, embedType string) *Element
```

EmbedEl creates an `<embed>` element. It's self-closing. Renamed to avoid keyword conflict.

#### `Favicon`

```go
func Favicon(href string, rel ...string) *Element
```

Favicon creates a `<link>` element for a favicon.

#### `Fieldset`

```go
func Fieldset(content ...HTMLElement) *Element
```

Fieldset creates a `<fieldset>` element.

#### `Figcaption`

```go
func Figcaption(content ...HTMLElement) *Element
```

Figcaption creates a `<figcaption>` element.

#### `Figure`

```go
func Figure(content ...HTMLElement) *Element
```

Figure creates a `<figure>` element.

#### `FileInput`

```go
func FileInput(name string) *Element
```

FileInput creates an `<input type="file">` field.

#### `Footer`

```go
func Footer(content ...HTMLElement) *Element
```

Footer creates a `<footer>` semantic element.

#### `Form`

```go
func Form(content ...HTMLElement) *Element
```

Form creates a `<form>` element.

#### `H1`

```go
func H1(content ...HTMLElement) *Element
```

H1 creates an `<h1>` heading element.

#### `H2`

```go
func H2(content ...HTMLElement) *Element
```

H2 creates an `<h2>` heading element.

#### `H3`

```go
func H3(content ...HTMLElement) *Element
```

H3 creates an `<h3>` heading element.

#### `H4`

```go
func H4(content ...HTMLElement) *Element
```

H4 creates an `<h4>` heading element.

#### `H5`

```go
func H5(content ...HTMLElement) *Element
```

H5 creates an `<h5>` heading element.

#### `H6`

```go
func H6(content ...HTMLElement) *Element
```

H6 creates an `<h6>` heading element.

#### `Head`

```go
func Head(content ...HTMLElement) *Element
```

Head creates a `<head>` element.

#### `Header`

```go
func Header(content ...HTMLElement) *Element
```

Header creates a `<header>` semantic element.

#### `HiddenInput`

```go
func HiddenInput(name string, value string) *Element
```

HiddenInput creates an `<input type="hidden">` field.

#### `Hr`

```go
func Hr() *Element
```

Hr creates a self-closing `<hr>` horizontal rule element.

#### `Html`

```go
func Html(content ...HTMLElement) *Element
```

Html creates an `<html>` element.

#### `I`

```go
func I(content ...HTMLElement) *Element
```

I creates an `<i>` idiomatic text element.

#### `Iframe`

```go
func Iframe(src string) *Element
```

Iframe creates an `<iframe>` element.

#### `Image`

```go
func Image(src, alt string) *Element
```

Image creates an `<img>` element (alias for Img).

#### `Img`

```go
func Img(src, alt string) *Element
```

Img creates a self-closing `<img>` element.

#### `InlineScript`

```go
func InlineScript(scriptContent string) *Element
```

InlineScript creates a `<script>` element with inline JavaScript content.

#### `Input`

```go
func Input(inputType string) *Element
```

Input creates a self-closing `<input>` element.

#### `Kbd`

```go
func Kbd(content ...HTMLElement) *Element
```

Kbd creates a `<kbd>` keyboard input element.

#### `Label`

```go
func Label(content ...HTMLElement) *Element
```

Label creates a `<label>` element.

#### `Legend`

```go
func Legend(content ...HTMLElement) *Element
```

Legend creates a `<legend>` element.

#### `Li`

```go
func Li(content ...HTMLElement) *Element
```

Li creates a `<li>` list item element.

#### `Link`

```go
func Link(href, textContent string) *Element
```

Link creates an `<a>` anchor element with href and text content.

#### `LinkTag`

```go
func LinkTag() *Element
```

LinkTag creates a generic `<link>` element. It's self-closing.

#### `Main`

```go
func Main(content ...HTMLElement) *Element
```

Main creates a `<main>` semantic element.

#### `Mark`

```go
func Mark(content ...HTMLElement) *Element
```

Mark creates a `<mark>` element.

#### `Meta`

```go
func Meta() *Element
```

Meta creates a generic `<meta>` element. It's self-closing.

#### `MetaCharset`

```go
func MetaCharset(charset string) *Element
```

MetaCharset creates a `<meta charset="...">` element.

#### `MetaNameContent`

```go
func MetaNameContent(name, contentVal string) *Element
```

MetaNameContent creates a `<meta name="..." content="...">` element.

#### `MetaPropertyContent`

```go
func MetaPropertyContent(property, contentVal string) *Element
```

MetaPropertyContent creates a `<meta property="..." content="...">` element.

#### `MetaViewport`

```go
func MetaViewport(contentVal string) *Element
```

MetaViewport creates a `<meta name="viewport" content="...">` element.

#### `MeterEl`

```go
func MeterEl(content ...HTMLElement) *Element
```

MeterEl creates a `<meter>` element. Renamed to avoid potential conflicts.

#### `MonthInput`

```go
func MonthInput(name string) *Element
```

MonthInput creates an `<input type="month">` field.

#### `Nav`

```go
func Nav(content ...HTMLElement) *Element
```

Nav creates a `<nav>` semantic element.

#### `NoScript`

```go
func NoScript(content ...HTMLElement) *Element
```

NoScript creates a `<noscript>` element.

#### `NumberInput`

```go
func NumberInput(name string) *Element
```

NumberInput creates an `<input type="number">` field.

#### `ObjectEl`

```go
func ObjectEl(content ...HTMLElement) *Element
```

ObjectEl creates an `<object>` element. Renamed to avoid keyword conflict.

#### `Ol`

```go
func Ol(content ...HTMLElement) *Element
```

Ol creates an `<ol>` ordered list element.

#### `Optgroup`

```go
func Optgroup(label string, content ...HTMLElement) *Element
```

Optgroup creates an `<optgroup>` element.

#### `Option`

```go
func Option(value string, content ...HTMLElement) *Element
```

Option creates an `<option>` element.

#### `OutputEl`

```go
func OutputEl(content ...HTMLElement) *Element
```

OutputEl creates an `<output>` element. Renamed to avoid potential conflicts.

#### `P`

```go
func P(content ...HTMLElement) *Element
```

P creates a `<p>` paragraph element.

#### `Param`

```go
func Param(name, value string) *Element
```

Param creates a `<param>` element. It's self-closing.

#### `PasswordInput`

```go
func PasswordInput(name string) *Element
```

PasswordInput creates an `<input type="password">` field.

#### `Pre`

```go
func Pre(content ...HTMLElement) *Element
```

Pre creates a `<pre>` element.

#### `Preload`

```go
func Preload(href string, asType string) *Element
```

Preload creates a `<link rel="preload">` element.

#### `ProgressEl`

```go
func ProgressEl(content ...HTMLElement) *Element
```

ProgressEl creates a `<progress>` element. Renamed to avoid potential conflicts.

#### `Q`

```go
func Q(content ...HTMLElement) *Element
```

Q creates a `<q>` inline quotation element.

#### `RadioInput`

```go
func RadioInput(name, value string) *Element
```

RadioInput creates an `<input type="radio">`.

#### `RangeInput`

```go
func RangeInput(name string) *Element
```

RangeInput creates an `<input type="range">` field.

#### `ResetButton`

```go
func ResetButton(text string) *Element
```

ResetButton creates a `<button type="reset">`.

#### `Samp`

```go
func Samp(content ...HTMLElement) *Element
```

Samp creates a `<samp>` sample output element.

#### `Script`

```go
func Script(src string) *Element
```

Script creates a `<script>` element for external JavaScript files.

#### `SearchInput`

```go
func SearchInput(name string) *Element
```

SearchInput creates an `<input type="search">` field.

#### `Section`

```go
func Section(content ...HTMLElement) *Element
```

Section creates a `<section>` semantic element.

#### `Select`

```go
func Select(content ...HTMLElement) *Element
```

Select creates a `<select>` dropdown element.

#### `Small`

```go
func Small(content ...HTMLElement) *Element
```

Small creates a `<small>` element.

#### `Source`

```go
func Source(src string, mediaType string) *Element
```

Source creates a `<source>` element. It's self-closing.

#### `Span`

```go
func Span(content ...HTMLElement) *Element
```

Span creates a `<span>` inline element.

#### `Strong`

```go
func Strong(content ...HTMLElement) *Element
```

Strong creates a `<strong>` element.

#### `StyleSheet`

```go
func StyleSheet(href string) *Element
```

StyleSheet creates a `<link rel="stylesheet">` element.

#### `StyleTag`

```go
func StyleTag(cssContent string) *Element
```

StyleTag creates a `<style>` element for embedding CSS.

#### `Sub`

```go
func Sub(content ...HTMLElement) *Element
```

Sub creates a `<sub>` subscript element.

#### `SubmitButton`

```go
func SubmitButton(text string) *Element
```

SubmitButton creates a `<button type="submit">`.

#### `Summary`

```go
func Summary(content ...HTMLElement) *Element
```

Summary creates a `<summary>` element.

#### `Sup`

```go
func Sup(content ...HTMLElement) *Element
```

Sup creates a `<sup>` superscript element.

#### `Table`

```go
func Table(content ...HTMLElement) *Element
```

Table creates a `<table>` element.

#### `Tbody`

```go
func Tbody(content ...HTMLElement) *Element
```

Tbody creates a `<tbody>` table body group element.

#### `Td`

```go
func Td(content ...HTMLElement) *Element
```

Td creates a `<td>` table data cell element.

#### `TelInput`

```go
func TelInput(name string) *Element
```

TelInput creates an `<input type="tel">` field.

#### `TextInput`

```go
func TextInput(name string) *Element
```

TextInput creates an `<input type="text">` field.

#### `Textarea`

```go
func Textarea(content ...HTMLElement) *Element
```

Textarea creates a `<textarea>` element.

#### `Th`

```go
func Th(content ...HTMLElement) *Element
```

Th creates a `<th>` table header cell element.

#### `Thead`

```go
func Thead(content ...HTMLElement) *Element
```

Thead creates a `<thead>` table header group element.

#### `TimeEl`

```go
func TimeEl(content ...HTMLElement) *Element
```

TimeEl creates a `<time>` element. Renamed to avoid potential conflicts.

#### `TimeInput`

```go
func TimeInput(name string) *Element
```

TimeInput creates an `<input type="time">` field.

#### `TitleEl`

```go
func TitleEl(titleText string) *Element
```

TitleEl creates a `<title>` element with the specified text.
Renamed from Title to TitleEl to avoid conflict with (*Element).Title method if it existed.

#### `Tr`

```go
func Tr(content ...HTMLElement) *Element
```

Tr creates a `<tr>` table row element.

#### `Track`

```go
func Track(kind, src, srclang string) *Element
```

Track creates a `<track>` element. It's self-closing.

#### `U`

```go
func U(content ...HTMLElement) *Element
```

U creates a `<u>` unarticulated annotation element.

#### `Ul`

```go
func Ul(content ...HTMLElement) *Element
```

Ul creates a `<ul>` unordered list element.

#### `UrlInput`

```go
func UrlInput(name string) *Element
```

UrlInput creates an `<input type="url">` field.

#### `VarEl`

```go
func VarEl(content ...HTMLElement) *Element
```

VarEl creates a `<var>` variable element. Renamed to avoid keyword conflict.

#### `Video`

```go
func Video(content ...HTMLElement) *Element
```

Video creates a `<video>` element.

#### `Wbr`

```go
func Wbr() *Element
```

Wbr creates a `<wbr>` word break opportunity element. It's self-closing.

#### `WeekInput`

```go
func WeekInput(name string) *Element
```

WeekInput creates an `<input type="week">` field.

#### Methods

#### `Add`

```go
func (e *Element) Add(children ...HTMLElement) *Element
```

Add appends child elements to this element.

#### `Attr`

```go
func (e *Element) Attr(key, value string) *Element
```

Attr sets an attribute on the element and returns the element for method chaining.

#### `BoolAttr`

```go
func (e *Element) BoolAttr(key string, present bool) *Element
```

BoolAttr sets or removes a boolean attribute on the element.
If present is true, the attribute is added (e.g., `<input disabled>`).
If present is false, the attribute is removed if it exists.

#### `Class`

```go
func (e *Element) Class(class string) *Element
```

Class sets the class attribute on the element.

#### `ID`

```go
func (e *Element) ID(id string) *Element
```

ID sets the id attribute on the element.

#### `Render`

```go
func (e *Element) Render() string
```

Render converts the element to its HTML string representation.
It handles both self-closing and container elements, attributes, content, and children.
The output is properly formatted HTML that can be sent to browsers.
Content and attribute values are HTML-escaped to prevent XSS, except for
specific tags like `<script>` and `<style>` whose content must be raw.

#### `Style`

```go
func (e *Element) Style(style string) *Element
```

Style sets the style attribute on the element.

#### `Text`

```go
func (e *Element) Text(text string) *Element
```

Text sets the text content of the element. This content is HTML-escaped during rendering.

### `EnforceContentTypeConfig`

```go
type EnforceContentTypeConfig struct {
	// AllowedTypes is a list of allowed Content-Type values (e.g., "application/json").
	// The check ignores parameters like "; charset=utf-8". Required.
	AllowedTypes []string
	// MethodsToCheck specifies which HTTP methods should have their Content-Type checked.
	// Defaults to POST, PUT, PATCH.
	MethodsToCheck []string
	// OnError allows custom handling when the Content-Type is missing or unsupported.
	// If nil, sends 400 Bad Request or 415 Unsupported Media Type.
	OnError func(w http.ResponseWriter, r *http.Request, err error)
}
```

EnforceContentTypeConfig holds configuration for EnforceContentTypeMiddleware.

### `Flag`

```go
type Flag interface {
	fmt.Stringer // For generating help text representation.

	// Apply binds the flag definition to a standard Go flag.FlagSet.
	// It configures the flag's name, aliases, usage, and default value.
	Apply(set *flag.FlagSet, cli *CLI) error
	// GetName returns the primary name of the flag (e.g., "config").
	GetName() string
	// IsRequired indicates whether the flag must be provided by the user.
	IsRequired() bool
	// Validate checks if a required flag was provided correctly after parsing.
	Validate(set *flag.FlagSet) error
	// GetAliases returns the alternative names for the flag.
	GetAliases() []string
}
```

Flag defines the interface for command-line flags.
Concrete types like StringFlag, IntFlag, BoolFlag implement this interface.

### `Float64Flag`

```go
type Float64Flag struct {
	// Name is the primary identifier for the flag (e.g., "rate"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "r"). Single letters are used as -r.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default float64
	// Required indicates if the user must provide this flag.
	Required bool
}
```

Float64Flag defines a flag that accepts a float64 value.

#### Methods

#### `Apply`

```go
func (f *Float64Flag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the float64 flag with the flag.FlagSet.

#### `GetAliases`

```go
func (f *Float64Flag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

#### `GetName`

```go
func (f *Float64Flag) GetName() string
```

GetName returns the primary name of the flag.

#### `IsRequired`

```go
func (f *Float64Flag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

#### `String`

```go
func (f *Float64Flag) String() string
```

String returns the help text representation of the flag.

#### `Validate`

```go
func (f *Float64Flag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided.

### `ForceHTTPSConfig`

```go
type ForceHTTPSConfig struct {
	// TargetHost overrides the host used in the redirect URL. If empty, uses r.Host.
	TargetHost string
	// TargetPort overrides the port used in the redirect URL. If 0, omits the port (standard 443).
	TargetPort int
	// RedirectCode is the HTTP status code for redirection. Defaults to http.StatusMovedPermanently (301).
	RedirectCode int
	// ForwardedProtoHeader is the header name to check for the original protocol (e.g., "X-Forwarded-Proto").
	// Defaults to "X-Forwarded-Proto".
	ForwardedProtoHeader string
	// TrustForwardedHeader explicitly enables trusting the ForwardedProtoHeader. Defaults to true.
	// Set to false if your proxy setup doesn't reliably set this header.
	TrustForwardedHeader *bool // Use pointer for explicit false vs unset
}
```

ForceHTTPSConfig holds configuration for ForceHTTPSMiddleware.

### `Group`

```go
type Group struct {
	// prefix is the URL path prefix applied to all routes in this group.
	prefix string
	// router is the parent router that actually stores the routes.
	router *Router
	// middlewares contains middleware functions applied only to routes in this group.
	middlewares []Middleware
}
```

Group is a lightweight helper that allows users to register a set of routes
that share a common prefix and/or middleware. It delegates to the parent router
while applying its own prefix and middleware chain.

#### Methods

#### `Delete`

```go
func (g *Group) Delete(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Delete registers a new route for HTTP DELETE requests within the group.

#### `DeleteFunc`

```go
func (g *Group) DeleteFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

DeleteFunc registers a new enhanced route for HTTP DELETE requests within the group.

#### `Get`

```go
func (g *Group) Get(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Get registers a new route for HTTP GET requests within the group.

#### `GetFunc`

```go
func (g *Group) GetFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

GetFunc registers a new enhanced route for HTTP GET requests within the group.

#### `Handle`

```go
func (g *Group) Handle(method, pattern string, handler http.HandlerFunc, opts ...*RouteOptions)
```

Handle registers a new route within the group, applying the group's prefix and middleware.
The route is ultimately registered with the parent router after transformations.

#### `HandleFunc`

```go
func (g *Group) HandleFunc(method, pattern string, handler HandlerFunc, opts ...*RouteOptions)
```

HandleFunc registers a new enhanced route within the group, applying the group's
prefix and middleware. Uses the enhanced handler signature for better error handling.

#### `Patch`

```go
func (g *Group) Patch(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Patch registers a new route for HTTP PATCH requests within the group.

#### `PatchFunc`

```go
func (g *Group) PatchFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PatchFunc registers a new enhanced route for HTTP PATCH requests within the group.

#### `Post`

```go
func (g *Group) Post(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Post registers a new route for HTTP POST requests within the group.

#### `PostFunc`

```go
func (g *Group) PostFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PostFunc registers a new enhanced route for HTTP POST requests within the group.

#### `Put`

```go
func (g *Group) Put(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Put registers a new route for HTTP PUT requests within the group.

#### `PutFunc`

```go
func (g *Group) PutFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PutFunc registers a new enhanced route for HTTP PUT requests within the group.

#### `Use`

```go
func (g *Group) Use(mws ...Middleware)
```

Use adds middleware functions to the group. These middleware functions apply
only to routes registered through the group, not to the parent router.

### `GzipConfig`

```go
type GzipConfig struct {
	// CompressionLevel specifies the gzip compression level.
	// Accepts values from gzip.BestSpeed (1) to gzip.BestCompression (9).
	// Defaults to gzip.DefaultCompression (-1).
	CompressionLevel int
	// AddVaryHeader indicates whether to add the "Vary: Accept-Encoding" header.
	// A nil value defaults to true. Set explicitly to false to disable.
	// Disabling is only recommended if caching behavior is fully understood.
	AddVaryHeader *bool
	// Logger specifies an optional logger for errors.
	// Defaults to log.Default().
	Logger *log.Logger
	// Pool specifies an optional sync.Pool for gzip.Writer reuse.
	// Can improve performance by reducing allocations.
	Pool *sync.Pool // Optional: Pool for gzip.Writers
}
```

GzipConfig holds configuration options for the GzipMiddleware.

### `HTMLDocument`

```go
type HTMLDocument struct {
	rootElement *Element // The root `<html>` element.
}
```

HTMLDocument represents a full HTML document, including the DOCTYPE.
Its Render method produces the complete HTML string.

#### Associated Functions

#### `Document`

```go
func Document(config DocumentConfig, bodyContent ...HTMLElement) *HTMLDocument
```

Document creates a complete HTML5 document structure, encapsulated in an HTMLDocument.
The returned HTMLDocument's Render method will produce the full HTML string,
including the DOCTYPE.
It uses DocumentConfig to customize the document's head and html attributes,
and accepts variadic arguments for the body content.
Sensible defaults are applied for common attributes and meta tags if not specified
in the config.

#### Methods

#### `Render`

```go
func (d *HTMLDocument) Render() string
```

Render converts the HTMLDocument to its full string representation,
prepending the HTML5 DOCTYPE declaration.

### `HTMLElement`

```go
type HTMLElement interface {
	Render() string
}
```

HTMLElement represents any HTML element that can be rendered to a string.
This interface allows for composition of complex HTML structures using
both predefined elements and custom implementations.

#### Associated Functions

#### `Text`

```go
func Text(text string) HTMLElement
```

Text creates a raw text node.

### `HandlerFunc`

```go
type HandlerFunc func(*ResponseContext) error
```

HandlerFunc is an enhanced handler function that receives a ResponseContext and returns an error.
This allows for cleaner error handling and response management compared to standard http.HandlerFunc.

### `HeaderObject`

```go
type HeaderObject struct {
	Description string        `json:"description,omitempty"`
	Schema      *SchemaObject `json:"schema"`
}
```

HeaderObject describes a response header in an Operation.

### `HealthCheckConfig`

```go
type HealthCheckConfig struct {
	// Path is the URL path for the health check endpoint. Defaults to "/healthz".
	Path string
	// Handler is the handler function to execute for the health check.
	// If nil, a default handler returning 200 OK with "OK" body is used.
	// You can provide a custom handler to check database connections, etc.
	Handler http.HandlerFunc
}
```

HealthCheckConfig holds configuration for HealthCheckMiddleware.

### `IPFilterConfig`

```go
type IPFilterConfig struct {
	// AllowedIPs is a list of allowed IPs or CIDRs.
	AllowedIPs []string
	// BlockedIPs is a list of blocked IPs or CIDRs. Takes precedence over AllowedIPs.
	BlockedIPs []string
	// BlockByDefault determines the behavior if an IP matches neither list.
	// If true, the IP is blocked unless explicitly in AllowedIPs (and not in BlockedIPs).
	// If false, the IP is allowed unless explicitly in BlockedIPs. Defaults to false.
	BlockByDefault bool
	// OnForbidden allows custom handling when an IP is forbidden.
	// If nil, sends 403 Forbidden.
	OnForbidden func(w http.ResponseWriter, r *http.Request)
	// Logger for potential IP parsing errors. Defaults to log.Default().
	Logger *log.Logger
}
```

IPFilterConfig holds configuration for IP filtering.

### `Info`

```go
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}
```

Info provides metadata about the API: title, version, and optional description.

### `IntFlag`

```go
type IntFlag struct {
	// Name is the primary identifier for the flag (e.g., "port"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "p"). Single letters are used as -p.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default int
	// Required indicates if the user must provide this flag.
	Required bool
}
```

IntFlag defines a flag that accepts an integer value.

#### Methods

#### `Apply`

```go
func (f *IntFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the int flag with the flag.FlagSet.

#### `GetAliases`

```go
func (f *IntFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

#### `GetName`

```go
func (f *IntFlag) GetName() string
```

GetName returns the primary name of the flag.

#### `IsRequired`

```go
func (f *IntFlag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

#### `String`

```go
func (f *IntFlag) String() string
```

String returns the help text representation of the flag.

#### `Validate`

```go
func (f *IntFlag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided.

### `LoggingConfig`

```go
type LoggingConfig struct {
	// Logger specifies the logger instance to use. Defaults to log.Default().
	Logger *log.Logger
	// LogRequestID determines whether to include the request ID (if available)
	// in log messages. Defaults to true.
	LogRequestID bool
	// RequestIDKey is the context key used to retrieve the request ID.
	// Defaults to the package's internal requestIDKey.
	RequestIDKey contextKey
}
```

LoggingConfig holds configuration for LoggingMiddleware.

### `MaintenanceModeConfig`

```go
type MaintenanceModeConfig struct {
	// EnabledFlag is a pointer to an atomic boolean. If true, maintenance mode is active. Required.
	EnabledFlag *atomic.Bool
	// AllowedIPs is a list of IPs or CIDRs that can bypass maintenance mode.
	AllowedIPs []string
	// StatusCode is the HTTP status code returned during maintenance. Defaults to 503.
	StatusCode int
	// RetryAfterSeconds sets the Retry-After header value in seconds. Defaults to 300.
	RetryAfterSeconds int
	// Message is the response body sent during maintenance. Defaults to a standard message.
	Message string
	// Logger for potential IP parsing errors. Defaults to log.Default().
	Logger *log.Logger
}
```

MaintenanceModeConfig holds configuration for MaintenanceModeMiddleware.

### `MaxRequestBodySizeConfig`

```go
type MaxRequestBodySizeConfig struct {
	// LimitBytes is the maximum allowed size of the request body in bytes. Required.
	LimitBytes int64
	// OnError allows custom handling when the body size limit is exceeded.
	// If nil, sends 413 Request Entity Too Large.
	OnError func(w http.ResponseWriter, r *http.Request)
}
```

MaxRequestBodySizeConfig holds configuration for MaxRequestBodySizeMiddleware.

### `MediaTypeObject`

```go
type MediaTypeObject struct {
	Schema *SchemaObject `json:"schema,omitempty"`
}
```

MediaTypeObject holds the schema defining the media type for a request or response.

### `MethodOverrideConfig`

```go
type MethodOverrideConfig struct {
	// HeaderName is the name of the header checked for the method override value.
	// Defaults to "X-HTTP-Method-Override".
	HeaderName string
	// FormFieldName is the name of the form field (used for POST requests)
	// checked for the method override value if the header is not present.
	// Set to "" to disable form field checking. Defaults to "_method".
	FormFieldName string
}
```

MethodOverrideConfig holds configuration for MethodOverrideMiddleware.

### `Middleware`

```go
type Middleware func(http.Handler) http.Handler
```

Middleware defines the function signature for middleware.
A middleware is a function that wraps an http.Handler, adding extra behavior
such as logging, authentication, or request modification.

#### Associated Functions

#### `BasicAuthMiddleware`

```go
func BasicAuthMiddleware(config BasicAuthConfig) Middleware
```

BasicAuthMiddleware provides simple HTTP Basic Authentication.

#### `CORSMiddleware`

```go
func CORSMiddleware(config CORSConfig) Middleware
```

CORSMiddleware sets Cross-Origin Resource Sharing headers.

#### `CSRFMiddleware`

```go
func CSRFMiddleware(config *CSRFConfig) Middleware
```

CSRFMiddleware provides Cross-Site Request Forgery protection.
It uses the "Double Submit Cookie" pattern. A random token is generated and
set in a secure, HttpOnly cookie. For unsafe methods (POST, PUT, etc.),
the middleware expects the same token to be present in a request header
(e.g., X-CSRF-Token) or form field, sent by the frontend JavaScript.

#### `CacheControlMiddleware`

```go
func CacheControlMiddleware(config CacheControlConfig) Middleware
```

CacheControlMiddleware sets the Cache-Control header for responses.

#### `ConcurrencyLimiterMiddleware`

```go
func ConcurrencyLimiterMiddleware(config ConcurrencyLimiterConfig) Middleware
```

ConcurrencyLimiterMiddleware limits the number of concurrent requests.

#### `ETagMiddleware`

```go
func ETagMiddleware(config *ETagConfig) Middleware
```

ETagMiddleware adds ETag headers to responses and handles If-None-Match
conditional requests, potentially returning a 304 Not Modified status.
Note: This middleware buffers the entire response body in memory to calculate
the ETag hash. This may be inefficient for very large responses.

#### `EnforceContentTypeMiddleware`

```go
func EnforceContentTypeMiddleware(config EnforceContentTypeConfig) Middleware
```

EnforceContentTypeMiddleware checks if the request's Content-Type header is allowed.

#### `ForceHTTPSMiddleware`

```go
func ForceHTTPSMiddleware(config ForceHTTPSConfig) Middleware
```

ForceHTTPSMiddleware redirects HTTP requests to HTTPS.

#### `GzipMiddleware`

```go
func GzipMiddleware(config *GzipConfig) Middleware
```

GzipMiddleware returns middleware that compresses response bodies using gzip
if the client indicates support via the Accept-Encoding header.

#### `HealthCheckMiddleware`

```go
func HealthCheckMiddleware(config *HealthCheckConfig) Middleware
```

HealthCheckMiddleware provides a simple health check endpoint.

#### `IPFilterMiddleware`

```go
func IPFilterMiddleware(config IPFilterConfig) Middleware
```

IPFilterMiddleware restricts access based on client IP address.

#### `LoggingMiddleware`

```go
func LoggingMiddleware(config *LoggingConfig) Middleware
```

LoggingMiddleware logs request details including method, path, status, size, and duration.

#### `MaintenanceModeMiddleware`

```go
func MaintenanceModeMiddleware(config MaintenanceModeConfig) Middleware
```

MaintenanceModeMiddleware returns a 503 Service Unavailable if enabled, allowing bypass for specific IPs.

#### `MaxRequestBodySizeMiddleware`

```go
func MaxRequestBodySizeMiddleware(config MaxRequestBodySizeConfig) Middleware
```

MaxRequestBodySizeMiddleware limits the size of incoming request bodies.

#### `MethodOverrideMiddleware`

```go
func MethodOverrideMiddleware(config *MethodOverrideConfig) Middleware
```

MethodOverrideMiddleware checks a header or form field to override the request method.

#### `RateLimitMiddleware`

```go
func RateLimitMiddleware(config RateLimiterConfig) Middleware
```

RateLimitMiddleware provides basic in-memory rate limiting.
Warning: This simple implementation has limitations:
- Memory usage can grow indefinitely without CleanupInterval set.
- Not suitable for distributed systems (limit is per instance).
- Accuracy decreases slightly under very high concurrency due to locking.

#### `RealIPMiddleware`

```go
func RealIPMiddleware(config RealIPConfig) Middleware
```

RealIPMiddleware extracts the client's real IP address from proxy headers.
Warning: Only use this if you have a trusted proxy setting these headers correctly.

#### `RecoveryMiddleware`

```go
func RecoveryMiddleware(config *RecoveryConfig) Middleware
```

RecoveryMiddleware recovers from panics in downstream handlers.

#### `RequestIDMiddleware`

```go
func RequestIDMiddleware(config *RequestIDConfig) Middleware
```

RequestIDMiddleware retrieves a request ID from a header or generates one.
It sets the ID in the response header and request context.

#### `SecurityHeadersMiddleware`

```go
func SecurityHeadersMiddleware(config SecurityHeadersConfig) Middleware
```

SecurityHeadersMiddleware sets common security headers.

#### `TimeoutMiddleware`

```go
func TimeoutMiddleware(config TimeoutConfig) Middleware
```

TimeoutMiddleware sets a maximum duration for handling requests.

#### `TrailingSlashRedirectMiddleware`

```go
func TrailingSlashRedirectMiddleware(config TrailingSlashRedirectConfig) Middleware
```

TrailingSlashRedirectMiddleware redirects requests to add or remove a trailing slash.

### `OpenAPI`

```go
type OpenAPI struct {
	OpenAPI    string               `json:"openapi"`
	Info       Info                 `json:"info"`
	Paths      map[string]*PathItem `json:"paths"`
	Components *Components          `json:"components,omitempty"`
}
```

OpenAPI is the root document object for an OpenAPI 3 specification.

#### Associated Functions

#### `GenerateOpenAPISpec`

```go
func GenerateOpenAPISpec(router *Router, config OpenAPIConfig) *OpenAPI
```

GenerateOpenAPISpec constructs an OpenAPI 3.0 specification from the given
router and configuration, including paths, operations, and components.

### `OpenAPIConfig`

```go
type OpenAPIConfig struct {
	Title       string
	Version     string
	Description string
}
```

OpenAPIConfig holds metadata for generating an OpenAPI specification.

### `Operation`

```go
type Operation struct {
	Tags        []string                   `json:"tags,omitempty"`
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	OperationID string                     `json:"operationId,omitempty"`
	Parameters  []ParameterObject          `json:"parameters,omitempty"`
	RequestBody *RequestBodyObject         `json:"requestBody,omitempty"`
	Responses   map[string]*ResponseObject `json:"responses"`
	Deprecated  bool                       `json:"deprecated,omitempty"`
}
```

Operation describes a single API operation on a path.

#### Associated Functions

#### `buildOperation`

```go
func buildOperation(route route, schemaCtx *schemaGenCtx) *Operation
```

buildOperation creates an OpenAPI OperationObject from a route and schema
generation context, adding parameters, requestBody, and responses.

### `ParameterObject`

```go
type ParameterObject struct {
	Name        string        `json:"name"`
	In          string        `json:"in"`
	Description string        `json:"description,omitempty"`
	Required    bool          `json:"required"`
	Deprecated  bool          `json:"deprecated,omitempty"`
	Schema      *SchemaObject `json:"schema"`
	Example     any           `json:"example,omitempty"`
}
```

ParameterObject describes a single parameter for an Operation or PathItem.

### `ParameterOption`

```go
type ParameterOption struct {
	Name        string
	In          string
	Description string
	Required    bool
	Schema      any
	Example     any
}
```

ParameterOption configures an Operation parameter.
Name and In ("query", "header", "path", "cookie") are required.
Required, Schema and Example further customize the generated parameter.

### `PathItem`

```go
type PathItem struct {
	Get         *Operation        `json:"get,omitempty"`
	Post        *Operation        `json:"post,omitempty"`
	Put         *Operation        `json:"put,omitempty"`
	Delete      *Operation        `json:"delete,omitempty"`
	Patch       *Operation        `json:"patch,omitempty"`
	Parameters  []ParameterObject `json:"parameters,omitempty"`
	Summary     string            `json:"summary,omitempty"`
	Description string            `json:"description,omitempty"`
}
```

PathItem describes the operations available on a single API path.
One of Get, Post, Put, Delete, Patch may be non-nil.

### `RateLimiterConfig`

```go
type RateLimiterConfig struct {
	// Requests is the maximum number of requests allowed within the Duration. Required.
	Requests int
	// Duration specifies the time window for the request limit. Required.
	Duration time.Duration
	// Burst allows temporary bursts exceeding the rate limit, up to this many requests.
	// Defaults to the value of Requests (no extra burst capacity).
	Burst int
	// KeyFunc extracts a unique key from the request to identify the client.
	// Defaults to using the client's IP address (r.RemoteAddr).
	KeyFunc func(r *http.Request) string
	// OnLimitExceeded allows custom handling when the rate limit is hit.
	// If nil, sends 429 Too Many Requests.
	OnLimitExceeded func(w http.ResponseWriter, r *http.Request)
	// CleanupInterval specifies how often to scan and remove old entries from memory.
	// If zero or negative, no automatic cleanup occurs (potential memory leak).
	// A value like 10*time.Minute is reasonable.
	CleanupInterval time.Duration
	// Logger for potential errors during key extraction or IP parsing. Defaults to log.Default().
	Logger *log.Logger
}
```

RateLimiterConfig holds configuration for the simple rate limiter.

### `RealIPConfig`

```go
type RealIPConfig struct {
	// TrustedProxyCIDRs is a list of CIDR notations for trusted proxies.
	// If the direct connection (r.RemoteAddr) is from one of these, proxy headers are trusted.
	TrustedProxyCIDRs []string
	// IPHeaders is an ordered list of header names to check for the client's IP.
	// The first non-empty, valid IP found is used.
	// Defaults to ["X-Forwarded-For", "X-Real-IP"].
	IPHeaders []string
	// StoreInContext determines whether to store the found real IP in the request context.
	// Defaults to true.
	StoreInContext bool
	// ContextKey is the key used if StoreInContext is true. Defaults to realIPKey.
	ContextKey contextKey
}
```

RealIPConfig holds configuration for RealIPMiddleware.

### `RecoveryConfig`

```go
type RecoveryConfig struct {
	// Logger specifies the logger instance for panic messages. Defaults to log.Default().
	Logger *log.Logger
	// LogRequestID determines whether to include the request ID (if available)
	// in log messages. Defaults to true.
	LogRequestID bool
	// RequestIDKey is the context key used to retrieve the request ID.
	// Defaults to the package's internal requestIDKey.
	RequestIDKey contextKey
	// RecoveryHandler allows custom logic to run after a panic is recovered.
	// It receives the response writer, request, and the recovered panic value.
	// If nil, a default 500 Internal Server Error response is sent.
	RecoveryHandler func(http.ResponseWriter, *http.Request, any)
}
```

RecoveryConfig holds configuration for RecoveryMiddleware.

### `RequestBodyObject`

```go
type RequestBodyObject struct {
	Description string                      `json:"description,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content"`
	Required    bool                        `json:"required,omitempty"`
}
```

RequestBodyObject describes a request body for an Operation.

### `RequestIDConfig`

```go
type RequestIDConfig struct {
	// HeaderName is the name of the HTTP header to check for an existing ID
	// and to set in the response. Defaults to "X-Request-ID".
	HeaderName string
	// ContextKey is the key used to store the request ID in the request context.
	// Defaults to the package's internal requestIDKey.
	ContextKey contextKey
	// Generator is a function that generates a new request ID if one is not
	// found in the header. Defaults to a nanosecond timestamp-based generator.
	// Consider using a UUID library if external dependencies are acceptable.
	Generator func() string
}
```

RequestIDConfig holds configuration for RequestIDMiddleware.

### `ResponseContext`

```go
type ResponseContext struct {
	// w is the underlying HTTP response writer for sending responses to the client.
	w http.ResponseWriter
	// r is the incoming HTTP request containing headers, body, and URL parameters.
	r *http.Request
	// router is the parent router instance used for URL parameter extraction.
	router *Router
}
```

ResponseContext provides helper methods for sending HTTP responses with reduced boilerplate.
It wraps the http.ResponseWriter and http.Request along with router context to provide
convenient methods for JSON, HTML, text responses, and automatic data binding.

#### Methods

#### `Bind`

```go
func (rc *ResponseContext) Bind(v any) error
```

Bind automatically unmarshals and validates request data into the provided struct.
It supports both JSON and form data, with automatic content-type detection.
The struct should be a pointer. Validation is performed if validation middleware is active.

#### `BindForm`

```go
func (rc *ResponseContext) BindForm(v any) error
```

BindForm binds form data into the provided struct using reflection.
The struct should be a pointer. Field names are matched using JSON tags or struct field names.
Supports string, bool, int, and float fields with automatic type conversion.

#### `BindJSON`

```go
func (rc *ResponseContext) BindJSON(v any) error
```

BindJSON unmarshals JSON request body into the provided struct.
The struct should be a pointer. Returns an error if the body is nil or JSON is invalid.

#### `BindValidated`

```go
func (rc *ResponseContext) BindValidated(v any) error
```

BindValidated binds and validates request data (JSON or form) with comprehensive validation.

#### `HTML`

```go
func (rc *ResponseContext) HTML(statusCode int, content HTMLElement) error
```

HTML sends an HTML response with the given status code and content.
It automatically sets the Content-Type header to "text/html; charset=utf-8"
and renders the provided HTMLElement to string.

#### `JSON`

```go
func (rc *ResponseContext) JSON(statusCode int, data any) error
```

JSON sends a JSON response with the given status code and data.
It automatically sets the Content-Type header to "application/json" and
handles JSON encoding. Returns an error if encoding fails.

#### `JSONError`

```go
func (rc *ResponseContext) JSONError(statusCode int, message string) error
```

JSONError sends a JSON error response with the given status code and message.
It creates a standardized error response format with an "error" field.

#### `Redirect`

```go
func (rc *ResponseContext) Redirect(statusCode int, url string) error
```

Redirect sends an HTTP redirect response with the given status code and URL.
Common status codes are 301 (permanent), 302 (temporary), and 307 (temporary, preserve method).

#### `Request`

```go
func (rc *ResponseContext) Request() *http.Request
```

Request returns the underlying http.Request for advanced request handling
when the ResponseContext helpers are not sufficient.

#### `Text`

```go
func (rc *ResponseContext) Text(statusCode int, text string) error
```

Text sends a plain text response with the given status code and text content.
It automatically sets the Content-Type header to "text/plain; charset=utf-8".

#### `URLParam`

```go
func (rc *ResponseContext) URLParam(key string) string
```

URLParam retrieves the URL parameter for the given key from the request context.
Returns an empty string if the parameter is not present or if the request
doesn't contain URL parameters.

#### `WantsJSON`

```go
func (rc *ResponseContext) WantsJSON() bool
```

WantsJSON returns true if the request expects a JSON response based on
Content-Type or Accept headers. Useful for dual HTML/JSON endpoints.

#### `Writer`

```go
func (rc *ResponseContext) Writer() http.ResponseWriter
```

Writer returns the underlying http.ResponseWriter for advanced response handling
when the ResponseContext helpers are not sufficient.

### `ResponseObject`

```go
type ResponseObject struct {
	Description string                      `json:"description"`
	Headers     map[string]*HeaderObject    `json:"headers,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content,omitempty"`
}
```

ResponseObject describes a single response in an Operation.

### `ResponseOption`

```go
type ResponseOption struct {
	Description string
	Body        any
}
```

ResponseOption configures a single HTTP response in an Operation.
Description is required; Body, if non-nil, is used to generate a schema.

### `RouteOptions`

```go
type RouteOptions struct {
	Tags        []string
	Summary     string
	Description string
	OperationID string
	Deprecated  bool
	RequestBody any
	Responses   map[int]ResponseOption
	Parameters  []ParameterOption
}
```

RouteOptions holds OpenAPI metadata for a single route.
Tags, Summary, Description, and OperationID map directly into the
corresponding Operation fields. RequestBody, Responses, and Parameters
drive schema generation for request bodies, responses, and parameters.

### `Router`

```go
type Router struct {
	// routes contains all registered routes with their patterns and handlers.
	routes []route
	// subrouters contains mounted child routers with their own base paths.
	subrouters []*Router
	// paramsKey is the context key used to store URL parameters in request context.
	paramsKey any
	// middlewares contains all registered middleware functions applied to routes.
	middlewares []Middleware
	// chain is the composed middleware chain function applied to all routes.
	chain func(http.Handler) http.Handler
	// basePath is the prefix path for this router, used in subrouters and groups.
	basePath string
	// notFoundHandler is the custom handler for 404 Not Found responses.
	notFoundHandler http.Handler
	// methodNotAllowedHandler is the custom handler for 405 Method Not Allowed responses.
	methodNotAllowedHandler http.Handler
}
```

Router is a minimal HTTP router that supports dynamic routes with regex validation
in path parameters and can mount subrouters. It provides middleware support,
custom error handlers, and both traditional and enhanced handler functions.

#### Associated Functions

#### `NewRouter`

```go
func NewRouter() *Router
```

NewRouter creates and returns a new Router instance with default configuration.

#### Methods

#### `Delete`

```go
func (r *Router) Delete(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Delete registers a new route for HTTP DELETE requests using the standard handler signature.

#### `DeleteFunc`

```go
func (r *Router) DeleteFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

DeleteFunc registers a new route for HTTP DELETE requests using the enhanced handler signature.

#### `Get`

```go
func (r *Router) Get(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Get registers a new route for HTTP GET requests using the standard handler signature.

#### `GetFunc`

```go
func (r *Router) GetFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

GetFunc registers a new route for HTTP GET requests using the enhanced handler signature.

#### `Group`

```go
func (r *Router) Group(prefix string, mws ...Middleware) *Group
```

Group creates and returns a new Group with the given prefix. A group is a lightweight
convenience wrapper that prefixes routes and can add its own middleware without
creating a separate router instance.

#### `Handle`

```go
func (r *Router) Handle(method, pattern string, handler http.HandlerFunc, opts ...*RouteOptions)
```

Handle registers a new route with the given HTTP method, URL pattern, and handler.
If the router has a non-empty basePath, it is automatically prepended to the pattern.
Optional RouteOptions can be provided for OpenAPI documentation.

#### `HandleFunc`

```go
func (r *Router) HandleFunc(method, pattern string, handler HandlerFunc, opts ...*RouteOptions)
```

HandleFunc registers a new enhanced route that receives a ResponseContext instead
of separate ResponseWriter and Request parameters. This enables cleaner error handling
and response management with automatic data binding capabilities.

#### `Patch`

```go
func (r *Router) Patch(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Patch registers a new route for HTTP PATCH requests using the standard handler signature.

#### `PatchFunc`

```go
func (r *Router) PatchFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PatchFunc registers a new route for HTTP PATCH requests using the enhanced handler signature.

#### `Post`

```go
func (r *Router) Post(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Post registers a new route for HTTP POST requests using the standard handler signature.

#### `PostFunc`

```go
func (r *Router) PostFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PostFunc registers a new route for HTTP POST requests using the enhanced handler signature.

#### `Put`

```go
func (r *Router) Put(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Put registers a new route for HTTP PUT requests using the standard handler signature.

#### `PutFunc`

```go
func (r *Router) PutFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PutFunc registers a new route for HTTP PUT requests using the enhanced handler signature.

#### `ServeHTTP`

```go
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

ServeHTTP implements the http.Handler interface.
It first checks subrouters based on their base path, then its own routes.
If a URL pattern matches but the HTTP method does not, a 405 Method Not Allowed is returned.
If no routes match, a 404 Not Found is returned.

#### `ServeOpenAPISpec`

```go
func (r *Router) ServeOpenAPISpec(path string, config OpenAPIConfig)
```

ServeOpenAPISpec makes your OpenAPI specification available at `path` (e.g. "/openapi.json")

#### `ServeSwaggerUI`

```go
func (r *Router) ServeSwaggerUI(prefix string)
```

ServeSwaggerUI shows swaggerUI at `prefix` (e.g. "/docs")

#### `SetMethodNotAllowedHandler`

```go
func (r *Router) SetMethodNotAllowedHandler(h http.Handler)
```

SetMethodNotAllowedHandler allows users to supply a custom http.Handler
for requests that match a route pattern but use an unsupported HTTP method.
If not set, a default 405 Method Not Allowed response is returned.

#### `SetNotFoundHandler`

```go
func (r *Router) SetNotFoundHandler(h http.Handler)
```

SetNotFoundHandler allows users to supply a custom http.Handler
for requests that do not match any route. If not set, the standard
http.NotFound handler is used.

#### `Static`

```go
func (r *Router) Static(urlPathPrefix string, subFS fs.FS)
```

Static mounts an fs.FS under the given URL prefix.

#### `Subrouter`

```go
func (r *Router) Subrouter(prefix string) *Router
```

Subrouter creates a new Router mounted at the given prefix.
The subrouter inherits the parent's context key, global middleware,
and error handlers. Its routes will automatically receive the combined prefix.

#### `URLParam`

```go
func (r *Router) URLParam(req *http.Request, key string) string
```

URLParam retrieves the URL parameter for the given key from the request context.
It returns an empty string if the parameter is not present. This method should
only be called within route handlers that were registered with this router.

#### `Use`

```go
func (r *Router) Use(mws ...Middleware)
```

Use registers one or more middleware functions that will be applied globally
to every matched route handler. Middleware is applied in the order it is registered.

#### `rebuildChain`

```go
func (r *Router) rebuildChain()
```

rebuildChain reconstructs the composed middleware chain based on the
currently registered middlewares. Middleware is applied in reverse order
so that the first registered middleware wraps the outermost layer.

### `SchemaObject`

```go
type SchemaObject struct {
	Ref                  string                   `json:"$ref,omitempty"`
	Type                 string                   `json:"type,omitempty"`
	Format               string                   `json:"format,omitempty"`
	Properties           map[string]*SchemaObject `json:"properties,omitempty"`
	Items                *SchemaObject            `json:"items,omitempty"`
	Description          string                   `json:"description,omitempty"`
	Example              any                      `json:"example,omitempty"`
	Required             []string                 `json:"required,omitempty"`
	Nullable             bool                     `json:"nullable,omitempty"`
	Enum                 any                      `json:"enum,omitempty"`
	AdditionalProperties *SchemaObject            `json:"additionalProperties,omitempty"`
}
```

SchemaObject represents an OpenAPI Schema or a reference to one.

#### Associated Functions

#### `generateSchema`

```go
func generateSchema(instance any, ctx *schemaGenCtx) *SchemaObject
```

generateSchema inspects an example instance via reflection and adds or reuses
a SchemaObject in the context for components. Supports structs, arrays,
maps, and basic Go types.

### `SecurityHeadersConfig`

```go
type SecurityHeadersConfig struct {
	// ContentTypeOptions sets the X-Content-Type-Options header.
	// Defaults to "nosniff". Set to "" to disable.
	ContentTypeOptions string
	// FrameOptions sets the X-Frame-Options header.
	// Defaults to "DENY". Other common values: "SAMEORIGIN". Set to "" to disable.
	FrameOptions string
	// ReferrerPolicy sets the Referrer-Policy header.
	// Defaults to "strict-origin-when-cross-origin". Set to "" to disable.
	ReferrerPolicy string
	// HSTSMaxAgeSeconds sets the max-age for Strict-Transport-Security (HSTS).
	// If > 0, HSTS header is set for HTTPS requests. Defaults to 0 (disabled).
	HSTSMaxAgeSeconds int
	// HSTSIncludeSubdomains adds `includeSubDomains` to HSTS header. Defaults to true if HSTSMaxAgeSeconds > 0.
	HSTSIncludeSubdomains *bool // Use pointer for explicit false vs unset
	// HSTSPreload adds `preload` to HSTS header. Use with caution. Defaults to false.
	HSTSPreload bool
	// ContentSecurityPolicy sets the Content-Security-Policy header.
	// Defaults to "". This policy is complex and highly site-specific.
	ContentSecurityPolicy string
	// PermissionsPolicy sets the Permissions-Policy header (formerly Feature-Policy).
	// Defaults to "". Example: "geolocation=(), microphone=()"
	PermissionsPolicy string
}
```

SecurityHeadersConfig holds configuration for SecurityHeadersMiddleware.

### `StringFlag`

```go
type StringFlag struct {
	// Name is the primary identifier for the flag (e.g., "output"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "o"). Single letters are used as -o.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default string
	// Required indicates if the user must provide this flag.
	Required bool
}
```

StringFlag defines a flag that accepts a string value.

#### Methods

#### `Apply`

```go
func (f *StringFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the string flag with the flag.FlagSet.

#### `GetAliases`

```go
func (f *StringFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

#### `GetName`

```go
func (f *StringFlag) GetName() string
```

GetName returns the primary name of the flag.

#### `IsRequired`

```go
func (f *StringFlag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

#### `String`

```go
func (f *StringFlag) String() string
```

String returns the help text representation of the flag.

#### `Validate`

```go
func (f *StringFlag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided and is not empty.

### `StringSliceFlag`

```go
type StringSliceFlag struct {
	// Name is the primary identifier for the flag (e.g., "tag"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "t"). Single letters are used as -t.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default []string
	// Required indicates if the user must provide this flag at least once.
	Required bool
}
```

StringSliceFlag defines a flag that accepts multiple string values.
The flag can be repeated on the command line (e.g., --tag foo --tag bar).

#### Methods

#### `Apply`

```go
func (f *StringSliceFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the string slice flag with the flag.FlagSet using a custom Value.

#### `GetAliases`

```go
func (f *StringSliceFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

#### `GetName`

```go
func (f *StringSliceFlag) GetName() string
```

GetName returns the primary name of the flag.

#### `IsRequired`

```go
func (f *StringSliceFlag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

#### `String`

```go
func (f *StringSliceFlag) String() string
```

String returns the help text representation of the flag.

#### `Validate`

```go
func (f *StringSliceFlag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided and resulted in a non-empty slice.

### `TimeoutConfig`

```go
type TimeoutConfig struct {
	// Duration is the maximum time allowed for the handler to process the request.
	Duration time.Duration
	// TimeoutMessage is the message sent in the response body on timeout.
	// Defaults to "Service timed out".
	TimeoutMessage string
	// TimeoutHandler allows custom logic to run on timeout. If nil, the default
	// http.TimeoutHandler behavior (503 Service Unavailable with message) is used.
	TimeoutHandler http.Handler
}
```

TimeoutConfig holds configuration for TimeoutMiddleware.

### `TrailingSlashRedirectConfig`

```go
type TrailingSlashRedirectConfig struct {
	// AddSlash enforces trailing slashes (redirects /path to /path/). Defaults to false (removes slash).
	AddSlash bool
	// RedirectCode is the HTTP status code used for redirection.
	// Defaults to http.StatusMovedPermanently (301). Use 308 for POST/PUT etc.
	RedirectCode int
}
```

TrailingSlashRedirectConfig holds configuration for TrailingSlashRedirectMiddleware.

### `ValidationErrors`

```go
type ValidationErrors []error
```

ValidationErrors collects one or more validation violations and implements error.

#### Associated Functions

#### `validateField`

```go
func validateField(
	fv reflect.Value,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors
```

validateField routes to the appropriate validator based on kind,
returning all errors for that field.

#### `validateNumericField`

```go
func validateNumericField(
	num float64,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors
```

validateNumericField applies min, max, and multipleOf checks on num,
returning every violation.

#### `validateSliceField`

```go
func validateSliceField(
	fv reflect.Value,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors
```

validateSliceField applies minItems, maxItems, uniqueItems and recurses
into elements, returning every violation.

#### `validateStringField`

```go
func validateStringField(
	s string,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors
```

validateStringField applies minlength, maxlength, pattern, enum and format
checks on s, returning every violation.

#### Methods

#### `Error`

```go
func (ve ValidationErrors) Error() string
```

Error joins all contained error messages into a single string.

### `bufferingResponseWriterInterceptor`

```go
type bufferingResponseWriterInterceptor struct {
	http.ResponseWriter
	statusCode  int
	size        int64
	wroteHeader bool
	buffer      *bytes.Buffer // Buffer to hold the response body
}
```

bufferingResponseWriterInterceptor wraps http.ResponseWriter to capture status, size,
and buffer the response body.

#### Associated Functions

#### `NewBufferingResponseWriterInterceptor`

```go
func NewBufferingResponseWriterInterceptor(w http.ResponseWriter) *bufferingResponseWriterInterceptor
```

NewBufferingResponseWriterInterceptor creates a new buffering interceptor.

#### Methods

#### `Body`

```go
func (w *bufferingResponseWriterInterceptor) Body() []byte
```

Body returns the buffered response body bytes.

#### `Flush`

```go
func (w *bufferingResponseWriterInterceptor) Flush()
```

Flush implements the http.Flusher interface if the underlying writer supports it.

#### `Size`

```go
func (w *bufferingResponseWriterInterceptor) Size() int64
```

Size returns the captured response size.

#### `StatusCode`

```go
func (w *bufferingResponseWriterInterceptor) StatusCode() int
```

StatusCode returns the captured status code.

#### `Write`

```go
func (w *bufferingResponseWriterInterceptor) Write(b []byte) (int, error)
```

Write captures the number of bytes written and writes to the internal buffer.
It also sets the status code implicitly if not already set.

#### `WriteCapturedData`

```go
func (w *bufferingResponseWriterInterceptor) WriteCapturedData() (int64, error)
```

WriteCapturedData writes the captured status code, headers, and buffered body
to the original ResponseWriter. Returns the number of bytes written from the body.

#### `WriteHeader`

```go
func (w *bufferingResponseWriterInterceptor) WriteHeader(statusCode int)
```

WriteHeader captures the status code. Does NOT write to underlying writer yet.

### `contextKey`

```go
type contextKey string
```

contextKey is a private type used for context keys to avoid collisions.

#### Associated Constants

#### `requestIDKey, realIPKey, basicAuthUserKey, csrfTokenKey`

```go
const (
	// requestIDKey is the context key used for storing the request ID.
	requestIDKey contextKey = "requestID"
	// realIPKey is the context key used for storing the real client IP.
	realIPKey contextKey = "realIP"
	// basicAuthUserKey is the context key used for storing the authenticated user.
	basicAuthUserKey contextKey = "basicAuthUser"
	// csrfTokenKey is the context key used for storing the generated CSRF token.
	csrfTokenKey contextKey = "csrfToken"
)
```

### `gzipResponseWriter`

```go
type gzipResponseWriter struct {
	http.ResponseWriter              // Underlying writer
	Writer              *gzip.Writer // gzip writer
}
```

gzipResponseWriter wraps http.ResponseWriter to provide gzip compression.
It implements http.ResponseWriter and optionally http.Flusher.

#### Methods

#### `Flush`

```go
func (w *gzipResponseWriter) Flush()
```

Flush implements http.Flusher if the underlying ResponseWriter supports it.

#### `Header`

```go
func (w *gzipResponseWriter) Header() http.Header
```

Header returns the header map of the underlying ResponseWriter.

#### `Write`

```go
func (w *gzipResponseWriter) Write(b []byte) (int, error)
```

Write writes the data to the gzip writer.

#### `WriteHeader`

```go
func (w *gzipResponseWriter) WriteHeader(statusCode int)
```

WriteHeader sends the HTTP status code using the underlying ResponseWriter.

### `responseWriterInterceptor`

```go
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode  int
	size        int64
	wroteHeader bool // Track if WriteHeader was called explicitly
}
```

responseWriterInterceptor wraps http.ResponseWriter to capture status code and size.

#### Associated Functions

#### `NewResponseWriterInterceptor`

```go
func NewResponseWriterInterceptor(w http.ResponseWriter) *responseWriterInterceptor
```

NewResponseWriterInterceptor creates a new interceptor.

#### Methods

#### `Flush`

```go
func (w *responseWriterInterceptor) Flush()
```

Flush implements the http.Flusher interface if the underlying writer supports it.

#### `Write`

```go
func (w *responseWriterInterceptor) Write(b []byte) (int, error)
```

Write captures the number of bytes written and calls the underlying Write.
It also calls WriteHeader implicitly if not already called.

#### `WriteHeader`

```go
func (w *responseWriterInterceptor) WriteHeader(statusCode int)
```

WriteHeader captures the status code.

### `route`

```go
type route struct {
	// method is the HTTP method (GET, POST, PUT, etc.) for this route.
	method string
	// handler is the HTTP handler function that processes requests to this route.
	handler http.HandlerFunc
	// segments contains the compiled URL pattern segments for path matching.
	segments []segment
	// options contains optional metadata for OpenAPI documentation and validation.
	options *RouteOptions
}
```

route represents an individual route with its compiled URL pattern and associated handler.
It stores the HTTP method, handler function, compiled path segments, and optional metadata.

### `schemaGenCtx`

```go
type schemaGenCtx struct {
	componentsSchemas map[string]*SchemaObject
	generatedNames    map[reflect.Type]string
}
```

schemaGenCtx tracks state during schema generation, to de-duplicate
component schemas and reuse references.

#### Associated Functions

#### `newSchemaGenCtx`

```go
func newSchemaGenCtx() *schemaGenCtx
```

newSchemaGenCtx creates and returns an empty schemaGenCtx to drive
generateSchema and collect component definitions.

### `segment`

```go
type segment struct {
	// isParam indicates whether this segment is a dynamic parameter ({name}) or literal text.
	isParam bool
	// literal contains the exact text to match when isParam is false.
	literal string
	// paramName is the parameter name when isParam is true (e.g., "id" from "{id}").
	paramName string
	// regex is the compiled regular expression for parameter validation, if specified.
	regex *regexp.Regexp
}
```

segment represents a part of the URL path. It may be a literal string or a dynamic
parameter with an optional regex pattern for validation.

#### Associated Functions

#### `compilePattern`

```go
func compilePattern(pattern string) ([]segment, error)
```

compilePattern converts a URL pattern string into a slice of segments.
Parameters are declared as {name} or {name:regex}. In the latter case,
the regex is precompiled and validated during route registration.

### `stringSliceValue`

```go
type stringSliceValue struct {
	destination *[]string // Pointer to the underlying slice
	hasBeenSet  bool      // Track if the flag was set at all (for default handling)
}
```

stringSliceValue implements the flag.Value interface for string slices.

#### Associated Functions

#### `newStringSliceValue`

```go
func newStringSliceValue(dest *[]string, defaults []string) *stringSliceValue
```

newStringSliceValue creates a new value for string slice flags.

#### Methods

#### `Get`

```go
func (s *stringSliceValue) Get() any
```

Get returns the underlying slice.

#### `Set`

```go
func (s *stringSliceValue) Set(value string) error
```

Set appends the provided value to the slice. Called by flag package for each occurrence.

#### `String`

```go
func (s *stringSliceValue) String() string
```

String returns a comma-separated representation of the slice (required by flag.Value).

### `textNode`

```go
type textNode struct {
	text string
}
```

textNode represents raw text content that should be rendered without HTML tags.
It implements HTMLElement to allow text to be mixed with HTML elements in compositions.

#### Methods

#### `Render`

```go
func (t *textNode) Render() string
```

Render returns the raw text content, HTML-escaped.

### `timeoutResponseWriter`

```go
type timeoutResponseWriter struct {
	http.ResponseWriter
	mu          sync.Mutex
	wroteHeader bool
}
```

timeoutResponseWriter wraps http.ResponseWriter to prevent panics from
superfluous WriteHeader calls after a timeout.

#### Methods

#### `Write`

```go
func (w *timeoutResponseWriter) Write(b []byte) (int, error)
```

#### `WriteHeader`

```go
func (w *timeoutResponseWriter) WriteHeader(statusCode int)
```

### `visitor`

```go
type visitor struct {
	tokens    float64   // Current number of available tokens
	lastToken time.Time // Time when tokens were last refilled
	lastSeen  time.Time // Timestamp of the last request seen (for cleanup)
}
```

visitor tracks request counts and timestamps for rate limiting.

{{ endinclude }}