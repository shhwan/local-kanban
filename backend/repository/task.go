package repository

import (
	"local-kanban/backend/model"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) FindAll(stageID, labelID uint) ([]model.Task, error) {
	var tasks []model.Task
	query := r.DB.Preload("Label").Preload("Stage").Preload("WorkLogs").Preload("Notes")

	if stageID > 0 {
		query = query.Where("stage_id = ?", stageID)
	}
	if labelID > 0 {
		query = query.Where("label_id = ?", labelID)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.
		Preload("Label").
		Preload("Stage").
		Preload("WorkLogs").
		Preload("Notes").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Create(task *model.Task) error {
	return r.DB.Create(task).Error
}

func (r *TaskRepository) Update(task *model.Task) error {
	return r.DB.Save(task).Error
}

func (r *TaskRepository) UpdateStage(id, stageID uint) error {
	return r.DB.Model(&model.Task{}).Where("id = ?", id).Update("stage_id", stageID).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.DB.Delete(&model.Task{}, id).Error
}

// CountByStageAndLabel はWIP制限チェック用。excludeTaskID > 0 の場合、自身を除外してカウントする。
func (r *TaskRepository) CountByStageAndLabel(stageID, labelID, excludeTaskID uint) (int64, error) {
	var count int64
	query := r.DB.Model(&model.Task{}).Where("stage_id = ? AND label_id = ?", stageID, labelID)

	if excludeTaskID > 0 {
		query = query.Where("id != ?", excludeTaskID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TaskRepository) CreateWorkLog(workLog *model.WorkLog) error {
	return r.DB.Create(workLog).Error
}

func (r *TaskRepository) CreateNote(note *model.Note) error {
	return r.DB.Create(note).Error
}
