package nova

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/mail"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ResponseContext provides helper methods for sending HTTP responses with reduced boilerplate.
// It wraps the http.ResponseWriter and http.Request along with router context to provide
// convenient methods for JSON, HTML, text responses, and automatic data binding.
type ResponseContext struct {
	// w is the underlying HTTP response writer for sending responses to the client.
	w http.ResponseWriter
	// r is the incoming HTTP request containing headers, body, and URL parameters.
	r *http.Request
	// router is the parent router instance used for URL parameter extraction.
	router *Router
}

// HandlerFunc is an enhanced handler function that receives a ResponseContext and returns an error.
// This allows for cleaner error handling and response management compared to standard http.HandlerFunc.
type HandlerFunc func(*ResponseContext) error

// Middleware defines the function signature for middleware.
// A middleware is a function that wraps an http.Handler, adding extra behavior
// such as logging, authentication, or request modification.
type Middleware func(http.Handler) http.Handler

// Router is a minimal HTTP router that supports dynamic routes with regex validation
// in path parameters and can mount subrouters. It provides middleware support,
// custom error handlers, and both traditional and enhanced handler functions.
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

// Group is a lightweight helper that allows users to register a set of routes
// that share a common prefix and/or middleware. It delegates to the parent router
// while applying its own prefix and middleware chain.
type Group struct {
	// prefix is the URL path prefix applied to all routes in this group.
	prefix string
	// router is the parent router that actually stores the routes.
	router *Router
	// middlewares contains middleware functions applied only to routes in this group.
	middlewares []Middleware
}

// route represents an individual route with its compiled URL pattern and associated handler.
// It stores the HTTP method, handler function, compiled path segments, and optional metadata.
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

// segment represents a part of the URL path. It may be a literal string or a dynamic
// parameter with an optional regex pattern for validation.
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

// Translation messages for validation errors.
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

// detectLanguage extracts the preferred language from Accept-Language header
func detectLanguage(acceptLanguage string) string {
	if acceptLanguage == "" {
		return "en"
	}

	// Parse Accept-Language header (simplified)
	languages := strings.SplitSeq(acceptLanguage, ",")
	for lang := range languages {
		// Remove quality value if present (e.g., "en;q=0.9")
		langCode := strings.TrimSpace(strings.Split(lang, ";")[0])
		// Extract just the language part (e.g., "en" from "en-US")
		langCode = strings.Split(langCode, "-")[0]

		if _, exists := validationMessages[langCode]; exists {
			return langCode
		}
	}

	return "en"
}

// getMessage gets a localized validation message
func getMessage(lang, key string, args ...any) string {
	if messages, exists := validationMessages[lang]; exists {
		if template, exists := messages[key]; exists {
			return fmt.Sprintf(template, args...)
		}
	}

	// Fallback to English
	if messages, exists := validationMessages["en"]; exists {
		if template, exists := messages[key]; exists {
			return fmt.Sprintf(template, args...)
		}
	}

	// Ultimate fallback
	return fmt.Sprintf("Validation error for field '%v'", args[0])
}

// rebuildChain reconstructs the composed middleware chain based on the
// currently registered middlewares. Middleware is applied in reverse order
// so that the first registered middleware wraps the outermost layer.
func (r *Router) rebuildChain() {
	r.chain = func(finalHandler http.Handler) http.Handler {
		for i := len(r.middlewares) - 1; i >= 0; i-- {
			finalHandler = r.middlewares[i](finalHandler)
		}
		return finalHandler
	}
}

// compilePattern converts a URL pattern string into a slice of segments.
// Parameters are declared as {name} or {name:regex}. In the latter case,
// the regex is precompiled and validated during route registration.
func compilePattern(pattern string) ([]segment, error) {
	trimmed := strings.Trim(pattern, "/")
	if trimmed == "" {
		return []segment{}, nil
	}

	parts := strings.Split(trimmed, "/")
	segs := make([]segment, 0, len(parts))
	for _, part := range parts {
		if len(part) >= 2 && part[0] == '{' && part[len(part)-1] == '}' {
			inner := part[1 : len(part)-1]
			subparts := strings.SplitN(inner, ":", 2)
			paramName := subparts[0]
			var re *regexp.Regexp
			if len(subparts) == 2 {
				// Compile the regex and add anchors.
				pat := "^" + subparts[1] + "$"
				rx, err := regexp.Compile(pat)
				if err != nil {
					return nil, err
				}
				re = rx
			}
			segs = append(segs, segment{
				isParam:   true,
				paramName: paramName,
				regex:     re,
			})
		} else {
			segs = append(segs, segment{
				isParam: false,
				literal: part,
			})
		}
	}
	return segs, nil
}

// joinPaths concatenates two URL path segments ensuring exactly one "/" between them.
// It handles edge cases where either path is empty and normalizes trailing/leading slashes.
func joinPaths(a, b string) string {
	a = strings.Trim(a, "/")
	b = strings.Trim(b, "/")
	if a == "" && b == "" {
		return ""
	}
	if a == "" {
		return b
	}
	if b == "" {
		return a
	}
	return a + "/" + b
}

// matchSegments checks if the given URL path matches the compiled segments.
// If the path matches, it returns true and a map of parameter names to their values.
// Parameter validation is performed using regex patterns if specified.
func matchSegments(path string, segments []segment) (bool, map[string]string) {
	trimmed := strings.Trim(path, "/")
	var pathParts []string
	if trimmed == "" {
		pathParts = []string{}
	} else {
		pathParts = strings.Split(trimmed, "/")
	}

	if len(pathParts) != len(segments) {
		return false, nil
	}

	params := make(map[string]string)
	for i, seg := range segments {
		part := pathParts[i]
		if seg.isParam {
			if seg.regex != nil && !seg.regex.MatchString(part) {
				return false, nil
			}
			params[seg.paramName] = part
		} else if seg.literal != part {
			return false, nil
		}
	}
	return true, params
}

// bindFormToStruct uses reflection to bind form values to struct fields.
// Field names are matched using JSON tags (without omitempty) or struct field names.
// Supports automatic type conversion for string, bool, int, and float types.
func bindFormToStruct(form url.Values, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := range rv.NumField() {
		field := rv.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := rt.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		name := strings.Split(jsonTag, ",")[0]
		if name == "" {
			name = fieldType.Name
		}

		value := form.Get(name)
		if err := setFieldValue(field, value); err != nil {
			return fmt.Errorf("field %s: %w", name, err)
		}
	}

	return nil
}

// setFieldValue sets a reflect.Value based on its type and the string value from form data.
// Handles automatic type conversion for common types including special checkbox handling.
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		// Handle checkbox values (HTML checkboxes send "on" when checked, nothing when unchecked)
		if value == "on" || value == "true" || value == "1" {
			field.SetBool(true)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value != "" {
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value != "" {
			i, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if value != "" {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		}
	case reflect.Slice:
		// Handle multiple form values (e.g., checkboxes with same name)
		if field.Type().Elem().Kind() == reflect.String {
			values := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(values), len(values))
			for i, v := range values {
				slice.Index(i).SetString(strings.TrimSpace(v))
			}
			field.Set(slice)
		}
	}
	return nil
}

// validateStruct walks a struct’s exported fields, applies all tag‐based
// validations, and returns every violation found rather than failing fast.
//
// val may be a struct or a pointer to struct. lang selects the locale for
// error messages.
func validateStruct(val any, lang string) error {
	v := reflect.ValueOf(val)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fmt.Errorf("nil value")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	var allErrs ValidationErrors
	t := v.Type()
	for _, f := range reflect.VisibleFields(t) {
		if f.PkgPath != "" { // skip unexported
			continue
		}
		parts := strings.Split(f.Tag.Get("json"), ",")
		name := parts[0]
		if name == "" || name == "-" {
			name = f.Name
		}
		fv := v.FieldByIndex(f.Index)
		required := !slices.Contains(parts[1:], "omitempty")
		custom := f.Tag.Get("error")

		if required && fv.IsZero() {
			if custom != "" {
				allErrs = append(allErrs, fmt.Errorf("%s", custom))
			} else {
				allErrs = append(allErrs,
					fmt.Errorf("%s", getMessage(lang, "required", name)))
			}
			continue
		}

		if errs := validateField(fv, f, name, lang, custom); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
	}

	if len(allErrs) > 0 {
		return allErrs
	}
	return nil
}

// validateField routes to the appropriate validator based on kind,
// returning all errors for that field.
func validateField(
	fv reflect.Value,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors {
	switch fv.Kind() {
	case reflect.String:
		return validateStringField(fv.String(), f, fieldName, lang, custom)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return validateNumericField(float64(fv.Int()), f, fieldName, lang, custom)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return validateNumericField(float64(fv.Uint()), f, fieldName, lang, custom)
	case reflect.Float32, reflect.Float64:
		return validateNumericField(fv.Float(), f, fieldName, lang, custom)
	case reflect.Slice, reflect.Array:
		return validateSliceField(fv, f, fieldName, lang, custom)
	case reflect.Struct:
		if err := validateStruct(fv.Addr().Interface(), lang); err != nil {
			return ValidationErrors{err}
		}
	case reflect.Pointer:
		if !fv.IsNil() {
			if err := validateStruct(fv.Interface(), lang); err != nil {
				return ValidationErrors{err}
			}
		}
	}
	return nil
}

// validateStringField applies minlength, maxlength, pattern, enum and format
// checks on s, returning every violation.
func validateStringField(
	s string,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors {
	var errs ValidationErrors

	// minlength / maxlength
	if minTag := f.Tag.Get("minlength"); minTag != "" {
		if min, _ := strconv.Atoi(minTag); len(s) < min {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "minlength", fieldName, min)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}
	if maxTag := f.Tag.Get("maxlength"); maxTag != "" {
		if max, _ := strconv.Atoi(maxTag); len(s) > max {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "maxlength", fieldName, max)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}

	// pattern
	if pattern := f.Tag.Get("pattern"); pattern != "" && s != "" {
		if matched, _ := regexp.MatchString(pattern, s); !matched {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "pattern", fieldName)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}

	// enum
	if enum := f.Tag.Get("enum"); enum != "" && s != "" {
		allowed := strings.Split(enum, "|")
		if !slices.Contains(allowed, s) {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "enum",
					fieldName, strings.Join(allowed, ", "))
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}

	// format‐based checks
	if s != "" {
		switch f.Tag.Get("format") {
		case "email":
			if _, err := mail.ParseAddress(s); err != nil {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "email", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "url":
			if u, err := url.ParseRequestURI(s); err != nil ||
				u.Scheme == "" || u.Host == "" {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "url", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "uuid":
			// RFC-4122 v4 UUID
			if matched, _ := regexp.MatchString(
				`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ABab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
				s,
			); !matched {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "uuid", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "date-time":
			if _, err := time.Parse(time.RFC3339, s); err != nil {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "date-time", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "date":
			if _, err := time.Parse("2006-01-02", s); err != nil {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "date", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "time":
			if _, err := time.Parse("15:04:05", s); err != nil {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "time", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "password":
			if len(s) < 8 {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "password", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "phone":
			if matched, _ := regexp.MatchString(
				`^[\+]?[1-9][\d\s\-\(\)]{7,15}$`,
				s,
			); !matched {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "phone", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "alphanumeric":
			if matched, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, s); !matched {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "alphanumeric", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "alpha":
			if matched, _ := regexp.MatchString(`^[A-Za-z]+$`, s); !matched {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "alpha", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		case "numeric":
			if matched, _ := regexp.MatchString(`^[0-9]+$`, s); !matched {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "numeric", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		}
	}

	return errs
}

// validateNumericField applies min, max, and multipleOf checks on num,
// returning every violation.
func validateNumericField(
	num float64,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors {
	var errs ValidationErrors

	if minTag := f.Tag.Get("min"); minTag != "" {
		if min, _ := strconv.ParseFloat(minTag, 64); num < min {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "min", fieldName, min)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}
	if maxTag := f.Tag.Get("max"); maxTag != "" {
		if max, _ := strconv.ParseFloat(maxTag, 64); num > max {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "max", fieldName, max)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}
	if multTag := f.Tag.Get("multipleOf"); multTag != "" {
		if mult, _ := strconv.ParseFloat(multTag, 64); mult != 0 {
			if rem := num - (mult * float64(int(num/mult))); rem != 0 {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "multipleOf",
						fieldName, mult)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
		}
	}
	return errs
}

// validateSliceField applies minItems, maxItems, uniqueItems and recurses
// into elements, returning every violation.
func validateSliceField(
	fv reflect.Value,
	f reflect.StructField,
	fieldName, lang, custom string,
) ValidationErrors {
	var errs ValidationErrors
	length := fv.Len()

	if minTag := f.Tag.Get("minItems"); minTag != "" {
		if min, _ := strconv.Atoi(minTag); length < min {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "minItems", fieldName, min)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}
	if maxTag := f.Tag.Get("maxItems"); maxTag != "" {
		if max, _ := strconv.Atoi(maxTag); length > max {
			msg := custom
			if msg == "" {
				msg = getMessage(lang, "maxItems", fieldName, max)
			}
			errs = append(errs, fmt.Errorf("%s", msg))
		}
	}
	if f.Tag.Get("uniqueItems") == "true" {
		seen := map[any]bool{}
		for i := range length {
			v := fv.Index(i).Interface()
			if seen[v] {
				msg := custom
				if msg == "" {
					msg = getMessage(lang, "uniqueItems", fieldName)
				}
				errs = append(errs, fmt.Errorf("%s", msg))
			}
			seen[v] = true
		}
	}
	for i := range length {
		item := fv.Index(i)
		if item.Kind() == reflect.Struct ||
			(item.Kind() == reflect.Pointer && !item.IsNil()) {
			if sub := validateStruct(item.Interface(), lang); sub != nil {
				errs = append(errs,
					fmt.Errorf("%s[%d]: %w", fieldName, i, sub))
			}
		}
	}
	return errs
}

// Static mounts an fs.FS under the given URL prefix.
func (r *Router) Static(urlPathPrefix string, subFS fs.FS) {
	prefix := "/" + strings.Trim(urlPathPrefix, "/")
	if prefix == "/" {
		prefix = ""
	}
	routePattern := path.Join(prefix, "{filepath:.*}")

	stripPrefix := prefix + "/"
	if stripPrefix == "/" {
		stripPrefix = "/"
	}

	fileServer := http.FileServer(http.FS(subFS))
	fsHandler := http.StripPrefix(stripPrefix, fileServer)
	hf := http.HandlerFunc(fsHandler.ServeHTTP)

	r.Handle(http.MethodGet, routePattern, hf)
	r.Handle(http.MethodHead, routePattern, hf)
}

// ValidationErrors collects one or more validation violations and implements error.
type ValidationErrors []error

// Error joins all contained error messages into a single string.
func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, err := range ve {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// JSON sends a JSON response with the given status code and data.
// It automatically sets the Content-Type header to "application/json" and
// handles JSON encoding. Returns an error if encoding fails.
func (rc *ResponseContext) JSON(statusCode int, data any) error {
	rc.w.Header().Set("Content-Type", "application/json")
	rc.w.WriteHeader(statusCode)
	return json.NewEncoder(rc.w).Encode(data)
}

// JSONError sends a JSON error response with the given status code and message.
// It creates a standardized error response format with an "error" field.
func (rc *ResponseContext) JSONError(statusCode int, message string) error {
	return rc.JSON(statusCode, map[string]string{"error": message})
}

// HTML sends an HTML response with the given status code and content.
// It automatically sets the Content-Type header to "text/html; charset=utf-8"
// and renders the provided HTMLElement to string.
func (rc *ResponseContext) HTML(statusCode int, content HTMLElement) error {
	rc.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	rc.w.WriteHeader(statusCode)
	_, err := rc.w.Write([]byte(content.Render()))
	return err
}

// Text sends a plain text response with the given status code and text content.
// It automatically sets the Content-Type header to "text/plain; charset=utf-8".
func (rc *ResponseContext) Text(statusCode int, text string) error {
	rc.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rc.w.WriteHeader(statusCode)
	_, err := rc.w.Write([]byte(text))
	return err
}

// Redirect sends an HTTP redirect response with the given status code and URL.
// Common status codes are 301 (permanent), 302 (temporary), and 307 (temporary, preserve method).
func (rc *ResponseContext) Redirect(statusCode int, url string) error {
	rc.w.Header().Set("Location", url)
	rc.w.WriteHeader(statusCode)
	return nil
}

// URLParam retrieves the URL parameter for the given key from the request context.
// Returns an empty string if the parameter is not present or if the request
// doesn't contain URL parameters.
func (rc *ResponseContext) URLParam(key string) string {
	return rc.router.URLParam(rc.r, key)
}

// Request returns the underlying http.Request for advanced request handling
// when the ResponseContext helpers are not sufficient.
func (rc *ResponseContext) Request() *http.Request {
	return rc.r
}

// Writer returns the underlying http.ResponseWriter for advanced response handling
// when the ResponseContext helpers are not sufficient.
func (rc *ResponseContext) Writer() http.ResponseWriter {
	return rc.w
}

// WantsJSON returns true if the request expects a JSON response based on
// Content-Type or Accept headers. Useful for dual HTML/JSON endpoints.
func (rc *ResponseContext) WantsJSON() bool {
	contentType := rc.r.Header.Get("Content-Type")
	accept := rc.r.Header.Get("Accept")
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(accept, "application/json")
}

// Bind automatically unmarshals and validates request data into the provided struct.
// It supports both JSON and form data, with automatic content-type detection.
// The struct should be a pointer. Validation is performed if validation middleware is active.
func (rc *ResponseContext) Bind(v any) error {
	contentType := rc.r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		return rc.BindJSON(v)
	}
	return rc.BindForm(v)
}

// BindJSON unmarshals JSON request body into the provided struct.
// The struct should be a pointer. Returns an error if the body is nil or JSON is invalid.
func (rc *ResponseContext) BindJSON(v any) error {
	if rc.r.Body == nil {
		return fmt.Errorf("request body is nil")
	}
	return json.NewDecoder(rc.r.Body).Decode(v)
}

// BindForm binds form data into the provided struct using reflection.
// The struct should be a pointer. Field names are matched using JSON tags or struct field names.
// Supports string, bool, int, and float fields with automatic type conversion.
func (rc *ResponseContext) BindForm(v any) error {
	if err := rc.r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	return bindFormToStruct(rc.r.Form, v)
}

// BindValidated binds and validates request data (JSON or form) with comprehensive validation.
func (rc *ResponseContext) BindValidated(v any) error {
	// Detect language from Accept-Language header
	lang := detectLanguage(rc.r.Header.Get("Accept-Language"))

	contentType := rc.r.Header.Get("Content-Type")

	// Bind the data
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(rc.r.Body).Decode(v); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
	} else {
		// Handle form data
		if err := rc.r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %w", err)
		}
		if err := bindFormToStruct(rc.r.Form, v); err != nil {
			return fmt.Errorf("failed to bind form data: %w", err)
		}
	}

	// Validate the bound data
	if err := validateStruct(v, lang); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// NewRouter creates and returns a new Router instance with default configuration.
func NewRouter() *Router {
	r := &Router{
		routes:      make([]route, 0),
		subrouters:  make([]*Router, 0),
		paramsKey:   new(struct{}),
		middlewares: make([]Middleware, 0),
		basePath:    "",
		chain: func(h http.Handler) http.Handler {
			return h
		},
	}
	return r
}

// SetNotFoundHandler allows users to supply a custom http.Handler
// for requests that do not match any route. If not set, the standard
// http.NotFound handler is used.
func (r *Router) SetNotFoundHandler(h http.Handler) {
	r.notFoundHandler = h
}

// SetMethodNotAllowedHandler allows users to supply a custom http.Handler
// for requests that match a route pattern but use an unsupported HTTP method.
// If not set, a default 405 Method Not Allowed response is returned.
func (r *Router) SetMethodNotAllowedHandler(h http.Handler) {
	r.methodNotAllowedHandler = h
}

// Use registers one or more middleware functions that will be applied globally
// to every matched route handler. Middleware is applied in the order it is registered.
func (r *Router) Use(mws ...Middleware) {
	r.middlewares = append(r.middlewares, mws...)
	r.rebuildChain()
}

// Handle registers a new route with the given HTTP method, URL pattern, and handler.
// If the router has a non-empty basePath, it is automatically prepended to the pattern.
// Optional RouteOptions can be provided for OpenAPI documentation.
func (r *Router) Handle(method, pattern string, handler http.HandlerFunc, opts ...*RouteOptions) {
	fullPattern := pattern

	if r.basePath != "" {
		fullPattern = joinPaths(r.basePath, pattern)
	}

	segs, err := compilePattern(fullPattern)
	if err != nil {
		panic("invalid route pattern: " + err.Error())
	}

	var routeOpts *RouteOptions
	if len(opts) > 0 {
		routeOpts = opts[0]
	}
	r.routes = append(r.routes, route{
		method:   method,
		handler:  handler,
		segments: segs,
		options:  routeOpts,
	})
}

// HandleFunc registers a new enhanced route that receives a ResponseContext instead
// of separate ResponseWriter and Request parameters. This enables cleaner error handling
// and response management with automatic data binding capabilities.
func (r *Router) HandleFunc(method, pattern string, handler HandlerFunc, opts ...*RouteOptions) {
	r.Handle(method, pattern, func(w http.ResponseWriter, req *http.Request) {
		rc := &ResponseContext{
			w:      w,
			r:      req,
			router: r,
		}

		if err := handler(rc); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}, opts...)
}

// Get registers a new route for HTTP GET requests using the standard handler signature.
func (r *Router) Get(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	r.Handle(http.MethodGet, pattern, handler, options...)
}

// GetFunc registers a new route for HTTP GET requests using the enhanced handler signature.
func (r *Router) GetFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	r.HandleFunc(http.MethodGet, pattern, handler, options...)
}

// Post registers a new route for HTTP POST requests using the standard handler signature.
func (r *Router) Post(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	r.Handle(http.MethodPost, pattern, handler, options...)
}

// PostFunc registers a new route for HTTP POST requests using the enhanced handler signature.
func (r *Router) PostFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	r.HandleFunc(http.MethodPost, pattern, handler, options...)
}

// Put registers a new route for HTTP PUT requests using the standard handler signature.
func (r *Router) Put(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	r.Handle(http.MethodPut, pattern, handler, options...)
}

// PutFunc registers a new route for HTTP PUT requests using the enhanced handler signature.
func (r *Router) PutFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	r.HandleFunc(http.MethodPut, pattern, handler, options...)
}

// Patch registers a new route for HTTP PATCH requests using the standard handler signature.
func (r *Router) Patch(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	r.Handle(http.MethodPatch, pattern, handler, options...)
}

// PatchFunc registers a new route for HTTP PATCH requests using the enhanced handler signature.
func (r *Router) PatchFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	r.HandleFunc(http.MethodPatch, pattern, handler, options...)
}

// Delete registers a new route for HTTP DELETE requests using the standard handler signature.
func (r *Router) Delete(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	r.Handle(http.MethodDelete, pattern, handler, options...)
}

// DeleteFunc registers a new route for HTTP DELETE requests using the enhanced handler signature.
func (r *Router) DeleteFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	r.HandleFunc(http.MethodDelete, pattern, handler, options...)
}

// Subrouter creates a new Router mounted at the given prefix.
// The subrouter inherits the parent's context key, global middleware,
// and error handlers. Its routes will automatically receive the combined prefix.
func (r *Router) Subrouter(prefix string) *Router {
	newRouter := &Router{
		routes:                  make([]route, 0),
		subrouters:              make([]*Router, 0),
		paramsKey:               r.paramsKey,
		middlewares:             slices.Clone(r.middlewares),
		basePath:                joinPaths(r.basePath, prefix),
		notFoundHandler:         r.notFoundHandler,
		methodNotAllowedHandler: r.methodNotAllowedHandler,
	}
	newRouter.rebuildChain()
	r.subrouters = append(r.subrouters, newRouter)
	return newRouter
}

// Group creates and returns a new Group with the given prefix. A group is a lightweight
// convenience wrapper that prefixes routes and can add its own middleware without
// creating a separate router instance.
func (r *Router) Group(prefix string, mws ...Middleware) *Group {
	fullPrefix := joinPaths(r.basePath, prefix)
	return &Group{
		prefix:      fullPrefix,
		router:      r,
		middlewares: mws,
	}
}

// ServeHTTP implements the http.Handler interface.
// It first checks subrouters based on their base path, then its own routes.
// If a URL pattern matches but the HTTP method does not, a 405 Method Not Allowed is returned.
// If no routes match, a 404 Not Found is returned.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestPath := req.URL.Path

	// Check mounted subrouters first.
	for _, sr := range r.subrouters {
		if sr.basePath != "" && strings.HasPrefix(requestPath, "/"+strings.Trim(sr.basePath, "/")) {
			sr.ServeHTTP(w, req)
			return
		}
	}

	var patternMatched bool
	for _, rt := range r.routes {
		if ok, params := matchSegments(req.URL.Path, rt.segments); ok {
			patternMatched = true
			if rt.method == req.Method {
				if len(params) > 0 {
					ctx := context.WithValue(req.Context(), r.paramsKey, params)
					req = req.WithContext(ctx)
				}
				finalHandler := r.chain(rt.handler)
				finalHandler.ServeHTTP(w, req)
				return
			}
		}
	}

	if patternMatched {
		if r.methodNotAllowedHandler != nil {
			r.methodNotAllowedHandler.ServeHTTP(w, req)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
		}
		return
	}

	if r.notFoundHandler != nil {
		r.notFoundHandler.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// URLParam retrieves the URL parameter for the given key from the request context.
// It returns an empty string if the parameter is not present. This method should
// only be called within route handlers that were registered with this router.
func (r *Router) URLParam(req *http.Request, key string) string {
	params, ok := req.Context().Value(r.paramsKey).(map[string]string)
	if !ok {
		return ""
	}
	return params[key]
}

// Use adds middleware functions to the group. These middleware functions apply
// only to routes registered through the group, not to the parent router.
func (g *Group) Use(mws ...Middleware) {
	g.middlewares = append(g.middlewares, mws...)
}

// Handle registers a new route within the group, applying the group's prefix and middleware.
// The route is ultimately registered with the parent router after transformations.
func (g *Group) Handle(method, pattern string, handler http.HandlerFunc, opts ...*RouteOptions) {
	fullPattern := joinPaths(g.prefix, pattern)

	// wrap the handler with the group's middleware
	h := http.Handler(handler)
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		h = g.middlewares[i](h)
	}

	// forward to the underlying router.Handle, including opts
	g.router.Handle(method, fullPattern,
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		},
		opts...,
	)
}

// HandleFunc registers a new enhanced route within the group, applying the group's
// prefix and middleware. Uses the enhanced handler signature for better error handling.
func (g *Group) HandleFunc(method, pattern string, handler HandlerFunc, opts ...*RouteOptions) {
	fullPattern := joinPaths(g.prefix, pattern)

	// Create enhanced handler
	enhancedHandler := func(w http.ResponseWriter, req *http.Request) {
		rc := &ResponseContext{
			w:      w,
			r:      req,
			router: g.router,
		}

		if err := handler(rc); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	// wrap with group middleware
	h := http.Handler(http.HandlerFunc(enhancedHandler))
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		h = g.middlewares[i](h)
	}

	g.router.Handle(method, fullPattern,
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		},
		opts...,
	)
}

// Get registers a new route for HTTP GET requests within the group.
func (g *Group) Get(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	g.Handle(http.MethodGet, pattern, handler, options...)
}

// GetFunc registers a new enhanced route for HTTP GET requests within the group.
func (g *Group) GetFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	g.HandleFunc(http.MethodGet, pattern, handler, options...)
}

// Post registers a new route for HTTP POST requests within the group.
func (g *Group) Post(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	g.Handle(http.MethodPost, pattern, handler, options...)
}

// PostFunc registers a new enhanced route for HTTP POST requests within the group.
func (g *Group) PostFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	g.HandleFunc(http.MethodPost, pattern, handler, options...)
}

// Put registers a new route for HTTP PUT requests within the group.
func (g *Group) Put(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	g.Handle(http.MethodPut, pattern, handler, options...)
}

// PutFunc registers a new enhanced route for HTTP PUT requests within the group.
func (g *Group) PutFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	g.HandleFunc(http.MethodPut, pattern, handler, options...)
}

// Patch registers a new route for HTTP PATCH requests within the group.
func (g *Group) Patch(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	g.Handle(http.MethodPatch, pattern, handler, options...)
}

// PatchFunc registers a new enhanced route for HTTP PATCH requests within the group.
func (g *Group) PatchFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	g.HandleFunc(http.MethodPatch, pattern, handler, options...)
}

// Delete registers a new route for HTTP DELETE requests within the group.
func (g *Group) Delete(pattern string, handler http.HandlerFunc, options ...*RouteOptions) {
	g.Handle(http.MethodDelete, pattern, handler, options...)
}

// DeleteFunc registers a new enhanced route for HTTP DELETE requests within the group.
func (g *Group) DeleteFunc(pattern string, handler HandlerFunc, options ...*RouteOptions) {
	g.HandleFunc(http.MethodDelete, pattern, handler, options...)
}
