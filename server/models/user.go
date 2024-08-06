package models

type User struct {
	UserID     int    `json:"user_id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	UserRole   string `json:"user_role"`
	JoinedDate string `json:"joined_date"`
}
