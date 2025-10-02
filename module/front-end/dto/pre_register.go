package dto

type PreRegisterDataBasePayload struct {
	FirstNameTH string `json:"first_name_th" validate:"omitempty,max=100,regexp=^[\u0E00-\u0E4C .-]+$"`
	LastNameTH  string `json:"last_name_th" validate:"omitempty,max=100,regexp=^[\u0E00-\u0E4C .-]*$"`
	TitleNameTH string `json:"title_name_th" validate:"omitempty,max=100,regexp=^[\u0E00-\u0E4C]*$"`
	Dob         string `json:"birth_date" validate:"regexp=^[0-9 .-]*$"`
	FirstNameEN string `json:"first_name_en" validate:"omitempty,max=100,regexp=^[a-zA-Z .-]+$"`
	LastNameEN  string `json:"last_name_en" validate:"omitempty,max=100,regexp=^[a-zA-Z .-]*$"`
	TitleNameEN string `json:"title_name_en" validate:"omitempty,max=100,regexp=^[a-zA-Z .-]*$"`
	Tel         string `json:"tel" validate:"regexp=^[0-9]*$"`
	Remark      string `json:"remark" validate:"omitempty,max=100,regexp=^[a-zA-Z .-]+$"`
	Email       string `json:"email" validate:"omitempty,email"`
	Username    string `json:"username" validate:"omitempty,max=50,regexp=^[a-zA-Z0-9_@ .-]+$"`
	NickName    string `json:"nick_name" validate:"omitempty,max=50,regexp=^.+$"`
	Password    string `json:"password" validate:"required,min=8,max=50"`
	UserLogin   string `json:"user_login"`
}
