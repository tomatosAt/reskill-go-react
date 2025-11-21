package services

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/module/front-end/mapper"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (s Service) GetProfilesService(ctx context.Context) (dto.ResponseUserProfiles, int, error) {
	ctx, span := s.repo.Trace(ctx, "svc.GetProfilesService", oteltrace.WithAttributes())
	defer span.End()
	cliams, _ := s.repo.GetAuthCtxRepo(ctx)
	var resp dto.ResponseUserProfiles
	getPreRegister, err := s.repo.GetPreRegisterByPreRegisterUidRepo(ctx, nil, cliams.PreRegisterUid, cliams.TransactionAuthUid)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Error("system error (Get PreRegisterData)", err)
			return resp, 500, errors.New("system error")
		}
		return resp, 500, errors.New("pre register data not found")
	}
	decryptData, _ := mapper.DecodePreRegisterConvernMapper(s.repo.AppCfg().Secret.EncryptKey, getPreRegister.PreRegister)
	resp = mapper.ResponseProfileMapper(getPreRegister.PreRegister.Id.String(), getPreRegister.PreRegister.UserID, getPreRegister.PreRegister.Tel, getPreRegister.PreRegister.Email, getPreRegister.PreRegister.NickName, decryptData)
	return resp, 200, nil
}
