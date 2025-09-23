package model

type TransactionAPILog struct {
	Model
	UserID        uint   `json:"user_id"`
	API           string `json:"api"`
	Status        string `json:"status"` // pending, success, failed
	PaymentMethod string `json:"payment_method"`
}
