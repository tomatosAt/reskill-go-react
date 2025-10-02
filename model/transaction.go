package model

type Transaction struct {
	Model
	UserID        string  `gorm:"column:user_id;size:255;"`
	TotalAmount   float64 `gorm:"column:total_amount;size:255;"`
	Status        string  `gorm:"column:status;size:255;"` // pending, success, failed
	PaymentMethod string  `gorm:"column:payment_method;size:255;"`
}
