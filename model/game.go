package model

type Game struct {
	Model
	Title    string  `gorm:"not null" json:"title"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Rating   float64 `json:"rating"`
	CoverURL string  `json:"cover_url"`
}
