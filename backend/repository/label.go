package repository

import (
	"local-kanban/backend/model"

	"gorm.io/gorm"
)

type LabelRepository struct {
	DB *gorm.DB
}

func NewLabelRepository(db *gorm.DB) *LabelRepository {
	return &LabelRepository{DB: db}
}

func (r *LabelRepository) FindAll() ([]model.Label, error) {
	var labels []model.Label
	if err := r.DB.Find(&labels).Error; err != nil {
		return nil, err
	}
	return labels, nil
}

func (r *LabelRepository) Create(label *model.Label) error {
	return r.DB.Create(label).Error
}
