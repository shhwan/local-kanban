package service

import (
	"local-kanban/backend/model"
	"local-kanban/backend/repository"
)

type LabelService struct {
	Repo *repository.LabelRepository
}

func NewLabelService(repo *repository.LabelRepository) *LabelService {
	return &LabelService{Repo: repo}
}

func (s *LabelService) GetAll() ([]model.Label, error) {
	return s.Repo.FindAll()
}

func (s *LabelService) Create(label *model.Label) error {
	return s.Repo.Create(label)
}
