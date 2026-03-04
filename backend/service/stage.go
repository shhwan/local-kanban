package service

import (
	"local-kanban/backend/model"
)

type StageRepo interface {
	FindAll() ([]model.Stage, error)
	FindByName(name string) (*model.Stage, error)
	Create(stage *model.Stage) error
}

type StageService struct {
	Repo StageRepo
}

func NewStageService(repo StageRepo) *StageService {
	return &StageService{Repo: repo}
}

func (s *StageService) GetAll() ([]model.Stage, error) {
	return s.Repo.FindAll()
}

func (s *StageService) Create(stage *model.Stage) error {
	return s.Repo.Create(stage)
}
