package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/module/front-end/mapper"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (s *Service) DecryptStringToStructService(ctx context.Context, code string) (dto.PreRegisterConvern, error) {
	_, span := s.repo.Trace(ctx, "svc.DecryptStringToStructService", oteltrace.WithAttributes())
	defer span.End()
	//  Decrypt code
	decryptedData, _ := util.DecryptAes256Ecb(s.repo.AppCfg().Secret.EncryptKey, code)
	// unmarshall แปลง code เป็น struct
	var preRegData dto.PreRegisterConvern
	if err := json.Unmarshal([]byte(decryptedData), &preRegData); err != nil {
		return preRegData, err
	}
	return preRegData, nil
}

func (s *Service) CreateSessionService(ctx context.Context, code string) (dto.Session, error) {
	ctx, span := s.repo.Trace(ctx, "svc.CreateSessionService", oteltrace.WithAttributes())
	defer span.End()
	var res dto.Session
	claimPreRegData, err := s.DecryptStringToStructService(ctx, code)
	if err != nil {
		return res, err
	}
	// ตรวจสอบ code หมดอายุ
	if claimPreRegData.ExpireAt.Before(time.Now()) {
		return res, errors.New("link expire")
	}
	// ตรวจสอบข้อมูลกับฐานข้อมูล
	getPreRegister, err := s.repo.GetPreRegisterByPreRegisterUidRepo(ctx, nil, claimPreRegData.PreRegisterUUID, claimPreRegData.TransactionAuthUUID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Error("system error (Get PreRegisterData)", err)
			return res, errors.New("system error")
		}
		return res, errors.New("pre register data not found")
	}
	logrus.Infoln("getPreRegister ->", getPreRegister)
	// : สร้าง session
	// : map struct to Create Session
	mapSession := mapper.CreateSessionMapper(getPreRegister.PreRegister.UserID, getPreRegister.PreRegister.RegisterStatus, getPreRegister.PreRegister.Password, getPreRegister.PreRegister.Id.String(), getPreRegister.Id.String())
	// : set session
	if err := s.repo.SetAuthSession(mapSession.Uid.String(), getPreRegister.PreRegisterUid, mapSession); err != nil {
		logrus.Error("set session error ->", err)
		return res, errors.New("system error")
	}
	// : generate access token
	privateKey := s.repo.AppCfg().Secret.PrivateKey
	token, sid, _ := util.GenerateNewAccessTokenRepo(mapSession.Uid.String(), getPreRegister.PreRegisterUid, privateKey)
	// : เพิ่ม refresh token
	refreshToken, err := util.GenerateRefreshToken()
	if err != nil {
		logrus.Error("generate refresh token error ->", err)
		return res, errors.New("system error")
	}
	// : set refresh token
	if err := s.repo.SetRefreshToken(refreshToken, sid); err != nil {
		logrus.Error("set session error ->", err)
		return res, errors.New("system error")
	}
	resSessionMapper := mapper.ResponseSessionMapper(mapSession.ExpiresAt, sid, token, refreshToken)
	return resSessionMapper, nil
}
