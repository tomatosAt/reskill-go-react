package model

type User struct {
	Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // hashed
	Role     string `gorm:"default:user" json:"role"`
}
