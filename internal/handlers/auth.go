package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your-secret-key") // Replace with a strong key

// User structure for login/register
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// In-memory user storage for testing (replace with database logic)
var users = map[string]string{} // username: password

// Generate JWT token
func GenerateToken(username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Register a new user
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if _, exists := users[user.Username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}
	users[user.Username] = user.Password
	w.Write([]byte("User registered successfully"))
}

// Login and generate token
func Login(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	storedPassword, exists := users[creds.Username]
	if !exists || storedPassword != creds.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Logout (invalidate token)
func Logout(w http.ResponseWriter, r *http.Request) {
	// This implementation only works if you have token blacklisting (stateful JWT)
	w.Write([]byte("Logout successful"))
}
