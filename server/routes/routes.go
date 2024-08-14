package routes

import (
	"tenant-management/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	// User routes
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/user/{userID}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/user/{userID}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{userID}", controllers.DeleteUser).Methods("DELETE")

	// Payment routes
	router.HandleFunc("/payments", controllers.GetAllPayments).Methods("GET")
	router.HandleFunc("/user/{userID}/payments", controllers.GetPaymentsByUsers).Methods("GET")
	router.HandleFunc("/user/{userID}/payment/{paymentID}", controllers.GetPayment).Methods("GET")
	router.HandleFunc("/user/{userID}/payments", controllers.CreatePayment).Methods("POST")
	router.HandleFunc("/payment/{paymentID}", controllers.UpdatePayment).Methods("PUT")
	router.HandleFunc("/payment/{paymentID}", controllers.DeletePayment).Methods("DELETE")

}
