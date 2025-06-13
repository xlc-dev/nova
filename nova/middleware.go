package nova

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// contextKey is a private type used for context keys to avoid collisions.
type contextKey string

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

// GetRequestID retrieves the request ID from the context, if available via
// RequestIDMiddleware using the default key.
func GetRequestID(ctx context.Context) string {
	return GetRequestIDWithKey(ctx, requestIDKey)
}

// GetRequestIDWithKey retrieves the request ID from the context using a specific key.
func GetRequestIDWithKey(ctx context.Context, key contextKey) string {
	if reqID, ok := ctx.Value(key).(string); ok {
		return reqID
	}
	return "unknown"
}

// GetRealIP retrieves the client's real IP from the context, if available via
// RealIPMiddleware using the default key.
func GetRealIP(ctx context.Context) string {
	return GetRealIPWithKey(ctx, realIPKey)
}

// GetRealIPWithKey retrieves the client's real IP from the context using a specific key.
func GetRealIPWithKey(ctx context.Context, key contextKey) string {
	if ip, ok := ctx.Value(key).(string); ok {
		return ip
	}
	return ""
}

// GetBasicAuthUser retrieves the authenticated username from the context, if available
// via BasicAuthMiddleware using the default key and if configured to store it.
func GetBasicAuthUser(ctx context.Context) string {
	return GetBasicAuthUserWithKey(ctx, basicAuthUserKey)
}

// GetBasicAuthUserWithKey retrieves the authenticated username from the context
// using a specific key.
func GetBasicAuthUserWithKey(ctx context.Context, key contextKey) string {
	if user, ok := ctx.Value(key).(string); ok {
		return user
	}
	return ""
}

// GetCSRFToken retrieves the CSRF token from the context, if available via
// CSRFMiddleware using the default key. This is the token expected in subsequent
// unsafe requests.
func GetCSRFToken(ctx context.Context) string {
	return GetCSRFTokenWithKey(ctx, csrfTokenKey)
}

// GetCSRFTokenWithKey retrieves the CSRF token from the context using a specific key.
func GetCSRFTokenWithKey(ctx context.Context, key contextKey) string {
	if token, ok := ctx.Value(key).(string); ok {
		return token
	}
	return ""
}

// responseWriterInterceptor wraps http.ResponseWriter to capture status code and size.
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode  int
	size        int64
	wroteHeader bool // Track if WriteHeader was called explicitly
}

// NewResponseWriterInterceptor creates a new interceptor.
func NewResponseWriterInterceptor(w http.ResponseWriter) *responseWriterInterceptor {
	return &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the status code.
func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.statusCode = statusCode
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the number of bytes written and calls the underlying Write.
// It also calls WriteHeader implicitly if not already called.
func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += int64(size)
	return size, err
}

// Flush implements the http.Flusher interface if the underlying writer supports it.
func (w *responseWriterInterceptor) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// bufferingResponseWriterInterceptor wraps http.ResponseWriter to capture status, size,
// and buffer the response body.
type bufferingResponseWriterInterceptor struct {
	http.ResponseWriter
	statusCode  int
	size        int64
	wroteHeader bool
	buffer      *bytes.Buffer // Buffer to hold the response body
}

// NewBufferingResponseWriterInterceptor creates a new buffering interceptor.
func NewBufferingResponseWriterInterceptor(w http.ResponseWriter) *bufferingResponseWriterInterceptor {
	return &bufferingResponseWriterInterceptor{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		buffer:         new(bytes.Buffer),
	}
}

// WriteHeader captures the status code. Does NOT write to underlying writer yet.
func (w *bufferingResponseWriterInterceptor) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.statusCode = statusCode
	w.wroteHeader = true
}

// Write captures the number of bytes written and writes to the internal buffer.
// It also sets the status code implicitly if not already set.
func (w *bufferingResponseWriterInterceptor) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.buffer.Write(b)
	w.size += int64(size)
	return size, err
}

// Flush implements the http.Flusher interface if the underlying writer supports it.
func (w *bufferingResponseWriterInterceptor) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		// We could potentially write the buffer and flush here,
		// but standard middleware pattern usually writes fully at the end.
		// For simplicity, we might just flush the underlying if needed,
		// though the buffer hasn't been written yet.
		flusher.Flush()
	}
}

// WriteCapturedData writes the captured status code, headers, and buffered body
// to the original ResponseWriter. Returns the number of bytes written from the body.
func (w *bufferingResponseWriterInterceptor) WriteCapturedData() (int64, error) {
	w.ResponseWriter.WriteHeader(w.statusCode)
	written, err := io.Copy(w.ResponseWriter, w.buffer)
	return written, err
}

// StatusCode returns the captured status code.
func (w *bufferingResponseWriterInterceptor) StatusCode() int {
	return w.statusCode
}

// Size returns the captured response size.
func (w *bufferingResponseWriterInterceptor) Size() int64 {
	return w.size
}

// Body returns the buffered response body bytes.
func (w *bufferingResponseWriterInterceptor) Body() []byte {
	return w.buffer.Bytes()
}

// CSRFConfig holds configuration for CSRFMiddleware.
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

// CSRFMiddleware provides Cross-Site Request Forgery protection.
// It uses the "Double Submit Cookie" pattern. A random token is generated and
// set in a secure, HttpOnly cookie. For unsafe methods (POST, PUT, etc.),
// the middleware expects the same token to be present in a request header
// (e.g., X-CSRF-Token) or form field, sent by the frontend JavaScript.
func CSRFMiddleware(config *CSRFConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &CSRFConfig{}
	}
	// Set defaults for all configuration fields
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}
	if cfg.FieldName == "" {
		cfg.FieldName = "csrf_token"
	}
	if cfg.HeaderName == "" {
		cfg.HeaderName = "X-CSRF-Token"
	}
	if cfg.CookieName == "" {
		cfg.CookieName = "_csrf"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = csrfTokenKey
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				http.StatusText(http.StatusForbidden),
				http.StatusForbidden,
			)
		}
	}
	if cfg.CookiePath == "" {
		cfg.CookiePath = "/"
	}
	if cfg.CookieMaxAge == 0 {
		cfg.CookieMaxAge = 12 * time.Hour
	}
	if cfg.CookieSameSite == 0 {
		cfg.CookieSameSite = http.SameSiteLaxMode
	}
	if cfg.TokenLength == 0 {
		cfg.TokenLength = 32
	}
	if cfg.SkipMethods == nil {
		cfg.SkipMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
	}

	// Create a set for quick lookups of methods to skip
	skipMethodSet := make(map[string]struct{}, len(cfg.SkipMethods))
	for _, m := range cfg.SkipMethods {
		skipMethodSet[strings.ToUpper(m)] = struct{}{}
	}

	// Generate a random base64 encoded token
	generateToken := func() (string, error) {
		b := make([]byte, cfg.TokenLength)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(b), nil
	}

	// Define the secure comparison function once to avoid recreating it
	compare := func(a, b string) bool {
		// Check length first to prevent timing attacks on unequal lengths
		if len(a) != len(b) {
			return false
		}
		return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always generate or retrieve the real token
			var realToken string
			cookie, err := r.Cookie(cfg.CookieName)

			// If cookie exists and has a value, use it. Otherwise, generate a new one
			if err == nil && cookie.Value != "" {
				realToken = cookie.Value
			} else {
				newToken, err := generateToken()
				if err != nil {
					cfg.Logger.Printf(
						"[ERROR] CSRF: Failed to generate token: %v",
						err,
					)
					// This is a server error, not a client one
					http.Error(
						w,
						"Internal Server Error",
						http.StatusInternalServerError,
					)
					return
				}
				realToken = newToken

				// Set the new token in a cookie on the response
				http.SetCookie(w, &http.Cookie{
					Name:     cfg.CookieName,
					Value:    realToken,
					Path:     cfg.CookiePath,
					Domain:   cfg.CookieDomain,
					MaxAge:   int(cfg.CookieMaxAge.Seconds()),
					Secure:   cfg.CookieSecure,
					HttpOnly: true,
					SameSite: cfg.CookieSameSite,
				})
			}

			// Store the real token in the context. This is useful for HTML
			// templates to embed the token in forms
			ctx := context.WithValue(r.Context(), cfg.ContextKey, realToken)
			r = r.WithContext(ctx)

			// For "safe" methods, we're done. Just call the next handler
			if _, skip := skipMethodSet[strings.ToUpper(r.Method)]; skip {
				next.ServeHTTP(w, r)
				return
			}

			// For "unsafe" methods, we must validate the token from the request
			sentToken := r.Header.Get(cfg.HeaderName)
			if sentToken == "" {
				// If not in header, try the form field. This requires parsing the form,
				// which consumes the request body. Be aware of this side effect
				if err := r.ParseForm(); err == nil {
					sentToken = r.PostFormValue(cfg.FieldName)
				}
			}

			// Perform the secure comparison
			// If the sent token is empty or doesn't match the real token, fail
			if !compare(sentToken, realToken) {
				cfg.Logger.Printf(
					"[WARN] CSRF: Invalid token for %s %s. Token was empty or mismatched.",
					r.Method,
					r.URL.Path,
				)
				cfg.ErrorHandler(w, r)
				return
			}

			// If validation passes, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// ETagConfig holds configuration for ETagMiddleware.
type ETagConfig struct {
	// Weak determines if weak ETags (prefixed with W/) should be generated.
	// Weak ETags indicate semantic equivalence, not byte-for-byte identity.
	// Defaults to false (strong ETags).
	Weak bool
	// SkipNoContent determines whether to skip ETag generation/checking for
	// responses with status 204 No Content. Defaults to true.
	SkipNoContent bool
}

// ETagMiddleware adds ETag headers to responses and handles If-None-Match
// conditional requests, potentially returning a 304 Not Modified status.
// Note: This middleware buffers the entire response body in memory to calculate
// the ETag hash. This may be inefficient for very large responses.
func ETagMiddleware(config *ETagConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &ETagConfig{
			// Default SkipNoContent to true when config is nil
			SkipNoContent: true,
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use the buffering interceptor to capture the response
			interceptor := NewBufferingResponseWriterInterceptor(w)

			// Call the next handler, which writes to the interceptor's buffer
			next.ServeHTTP(interceptor, r)

			// Get response data from interceptor
			statusCode := interceptor.StatusCode()
			body := interceptor.Body()

			// Option to skip ETag for 204 No Content
			if cfg.SkipNoContent && statusCode == http.StatusNoContent {
				interceptor.WriteCapturedData() // Write original 204 response
				return
			}

			// Only generate ETag for successful responses (2xx) with a body
			if statusCode >= 200 && statusCode < 300 && len(body) > 0 {
				// Calculate ETag (SHA1 hash of the body)
				hasher := sha1.New()
				hasher.Write(body)
				hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

				// Format ETag (strong or weak)
				etag := fmt.Sprintf("\"%s\"", hash) // Strong ETag format
				if cfg.Weak {
					etag = "W/" + etag
				}

				// Check If-None-Match request header for a match
				ifNoneMatch := r.Header.Get("If-None-Match")
				if ifNoneMatch != "" {
					match := false
					if ifNoneMatch == "*" {
						match = true
					} else {
						// Parse comma-separated list of ETags
						for _, token := range strings.Split(ifNoneMatch, ",") {
							if strings.TrimSpace(token) == etag {
								match = true
								break
							}
						}
					}

					if match {
						interceptor.Header().Del("Content-Type")
						interceptor.Header().Del("Content-Length")
						interceptor.Header().Del("Transfer-Encoding")
						interceptor.WriteHeader(http.StatusNotModified)
						interceptor.buffer.Reset()
						interceptor.size = 0
						interceptor.WriteCapturedData()
						return
					}
				}

				// No match or no If-None-Match header, set ETag on the response
				interceptor.Header().Set("ETag", etag)
			}

			// Write the original (or modified) response data
			interceptor.WriteCapturedData()
		})
	}
}

// HealthCheckConfig holds configuration for HealthCheckMiddleware.
type HealthCheckConfig struct {
	// Path is the URL path for the health check endpoint. Defaults to "/healthz".
	Path string
	// Handler is the handler function to execute for the health check.
	// If nil, a default handler returning 200 OK with "OK" body is used.
	// You can provide a custom handler to check database connections, etc.
	Handler http.HandlerFunc
}

// HealthCheckMiddleware provides a simple health check endpoint.
func HealthCheckMiddleware(config *HealthCheckConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &HealthCheckConfig{}
	}
	if cfg.Path == "" {
		cfg.Path = "/healthz" // Common default
	}
	if cfg.Handler == nil {
		cfg.Handler = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request path matches the health check path
			if r.URL.Path == cfg.Path {
				// If it matches, execute the health check handler and stop processing
				cfg.Handler(w, r)
				return // Don't call the next handler in the chain
			}
			// If it doesn't match, pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingConfig holds configuration for LoggingMiddleware.
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

// LoggingMiddleware logs request details including method, path, status, size, and duration.
func LoggingMiddleware(config *LoggingConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &LoggingConfig{}
	}
	if cfg.Logger == nil {
		cfg.Logger = log.Default() // Use standard logger if none provided
	}
	if cfg.RequestIDKey == "" {
		cfg.RequestIDKey = requestIDKey
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestIDStr := ""
			if cfg.LogRequestID {
				reqID := GetRequestIDWithKey(r.Context(), cfg.RequestIDKey)
				if reqID != "unknown" {
					requestIDStr = fmt.Sprintf("[%s] ", reqID)
				}
			}

			cfg.Logger.Printf("[INFO] %sStarted %s %s from %s",
				requestIDStr, r.Method, r.URL.Path, r.RemoteAddr)

			// Wrap the response writer
			interceptor := NewResponseWriterInterceptor(w)

			next.ServeHTTP(interceptor, r)

			duration := time.Since(start)
			cfg.Logger.Printf(
				"[INFO] %sCompleted %s %s %d %s (%d bytes) in %v",
				requestIDStr,
				r.Method,
				r.URL.Path,
				interceptor.statusCode,
				http.StatusText(interceptor.statusCode),
				interceptor.size,
				duration,
			)
		})
	}
}

// RecoveryConfig holds configuration for RecoveryMiddleware.
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

// RecoveryMiddleware recovers from panics in downstream handlers.
func RecoveryMiddleware(config *RecoveryConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &RecoveryConfig{}
	}
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}
	if cfg.RequestIDKey == "" {
		cfg.RequestIDKey = requestIDKey
	}

	defaultHandler := func(w http.ResponseWriter, r *http.Request, err any) {
		// Check if headers were already sent (e.g., by a streaming handler before panic)
		// This is a best-effort check.
		if _, ok := w.(interface{ Status() int }); !ok { // Simple check if it's our interceptor
			if w.Header().Get("Content-Type") == "" { // Check if common headers are set
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}
	}

	handler := cfg.RecoveryHandler
	if handler == nil {
		handler = defaultHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestIDStr := ""
					if cfg.LogRequestID {
						reqID := GetRequestIDWithKey(r.Context(), cfg.RequestIDKey)
						if reqID != "unknown" {
							requestIDStr = fmt.Sprintf("[%s] ", reqID)
						}
					}
					cfg.Logger.Printf(
						"[ERROR] %sRecovered from panic: %v\n%s",
						requestIDStr,
						err,
						debug.Stack(),
					)
					handler(w, r, err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDConfig holds configuration for RequestIDMiddleware.
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

// RequestIDMiddleware retrieves a request ID from a header or generates one.
// It sets the ID in the response header and request context.
func RequestIDMiddleware(config *RequestIDConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &RequestIDConfig{}
	}
	if cfg.HeaderName == "" {
		cfg.HeaderName = "X-Request-ID"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = requestIDKey
	}
	if cfg.Generator == nil {
		cfg.Generator = func() string {
			// Simple default generator
			return strconv.FormatInt(time.Now().UnixNano(), 16)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get(cfg.HeaderName)
			if reqID == "" {
				reqID = cfg.Generator()
			}
			w.Header().Set(cfg.HeaderName, reqID)
			ctx := context.WithValue(r.Context(), cfg.ContextKey, reqID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORSConfig holds configuration for CORSMiddleware.
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

// CORSMiddleware sets Cross-Origin Resource Sharing headers.
func CORSMiddleware(config CORSConfig) Middleware {
	// Set defaults
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{"Content-Type", "Authorization", "X-Request-ID"}
	}
	if config.MaxAgeSeconds == 0 {
		config.MaxAgeSeconds = 86400 // Default to 24 hours
	}

	allowAllOrigins := false
	allowedOriginsMap := make(map[string]struct{})
	for _, origin := range config.AllowedOrigins {
		if origin == "*" {
			allowAllOrigins = true
			break // No need to process others if * is present
		}
		allowedOriginsMap[strings.ToLower(origin)] = struct{}{}
	}

	methods := strings.Join(config.AllowedMethods, ", ")
	headers := strings.Join(config.AllowedHeaders, ", ")
	expose := strings.Join(config.ExposedHeaders, ", ")
	maxAge := ""
	if config.MaxAgeSeconds > 0 {
		maxAge = strconv.Itoa(config.MaxAgeSeconds)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowedOrigin := ""

			if allowAllOrigins {
				// Cannot set AllowCredentials to true with wildcard origin
				if config.AllowCredentials {
					// If credentials are required, reflect the specific origin
					if origin != "" {
						allowedOrigin = origin
					}
				} else {
					allowedOrigin = "*"
				}
			} else if origin != "" {
				if _, ok := allowedOriginsMap[strings.ToLower(origin)]; ok {
					allowedOrigin = origin
				}
			}

			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			}
			// Vary header is important for caching proxies
			w.Header().Add("Vary", "Origin")

			if r.Method == http.MethodOptions {
				// Handle preflight request
				w.Header().Add("Vary", "Access-Control-Request-Method")
				w.Header().Add("Vary", "Access-Control-Request-Headers")
				w.Header().Set("Access-Control-Allow-Methods", methods)
				w.Header().Set("Access-Control-Allow-Headers", headers)
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				if maxAge != "" {
					w.Header().Set("Access-Control-Max-Age", maxAge)
				}
				if expose != "" {
					w.Header().Set("Access-Control-Expose-Headers", expose)
				}
				w.WriteHeader(http.StatusNoContent) // Preflight should return 204
				return
			}

			// For actual requests
			if config.AllowCredentials && allowedOrigin != "*" && allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if expose != "" {
				w.Header().Set("Access-Control-Expose-Headers", expose)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeadersConfig holds configuration for SecurityHeadersMiddleware.
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

// SecurityHeadersMiddleware sets common security headers.
func SecurityHeadersMiddleware(config SecurityHeadersConfig) Middleware {
	// Set defaults
	if config.ContentTypeOptions == "" {
		config.ContentTypeOptions = "nosniff"
	}
	if config.FrameOptions == "" {
		config.FrameOptions = "DENY"
	}
	if config.ReferrerPolicy == "" {
		config.ReferrerPolicy = "strict-origin-when-cross-origin"
	}

	hstsValue := ""
	if config.HSTSMaxAgeSeconds > 0 {
		hstsValue = "max-age=" + strconv.Itoa(config.HSTSMaxAgeSeconds)
		includeSubdomains := true // Default to true if HSTS is enabled
		if config.HSTSIncludeSubdomains != nil {
			includeSubdomains = *config.HSTSIncludeSubdomains
		}
		if includeSubdomains {
			hstsValue += "; includeSubDomains"
		}
		if config.HSTSPreload {
			hstsValue += "; preload"
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hdr := w.Header()
			if config.ContentTypeOptions != "" {
				hdr.Set("X-Content-Type-Options", config.ContentTypeOptions)
			}
			if config.FrameOptions != "" {
				hdr.Set("X-Frame-Options", config.FrameOptions)
			}
			if config.ReferrerPolicy != "" {
				hdr.Set("Referrer-Policy", config.ReferrerPolicy)
			}
			if config.ContentSecurityPolicy != "" {
				hdr.Set("Content-Security-Policy", config.ContentSecurityPolicy)
			}
			if config.PermissionsPolicy != "" {
				hdr.Set("Permissions-Policy", config.PermissionsPolicy)
			}
			// Only set HSTS on HTTPS connections
			if hstsValue != "" && r.TLS != nil {
				hdr.Set("Strict-Transport-Security", hstsValue)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// timeoutResponseWriter wraps http.ResponseWriter to prevent panics from
// superfluous WriteHeader calls after a timeout.
type timeoutResponseWriter struct {
	http.ResponseWriter
	mu          sync.Mutex
	wroteHeader bool
}

func (w *timeoutResponseWriter) WriteHeader(statusCode int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.wroteHeader {
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
}

func (w *timeoutResponseWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.wroteHeader {
		// Implicitly write 200 OK if not already written
		w.ResponseWriter.WriteHeader(http.StatusOK)
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(b)
}

// TimeoutConfig holds configuration for TimeoutMiddleware.
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

// TimeoutMiddleware sets a maximum duration for handling requests.
func TimeoutMiddleware(config TimeoutConfig) Middleware {
	if config.Duration <= 0 {
		// No timeout, return identity middleware
		return func(next http.Handler) http.Handler { return next }
	}
	if config.TimeoutMessage == "" {
		config.TimeoutMessage = "Service timed out"
	}

	return func(next http.Handler) http.Handler {
		if config.TimeoutHandler != nil {
			// Use custom timeout handler if provided
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, cancel := context.WithTimeout(r.Context(), config.Duration)
				defer cancel()

				done := make(chan struct{})
				panicChan := make(chan any, 1)
				tw := &timeoutResponseWriter{ResponseWriter: w}

				go func() {
					defer func() {
						if p := recover(); p != nil {
							panicChan <- p
						}
						close(done)
					}()
					next.ServeHTTP(tw, r.WithContext(ctx))
				}()

				select {
				case <-done:
					// Handler finished within timeout. Re-panic if it panicked
					select {
					case p := <-panicChan:
						panic(p)
					default:
						// No panic.
					}
				case <-ctx.Done():
					// Timeout occurred.
					// Ensure we don't call the timeout handler if the main handler
					// already wrote the headers
					tw.mu.Lock()
					if tw.wroteHeader {
						tw.mu.Unlock()
						return
					}
					config.TimeoutHandler.ServeHTTP(w, r)
					tw.wroteHeader = true
					tw.mu.Unlock()
				}
			})
		} else {
			// Use standard http.TimeoutHandler for the simple case
			return http.TimeoutHandler(next, config.Duration, config.TimeoutMessage)
		}
	}
}

// AuthValidator is a function type that validates the provided username and password.
// It returns true if the credentials are valid.
type AuthValidator func(username, password string) bool

// BasicAuthConfig holds configuration for BasicAuthMiddleware.
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

// BasicAuthMiddleware provides simple HTTP Basic Authentication.
func BasicAuthMiddleware(config BasicAuthConfig) Middleware {
	if config.Realm == "" {
		config.Realm = "Restricted"
	}
	if config.Validator == nil {
		panic("BasicAuthMiddleware: Validator function is required")
	}
	if config.ContextKey == "" {
		config.ContextKey = basicAuthUserKey
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || !config.Validator(username, password) {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+config.Realm+`"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			req := r
			if config.StoreUserInContext {
				ctx := context.WithValue(r.Context(), config.ContextKey, username)
				req = r.WithContext(ctx)
			}
			next.ServeHTTP(w, req)
		})
	}
}

// MethodOverrideConfig holds configuration for MethodOverrideMiddleware.
type MethodOverrideConfig struct {
	// HeaderName is the name of the header checked for the method override value.
	// Defaults to "X-HTTP-Method-Override".
	HeaderName string
	// FormFieldName is the name of the form field (used for POST requests)
	// checked for the method override value if the header is not present.
	// Set to "" to disable form field checking. Defaults to "_method".
	FormFieldName string
}

// MethodOverrideMiddleware checks a header or form field to override the request method.
func MethodOverrideMiddleware(config *MethodOverrideConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &MethodOverrideConfig{}
	}
	if cfg.HeaderName == "" {
		cfg.HeaderName = "X-HTTP-Method-Override"
	}
	if cfg.FormFieldName == "" {
		// Allow disabling form field check explicitly
		if config != nil && config.FormFieldName != "" {
			cfg.FormFieldName = ""
		} else {
			cfg.FormFieldName = "_method" // Default form field name
		}
	}

	// A set of methods that are safe to allow for overriding.
	allowedOverrideMethods := map[string]struct{}{
		http.MethodPut:    {},
		http.MethodPatch:  {},
		http.MethodDelete: {},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			overrideMethod := ""
			// Check header first
			if headerVal := r.Header.Get(cfg.HeaderName); headerVal != "" {
				overrideMethod = headerVal
			} else if cfg.FormFieldName != "" && r.Method == http.MethodPost {
				// Check form field only for POST requests
				if err := r.ParseForm(); err == nil {
					if formVal := r.PostFormValue(cfg.FormFieldName); formVal != "" {
						overrideMethod = formVal
					}
				}
			}

			if overrideMethod != "" {
				// Validate that the override method is one we allow.
				// This prevents clients from setting arbitrary methods.
				method := strings.ToUpper(overrideMethod)
				if _, ok := allowedOverrideMethods[method]; ok {
					r.Method = method
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// EnforceContentTypeConfig holds configuration for EnforceContentTypeMiddleware.
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

// ErrMissingContentType indicates the Content-Type header was absent.
var ErrMissingContentType = fmt.Errorf("missing Content-Type header")

// ErrUnsupportedContentType indicates the Content-Type was not in the allowed list.
var ErrUnsupportedContentType = fmt.Errorf("unsupported Content-Type")

// EnforceContentTypeMiddleware checks if the request's Content-Type header is allowed.
func EnforceContentTypeMiddleware(config EnforceContentTypeConfig) Middleware {
	if len(config.AllowedTypes) == 0 {
		panic("EnforceContentTypeMiddleware: AllowedTypes cannot be empty")
	}

	allowedMap := make(map[string]struct{}, len(config.AllowedTypes))
	for _, t := range config.AllowedTypes {
		allowedMap[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}

	methodsMap := make(map[string]struct{})
	if len(config.MethodsToCheck) == 0 {
		methodsMap[http.MethodPost] = struct{}{}
		methodsMap[http.MethodPut] = struct{}{}
		methodsMap[http.MethodPatch] = struct{}{}
	} else {
		for _, m := range config.MethodsToCheck {
			methodsMap[strings.ToUpper(m)] = struct{}{}
		}
	}

	onError := config.OnError
	if onError == nil {
		onError = func(w http.ResponseWriter, r *http.Request, err error) {
			if err == ErrMissingContentType {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else if err == ErrUnsupportedContentType {
				http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			} else {
				// Should not happen with default usage, but handle defensively
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if method requires Content-Type validation
			if _, shouldCheck := methodsMap[r.Method]; !shouldCheck {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")
			if contentType == "" {
				onError(w, r, ErrMissingContentType)
				return
			}

			// Parse media type (e.g., "application/json; charset=utf-8" -> "application/json")
			mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))

			if _, ok := allowedMap[mediaType]; !ok {
				onError(w, r, ErrUnsupportedContentType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CacheControlConfig holds configuration for CacheControlMiddleware.
type CacheControlConfig struct {
	// CacheControlValue is the string to set for the Cache-Control header. Required.
	// Example values: "no-store", "no-cache", "public, max-age=3600"
	CacheControlValue string
}

// CacheControlMiddleware sets the Cache-Control header for responses.
func CacheControlMiddleware(config CacheControlConfig) Middleware {
	if config.CacheControlValue == "" {
		panic("CacheControlMiddleware: CacheControlValue is required")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", config.CacheControlValue)
			next.ServeHTTP(w, r)
		})
	}
}

// gzipResponseWriter wraps http.ResponseWriter to provide gzip compression.
// It implements http.ResponseWriter and optionally http.Flusher.
type gzipResponseWriter struct {
	http.ResponseWriter              // Underlying writer
	Writer              *gzip.Writer // gzip writer
}

// Header returns the header map of the underlying ResponseWriter.
func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader sends the HTTP status code using the underlying ResponseWriter.
func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write writes the data to the gzip writer.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Flush implements http.Flusher if the underlying ResponseWriter supports it.
func (w *gzipResponseWriter) Flush() {
	if err := w.Writer.Flush(); err != nil {
		log.Printf("[DEBUG] GzipMiddleware: Error flushing gzip writer: %v", err)
	}
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// GzipConfig holds configuration options for the GzipMiddleware.
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

// GzipMiddleware returns middleware that compresses response bodies using gzip
// if the client indicates support via the Accept-Encoding header.
func GzipMiddleware(config *GzipConfig) Middleware {
	cfg := config
	if cfg == nil {
		cfg = &GzipConfig{}
	}

	// Set defaults
	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
	}

	level := cfg.CompressionLevel
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	addVary := true // Default to true
	if cfg.AddVaryHeader != nil {
		addVary = *cfg.AddVaryHeader // Use explicitly set value if provided
	}

	// Initialize pool if not provided
	pool := cfg.Pool
	if pool == nil {
		pool = &sync.Pool{
			New: func() any {
				gw, _ := gzip.NewWriterLevel(io.Discard, level)
				return gw
			},
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client accepts gzip
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r) // Pass through if gzip not accepted
				return
			}

			// Client accepts gzip, set response headers
			w.Header().Set("Content-Encoding", "gzip")
			if addVary {
				w.Header().Add("Vary", "Accept-Encoding")
			}

			// Get a gzip writer from the pool
			gz := pool.Get().(*gzip.Writer)
			gz.Reset(w) // Reset writer to use the current response writer

			// Defer closing and returning the writer to the pool
			defer func() {
				if err := gz.Close(); err != nil {
					// Log error during close, as it might indicate incomplete write
					logger.Printf("[ERROR] GzipMiddleware: Failed to close gzip writer: %v", err)
				}
				pool.Put(gz)
			}()

			// Create the wrapped response writer
			gzw := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gz,
			}

			// Call the next handler with the wrapped gzip writer
			next.ServeHTTP(gzw, r)
		})
	}
}

// RealIPConfig holds configuration for RealIPMiddleware.
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

// RealIPMiddleware extracts the client's real IP address from proxy headers.
// Warning: Only use this if you have a trusted proxy setting these headers correctly.
func RealIPMiddleware(config RealIPConfig) Middleware {
	if len(config.IPHeaders) == 0 {
		config.IPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	}
	if config.ContextKey == "" {
		config.ContextKey = realIPKey
	}

	var trustedNets []*net.IPNet
	for _, cidr := range config.TrustedProxyCIDRs {
		_, ipNet, err := net.ParseCIDR(strings.TrimSpace(cidr))
		if err == nil {
			trustedNets = append(trustedNets, ipNet)
		} else {
			log.Printf("[WARN] RealIPMiddleware: Invalid CIDR for trusted proxy: %s", cidr)
		}
	}

	isTrusted := func(ip net.IP) bool {
		if ip == nil || len(trustedNets) == 0 {
			return false // No trusted proxies configured or invalid IP
		}
		for _, network := range trustedNets {
			if network.Contains(ip) {
				return true
			}
		}
		return false
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remoteIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				remoteIPStr = r.RemoteAddr // Handle cases without port (e.g., Unix sockets)
			}
			parsedRemoteIP := net.ParseIP(remoteIPStr)
			originalRemoteAddr := r.RemoteAddr // Keep original if needed

			realIP := ""

			// Only trust headers if the direct connection is from a trusted proxy
			if isTrusted(parsedRemoteIP) {
				for _, headerName := range config.IPHeaders {
					headerValue := r.Header.Get(headerName)
					if headerValue == "" {
						continue
					}

					// X-Forwarded-For can be a list: client, proxy1, proxy2
					// We want the first one (client).
					ips := strings.Split(headerValue, ",")
					candidateIP := strings.TrimSpace(ips[0])

					parsedCandidate := net.ParseIP(candidateIP)
					if parsedCandidate != nil {
						realIP = candidateIP
						break // Found a valid IP, stop checking headers
					}
				}
			}

			req := r
			if realIP != "" {
				// Successfully found IP via trusted header
				r.RemoteAddr = net.JoinHostPort(realIP, "0") // Update RemoteAddr (use dummy port)
				if config.StoreInContext {
					ctx := context.WithValue(r.Context(), config.ContextKey, realIP)
					req = r.WithContext(ctx)
				}
			} else {
				// Could not find via header, or proxy not trusted. Use original RemoteAddr IP.
				// Store the original parsed IP (without port) in context if configured.
				if config.StoreInContext && parsedRemoteIP != nil {
					ctx := context.WithValue(r.Context(), config.ContextKey, parsedRemoteIP.String())
					req = r.WithContext(ctx)
				}
				// Ensure RemoteAddr is restored if it was temporarily changed
				r.RemoteAddr = originalRemoteAddr
			}

			next.ServeHTTP(w, req)
		})
	}
}

// MaxRequestBodySizeConfig holds configuration for MaxRequestBodySizeMiddleware.
type MaxRequestBodySizeConfig struct {
	// LimitBytes is the maximum allowed size of the request body in bytes. Required.
	LimitBytes int64
	// OnError allows custom handling when the body size limit is exceeded.
	// If nil, sends 413 Request Entity Too Large.
	OnError func(w http.ResponseWriter, r *http.Request)
}

// MaxRequestBodySizeMiddleware limits the size of incoming request bodies.
func MaxRequestBodySizeMiddleware(config MaxRequestBodySizeConfig) Middleware {
	if config.LimitBytes <= 0 {
		panic("MaxRequestBodySizeMiddleware: LimitBytes must be positive")
	}

	onError := config.OnError
	if onError == nil {
		onError = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check ContentLength header first for quick rejection
			if r.ContentLength > config.LimitBytes {
				onError(w, r)
				return
			}
			// Wrap the body with MaxBytesReader
			r.Body = http.MaxBytesReader(w, r.Body, config.LimitBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// TrailingSlashRedirectConfig holds configuration for TrailingSlashRedirectMiddleware.
type TrailingSlashRedirectConfig struct {
	// AddSlash enforces trailing slashes (redirects /path to /path/). Defaults to false (removes slash).
	AddSlash bool
	// RedirectCode is the HTTP status code used for redirection.
	// Defaults to http.StatusMovedPermanently (301). Use 308 for POST/PUT etc.
	RedirectCode int
}

// TrailingSlashRedirectMiddleware redirects requests to add or remove a trailing slash.
func TrailingSlashRedirectMiddleware(config TrailingSlashRedirectConfig) Middleware {
	if config.RedirectCode == 0 {
		config.RedirectCode = http.StatusMovedPermanently
	}
	if config.RedirectCode < 300 || config.RedirectCode >= 400 {
		log.Printf("[WARN] TrailingSlashRedirectMiddleware: Invalid redirect code %d, using 301", config.RedirectCode)
		config.RedirectCode = http.StatusMovedPermanently
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			// Ignore root path "/" and paths with extensions (likely files)
			if len(path) > 1 && !strings.Contains(path[strings.LastIndex(path, "/"):], ".") {
				hasSlash := strings.HasSuffix(path, "/")
				shouldHaveSlash := config.AddSlash

				if hasSlash != shouldHaveSlash {
					newPath := path
					if shouldHaveSlash {
						newPath += "/"
					} else {
						newPath = strings.TrimSuffix(path, "/")
					}
					// Preserve query string
					r.URL.Path = newPath
					http.Redirect(w, r, r.URL.String(), config.RedirectCode)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ForceHTTPSConfig holds configuration for ForceHTTPSMiddleware.
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

// ForceHTTPSMiddleware redirects HTTP requests to HTTPS.
func ForceHTTPSMiddleware(config ForceHTTPSConfig) Middleware {
	if config.RedirectCode == 0 {
		config.RedirectCode = http.StatusMovedPermanently
	}
	if config.RedirectCode < 300 || config.RedirectCode >= 400 {
		log.Printf("[WARN] ForceHTTPSMiddleware: Invalid redirect code %d, using 301", config.RedirectCode)
		config.RedirectCode = http.StatusMovedPermanently
	}
	if config.ForwardedProtoHeader == "" {
		config.ForwardedProtoHeader = "X-Forwarded-Proto"
	}
	trustHeader := true
	if config.TrustForwardedHeader != nil {
		trustHeader = *config.TrustForwardedHeader
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isHTTPS := r.TLS != nil
			if !isHTTPS && trustHeader {
				proto := r.Header.Get(config.ForwardedProtoHeader)
				isHTTPS = strings.EqualFold(proto, "https")
			}

			if !isHTTPS {
				targetHost := config.TargetHost
				if targetHost == "" {
					// Split host and port, ignore error if port isn't present
					reqHost, _, splitErr := net.SplitHostPort(r.Host)
					if splitErr != nil {
						targetHost = r.Host // Use full host if SplitHostPort fails (e.g., no port)
					} else {
						targetHost = reqHost
					}
				}

				targetURL := "https://" + targetHost
				if config.TargetPort > 0 && config.TargetPort != 443 {
					targetURL += ":" + strconv.Itoa(config.TargetPort)
				}
				targetURL += r.URL.RequestURI() // Keep path and query string

				w.Header().Set("Connection", "close")
				http.Redirect(w, r, targetURL, config.RedirectCode)
				return
			}

			// If already HTTPS, proceed. HSTS might be set by SecurityHeadersMiddleware.
			next.ServeHTTP(w, r)
		})
	}
}

// ConcurrencyLimiterConfig holds configuration for ConcurrencyLimiterMiddleware.
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

// ConcurrencyLimiterMiddleware limits the number of concurrent requests.
func ConcurrencyLimiterMiddleware(config ConcurrencyLimiterConfig) Middleware {
	if config.MaxConcurrent <= 0 {
		panic("ConcurrencyLimiterMiddleware: MaxConcurrent must be positive")
	}

	// Use a buffered channel as a semaphore
	sem := make(chan struct{}, config.MaxConcurrent)

	onLimitExceeded := config.OnLimitExceeded
	if onLimitExceeded == nil {
		onLimitExceeded = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service temporarily unavailable (concurrency limit reached)", http.StatusServiceUnavailable)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.WaitTimeout > 0 {
				// Try to acquire semaphore with timeout
				ctx, cancel := context.WithTimeout(r.Context(), config.WaitTimeout)
				defer cancel()
				select {
				case sem <- struct{}{}: // Acquired successfully
					defer func() {
						<-sem // Release semaphore
					}()
					next.ServeHTTP(w, r.WithContext(ctx)) // Pass potentially timed-out context
				case <-ctx.Done(): // Timed out waiting for semaphore
					onLimitExceeded(w, r)
					return
				}
			} else {
				// Wait indefinitely
				sem <- struct{}{} // Acquire semaphore (blocks if full)
				defer func() {
					<-sem // Release semaphore
				}()
				next.ServeHTTP(w, r)
			}
		})
	}
}

// MaintenanceModeConfig holds configuration for MaintenanceModeMiddleware.
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

// MaintenanceModeMiddleware returns a 503 Service Unavailable if enabled, allowing bypass for specific IPs.
func MaintenanceModeMiddleware(config MaintenanceModeConfig) Middleware {
	if config.EnabledFlag == nil {
		panic("MaintenanceModeMiddleware: EnabledFlag is required")
	}
	if config.StatusCode == 0 {
		config.StatusCode = http.StatusServiceUnavailable
	}
	if config.RetryAfterSeconds <= 0 {
		config.RetryAfterSeconds = 300 // 5 minutes
	}
	if config.Message == "" {
		config.Message = "Service temporarily unavailable due to maintenance"
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}

	allowedMap := make(map[string]struct{})
	for _, ipStr := range config.AllowedIPs {
		ipStr = strings.TrimSpace(ipStr)
		ip := net.ParseIP(ipStr)
		if ip != nil {
			allowedMap[ip.String()] = struct{}{}
		} else {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err == nil {
				// Store CIDR string representation for matching later
				allowedMap[ipNet.String()] = struct{}{}
			} else {
				config.Logger.Printf("[WARN] MaintenanceModeMiddleware: Invalid allowed IP/CIDR: %s", ipStr)
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EnabledFlag.Load() { // Check the atomic flag
				clientIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					clientIPStr = r.RemoteAddr // Handle cases without port
				}
				clientIP := net.ParseIP(clientIPStr)

				isAllowed := false
				if clientIP != nil {
					// Check direct IP match
					if _, ok := allowedMap[clientIP.String()]; ok {
						isAllowed = true
					} else {
						// Check CIDR matches
						for cidrStr := range allowedMap {
							if strings.Contains(cidrStr, "/") { // Simple check for CIDR
								_, ipNet, _ := net.ParseCIDR(cidrStr) // Ignore error, already validated
								if ipNet != nil && ipNet.Contains(clientIP) {
									isAllowed = true
									break
								}
							}
						}
					}
				}

				if !isAllowed {
					// Return maintenance response
					w.Header().Set("Retry-After", strconv.Itoa(config.RetryAfterSeconds))
					http.Error(w, config.Message, config.StatusCode)
					return
				}
				// Allowed IP, proceed even in maintenance mode
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IPFilterConfig holds configuration for IP filtering.
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

// IPFilterMiddleware restricts access based on client IP address.
func IPFilterMiddleware(config IPFilterConfig) Middleware {
	if config.Logger == nil {
		config.Logger = log.Default()
	}

	parseCIDRList := func(list []string) []*net.IPNet {
		nets := []*net.IPNet{}
		for _, entry := range list {
			entry = strings.TrimSpace(entry)
			ip := net.ParseIP(entry)
			if ip != nil {
				// Treat single IP as /32 or /128
				mask := net.CIDRMask(128, 128)
				if ip.To4() != nil {
					mask = net.CIDRMask(32, 32)
				}
				nets = append(nets, &net.IPNet{IP: ip, Mask: mask})
			} else {
				_, ipNet, err := net.ParseCIDR(entry)
				if err == nil {
					nets = append(nets, ipNet)
				} else {
					config.Logger.Printf(
						"[WARN] IPFilterMiddleware: Invalid IP/CIDR in filter list: %s",
						entry,
					)
				}
			}
		}
		return nets
	}

	allowedNets := parseCIDRList(config.AllowedIPs)
	blockedNets := parseCIDRList(config.BlockedIPs)

	onForbidden := config.OnForbidden
	if onForbidden == nil {
		onForbidden = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				clientIPStr = r.RemoteAddr // Handle cases without port
			}
			clientIP := net.ParseIP(clientIPStr)
			if clientIP == nil {
				config.Logger.Printf(
					"[WARN] IPFilterMiddleware: Could not parse client IP: %s",
					clientIPStr,
				)
				onForbidden(w, r)
				return
			}

			// Check if the IP is in a blocked network.
			for _, network := range blockedNets {
				if network.Contains(clientIP) {
					onForbidden(w, r)
					return
				}
			}

			// Check if the IP is in an allowed network.
			isAllowed := false
			for _, network := range allowedNets {
				if network.Contains(clientIP) {
					isAllowed = true
					break
				}
			}

			// Decide based on the BlockByDefault setting.
			// If we block by default, the IP must be explicitly allowed.
			// If we don't block by default, the IP is allowed as long as it wasn't blocked.
			if config.BlockByDefault && !isAllowed {
				onForbidden(w, r)
				return
			}

			// If we reach here, the request is permitted.
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiterConfig holds configuration for the simple rate limiter.
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

// visitor tracks request counts and timestamps for rate limiting.
type visitor struct {
	tokens    float64   // Current number of available tokens
	lastToken time.Time // Time when tokens were last refilled
	lastSeen  time.Time // Timestamp of the last request seen (for cleanup)
}

// RateLimitMiddleware provides basic in-memory rate limiting.
// Warning: This simple implementation has limitations:
// - Memory usage can grow indefinitely without CleanupInterval set.
// - Not suitable for distributed systems (limit is per instance).
// - Accuracy decreases slightly under very high concurrency due to locking.
func RateLimitMiddleware(config RateLimiterConfig) Middleware {
	if config.Requests <= 0 || config.Duration <= 0 {
		panic("RateLimitMiddleware: Requests and Duration must be positive")
	}
	if config.Burst <= 0 {
		config.Burst = config.Requests
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}

	keyFunc := config.KeyFunc
	if keyFunc == nil {
		keyFunc = func(r *http.Request) string {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return r.RemoteAddr
			}
			return ip
		}
	}

	onLimitExceeded := config.OnLimitExceeded
	if onLimitExceeded == nil {
		onLimitExceeded = func(w http.ResponseWriter, r *http.Request) {
			// Suggest retrying after the window duration passes
			retryAfter := strconv.Itoa(int(config.Duration.Seconds()))
			w.Header().Set("Retry-After", retryAfter)
			http.Error(
				w,
				http.StatusText(http.StatusTooManyRequests),
				http.StatusTooManyRequests,
			)
		}
	}

	visitors := make(map[string]*visitor)
	var mu sync.Mutex

	// Token bucket parameters
	rate := float64(config.Requests) / config.Duration.Seconds()
	burst := float64(config.Burst)

	// Start cleanup goroutine if interval is set
	if config.CleanupInterval > 0 {
		go func() {
			ticker := time.NewTicker(config.CleanupInterval)
			defer ticker.Stop()
			for range ticker.C {
				mu.Lock()
				now := time.Now()
				for key, v := range visitors {
					// Remove entries idle for more than ~2x the duration (heuristic)
					if now.Sub(v.lastSeen) > config.Duration*2 {
						delete(visitors, key)
					}
				}
				mu.Unlock()
			}
		}()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			if key == "" {
				config.Logger.Printf(
					"[WARN] RateLimitMiddleware: Could not determine rate limit key for request %s %s",
					r.Method,
					r.URL.Path,
				)
				next.ServeHTTP(w, r)
				return
			}

			mu.Lock()
			v, exists := visitors[key]
			now := time.Now()

			if !exists {
				v = &visitor{tokens: burst, lastToken: now}
				visitors[key] = v
			}

			elapsed := now.Sub(v.lastToken)
			tokensToAdd := elapsed.Seconds() * rate
			v.tokens += tokensToAdd
			v.lastToken = now

			if v.tokens > burst {
				v.tokens = burst
			}

			allowed := false
			if v.tokens >= 1 {
				v.tokens--
				allowed = true
				v.lastSeen = now
			}

			mu.Unlock()

			if !allowed {
				onLimitExceeded(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
