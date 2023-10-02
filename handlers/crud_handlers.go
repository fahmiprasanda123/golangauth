// handlers/crud_handlers.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"auth/auth"
	"auth/crud"
)

// CreateItemHandler handles the creation of a new item.
func CreateItemHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse user claims from context
		claims, ok := r.Context().Value("user").(*auth.Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse JSON request body
		var newItem crud.Item
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create item
		err := crud.CreateItem(db, claims.UserID, newItem.Name, newItem.Description)
		if err != nil {
			http.Error(w, "Error creating item", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Item created successfully"))
	})
}

// GetItemsHandler retrieves items for the authenticated user.
func GetItemsHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse user claims from context
		claims, ok := r.Context().Value("user").(*auth.Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get items for the authenticated user
		items, err := crud.GetItems(db, claims.UserID)
		if err != nil {
			http.Error(w, "Error getting items", http.StatusInternalServerError)
			return
		}

		// Convert items to JSON and send response
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(items)
		w.Write(jsonResponse)
	})
}

// UpdateItemHandler handles the update of an existing item.
// func UpdateItemHandler(db *sql.DB) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Parse user claims from context
// 		claims, ok := r.Context().Value("user").(*auth.Claims)
// 		if !ok {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Parse JSON request body
// 		var updatedItem crud.Item
// 		if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}

// 		// Update item
// 		err := crud.UpdateItem(db, updatedItem.ID, updatedItem.Name, updatedItem.Description)

// 		if err != nil {
// 			http.Error(w, "Error updating item", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Write([]byte("Item updated successfully"))
// 	})
// }

// // DeleteItemHandler handles the deletion of an item.
// func DeleteItemHandler(db *sql.DB) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Parse user claims from context
// 		claims, ok := r.Context().Value("user").(*auth.Claims)
// 		if !ok {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Parse JSON request body
// 		var deleteRequest struct {
// 			ItemID int `json:"item_id"`
// 		}
// 		if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}

// 		// Delete item
// 		err := crud.DeleteItem(db, deleteRequest.ItemID)
// 		if err != nil {
// 			http.Error(w, "Error deleting item", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Write([]byte("Item deleted successfully"))
// 	})
// }
