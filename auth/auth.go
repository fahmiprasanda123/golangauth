// auth/auth.go
package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"database/sql"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "your_secret_key"

// RegisterUser registers a new user.
func RegisterUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hashedPassword))

	return err
}

func LoginUser(db *sql.DB, username, password string) (string, error) {
	var userID int
	var storedPassword string

	err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &storedPassword)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := generateToken(userID, username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Claims represents the token claims.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// ValidateToken validates the JWT token.
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	// Check for errors
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	return claims, nil
}

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// The token should be in the format "Bearer <token>"
		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the token
		claims, err := ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user claims to the request context (even if not used directly)
		_ = claims

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func generateToken(userID int, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
