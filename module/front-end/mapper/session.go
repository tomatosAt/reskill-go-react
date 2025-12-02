package mapper

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
)

func CreateSessionMapper(userID, registerStatus, hashPassword string, preRegisterUuid, transactionAuthUuid string) dto.AuthSession {
	sessionId := uuid.NewV4()
	return dto.AuthSession{
		Uid:                sessionId,
		PreRegisterUid:     preRegisterUuid,
		UserID:             userID,
		TokenType:          "Bearer",
		RegisterStatus:     registerStatus,
		TransactionAuthUid: transactionAuthUuid,
		HashPassword:       hashPassword,
		ExpiresAt:          time.Now().Add(config.CachingShortDuration),
	}
}

func ResponseSessionMapper(expTime time.Time, sID, token, refreshToken string) dto.Session {
	expiresAt := time.Now().Add(config.CachingShortDuration)
	expiresAtRFC1123 := expiresAt.UTC().Format(time.RFC1123)

	return dto.Session{
		TokenType:    "Bearer",
		AccessToken:  token,
		RefreshToken: refreshToken,
		SessionID:    sID,
		ExpiresAt:    expiresAtRFC1123,
	}
}
