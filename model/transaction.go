package model

type Transaction struct {
	Model
	UserID        uint    `json:"user_id"`
	TotalAmount   float64 `json:"total_amount"`
	Status        string  `json:"status"` // pending, success, failed
	PaymentMethod string  `json:"payment_method"`
}
