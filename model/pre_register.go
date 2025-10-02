package model

type PreRegister struct {
	Model
	Username       string `gorm:"column:username;size:255;index"`
	TitleTh        string `gorm:"column:title_th;size:255;"`
	TitleEng       string `gorm:"column:title_eng;size:255;"`
	Email          string `gorm:"column:email;size:255;index"`
	Tel            string `gorm:"column:tel;size:10;index"`
	FirstNameTh    string `gorm:"column:first_name_th;size:255;"`
	LastNameTh     string `gorm:"column:last_name_th;size:255;"`
	FirstNameEng   string `gorm:"column:first_name_eng;size:255;"`
	LastNameEng    string `gorm:"column:last_name_eng;size:255;"`
	NickName       string `gorm:"column:nick_name;size:50;"`
	Password       string `gorm:"column:password;size:255;"`
	RegisterStatus string `gorm:"column:register_status;size:255;"`
}
