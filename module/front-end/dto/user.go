package dto

type ResponseUserProfiles struct {
	Id           string `json:"id"`
	UserID       string `json:"user_id"`
	FirstNameTh  string `json:"first_name_th"`
	FirstNameEng string `json:"first_name_eng"`
	LastNameTh   string `json:"last_name_th"`
	LastNameEng  string `json:"last_name_eng"`
	MobileNo     string `json:"mobile_no"`
	Email        string `json:"email"`
	TitleTh      string `json:"title_th"`
	TitleEng     string `json:"title_eng"`
	NickName     string `json:"nick_name"`
}
