package repository

import (
	"local-kanban/backend/model"

	"gorm.io/gorm"
)

type StageRepository struct {
	DB *gorm.DB
}

func NewStageRepository(db *gorm.DB) *StageRepository {
	return &StageRepository{DB: db}
}

func (r *StageRepository) FindAll() ([]model.Stage, error) {
	var stages []model.Stage
	if err := r.DB.Order("position ASC").Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}

// FindByName はWIP制限チェック時に"DOING"ステージを特定するために使用
func (r *StageRepository) FindByName(name string) (*model.Stage, error) {
	var stage model.Stage
	if err := r.DB.Where("name = ?", name).First(&stage).Error; err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *StageRepository) Create(stage *model.Stage) error {
	return r.DB.Create(stage).Error
}
