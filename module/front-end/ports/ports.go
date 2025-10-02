package ports

import "github.com/tomatosAt/reskill-go-react/module/front-end/dto"

type Repository interface {
}

type Service interface {
	// pre register
	DOBCheckSVC(dob string) (string, error)
	TitleNameCheckSVC(titleNameTH, titleNameENG string) error
	PreRegisterSVC(data dto.PreRegisterDataBasePayload) error
	CheckFormatPreRegisterSVC(data *dto.PreRegisterDataBasePayload) error
}
