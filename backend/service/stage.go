package service

import (
	"local-kanban/backend/model"
	"local-kanban/backend/repository"
)

type StageService struct {
	Repo *repository.StageRepository
}

func NewStageService(repo *repository.StageRepository) *StageService {
	return &StageService{Repo: repo}
}

func (s *StageService) GetAll() ([]model.Stage, error) {
	return s.Repo.FindAll()
}

func (s *StageService) Create(stage *model.Stage) error {
	return s.Repo.Create(stage)
}
