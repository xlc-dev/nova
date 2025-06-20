package nova

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"
)

// OpenAPIConfig holds metadata for generating an OpenAPI specification.
type OpenAPIConfig struct {
	Title       string
	Version     string
	Description string
}

// RouteOptions holds OpenAPI metadata for a single route.
// Tags, Summary, Description, and OperationID map directly into the
// corresponding Operation fields. RequestBody, Responses, and Parameters
// drive schema generation for request bodies, responses, and parameters.
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

// ResponseOption configures a single HTTP response in an Operation.
// Description is required; Body, if non-nil, is used to generate a schema.
type ResponseOption struct {
	Description string
	Body        any
}

// ParameterOption configures an Operation parameter.
// Name and In ("query", "header", "path", "cookie") are required.
// Required, Schema and Example further customize the generated parameter.
type ParameterOption struct {
	Name        string
	In          string
	Description string
	Required    bool
	Schema      any
	Example     any
}

// OpenAPI is the root document object for an OpenAPI 3 specification.
type OpenAPI struct {
	OpenAPI    string               `json:"openapi"`
	Info       Info                 `json:"info"`
	Paths      map[string]*PathItem `json:"paths"`
	Components *Components          `json:"components,omitempty"`
}

// Info provides metadata about the API: title, version, and optional description.
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// PathItem describes the operations available on a single API path.
// One of Get, Post, Put, Delete, Patch may be non-nil.
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

// Operation describes a single API operation on a path.
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

// ParameterObject describes a single parameter for an Operation or PathItem.
type ParameterObject struct {
	Name        string        `json:"name"`
	In          string        `json:"in"`
	Description string        `json:"description,omitempty"`
	Required    bool          `json:"required"`
	Deprecated  bool          `json:"deprecated,omitempty"`
	Schema      *SchemaObject `json:"schema"`
	Example     any           `json:"example,omitempty"`
}

// RequestBodyObject describes a request body for an Operation.
type RequestBodyObject struct {
	Description string                      `json:"description,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content"`
	Required    bool                        `json:"required,omitempty"`
}

// ResponseObject describes a single response in an Operation.
type ResponseObject struct {
	Description string                      `json:"description"`
	Headers     map[string]*HeaderObject    `json:"headers,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content,omitempty"`
}

// MediaTypeObject holds the schema defining the media type for a request or response.
type MediaTypeObject struct {
	Schema *SchemaObject `json:"schema,omitempty"`
}

// SchemaObject represents an OpenAPI Schema or a reference to one.
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

// Components holds reusable schema definitions for the OpenAPI spec.
type Components struct {
	Schemas map[string]*SchemaObject `json:"schemas,omitempty"`
}

// HeaderObject describes a response header in an Operation.
type HeaderObject struct {
	Description string        `json:"description,omitempty"`
	Schema      *SchemaObject `json:"schema"`
}

// schemaGenCtx tracks state during schema generation, to de-duplicate
// component schemas and reuse references.
type schemaGenCtx struct {
	componentsSchemas map[string]*SchemaObject
	generatedNames    map[reflect.Type]string
}

// newSchemaGenCtx creates and returns an empty schemaGenCtx to drive
// generateSchema and collect component definitions.
func newSchemaGenCtx() *schemaGenCtx {
	return &schemaGenCtx{
		componentsSchemas: make(map[string]*SchemaObject),
		generatedNames:    make(map[reflect.Type]string),
	}
}

// collectRoutes traverses the router hierarchy starting at r, and populates
// the OpenAPI spec with PathItems and Operations based on registered routes.
func collectRoutes(r *Router, spec *OpenAPI, schemaCtx *schemaGenCtx, parentPath string) {
	for _, route := range r.routes {
		fullPath := buildPathString(route.segments, parentPath)
		if fullPath == "" {
			fullPath = "/"
		} else if !strings.HasPrefix(fullPath, "/") {
			fullPath = "/" + fullPath
		}

		pathItem, exists := spec.Paths[fullPath]
		if !exists {
			pathItem = &PathItem{}
			spec.Paths[fullPath] = pathItem
		}

		op := buildOperation(route, schemaCtx)

		switch route.method {
		case http.MethodGet:
			pathItem.Get = op
		case http.MethodPost:
			pathItem.Post = op
		case http.MethodPut:
			pathItem.Put = op
		case http.MethodDelete:
			pathItem.Delete = op
		case http.MethodPatch:
			pathItem.Patch = op
		}
	}

	for _, sr := range r.subrouters {
		collectRoutes(sr, spec, schemaCtx, sr.basePath)
	}
}

// buildPathString constructs a URI path string from segments and an optional
// parent prefix. Parameters are wrapped in curly braces, e.g., "{id}".
func buildPathString(segments []segment, prefix string) string {
	var parts []string
	if prefix != "" {
		parts = append(parts, strings.Trim(prefix, "/"))
	}
	for _, seg := range segments {
		if seg.isParam {
			parts = append(parts, fmt.Sprintf("{%s}", seg.paramName))
		} else {
			parts = append(parts, seg.literal)
		}
	}
	finalParts := []string{}
	for _, p := range parts {
		if p != "" {
			finalParts = append(finalParts, p)
		}
	}
	return strings.Join(finalParts, "/")
}

// buildOperation creates an OpenAPI OperationObject from a route and schema
// generation context, adding parameters, requestBody, and responses.
func buildOperation(route route, schemaCtx *schemaGenCtx) *Operation {
	op := &Operation{
		Responses:  make(map[string]*ResponseObject),
		Parameters: []ParameterObject{},
	}

	if route.options != nil {
		opts := route.options
		op.Tags = opts.Tags
		op.Summary = opts.Summary
		op.Description = opts.Description
		op.OperationID = opts.OperationID
		op.Deprecated = opts.Deprecated

		if opts.RequestBody != nil {
			op.RequestBody = &RequestBodyObject{
				Required: true,
				Content: map[string]*MediaTypeObject{
					"application/json": {
						Schema: generateSchema(opts.RequestBody, schemaCtx),
					},
				},
			}
		}

		for statusCode, respOpt := range opts.Responses {
			resp := &ResponseObject{Description: respOpt.Description}
			if respOpt.Body != nil {
				resp.Content = map[string]*MediaTypeObject{
					"application/json": {
						Schema: generateSchema(respOpt.Body, schemaCtx),
					},
				}
			}
			op.Responses[fmt.Sprintf("%d", statusCode)] = resp
		}

		for _, paramOpt := range opts.Parameters {
			required := paramOpt.Required
			if paramOpt.In == "path" {
				required = true
			}
			paramObj := ParameterObject{
				Name:        paramOpt.Name,
				In:          paramOpt.In,
				Description: paramOpt.Description,
				Required:    required,
				Example:     paramOpt.Example,
			}
			if paramOpt.Schema != nil {
				paramObj.Schema = generateSchema(paramOpt.Schema, schemaCtx)
			} else {
				paramObj.Schema = &SchemaObject{Type: "string"}
			}
			op.Parameters = append(op.Parameters, paramObj)
		}
	}

	// Ensure path parameters are included
	existingParams := make(map[string]bool)
	for _, p := range op.Parameters {
		existingParams[p.Name] = true
	}

	for _, seg := range route.segments {
		if seg.isParam && !existingParams[seg.paramName] {
			param := ParameterObject{
				Name:     seg.paramName,
				In:       "path",
				Required: true,
				Schema:   &SchemaObject{Type: "string"},
			}
			if route.options != nil {
				for _, pOpt := range route.options.Parameters {
					if pOpt.Name == seg.paramName && pOpt.In == "path" {
						param.Description = pOpt.Description
						param.Example = pOpt.Example
						if pOpt.Schema != nil {
							param.Schema = generateSchema(pOpt.Schema, schemaCtx)
						}
						break
					}
				}
			}
			op.Parameters = append(op.Parameters, param)
		}
	}

	// Default response if none specified
	if len(op.Responses) == 0 {
		op.Responses["200"] = &ResponseObject{Description: "OK"}
	}
	if len(op.Parameters) == 0 {
		op.Parameters = nil
	}
	return op
}

// generateSchema inspects an example instance via reflection and adds or reuses
// a SchemaObject in the context for components. Supports structs, arrays,
// maps, and basic Go types.
func generateSchema(instance any, ctx *schemaGenCtx) *SchemaObject {
	if instance == nil {
		return nil
	}
	val := reflect.ValueOf(instance)
	typ := val.Type()

	// Dereference pointer
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			val = reflect.New(typ.Elem())
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	// Reuse existing named schema for struct types
	if name, exists := ctx.generatedNames[typ]; exists && typ.Kind() == reflect.Struct {
		return &SchemaObject{Ref: "#/components/schemas/" + name}
	}

	schema := &SchemaObject{}

	switch typ.Kind() {
	case reflect.Struct:
		// Special-case time.Time
		if typ == reflect.TypeOf(time.Time{}) {
			schema.Type = "string"
			schema.Format = "date-time"
			return schema
		}

		// Assign a unique name and reserve slot
		schemaName := typ.Name()
		if schemaName == "" {
			schemaName = fmt.Sprintf("AnonymousStruct%d", len(ctx.componentsSchemas)+1)
		}
		originalName := schemaName
		counter := 1
		for {
			if _, used := ctx.componentsSchemas[schemaName]; !used {
				break
			}
			schemaName = fmt.Sprintf("%s%d", originalName, counter)
			counter++
		}
		ctx.generatedNames[typ] = schemaName
		ctx.componentsSchemas[schemaName] = nil // reserve

		// Build object properties
		schema.Type = "object"
		schema.Properties = make(map[string]*SchemaObject)
		var requiredFields []string

		for _, field := range reflect.VisibleFields(typ) {
			// skip unexported
			if field.PkgPath != "" {
				continue
			}

			// JSON tag handling
			tag := field.Tag.Get("json")
			parts := strings.Split(tag, ",")
			// skip explicit "-"
			if parts[0] == "-" {
				continue
			}

			fieldName := field.Name
			if parts[0] != "" {
				fieldName = parts[0]
			}

			// required if no "omitempty"
			if !slices.Contains(parts[1:], "omitempty") {
				requiredFields = append(requiredFields, fieldName)
			}

			desc := field.Tag.Get("description")
			example := field.Tag.Get("example")

			// pull the actual value via FieldByIndex (handles embedded)
			fv := val.FieldByIndex(field.Index)
			fieldSchema := generateSchema(fv.Interface(), ctx)

			if desc != "" {
				fieldSchema.Description = desc
			}
			if example != "" {
				fieldSchema.Example = example
			}
			schema.Properties[fieldName] = fieldSchema
		}

		if len(requiredFields) > 0 {
			schema.Required = requiredFields
		}
		if len(schema.Properties) == 0 {
			schema.Properties = nil
		}

		// finalize named schema
		ctx.componentsSchemas[schemaName] = schema
		return &SchemaObject{Ref: "#/components/schemas/" + schemaName}

	case reflect.Slice, reflect.Array:
		schema.Type = "array"
		elemType := typ.Elem()
		schema.Items = generateSchema(reflect.Zero(elemType).Interface(), ctx)

	case reflect.Map:
		schema.Type = "object"
		if typ.Key().Kind() != reflect.String {
			slog.Warn("OpenAPI schema generation: Map keys must be strings",
				"type", typ.String())
			schema.AdditionalProperties = &SchemaObject{Type: "object"}
			return schema
		}
		schema.AdditionalProperties = generateSchema(
			reflect.Zero(typ.Elem()).Interface(), ctx,
		)

	case reflect.String:
		schema.Type = "string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		schema.Type = "integer"
		schema.Format = "int32"

	case reflect.Int64:
		schema.Type = "integer"
		schema.Format = "int64"

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		schema.Type = "integer"
		schema.Format = "int64"

	case reflect.Float32:
		schema.Type = "number"
		schema.Format = "float"

	case reflect.Float64:
		schema.Type = "number"
		schema.Format = "double"

	case reflect.Bool:
		schema.Type = "boolean"

	case reflect.Interface:
		schema.Type = "object"
		schema.AdditionalProperties = &SchemaObject{}

	default:
		slog.Warn("OpenAPI schema generation: Unsupported type",
			"kind", typ.Kind().String())
		schema.Type = "object"
		schema.Description = fmt.Sprintf("Unsupported type: %s", typ.Kind())
	}

	return schema
}

// GenerateOpenAPISpec constructs an OpenAPI 3.0 specification from the given
// router and configuration, including paths, operations, and components.
func GenerateOpenAPISpec(router *Router, config OpenAPIConfig) *OpenAPI {
	spec := &OpenAPI{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       config.Title,
			Version:     config.Version,
			Description: config.Description,
		},
		Paths:      make(map[string]*PathItem),
		Components: &Components{Schemas: make(map[string]*SchemaObject)},
	}

	schemaCtx := newSchemaGenCtx()
	collectRoutes(router, spec, schemaCtx, "")

	if len(schemaCtx.componentsSchemas) > 0 {
		spec.Components.Schemas = schemaCtx.componentsSchemas
	} else {
		spec.Components = nil
	}

	return spec
}

// ServeOpenAPISpec makes your OpenAPI specification available at `path` (e.g. "/openapi.json")
func (r *Router) ServeOpenAPISpec(path string, config OpenAPIConfig) {
	spec := GenerateOpenAPISpec(r, config)

	specJSON, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal OpenAPI spec: %v", err))
	}

	r.Handle(http.MethodGet, path, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(specJSON)
	})

	slog.Info("OpenAPI specification served", "path", path)
}

//go:embed swagger-ui/*
var swaggerUIFS embed.FS

// ServeSwaggerUI shows swaggerUI at `prefix` (e.g. "/docs")
func (r *Router) ServeSwaggerUI(prefix string) {
	clean := strings.TrimSuffix(prefix, "/")

	sub, err := fs.Sub(swaggerUIFS, "swagger-ui")
	if err != nil {
		panic("failed to locate embedded swagger-ui assets: " + err.Error())
	}
	fsys := http.FS(sub)

	// Serve index.html at /docs
	r.Handle(http.MethodGet, clean, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f, err := fsys.Open("index.html")
		if err != nil {
			http.Error(w, "swagger-ui index.html not found", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		raw, _ := io.ReadAll(f)

		// Inject <base href="/docs/"> right after <head>
		inject := []byte(`<base href="` + clean + `/">`)
		adjusted := bytes.Replace(raw,
			[]byte("<head>"),
			append([]byte("<head>"), inject...),
			1,
		)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(adjusted)
	}))

	// Serve the static assets at /docs/{file}
	r.Handle(http.MethodGet, clean+"/{file}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		name := r.URLParam(req, "file")
		f, err := fsys.Open(name)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		defer f.Close()

		// set the real Content-Type
		if ct := mime.TypeByExtension(filepath.Ext(name)); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		io.Copy(w, f)
	}))

	slog.Info("Swagger UI served", "prefix", clean)
}
