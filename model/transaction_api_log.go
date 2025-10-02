package model

type TransactionAPILog struct {
	Model
	UserID        string `gorm:"column:user_id;size:255;"`
	API           string `gorm:"column:api;size:255;"`
	Status        string `gorm:"column:status;size:255;"` // pending, success, failed
	PaymentMethod string `gorm:"column:payment_method;size:255;"`
}
