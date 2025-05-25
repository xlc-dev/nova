# Middleware

Nova has ready-to-use middleware for `net/http`. In this document, you will find the list of built-in middleware and how to use them.

## Table of Contents

1.  [What is Middleware?](#what-is-middleware)
2.  [Custom Middleware](#custom-middleware)
3.  [Built-in Middleware](#built-in-middleware)
    - [LoggingMiddleware](#loggingmiddleware)
    - [RecoveryMiddleware](#recoverymiddleware)
    - [RequestIDMiddleware](#requestidmiddleware)
    - [CORSMiddleware](#corsmiddleware)
    - [SecurityHeadersMiddleware](#securityheadersmiddleware)
    - [TimeoutMiddleware](#timeoutmiddleware)
    - [BasicAuthMiddleware](#basicauthmiddleware)
    - [MethodOverrideMiddleware](#methodoverridemiddleware)
    - [EnforceContentTypeMiddleware](#enforcecontenttypemiddleware)
    - [CacheControlMiddleware](#cachecontrolmiddleware)
    - [GzipMiddleware](#gzipmiddleware)
    - [CSRFMiddleware](#csrfmiddleware)
    - [ETagMiddleware](#ETagMiddleware)
    - [HealthCheckMiddleware](#healthcheckmiddleware)
    - [RealIPMiddleware](#realipmiddleware)
    - [MaxRequestBodySizeMiddleware](#maxrequestbodysizemiddleware)
    - [TrailingSlashRedirectMiddleware](#trailingslashredirectmiddleware)
    - [ForceHTTPSMiddleware](#forcehttpsmiddleware)
    - [ConcurrencyLimiterMiddleware](#concurrencylimitermiddleware)
    - [MaintenanceModeMiddleware](#maintenancemodemiddleware)
    - [IPFilterMiddleware](#ipfiltermiddleware)
    - [RateLimitMiddleware](#ratelimitmiddleware)

## What is Middleware?

In the context of Nova's web handling capabilities (using Go's standard `net/http` package), **middleware** refers to a function that wraps an `http.Handler`. It follows the standard Go pattern: a function that takes an `http.Handler` and returns a new `http.Handler`.

```go
// Middleware defines the function signature for middleware.
// A middleware is a function that wraps an http.Handler, adding extra behavior.
type Middleware func(http.Handler) http.Handler
```

Middleware functions sit between the server's routing logic (like `nova.Router`) and your final request handler (`http.HandlerFunc`). They provide a way to process the `http.Request` and `http.ResponseWriter` or perform actions _before_ or _after_ your main handler logic runs.

Think of it like layers processing an incoming HTTP request:

1.  An HTTP request arrives.
2.  Nova's router directs the request towards the appropriate handler.
3.  If middleware is applied (e.g., via `router.Use(...)`), the request passes through each middleware function in sequence.
4.  Each middleware can:
    - Examine or modify the `http.Request` (e.g., add context values, parse headers).
    - Wrap the `http.ResponseWriter` to intercept or modify the response (e.g., capture status code, compress data).
    - Perform tasks like logging, timing, authentication checks, authorization, rate limiting, or setting common headers.
    - Decide whether to pass control to the `next` handler in the chain by calling `next.ServeHTTP(w, r)`.
    - Perform tasks _after_ the `next` handler has completed (e.g., logging the response status, cleanup).
5.  Finally, the core `http.HandlerFunc` for the route is executed (if the middleware chain allowed it).

**Key Benefits of Using Middleware:**

- **Separation of Concerns:** Keeps cross-cutting logic (like logging, authentication, compression, security headers) separate from your core request handling logic, making handlers cleaner and more focused.
- **Reusability:** Write common pre-processing or post-processing logic once as middleware and apply it to multiple routes or groups of routes easily using `router.Use(...)` or `group.Use(...)`.
- **Composability:** Chain multiple, small middleware functions together to build complex request processing pipelines in a modular and maintainable way.

This pattern is fundamental to building web applications and APIs in Go.

## Custom Middleware

You can easily create your own middleware by defining a function that matches the `nova.Middleware` type signature:

```go
type Middleware func(http.Handler) http.Handler
```

Here's a simple example of a custom middleware that adds a custom header to every response:

```go
// CustomHeaderMiddleware adds a "X-Custom-Header" to responses.
func CustomHeaderMiddleware(headerName, headerValue string) nova.Middleware {
    // Return the actual middleware function
    return func(next http.Handler) http.Handler {
        // Return the HandlerFunc that wraps the next handler
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Set the header before calling the next handler
            w.Header().Set(headerName, headerValue)

            // Call the next middleware or handler in the chain
            next.ServeHTTP(w, r)

            // Could also perform actions *after* the handler runs here
        })
    }
}

// And then somewhere else where you want to use it on a router/group/subrouter:
router.Use(CustomHeaderMiddleware("X-App-Version", "1.2.3"))
```

## Built-in Middleware

Nova provides a collection of standard `net/http` middleware. Each middleware is typically configured using a specific `Config` struct and applied using `router.Use(...)` for global application or `group.Use(...)` for group-specific application.

---

### LoggingMiddleware

- **Description:** Logs incoming requests (start) and outgoing responses (completion), including method, path, remote address, status code, response size, and duration.
- **Configuration:** `nova.LoggingConfig`
  - `Logger *log.Logger`: Logger instance (defaults to `log.Default()`).
  - `LogRequestID bool`: Include request ID in logs (defaults to true). Requires `RequestIDMiddleware`.
  - `RequestIDKey contextKey`: Context key for request ID (defaults to internal key).

#### Example

```go
func main() {
	router := nova.NewRouter()
	customLogger := log.New(os.Stdout, "[ACCESS] ", log.LstdFlags)

	// Apply logging middleware
	router.Use(nova.LoggingMiddleware(&nova.LoggingConfig{
		Logger:       customLogger,
		LogRequestID: true, // Optional: Explicitly true (default)
	}))
	// Apply RequestID middleware if LogRequestID is true
	router.Use(nova.RequestIDMiddleware(nil)) // Use defaults
}
```

---

### RecoveryMiddleware

- **Description:** Recovers from panics in downstream handlers/middleware, logs the panic, and sends a 500 Internal Server Error response (or calls a custom handler).
- **Configuration:** `nova.RecoveryConfig`
  - `Logger *log.Logger`: Logger for panic messages (defaults to `log.Default()`).
  - `LogRequestID bool`: Include request ID in panic logs (defaults to true).
  - `RequestIDKey contextKey`: Context key for request ID (defaults to internal key).
  - `RecoveryHandler func(http.ResponseWriter, *http.Request, interface{})`: Custom function to handle recovered panics.

#### Example

```go
func customPanicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	reqID := nova.GetRequestID(r.Context()) // Get request ID if available
	log.Printf("[PANIC RECOVERY][%s] Custom handler called: %v", reqID, err)
	http.Error(w, "Something went terribly wrong!", http.StatusInternalServerError)
}

func main() {
	router := nova.NewRouter()

	// Apply RequestID first so Recovery can log it
	router.Use(nova.RequestIDMiddleware(nil))

	// Apply recovery middleware (must be early in the chain, so it can set the recover mechanism)
	router.Use(nova.RecoveryMiddleware(&nova.RecoveryConfig{
		RecoveryHandler: customPanicHandler,
		LogRequestID:    true,
	}))

	router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("Simulating a handler panic!")
	})

	router.Get("/safe", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is fine."))
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
```

---

### RequestIDMiddleware

- **Description:** Assigns a unique ID to each request (from header or generated), sets it in the response header, and adds it to the request context.
- **Configuration:** `nova.RequestIDConfig`
  - `HeaderName string`: Header to check/set (defaults to `X-Request-ID`).
  - `ContextKey contextKey`: Context key for storage (defaults to internal key).
  - `Generator func() string`: ID generation function (defaults to timestamp-based).
- **Context Helper:** `nova.GetRequestID(ctx context.Context)` retrieves the ID.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply RequestID middleware globally
	router.Use(nova.RequestIDMiddleware(&nova.RequestIDConfig{
		HeaderName: "X-Trace-ID", // Custom header name
		Generator: func() string { // Custom generator (e.g., UUID)
			return uuid.NewString() // From github.com/google/uuid
		},
	}))

	router.Get("/whoami", func(w http.ResponseWriter, r *http.Request) {
		reqID := nova.GetRequestID(r.Context())
		fmt.Fprintf(w, "Your request ID is: %s", reqID)
	})
}
```

---

### CORSMiddleware

- **Description:** Handles Cross-Origin Resource Sharing (CORS) by setting `Access-Control-*` headers and managing preflight `OPTIONS` requests.
- **Configuration:** `nova.CORSConfig`
  - `AllowedOrigins []string`: Allowed origin domains (e.g., `["http://localhost:3000"]`, `["*"]`). Defaults to none.
  - `AllowedMethods []string`: Allowed HTTP methods (defaults to GET, POST, PUT, DELETE, PATCH, OPTIONS).
  - `AllowedHeaders []string`: Allowed request headers (defaults to Content-Type, Authorization, X-Request-ID). `"*"` allows any.
  - `ExposedHeaders []string`: Response headers accessible to the client script. Defaults to none.
  - `AllowCredentials bool`: Allow cookies/auth headers with requests (defaults to false). Cannot be true if `AllowedOrigins` is `["*"]`.
  - `MaxAgeSeconds int`: Cache duration for preflight results (defaults to 86400).

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply CORS middleware globally
	// OPTIONS requests are handled automatically by the middleware
	router.Use(nova.CORSMiddleware(nova.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000", "https://my-frontend.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Custom-Response-Header"},
		AllowCredentials: true,
		MaxAgeSeconds:    3600, // 1 hour
	}))

	router.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Response-Header", "SomeValue")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Data from API"}`))
	})
}
```

---

### SecurityHeadersMiddleware

- **Description:** Sets various HTTP security headers (CSP, HSTS, X-Frame-Options, etc.) to enhance application security.
- **Configuration:** `nova.SecurityHeadersConfig` (all fields optional, have sensible defaults)
  - `ContentTypeOptions string`: Defaults to `nosniff`.
  - `FrameOptions string`: Defaults to `DENY`.
  - `XSSProtection string`: Defaults to `1; mode=block`.
  - `ReferrerPolicy string`: Defaults to `strict-origin-when-cross-origin`.
  - `HSTSMaxAgeSeconds int`: Enables HSTS if > 0 (defaults to 0).
  - `HSTSIncludeSubdomains *bool`: Defaults to true if HSTS enabled.
  - `HSTSPreload bool`: Defaults to false.
  - `ContentSecurityPolicy string`: Defaults to `""`. **Important:** Define a policy specific to your app.
  - `PermissionsPolicy string`: Defaults to `""`.

#### Example

```go
func main() {
	router := nova.NewRouter()

	hstsEnabled := true // Example: Enable HSTS

	// Apply Security Headers middleware globally
	router.Use(nova.SecurityHeadersMiddleware(nova.SecurityHeadersConfig{
		// Example: Define a basic Content Security Policy
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; object-src 'none';",
		// Example: Enable HSTS for 1 year
		HSTSMaxAgeSeconds:     int(365 * 24 * time.Hour.Seconds()),
		HSTSIncludeSubdomains: &hstsEnabled, // Explicitly enable (default if HSTSMaxAgeSeconds > 0)
		// HSTSPreload: true, // Use with caution after submission
		// PermissionsPolicy: "geolocation=(), microphone=()", // Example policy
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Secure Page</h1>"))
	})
}
```

---

### TimeoutMiddleware

- **Description:** Enforces a maximum duration for request handling. Cancels the request context and returns 503 Service Unavailable (or calls custom handler) on timeout.
- **Configuration:** `nova.TimeoutConfig`
  - `Duration time.Duration`: Maximum processing time (required if used).
  - `TimeoutMessage string`: Response body on timeout (defaults to "Service timed out").
  - `TimeoutHandler http.Handler`: Custom handler for timeout events.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Timeout middleware globally
	router.Use(nova.TimeoutMiddleware(nova.TimeoutConfig{
		Duration: 5 * time.Second, // Set a 5-second timeout
		// TimeoutMessage: "Request took too long!", // Optional custom message
		// TimeoutHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 	http.Error(w, "Custom timeout response", http.StatusGatewayTimeout)
		// }), // Optional custom handler
	}))

	router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Slow handler started...")
		select {
		case <-time.After(10 * time.Second): // Simulate work longer than timeout
			log.Println("Slow handler finished (too late).")
			w.Write([]byte("Finished slow task."))
		case <-r.Context().Done(): // Context cancelled by timeout middleware
			log.Println("Slow handler cancelled by timeout.")
			// Note: Cannot write response here, middleware handles it.
			return
		}
	})

	router.Get("/fast", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Finished quickly."))
	})
}
```

---

### BasicAuthMiddleware

- **Description:** Protects routes using HTTP Basic Authentication, validating credentials with a provided function.
- **Configuration:** `nova.BasicAuthConfig`
  - `Realm string`: Authentication realm (defaults to "Restricted").
  - `Validator AuthValidator`: Function `func(user, pass string) bool` (required).
  - `StoreUserInContext bool`: Store authenticated user in context (defaults to false).
  - `ContextKey contextKey`: Context key for user (defaults to internal key).
- **Context Helper:** `nova.GetBasicAuthUser(ctx context.Context)` retrieves the user if stored.

#### Example

```go
// Simple validator function (replace with real logic)
func myAuthValidator(username, password string) bool {
	// WARNING: Hardcoded credentials are insecure! Use a proper check.
	return username == "admin" && password == "password123"
}

func main() {
	router := nova.NewRouter()

	// Apply Basic Auth middleware globally or to a group
	authMiddleware := nova.BasicAuthMiddleware(nova.BasicAuthConfig{
		Validator:          myAuthValidator,
		Realm:              "My Protected Area",
		StoreUserInContext: true, // Store the username
	})
	router.Use(authMiddleware) // Apply globally

	// Alternatively, apply to a group:
	// adminGroup := router.Group("/admin")
	// adminGroup.Use(authMiddleware)
	// adminGroup.Get("/dashboard", ...)

	router.Get("/secure", func(w http.ResponseWriter, r *http.Request) {
		user := nova.GetBasicAuthUser(r.Context())
		fmt.Fprintf(w, "Welcome, authenticated user: %s", user)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
```

---

### MethodOverrideMiddleware

- **Description:** Allows overriding the HTTP method via a header (`X-HTTP-Method-Override`) or form field (`_method` for POST requests).
- **Configuration:** `nova.MethodOverrideConfig`
  - `HeaderName string`: Header to check (defaults to `X-HTTP-Method-Override`).
  - `FormFieldName string`: Form field to check (defaults to `_method`). Set to `""` to disable form field check.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Method Override middleware globally
	router.Use(nova.MethodOverrideMiddleware(nil)) // Use defaults

	// Handler that might receive overridden methods
	router.Post("/resource", func(w http.ResponseWriter, r *http.Request) {
		// r.Method will be the overridden method (e.g., PUT, DELETE)
		fmt.Fprintf(w, "Handling resource with method: %s", r.Method)
	})
	// Need corresponding handlers for the methods you expect to be overridden
	router.Put("/resource", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Handling resource with method: PUT (via override or direct)")
	})
	router.Delete("/resource", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Handling resource with method: DELETE (via override or direct)")
	})


	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
// Example client request using header:
// curl -X POST -H "X-HTTP-Method-Override: DELETE" http://localhost:8080/resource

// Example client request using form field (HTML form):
// <form action="/resource" method="POST">
//   <input type="hidden" name="_method" value="PUT">
//   <button type="submit">Update Resource</button>
// </form>
```

---

### EnforceContentTypeMiddleware

- **Description:** Ensures requests for specific methods (POST, PUT, PATCH by default) have a `Content-Type` header matching an allowed list.
- **Configuration:** `nova.EnforceContentTypeConfig`
  - `AllowedTypes []string`: List of allowed Content-Types (e.g., `["application/json"]`) (required).
  - `MethodsToCheck []string`: Methods to check (defaults to POST, PUT, PATCH).
  - `OnError func(w http.ResponseWriter, r *http.Request, err error)`: Custom error handler.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Content-Type enforcement globally or to relevant groups/routes
	router.Use(nova.EnforceContentTypeMiddleware(nova.EnforceContentTypeConfig{
		AllowedTypes: []string{"application/json", "application/xml"},
		// MethodsToCheck: []string{"POST"}, // Optional: Only check POST
	}))

	router.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		// Handler logic assumes Content-Type is valid here
		w.Write([]byte("Data submitted successfully."))
	})

	router.Get("/fetch", func(w http.ResponseWriter, r *http.Request) {
		// GET requests are not checked by default
		w.Write([]byte("Data fetched."))
	})
}
// Example failing request:
// curl -X POST -d 'data' http://localhost:8080/submit
// -> 400 Bad Request (Missing Content-Type)

// Example failing request:
// curl -X POST -H "Content-Type: text/plain" -d 'data' http://localhost:8080/submit
// -> 415 Unsupported Media Type

// Example successful request:
// curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' http://localhost:8080/submit
// -> 200 OK
```

---

### CacheControlMiddleware

- **Description:** Sets the `Cache-Control` header on all responses.
- **Configuration:** `nova.CacheControlConfig`
  - `CacheControlValue string`: The value for the header (e.g., "no-store") (required).

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Cache-Control middleware globally
	router.Use(nova.CacheControlMiddleware(nova.CacheControlConfig{
		CacheControlValue: "no-store, no-cache, must-revalidate",
	}))

	router.Get("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
```

---

### GzipMiddleware

- **Description:** Compresses response bodies using gzip if the client supports it (`Accept-Encoding: gzip`). Uses a `sync.Pool` for `gzip.Writer` reuse to improve performance.
- **Configuration:** `nova.GzipConfig`
  - `CompressionLevel int`: Gzip level (e.g., `gzip.BestSpeed`, `gzip.DefaultCompression`, `gzip.BestCompression`). Defaults to `gzip.DefaultCompression` (-1).
  - `AddVaryHeader *bool`: Adds `Vary: Accept-Encoding` header. Defaults to `true`. Use `new(bool)` to set explicitly (e.g., `AddVaryHeader: new(bool) // false`).
  - `Logger *log.Logger`: Optional logger for errors. Defaults to `log.Default()`.
  - `Pool *sync.Pool`: Optional `sync.Pool` for `gzip.Writer` reuse. Defaults to an internal pool.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Custom pool example (optional)
	gzipPool := &sync.Pool{
		New: func() interface{} {
			gw, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
			return gw
		},
	}

	// Apply Gzip middleware globally
	router.Use(nova.GzipMiddleware(&nova.GzipConfig{
		CompressionLevel: gzip.BestSpeed, // Optional: Prioritize speed
		// AddVaryHeader:    func() *bool { b := false; return &b }(), // Explicitly disable Vary
		Pool:             gzipPool,       // Optional: Use custom pool
	}))

	router.Get("/large-data", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a large response
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		for i := 0; i < 1000; i++ {
			w.Write([]byte("This is line number "))
			w.Write([]byte(http.StatusText(i + 200)))
			w.Write([]byte("\n"))
		}
}

// Example request:
// curl -H "Accept-Encoding: gzip" http://localhost:8080/large-data --output - | gunzip
```

---

### CSRFMiddleware

- **Description:** Provides Cross-Site Request Forgery (CSRF) protection using the Double Submit Cookie pattern. It sets a secure, HttpOnly cookie and expects a matching token in a header or form field for unsafe HTTP methods (POST, PUT, DELETE, etc.).
- **Configuration:** `nova.CSRFConfig`
  - `Logger *log.Logger`: Optional logger. Defaults to `log.Default()`.
  - `FieldName string`: Form field name for the token. Defaults to `"csrf_token"`.
  - `HeaderName string`: HTTP header name for the token. Defaults to `"X-CSRF-Token"`.
  - `CookieName string`: Name of the HttpOnly cookie storing the secret. Defaults to `"_csrf"`.
  - `ContextKey contextKey`: Context key to store the expected token. Defaults to internal package key.
  - `ErrorHandler http.HandlerFunc`: Handler called on CSRF failure. Defaults to 403 Forbidden.
  - `CookiePath string`: Path for the CSRF cookie. Defaults to `"/"`.
  - `CookieDomain string`: Domain for the CSRF cookie. Defaults to `""`.
  - `CookieMaxAge time.Duration`: Max age of the cookie. Defaults to 12 hours.
  - `CookieSecure bool`: Secure flag for the cookie (requires HTTPS). Defaults to `false`. **Set to `true` in production.**
  - `CookieSameSite http.SameSite`: SameSite attribute. Defaults to `http.SameSiteLaxMode`.
  - `TokenLength int`: Byte length of the generated token. Defaults to 32.
  - `SkipMethods []string`: HTTP methods exempt from checks. Defaults to `["GET", "HEAD", "OPTIONS", "TRACE"]`.

#### Example

```go
// Simple template to include CSRF token
var formTmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<body>
	<h2>CSRF Demo Form</h2>
	<form method="POST" action="/submit">
		<input type="hidden" name="csrf_token" value="{{.}}">
		<label for="data">Data:</label>
		<input type="text" id="data" name="data">
		<button type="submit">Submit</button>
	</form>
</body>
</html>
`))

func main() {
	router := nova.NewRouter()

	// Apply CSRF middleware
	router.Use(nova.CSRFMiddleware(&nova.CSRFConfig{
		CookieSecure:   false, // Set to true if using HTTPS
		CookieSameSite: http.SameSiteStrictMode, // Often preferred
		// ErrorHandler: func(w http.ResponseWriter, r *http.Request) { ... }, // Custom error
	}))

	// Handler to display the form
	router.Get("/form", func(w http.ResponseWriter, r *http.Request) {
		// Get the token set by the middleware via context helper
		csrfToken := nova.GetCSRFToken(r.Context())
		w.Header().Set("Content-Type", "text/html")
		formTmpl.Execute(w, csrfToken) // Pass token to template
	})

	// Handler to process the form submission
	router.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		data := r.FormValue("data")
		fmt.Fprintf(w, "Received data: %s", data)
	})
}

// Example request (simulating a valid POST after getting the form):
// 1. Get form to get cookie and token: curl -c cookies.txt http://localhost:8080/form
// 2. Extract token from HTML output (e.g., TOKEN_VALUE)
// 3. Make POST request with cookie and token header:
//    curl -b cookies.txt -X POST -H "X-CSRF-Token: TOKEN_VALUE" -d "data=hello" http://localhost:8080/submit
// Example request (simulating a failed POST - missing token):
// curl -b cookies.txt -X POST -d "data=hello" http://localhost:8080/submit
// Response: Forbidden
```

---

### ETagMiddleware

- **Description:** Adds an `ETag` header to successful responses based on a hash of the response body. Handles `If-None-Match` conditional requests, potentially returning a `304 Not Modified` status without the response body if the client's cached ETag matches. **Note:** This middleware buffers the entire response body in memory to calculate the hash, which may be unsuitable for very large responses.
- **Configuration:** `nova.ETagConfig`
  - `Weak bool`: Generate weak ETags (prefixed with `W/`). Defaults to `false` (strong ETags).
  - `SkipNoContent bool`: Skip ETag generation/checking for `204 No Content` responses. Defaults to `true`.

#### Example

```go
var lastModified = time.Now()
var responseBody = "Initial content"

func main() {
	router := nova.NewRouter()

	// Apply ETag middleware
	router.Use(nova.ETagMiddleware(&nova.ETagConfig{
		// Weak: true, // Optional: Use weak ETags
	}))

	router.Get("/content", func(w http.ResponseWriter, r *http.Request) {
		// Simulate content that might change
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Cache-Control", "max-age=60") // Advise caching
		fmt.Fprint(w, responseBody)
	})

	// Example route to change the content
	router.Get("/update", func(w http.ResponseWriter, r *http.Request) {
		responseBody = fmt.Sprintf("Content updated at %s", time.Now())
		lastModified = time.Now()
		fmt.Fprint(w, "Content updated. Try /content again.")
	})
}

// Example requests:
// 1. First request: curl -v http://localhost:8080/content
//    (Note the ETag header in the response, e.g., ETag: "HASH_VALUE")
// 2. Second request (client sends If-None-Match):
//    curl -v -H 'If-None-Match: "HASH_VALUE"' http://localhost:8080/content
//    (Response should be 304 Not Modified with an empty body)
// 3. Update content: curl http://localhost:8080/update
// 4. Request again with old ETag:
//    curl -v -H 'If-None-Match: "HASH_VALUE"' http://localhost:8080/content
//    (Response should be 200 OK with the new content and a *new* ETag)
```

---

### HealthCheckMiddleware

- **Description:** Provides a dedicated health check endpoint (e.g., `/healthz`). Requests to this path are handled directly by the middleware (returning a status, often 200 OK), bypassing subsequent middleware and application handlers. Useful for load balancers and monitoring systems.
- **Configuration:** `nova.HealthCheckConfig`
  - `Path string`: The URL path for the health check endpoint. Defaults to `"/healthz"`.
  - `Handler http.HandlerFunc`: The handler function to execute for the health check. If nil, a default handler returning 200 OK with "OK" body is used. Can be customized to check database connections, etc.

#### Example

```go
// Example custom health check handler
func customHealthHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate checking a dependency (e.g., database)
	time.Sleep(10 * time.Millisecond)
	isHealthy := true // Replace with actual check logic

	if isHealthy {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "UP"}`)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, `{"status": "DOWN", "reason": "database connection failed"}`)
	}
}

func main() {
	router := nova.NewRouter()

	// Apply HealthCheck middleware (often placed early, but after recovery/logging if needed)
	router.Use(nova.HealthCheckMiddleware(&nova.HealthCheckConfig{
		Path:    "/status",             // Optional: Custom path
		Handler: customHealthHandler, // Optional: Custom check logic
	}))

	// Other middleware and routes...
	router.Use(nova.LoggingMiddleware(nil)) // Example: Log other requests

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Main application page")
	})
}

// Example requests:
// curl http://localhost:8080/status
// Response (with custom handler): {"status": "UP"}
// curl http://localhost:8080/
// Response: Main application page (Logs will show this request, but not the /status request if logging is after health check)
```

---

### RealIPMiddleware

- **Description:** Extracts the client's real IP address from trusted proxy headers (e.g., `X-Forwarded-For`, `X-Real-IP`). **Warning:** Only use behind a trusted proxy.
- **Configuration:** `nova.RealIPConfig`
  - `TrustedProxyCIDRs []string`: CIDR ranges of trusted proxies (e.g., `["10.0.0.0/8", "192.168.1.1/32"]`). Required for header trusting.
  - `IPHeaders []string`: Headers to check in order (defaults to `X-Forwarded-For`, `X-Real-IP`).
  - `StoreInContext bool`: Store the real IP in context (defaults to true).
  - `ContextKey contextKey`: Context key for IP (defaults to internal key).
- **Context Helper:** `nova.GetRealIP(ctx context.Context)` retrieves the IP.
- **Usage:**

```go
func main() {
	router := nova.NewRouter()

	// Apply RealIP middleware globally (ensure it runs early)
	router.Use(nova.RealIPMiddleware(nova.RealIPConfig{
		// IMPORTANT: Only list CIDRs of proxies you TRUST
		TrustedProxyCIDRs: []string{"127.0.0.1/32", "::1/128"}, // Example: Trust localhost proxy
		IPHeaders:         []string{"X-Forwarded-For", "X-Real-IP"}, // Default
		StoreInContext:    true, // Default
	}))

	// Apply Logging middleware *after* RealIP so logs show the real IP
	router.Use(nova.LoggingMiddleware(nil))

	router.Get("/ip", func(w http.ResponseWriter, r *http.Request) {
		realIP := nova.GetRealIP(r.Context())
		// r.RemoteAddr might also be updated (with port 0) if IP found via header
		fmt.Fprintf(w, "Your detected IP: %s\n", realIP)
		fmt.Fprintf(w, "Request RemoteAddr: %s\n", r.RemoteAddr)
	})
}
// Example request (simulating proxy):
// curl -H "X-Forwarded-For: 1.2.3.4" http://localhost:8080/ip
```

---

### MaxRequestBodySizeMiddleware

- **Description:** Limits the size of incoming request bodies using `http.MaxBytesReader`.
- **Configuration:** `nova.MaxRequestBodySizeConfig`
  - `LimitBytes int64`: Maximum body size in bytes (required).
  - `OnError func(w http.ResponseWriter, r *http.Request)`: Custom error handler (defaults to 413 response).

#### Example

```go

func main() {
	router := nova.NewRouter()

	// Apply Max Body Size middleware globally or to upload routes
	router.Use(nova.MaxRequestBodySizeMiddleware(nova.MaxRequestBodySizeConfig{
		LimitBytes: 1 * 1024 * 1024, // 1 MB limit
		// OnError: func(w http.ResponseWriter, r *http.Request) {
		// 	http.Error(w, "Request body too large!", http.StatusRequestEntityTooLarge)
		// }, // Optional custom handler
	}))

	router.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Attempt to read the body. MaxBytesReader will return an error
		// if the limit is exceeded during the read.
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			// Error might be due to size limit or other read issues.
			// The middleware's default OnError usually handles the 413 response
			// before the handler even gets here if ContentLength is too large.
			// If reading fails *during* the stream due to limit, this handler sees error.
			log.Printf("Error reading body: %v", err)
			// Check if it was a MaxBytesError (though http package might not export it easily)
			// A simple check:
			if r.ContentLength == -1 { // If ContentLength wasn't known beforehand
				// Assume error might be due to limit exceeded during read
				// The http.MaxBytesReader already wrote the 413 error response
				return
			}
			// Otherwise, handle other potential read errors
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		log.Printf("Received %d bytes", len(bodyBytes))
		w.Write([]byte("Upload received successfully."))
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
// Example request exceeding limit:
// curl -X POST --data-binary @large_file.dat http://localhost:8080/upload
// -> 413 Request Entity Too Large
```

---

### TrailingSlashRedirectMiddleware

- **Description:** Redirects requests to add or remove a trailing slash from the URL path for consistency.
- **Configuration:** `nova.TrailingSlashRedirectConfig`
  - `AddSlash bool`: Enforce trailing slash (defaults to false - removes slash).
  - `RedirectCode int`: HTTP redirect status code (defaults to 301). Use 308 for POST/PUT etc. to preserve method.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Trailing Slash middleware globally (usually early)
	router.Use(nova.TrailingSlashRedirectMiddleware(nova.TrailingSlashRedirectConfig{
		AddSlash:     false, // Default: remove trailing slash
		RedirectCode: http.StatusMovedPermanently, // Default: 301
		// Or to enforce slashes:
		// AddSlash: true,
		// RedirectCode: http.StatusPermanentRedirect, // 308 to preserve method
	}))

	// Define routes WITHOUT the trailing slash (if AddSlash is false)
	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("List of users"))
	})
	router.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("List of products"))
	})
}
// Example request:
// curl -L http://localhost:8080/users/
// -> Redirects (301) to http://localhost:8080/users
// -> Responds with "List of users"
```

---

### ForceHTTPSMiddleware

- **Description:** Redirects incoming HTTP requests to their HTTPS equivalent.
- **Configuration:** `nova.ForceHTTPSConfig`
  - `TargetHost string`: Override host in redirect URL (defaults to request host).
  - `TargetPort int`: Override port in redirect URL (defaults to standard 443).
  - `RedirectCode int`: Redirect status code (defaults to 301).
  - `ForwardedProtoHeader string`: Header to check for original protocol (defaults to `X-Forwarded-Proto`).
  - `TrustForwardedHeader *bool`: Trust the forwarded header (defaults to true). Set false if proxy doesn't set it reliably.

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Force HTTPS middleware globally (very early)
	trustHeader := true
	router.Use(nova.ForceHTTPSMiddleware(nova.ForceHTTPSConfig{
		// RedirectCode: http.StatusPermanentRedirect, // Use 308 if needed
		// ForwardedProtoHeader: "X-Scheme", // If your proxy uses a different header
		TrustForwardedHeader: &trustHeader, // Default is true
	}))

	// Add other middleware like HSTS *after* ForceHTTPS potentially
	router.Use(nova.SecurityHeadersMiddleware(nova.SecurityHeadersConfig{
		HSTSMaxAgeSeconds: 31536000, // Example: 1 year HSTS
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the secure site!"))
	})
}
// Example request:
// curl http://localhost:8080
// -> Redirects (301) to https://localhost:8080 (or https://localhost if port 443)
```

---

### ConcurrencyLimiterMiddleware

- **Description:** Limits the number of requests processed concurrently using a semaphore.
- **Configuration:** `nova.ConcurrencyLimiterConfig`
  - `MaxConcurrent int`: Max concurrent requests (required).
  - `WaitTimeout time.Duration`: Max time to wait for a slot (0 = wait forever).
  - `OnLimitExceeded func(w http.ResponseWriter, r *http.Request)`: Custom handler for limit exceeded (defaults to 503).

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Concurrency Limiter middleware globally or to heavy routes
	router.Use(nova.ConcurrencyLimiterMiddleware(nova.ConcurrencyLimiterConfig{
		MaxConcurrent: 10, // Allow only 10 requests at a time
		WaitTimeout:   2 * time.Second, // Wait max 2s for a slot
		// OnLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
		// 	http.Error(w, "Too busy, try later", http.StatusServiceUnavailable)
		// }, // Optional custom handler
	}))

	router.Get("/process", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Processing request...")
		time.Sleep(5 * time.Second) // Simulate work
		log.Println("Finished processing.")
		w.Write([]byte("Processing complete."))
	})
}
// Example: Run 15 concurrent requests:
// for i in {1..15}; do curl http://localhost:8080/process & done
// -> First 10 start immediately, next ~5 wait up to 2s. Some might get 503.
```

---

### MaintenanceModeMiddleware

- **Description:** Returns 503 Service Unavailable if enabled via an atomic flag, allowing bypass for specified IPs/CIDRs.
- **Configuration:** `nova.MaintenanceModeConfig`
  - `EnabledFlag *atomic.Bool`: Pointer to the control flag (required).
  - `AllowedIPs []string`: IPs/CIDRs that bypass maintenance (e.g., `["192.168.1.100", "10.0.0.0/8"]`).
  - `StatusCode int`: Status code during maintenance (defaults to 503).
  - `RetryAfterSeconds int`: Value for `Retry-After` header (defaults to 300).
  - `Message string`: Response body during maintenance.
  - `Logger *log.Logger`: Logger for errors (defaults to `log.Default()`).

#### Example

```go
var maintenanceEnabled atomic.Bool // The control flag

func main() {
	router := nova.NewRouter()

	// Set initial state (e.g., false = not in maintenance)
	maintenanceEnabled.Store(false)

	// Apply Maintenance Mode middleware globally (very early)
	router.Use(nova.MaintenanceModeMiddleware(nova.MaintenanceModeConfig{
		EnabledFlag:       &maintenanceEnabled,
		AllowedIPs:        []string{"127.0.0.1", "::1"}, // Allow localhost bypass
		RetryAfterSeconds: 600,                          // 10 minutes
		Message:           "Down for scheduled maintenance. Please try again later.",
	}))

	// Example route to toggle maintenance mode (in real app, use signals or admin API)
	router.Get("/admin/maintenance/on", func(w http.ResponseWriter, r *http.Request) {
		maintenanceEnabled.Store(true)
		w.Write([]byte("Maintenance mode ENABLED"))
		log.Println("Maintenance mode ENABLED")
	})
	router.Get("/admin/maintenance/off", func(w http.ResponseWriter, r *http.Request) {
		maintenanceEnabled.Store(false)
		w.Write([]byte("Maintenance mode DISABLED"))
		log.Println("Maintenance mode DISABLED")
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Application is running normally."))
	})
}
// Example:
// curl http://localhost:8080/ -> Shows normal page
// curl http://localhost:8080/admin/maintenance/on -> Enables maintenance
// curl http://<your_external_ip>:8080/ -> Shows 503 maintenance page
// curl http://localhost:8080/ -> Still shows normal page (due to AllowedIPs)
// curl http://localhost:8080/admin/maintenance/off -> Disables maintenance
```

---

### IPFilterMiddleware

- **Description:** Restricts access based on client IP using allow/block lists (CIDR supported).
- **Configuration:** `nova.IPFilterConfig`
  - `AllowedIPs []string`: Allowed IPs/CIDRs.
  - `BlockedIPs []string`: Blocked IPs/CIDRs (takes precedence).
  - `BlockByDefault bool`: Block IPs not matching any list (defaults to false - allow unless blocked).
  - `OnForbidden func(w http.ResponseWriter, r *http.Request)`: Custom handler for forbidden IPs (defaults to 403).
  - `Logger *log.Logger`: Logger for errors (defaults to `log.Default()`).

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply IP Filter middleware globally or to specific areas
	router.Use(nova.IPFilterMiddleware(nova.IPFilterConfig{
		// Example 1: Allow only specific IPs/ranges
		// AllowedIPs:     []string{"192.168.1.0/24", "10.0.0.5"},
		// BlockByDefault: true, // Block anything not in AllowedIPs

		// Example 2: Block specific IPs, allow others
		BlockedIPs:     []string{"1.2.3.4", "5.6.7.0/24"},
		BlockByDefault: false, // Default: Allow unless explicitly blocked

		// OnForbidden: func(w http.ResponseWriter, r *http.Request) {
		// 	http.Error(w, "Access denied from your IP.", http.StatusForbidden)
		// }, // Optional custom handler
	}))

	router.Get("/sensitive-data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is sensitive information."))
	})
}
```

---

### RateLimitMiddleware

- **Description:** Simple in-memory rate limiting using a token bucket algorithm per client IP (by default). **Warning:** Basic, single-instance only, potential memory growth without cleanup.
- **Configuration:** `nova.RateLimiterConfig`
  - `Requests int`: Max requests per duration (required).
  - `Duration time.Duration`: Time window (required).
  - `Burst int`: Allowed burst size (defaults to `Requests`).
  - `KeyFunc func(r *http.Request) string`: Function to get client key (defaults to IP).
  - `OnLimitExceeded func(w http.ResponseWriter, r *http.Request)`: Custom handler for limit (defaults to 429).
  - `CleanupInterval time.Duration`: How often to clean old entries (0 = no cleanup).
  - `Logger *log.Logger`: Logger for errors (defaults to `log.Default()`).

#### Example

```go
func main() {
	router := nova.NewRouter()

	// Apply Rate Limiter middleware globally or to specific APIs
	router.Use(nova.RateLimitMiddleware(nova.RateLimiterConfig{
		Requests: 5,               // Allow 5 requests...
		Duration: 1 * time.Minute, // ...per minute
		Burst:    10,              // Allow initial burst of 10 requests
		// KeyFunc: func(r *http.Request) string { // Optional: Limit by API key header
		// 	key := r.Header.Get("X-API-Key")
		// 	if key == "" { return r.RemoteAddr } // Fallback to IP if no key
		// 	return key
		// },
		CleanupInterval: 10 * time.Minute, // Clean up old entries every 10 mins
		// OnLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
		// 	w.WriteHeader(http.StatusTooManyRequests)
		// 	w.Write([]byte("Rate limit exceeded. Please try again later."))
		// }, // Optional custom handler
	}))

	router.Get("/api/resource", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Resource data"))
	})
}
// Example: Hit the endpoint repeatedly
// for i in {1..15}; do curl -I http://localhost:8080/api/resource; sleep 0.1; done
// -> First ~10 requests get 200 OK, subsequent ones get 429 Too Many Requests
```
