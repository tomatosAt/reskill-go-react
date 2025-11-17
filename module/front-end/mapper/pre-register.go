package mapper

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tomatosAt/reskill-go-react/model"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
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
