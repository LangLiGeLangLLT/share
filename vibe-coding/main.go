package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Item represents a simple data model
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage
var (
	items   = make(map[int]Item)
	nextID  = 1
	itemsMu sync.RWMutex
)

// Create - POST /items
func createItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	item.ID = nextID
	nextID++
	items[item.ID] = item
	itemsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Read All - GET /items
func getAllItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemsMu.RLock()
	itemList := make([]Item, 0, len(items))
	for _, item := range items {
		itemList = append(itemList, item)
	}
	itemsMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemList)
}

// Read One - GET /items/{id}
func getItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	itemsMu.RLock()
	item, exists := items[id]
	itemsMu.RUnlock()

	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Update - PUT /items/{id}
func updateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	_, exists := items[id]
	if !exists {
		itemsMu.Unlock()
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	updatedItem.ID = id
	items[id] = updatedItem
	itemsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

// Delete - DELETE /items/{id}
func deleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	_, exists := items[id]
	if !exists {
		itemsMu.Unlock()
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	delete(items, id)
	itemsMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// Root handler
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	// Root endpoint
	http.HandleFunc("/", helloWorldHandler)

	// CRUD endpoints
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllItems(w, r)
		case http.MethodPost:
			createItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItem(w, r)
		case http.MethodPut:
			updateItem(w, r)
		case http.MethodDelete:
			deleteItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := ":8000"
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  GET    /           - Hello World")
	fmt.Println("  POST   /items      - Create new item")
	fmt.Println("  GET    /items      - Get all items")
	fmt.Println("  GET    /items/?id=X - Get item by ID")
	fmt.Println("  PUT    /items/?id=X - Update item by ID")
	fmt.Println("  DELETE /items/?id=X - Delete item by ID")
	fmt.Println()

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}