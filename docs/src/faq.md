# FAQ

## How do I get JSON data sent in the _body_ of a request (e.g., POST, PUT)?

```go
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MyData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func handleJsonRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var data MyData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received data: Name=%s, Value=%d\n", data.Name, data.Value)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "received_name": data.Name}
	json.NewEncoder(w).Encode(response)
}
```

## How do I get data from a GET request? Can I send JSON in a GET request?

```go
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ItemData struct {
	ItemName string `json:"itemName"`
	ItemID   int    `json:"itemId"`
}

func handleGetWithQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	queryParams := r.URL.Query() // This returns a url.Values map (map[string][]string)

	name := queryParams.Get("name") // .Get() returns the first value for the key
	idStr := queryParams.Get("id")

	if name == "" || idStr == "" {
		http.Error(w, "Missing 'name' or 'id' query parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter, must be an integer", http.StatusBadRequest)
		return
	}

	item := ItemData{
		ItemName: name,
		ItemID:   id,
	}

	log.Printf("Processed GET request: ItemName=%s, ItemID=%d\n", item.ItemName, item.ItemID)

	// Send a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
```
