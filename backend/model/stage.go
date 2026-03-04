package model

type Stage struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;size:20;uniqueIndex"`
	Position int    `json:"position" gorm:"not null"`
}
