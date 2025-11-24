package mapper

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tomatosAt/reskill-go-react/model"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func InsertPreRegisterMapper(payload dto.PreRegisterDataBasePayload, password string, data ...string) model.PreRegister {
	return model.PreRegister{
		FirstNameTh:    data[0],
		LastNameTh:     data[1],
		FirstNameEng:   data[2],
		LastNameEng:    data[3],
		TitleTh:        data[4],
		TitleEng:       data[5],
		Password:       password,
		Email:          payload.Email,
		Tel:            payload.Tel,
		NickName:       payload.NickName,
		Username:       payload.Username,
		RegisterStatus: "pre-register",
		Model: model.Model{
			CreatedBy: "Admin",
			UpdatedBy: "Admin",
		},
	}
}

func InsertTransactionAuthMapper(action, status, remark string, preRegisterUUID string) model.TransactionAuth {
	return model.TransactionAuth{
		Action:         action,
		StatusRegister: status,
		Remark:         remark,
		PreRegisterUid: preRegisterUUID,
		Model: model.Model{
			CreatedBy: "Admin",
			UpdatedBy: "Admin",
		},
	}
}

func PreRegisterToConvernMapper(preRegisterUUID, transactionAuthUUID uuid.UUID) dto.PreRegisterConvern {
	timeNow := time.Now()
	return dto.PreRegisterConvern{
		PreRegisterUUID:     preRegisterUUID.String(),
		TransactionAuthUUID: transactionAuthUUID.String(),
		ExpireAt:            timeNow.Add(time.Hour),
	}
}

func ConvernToByteMapper(v any) ([]byte, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ResponsePreRegisterMapper(code string) dto.ResponsePreRegister {
	return dto.ResponsePreRegister{
		Code: code,
	}
}

func DecodePreRegisterConvernMapper(key string, data model.PreRegister) (model.PreRegister, error) {
	decryptedData, _ := util.DecryptList(key, data.FirstNameTh, data.LastNameTh, data.FirstNameEng, data.LastNameEng, data.TitleTh, data.TitleEng)
	return model.PreRegister{
		FirstNameTh:  decryptedData[0],
		LastNameTh:   decryptedData[1],
		FirstNameEng: decryptedData[2],
		LastNameEng:  decryptedData[3],
		TitleTh:      decryptedData[4],
		TitleEng:     decryptedData[5],
	}, nil
}

func ResponseProfileMapper(id, userID, tel, email, nickname string, datadecryp model.PreRegister) dto.ResponseUserProfiles {
	return dto.ResponseUserProfiles{
		Id:           id,
		UserID:       userID,
		FirstNameTh:  datadecryp.FirstNameTh,
		LastNameTh:   datadecryp.LastNameTh,
		TitleTh:      datadecryp.TitleTh,
		FirstNameEng: datadecryp.FirstNameEng,
		LastNameEng:  datadecryp.LastNameEng,
		TitleEng:     datadecryp.TitleEng,
		MobileNo:     tel,
		Email:        email,
		NickName:     nickname,
	}
}

func InsertUsersAuthMapper(username, password, role string) model.User {
	return model.User{
		Username: username,
		Password: password,
		Role:     role,
		Model: model.Model{
			CreatedBy: "Admin",
			UpdatedBy: "Admin",
		},
	}
}
