# Reference

## Table of Contents

- [Variables](#variables)
  - [`ErrMissingContentType`](#errmissingcontenttype)
  - [`ErrUnsupportedContentType`](#errunsupportedcontenttype)
- [Functions](#functions)
  - [`CreateNewMigration`](#createnewmigration)
  - [`GetBasicAuthUser`](#getbasicauthuser)
  - [`GetBasicAuthUserWithKey`](#getbasicauthuserwithkey)
  - [`GetCSRFToken`](#getcsrftoken)
  - [`GetCSRFTokenWithKey`](#getcsrftokenwithkey)
  - [`GetRealIP`](#getrealip)
  - [`GetRealIPWithKey`](#getrealipwithkey)
  - [`GetRequestID`](#getrequestid)
  - [`GetRequestIDWithKey`](#getrequestidwithkey)
  - [`LoadDotenv`](#loaddotenv)
  - [`MigrateDown`](#migratedown)
  - [`MigrateUp`](#migrateup)
  - [`NewBufferingResponseWriterInterceptor`](#newbufferingresponsewriterinterceptor)
  - [`NewResponseWriterInterceptor`](#newresponsewriterinterceptor)
  - [`Serve`](#serve)
- [Types](#types)
  - [`ActionFunc`](#actionfunc)
  - [`AuthValidator`](#authvalidator)
  - [`BasicAuthConfig`](#basicauthconfig)
  - [`BoolFlag`](#boolflag)
    - [`BoolFlag.Apply`](#boolflagapply)
    - [`BoolFlag.GetAliases`](#boolflaggetaliases)
    - [`BoolFlag.GetName`](#boolflaggetname)
    - [`BoolFlag.IsRequired`](#boolflagisrequired)
    - [`BoolFlag.String`](#boolflagstring)
    - [`BoolFlag.Validate`](#boolflagvalidate)
  - [`CLI`](#cli)
    - [`NewCLI`](#newcli)
    - [`CLI.AddCommand`](#cliaddcommand)
    - [`CLI.Run`](#clirun)
    - [`CLI.ShowHelp`](#clishowhelp)
  - [`CORSConfig`](#corsconfig)
  - [`CSRFConfig`](#csrfconfig)
  - [`CacheControlConfig`](#cachecontrolconfig)
  - [`Command`](#command)
    - [`Command.ShowHelp`](#commandshowhelp)
  - [`Components`](#components)
  - [`ConcurrencyLimiterConfig`](#concurrencylimiterconfig)
  - [`Context`](#context)
    - [`Context.Args`](#contextargs)
    - [`Context.Bool`](#contextbool)
    - [`Context.Float64`](#contextfloat64)
    - [`Context.Int`](#contextint)
    - [`Context.String`](#contextstring)
    - [`Context.StringSlice`](#contextstringslice)
  - [`DocumentConfig`](#documentconfig)
  - [`ETagConfig`](#etagconfig)
  - [`Element`](#element)
    - [`A`](#a)
    - [`Abbr`](#abbr)
    - [`Address`](#address)
    - [`Article`](#article)
    - [`Aside`](#aside)
    - [`Audio`](#audio)
    - [`B`](#b)
    - [`Base`](#base)
    - [`Blockquote`](#blockquote)
    - [`Body`](#body)
    - [`Br`](#br)
    - [`Button`](#button)
    - [`ButtonInput`](#buttoninput)
    - [`Caption`](#caption)
    - [`CheckboxInput`](#checkboxinput)
    - [`Cite`](#cite)
    - [`Code`](#code)
    - [`Col`](#col)
    - [`Colgroup`](#colgroup)
    - [`ColorInput`](#colorinput)
    - [`Datalist`](#datalist)
    - [`DateInput`](#dateinput)
    - [`DateTimeLocalInput`](#datetimelocalinput)
    - [`Details`](#details)
    - [`Dfn`](#dfn)
    - [`DialogEl`](#dialogel)
    - [`Div`](#div)
    - [`Em`](#em)
    - [`EmailInput`](#emailinput)
    - [`EmbedEl`](#embedel)
    - [`Favicon`](#favicon)
    - [`Fieldset`](#fieldset)
    - [`Figcaption`](#figcaption)
    - [`Figure`](#figure)
    - [`FileInput`](#fileinput)
    - [`Footer`](#footer)
    - [`Form`](#form)
    - [`H1`](#h1)
    - [`H2`](#h2)
    - [`H3`](#h3)
    - [`H4`](#h4)
    - [`H5`](#h5)
    - [`H6`](#h6)
    - [`Head`](#head)
    - [`Header`](#header)
    - [`HiddenInput`](#hiddeninput)
    - [`Hr`](#hr)
    - [`Html`](#html)
    - [`I`](#i)
    - [`Iframe`](#iframe)
    - [`Image`](#image)
    - [`Img`](#img)
    - [`InlineScript`](#inlinescript)
    - [`Input`](#input)
    - [`Kbd`](#kbd)
    - [`Label`](#label)
    - [`Legend`](#legend)
    - [`Li`](#li)
    - [`Link`](#link)
    - [`LinkTag`](#linktag)
    - [`Main`](#main)
    - [`Mark`](#mark)
    - [`Meta`](#meta)
    - [`MetaCharset`](#metacharset)
    - [`MetaNameContent`](#metanamecontent)
    - [`MetaPropertyContent`](#metapropertycontent)
    - [`MetaViewport`](#metaviewport)
    - [`MeterEl`](#meterel)
    - [`MonthInput`](#monthinput)
    - [`Nav`](#nav)
    - [`NoScript`](#noscript)
    - [`NumberInput`](#numberinput)
    - [`ObjectEl`](#objectel)
    - [`Ol`](#ol)
    - [`Optgroup`](#optgroup)
    - [`Option`](#option)
    - [`OutputEl`](#outputel)
    - [`P`](#p)
    - [`Param`](#param)
    - [`PasswordInput`](#passwordinput)
    - [`Pre`](#pre)
    - [`Preload`](#preload)
    - [`ProgressEl`](#progressel)
    - [`Q`](#q)
    - [`RadioInput`](#radioinput)
    - [`RangeInput`](#rangeinput)
    - [`ResetButton`](#resetbutton)
    - [`Samp`](#samp)
    - [`Script`](#script)
    - [`SearchInput`](#searchinput)
    - [`Section`](#section)
    - [`Select`](#select)
    - [`Small`](#small)
    - [`Source`](#source)
    - [`Span`](#span)
    - [`Strong`](#strong)
    - [`StyleSheet`](#stylesheet)
    - [`StyleTag`](#styletag)
    - [`Sub`](#sub)
    - [`SubmitButton`](#submitbutton)
    - [`Summary`](#summary)
    - [`Sup`](#sup)
    - [`Table`](#table)
    - [`Tbody`](#tbody)
    - [`Td`](#td)
    - [`TelInput`](#telinput)
    - [`TextInput`](#textinput)
    - [`Textarea`](#textarea)
    - [`Th`](#th)
    - [`Thead`](#thead)
    - [`TimeEl`](#timeel)
    - [`TimeInput`](#timeinput)
    - [`TitleEl`](#titleel)
    - [`Tr`](#tr)
    - [`Track`](#track)
    - [`U`](#u)
    - [`Ul`](#ul)
    - [`UrlInput`](#urlinput)
    - [`VarEl`](#varel)
    - [`Video`](#video)
    - [`Wbr`](#wbr)
    - [`WeekInput`](#weekinput)
    - [`Element.Add`](#elementadd)
    - [`Element.Attr`](#elementattr)
    - [`Element.BoolAttr`](#elementboolattr)
    - [`Element.Class`](#elementclass)
    - [`Element.ID`](#elementid)
    - [`Element.Render`](#elementrender)
    - [`Element.Style`](#elementstyle)
    - [`Element.Text`](#elementtext)
  - [`EnforceContentTypeConfig`](#enforcecontenttypeconfig)
  - [`Flag`](#flag)
  - [`Float64Flag`](#float64flag)
    - [`Float64Flag.Apply`](#float64flagapply)
    - [`Float64Flag.GetAliases`](#float64flaggetaliases)
    - [`Float64Flag.GetName`](#float64flaggetname)
    - [`Float64Flag.IsRequired`](#float64flagisrequired)
    - [`Float64Flag.String`](#float64flagstring)
    - [`Float64Flag.Validate`](#float64flagvalidate)
  - [`ForceHTTPSConfig`](#forcehttpsconfig)
  - [`Group`](#group)
    - [`Group.Delete`](#groupdelete)
    - [`Group.DeleteFunc`](#groupdeletefunc)
    - [`Group.Get`](#groupget)
    - [`Group.GetFunc`](#groupgetfunc)
    - [`Group.Handle`](#grouphandle)
    - [`Group.HandleFunc`](#grouphandlefunc)
    - [`Group.Patch`](#grouppatch)
    - [`Group.PatchFunc`](#grouppatchfunc)
    - [`Group.Post`](#grouppost)
    - [`Group.PostFunc`](#grouppostfunc)
    - [`Group.Put`](#groupput)
    - [`Group.PutFunc`](#groupputfunc)
    - [`Group.Use`](#groupuse)
  - [`GzipConfig`](#gzipconfig)
  - [`HTMLDocument`](#htmldocument)
    - [`Document`](#document)
    - [`HTMLDocument.Render`](#htmldocumentrender)
  - [`HTMLElement`](#htmlelement)
    - [`Text`](#text)
  - [`HandlerFunc`](#handlerfunc)
  - [`HeaderObject`](#headerobject)
  - [`HealthCheckConfig`](#healthcheckconfig)
  - [`IPFilterConfig`](#ipfilterconfig)
  - [`Info`](#info)
  - [`IntFlag`](#intflag)
    - [`IntFlag.Apply`](#intflagapply)
    - [`IntFlag.GetAliases`](#intflaggetaliases)
    - [`IntFlag.GetName`](#intflaggetname)
    - [`IntFlag.IsRequired`](#intflagisrequired)
    - [`IntFlag.String`](#intflagstring)
    - [`IntFlag.Validate`](#intflagvalidate)
  - [`LoggingConfig`](#loggingconfig)
  - [`MaintenanceModeConfig`](#maintenancemodeconfig)
  - [`MaxRequestBodySizeConfig`](#maxrequestbodysizeconfig)
  - [`MediaTypeObject`](#mediatypeobject)
  - [`MethodOverrideConfig`](#methodoverrideconfig)
  - [`Middleware`](#middleware)
    - [`BasicAuthMiddleware`](#basicauthmiddleware)
    - [`CORSMiddleware`](#corsmiddleware)
    - [`CSRFMiddleware`](#csrfmiddleware)
    - [`CacheControlMiddleware`](#cachecontrolmiddleware)
    - [`ConcurrencyLimiterMiddleware`](#concurrencylimitermiddleware)
    - [`ETagMiddleware`](#etagmiddleware)
    - [`EnforceContentTypeMiddleware`](#enforcecontenttypemiddleware)
    - [`ForceHTTPSMiddleware`](#forcehttpsmiddleware)
    - [`GzipMiddleware`](#gzipmiddleware)
    - [`HealthCheckMiddleware`](#healthcheckmiddleware)
    - [`IPFilterMiddleware`](#ipfiltermiddleware)
    - [`LoggingMiddleware`](#loggingmiddleware)
    - [`MaintenanceModeMiddleware`](#maintenancemodemiddleware)
    - [`MaxRequestBodySizeMiddleware`](#maxrequestbodysizemiddleware)
    - [`MethodOverrideMiddleware`](#methodoverridemiddleware)
    - [`RateLimitMiddleware`](#ratelimitmiddleware)
    - [`RealIPMiddleware`](#realipmiddleware)
    - [`RecoveryMiddleware`](#recoverymiddleware)
    - [`RequestIDMiddleware`](#requestidmiddleware)
    - [`SecurityHeadersMiddleware`](#securityheadersmiddleware)
    - [`TimeoutMiddleware`](#timeoutmiddleware)
    - [`TrailingSlashRedirectMiddleware`](#trailingslashredirectmiddleware)
  - [`OpenAPI`](#openapi)
    - [`GenerateOpenAPISpec`](#generateopenapispec)
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
    - [`ResponseContext.Bind`](#responsecontextbind)
    - [`ResponseContext.BindForm`](#responsecontextbindform)
    - [`ResponseContext.BindJSON`](#responsecontextbindjson)
    - [`ResponseContext.BindValidated`](#responsecontextbindvalidated)
    - [`ResponseContext.HTML`](#responsecontexthtml)
    - [`ResponseContext.JSON`](#responsecontextjson)
    - [`ResponseContext.JSONError`](#responsecontextjsonerror)
    - [`ResponseContext.Redirect`](#responsecontextredirect)
    - [`ResponseContext.Request`](#responsecontextrequest)
    - [`ResponseContext.Text`](#responsecontexttext)
    - [`ResponseContext.URLParam`](#responsecontexturlparam)
    - [`ResponseContext.WantsJSON`](#responsecontextwantsjson)
    - [`ResponseContext.Writer`](#responsecontextwriter)
  - [`ResponseObject`](#responseobject)
  - [`ResponseOption`](#responseoption)
  - [`RouteOptions`](#routeoptions)
  - [`Router`](#router)
    - [`NewRouter`](#newrouter)
    - [`Router.Delete`](#routerdelete)
    - [`Router.DeleteFunc`](#routerdeletefunc)
    - [`Router.Get`](#routerget)
    - [`Router.GetFunc`](#routergetfunc)
    - [`Router.Group`](#routergroup)
    - [`Router.Handle`](#routerhandle)
    - [`Router.HandleFunc`](#routerhandlefunc)
    - [`Router.Patch`](#routerpatch)
    - [`Router.PatchFunc`](#routerpatchfunc)
    - [`Router.Post`](#routerpost)
    - [`Router.PostFunc`](#routerpostfunc)
    - [`Router.Put`](#routerput)
    - [`Router.PutFunc`](#routerputfunc)
    - [`Router.ServeHTTP`](#routerservehttp)
    - [`Router.ServeOpenAPISpec`](#routerserveopenapispec)
    - [`Router.ServeSwaggerUI`](#routerserveswaggerui)
    - [`Router.SetMethodNotAllowedHandler`](#routersetmethodnotallowedhandler)
    - [`Router.SetNotFoundHandler`](#routersetnotfoundhandler)
    - [`Router.Static`](#routerstatic)
    - [`Router.Subrouter`](#routersubrouter)
    - [`Router.URLParam`](#routerurlparam)
    - [`Router.Use`](#routeruse)
  - [`SchemaObject`](#schemaobject)
  - [`SecurityHeadersConfig`](#securityheadersconfig)
  - [`StringFlag`](#stringflag)
    - [`StringFlag.Apply`](#stringflagapply)
    - [`StringFlag.GetAliases`](#stringflaggetaliases)
    - [`StringFlag.GetName`](#stringflaggetname)
    - [`StringFlag.IsRequired`](#stringflagisrequired)
    - [`StringFlag.String`](#stringflagstring)
    - [`StringFlag.Validate`](#stringflagvalidate)
  - [`StringSliceFlag`](#stringsliceflag)
    - [`StringSliceFlag.Apply`](#stringsliceflagapply)
    - [`StringSliceFlag.GetAliases`](#stringsliceflaggetaliases)
    - [`StringSliceFlag.GetName`](#stringsliceflaggetname)
    - [`StringSliceFlag.IsRequired`](#stringsliceflagisrequired)
    - [`StringSliceFlag.String`](#stringsliceflagstring)
    - [`StringSliceFlag.Validate`](#stringsliceflagvalidate)
  - [`TimeoutConfig`](#timeoutconfig)
  - [`TrailingSlashRedirectConfig`](#trailingslashredirectconfig)
  - [`ValidationErrors`](#validationerrors)
    - [`ValidationErrors.Error`](#validationerrorserror)


## Variables

### `ErrMissingContentType` {#errmissingcontenttype}

```go
var ErrMissingContentType = fmt.Errorf("missing Content-Type header")
```

ErrMissingContentType indicates the Content-Type header was absent.

---

### `ErrUnsupportedContentType` {#errunsupportedcontenttype}

```go
var ErrUnsupportedContentType = fmt.Errorf("unsupported Content-Type")
```

ErrUnsupportedContentType indicates the Content-Type was not in the allowed list.

---

## Functions

### `CreateNewMigration` {#createnewmigration}

```go
func CreateNewMigration(name string) error
```

CreateNewMigration creates a new migration file with a basic "up" and "down" template. The new file is named using the current Unix timestamp followed by an underscore and the provided name. The migration file is saved in the migrationsFolder and contains two sections separated by "-- migrate:up" (for the migration) and "-- migrate:down" (for the rollback).

Parameters:

```text
name - A descriptive name for the migration (e.g., "create_users_table").
```

Returns an error if the file creation fails, otherwise nil.

---

### `GetBasicAuthUser` {#getbasicauthuser}

```go
func GetBasicAuthUser(ctx context.Context) string
```

GetBasicAuthUser retrieves the authenticated username from the context, if available via BasicAuthMiddleware using the default key and if configured to store it.

---

### `GetBasicAuthUserWithKey` {#getbasicauthuserwithkey}

```go
func GetBasicAuthUserWithKey(ctx context.Context, key contextKey) string
```

GetBasicAuthUserWithKey retrieves the authenticated username from the context using a specific key.

---

### `GetCSRFToken` {#getcsrftoken}

```go
func GetCSRFToken(ctx context.Context) string
```

GetCSRFToken retrieves the CSRF token from the context, if available via CSRFMiddleware using the default key. This is the token expected in subsequent unsafe requests.

---

### `GetCSRFTokenWithKey` {#getcsrftokenwithkey}

```go
func GetCSRFTokenWithKey(ctx context.Context, key contextKey) string
```

GetCSRFTokenWithKey retrieves the CSRF token from the context using a specific key.

---

### `GetRealIP` {#getrealip}

```go
func GetRealIP(ctx context.Context) string
```

GetRealIP retrieves the client's real IP from the context, if available via RealIPMiddleware using the default key.

---

### `GetRealIPWithKey` {#getrealipwithkey}

```go
func GetRealIPWithKey(ctx context.Context, key contextKey) string
```

GetRealIPWithKey retrieves the client's real IP from the context using a specific key.

---

### `GetRequestID` {#getrequestid}

```go
func GetRequestID(ctx context.Context) string
```

GetRequestID retrieves the request ID from the context, if available via RequestIDMiddleware using the default key.

---

### `GetRequestIDWithKey` {#getrequestidwithkey}

```go
func GetRequestIDWithKey(ctx context.Context, key contextKey) string
```

GetRequestIDWithKey retrieves the request ID from the context using a specific key.

---

### `LoadDotenv` {#loaddotenv}

```go
func LoadDotenv(paths ...string) error
```

LoadDotenv loads variables from a .env file, expands them, and sets them as environment variables. If no path is provided, ".env" is used. If the specified file doesn't exist, it is silently ignored.

---

### `MigrateDown` {#migratedown}

```go
func MigrateDown(db *sql.DB, steps int) error
```

MigrateDown rolls back migrations on the provided database. It reads migration files from the migrations folder, sorts them in descending order, and applies the rollback (down) statements for each migration file where the migration version is less than or equal to the current version. The parameter steps indicates how many migrations to roll back: if steps is 0 the function rolls back one migration by default.

Parameters:

```text
db    - The database handle (from database/sql).
steps - The number of migrations to roll back (0 means 1 migration).
```

Returns an error if rollback fails, otherwise nil.

---

### `MigrateUp` {#migrateup}

```go
func MigrateUp(db *sql.DB, steps int) error
```

MigrateUp applies pending migrations to the provided database. It reads migration files from the migrations folder, sorts them in ascending order, and applies each migration with a version greater than the current database version. The parameter steps indicates how many migrations to apply: if steps is 0 the function applies all pending migrations.

Parameters:

```text
db    - The database handle (from database/sql).
steps - The maximum number of migrations to apply (0 means apply all).
```

Returns an error if migration fails, otherwise nil.

---

### `NewBufferingResponseWriterInterceptor` {#newbufferingresponsewriterinterceptor}

```go
func NewBufferingResponseWriterInterceptor(w http.ResponseWriter) *bufferingResponseWriterInterceptor
```

NewBufferingResponseWriterInterceptor creates a new buffering interceptor.

---

### `NewResponseWriterInterceptor` {#newresponsewriterinterceptor}

```go
func NewResponseWriterInterceptor(w http.ResponseWriter) *responseWriterInterceptor
```

NewResponseWriterInterceptor creates a new interceptor.

---

### `Serve` {#serve}

```go
func Serve(ctx *Context, router http.Handler) error
```

Serve launches the web server (concurrently) with graceful shutdown and live reloading (if enabled). It wraps key goroutines with recovery blocks to avoid crashes due to unexpected errors. It also allows for logging customization via context options.

---

## Types

### `ActionFunc` {#actionfunc}

```go
type ActionFunc func(ctx *Context) error
```

ActionFunc defines the function signature for CLI actions (both global and command-specific). It receives a Context containing parsed flags and arguments.

---

### `AuthValidator` {#authvalidator}

```go
type AuthValidator func(username, password string) bool
```

AuthValidator is a function type that validates the provided username and password. It returns true if the credentials are valid.

---

### `BasicAuthConfig` {#basicauthconfig}

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

---

### `BoolFlag` {#boolflag}

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

---

#### Methods

### `Apply` {#apply}

```go
func (f *BoolFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the bool flag with the flag.FlagSet.

---

### `GetAliases` {#getaliases}

```go
func (f *BoolFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

---

### `GetName` {#getname}

```go
func (f *BoolFlag) GetName() string
```

GetName returns the primary name of the flag.

---

### `IsRequired` {#isrequired}

```go
func (f *BoolFlag) IsRequired() bool
```

IsRequired always returns false for boolean flags.

---

### `String` {#string}

```go
func (f *BoolFlag) String() string
```

String returns the help text representation of the flag.

---

### `Validate` {#validate}

```go
func (f *BoolFlag) Validate(set *flag.FlagSet) error
```

Validate always returns nil for boolean flags.

---

### `CLI` {#cli}

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
	// contains filtered or unexported fields
}
```

CLI represents the command-line application structure. It holds the application's metadata, commands, global flags, and the main action. Name and Version are required fields when creating a CLI instance via NewCLI.

---

#### Associated Functions

### `NewCLI` {#newcli}

```go
func NewCLI(cli *CLI) (*CLI, error)
```

NewCLI creates and validates a new CLI application instance based on the provided configuration. The Name and Version fields in the input CLI struct are required. It checks for conflicts with reserved names/aliases (help command, h alias, version flag, v alias) and basic flag/command requirements. Returns the validated CLI instance or an error if validation fails.

---

#### Methods

### `AddCommand` {#addcommand}

```go
func (c *CLI) AddCommand(cmd *Command) error
```

AddCommand registers a new command with the application *after* initial NewCLI() validation. It performs the same conflict checks as NewCLI(). It is generally recommended to define all commands and flags within the struct passed to NewCLI().

---

### `Run` {#run}

```go
func (c *CLI) Run(arguments []string) error
```

Run executes the CLI application based on the provided command-line arguments. Call NewCLI() to create and validate the CLI instance before calling Run. Run parses flags, handles the built-in version flag and help command, validates required flags, and executes the appropriate action (global or command-specific).

---

### `ShowHelp` {#showhelp}

```go
func (c *CLI) ShowHelp(w io.Writer)
```

ShowHelp prints the main help message for the application to the specified writer.

---

### `CORSConfig` {#corsconfig}

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

---

### `CSRFConfig` {#csrfconfig}

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

---

### `CacheControlConfig` {#cachecontrolconfig}

```go
type CacheControlConfig struct {
	// CacheControlValue is the string to set for the Cache-Control header. Required.
	// Example values: "no-store", "no-cache", "public, max-age=3600"
	CacheControlValue string
}
```

CacheControlConfig holds configuration for CacheControlMiddleware.

---

### `Command` {#command}

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

Command defines a specific action the CLI application can perform. It includes metadata, flags specific to the command, and the action function. Commands named "help" or aliased "h" are reserved.

---

#### Methods

### `ShowHelp` {#showhelp}

```go
func (cmd *Command) ShowHelp(w io.Writer, appName string)
```

ShowHelp prints the help message for a specific command to the specified writer.

---

### `Components` {#components}

```go
type Components struct {
	Schemas map[string]*SchemaObject `json:"schemas,omitempty"`
}
```

Components holds reusable schema definitions for the OpenAPI spec.

---

### `ConcurrencyLimiterConfig` {#concurrencylimiterconfig}

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

---

### `Context` {#context}

```go
type Context struct {
	// CLI points to the parent application instance.
	CLI *CLI
	// Command points to the specific command being executed (nil for the global action).
	Command *Command
	// contains filtered or unexported fields
}
```

Context provides access to parsed flags, arguments, and application/command metadata within an ActionFunc.

---

#### Methods

### `Args` {#args}

```go
func (c *Context) Args() []string
```

Args returns the non-flag arguments remaining after parsing for the current context.

---

### `Bool` {#bool}

```go
func (c *Context) Bool(name string) bool
```

Bool returns the boolean value of a flag specified by name. It checks command flags first, then global flags. Returns false if not found or type mismatch.

---

### `Float64` {#float64}

```go
func (c *Context) Float64(name string) float64
```

Float64 returns the float64 value of a flag specified by name. It checks command flags first, then global flags. Returns 0.0 if not found or type mismatch.

---

### `Int` {#int}

```go
func (c *Context) Int(name string) int
```

Int returns the integer value of a flag specified by name. It checks command flags first, then global flags. Returns 0 if not found or type mismatch.

---

### `String` {#string}

```go
func (c *Context) String(name string) string
```

String returns the string value of a flag specified by name. It checks command flags first, then global flags. Returns "" if not found or type mismatch.

---

### `StringSlice` {#stringslice}

```go
func (c *Context) StringSlice(name string) []string
```

StringSlice returns the []string value of a flag specified by name. It checks command flags first, then global flags. Returns nil if not found or type mismatch.

---

### `DocumentConfig` {#documentconfig}

```go
type DocumentConfig struct {
	Lang        string        // Lang attribute for <html> tag, defaults to "en".
	Title       string        // Content for <title> tag, defaults to "Document".
	Charset     string        // Charset for <meta charset>, defaults to "utf-8".
	Viewport    string        // Content for <meta name="viewport">, defaults to "width=device-width, initial-scale=1".
	Description string        // Content for <meta name="description">. If empty, tag is omitted.
	Keywords    string        // Content for <meta name="keywords">. If empty, tag is omitted.
	Author      string        // Content for <meta name="author">. If empty, tag is omitted.
	HeadExtras  []HTMLElement // Additional HTMLElements to be included in the <head> section.
}
```

DocumentConfig provides configuration options for creating a new HTML document. Fields left as zero-values (e.g., empty strings) will use sensible defaults or be omitted if optional (like Description, Keywords, Author).

---

### `ETagConfig` {#etagconfig}

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

---

### `Element` {#element}

```go
type Element struct {
	// contains filtered or unexported fields
}
```

Element represents an HTML element with attributes, content, and child elements. It supports both self-closing elements (like <img>) and container elements (like <div>). Elements can be chained using fluent API methods for convenient construction.

---

#### Associated Functions

### `A` {#a}

```go
func A(href string, content ...HTMLElement) *Element
```

A creates an <a> anchor element.

---

### `Abbr` {#abbr}

```go
func Abbr(content ...HTMLElement) *Element
```

Abbr creates an <abbr> abbreviation element.

---

### `Address` {#address}

```go
func Address(content ...HTMLElement) *Element
```

Address creates an <address> semantic element.

---

### `Article` {#article}

```go
func Article(content ...HTMLElement) *Element
```

Article creates an <article> semantic element.

---

### `Aside` {#aside}

```go
func Aside(content ...HTMLElement) *Element
```

Aside creates an <aside> semantic element.

---

### `Audio` {#audio}

```go
func Audio(content ...HTMLElement) *Element
```

Audio creates an <audio> element.

---

### `B` {#b}

```go
func B(content ...HTMLElement) *Element
```

B creates a <b> element for stylistically offset text.

---

### `Base` {#base}

```go
func Base(href string) *Element
```

Base creates a <base> element.

---

### `Blockquote` {#blockquote}

```go
func Blockquote(content ...HTMLElement) *Element
```

Blockquote creates a <blockquote> element.

---

### `Body` {#body}

```go
func Body(content ...HTMLElement) *Element
```

Body creates a <body> element.

---

### `Br` {#br}

```go
func Br() *Element
```

Br creates a self-closing <br> line break element.

---

### `Button` {#button}

```go
func Button(content ...HTMLElement) *Element
```

Button creates a <button> element.

---

### `ButtonInput` {#buttoninput}

```go
func ButtonInput(valueText string) *Element
```

ButtonInput creates an <input type="button">.

---

### `Caption` {#caption}

```go
func Caption(content ...HTMLElement) *Element
```

Caption creates a <caption> element for a table.

---

### `CheckboxInput` {#checkboxinput}

```go
func CheckboxInput(name string) *Element
```

CheckboxInput creates an <input type="checkbox">.

---

### `Cite` {#cite}

```go
func Cite(content ...HTMLElement) *Element
```

Cite creates a <cite> element.

---

### `Code` {#code}

```go
func Code(content ...HTMLElement) *Element
```

Code creates a <code> element.

---

### `Col` {#col}

```go
func Col() *Element
```

Col creates a <col> element. It's self-closing.

---

### `Colgroup` {#colgroup}

```go
func Colgroup(content ...HTMLElement) *Element
```

Colgroup creates a <colgroup> element.

---

### `ColorInput` {#colorinput}

```go
func ColorInput(name string) *Element
```

ColorInput creates an <input type="color"> field.

---

### `Datalist` {#datalist}

```go
func Datalist(id string, content ...HTMLElement) *Element
```

Datalist creates a <datalist> element.

---

### `DateInput` {#dateinput}

```go
func DateInput(name string) *Element
```

DateInput creates an <input type="date"> field.

---

### `DateTimeLocalInput` {#datetimelocalinput}

```go
func DateTimeLocalInput(name string) *Element
```

DateTimeLocalInput creates an <input type="datetime-local"> field.

---

### `Details` {#details}

```go
func Details(content ...HTMLElement) *Element
```

Details creates a <details> element.

---

### `Dfn` {#dfn}

```go
func Dfn(content ...HTMLElement) *Element
```

Dfn creates a <dfn> definition element.

---

### `DialogEl` {#dialogel}

```go
func DialogEl(content ...HTMLElement) *Element
```

DialogEl creates a <dialog> element. Renamed to avoid potential conflicts.

---

### `Div` {#div}

```go
func Div(content ...HTMLElement) *Element
```

Div creates a <div> element.

---

### `Em` {#em}

```go
func Em(content ...HTMLElement) *Element
```

Em creates an <em> emphasis element.

---

### `EmailInput` {#emailinput}

```go
func EmailInput(name string) *Element
```

EmailInput creates an <input type="email"> field.

---

### `EmbedEl` {#embedel}

```go
func EmbedEl(src string, embedType string) *Element
```

EmbedEl creates an <embed> element. It's self-closing. Renamed to avoid keyword conflict.

---

### `Favicon` {#favicon}

```go
func Favicon(href string, rel ...string) *Element
```

Favicon creates a <link> element for a favicon.

---

### `Fieldset` {#fieldset}

```go
func Fieldset(content ...HTMLElement) *Element
```

Fieldset creates a <fieldset> element.

---

### `Figcaption` {#figcaption}

```go
func Figcaption(content ...HTMLElement) *Element
```

Figcaption creates a <figcaption> element.

---

### `Figure` {#figure}

```go
func Figure(content ...HTMLElement) *Element
```

Figure creates a <figure> element.

---

### `FileInput` {#fileinput}

```go
func FileInput(name string) *Element
```

FileInput creates an <input type="file"> field.

---

### `Footer` {#footer}

```go
func Footer(content ...HTMLElement) *Element
```

Footer creates a <footer> semantic element.

---

### `Form` {#form}

```go
func Form(content ...HTMLElement) *Element
```

Form creates a <form> element.

---

### `H1` {#h1}

```go
func H1(content ...HTMLElement) *Element
```

H1 creates an <h1> heading element.

---

### `H2` {#h2}

```go
func H2(content ...HTMLElement) *Element
```

H2 creates an <h2> heading element.

---

### `H3` {#h3}

```go
func H3(content ...HTMLElement) *Element
```

H3 creates an <h3> heading element.

---

### `H4` {#h4}

```go
func H4(content ...HTMLElement) *Element
```

H4 creates an <h4> heading element.

---

### `H5` {#h5}

```go
func H5(content ...HTMLElement) *Element
```

H5 creates an <h5> heading element.

---

### `H6` {#h6}

```go
func H6(content ...HTMLElement) *Element
```

H6 creates an <h6> heading element.

---

### `Head` {#head}

```go
func Head(content ...HTMLElement) *Element
```

Head creates a <head> element.

---

### `Header` {#header}

```go
func Header(content ...HTMLElement) *Element
```

Header creates a <header> semantic element.

---

### `HiddenInput` {#hiddeninput}

```go
func HiddenInput(name string, value string) *Element
```

HiddenInput creates an <input type="hidden"> field.

---

### `Hr` {#hr}

```go
func Hr() *Element
```

Hr creates a self-closing <hr> horizontal rule element.

---

### `Html` {#html}

```go
func Html(content ...HTMLElement) *Element
```

Html creates an <html> element.

---

### `I` {#i}

```go
func I(content ...HTMLElement) *Element
```

I creates an <i> idiomatic text element.

---

### `Iframe` {#iframe}

```go
func Iframe(src string) *Element
```

Iframe creates an <iframe> element.

---

### `Image` {#image}

```go
func Image(src, alt string) *Element
```

Image creates an <img> element (alias for Img).

---

### `Img` {#img}

```go
func Img(src, alt string) *Element
```

Img creates a self-closing <img> element.

---

### `InlineScript` {#inlinescript}

```go
func InlineScript(scriptContent string) *Element
```

InlineScript creates a <script> element with inline JavaScript content.

---

### `Input` {#input}

```go
func Input(inputType string) *Element
```

Input creates a self-closing <input> element.

---

### `Kbd` {#kbd}

```go
func Kbd(content ...HTMLElement) *Element
```

Kbd creates a <kbd> keyboard input element.

---

### `Label` {#label}

```go
func Label(content ...HTMLElement) *Element
```

Label creates a <label> element.

---

### `Legend` {#legend}

```go
func Legend(content ...HTMLElement) *Element
```

Legend creates a <legend> element.

---

### `Li` {#li}

```go
func Li(content ...HTMLElement) *Element
```

Li creates a <li> list item element.

---

### `Link` {#link}

```go
func Link(href, textContent string) *Element
```

Link creates an <a> anchor element with href and text content.

---

### `LinkTag` {#linktag}

```go
func LinkTag() *Element
```

LinkTag creates a generic <link> element. It's self-closing.

---

### `Main` {#main}

```go
func Main(content ...HTMLElement) *Element
```

Main creates a <main> semantic element.

---

### `Mark` {#mark}

```go
func Mark(content ...HTMLElement) *Element
```

Mark creates a <mark> element.

---

### `Meta` {#meta}

```go
func Meta() *Element
```

Meta creates a generic <meta> element. It's self-closing.

---

### `MetaCharset` {#metacharset}

```go
func MetaCharset(charset string) *Element
```

MetaCharset creates a <meta charset="..."> element.

---

### `MetaNameContent` {#metanamecontent}

```go
func MetaNameContent(name, contentVal string) *Element
```

MetaNameContent creates a <meta name="..." content="..."> element.

---

### `MetaPropertyContent` {#metapropertycontent}

```go
func MetaPropertyContent(property, contentVal string) *Element
```

MetaPropertyContent creates a <meta property="..." content="..."> element.

---

### `MetaViewport` {#metaviewport}

```go
func MetaViewport(contentVal string) *Element
```

MetaViewport creates a <meta name="viewport" content="..."> element.

---

### `MeterEl` {#meterel}

```go
func MeterEl(content ...HTMLElement) *Element
```

MeterEl creates a <meter> element. Renamed to avoid potential conflicts.

---

### `MonthInput` {#monthinput}

```go
func MonthInput(name string) *Element
```

MonthInput creates an <input type="month"> field.

---

### `Nav` {#nav}

```go
func Nav(content ...HTMLElement) *Element
```

Nav creates a <nav> semantic element.

---

### `NoScript` {#noscript}

```go
func NoScript(content ...HTMLElement) *Element
```

NoScript creates a <noscript> element.

---

### `NumberInput` {#numberinput}

```go
func NumberInput(name string) *Element
```

NumberInput creates an <input type="number"> field.

---

### `ObjectEl` {#objectel}

```go
func ObjectEl(content ...HTMLElement) *Element
```

ObjectEl creates an <object> element. Renamed to avoid keyword conflict.

---

### `Ol` {#ol}

```go
func Ol(content ...HTMLElement) *Element
```

Ol creates an <ol> ordered list element.

---

### `Optgroup` {#optgroup}

```go
func Optgroup(label string, content ...HTMLElement) *Element
```

Optgroup creates an <optgroup> element.

---

### `Option` {#option}

```go
func Option(value string, content ...HTMLElement) *Element
```

Option creates an <option> element.

---

### `OutputEl` {#outputel}

```go
func OutputEl(content ...HTMLElement) *Element
```

OutputEl creates an <output> element. Renamed to avoid potential conflicts.

---

### `P` {#p}

```go
func P(content ...HTMLElement) *Element
```

P creates a <p> paragraph element.

---

### `Param` {#param}

```go
func Param(name, value string) *Element
```

Param creates a <param> element. It's self-closing.

---

### `PasswordInput` {#passwordinput}

```go
func PasswordInput(name string) *Element
```

PasswordInput creates an <input type="password"> field.

---

### `Pre` {#pre}

```go
func Pre(content ...HTMLElement) *Element
```

Pre creates a <pre> element.

---

### `Preload` {#preload}

```go
func Preload(href string, asType string) *Element
```

Preload creates a <link rel="preload"> element.

---

### `ProgressEl` {#progressel}

```go
func ProgressEl(content ...HTMLElement) *Element
```

ProgressEl creates a <progress> element. Renamed to avoid potential conflicts.

---

### `Q` {#q}

```go
func Q(content ...HTMLElement) *Element
```

Q creates a <q> inline quotation element.

---

### `RadioInput` {#radioinput}

```go
func RadioInput(name, value string) *Element
```

RadioInput creates an <input type="radio">.

---

### `RangeInput` {#rangeinput}

```go
func RangeInput(name string) *Element
```

RangeInput creates an <input type="range"> field.

---

### `ResetButton` {#resetbutton}

```go
func ResetButton(text string) *Element
```

ResetButton creates a <button type="reset">.

---

### `Samp` {#samp}

```go
func Samp(content ...HTMLElement) *Element
```

Samp creates a <samp> sample output element.

---

### `Script` {#script}

```go
func Script(src string) *Element
```

Script creates a <script> element for external JavaScript files.

---

### `SearchInput` {#searchinput}

```go
func SearchInput(name string) *Element
```

SearchInput creates an <input type="search"> field.

---

### `Section` {#section}

```go
func Section(content ...HTMLElement) *Element
```

Section creates a <section> semantic element.

---

### `Select` {#select}

```go
func Select(content ...HTMLElement) *Element
```

Select creates a <select> dropdown element.

---

### `Small` {#small}

```go
func Small(content ...HTMLElement) *Element
```

Small creates a <small> element.

---

### `Source` {#source}

```go
func Source(src string, mediaType string) *Element
```

Source creates a <source> element. It's self-closing.

---

### `Span` {#span}

```go
func Span(content ...HTMLElement) *Element
```

Span creates a <span> inline element.

---

### `Strong` {#strong}

```go
func Strong(content ...HTMLElement) *Element
```

Strong creates a <strong> element.

---

### `StyleSheet` {#stylesheet}

```go
func StyleSheet(href string) *Element
```

StyleSheet creates a <link rel="stylesheet"> element.

---

### `StyleTag` {#styletag}

```go
func StyleTag(cssContent string) *Element
```

StyleTag creates a <style> element for embedding CSS.

---

### `Sub` {#sub}

```go
func Sub(content ...HTMLElement) *Element
```

Sub creates a <sub> subscript element.

---

### `SubmitButton` {#submitbutton}

```go
func SubmitButton(text string) *Element
```

SubmitButton creates a <button type="submit">.

---

### `Summary` {#summary}

```go
func Summary(content ...HTMLElement) *Element
```

Summary creates a <summary> element.

---

### `Sup` {#sup}

```go
func Sup(content ...HTMLElement) *Element
```

Sup creates a <sup> superscript element.

---

### `Table` {#table}

```go
func Table(content ...HTMLElement) *Element
```

Table creates a <table> element.

---

### `Tbody` {#tbody}

```go
func Tbody(content ...HTMLElement) *Element
```

Tbody creates a <tbody> table body group element.

---

### `Td` {#td}

```go
func Td(content ...HTMLElement) *Element
```

Td creates a <td> table data cell element.

---

### `TelInput` {#telinput}

```go
func TelInput(name string) *Element
```

TelInput creates an <input type="tel"> field.

---

### `TextInput` {#textinput}

```go
func TextInput(name string) *Element
```

TextInput creates an <input type="text"> field.

---

### `Textarea` {#textarea}

```go
func Textarea(content ...HTMLElement) *Element
```

Textarea creates a <textarea> element.

---

### `Th` {#th}

```go
func Th(content ...HTMLElement) *Element
```

Th creates a <th> table header cell element.

---

### `Thead` {#thead}

```go
func Thead(content ...HTMLElement) *Element
```

Thead creates a <thead> table header group element.

---

### `TimeEl` {#timeel}

```go
func TimeEl(content ...HTMLElement) *Element
```

TimeEl creates a <time> element. Renamed to avoid potential conflicts.

---

### `TimeInput` {#timeinput}

```go
func TimeInput(name string) *Element
```

TimeInput creates an <input type="time"> field.

---

### `TitleEl` {#titleel}

```go
func TitleEl(titleText string) *Element
```

TitleEl creates a <title> element with the specified text. Renamed from Title to TitleEl to avoid conflict with (*Element).Title method if it existed.

---

### `Tr` {#tr}

```go
func Tr(content ...HTMLElement) *Element
```

Tr creates a <tr> table row element.

---

### `Track` {#track}

```go
func Track(kind, src, srclang string) *Element
```

Track creates a <track> element. It's self-closing.

---

### `U` {#u}

```go
func U(content ...HTMLElement) *Element
```

U creates a <u> unarticulated annotation element.

---

### `Ul` {#ul}

```go
func Ul(content ...HTMLElement) *Element
```

Ul creates a <ul> unordered list element.

---

### `UrlInput` {#urlinput}

```go
func UrlInput(name string) *Element
```

UrlInput creates an <input type="url"> field.

---

### `VarEl` {#varel}

```go
func VarEl(content ...HTMLElement) *Element
```

VarEl creates a <var> variable element. Renamed to avoid keyword conflict.

---

### `Video` {#video}

```go
func Video(content ...HTMLElement) *Element
```

Video creates a <video> element.

---

### `Wbr` {#wbr}

```go
func Wbr() *Element
```

Wbr creates a <wbr> word break opportunity element. It's self-closing.

---

### `WeekInput` {#weekinput}

```go
func WeekInput(name string) *Element
```

WeekInput creates an <input type="week"> field.

---

#### Methods

### `Add` {#add}

```go
func (e *Element) Add(children ...HTMLElement) *Element
```

Add appends child elements to this element.

---

### `Attr` {#attr}

```go
func (e *Element) Attr(key, value string) *Element
```

Attr sets an attribute on the element and returns the element for method chaining.

---

### `BoolAttr` {#boolattr}

```go
func (e *Element) BoolAttr(key string, present bool) *Element
```

BoolAttr sets or removes a boolean attribute on the element. If present is true, the attribute is added (e.g., <input disabled>). If present is false, the attribute is removed if it exists.

---

### `Class` {#class}

```go
func (e *Element) Class(class string) *Element
```

Class sets the class attribute on the element.

---

### `ID` {#id}

```go
func (e *Element) ID(id string) *Element
```

ID sets the id attribute on the element.

---

### `Render` {#render}

```go
func (e *Element) Render() string
```

Render converts the element to its HTML string representation. It handles both self-closing and container elements, attributes, content, and children. The output is properly formatted HTML that can be sent to browsers. Content and attribute values are HTML-escaped to prevent XSS, except for specific tags like <script> and <style> whose content must be raw.

---

### `Style` {#style}

```go
func (e *Element) Style(style string) *Element
```

Style sets the style attribute on the element.

---

### `Text` {#text}

```go
func (e *Element) Text(text string) *Element
```

Text sets the text content of the element. This content is HTML-escaped during rendering.

---

### `EnforceContentTypeConfig` {#enforcecontenttypeconfig}

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

---

### `Flag` {#flag}

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

Flag defines the interface for command-line flags. Concrete types like StringFlag, IntFlag, BoolFlag implement this interface.

---

### `Float64Flag` {#float64flag}

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

---

#### Methods

### `Apply` {#apply}

```go
func (f *Float64Flag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the float64 flag with the flag.FlagSet.

---

### `GetAliases` {#getaliases}

```go
func (f *Float64Flag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

---

### `GetName` {#getname}

```go
func (f *Float64Flag) GetName() string
```

GetName returns the primary name of the flag.

---

### `IsRequired` {#isrequired}

```go
func (f *Float64Flag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

---

### `String` {#string}

```go
func (f *Float64Flag) String() string
```

String returns the help text representation of the flag.

---

### `Validate` {#validate}

```go
func (f *Float64Flag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided.

---

### `ForceHTTPSConfig` {#forcehttpsconfig}

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

---

### `Group` {#group}

```go
type Group struct {
	// contains filtered or unexported fields
}
```

Group is a lightweight helper that allows users to register a set of routes that share a common prefix and/or middleware. It delegates to the parent router while applying its own prefix and middleware chain.

---

#### Methods

### `Delete` {#delete}

```go
func (g *Group) Delete(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Delete registers a new route for HTTP DELETE requests within the group.

---

### `DeleteFunc` {#deletefunc}

```go
func (g *Group) DeleteFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

DeleteFunc registers a new enhanced route for HTTP DELETE requests within the group.

---

### `Get` {#get}

```go
func (g *Group) Get(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Get registers a new route for HTTP GET requests within the group.

---

### `GetFunc` {#getfunc}

```go
func (g *Group) GetFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

GetFunc registers a new enhanced route for HTTP GET requests within the group.

---

### `Handle` {#handle}

```go
func (g *Group) Handle(method, pattern string, handler http.HandlerFunc, opts ...*RouteOptions)
```

Handle registers a new route within the group, applying the group's prefix and middleware. The route is ultimately registered with the parent router after transformations.

---

### `HandleFunc` {#handlefunc}

```go
func (g *Group) HandleFunc(method, pattern string, handler HandlerFunc, opts ...*RouteOptions)
```

HandleFunc registers a new enhanced route within the group, applying the group's prefix and middleware. Uses the enhanced handler signature for better error handling.

---

### `Patch` {#patch}

```go
func (g *Group) Patch(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Patch registers a new route for HTTP PATCH requests within the group.

---

### `PatchFunc` {#patchfunc}

```go
func (g *Group) PatchFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PatchFunc registers a new enhanced route for HTTP PATCH requests within the group.

---

### `Post` {#post}

```go
func (g *Group) Post(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Post registers a new route for HTTP POST requests within the group.

---

### `PostFunc` {#postfunc}

```go
func (g *Group) PostFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PostFunc registers a new enhanced route for HTTP POST requests within the group.

---

### `Put` {#put}

```go
func (g *Group) Put(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Put registers a new route for HTTP PUT requests within the group.

---

### `PutFunc` {#putfunc}

```go
func (g *Group) PutFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PutFunc registers a new enhanced route for HTTP PUT requests within the group.

---

### `Use` {#use}

```go
func (g *Group) Use(mws ...Middleware)
```

Use adds middleware functions to the group. These middleware functions apply only to routes registered through the group, not to the parent router.

---

### `GzipConfig` {#gzipconfig}

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

---

### `HTMLDocument` {#htmldocument}

```go
type HTMLDocument struct {
	// contains filtered or unexported fields
}
```

HTMLDocument represents a full HTML document, including the DOCTYPE. Its Render method produces the complete HTML string.

---

#### Associated Functions

### `Document` {#document}

```go
func Document(config DocumentConfig, bodyContent ...HTMLElement) *HTMLDocument
```

Document creates a complete HTML5 document structure, encapsulated in an HTMLDocument. The returned HTMLDocument's Render method will produce the full HTML string, including the DOCTYPE.

It uses DocumentConfig to customize the document's head and html attributes, and accepts variadic arguments for the body content. Sensible defaults are applied for common attributes and meta tags if not specified in the config.

---

#### Methods

### `Render` {#render}

```go
func (d *HTMLDocument) Render() string
```

Render converts the HTMLDocument to its full string representation, prepending the HTML5 DOCTYPE declaration.

---

### `HTMLElement` {#htmlelement}

```go
type HTMLElement interface {
	Render() string
}
```

HTMLElement represents any HTML element that can be rendered to a string. This interface allows for composition of complex HTML structures using both predefined elements and custom implementations.

---

#### Associated Functions

### `Text` {#text}

```go
func Text(text string) HTMLElement
```

Text creates a raw text node.

---

### `HandlerFunc` {#handlerfunc}

```go
type HandlerFunc func(*ResponseContext) error
```

HandlerFunc is an enhanced handler function that receives a ResponseContext and returns an error. This allows for cleaner error handling and response management compared to standard http.HandlerFunc.

---

### `HeaderObject` {#headerobject}

```go
type HeaderObject struct {
	Description string        `json:"description,omitempty"`
	Schema      *SchemaObject `json:"schema"`
}
```

HeaderObject describes a response header in an Operation.

---

### `HealthCheckConfig` {#healthcheckconfig}

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

---

### `IPFilterConfig` {#ipfilterconfig}

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

---

### `Info` {#info}

```go
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}
```

Info provides metadata about the API: title, version, and optional description.

---

### `IntFlag` {#intflag}

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

---

#### Methods

### `Apply` {#apply}

```go
func (f *IntFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the int flag with the flag.FlagSet.

---

### `GetAliases` {#getaliases}

```go
func (f *IntFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

---

### `GetName` {#getname}

```go
func (f *IntFlag) GetName() string
```

GetName returns the primary name of the flag.

---

### `IsRequired` {#isrequired}

```go
func (f *IntFlag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

---

### `String` {#string}

```go
func (f *IntFlag) String() string
```

String returns the help text representation of the flag.

---

### `Validate` {#validate}

```go
func (f *IntFlag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided.

---

### `LoggingConfig` {#loggingconfig}

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

---

### `MaintenanceModeConfig` {#maintenancemodeconfig}

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

---

### `MaxRequestBodySizeConfig` {#maxrequestbodysizeconfig}

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

---

### `MediaTypeObject` {#mediatypeobject}

```go
type MediaTypeObject struct {
	Schema *SchemaObject `json:"schema,omitempty"`
}
```

MediaTypeObject holds the schema defining the media type for a request or response.

---

### `MethodOverrideConfig` {#methodoverrideconfig}

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

---

### `Middleware` {#middleware}

```go
type Middleware func(http.Handler) http.Handler
```

Middleware defines the function signature for middleware. A middleware is a function that wraps an http.Handler, adding extra behavior such as logging, authentication, or request modification.

---

#### Associated Functions

### `BasicAuthMiddleware` {#basicauthmiddleware}

```go
func BasicAuthMiddleware(config BasicAuthConfig) Middleware
```

BasicAuthMiddleware provides simple HTTP Basic Authentication.

---

### `CORSMiddleware` {#corsmiddleware}

```go
func CORSMiddleware(config CORSConfig) Middleware
```

CORSMiddleware sets Cross-Origin Resource Sharing headers.

---

### `CSRFMiddleware` {#csrfmiddleware}

```go
func CSRFMiddleware(config *CSRFConfig) Middleware
```

CSRFMiddleware provides Cross-Site Request Forgery protection. It uses the "Double Submit Cookie" pattern. A random token is generated and set in a secure, HttpOnly cookie. For unsafe methods (POST, PUT, etc.), the middleware expects the same token to be present in a request header (e.g., X-CSRF-Token) or form field, sent by the frontend JavaScript.

---

### `CacheControlMiddleware` {#cachecontrolmiddleware}

```go
func CacheControlMiddleware(config CacheControlConfig) Middleware
```

CacheControlMiddleware sets the Cache-Control header for responses.

---

### `ConcurrencyLimiterMiddleware` {#concurrencylimitermiddleware}

```go
func ConcurrencyLimiterMiddleware(config ConcurrencyLimiterConfig) Middleware
```

ConcurrencyLimiterMiddleware limits the number of concurrent requests.

---

### `ETagMiddleware` {#etagmiddleware}

```go
func ETagMiddleware(config *ETagConfig) Middleware
```

ETagMiddleware adds ETag headers to responses and handles If-None-Match conditional requests, potentially returning a 304 Not Modified status. Note: This middleware buffers the entire response body in memory to calculate the ETag hash. This may be inefficient for very large responses.

---

### `EnforceContentTypeMiddleware` {#enforcecontenttypemiddleware}

```go
func EnforceContentTypeMiddleware(config EnforceContentTypeConfig) Middleware
```

EnforceContentTypeMiddleware checks if the request's Content-Type header is allowed.

---

### `ForceHTTPSMiddleware` {#forcehttpsmiddleware}

```go
func ForceHTTPSMiddleware(config ForceHTTPSConfig) Middleware
```

ForceHTTPSMiddleware redirects HTTP requests to HTTPS.

---

### `GzipMiddleware` {#gzipmiddleware}

```go
func GzipMiddleware(config *GzipConfig) Middleware
```

GzipMiddleware returns middleware that compresses response bodies using gzip if the client indicates support via the Accept-Encoding header.

---

### `HealthCheckMiddleware` {#healthcheckmiddleware}

```go
func HealthCheckMiddleware(config *HealthCheckConfig) Middleware
```

HealthCheckMiddleware provides a simple health check endpoint.

---

### `IPFilterMiddleware` {#ipfiltermiddleware}

```go
func IPFilterMiddleware(config IPFilterConfig) Middleware
```

IPFilterMiddleware restricts access based on client IP address.

---

### `LoggingMiddleware` {#loggingmiddleware}

```go
func LoggingMiddleware(config *LoggingConfig) Middleware
```

LoggingMiddleware logs request details including method, path, status, size, and duration.

---

### `MaintenanceModeMiddleware` {#maintenancemodemiddleware}

```go
func MaintenanceModeMiddleware(config MaintenanceModeConfig) Middleware
```

MaintenanceModeMiddleware returns a 503 Service Unavailable if enabled, allowing bypass for specific IPs.

---

### `MaxRequestBodySizeMiddleware` {#maxrequestbodysizemiddleware}

```go
func MaxRequestBodySizeMiddleware(config MaxRequestBodySizeConfig) Middleware
```

MaxRequestBodySizeMiddleware limits the size of incoming request bodies.

---

### `MethodOverrideMiddleware` {#methodoverridemiddleware}

```go
func MethodOverrideMiddleware(config *MethodOverrideConfig) Middleware
```

MethodOverrideMiddleware checks a header or form field to override the request method.

---

### `RateLimitMiddleware` {#ratelimitmiddleware}

```go
func RateLimitMiddleware(config RateLimiterConfig) Middleware
```

RateLimitMiddleware provides basic in-memory rate limiting. Warning: This simple implementation has limitations:
 - Memory usage can grow indefinitely without CleanupInterval set.
 - Not suitable for distributed systems (limit is per instance).
 - Accuracy decreases slightly under very high concurrency due to locking.

---

### `RealIPMiddleware` {#realipmiddleware}

```go
func RealIPMiddleware(config RealIPConfig) Middleware
```

RealIPMiddleware extracts the client's real IP address from proxy headers. Warning: Only use this if you have a trusted proxy setting these headers correctly.

---

### `RecoveryMiddleware` {#recoverymiddleware}

```go
func RecoveryMiddleware(config *RecoveryConfig) Middleware
```

RecoveryMiddleware recovers from panics in downstream handlers.

---

### `RequestIDMiddleware` {#requestidmiddleware}

```go
func RequestIDMiddleware(config *RequestIDConfig) Middleware
```

RequestIDMiddleware retrieves a request ID from a header or generates one. It sets the ID in the response header and request context.

---

### `SecurityHeadersMiddleware` {#securityheadersmiddleware}

```go
func SecurityHeadersMiddleware(config SecurityHeadersConfig) Middleware
```

SecurityHeadersMiddleware sets common security headers.

---

### `TimeoutMiddleware` {#timeoutmiddleware}

```go
func TimeoutMiddleware(config TimeoutConfig) Middleware
```

TimeoutMiddleware sets a maximum duration for handling requests.

---

### `TrailingSlashRedirectMiddleware` {#trailingslashredirectmiddleware}

```go
func TrailingSlashRedirectMiddleware(config TrailingSlashRedirectConfig) Middleware
```

TrailingSlashRedirectMiddleware redirects requests to add or remove a trailing slash.

---

### `OpenAPI` {#openapi}

```go
type OpenAPI struct {
	OpenAPI    string               `json:"openapi"`
	Info       Info                 `json:"info"`
	Paths      map[string]*PathItem `json:"paths"`
	Components *Components          `json:"components,omitempty"`
}
```

OpenAPI is the root document object for an OpenAPI 3 specification.

---

#### Associated Functions

### `GenerateOpenAPISpec` {#generateopenapispec}

```go
func GenerateOpenAPISpec(router *Router, config OpenAPIConfig) *OpenAPI
```

GenerateOpenAPISpec constructs an OpenAPI 3.0 specification from the given router and configuration, including paths, operations, and components.

---

### `OpenAPIConfig` {#openapiconfig}

```go
type OpenAPIConfig struct {
	Title       string // Title of the API
	Version     string // Version of the API
	Description string // Description of the API
}
```

OpenAPIConfig holds metadata for generating an OpenAPI specification.

---

### `Operation` {#operation}

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

---

### `ParameterObject` {#parameterobject}

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

---

### `ParameterOption` {#parameteroption}

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

ParameterOption configures an Operation parameter. Name and In ("query", "header", "path", "cookie") are required. Required, Schema and Example further customize the generated parameter.

---

### `PathItem` {#pathitem}

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

PathItem describes the operations available on a single API path. One of Get, Post, Put, Delete, Patch may be non-nil.

---

### `RateLimiterConfig` {#ratelimiterconfig}

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

---

### `RealIPConfig` {#realipconfig}

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

---

### `RecoveryConfig` {#recoveryconfig}

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

---

### `RequestBodyObject` {#requestbodyobject}

```go
type RequestBodyObject struct {
	Description string                      `json:"description,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content"`
	Required    bool                        `json:"required,omitempty"`
}
```

RequestBodyObject describes a request body for an Operation.

---

### `RequestIDConfig` {#requestidconfig}

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

---

### `ResponseContext` {#responsecontext}

```go
type ResponseContext struct {
	// contains filtered or unexported fields
}
```

ResponseContext provides helper methods for sending HTTP responses with reduced boilerplate. It wraps the http.ResponseWriter and http.Request along with router context to provide convenient methods for JSON, HTML, text responses, and automatic data binding.

---

#### Methods

### `Bind` {#bind}

```go
func (rc *ResponseContext) Bind(v any) error
```

Bind automatically unmarshals and validates request data into the provided struct. It supports both JSON and form data, with automatic content-type detection. The struct should be a pointer. Validation is performed if validation middleware is active.

---

### `BindForm` {#bindform}

```go
func (rc *ResponseContext) BindForm(v any) error
```

BindForm binds form data into the provided struct using reflection. The struct should be a pointer. Field names are matched using JSON tags or struct field names. Supports string, bool, int, and float fields with automatic type conversion.

---

### `BindJSON` {#bindjson}

```go
func (rc *ResponseContext) BindJSON(v any) error
```

BindJSON unmarshals JSON request body into the provided struct. The struct should be a pointer. Returns an error if the body is nil or JSON is invalid.

---

### `BindValidated` {#bindvalidated}

```go
func (rc *ResponseContext) BindValidated(v any) error
```

BindValidated binds and validates request data (JSON or form) with comprehensive validation

---

### `HTML` {#html}

```go
func (rc *ResponseContext) HTML(statusCode int, content HTMLElement) error
```

HTML sends an HTML response with the given status code and content. It automatically sets the Content-Type header to "text/html; charset=utf-8" and renders the provided HTMLElement to string.

---

### `JSON` {#json}

```go
func (rc *ResponseContext) JSON(statusCode int, data any) error
```

JSON sends a JSON response with the given status code and data. It automatically sets the Content-Type header to "application/json" and handles JSON encoding. Returns an error if encoding fails.

---

### `JSONError` {#jsonerror}

```go
func (rc *ResponseContext) JSONError(statusCode int, message string) error
```

JSONError sends a JSON error response with the given status code and message. It creates a standardized error response format with an "error" field.

---

### `Redirect` {#redirect}

```go
func (rc *ResponseContext) Redirect(statusCode int, url string) error
```

Redirect sends an HTTP redirect response with the given status code and URL. Common status codes are 301 (permanent), 302 (temporary), and 307 (temporary, preserve method).

---

### `Request` {#request}

```go
func (rc *ResponseContext) Request() *http.Request
```

Request returns the underlying http.Request for advanced request handling when the ResponseContext helpers are not sufficient.

---

### `Text` {#text}

```go
func (rc *ResponseContext) Text(statusCode int, text string) error
```

Text sends a plain text response with the given status code and text content. It automatically sets the Content-Type header to "text/plain; charset=utf-8".

---

### `URLParam` {#urlparam}

```go
func (rc *ResponseContext) URLParam(key string) string
```

URLParam retrieves the URL parameter for the given key from the request context. Returns an empty string if the parameter is not present or if the request doesn't contain URL parameters.

---

### `WantsJSON` {#wantsjson}

```go
func (rc *ResponseContext) WantsJSON() bool
```

WantsJSON returns true if the request expects a JSON response based on Content-Type or Accept headers. Useful for dual HTML/JSON endpoints.

---

### `Writer` {#writer}

```go
func (rc *ResponseContext) Writer() http.ResponseWriter
```

Writer returns the underlying http.ResponseWriter for advanced response handling when the ResponseContext helpers are not sufficient.

---

### `ResponseObject` {#responseobject}

```go
type ResponseObject struct {
	Description string                      `json:"description"`
	Headers     map[string]*HeaderObject    `json:"headers,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content,omitempty"`
}
```

ResponseObject describes a single response in an Operation.

---

### `ResponseOption` {#responseoption}

```go
type ResponseOption struct {
	Description string
	Body        any
}
```

ResponseOption configures a single HTTP response in an Operation. Description is required; Body, if non-nil, is used to generate a schema.

---

### `RouteOptions` {#routeoptions}

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

RouteOptions holds OpenAPI metadata for a single route. Tags, Summary, Description, and OperationID map directly into the corresponding Operation fields. RequestBody, Responses, and Parameters drive schema generation for request bodies, responses, and parameters.

---

### `Router` {#router}

```go
type Router struct {
	// contains filtered or unexported fields
}
```

Router is a minimal HTTP router that supports dynamic routes with regex validation in path parameters and can mount subrouters. It provides middleware support, custom error handlers, and both traditional and enhanced handler functions.

---

#### Associated Functions

### `NewRouter` {#newrouter}

```go
func NewRouter() *Router
```

NewRouter creates and returns a new Router instance with default configuration.

---

#### Methods

### `Delete` {#delete}

```go
func (r *Router) Delete(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Delete registers a new route for HTTP DELETE requests using the standard handler signature.

---

### `DeleteFunc` {#deletefunc}

```go
func (r *Router) DeleteFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

DeleteFunc registers a new route for HTTP DELETE requests using the enhanced handler signature.

---

### `Get` {#get}

```go
func (r *Router) Get(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Get registers a new route for HTTP GET requests using the standard handler signature.

---

### `GetFunc` {#getfunc}

```go
func (r *Router) GetFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

GetFunc registers a new route for HTTP GET requests using the enhanced handler signature.

---

### `Group` {#group}

```go
func (r *Router) Group(prefix string, mws ...Middleware) *Group
```

Group creates and returns a new Group with the given prefix. A group is a lightweight convenience wrapper that prefixes routes and can add its own middleware without creating a separate router instance.

---

### `Handle` {#handle}

```go
func (r *Router) Handle(method, pattern string, handler http.HandlerFunc, opts ...*RouteOptions)
```

Handle registers a new route with the given HTTP method, URL pattern, and handler. If the router has a non-empty basePath, it is automatically prepended to the pattern. Optional RouteOptions can be provided for OpenAPI documentation.

---

### `HandleFunc` {#handlefunc}

```go
func (r *Router) HandleFunc(method, pattern string, handler HandlerFunc, opts ...*RouteOptions)
```

HandleFunc registers a new enhanced route that receives a ResponseContext instead of separate ResponseWriter and Request parameters. This enables cleaner error handling and response management with automatic data binding capabilities.

---

### `Patch` {#patch}

```go
func (r *Router) Patch(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Patch registers a new route for HTTP PATCH requests using the standard handler signature.

---

### `PatchFunc` {#patchfunc}

```go
func (r *Router) PatchFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PatchFunc registers a new route for HTTP PATCH requests using the enhanced handler signature.

---

### `Post` {#post}

```go
func (r *Router) Post(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Post registers a new route for HTTP POST requests using the standard handler signature.

---

### `PostFunc` {#postfunc}

```go
func (r *Router) PostFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PostFunc registers a new route for HTTP POST requests using the enhanced handler signature.

---

### `Put` {#put}

```go
func (r *Router) Put(pattern string, handler http.HandlerFunc, options ...*RouteOptions)
```

Put registers a new route for HTTP PUT requests using the standard handler signature.

---

### `PutFunc` {#putfunc}

```go
func (r *Router) PutFunc(pattern string, handler HandlerFunc, options ...*RouteOptions)
```

PutFunc registers a new route for HTTP PUT requests using the enhanced handler signature.

---

### `ServeHTTP` {#servehttp}

```go
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

ServeHTTP implements the http.Handler interface. It first checks subrouters based on their base path, then its own routes. If a URL pattern matches but the HTTP method does not, a 405 Method Not Allowed is returned. If no routes match, a 404 Not Found is returned.

---

### `ServeOpenAPISpec` {#serveopenapispec}

```go
func (r *Router) ServeOpenAPISpec(path string, config OpenAPIConfig)
```

ServeOpenAPISpec makes your OpenAPI specification available at `path` (e.g. "/openapi.json")

---

### `ServeSwaggerUI` {#serveswaggerui}

```go
func (r *Router) ServeSwaggerUI(prefix string)
```

ServeSwaggerUI shows swaggerUI at `prefix` (e.g. "/docs")

---

### `SetMethodNotAllowedHandler` {#setmethodnotallowedhandler}

```go
func (r *Router) SetMethodNotAllowedHandler(h http.Handler)
```

SetMethodNotAllowedHandler allows users to supply a custom http.Handler for requests that match a route pattern but use an unsupported HTTP method. If not set, a default 405 Method Not Allowed response is returned.

---

### `SetNotFoundHandler` {#setnotfoundhandler}

```go
func (r *Router) SetNotFoundHandler(h http.Handler)
```

SetNotFoundHandler allows users to supply a custom http.Handler for requests that do not match any route. If not set, the standard http.NotFound handler is used.

---

### `Static` {#static}

```go
func (r *Router) Static(urlPathPrefix string, subFS fs.FS)
```

Static mounts an fs.FS under the given URL prefix.

---

### `Subrouter` {#subrouter}

```go
func (r *Router) Subrouter(prefix string) *Router
```

Subrouter creates a new Router mounted at the given prefix. The subrouter inherits the parent's context key, global middleware, and error handlers. Its routes will automatically receive the combined prefix.

---

### `URLParam` {#urlparam}

```go
func (r *Router) URLParam(req *http.Request, key string) string
```

URLParam retrieves the URL parameter for the given key from the request context. It returns an empty string if the parameter is not present. This method should only be called within route handlers that were registered with this router.

---

### `Use` {#use}

```go
func (r *Router) Use(mws ...Middleware)
```

Use registers one or more middleware functions that will be applied globally to every matched route handler. Middleware is applied in the order it is registered.

---

### `SchemaObject` {#schemaobject}

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

---

### `SecurityHeadersConfig` {#securityheadersconfig}

```go
type SecurityHeadersConfig struct {
	// ContentTypeOptions sets the X-Content-Type-Options header.
	// Defaults to "nosniff". Set to "" to disable.
	ContentTypeOptions string
	// FrameOptions sets the X-Frame-Options header.
	// Defaults to "DENY". Other common values: "SAMEORIGIN". Set to "" to disable.
	FrameOptions string
	// XSSProtection sets the X-XSS-Protection header.
	// Defaults to "1; mode=block". Set to "" to disable.
	XSSProtection string
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

---

### `StringFlag` {#stringflag}

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

---

#### Methods

### `Apply` {#apply}

```go
func (f *StringFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the string flag with the flag.FlagSet.

---

### `GetAliases` {#getaliases}

```go
func (f *StringFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

---

### `GetName` {#getname}

```go
func (f *StringFlag) GetName() string
```

GetName returns the primary name of the flag.

---

### `IsRequired` {#isrequired}

```go
func (f *StringFlag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

---

### `String` {#string}

```go
func (f *StringFlag) String() string
```

String returns the help text representation of the flag.

---

### `Validate` {#validate}

```go
func (f *StringFlag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided and is not empty.

---

### `StringSliceFlag` {#stringsliceflag}

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

StringSliceFlag defines a flag that accepts multiple string values. The flag can be repeated on the command line (e.g., --tag foo --tag bar).

---

#### Methods

### `Apply` {#apply}

```go
func (f *StringSliceFlag) Apply(set *flag.FlagSet, cli *CLI) error
```

Apply registers the string slice flag with the flag.FlagSet using a custom Value.

---

### `GetAliases` {#getaliases}

```go
func (f *StringSliceFlag) GetAliases() []string
```

GetAliases returns the aliases for the flag.

---

### `GetName` {#getname}

```go
func (f *StringSliceFlag) GetName() string
```

GetName returns the primary name of the flag.

---

### `IsRequired` {#isrequired}

```go
func (f *StringSliceFlag) IsRequired() bool
```

IsRequired returns true if the flag must be provided.

---

### `String` {#string}

```go
func (f *StringSliceFlag) String() string
```

String returns the help text representation of the flag.

---

### `Validate` {#validate}

```go
func (f *StringSliceFlag) Validate(set *flag.FlagSet) error
```

Validate checks if the required flag was provided and resulted in a non-empty slice.

---

### `TimeoutConfig` {#timeoutconfig}

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

---

### `TrailingSlashRedirectConfig` {#trailingslashredirectconfig}

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

---

### `ValidationErrors` {#validationerrors}

```go
type ValidationErrors []error
```

ValidationErrors collects one or more validation violations and implements error.

---

#### Methods

### `Error` {#error}

```go
func (ve ValidationErrors) Error() string
```

Error joins all contained error messages into a single string.

---

