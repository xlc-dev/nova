package nova

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

// TestCompilePattern tests that compilePattern splits a pattern string into
// the correct sequence of segments, handling literals, parameters, and regex.
func TestCompilePattern(t *testing.T) {
	tests := []struct {
		pattern string
		want    []segment
	}{
		{"", []segment{}},
		{"/", []segment{}},
		{"foo/bar", []segment{
			{isParam: false, literal: "foo"},
			{isParam: false, literal: "bar"},
		}},
		{"foo/{id}", []segment{
			{isParam: false, literal: "foo"},
			{isParam: true, paramName: "id"},
		}},
		{"{name:[a-z]+}", []segment{
			{
				isParam:   true,
				paramName: "name",
				regex:     regexp.MustCompile("^[a-z]+$"),
			},
		}},
	}

	for _, tt := range tests {
		segs, err := compilePattern(tt.pattern)
		if err != nil {
			t.Errorf("compilePattern(%q) error: %v", tt.pattern, err)
			continue
		}
		if len(segs) != len(tt.want) {
			t.Errorf("compilePattern(%q): got %d segments, want %d",
				tt.pattern, len(segs), len(tt.want))
			continue
		}
		for i, got := range segs {
			want := tt.want[i]
			if got.isParam != want.isParam {
				t.Errorf("%q seg[%d].isParam = %v, want %v",
					tt.pattern, i, got.isParam, want.isParam)
			}
			if got.literal != want.literal {
				t.Errorf("%q seg[%d].literal = %q, want %q",
					tt.pattern, i, got.literal, want.literal)
			}
			if got.paramName != want.paramName {
				t.Errorf("%q seg[%d].paramName = %q, want %q",
					tt.pattern, i, got.paramName, want.paramName)
			}
			if want.regex != nil {
				if got.regex == nil || got.regex.String() != want.regex.String() {
					t.Errorf("%q seg[%d].regex = %v, want %v",
						tt.pattern, i, got.regex, want.regex)
				}
			}
		}
	}
}

// TestJoinPaths tests that joinPaths concatenates two path parts correctly,
// trimming and normalizing slashes.
func TestJoinPaths(t *testing.T) {
	cases := []struct{ a, b, want string }{
		{"", "", ""},
		{"", "foo", "foo"},
		{"foo", "", "foo"},
		{"foo/", "/bar", "foo/bar"},
		{"/foo/", "/bar/", "foo/bar"},
	}
	for _, c := range cases {
		got := joinPaths(c.a, c.b)
		if got != c.want {
			t.Errorf("joinPaths(%q,%q) = %q; want %q",
				c.a, c.b, got, c.want)
		}
	}
}

// TestMatchSegments tests that matchSegments correctly matches a request path
// against compiled segments, extracting parameters and enforcing regex rules.
func TestMatchSegments(t *testing.T) {
	// static
	segs, _ := compilePattern("foo/bar")
	ok, params := matchSegments("/foo/bar", segs)
	if !ok || len(params) != 0 {
		t.Error("static match failed")
	}
	// param
	segs, _ = compilePattern("foo/{id}")
	ok, params = matchSegments("/foo/123", segs)
	if !ok || params["id"] != "123" {
		t.Errorf("param match failed, got %v", params)
	}
	// regex fail
	segs, _ = compilePattern("foo/{id:[0-9]+}")
	ok, _ = matchSegments("/foo/abc", segs)
	if ok {
		t.Error("expected regex match to fail for abc")
	}
}

// TestBasicRouting tests that a simple GET route returns the correct status
// code and body.
func TestBasicRouting(t *testing.T) {
	r := NewRouter()
	r.Get("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("hi"))
	})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("GET", "/hello", nil))
	if rr.Code != http.StatusTeapot || rr.Body.String() != "hi" {
		t.Errorf("GET /hello = %d,%q; want %d,hi",
			rr.Code, rr.Body.String(), http.StatusTeapot)
	}
}

// TestNotFound tests that an unmatched route returns HTTP 404 Not Found.
func TestNotFound(t *testing.T) {
	r := NewRouter()
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("GET", "/nope", nil))
	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /nope = %d; want %d",
			rr.Code, http.StatusNotFound)
	}
}

// TestMethodNotAllowed tests that a request with an invalid method returns
// HTTP 405 Method Not Allowed.
func TestMethodNotAllowed(t *testing.T) {
	r := NewRouter()
	r.Get("/foo", func(w http.ResponseWriter, req *http.Request) {})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("POST", "/foo", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /foo = %d; want %d",
			rr.Code, http.StatusMethodNotAllowed)
	}
}

// TestCustomNotFoundHandler tests that setting a custom 404 handler overrides
// the default Not Found response.
func TestCustomNotFoundHandler(t *testing.T) {
	r := NewRouter()
	r.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("GET", "/x", nil))
	if rr.Code != http.StatusTeapot {
		t.Errorf("custom 404 = %d; want %d",
			rr.Code, http.StatusTeapot)
	}
}

// TestCustomMethodNotAllowed tests that setting a custom 405 handler overrides
// the default Method Not Allowed response.
func TestCustomMethodNotAllowed(t *testing.T) {
	r := NewRouter()
	r.Get("/foo", func(w http.ResponseWriter, req *http.Request) {})
	r.SetMethodNotAllowedHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(499)
	}))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("PUT", "/foo", nil))
	if rr.Code != 499 {
		t.Errorf("custom 405 = %d; want 499", rr.Code)
	}
}

// TestURLParam tests that URLParam correctly retrieves path parameters
// from the request context.
func TestURLParam(t *testing.T) {
	r := NewRouter()
	r.Get("/item/{id}", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(r.URLParam(req, "id")))
	})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("GET", "/item/42", nil))
	if rr.Body.String() != "42" {
		t.Errorf("URLParam = %q; want 42", rr.Body.String())
	}
}

// TestSubrouter tests that a subrouter with a base path correctly matches
// and handles nested routes.
func TestSubrouter(t *testing.T) {
	r := NewRouter()
	sr := r.Subrouter("/api")
	sr.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("pong"))
	})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("GET", "/api/ping", nil))
	if rr.Body.String() != "pong" {
		t.Errorf("subrouter = %q; want pong", rr.Body.String())
	}
}

// TestGroupPrefixAndMiddleware tests that grouping routes under a prefix and
// applying middleware works as expected, modifying headers and invoking handlers.
func TestGroupPrefixAndMiddleware(t *testing.T) {
	r := NewRouter()
	g := r.Group("/g", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-G", "1")
			next.ServeHTTP(w, req)
		})
	})
	g.Get("/test", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest("GET", "/g/test", nil))
	if rr.Header().Get("X-G") != "1" || rr.Body.String() != "ok" {
		t.Errorf("group prefix/middleware failed: header=%q body=%q",
			rr.Header().Get("X-G"), rr.Body.String())
	}
}
