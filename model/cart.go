package model

import "time"

type Cart struct {
	Model
	UserID  uint      `json:"user_id"`
	GameID  uint      `json:"game_id"`
	AddedAt time.Time `json:"added_at"`
}
