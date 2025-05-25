package nova

import "reflect"

// OpenAPIConfig holds metadata for generating an OpenAPI specification.
type OpenAPIConfig struct {
	Title       string // Title of the API
	Version     string // Version of the API
	Description string // Description of the API
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
