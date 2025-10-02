package model

import "time"

type Purchase struct {
	Model
	UserID        string    `gorm:"column:user_id;size:255;"`
	GameID        string    `gorm:"column:game_id;size:255;"`
	TransactionID *string   `gorm:"column:transaction_id;size:255;"`
	PurchasedAt   time.Time `gorm:"column:purchased_at;size:255;"`
}
