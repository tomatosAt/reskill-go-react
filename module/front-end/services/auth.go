package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/module/front-end/mapper"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (s *Service) LoginUserPassService(ctx context.Context, userPass dto.UserPasswordPayload) (*dto.ResponsePreRegister, int, error) {
	ctx, span := s.repo.Trace(ctx, "svc.LoginUserPassService", oteltrace.WithAttributes())
	defer span.End()
	var result dto.ResponsePreRegister
	userPass.Password = util.HashSHA256(userPass.Password)
	// Check UserPass พวก type ตรงมั้ย format ถูกมั้ย
	if err := s.CheckFormatUsernameSVC(ctx, &userPass); err != nil {
		return &result, http.StatusBadRequest, err
	}
	// Logic เช็ค สถานะ user ว่าถูกต้องมั้ย => มีสถานะอะไรบ้าง (login ได้ปกติ, ลงทะเบียน,username ผิด, user ถูกระงับ)
	status := s.CheckFStatusUsernameSVC(ctx, &userPass)
	if status != "login" {
		if status == "register" {
			return &result, http.StatusBadRequest, errors.New("user status is register") // ให้หน้าบ้านไปลงทะเบียน
		} else {
			return &result, http.StatusBadRequest, errors.New("user status is suspend")
		}
	}
	tx := s.repo.DB().Ctx().Begin()
	preRegisterRepo, err := s.repo.GetPreRegisterByUserPassEmailRepo(ctx, tx, userPass.Username, userPass.Password, userPass.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return &result, http.StatusInternalServerError, err
		}
	}
	// Get transaction auth เพื่อ อัปเดตสถานะการ login
	transactionAuthMapper := mapper.InsertTransactionAuthMapper(status, "complete", status, preRegisterRepo.Id.String())
	transactionAuthRepo, err := s.repo.InsertTransactionAuthRepo(ctx, tx, transactionAuthMapper)
	if err != nil {
		logrus.Errorln("[LoginUserPassService] InsertTransactionAuthRepo err->", err)
		tx.Rollback()
		util.RecordSpanError(span, err, "repo.InsertTransactionAuthRepo")
		return &result, http.StatusInternalServerError, errors.New("server error")
	}
	// Table user
	userAuth, err := s.repo.GetUsersAuthRepo(ctx, tx, userPass.Username, userPass.Password)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return &result, http.StatusInternalServerError, err
		} else {
			//  สร้าง Table user
			userAuthMapper := mapper.InsertUsersAuthMapper(userPass.Username, userPass.Password, "user")
			userAuth, err = s.repo.InsertUsersAuthRepo(ctx, tx, userAuthMapper)
			if err != nil {
				logrus.Errorln("[LoginUserPassService] InsertUsersAuthRepo err->", err)
				tx.Rollback()
				util.RecordSpanError(span, err, "repo.InsertUsersAuthRepo")
				return &result, http.StatusInternalServerError, errors.New("server error")
			}
		}
	}
	// Get PreRegister auth เพื่อ อัปเดตสถานะการ login
	if err := s.repo.UpdatePreRegisterByUserPassEmailRepo(ctx, tx, preRegisterRepo.Id.String(), map[string]interface{}{"register_status": "login", "user_id": userAuth.Id}); err != nil {
		logrus.Errorln("[LoginUserPassService] UpdatePreRegisterByUserPassEmailRepo err->", err)
		tx.Rollback()
		util.RecordSpanError(span, err, "repo.UpdatePreRegisterByUserPassEmailRepo")
		return &result, http.StatusInternalServerError, errors.New("server error")
	}
	tx.Commit()
	convernPreRegister := mapper.PreRegisterToConvernMapper(preRegisterRepo.Id, transactionAuthRepo.Id)
	convertByte, err := mapper.ConvernToByteMapper(convernPreRegister)
	if err != nil {
		logrus.Errorln("[LoginUserPassService] ConvernToByteMapper err->", err)
		util.RecordSpanError(span, err, "repo.mapper.ConvernToByteMapper")
	}
	encryptedConvertByte, _ := util.EncryptAes256Ecb(s.repo.AppCfg().Secret.EncryptKey, string(convertByte))
	result = mapper.ResponsePreRegisterMapper(string(encryptedConvertByte))
	return &result, http.StatusOK, nil
}

func (s *Service) CheckFormatUsernameSVC(ctx context.Context, data *dto.UserPasswordPayload) error {
	_, span := s.repo.Trace(ctx, "svc.CheckFormatUsernameSVC", oteltrace.WithAttributes())
	defer span.End()
	// Input validation
	if _, err := util.ValidatorStruct(data); err != nil {
		return err
	}
	return nil
}

func (s *Service) CheckFStatusUsernameSVC(ctx context.Context, data *dto.UserPasswordPayload) string {
	ctx, span := s.repo.Trace(ctx, "svc.CheckFStatusUsernameSVC", oteltrace.WithAttributes())
	defer span.End()
	// TODO Logic เช็ค สถานะ user ว่าถูกต้องมั้ย => มีสถานะอะไรบ้าง (login ได้ปกติ, ลงทะเบียน,username ผิด, user ถูกระงับ)
	status := "login"
	// repo
	_, err := s.repo.GetPreRegisterByUserPassEmailRepo(ctx, nil, data.Username, data.Password, data.Email)
	if err != nil {
		// อาจจะทำเพิ่ม กรณีถ้า ไป search เจอบน table user ที่ถูกระงับ status = "suspend"
		// ถ้าไม่เจอใน table user ที่ถูกระงับ ให้ถือว่าเป็น user ใหม่
		status = "register"
		return status
	}
	// ถ้าหาเจอ แสดงว่า username,password,email ถูกต้อง เป็น user ที่ลงทะเบียนไว้ ให้ทำงานต่อ login
	return status
}

func (s *Service) LogoutService(ctx context.Context) (int, error) {
	ctx, span := s.repo.Trace(ctx, "svc.LogoutService", oteltrace.WithAttributes())
	defer span.End()
	cliams, _ := s.repo.GetAuthCtxRepo(ctx)
	if err := s.repo.ClearAuthSession(cliams.Uid.String(), cliams.PreRegisterUid); err != nil {
		logrus.Error("ClearAuthSession error ->", err)
		return http.StatusInternalServerError, errors.New("system error")
	}
	return http.StatusOK, nil
}
