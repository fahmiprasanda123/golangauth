// crud/item.go
package crud

import (
	"database/sql"
)

// Item represents an item in the system.
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateItem creates a new item.
func CreateItem(db *sql.DB, userID int, name, description string) error {
	_, err := db.Exec("INSERT INTO items (name, description, user_id) VALUES ($1, $2, $3)", name, description, userID)
	return err
}

// GetItems retrieves items for a user.
func GetItems(db *sql.DB, userID int) ([]Item, error) {
	rows, err := db.Query("SELECT id, name, description FROM items WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// UpdateItem updates an existing item.
func UpdateItem(db *sql.DB, itemID int, name, description string) error {
	_, err := db.Exec("UPDATE items SET name=$1, description=$2 WHERE id=$3", name, description, itemID)
	return err
}

// DeleteItem deletes an item.
func DeleteItem(db *sql.DB, itemID int) error {
	_, err := db.Exec("DELETE FROM items WHERE id=$1", itemID)
	return err
}
