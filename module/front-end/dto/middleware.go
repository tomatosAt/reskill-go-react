package dto

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AuthSession struct {
	Uid                uuid.UUID `json:"uid"`
	PreRegisterUid     string    `json:"pre_register_uid"`
	UserID             string    `json:"user_id"`
	TokenType          string    `json:"token_type"`
	RegisterStatus     string    `json:"register_status"`
	HashPassword       string    `json:"hash_password"`
	TransactionAuthUid string    `json:"transaction_auth_uid"`
	ExpiresAt          time.Time `json:"expires_at"`
}

type Session struct {
	SessionID    string `json:"sid"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}
