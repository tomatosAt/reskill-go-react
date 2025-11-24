package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/module/front-end/mapper"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (s *Service) DOBCheckSVC(ctx context.Context, dob string) (string, error) {
	ctx, span := s.repo.Trace(ctx, "svc.CheckFormatPreRegisterSVC", oteltrace.WithAttributes())
	defer span.End()
	_, errDob := time.Parse("2006-01-02", dob)
	if errDob != nil {
		return "", errors.New("birth_date must be YYYY-MM-DD format")
	}
	return dob, nil
}

func (s *Service) TitleNameCheckSVC(ctx context.Context, titleNameTH, titleNameENG string) error {
	ctx, span := s.repo.Trace(ctx, "svc.CheckFormatPreRegisterSVC", oteltrace.WithAttributes())
	defer span.End()
	if titleNameENG != "" && titleNameTH != "" {
		titleNameList := map[string]string{
			"Miss":   "นางสาว",
			"Mr":     "นาย",
			"Mrs":    "นาง",
			"Master": "เด็กชาย",
		}
		if _, ok := titleNameList[titleNameENG]; !ok {
			return errors.New("title name  must be 'Miss' or 'Mr'  or 'Mrs' or 'Master'")
		}
		if titleNameList[titleNameENG] != titleNameTH {
			return errors.New("title name Eng and Thai  not match ")
		} else if titleNameENG == "Miss" && titleNameTH == "เด็กหญิง" {
			return errors.New("title name Eng and Thai  not match ")
		}
	}
	return nil
}

func (s *Service) CheckFormatPreRegisterSVC(ctx context.Context, data *dto.PreRegisterDataBasePayload) error {
	ctx, span := s.repo.Trace(ctx, "svc.CheckFormatPreRegisterSVC", oteltrace.WithAttributes())
	defer span.End()
	// Input validation
	if _, err := util.ValidatorStruct(data); err != nil {
		return err
	}
	// Check birth_date format
	if _, err := s.DOBCheckSVC(ctx, data.Dob); err != nil {
		return err
	}
	// Check คำนำหน้าชื่อ
	if err := s.TitleNameCheckSVC(ctx, data.TitleNameTH, data.TitleNameEN); err != nil {
		return err
	}
	return nil
}

func (s *Service) PreRegisterSVC(ctx context.Context, data dto.PreRegisterDataBasePayload) (*dto.ResponsePreRegister, int, error) {
	ctx, span := s.repo.Trace(ctx, "svc.PreRegisterSVC", oteltrace.WithAttributes())
	defer span.End()
	// Process เก็บข้อมูล
	var result dto.ResponsePreRegister
	// เช็คข้อมูล pre-register ว่า username ซ้ำมั้ย
	repeatUsername := s.repo.RepeatByUsernameRepo(ctx, nil, data.Username)
	if repeatUsername {
		return &result, http.StatusBadRequest, errors.New("username already exists")
	}
	tx := s.repo.DB().Ctx().Begin()
	preRegisterRepo, err := s.repo.GetPreRegisterByEmailTelNoRepo(ctx, tx, data.Email, data.Tel)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return &result, http.StatusInternalServerError, err
		}
		//  Encrpy ชื่อ นามสกุล
		encryptedData, _ := util.EncryptList(s.repo.AppCfg().Secret.EncryptKey, data.FirstNameTH, data.LastNameTH, data.FirstNameEN, data.LastNameEN, data.TitleNameTH, data.TitleNameEN)
		hashPassword := util.HashSHA256(data.Password)
		mapperPreRegister := mapper.InsertPreRegisterMapper(data, hashPassword, encryptedData...)
		preRegisterRepo, err = s.repo.InsertPreRegisterRepo(ctx, tx, mapperPreRegister)
		if err != nil {
			tx.Rollback()
			util.RecordSpanError(span, err, "repo.InsertPreRegisterRepo")
			return &result, http.StatusInternalServerError, errors.New("server error")
		}
	}
	// : transaction auth
	transactionAuthMapper := mapper.InsertTransactionAuthMapper("pre-register", "pending", "pre-register", preRegisterRepo.Id.String())
	transactionAuthRepo, err := s.repo.InsertTransactionAuthRepo(ctx, tx, transactionAuthMapper)
	if err != nil {
		tx.Rollback()
		util.RecordSpanError(span, err, "repo.InsertTransactionAuthRepo")
		return &result, http.StatusInternalServerError, errors.New("server error")
	}
	tx.Commit()
	convernPreRegister := mapper.PreRegisterToConvernMapper(preRegisterRepo.Id, transactionAuthRepo.Id)
	convertByte, err := mapper.ConvernToByteMapper(convernPreRegister)
	if err != nil {
		logrus.Errorln("PreRegisterSVC ConvernToByteMapper err->", err)
		util.RecordSpanError(span, err, "repo.mapper.ConvernToByteMapper")
	}
	encryptedConvertByte, _ := util.EncryptAes256Ecb(s.repo.AppCfg().Secret.EncryptKey, string(convertByte))
	result = mapper.ResponsePreRegisterMapper(string(encryptedConvertByte))
	return &result, http.StatusOK, nil
}
