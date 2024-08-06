package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"tenant-management/database"
	"tenant-management/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.UserRole, &user.JoinedDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response := map[string]interface{}{
		"message": "Users retrieved successfully",
		"user":    users,
	}
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	var user models.User
	row := database.DB.QueryRow("SELECT * FROM users WHERE user_id = ?", userID)

	if err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.UserRole, &user.JoinedDate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User retrieved successfully",
		"user":    user,
	}
	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	// Insert the user into the database
	_, err := database.DB.Exec("INSERT INTO users (firstname, lastname, email, password, userRole, joinedDate) VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.UserRole, user.JoinedDate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message and user details
	response := map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	// Query the user by email and password
	row := database.DB.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", user.Email, user.Password)

	if err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.UserRole, &user.JoinedDate); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Println("User not found")
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID,
		"email":  user.Email,
		"role":   user.UserRole,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Respond with the token and success message
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
	}
	json.NewEncoder(w).Encode(response)
}

func LogOutUser(w http.ResponseWriter, r *http.Request) {
	// For JWT-based authentication, logout can simply be a client-side action where the token is discarded.
	// Respond with a success message
	response := map[string]interface{}{
		"message": "Logout successful",
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	_, err := database.DB.Exec("UPDATE users SET firstname = ?, lastname = ?, email = ?, password = ?, userRole = ?, joinedDate = ? WHERE user_id = ?",
		user.FirstName, user.LastName, user.Email, user.Password, user.UserRole, user.JoinedDate, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	_, err := database.DB.Exec("DELETE FROM users WHERE user_id = ?", userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"message": "User deleted successfully",
		"userID":  userID,
	}
	json.NewEncoder(w).Encode(response)
}
