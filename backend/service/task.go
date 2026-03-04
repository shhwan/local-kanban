package service

import (
	"errors"
	"local-kanban/backend/model"
	"local-kanban/backend/repository"
)

const WIPLimit = 2
const DOINGStageName = "DOING"

var ErrWIPLimitExceeded = errors.New("WIP制限超過: DOINGステージにはラベルごとに最大2つまでのタスクしか配置できません")

type TaskService struct {
	TaskRepo  *repository.TaskRepository
	StageRepo *repository.StageRepository
}

func NewTaskService(taskRepo *repository.TaskRepository, stageRepo *repository.StageRepository) *TaskService {
	return &TaskService{
		TaskRepo:  taskRepo,
		StageRepo: stageRepo,
	}
}

func (s *TaskService) GetAll(stageID, labelID uint) ([]model.Task, error) {
	return s.TaskRepo.FindAll(stageID, labelID)
}

func (s *TaskService) GetByID(id uint) (*model.Task, error) {
	return s.TaskRepo.FindByID(id)
}

func (s *TaskService) Create(task *model.Task) error {
	if err := s.checkWIPLimit(task.StageID, task.LabelID, 0); err != nil {
		return err
	}
	return s.TaskRepo.Create(task)
}

func (s *TaskService) Update(task *model.Task) error {
	if err := s.checkWIPLimit(task.StageID, task.LabelID, task.ID); err != nil {
		return err
	}
	return s.TaskRepo.Update(task)
}

func (s *TaskService) ChangeStage(taskID, newStageID uint) error {
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		return err
	}

	if err := s.checkWIPLimit(newStageID, task.LabelID, taskID); err != nil {
		return err
	}

	return s.TaskRepo.UpdateStage(taskID, newStageID)
}

func (s *TaskService) Delete(id uint) error {
	return s.TaskRepo.Delete(id)
}

func (s *TaskService) AddWorkLog(workLog *model.WorkLog) error {
	return s.TaskRepo.CreateWorkLog(workLog)
}

func (s *TaskService) AddNote(note *model.Note) error {
	return s.TaskRepo.CreateNote(note)
}

// checkWIPLimit はDOINGステージのWIP制限をチェックする。
// DOINGステージが未作成の場合はスキップする。
func (s *TaskService) checkWIPLimit(stageID, labelID, excludeTaskID uint) error {
	doingStage, err := s.StageRepo.FindByName(DOINGStageName)
	if err != nil {
		return nil // DOINGステージが存在しない場合はチェック不要
	}

	if stageID != doingStage.ID {
		return nil
	}

	count, err := s.TaskRepo.CountByStageAndLabel(stageID, labelID, excludeTaskID)
	if err != nil {
		return err
	}

	if count >= WIPLimit {
		return ErrWIPLimitExceeded
	}

	return nil
}
