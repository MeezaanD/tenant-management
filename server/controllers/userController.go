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

	row := database.DB.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", user.Email, user.Password)

	if err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.UserRole, &user.JoinedDate); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Println("User not found")
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":     user.UserID,
		"firstname":  user.FirstName,
		"lastname":   user.LastName,
		"email":      user.Email,
		"userRole":   user.UserRole,
		"joinedDate": user.JoinedDate,
		"exp":        time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

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
	params := mux.Vars(r)
	userID := params["userID"]

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	_, err = database.DB.Exec("UPDATE users SET firstname=?, lastname=?, email=?, userRole=? WHERE user_id=?", 
		user.FirstName, user.LastName, user.Email, user.UserRole, userID)

	if err != nil {
		log.Printf("Error updating user with ID %s: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	_, err := database.DB.Exec("DELETE FROM users WHERE user_id = ?", userID)
	if err != nil {
		if err.Error() == "Error 1451: Cannot delete or update a parent row: a foreign key constraint fails" {
			log.Printf("Cannot delete user %s: associated payments exist", userID) 
			w.WriteHeader(http.StatusConflict)                                     
			json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete user with associated payments"})
			return
		}

		// Log the exact error for debugging
		log.Printf("Error deleting user %s: %v", userID, err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete user"})
		return
	}

	// Successfully deleted the user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
