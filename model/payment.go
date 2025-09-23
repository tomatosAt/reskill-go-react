package model

import "time"

type Purchase struct {
	Model
	UserID        uint      `json:"user_id"`
	GameID        uint      `json:"game_id"`
	TransactionID *uint     `json:"transaction_id"` // optional
	PurchasedAt   time.Time `json:"purchased_at"`
}
