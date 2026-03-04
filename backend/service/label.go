package service

import (
	"local-kanban/backend/model"
)

type LabelRepo interface {
	FindAll() ([]model.Label, error)
	Create(label *model.Label) error
}

type LabelService struct {
	Repo LabelRepo
}

func NewLabelService(repo LabelRepo) *LabelService {
	return &LabelService{Repo: repo}
}

func (s *LabelService) GetAll() ([]model.Label, error) {
	return s.Repo.FindAll()
}

func (s *LabelService) Create(label *model.Label) error {
	return s.Repo.Create(label)
}
