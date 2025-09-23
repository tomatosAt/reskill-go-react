package model

type PreRegister struct {
	Model
	Username     string `gorm:"column:username;size:255;"`
	TitleTh      string `gorm:"column:title_th;size:255;"`
	TitleEng     string `gorm:"column:title_eng;size:255;"`
	Email        string `gorm:"column:email;size:255;"`
	EmailOneId   string `gorm:"column:email_one_id;size:255;"`
	Tel          string `gorm:"column:tel;size:10;"`
	FirstNameTh  string `gorm:"column:first_name_th;size:255;"`
	LastNameTh   string `gorm:"column:last_name_th;size:255;"`
	FirstNameEng string `gorm:"column:first_name_eng;size:255;"`
	LastNameEng  string `gorm:"column:last_name_eng;size:255;"`
	NickName     string `gorm:"column:nick_name;size:50;"`
}
