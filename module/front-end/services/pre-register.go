package services

import (
	"errors"
	"time"

	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (s *Service) DOBCheckSVC(dob string) (string, error) {
	_, errDob := time.Parse("2006-01-02", dob)
	if errDob != nil {
		return "", errors.New("birth_date must be YYYY-MM-DD format")
	}
	return dob, nil
}

func (s *Service) TitleNameCheckSVC(titleNameTH, titleNameENG string) error {
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

func (s *Service) CheckFormatPreRegisterSVC(data *dto.PreRegisterDataBasePayload) error {
	// Input validation
	if _, err := util.ValidatorStruct(data); err != nil {
		return err
	}
	// Check birth_date format
	if _, err := s.DOBCheckSVC(data.Dob); err != nil {
		return err
	}
	// Check คำนำหน้าชื่อ
	if err := s.TitleNameCheckSVC(data.TitleNameTH, data.TitleNameEN); err != nil {
		return err
	}
	return nil
}

func (s *Service) PreRegisterSVC(data dto.PreRegisterDataBasePayload) error {

	return nil
}
