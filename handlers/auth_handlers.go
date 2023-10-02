// handlers/auth_handlers.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"auth/auth"
)

// RegisterUserHandler handles user registration.
func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse JSON request body
		var requestBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Register user
		err := auth.RegisterUser(db, requestBody.Username, requestBody.Password)
		if err != nil {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("User registered successfully"))
	}
}

// LoginUserHandler handles user login.
func LoginUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse JSON request body
		var requestBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Login user and generate token
		token, err := auth.LoginUser(db, requestBody.Username, requestBody.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Send token in response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"token": token}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
	}
}
