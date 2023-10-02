// main.go
package main

import (
	"database/sql"
	"log"
	"net/http"

	"auth/db"
	"auth/handlers"
)

var dbInstance *sql.DB

func main() {
	// Initialize database connection
	dbInstance, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbInstance.Close()

	// Routes
	http.HandleFunc("/register", handlers.RegisterUserHandler(dbInstance))
	http.HandleFunc("/login", handlers.LoginUserHandler(dbInstance))

	// Protected CRUD routes
	http.Handle("/create-item", handlers.CreateItemHandler(dbInstance))
	http.Handle("/get-items", handlers.GetItemsHandler(dbInstance))
	// http.Handle("/update-item", handlers.UpdateItemHandler(dbInstance))
	// http.Handle("/delete-item", handlers.DeleteItemHandler(dbInstance))

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
