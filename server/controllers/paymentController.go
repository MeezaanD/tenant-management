package controllers

import (
	"encoding/json"
	"net/http"
	"tenant-management/database"
	"tenant-management/models"

	"github.com/gorilla/mux"
)

func GetAllPayments(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM payments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.PaymentID, &payment.UserID, &payment.AgreedAmount, &payment.PaidAmount, &payment.RemainingAmount, &payment.ProofOfPayment, &payment.PaidOnTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments = append(payments, payment)
	}

	// Respond with all payments
	response := map[string]interface{}{
		"message":  "All payments retrieved successfully",
		"payments": payments,
	}
	json.NewEncoder(w).Encode(response)
}

func GetPaymentsByUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM payments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.PaymentID, &payment.UserID, &payment.AgreedAmount, &payment.PaidAmount, &payment.RemainingAmount, &payment.ProofOfPayment, &payment.PaidOnTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments = append(payments, payment)
	}

	json.NewEncoder(w).Encode(payments)
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	paymentID := params["paymentID"]

	var payment models.Payment
	row := database.DB.QueryRow("SELECT * FROM payments WHERE payment_id = ?", paymentID)

	if err := row.Scan(&payment.PaymentID, &payment.UserID, &payment.AgreedAmount, &payment.PaidAmount, &payment.RemainingAmount, &payment.ProofOfPayment, &payment.PaidOnTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	var payment models.Payment
	json.NewDecoder(r.Body).Decode(&payment)

	// Check if the user_id exists in the users table
	var userExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)", userID).Scan(&userExists)

	if err != nil || !userExists {
		http.Error(w, "User ID does not exist", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO payments (user_id, agreed_amount, paid_amount, remaining_amount, proof_of_payment, paid_ontime) VALUES (?, ?, ?, ?, ?, ?)",
		userID, payment.AgreedAmount, payment.PaidAmount, payment.RemainingAmount, payment.ProofOfPayment, payment.PaidOnTime)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"message": "Payment created successfully",
		"payment": payment,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	json.NewDecoder(r.Body).Decode(&payment)

	_, err := database.DB.Exec("UPDATE payments SET user_id = ?, agreed_amount = ?, paid_amount = ?, remaining_amount = ?, proof_of_payment = ?, paid_ontime = ? WHERE payment_id = ?",
		payment.UserID, payment.AgreedAmount, payment.PaidAmount, payment.RemainingAmount, payment.ProofOfPayment, payment.PaidOnTime, payment.PaymentID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"message": "Payment updated successfully",
		"payment": payment,
	}
	json.NewEncoder(w).Encode(response)
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	json.NewDecoder(r.Body).Decode(&payment)

	_, err := database.DB.Exec("DELETE FROM payments WHERE payment_id = ?", payment.PaymentID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"message": "Payment deleted successfully",
		"payment": payment,
	}
	json.NewEncoder(w).Encode(response)
}
