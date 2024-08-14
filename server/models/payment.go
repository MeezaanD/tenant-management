package models

type Payment struct {
	PaymentID       int    `json:"payment_id"`
	UserID          int    `json:"user_id"`
	FirstName       string `json:"firstname"`
	AgreedAmount    int    `json:"agreed_amount"`
	PaidAmount      int    `json:"paid_amount"`
	RemainingAmount int    `json:"remaining_amount"`
	ProofOfPayment  string `json:"proof_of_payment"`
	PaidOnTime      bool   `json:"paid_ontime"`
}
