{{ title: Nova - FAQ }}

{{ include-block: doc.html markdown="true" }}

# FAQ

### How do I get JSON data from the body of a request (e.g., POST, PUT)?

With `nova`, you can use the `rc.Bind()` method on the `ResponseContext` to automatically decode JSON from the request body into your struct. This avoids manual decoding and error handling boilerplate.

Here is a complete example using `PostFunc` and `rc.Bind()`:

```go
package main

import (
	"log"
	"net/http"

    "github.com/xlc-dev/nova/nova"
)

// MyData defines the structure of our expected JSON payload.
type MyData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// handleJsonRequest uses nova's ResponseContext to simplify data binding and response.
func handleJsonRequest(rc *nova.ResponseContext) error {
	var data MyData

	// Bind the incoming JSON request body to the 'data' struct.
	if err := rc.BindJSON(&data); err != nil {
		// If binding fails (e.g., malformed JSON), return a 400 Bad Request.
        log.Printf("Error binding JSON: %v", err)
		return rc.JSONError(http.StatusBadRequest, "Invalid JSON payload")
	}

	log.Printf("Received data: Name=%s, Value=%d\n", data.Name, data.Value)

	// Send a success response using the JSON helper.
    response := map[string]any{
		"status":        "success",
		"received_name": data.Name,
	}
	return rc.JSON(http.StatusOK, response)
}

func main() {
	router := nova.NewRouter()
	router.PostFunc("/data", handleJsonRequest)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
```

### How do I get data from a GET request?

GET requests typically pass data in two ways: **URL query parameters** (e.g., `?id=123`) or **path parameters** (e.g., `/items/123`). `nova` makes handling both easy.

#### 1. Reading Query Parameters

For query parameters, you can access the standard `http.Request` object via `rc.Request()` and use its `URL.Query()` method.

```go
// handleGetWithQuery processes data from URL query parameters.
// Example request: GET /items?name=widget&id=123
func handleGetWithQuery(rc *nova.ResponseContext) error {
	// Access the underlying request to get query parameters.
	queryParams := rc.Request().URL.Query()

	name := queryParams.Get("name")
	idStr := queryParams.Get("id")

	if name == "" || idStr == "" {
		return rc.JSONError(http.StatusBadRequest, "Missing 'name' or 'id' query parameter")
	}

	log.Printf("Processed query: name=%s, id=%s\n", name, idStr)

	response := map[string]string{
		"item_name": name,
		"item_id":   idStr,
	}
	return rc.JSON(http.StatusOK, response)
}

// In main():
// router.GetFunc("/items", handleGetWithQuery)
```

#### 2. Reading Path Parameters

For path parameters defined in your route (e.g., `/{id}`), `nova` provides the much cleaner `rc.URLParam()` helper.

```go
// handleGetWithPathVar processes data from a URL path parameter.
// Example request: GET /users/42
func handleGetWithPathVar(rc *nova.ResponseContext) error {
	// Use the URLParam helper to get the 'id' from the path.
	userID := rc.URLParam("id")

	log.Printf("Processed request for user ID: %s\n", userID)

	response := map[string]string{
		"status":  "found",
		"user_id": userID,
	}
	return rc.JSON(http.StatusOK, response)
}

// In main():
// router.GetFunc("/users/{id}", handleGetWithPathVar)
```

### How do I validate incoming request data?

Nova has powerful, built-in validation. Instead of `rc.Bind()`, use **`rc.BindValidated()`**. It automatically binds the data and then runs validations based on the struct tags you've defined (`required`, `minlength`, `format`, etc.).

If validation fails, it returns a detailed error, which you can send back to the client.

```go
package main

import (
	"log"
	"net/http"

	"github.com/xlc-dev/nova/nova"
)

// UserSignup defines a struct with validation tags.
type UserSignup struct {
	Username string `json:"username" minlength:"3" maxlength:"20" format:"alphanumeric"`
	Email    string `json:"email" format:"email"`
	// omitempty is not used, so this field is required by default.
	Password string `json:"password" format:"password"`
}

// handleUserSignup binds and validates the request body in one step.
func handleUserSignup(rc *nova.ResponseContext) error {
	var user UserSignup

	// Bind AND validate the incoming JSON.
	if err := rc.BindValidated(&user); err != nil {
		// e.g., "validation failed: Field 'username' must contain only alphanumeric characters"
		log.Printf("Validation failed: %v", err)
		return rc.JSONError(http.StatusBadRequest, err.Error())
	}

	log.Printf("Successfully validated and created user: %s", user.Username)

	response := map[string]string{
		"status":   "user_created",
		"username": user.Username,
	}
	return rc.JSON(http.StatusCreated, response)
}

func main() {
	router := nova.NewRouter()
	router.PostFunc("/signup", handleUserSignup)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
```

{{ endinclude }}
