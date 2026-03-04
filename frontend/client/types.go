package client

import "time"

type Task struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	LabelID   uint      `json:"label_id"`
	Label     Label     `json:"label"`
	StageID   uint      `json:"stage_id"`
	Stage     Stage     `json:"stage"`
	WorkLogs  []WorkLog `json:"work_logs"`
	Notes     []Note    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Label struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Stage struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type WorkLog struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Note struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
