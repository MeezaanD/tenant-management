package controllers

import (
	"encoding/json"
	"net/http"
	"tenant-management/database"
	"tenant-management/models"

	"github.com/gorilla/mux"
)

func GetAllPayments(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT p.payment_id, p.user_id, u.firstname, p.agreed_amount, p.paid_amount, 
		       p.remaining_amount, p.proof_of_payment, p.paid_ontime
		FROM payments p
		LEFT JOIN users u ON p.user_id = u.user_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.PaymentID, &payment.UserID, &payment.FirstName, &payment.AgreedAmount, &payment.PaidAmount, &payment.RemainingAmount, &payment.ProofOfPayment, &payment.PaidOnTime); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments = append(payments, payment)
	}

	response := map[string]interface{}{
		"message":  "All payments retrieved successfully",
		"payments": payments,
	}
	json.NewEncoder(w).Encode(response)
}

func GetPaymentsByUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	rows, err := database.DB.Query(`
		SELECT p.payment_id, p.user_id, u.firstname, p.agreed_amount, p.paid_amount, 
		       p.remaining_amount, p.proof_of_payment, p.paid_ontime
		FROM payments p
		LEFT JOIN users u ON p.user_id = u.user_id
		WHERE p.user_id = ?
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.PaymentID, &payment.UserID, &payment.FirstName, &payment.AgreedAmount, &payment.PaidAmount, &payment.RemainingAmount, &payment.ProofOfPayment, &payment.PaidOnTime); err != nil {
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

	response := map[string]interface{}{
		"message": "Payment created successfully",
		"payment": payment,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	paymentID := params["paymentID"]

	var payment models.Payment

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE payments SET user_id = ?, agreed_amount = ?, paid_amount = ?, remaining_amount = ?, proof_of_payment = ?, paid_ontime = ? WHERE payment_id = ?",
		payment.UserID, payment.AgreedAmount, payment.PaidAmount, payment.RemainingAmount, payment.ProofOfPayment, payment.PaidOnTime, paymentID)

	if err != nil {
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Payment updated successfully",
		"payment": payment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	paymentID := params["paymentID"]

	_, err := database.DB.Exec("DELETE FROM payments WHERE payment_id = ?", paymentID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":    "Payment deleted successfully",
		"payment_id": paymentID,
	}
	json.NewEncoder(w).Encode(response)
}
