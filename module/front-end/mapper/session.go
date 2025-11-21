package mapper

import (
	uuid "github.com/satori/go.uuid"
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
	}
}

func ResponseSessionMapper(sID, token string) dto.Session {
	return dto.Session{
		TokenType: "Bearer",
		Token:     token,
		SessionID: sID,
	}
}
