package service

import (
	"errors"
	"local-kanban/backend/model"
	"testing"
)

// --- モック ---

type mockTaskRepo struct {
	countResult int64
	countErr    error
	findTask    *model.Task
	findErr     error
	createErr   error
	updateErr   error
}

func (m *mockTaskRepo) FindAll(stageID, labelID uint) ([]model.Task, error) { return nil, nil }
func (m *mockTaskRepo) FindByID(id uint) (*model.Task, error)               { return m.findTask, m.findErr }
func (m *mockTaskRepo) Create(task *model.Task) error                       { return m.createErr }
func (m *mockTaskRepo) Update(task *model.Task) error                       { return m.updateErr }
func (m *mockTaskRepo) UpdateStage(id, stageID uint) error                  { return nil }
func (m *mockTaskRepo) Delete(id uint) error                                { return nil }
func (m *mockTaskRepo) CreateWorkLog(workLog *model.WorkLog) error          { return nil }
func (m *mockTaskRepo) CreateNote(note *model.Note) error                   { return nil }

func (m *mockTaskRepo) CountByStageAndLabel(stageID, labelID, excludeTaskID uint) (int64, error) {
	return m.countResult, m.countErr
}

type mockStageFinder struct {
	stage *model.Stage
	err   error
}

func (m *mockStageFinder) FindByName(name string) (*model.Stage, error) {
	return m.stage, m.err
}

// --- テスト ---

func TestCreate_WIPLimit(t *testing.T) {
	doingStage := &model.Stage{Name: "DOING"}
	doingStage.ID = 1

	tests := []struct {
		name      string
		count     int64
		stageID   uint
		wantError bool
	}{
		{"DOINGに0件 → 作成OK", 0, 1, false},
		{"DOINGに4件 → 作成OK", 4, 1, false},
		{"DOINGに5件 → WIP制限で拒否", 5, 1, true},
		{"DOING以外 → 制限なし", 10, 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(
				&mockTaskRepo{countResult: tt.count},
				&mockStageFinder{stage: doingStage},
			)

			task := &model.Task{StageID: tt.stageID, LabelID: 1}
			err := svc.Create(task)

			if tt.wantError && !errors.Is(err, ErrWIPLimitExceeded) {
				t.Errorf("WIP制限エラーを期待したが、got: %v", err)
			}
			if !tt.wantError && err != nil {
				t.Errorf("エラーなしを期待したが、got: %v", err)
			}
		})
	}
}

func TestChangeStage_WIPLimit(t *testing.T) {
	doingStage := &model.Stage{Name: "DOING"}
	doingStage.ID = 1

	existingTask := &model.Task{LabelID: 1}
	existingTask.ID = 10

	tests := []struct {
		name       string
		count      int64
		newStageID uint
		wantError  bool
	}{
		{"DOINGに移動、4件 → OK", 4, 1, false},
		{"DOINGに移動、5件 → WIP制限で拒否", 5, 1, true},
		{"DOING以外に移動 → 制限なし", 10, 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(
				&mockTaskRepo{countResult: tt.count, findTask: existingTask},
				&mockStageFinder{stage: doingStage},
			)

			err := svc.ChangeStage(10, tt.newStageID)

			if tt.wantError && !errors.Is(err, ErrWIPLimitExceeded) {
				t.Errorf("WIP制限エラーを期待したが、got: %v", err)
			}
			if !tt.wantError && err != nil {
				t.Errorf("エラーなしを期待したが、got: %v", err)
			}
		})
	}
}

func TestCreate_DOINGStageNotExist(t *testing.T) {
	svc := NewTaskService(
		&mockTaskRepo{},
		&mockStageFinder{err: errors.New("record not found")},
	)

	task := &model.Task{StageID: 1, LabelID: 1}
	err := svc.Create(task)

	if err != nil {
		t.Errorf("DOINGステージ未作成時はチェックスキップすべきだが、got: %v", err)
	}
}
