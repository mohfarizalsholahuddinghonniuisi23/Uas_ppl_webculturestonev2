package models

type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"` // BUG-02 FIX: json:"-" memastikan hash password tidak pernah tampil di response API
}
