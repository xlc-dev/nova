{{ title: Nova - Router }}

{{ include-block: doc.html markdown="true" }}

# Router

Nova provides a minimal, zero-dependency HTTP router and a programmatic HTML generation engine, built upon Go's standard `net/http` package. It simplifies web development in Go by offering:

- **Dynamic Route Parameters:** Define routes with named parameters (e.g., `/users/{id}`).
- **Regex Validation in Paths:** Optionally constrain parameters using regular expressions (e.g., `/users/{id:[0-9]+}`).
- **Programmatic HTML Generation:** A fluent API for building HTML documents and elements directly in Go code, featuring automatic escaping and comprehensive tag support.
- **Composable Middleware:** Apply middleware globally, to groups of routes, or individually.
- **Route Grouping:** Organize related routes under a common path prefix and shared middleware.
- **Subrouters:** Mount independent router instances under specific path prefixes for modular applications.
- **Standard `http.Handler` Compatibility:** Integrates seamlessly with standard Go HTTP handlers and middleware.
- **Customizable 404/405 Handlers:** Provide your own handlers for "Not Found" and "Method Not Allowed" responses.
- **Context-Based Parameter Access:** Retrieve URL parameters easily within your handlers.
- **Enhanced Handler Signature (`HandlerFunc`):** Use with a `ResponseContext` for cleaner error handling, convenient response helpers, and automatic data binding/validation.
- **Rich `ResponseContext` Helpers:** Methods for sending JSON, programmatically generated HTML, text, and redirect responses.
- **Data Binding:** Automatic binding of JSON and form data to Go structs.
- **Struct Validation & Localization:** Built-in validation for bound data using struct tags, with multi-language error messages (EN, ES, FR, DE, NL supported).
- **Static File Serving:** Serve static files and directories easily.
- **Server Management:** Integrated server utilities for graceful shutdown, live reloading, and configurable logging.

## Table of Contents

1.  [Getting Started](#getting-started)
2.  [Core Concepts](#core-concepts)
    - [The `Router` Struct](#the-router-struct)
    - [The `ResponseContext` Struct](#the-responsecontext-struct)
    - [The `Group` Struct](#the-group-struct)
    - [Route Matching](#route-matching)
3.  [Defining Routes](#defining-routes)
    - [Standard Handlers](#standard-handlers)
    - [Enhanced Handlers (`HandlerFunc`)](#enhanced-handlers-handlerfunc)
    - [Route Options](#route-options)
4.  [Route Parameters](#route-parameters)
    - [Basic Parameters](#basic-parameters)
    - [Regex Constrained Parameters](#regex-constrained-parameters)
    - [Accessing Parameters (`URLParam`)](#accessing-parameters-urlparam)
5.  [Middleware](#middleware)
    - [Global Middleware (`router.Use`)](#global-middleware-routeruse)
    - [Group Middleware (`group.Use`)](#group-middleware-groupuse)
    - [Execution Order](#execution-order)
6.  [Route Groups](#route-groups)
7.  [Subrouters](#subrouters)
8.  [Custom Error Handlers](#custom-error-handlers)
    - [Not Found (404)](#not-found-404)
    - [Method Not Allowed (405)](#method-not-allowed-405)
9.  [Serving Static Files](#serving-static-files)
10. [Programmatic HTML Generation](#programmatic-html-generation)
    - [Overview](#overview)
    - [The `HTMLElement` Interface](#the-htmlelement-interface)
    - [Creating Elements (`nova.Div`, `nova.P`, etc.)](#creating-elements-novadiv-novap-etc)
    - [Adding Content and Children](#adding-content-and-children)
    - [Setting Attributes (`.Attr()`, `.Class()`, `.ID()`)](#setting-attributes-attr-class-id)
    - [Self-Closing Tags](#self-closing-tags)
    - [Text Nodes (`nova.Text()`)](#text-nodes-novatext)
    - [Building Full HTML Documents (`nova.Document()`)](#building-full-html-documents-novadocument)
    - [HTML Escaping](#html-escaping)
11. [Response Helpers (`ResponseContext`)](#response-helpers-responsecontext)
    - [JSON Responses](#json-responses)
    - [HTML Responses (`ctx.HTML()`)](#html-responses-ctxhtml)
    - [Text Responses](#text-responses)
    - [Redirects](#redirects)
    - [Accessing Underlying Writer/Request](#accessing-underlying-writerrequest)
    - [Content Negotiation (`WantsJSON`)](#content-negotiation-wantsjson)
12. [Data Binding and Validation](#data-binding-and-validation)
    - [Binding Request Data](#binding-request-data)
    - [Validating Structs](#validating-structs)
    - [Supported Validation Tags](#supported-validation-tags)
    - [Localization](#localization)
13. [Server Management (`nova.Serve`)](#server-management-novaserve)
14. [Full Example](#full-example)

## 1. Getting Started

Here's a minimal example to create a router, serve a simple HTML page, and start an HTTP server using Nova's integrated server utilities:

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/xlc-dev/nova/nova"
)

func main() {
	router := nova.NewRouter()

	router.GetFunc("/", func(ctx *nova.ResponseContext) error {
		// Use Nova's programmatic HTML generation
		page := nova.Document(
			nova.DocumentConfig{Title: "Hello Nova!"}, // Basic document configuration
			nova.H1().Text("Welcome to the Nova Framework"),
			nova.P().Text("This page is rendered programmatically using Nova's HTML engine."),
		)
		return ctx.HTML(http.StatusOK, page) // Send the HTML page as response
	})

	// Setup CLI for server commands (like --port, --watch, etc.)
	cli, err := nova.NewCLI(&nova.CLI{
		Name:    "minimal-nova-app",
		Version: "0.0.1",
		Action: func(cliCtx *nova.Context) error {
			// nova.Serve handles server lifecycle, logging, and live reload.
			return nova.Serve(cliCtx, router)
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Run the application
	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Build and run:

```bash
go build -o app
./app serve # Or simply ./app if 'serve' is the default action
```

Test the route in your browser or with `curl`:

```bash
curl http://localhost:8080/
# Output will be the rendered HTML page.
```

## Core Concepts

### The `Router` Struct

The `nova.Router` is the main component for defining routes and handling HTTP requests. You create one using `nova.NewRouter()`. It implements the standard `http.Handler` interface, so it can be passed directly to `http.ListenAndServe` or, more commonly, to Nova's `Serve` function for enhanced server management.

Key responsibilities:

- Registering routes with standard or enhanced handlers.
- Applying global middleware.
- Managing route groups and subrouters.
- Matching incoming requests to registered routes.
- Extracting URL parameters.
- Dispatching requests to appropriate handlers.

### The `ResponseContext` Struct

When using enhanced handlers (`HandlerFunc`), your function receives a `*nova.ResponseContext`. This struct wraps the standard `http.ResponseWriter` and `http.Request`, providing a rich set of helper methods for:

- Sending various types of responses (JSON, HTML, text, redirects).
- Accessing URL parameters.
- Binding request data (JSON, form) to Go structs.
- Performing validation on bound data.
- Accessing the underlying `http.Request` and `http.ResponseWriter`.

### The `Group` Struct

A `nova.Group` allows you to define a set of routes that share a common URL prefix and/or a specific stack of middleware. It's a lightweight helper created via `router.Group("/prefix", optionalMiddleware...)`. Routes added to the group automatically inherit its prefix and middleware.

### Route Matching

Nova matches routes based on the request path segments:

1.  The incoming request path is split by `/`.
2.  The router iterates through its registered routes. Subrouters are checked first if their base path (e.g., `/admin`) is a prefix of the request path.
3.  For each route, it compares the request path segments with the route's pre-compiled segments.
4.  Literal segments must match exactly.
5.  Parameter segments (`{name}`) match any value in that position and capture it.
6.  Regex-constrained parameter segments (`{name:regex}`) must match the provided regular expression (the regex is automatically anchored with `^` and `$`).
7.  If a route pattern matches the path:
    - If the HTTP method also matches, the handler (with middleware) is executed. URL parameters are added to the request context.
    - If the HTTP method _doesn't_ match, a 405 Method Not Allowed response is sent (using the custom handler if set).
8.  If no route pattern matches the path, a 404 Not Found response is sent (using the custom handler if set).

## Defining Routes

Nova supports two primary handler signatures for defining routes.

### Standard Handlers

These use the standard Go `http.HandlerFunc` signature: `func(w http.ResponseWriter, r *http.Request)`. You can register them using HTTP method-specific functions on a `Router` or `Group` instance.

```go
func listItemsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "List of items")
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "Item created")
}

router.Get("/items", listItemsHandler)
router.Post("/items", createItemHandler)
```

Other standard handler registration methods include `Put`, `Patch`, `Delete`, and the generic `Handle(method, pattern, handler)`.

### Enhanced Handlers (`HandlerFunc`)

For cleaner error handling, convenient response helpers, and automatic data binding/validation, Nova provides an enhanced handler signature: `func(ctx *nova.ResponseContext) error`.

```go
func greetUserHandler(ctx *nova.ResponseContext) error {
    userName := ctx.URLParam("name")
    if userName == "" {
        return ctx.Text(http.StatusBadRequest, "Name parameter is missing.")
    }
    greeting := fmt.Sprintf("Hello, %s!", userName)
    return ctx.Text(http.StatusOK, greeting)
}

router.GetFunc("/greet/{name}", greetUserHandler)
```

Enhanced handler registration methods include `GetFunc`, `PostFunc`, `PutFunc`, `PatchFunc`, `DeleteFunc`, and the generic `HandleFunc(method, pattern, handler)`. Returning an error from an enhanced handler will typically result in a 500 Internal Server Error response, though this can be customized.

### Route Options

All route registration methods (`Handle`, `HandleFunc`, `Get`, `GetFunc`, etc.) accept optional `*RouteOptions` as the last argument. The `RouteOptions` struct itself is **not defined by Nova**; it's intended for users to define if they need to pass metadata associated with a route, for example, for OpenAPI documentation generation or other custom processing.

```go
// Example: User-defined RouteOptions for OpenAPI
type MyRouteOptions struct {
    Summary     string
    Description string
    Tags        []string
    Deprecated  bool
}

func getUserProfile(ctx *nova.ResponseContext) error { /* ... */ }

routeOpts := &MyRouteOptions{
    Summary: "Get user profile",
    Tags:    []string{"users", "profile"},
}
router.GetFunc("/users/{id}/profile", getUserProfile, routeOpts) // Pass your custom options
```

Nova's router will store this pointer, but it's up to other parts of your application or third-party tools to interpret these options.

## Route Parameters

Parameters allow parts of the URL path to be dynamic and captured for use in your handlers.

### Basic Parameters

Define parameters using curly braces: `{name}`. The value captured for this segment will be available via the parameter name.

```go
router.GetFunc("/store/{category}/items/{itemID}", func(ctx *nova.ResponseContext) error {
    category := ctx.URLParam("category")
    itemID := ctx.URLParam("itemID")
    return ctx.Text(http.StatusOK, fmt.Sprintf("Store Category: %s, Item ID: %s", category, itemID))
})
```

### Regex Constrained Parameters

You can add validation to a parameter by appending a colon and a Go regular expression within the curly braces: `{name:regex}`. The regex is automatically anchored with `^` and `$` by the router. If the path segment does not match the regex, the route will not be considered a match.

```go
// Matches /users/123 but not /users/abc
router.GetFunc("/users/{id:[0-9]+}", func(ctx *nova.ResponseContext) error {
    userID := ctx.URLParam("id") // userID is guaranteed to be a sequence of digits
    return ctx.Text(http.StatusOK, "Fetching user with numeric ID: "+userID)
})

// Matches /files/image.jpg but not /files/document.pdf if you want specific extensions
router.GetFunc("/files/{filename:[a-zA-Z0-9_]+\\.(jpg|png|gif)}", func(ctx *nova.ResponseContext) error {
    filename := ctx.URLParam("filename")
    return ctx.Text(http.StatusOK, "Fetching image file: "+filename)
})
```

### Accessing Parameters (`URLParam`)

- **Inside an enhanced handler (`HandlerFunc`):** Use `ctx.URLParam("key")` to retrieve the value of a captured parameter.

```go
name := ctx.URLParam("name")
```

- **Inside a standard handler (`http.HandlerFunc`):** You can use `router.URLParam(req, "key")`. This requires having access to the `router` instance within the handler.

```go
// Assuming 'router' is accessible, e.g., via a global variable or closure
// userID := router.URLParam(r, "id")
```

Using the `ResponseContext` method within an enhanced handler is generally preferred for cleaner code.

## Middleware

Middleware provides a way to add cross-cutting concerns (like logging, authentication, compression, CORS) to your request handling pipeline.

- **Type:** `type Middleware func(http.Handler) http.Handler`
- **Functionality:** A middleware function takes an `http.Handler` (the "next" handler in the chain) and returns a new `http.Handler`. This returned handler typically performs some action before and/or after calling the `ServeHTTP` method of the "next" handler.

```go
// Example: Logging Middleware
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)

        next.ServeHTTP(w, r) // Call the next handler in the chain

        log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(startTime))
    })
}
```

### Global Middleware (`router.Use`)

Middleware added via `router.Use(mws...)` applies to _all_ routes handled by that router instance _and_ any subrouters or groups created from it _after_ the `Use` call. Middleware functions are applied in the order they are added.

```go
router := nova.NewRouter()
router.Use(loggingMiddleware)       // Applied first
router.Use(panicRecoveryMiddleware) // Applied second (wraps loggingMiddleware's output)
```

### Group Middleware (`group.Use`)

Middleware added via `group.Use(mws...)` applies only to routes registered _through that specific group instance_. It runs _after_ any global middleware defined on the parent router. You can also pass middleware directly when creating the group.

```go
adminGroup := router.Group("/admin")
adminGroup.Use(requireAdminAuthMiddleware) // Applies only to routes in adminGroup

// Alternatively, pass middleware at group creation:
apiGroup := router.Group("/api", apiAuthMiddleware, rateLimitingMiddleware)

apiGroup.GetFunc("/data", func(ctx *nova.ResponseContext) error {
    // This handler will have global, apiAuth, and rateLimiting middleware applied.
    return ctx.JSON(http.StatusOK, map[string]string{"message": "Sensitive data"})
})
```

### Execution Order

Middleware execution follows a standard "onion" model:

1.  Request comes in.
2.  Global middleware (added via `router.Use`) executes. The last one added wraps the ones added before it, so it executes "first" on the way in.
3.  Group middleware (added via `group.Use` or `router.Group`) executes, similarly in a LIFO (Last-In, First-Out) wrapping order for that group.
4.  The route's specific handler executes.
5.  The response travels back out through the middleware in the reverse order of execution on the way in (FIFO relative to addition).

## Route Groups

Groups simplify managing routes that share a common URL prefix and/or a common set of middleware. Create a group using `router.Group(prefix, optionalMiddleware...)`.

```go
router := nova.NewRouter()

// API v1 routes
v1 := router.Group("/api/v1", apiVersionMiddleware("v1"))
v1.GetFunc("/users", listV1UsersHandler)         // Path: /api/v1/users
v1.PostFunc("/products", createV1ProductHandler) // Path: /api/v1/products

// API v2 routes with additional authentication
v2AuthMiddleware := func(next http.Handler) http.Handler { /* ... */ return next }
v2 := router.Group("/api/v2", apiVersionMiddleware("v2"), v2AuthMiddleware)
v2.GetFunc("/users", listV2UsersHandler) // Path: /api/v2/users

// Public website section
web := router.Group("") // No prefix, can be used to apply middleware to a set of top-level routes
web.Use(commonWebMiddleware)
web.GetFunc("/about", aboutPageHandler) // Path: /about
```

Routes defined within a group automatically have the group's prefix prepended to their pattern.

## Subrouters

Subrouters allow you to mount a completely separate `nova.Router` instance at a specific URL prefix. This is useful for modularizing large applications where different sections might have entirely different routing logic, middleware stacks, or even error handlers (though error handlers are inherited by default).

A subrouter created from a parent router inherits:

- The parent's `paramsKey` for context.
- A _clone_ of the parent's global middleware stack _at the time of subrouter creation_. Subsequent changes to the parent's global middleware do not affect already created subrouters.
- The parent's `notFoundHandler` and `methodNotAllowedHandler` by default. These can be overridden on the subrouter.
- The subrouter's `basePath` will be the parent's `basePath` joined with the prefix provided during subrouter creation.

```go
mainRouter := nova.NewRouter()
mainRouter.Use(globalLoggingMiddleware)

// Create a subrouter for an admin section
adminRouter := mainRouter.Subrouter("/admin")
adminRouter.Use(adminAuthenticationMiddleware) // Middleware specific to adminRouter

adminRouter.GetFunc("/dashboard", func(ctx *nova.ResponseContext) error {
    // Path: /admin/dashboard
    // Will have globalLoggingMiddleware and adminAuthenticationMiddleware applied.
    return ctx.HTML(http.StatusOK, nova.Document(nova.DocumentConfig{Title: "Admin Dashboard"}, nova.H1().Text("Admin Panel")))
})

adminRouter.GetFunc("/users", func(ctx *nova.ResponseContext) error {
    // Path: /admin/users
    return ctx.Text(http.StatusOK, "Admin User List")
})

// Another subrouter for a public API
publicApiRouter := mainRouter.Subrouter("/public-api")
publicApiRouter.GetFunc("/version", func(ctx *nova.ResponseContext) error {
    // Path: /public-api/version
    // Will only have globalLoggingMiddleware applied.
    return ctx.JSON(http.StatusOK, map[string]string{"version": "1.0"})
})

// mainRouter.ServeHTTP will delegate to adminRouter or publicApiRouter if the path matches their basePath.
```

When a request comes in, the main router first checks if the request path matches the `basePath` of any of its subrouters. If a match is found, the request is delegated to that subrouter's `ServeHTTP` method.

## Custom Error Handlers

You can customize the responses for 404 (Not Found) and 405 (Method Not Allowed) errors by providing your own `http.Handler`.

### Not Found (404)

This handler is invoked when no route matches the requested URL path.

```go
func customNotFoundHandler(w http.ResponseWriter, r *http.Request) {
    // Using ResponseContext for convenience, even in a standard handler
    rc := nova.ResponseContext{W: w, R: r} // Router field will be nil, but basic W/R ops are fine

    errorPage := nova.Document(
        nova.DocumentConfig{Title: "404 - Page Not Found"},
        nova.H1().Text("Oops! Page Not Found"),
        nova.P(nova.Text(fmt.Sprintf("The page you requested (%s) could not be found.", r.URL.Path))),
        nova.A("/", nova.Text("Go to Homepage")),
    )
    // Manually set status and content type for HTML
    w.Header().Set("Content-Type", "text/html; charset=utf-f-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, errorPage.Render())
}

router := nova.NewRouter()
router.SetNotFoundHandler(http.HandlerFunc(customNotFoundHandler))
```

If using an enhanced handler with `ResponseContext` is preferred for error pages, you'd typically handle this within your application logic or a dedicated error handling middleware if a route isn't found by Nova's router. The `SetNotFoundHandler` expects a standard `http.Handler`.

### Method Not Allowed (405)

This handler is invoked when a route pattern matches the URL path, but not the HTTP method (e.g., a POST request to a GET-only route).

```go
func customMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
    // It's good practice to set the 'Allow' header with permitted methods.
    // This information isn't directly available from Nova's router to this handler by default.
    // You might need to determine allowed methods based on your routing setup if you want to be precise.
    w.Header().Set("Allow", "GET, POST") // Example

    rc := nova.ResponseContext{W: w, R: r}
    rc.JSONError(http.StatusMethodNotAllowed,
        fmt.Sprintf("Method %s is not allowed for the resource %s.", r.Method, r.URL.Path))
}

router := nova.NewRouter()
router.SetMethodNotAllowedHandler(http.HandlerFunc(customMethodNotAllowedHandler))
```

## Serving Static Files

Nova allows you to serve static files (like CSS, JavaScript, images) from an `fs.FS` (such as one created from `embed.FS` or `os.DirFS`) under a specified URL prefix.

The `router.Static(urlPathPrefix string, subFS fs.FS)` method sets up routes for GET and HEAD requests to serve files from `subFS` under the given `urlPathPrefix`.

```go
import (
    "embed"
    "io/fs"
    "log"
    "net/http"

    "github.com/xlc-dev/nova/nova"
)

//go:embed assets/*
var embeddedAssets embed.FS

func main() {
    router := nova.NewRouter()

    // Create a sub-filesystem for the 'public' directory within 'assets'.
    // This is necessary if your files are in a subdirectory of the embedded FS.
    // If 'assets' itself is the root of your static files, you can use 'embeddedAssets' directly.
    publicFilesFS, err := fs.Sub(embeddedAssets, "assets/public")
    if err != nil {
        log.Fatalf("Failed to create sub FS for static files: %v", err)
    }

    // Serve files from the 'publicFilesFS' (i.e., 'assets/public' directory)
    // under the URL prefix '/static'.
    // Example: A file at 'assets/public/css/style.css' will be accessible at '/static/css/style.css'.
    // Example: A file at 'assets/public/images/logo.png' will be accessible at '/static/images/logo.png'.
    router.Static("/static", publicFilesFS)

    // Example route that might use these static files
    router.GetFunc("/", func(ctx *nova.ResponseContext) error {
        page := nova.Document(
            nova.DocumentConfig{
                Title: "Static Files Example",
                HeadExtras: []nova.HTMLElement{
                    nova.StyleSheet("/static/css/main.css"),
                },
            },
            nova.H1().Text("Page with Static Assets"),
            nova.Img("/static/images/banner.jpg", "Banner Image"),
            nova.Script("/static/js/app.js").Attr("defer","true"),
        )
        return ctx.HTML(http.StatusOK, page)
    })

    // rest of your server setup (CLI, nova.Serve)
}
```

This setup uses `http.FileServer` and `http.StripPrefix` internally to serve the files efficiently.

## Programmatic HTML Generation

Nova includes a powerful and fluent API for generating HTML programmatically within your Go code. This allows for type-safe construction of HTML structures without relying on traditional template files.

### Overview

The system is built around the `HTMLElement` interface. You construct HTML by creating `Element` instances (representing tags like `<div>`, `<p>`, etc.) or `textNode` instances, and composing them. Numerous helper functions (e.g., `nova.Div()`, `nova.P()`, `nova.Img()`) simplify element creation.

### The `HTMLElement` Interface

This is the core of the HTML generation system:

```go
type HTMLElement interface {
    Render() string
}
```

Any type that implements this interface can be rendered as part of an HTML structure and can be passed to `ctx.HTML()`. Both `*nova.Element` and `*nova.HTMLDocument` (and `textNode`) implement this interface.

### Creating Elements (`nova.Div`, `nova.P`, etc.)

Nova provides helper functions for most standard HTML tags. These functions return an `*nova.Element`, allowing for method chaining.

```go
import "github.com/xlc-dev/nova/nova"

myDiv := nova.Div()
myParagraph := nova.P()
myImage := nova.Img("path/to/image.jpg", "Alternative text") // src, alt
myLink := nova.A("https://example.com", nova.Text("Click me")) // href, content
```

A comprehensive list of helpers includes `Html`, `Head`, `Body`, `TitleEl`, `Meta`, `LinkTag`, `Script`, `StyleTag`, `H1`-`H6`, `Table`, `Form`, `Input`, `Button`, and many more covering semantic, form, and media elements.

### Adding Content and Children

Elements can have direct text content or child `HTMLElement`s.

- **Direct Text Content:** Use the `.Text(string) *Element` method. This content is HTML-escaped during rendering (except for `<script>` and `<style>` tags).

```go
titleHeader := nova.H1().Text("Page Title")
// Renders: <h1>Page Title</h1>
```

- **Child Elements:** Pass `HTMLElement`s as variadic arguments to the element creation function (e.g., `nova.Div(child1, child2)`) or use the `.Add(children ...HTMLElement) *Element` method.

```go
container := nova.Div(
  nova.H2().Text("Subtitle"),
  nova.P().Text("Some paragraph text."),
)
// Or using .Add():
anotherContainer := nova.Div().Add(
  nova.Span().Text("Part 1"),
  nova.Span().Text("Part 2"),
)
```

### Setting Attributes (`.Attr()`, `.Class()`, `.ID()`)

Attributes are set using fluent methods on an `*nova.Element`:

- `.Attr(key, value string) *Element`: Sets a generic attribute.
- `.Class(class string) *Element`: Sets the `class` attribute.
- `.ID(id string) *Element`: Sets the `id` attribute.
- `.Style(style string) *Element`: Sets the `style` attribute.
- `.BoolAttr(key string, present bool) *Element`: Adds a boolean attribute (e.g., `<input disabled>`) if `present` is true. The attribute is rendered as `key="key"`. If `present` is false, the attribute is removed if it exists.

```go
styledDiv := nova.Div().
    ID("main-content").
    Class("container theme-dark").
    Style("border: 1px solid red; padding: 10px;").
    Attr("data-custom", "my-value")

inputField := nova.Input("text").BoolAttr("required", true).Attr("placeholder", "Enter name")
// Renders: <input type="text" required="required" placeholder="Enter name" />
```

### Self-Closing Tags

Some elements, like `<img>`, `<br>`, `<hr>`, `<input>`, `<meta>`, and `<link>`, are self-closing. The helper functions for these tags create elements that render correctly with `/>`.

```go
separator := nova.Hr() // Renders: <hr />
iconLink := nova.LinkTag().Attr("rel", "icon").Attr("href", "/favicon.ico") // Renders: <link rel="icon" href="/favicon.ico" />
charsetMeta := nova.MetaCharset("UTF-16") // Renders: <meta charset="UTF-16" />
```

### Text Nodes (`nova.Text()`)

For adding plain text that should be HTML-escaped when used as a child of another element, use `nova.Text(content string) HTMLElement`. This returns a `textNode` which implements `HTMLElement`.

```go
paragraphWithMixedContent := nova.P(
    nova.Text("This is "),
    nova.Strong().Text("bold"),
    nova.Text(" text, and this is "),
    nova.Em().Text("italic"),
    nova.Text("."),
)
// Renders: <p>This is <strong>bold</strong> text, and this is <em>italic</em>.</p>
```

### Building Full HTML Documents (`nova.Document()`)

To create a complete HTML5 page, use the `nova.Document(config DocumentConfig, bodyContent ...HTMLElement) *HTMLDocument` function. It takes a `nova.DocumentConfig` for head customizations and variadic `HTMLElement`s for the body content. The returned `*HTMLDocument` also implements `HTMLElement`.

```go
// In your handler:
page := nova.Document(
    nova.DocumentConfig{
        Lang:        "en-US",
        Title:       "My Awesome App",
        Description: "A fantastic application built with Nova.",
        Keywords:    "nova, web, go, framework",
        Author:      "Awesome Dev",
        HeadExtras: []nova.HTMLElement{ // Add custom links, scripts, meta tags to head
            nova.StyleSheet("/css/theme.css"),
            nova.Script("/js/app.js").Attr("defer", "defer"), // Note: Script tag is not self-closing
            nova.MetaNameContent("robots", "index, follow"),
        },
    },
    // Body content starts here
    nova.Header(
        nova.H1().Text("Welcome to My Awesome App!"),
    ),
    nova.Main(
        nova.P().Text("This is the main content of the page, demonstrating a full document structure."),
        nova.Img("/images/logo.png", "App Logo").Class("logo"),
    ),
    nova.Footer(
        nova.P().Text("© 2025 My Company"),
    ),
)
// err := ctx.HTML(http.StatusOK, page)
```

`DocumentConfig` allows customization of:

- `Lang`: `<html>` lang attribute (default "en").
- `Title`: `<title>` tag content (default "Document").
- `Charset`: `<meta charset>` (default "utf-8").
- `Viewport`: `<meta name="viewport">` (default "width=device-width, initial-scale=1").
- `Description`, `Keywords`, `Author`: For respective meta tags (omitted if empty).
- `HeadExtras`: Slice of `HTMLElement`s to add to the `<head>` section (e.g., additional stylesheets, scripts, meta tags).

The `nova.Document()` function constructs the `<!DOCTYPE html>`, `<html>`, `<head>` (with specified meta tags and title), and `<body>` structure for you.

### HTML Escaping

By default, text content set via `.Text()` or `nova.Text()` and attribute values are HTML-escaped using `html.EscapeString` to prevent XSS vulnerabilities.

**Exception:** The direct content of `<script>` and `<style>` elements (set via `.Text()` or passed to `InlineScript()` / `StyleTag()`) is **not** escaped. This allows you to embed raw JavaScript and CSS directly.

```go
inlineJS := nova.InlineScript("if (a < b && c > d) { console.log('Condition met'); }")
// Renders: <script>if (a < b && c > d) { console.log('Condition met'); }</script> (content is raw)

inlineCSS := nova.StyleTag("body > p { font-weight: bold; color: #333; }")
// Renders: <style>body > p { font-weight: bold; color: #333; }</style> (content is raw)
```

Be cautious when embedding user-provided data into inline scripts or styles.

## Response Helpers (`ResponseContext`)

These methods are available on `*nova.ResponseContext` within `HandlerFunc` (enhanced handlers).

### JSON Responses

- `JSON(statusCode int, data any) error`: Encodes `data` to JSON and sends it. Sets `Content-Type: application/json`.
- `JSONError(statusCode int, message string) error`: Sends a standardized JSON error: `{"error": "message"}`.

```go
type Item struct { ID string `json:"id"`; Name string `json:"name"` }
func getItemHandler(ctx *nova.ResponseContext) error {
    item := Item{ID: "item123", Name: "Example Item"}
    return ctx.JSON(http.StatusOK, item)
}
```

### HTML Responses (`ctx.HTML()`)

- `HTML(statusCode int, content HTMLElement) error`: Sends an HTML response. Sets `Content-Type: text/html; charset=utf-8`.
  The `content` argument must be a type that implements the `nova.HTMLElement` interface (like `*nova.Element`, `*nova.HTMLDocument`, or a custom type that has a `Render() string` method).

```go
// In a handler:
func serveSimpleHTMLPage(ctx *nova.ResponseContext) error {
    // Using Nova's programmatic HTML builders
    myPageContent := nova.Div().
        Class("container").
        Add(
            nova.H1().Text("Hello from Nova!"),
            nova.P().Text("This HTML was generated in Go using Nova's fluent API."),
            nova.A("https://example.com", nova.Text("Learn more about Nova")),
        )

    // If you need a full document structure:
    fullDoc := nova.Document(
        nova.DocumentConfig{Title: "My Nova Page"},
        myPageContent, // Add the div as body content
    )

    return ctx.HTML(http.StatusOK, fullDoc) // Send the full document
}
```

### Text Responses

- `Text(statusCode int, text string) error`: Sends a plain text response. Sets `Content-Type: text/plain; charset=utf-8`.

```go
func healthCheckHandler(ctx *nova.ResponseContext) error {
    return ctx.Text(http.StatusOK, "Server is healthy and running.")
}
```

### Redirects

- `Redirect(statusCode int, url string) error`: Sends an HTTP redirect. Sets the `Location` header and writes the status code.
  Common status codes: `http.StatusMovedPermanently` (301), `http.StatusFound` (302), `http.StatusTemporaryRedirect` (307).

```go
func oldPathHandler(ctx *nova.ResponseContext) error {
    // Permanent redirect from an old path to a new one
    return ctx.Redirect(http.StatusMovedPermanently, "/new-path-for-this-resource")
}

func loginRequiredHandler(ctx *nova.ResponseContext) error {
    // Temporary redirect to login page if user is not authenticated
    return ctx.Redirect(http.StatusFound, "/login?returnTo="+ctx.Request().URL.Path)
}
```

### Accessing Underlying Writer/Request

- `Request() *http.Request`: Returns the underlying `*http.Request` instance.
- `Writer() http.ResponseWriter`: Returns the underlying `http.ResponseWriter` instance.
  These are useful for advanced scenarios or when integrating with third-party libraries that expect these standard Go types.

```go
func advancedHandler(ctx *nova.ResponseContext) error {
    // Access raw request for specific header not covered by helpers
    apiKey := ctx.Request().Header.Get("X-API-Key")
    if apiKey == "" {
        return ctx.JSONError(http.StatusUnauthorized, "API Key required")
    }

    // Use raw writer for something like streaming a response
    // ctx.Writer().Header().Set("Content-Type", "application/octet-stream")
    // flusher, ok := ctx.Writer().(http.Flusher)
    // if !ok { /* handle error */ }
    // ... stream data ...
    return nil // Or an error if streaming fails
}
```

### Content Negotiation (`WantsJSON`)

- `WantsJSON() bool`: Returns `true` if the request's `Content-Type` or `Accept` header suggests a preference for JSON (i.e., contains "application/json"). This is useful for creating endpoints that can serve multiple content types based on client preference.

```go
func versatileDataHandler(ctx *nova.ResponseContext) error {
    data := map[string]string{"id": "data123", "value": "Some important information"}

    if ctx.WantsJSON() {
        return ctx.JSON(http.StatusOK, data)
    }

    // Fallback to HTML representation
    htmlContent := nova.Document(
        nova.DocumentConfig{Title: "Data Details"},
        nova.H1().Text("Data Item: "+data["id"]),
        nova.Dl( // Definition List example
            nova.Dt().Text("ID"),
            nova.Dd().Text(data["id"]),
            nova.Dt().Text("Value"),
            nova.Dd().Text(data["value"]),
        ),
    )
    return ctx.HTML(http.StatusOK, htmlContent)
}
```

## Data Binding and Validation

Nova simplifies handling incoming request data (JSON, forms) and validating it using struct tags.

### Binding Request Data

Use `ResponseContext` methods within an enhanced handler (`HandlerFunc`):

- `ctx.Bind(v any) error`: Automatically detects `Content-Type` (supports "application/json" and URL-encoded form data) and binds the request data to the provided struct `v` (which must be a pointer).
- `ctx.BindJSON(v any) error`: Specifically for binding JSON request bodies to `v`.
- `ctx.BindForm(v any) error`: Specifically for parsing URL-encoded form data (from `r.Form`) and binding it to `v`.
  - For form binding, it supports `string`, `bool` (recognizes "on", "true", "1" as true, otherwise false if field is bool), numeric types (`int*`, `uint*`, `float*`), and `[]string` (from comma-separated values in a single field or multiple form fields with the same name).
  - Field name matching uses `json` struct tags (e.g., `json:"user_name"`), falling back to the struct field names if no `json` tag is present or if the tag is `"-"`.

```go
type CreateUserInput struct {
    Username string   `json:"username"` // Used for both JSON and form field name matching
    Email    string   `json:"email"`
    Age      int      `json:"age,omitempty"`
    Tags     []string `json:"tags"` // For forms, can be 'tag1,tag2' or multiple 'tags=tag1&tags=tag2'
}

func handleCreateUser(ctx *nova.ResponseContext) error {
    var input CreateUserInput
    if err := ctx.Bind(&input); err != nil { // Binds JSON or Form data
        return ctx.JSONError(http.StatusBadRequest, "Invalid input data: "+err.Error())
    }
    // Process input (e.g., save to database)
    log.Printf("User data received: %+v", input)
    return ctx.JSON(http.StatusCreated, input)
}
```

### Validating Structs

- `ctx.BindValidated(v any) error`: This is the most convenient method. It first binds the request data (JSON or form) to the struct `v` (pointer) and then validates `v` using rules defined in struct tags.
- If validation fails, it returns a `nova.ValidationErrors` type (which is `[]error`). The `Error()` method of `ValidationErrors` returns a semicolon-separated string of all validation messages.
- The `error:"custom message"` tag on a struct field allows you to specify a custom error message for _any_ validation failure on that field, overriding the default localized message for that specific field's validation.
- If a field in a struct is intended to be **required**, ensure its `json` tag does **not** include `omitempty`. If such a field is its zero value after binding (e.g., empty string for `string`, 0 for `int`), it will trigger a "required" validation error.

```go
type SignupInput struct {
    FullName string `json:"fullName" minlength:"2" maxlength:"50"`
    Email    string `json:"email" format:"email"` // Implicitly required if no 'omitempty'
    Age      int    `json:"age,omitempty" min:"18" max:"120"` // Optional due to omitempty
    Password string `json:"password" format:"password" error:"Your password is too weak, please choose a stronger one."`
}

func handleSignup(ctx *nova.ResponseContext) error {
    var input SignupInput
    if err := ctx.BindValidated(&input); err != nil {
        // err will be of type nova.ValidationErrors
        return ctx.JSONError(http.StatusBadRequest, err.Error())
    }
    // Process valid input (e.g., create user account)
    log.Printf("Signup successful for: %+v", input)
    return ctx.JSON(http.StatusOK, map[string]string{"message": "Signup successful!"})
}
```

Validation is recursive: if a struct field is itself a struct or a pointer to a struct, `validateStruct` will be called on it. For slices of structs, each element in the slice is validated.

### Supported Validation Tags

- `required`: (Implicit) If a field's `json` tag does not contain `omitempty` and the field has its zero value after binding.
- `minlength:"<value>"`: Minimum string length (for `string` type).
- `maxlength:"<value>"`: Maximum string length (for `string` type).
- `min:"<value>"`: Minimum numeric value (for `int*`, `uint*`, `float*` types).
- `max:"<value>"`: Maximum numeric value (for `int*`, `uint*`, `float*` types).
- `pattern:"<regex>"`: String must match the Go regular expression (for `string` type).
- `enum:"<val1>|<val2>|..."`: String must be one of the specified pipe-separated values (for `string` type).
- `format:"<type>"`: Predefined format validation for strings. Supported types:
  - `email`: Validates using `mail.ParseAddress`.
  - `url`: Validates using `url.ParseRequestURI`, ensuring scheme and host are present.
  - `uuid`: Validates against RFC-4122 v4 UUID format.
  - `date-time`: Validates against `time.RFC3339` format (e.g., `2023-10-26T10:00:00Z`).
  - `date`: Validates against `YYYY-MM-DD` format (e.g., `2023-10-26`).
  - `time`: Validates against `HH:MM:SS` format (e.g., `10:00:00`).
  - `password`: Basic check, ensures string length is at least 8 characters.
  - `phone`: Basic international phone number pattern `^[\+]?[1-9][\d\s\-\(\)]{7,15}$`.
  - `alphanumeric`: Contains only A-Z, a-z, 0-9.
  - `alpha`: Contains only A-Z, a-z.
  - `numeric`: Contains only 0-9.
- `multipleOf:"<value>"`: Number must be a multiple of the value (for numeric types).
- `minItems:"<value>"`: Minimum number of items in a slice/array.
- `maxItems:"<value>"`: Maximum number of items in a slice/array.
- `uniqueItems:"true"`: All items in a slice/array must be unique. (Compares underlying values).
- `error:"<custom_message>"`: Overrides the default/localized validation error message for _any_ validation rule that fails on this specific field.

### Localization

Validation error messages are automatically localized based on the `Accept-Language` HTTP header in the request.

- Supported languages (defined in `validationMessages` map): English (`en`), Spanish (`es`), French (`fr`), German (`de`), Dutch (`nl`).
- English (`en`) serves as the fallback if a requested language or a specific message key within a language is not found.
- The `detectLanguage` internal function parses the `Accept-Language` header (simplified parsing) to pick the best available match from the supported languages.

## Server Management (`nova.Serve`)

Nova provides a `Serve(ctx *nova.Context, router http.Handler) error` function to simplify server startup, management, and add features like live reloading and configurable logging. It's typically used as the action for a `nova.CLI` command.

The `*nova.Context` argument provides access to parsed command-line flags and configuration.

Key features of `nova.Serve`:

- **Server Startup:** Listens on a host and port. These can be configured via `*nova.Context` (typically from CLI flags).
  - `host` (string, default: "localhost", e.g., `--host 0.0.0.0`)
  - `port` (int, default: 8080, e.g., `--port 3000`)
- **Graceful Shutdown:** Catches `SIGINT` (Ctrl+C) and `SIGTERM` signals to shut down the HTTP server gracefully within a 5-second timeout.
- **Live Reloading (Hot Reload):**
  - Enable by setting `watch: true` in the `*nova.Context` (e.g., via a `--watch` CLI flag).
  - Monitors file changes in the current working directory (or a specified directory).
  - Watched file extensions are configurable via an `extensions` string in `*nova.Context` (comma-separated, default: ".go", e.g., `--extensions .go,.html,.css`).
  - When a change in a watched file is detected, `nova.Serve` rebuilds the application (`go build -o <current_executable_path> .`) and then re-executes the new binary, effectively restarting the server with the new code.
- **Configurable Logging (slog):**
  - Log output format can be set via `log_format` in `*nova.Context` (string, "text" or "json", default: "text", e.g., `--log-format json`).
  - Log level can be set via `log_level` in `*nova.Context` (string, "debug", "info", "warn"/"warning", "error", default: "info", e.g., `--log-level debug`).
  - This configures the global `slog.Default` logger.
- **Verbose Mode:**
  - Setting `verbose: true` in `*nova.Context` (e.g., via a `--verbose` CLI flag) enables more detailed logging from the `Serve` function itself, such as file watcher activity and recompilation commands.

```go
// Example usage within main.go
// Assuming cli is a *nova.CLI instance
cli.Action = func(cliCtx *nova.Context) error {
    router := setupMyApplicationRouter() // Your function to get the configured router
    slog.Info("Application starting...") // Uses slog, configured by nova.Serve
    return nova.Serve(cliCtx, router)
}
```

## Full Example

This example demonstrates routing, middleware, programmatic HTML generation, static file serving, data binding with validation, and server management via `nova.Serve` and `nova.CLI`.

```go
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/xlc-dev/nova/nova"
)

//go:embed static/*
var embeddedStaticFiles embed.FS

// Middleware
func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Use slog for structured logging, which nova.Serve can configure
		slog.Info("Request received", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
		next.ServeHTTP(w, r)
		slog.Info("Request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

func simpleAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth-Token") != "secret-token" {
			// Use ResponseContext for consistent error responses
			rc := nova.ResponseContext{W: w, R: r}
			rc.JSONError(http.StatusUnauthorized, "Unauthorized: Missing or invalid X-Auth-Token")
			return
		}
		ctx := context.WithValue(r.Context(), "user", "AuthenticatedUser")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTML Page Components (Example)
func pageLayout(config nova.DocumentConfig, bodyContent ...nova.HTMLElement) nova.HTMLElement {
	defaultConfig := nova.DocumentConfig{
		Lang:     "en",
		Title:    "Nova App",
		Charset:  "UTF-8",
		Viewport: "width=device-width, initial-scale=1.0",
		HeadExtras: []nova.HTMLElement{
			nova.StyleSheet("/assets/css/style.css"), // Path to static CSS
			nova.Script("/assets/js/main.js").Attr("defer", "true"),
		},
	}
	// Override defaults with provided config
	if config.Title != "" {
		defaultConfig.Title = config.Title
	}
	if config.Lang != "" {
		defaultConfig.Lang = config.Lang
	}
	// Append, don't overwrite, HeadExtras
	defaultConfig.HeadExtras = append(defaultConfig.HeadExtras, config.HeadExtras...)

	allBodyContent := []nova.HTMLElement{
		nova.Header(
			nova.Nav(nova.A("/", nova.Text("Home")), nova.Text(" | "), nova.A("/contact", nova.Text("Contact"))).Class("main-nav"),
		).Class("site-header"),
	}
	allBodyContent = append(allBodyContent, bodyContent...)
	allBodyContent = append(allBodyContent,
		nova.Footer(nova.P().Text(fmt.Sprintf("© %d Nova Example Inc.", time.Now().Year()))).Class("site-footer"),
	)

	return nova.Document(defaultConfig, allBodyContent...)
}

// Handlers
func handleHomepage(ctx *nova.ResponseContext) error {
	slog.Info("Handling homepage request")
	content := nova.Main(
		nova.H1().Text("Welcome to the Nova Framework Showcase!"),
		nova.P().Text("This page is dynamically rendered using Nova's programmatic HTML engine."),
		nova.Img("/assets/images/banner.png", "Nova Banner").Style("max-width:100%; height:auto;"),
	)
	return ctx.HTML(http.StatusOK, pageLayout(nova.DocumentConfig{Title: "Homepage"}, content))
}

func handleContactPage(ctx *nova.ResponseContext) error {
	slog.Info("Handling contact page request")
	contactForm := nova.Form(
		nova.H2().Text("Contact Us"),
		nova.P(
			nova.Label(nova.Text("Your Name: ")).Attr("for", "name"),
			nova.TextInput("name").ID("name").Attr("required", "true"),
		),
		nova.P(
			nova.Label(nova.Text("Your Email: ")).Attr("for", "email"),
			nova.EmailInput("email").ID("email").Attr("required", "true"),
		),
		nova.P(
			nova.Label(nova.Text("Message: ")).Attr("for", "message"),
			nova.Textarea().ID("message").Attr("name", "message").Attr("rows", "5").Attr("required", "true"),
		),
		nova.SubmitButton("Send Message"),
	).Attr("method", "POST").Attr("action", "/contact-submit")

	content := nova.Main(contactForm)
	return ctx.HTML(http.StatusOK, pageLayout(nova.DocumentConfig{Title: "Contact Us"}, content))
}

type ContactFormInput struct {
	Name    string `json:"name" minlength:"2" error:"Name must be at least 2 characters."`
	Email   string `json:"email" format:"email"`
	Message string `json:"message" minlength:"10" error:"Message is too short."`
}

func handleContactSubmit(ctx *nova.ResponseContext) error {
	var input ContactFormInput
	if err := ctx.BindValidated(&input); err != nil {
		slog.Warn("Contact form validation failed", "errors", err.Error())
		// For a real app, you'd re-render the form with errors.
		// Here, we'll just return a JSON error for simplicity.
		return ctx.JSONError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	slog.Info("Contact form submitted", "name", input.Name, "email", input.Email)

	thankYouMessage := nova.Main(
		nova.H1().Text("Thank You!"),
		nova.P().Text("Your message has been received. We will get back to you shortly."),
		nova.A("/", nova.Text("Return to Homepage")),
	)
	return ctx.HTML(http.StatusOK, pageLayout(nova.DocumentConfig{Title: "Message Sent"}, thankYouMessage))
}

func handleApiGetData(ctx *nova.ResponseContext) error {
	user, _ := ctx.Request().Context().Value("user").(string)
	slog.Info("API: GetData requested", "authenticated_user", user)
	data := map[string]any{
		"message":     "This is protected data from the API.",
		"timestamp":   time.Now().Format(time.RFC3339),
		"currentUser": user,
	}
	return ctx.JSON(http.StatusOK, data)
}

func main() {
	// Router Setup
	router := nova.NewRouter()

	// Global Middleware
	router.Use(requestLoggingMiddleware)

	// Serve static files from embedded 'static' directory under '/assets' URL path,
    // e.g., /assets/css/style.css
	staticFilesRoot, err := fs.Sub(embeddedStaticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub FS for static files: %v", err)
	}
	router.Static("/assets", staticFilesRoot)

	// Public Routes
	router.GetFunc("/", handleHomepage)
	router.GetFunc("/contact", handleContactPage)
	router.PostFunc("/contact-submit", handleContactSubmit)

	// API Routes with Authentication
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(simpleAuthMiddleware)
	apiGroup.GetFunc("/data", handleApiGetData)

	// CLI Setup
	cliApp, err := nova.NewCLI(&nova.CLI{
		Name:        "NovaFullApp",
		Version:     "1.0.0",
		Description: "A full example application using the Nova framework.",
		Action: func(cliCtx *nova.Context) error {
			// This action is called when 'myfullapp serve' (or just 'myfullapp') is run.
			// nova.Serve will use cliCtx for port, host, watch, log settings.
			slog.Info("Starting Nova application server...", "version", cliCtx.App.Version)
			return nova.Serve(cliCtx, router)
		},
		// Add more commands here if needed
	})

	if err != nil {
		log.Fatalf("Failed to initialize CLI: %v", err)
	}

	// Run the CLI application
	if err := cliApp.Run(os.Args); err != nil {
		slog.Error("Application exited with error", "error", err)
		os.Exit(1)
	}
}
```

{{ endinclude }}
