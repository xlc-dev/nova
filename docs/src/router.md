# Router

Nova includes a minimal, zero-dependency HTTP router built upon Go's standard `net/http` package. It simplifies routing by providing:

- **Dynamic Route Parameters:** Define routes with named parameters (e.g., `/users/{id}`).
- **Regex:** Optionally constrain parameters using regular expressions (e.g., `/users/{id:[0-9]+}`).
- **Composable Middleware:** Apply middleware globally, to groups of routes, or individually.
- **Route Grouping:** Organize related routes under a common path prefix and shared middleware.
- **Subrouters:** Mount independent router instances under specific path prefixes for modular applications.
- **Standard `http.Handler` Compatibility:** Integrates seamlessly with standard Go HTTP handlers and middleware.
- **Customizable 404/405 Handlers:** Provide your own handlers for "Not Found" and "Method Not Allowed" responses.
- **Context-Based Parameter Access:** Retrieve URL parameters easily within your handlers via the request context.
- **Enhanced Handler Signature:** Use `HandlerFunc` with a `ResponseContext` for cleaner error handling, response helpers, and automatic data binding/validation.
- **Validation & Localization:** Built-in struct validation with multi-language error messages for form and JSON input.

## Table of Contents

1.  [Getting Started](#getting-started)
2.  [Core Concepts](#core-concepts)
    - [The `Router` Struct](#the-router-struct)
    - [The `Group` Struct](#the-group-struct)
    - [Route Matching](#route-matching)
3.  [Defining Routes](#defining-routes)
    - [Standard Handlers](#standard-handlers)
    - [Enhanced Handlers (`HandlerFunc`)](#enhanced-handlers-handlerfunc)
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
10. [Validation & Localization](#validation--localization)
11. [Full Example](#full-example)

---

## Getting Started

Here's a minimal example to create a router and start an HTTP server:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xlc-dev/nova/nova"
)

func main() {
	router := nova.NewRouter()

	router.GetFunc("/", func(ctx *nova.ResponseContext) error {
		return ctx.Text(http.StatusOK, "Hello, Nova router!")
	})

	cli, err := nova.NewCLI(&nova.CLI{
		Name:    "minimal",
		Version: "0.0.1",
		Action: func(ctx *nova.Context) error {
			return nova.Serve(ctx, router)
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Build and run:

```bash
go build -o app # Or use make!
./app
```

Test the routes:

```bash
curl http://localhost:8080/
# Output: Hello, Nova router!
```

---

## Core Concepts

### The `Router` Struct

The `nova.Router` is the main component. You create one using `nova.NewRouter()`. It implements the standard `http.Handler` interface, so it can be passed directly to `http.ListenAndServe`.

Key responsibilities:

- Registering routes (`Get`, `Post`, `Handle`, etc.).
- Registering enhanced handler routes (`GetFunc`, `PostFunc`, etc.) with `ResponseContext`.
- Applying global middleware (`Use`).
- Managing route groups (`Group`) and subrouters (`Subrouter`).
- Matching incoming requests to registered routes.
- Extracting URL parameters.
- Dispatching requests to the appropriate handler, wrapped in middleware.
- Handling 404 (Not Found) and 405 (Method Not Allowed) errors.

### The `Group` Struct

A `nova.Group` allows you to define a set of routes that share a common URL prefix and/or middleware stack. It's a lightweight helper created via `router.Group("/prefix", optionalMiddleware...)`. Routes added to the group (`group.Get`, `group.Post`, etc.) automatically inherit the prefix and middleware.

### Route Matching

Nova matches routes based on the request path segments:

1.  The incoming request path is split by `/`.
2.  The router iterates through its registered routes.
3.  For each route, it compares the request path segments with the route's pre-compiled segments.
4.  Literal segments must match exactly.
5.  Parameter segments (`{name}`) match any value in that position and capture it.
6.  Regex-constrained parameter segments (`{name:regex}`) must match the provided regular expression.
7.  If a route pattern matches the path:
    - If the HTTP method also matches, the handler (with middleware) is executed. URL parameters are added to the request context.
    - If the HTTP method _doesn't_ match, a 405 Method Not Allowed response is sent (using the custom handler if set).
8.  If no route pattern matches the path, a 404 Not Found response is sent (using the custom handler if set).
9.  Subrouters are checked first based on their prefix before checking the parent router's direct routes.

---

## Defining Routes

Nova supports two handler signatures:

### Standard Handlers

Use the HTTP method helper functions on a `Router` or `Group` instance.

```go
router.Get("/posts", listPostsHandler)
router.Post("/posts", createPostHandler)
router.Put("/posts/{id:[0-9]+}", updatePostHandler)
router.Delete("/posts/{id:[0-9]+}", deletePostHandler)
```

### Enhanced Handlers (`HandlerFunc`)

For cleaner error handling, response helpers, and automatic data binding/validation, use the enhanced handler signature:

```go
router.GetFunc("/hello/{name}", func(ctx *nova.ResponseContext) error {
	name := ctx.URLParam("name")
	return ctx.Text(http.StatusOK, "Hello, "+name+"!")
})
```

You can use `.GetFunc`, `.PostFunc`, `.PutFunc`, `.PatchFunc`, `.DeleteFunc` on both `Router` and `Group`.

---

## Route Parameters

Parameters allow parts of the URL path to be dynamic and captured for use in your handlers.

### Basic Parameters

Define parameters using curly braces: `{name}`.

```go
router.GetFunc("/products/{category}", func(ctx *nova.ResponseContext) error {
	category := ctx.URLParam("category")
	return ctx.Text(http.StatusOK, "Fetching products in category: "+category)
})
```

### Regex Constrained Parameters

Add validation by appending a colon and a Go regular expression within the curly braces: `{name:regex}`.

```go
router.GetFunc("/articles/{id:[0-9]+}", func(ctx *nova.ResponseContext) error {
	id := ctx.URLParam("id")
	return ctx.Text(http.StatusOK, "Fetching article with ID: "+id)
})
```

### Accessing Parameters (`URLParam`)

Inside your handler, use `ctx.URLParam("key")` to retrieve the value of a captured parameter.

---

## Middleware

Middleware provides a way to add cross-cutting concerns (like logging, authentication, compression, CORS) to your request handling pipeline.

- **Type:** `type Middleware func(http.Handler) http.Handler`
- **Functionality:** A middleware function takes an `http.Handler` (the "next" handler in the chain) and returns a new `http.Handler`.

### Global Middleware (`router.Use`)

Middleware added via `router.Use(mws...)` applies to _all_ routes handled by that router instance _and_ any subrouters or groups created from it _after_ the `Use` call.

### Group Middleware (`group.Use`)

Middleware added via `group.Use(mws...)` applies only to routes registered _through that specific group instance_. It runs _after_ any global middleware defined on the parent router.

You can also pass middleware directly when creating the group:

```go
adminGroup := router.Group("/admin", requireAdminMiddleware)
adminGroup.GetFunc("/", adminDashboardHandler)
```

### Execution Order

Middleware execution follows a standard "onion" model:

1.  Request comes in.
2.  Global middleware (added via `router.Use`) executes in the reverse order they were added (last added runs first).
3.  Group middleware (added via `group.Use` or `router.Group`) executes in the reverse order they were added to the group.
4.  The route's specific handler executes.
5.  Response travels back out through the middleware in the forward order they were added.

---

## Route Groups

Groups simplify managing routes with common prefixes and middleware. Create a group using `router.Group(prefix, optionalMiddleware...)`.

```go
v1 := router.Group("/api/v1")
v1.GetFunc("/users", listV1Users)
```

---

## Subrouters

Subrouters allow mounting a completely separate `nova.Router` instance at a specific prefix.

```go
mainRouter := nova.NewRouter()
adminRouter := mainRouter.Subrouter("/admin")
adminRouter.GetFunc("/", adminDashboardHandler)
```

---

## Custom Error Handlers

You can customize the responses for 404 and 405 errors.

### Not Found (404)

```go
func customNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, `{"error": "Resource not found", "path": "`+r.URL.Path+`"}`)
}

router := nova.NewRouter()
router.SetNotFoundHandler(http.HandlerFunc(customNotFound))
```

### Method Not Allowed (405)

```go
func customMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST")
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Method %s not allowed for %s\n", r.Method, r.URL.Path)
}

router := nova.NewRouter()
router.SetMethodNotAllowedHandler(http.HandlerFunc(customMethodNotAllowed))
```

---

## Serving Static Files

Serve static files from an `fs.FS` under a URL prefix:

```go
import "embed"

//go:embed static/*
var staticFS embed.FS

router.Static("/static", staticFS)
```

This will serve files at `/static/*`.

---

## Validation & Localization

Nova provides built-in struct validation for both JSON and form data, with multi-language error messages.

### Binding and Validating Input

Use the enhanced handler signature and `ResponseContext` helpers:

```go
type UserInput struct {
	Email    string `json:"email" minlength:"5" format:"email"`
	Password string `json:"password" format:"password"`
}

router.PostFunc("/signup", func(ctx *nova.ResponseContext) error {
	var input UserInput
	if err := ctx.BindValidated(&input); err != nil {
		return ctx.JSONError(http.StatusBadRequest, err.Error())
	}
	// Use input.Email and input.Password
	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
})
```

- `ctx.Bind(&input)` binds JSON or form data (no validation).
- `ctx.BindValidated(&input)` binds and validates, returning localized error messages based on the `Accept-Language` header.

Supported validation tags include: `required`, `minlength`, `maxlength`, `pattern`, `enum`, `format` (email, url, uuid, date, time, date-time, password, phone, alphanumeric, alpha, numeric), `min`, `max`, `multipleOf`, `minItems`, `maxItems`, `uniqueItems`.

---

## Full Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xlc-dev/nova/nova"
)

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), "requestID", reqID)
		log.Printf("[%s] Incoming request: %s %s", reqID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleHome(ctx *nova.ResponseContext) error {
	return ctx.Text(http.StatusOK, "Welcome Home!")
}

func handleGetUser(ctx *nova.ResponseContext) error {
	userID := ctx.URLParam("id")
	return ctx.Text(http.StatusOK, fmt.Sprintf(
		"Profile for user ID: %s (Request ID: %v)\n",
		userID, ctx.Request().Context().Value("requestID"),
	))
}

func handleGetArticles(ctx *nova.ResponseContext) error {
	return ctx.Text(http.StatusOK, "List of articles")
}

func handleCreateArticle(ctx *nova.ResponseContext) error {
	ctx.Writer().WriteHeader(http.StatusCreated)
	return ctx.Text(http.StatusCreated, "Article created successfully")
}

var router *nova.Router

func main() {
	router = nova.NewRouter()
	router.Use(requestIDMiddleware)

	router.GetFunc("/", handleHome)
	router.GetFunc("/users/{id:[0-9]+}", handleGetUser)

	api := router.Group("/api/v1")
	api.Use(authMiddleware)

	api.GetFunc("/articles", handleGetArticles)
	api.PostFunc("/articles", handleCreateArticle)

	router.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Oops! Route '%s' not found.", r.URL.Path)
	}))

	cli, err := nova.NewCLI(&nova.CLI{
		Name:    "example",
		Version: "0.0.1",
		Action: func(ctx *nova.Context) error {
			return nova.Serve(ctx, router)
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```
