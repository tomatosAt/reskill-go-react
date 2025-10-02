package model

type Game struct {
	Model
	Title    string  `gorm:"column:title;size:255;"`
	Category string  `gorm:"column:category;size:255;"`
	Price    float64 `gorm:"column:price;"`
	Rating   float64 `gorm:"column:rating;"`
	CoverURL string  `gorm:"column:cover_url;size:255;"`
}
