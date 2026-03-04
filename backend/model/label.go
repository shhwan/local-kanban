package model

type Label struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"not null;size:20;uniqueIndex"`
	Color string `json:"color" gorm:"not null;size:7"` // #RRGGBB
}
