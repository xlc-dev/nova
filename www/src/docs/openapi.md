{{ title: Nova - OpenAPI }}

{{ include-block: doc.html markdown="true" }}

# OpenAPI

Nova's OpenAPI support automatically generates and serves an OpenAPI 3.0 specification for your routes, plus an embedded Swagger UI. It provides:

- **Automatic Spec Generation:** Reflects Go structs, route definitions, parameters, request bodies, and responses into a valid OpenAPI 3.0 document.
- **Component Schemas:** Deduplicates and reuses schema definitions for structs, arrays, maps, and basic types.
- **Path & Operation Mapping:** Converts your router’s methods (`Get`, `Post`, etc.) into OpenAPI `paths` with `OperationObject`.
- **Parameter Inference:** Generates path, query, header, and cookie parameters, including required flags and examples.
- **Request & Response Bodies:** Maps `RequestBody` and `Responses` options to JSON media types and schemas.
- **Servers Configuration:** Embeds one or more server definitions (URLs) into the spec.
- **Spec Endpoint:** Serves the JSON spec at a configurable path (e.g. `/openapi.json`).
- **Swagger UI:** Embeds Swagger UI assets under a prefix (e.g. `/docs`), complete with static asset handling.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Core Concepts](#core-concepts)
   - [OpenAPIConfig](#openapiconfig)
   - [RouteOptions and ResponseOption](#routeoptions--responseoption)
   - [Schema Generation](#schema-generation)

3. [Registering and Serving the Spec](#registering--serving-the-spec)
4. [Serving Swagger UI](#serving-swagger-ui)
5. [Full Example](#full-example)

## Getting Started

To enable OpenAPI support in your Nova-powered app:

1. **Import Nova’s router and OpenAPI helpers:**

```go
import "github.com/xlc-dev/nova/nova"
```

2. **Define your routes** with `RouteOptions` (including tags, summary, parameters, request/response schemas).
3. **Call**:

```go
router.ServeOpenAPISpec("/openapi.json", nova.OpenAPIConfig{
 Title:       "My Nova API",
 Version:     "1.0.0",
 Description: "This is a sample API built with Nova.",
 Servers:     []nova.Server{{URL: fmt.Sprintf("http://%s:%d", host, port)}},
})
router.ServeSwaggerUI("/docs")
```

4. **Start** your server as usual: `nova.Serve(ctx, router)`.

## Core Concepts

### OpenAPIConfig

Holds top‐level metadata for the spec:

```go
type OpenAPIConfig struct {
  Title       string    // API title (required)
  Version     string    // Spec version (required)
  Description string    // Optional description
  Servers     []Server  // List of servers (URL + optional description)
}
```

### RouteOptions and ResponseOption

Attach OpenAPI metadata when registering routes:

```go
type RouteOptions struct {
  Tags        []string            // Operation tags (grouping)
  Summary     string              // Short summary
  Description string              // Detailed description
  OperationID string              // Unique operation ID
  Deprecated  bool                // Mark as deprecated
  RequestBody interface{}         // Request schema
  Responses   map[int]ResponseOption
  Parameters  []ParameterOption
}

type ResponseOption struct {
  Description string      // Response description
  Body        interface{} // Schema
}
```

- **Responses**: map HTTP status codes to `ResponseOption`.
- **Parameters**: define additional `in:path|query|header|cookie` parameters.

### Schema Generation

Nova inspects Go types via reflection to build JSON Schemas:

- **Structs:**
  - Generates `components.schemas` entries.
  - Honors `json:"..."`, `description:"..."`, `example:"..."` tags.
  - Marks fields non‐nullable or required unless `omitempty`.

- **Primitives:** maps Go kinds to `type`/`format` (e.g., `time.Time` → `string` + `date‐time`).
- **Arrays & Slices:** `type: array` + `items`.
- **Maps:** `type: object` + `additionalProperties`.
- **References:** reuses named schemas when the same struct type appears multiple times.

## Registering and Serving the Spec

```go
// After defining routes on your `router`:
router.ServeOpenAPISpec("/openapi.json", nova.OpenAPIConfig{
  Title:       "My API",
  Version:     "2.0.0",
  Description: "Auto‐generated OpenAPI spec",
  Servers: []nova.Server{
    {URL: "https://api.example.com"},
  },
})
```

- **Endpoint:** performs a `GET /openapi.json`, returning JSON with `Content-Type: application/json`.
- **Internal:** calls `GenerateOpenAPISpec(router, config)` under the hood.

## Serving Swagger UI

Nova embeds the official Swagger UI and serves it statically:

```go
router.ServeSwaggerUI("/docs")
```

- **Access:**
  - `GET /docs` → redirects to `/docs/` and serves `index.html`.
  - `GET /docs/{file}` → serves CSS, JS, and favicon from embedded assets.

- **Customization:** To point the UI at your spec URL, pass query parameters to the `/docs/index.html` link (e.g. `?url=/openapi.json`).

## Full Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/xlc-dev/nova/nova"
)

func main() {
	cli, err := nova.NewCLI(&nova.CLI{
		Name:        "item-api",
		Version:     "1.0.0",
		Description: "A simple item API with OpenAPI & Swagger UI",
		GlobalFlags: []nova.Flag{},
		Action: func(ctx *nova.Context) error {
			router := nova.NewRouter()

			type Item struct {
				ID        int       `json:"id" description:"Item ID"`
				Name      string    `json:"name" example:"Widget"`
				CreatedAt time.Time `json:"createdAt" description:"Creation timestamp"`
			}

			router.Get("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
				idStr := router.URLParam(r, "id")
				id, _ := strconv.Atoi(idStr)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w,
					`{"id":%d,"name":"Item%d","createdAt":"%q"}`,
					id, id, time.Now().Format(time.RFC3339),
				)
			}, &nova.RouteOptions{
				Tags:    []string{"Items"},
				Summary: "Retrieve an item by ID",
				Parameters: []nova.ParameterOption{{
					Name:        "id",
					In:          "path",
					Description: "ID of the item to retrieve",
					Schema:      int(0),
				}},
				Responses: map[int]nova.ResponseOption{
					200: {Description: "Success", Body: &Item{}},
					404: {Description: "Not found"},
				},
			})

			router.ServeOpenAPISpec("/openapi.json", nova.OpenAPIConfig{
				Title:       "Item API",
				Version:     "1.0.0",
				Description: "API for managing items",
				Servers: []nova.Server{
					{URL: fmt.Sprintf("http://%s:%d", "localhost", 8080)},
				},
			})

			router.ServeSwaggerUI("/docs")

			return nova.Serve(ctx, router)
		},
	})
	if err != nil {
		log.Fatalf("failed to initialize CLI: %v", err)
	}
	if err := cli.Run(os.Args); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}
```

{{ endinclude }}
