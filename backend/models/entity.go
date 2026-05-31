package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string
	Description string
	ImageURL    string
}

type ProductCategory struct {
	gorm.Model
	Name string `json:"Name"`
}

type Product struct {
	gorm.Model
	Name        string
	Description string
	Category    string
	ImageURL    string
	Quality     string
	Size        string
	Finishing   string
}

type Portfolio struct {
	gorm.Model
	Title       string `json:"Title"`
	Description string `json:"Description"`
	ImageURL    string `json:"ImageURL"`
}

type Message struct {
	gorm.Model
	Name    string
	Phone   string
	Email   string
	Message string
	IsRead  bool
}

type Visitor struct {
	ID    uint `gorm:"primarykey"`
	Date  string
	Count int
}

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

type Testimoni struct {
	gorm.Model
	ClientName      string     `gorm:"not null" json:"ClientName"`
	TestimonialText string     `gorm:"type:text;not null" json:"TestimonialText"`
	ImagePath       string     `gorm:"default:''" json:"ImagePath"`
	VideoPath       string     `gorm:"default:''" json:"VideoPath"`
	PortfolioID     *uint      `gorm:"constraint:OnDelete:SET NULL;type:integer" json:"PortfolioID"`
	Portfolio       *Portfolio `gorm:"foreignKey:PortfolioID" json:"Portfolio"`
}
