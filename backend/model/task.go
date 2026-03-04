package model

import "time"

type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null;size:255"`
	LabelID   uint      `json:"label_id" gorm:"not null"`
	Label     Label     `json:"label" gorm:"foreignKey:LabelID"`
	StageID   uint      `json:"stage_id" gorm:"not null"`
	Stage     Stage     `json:"stage" gorm:"foreignKey:StageID"`
	WorkLogs  []WorkLog `json:"work_logs" gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	Notes     []Note    `json:"notes" gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
