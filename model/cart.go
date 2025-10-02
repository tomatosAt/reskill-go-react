package model

import "time"

type Cart struct {
	Model
	UserID  string    `gorm:"column:"user_id"`
	GameID  string    `gorm:"column:"game_id"`
	AddedAt time.Time `gorm:"column:"added_at"`
}
