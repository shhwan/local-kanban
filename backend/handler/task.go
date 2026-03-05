package handler

import (
	"errors"
	"net/http"
	"strconv"

	"local-kanban/backend/model"
	"local-kanban/backend/service"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	Service *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{Service: svc}
}

func (h *TaskHandler) GetAll(c echo.Context) error {
	var stageID, labelID uint
	if v := c.QueryParam("stage_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "stage_idが不正です"})
		}
		stageID = uint(id)
	}

	if v := c.QueryParam("label_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "label_idが不正です"})
		}
		labelID = uint(id)
	}

	tasks, err := h.Service.GetAll(stageID, labelID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスク一覧の取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
	}

	task, err := h.Service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "タスクが見つかりません"})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Create(c echo.Context) error {
	var task model.Task
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}

	if task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "タイトルは必須です"})
	}

	if err := h.Service.Create(&task); err != nil {
		if errors.Is(err, service.ErrWIPLimitExceeded) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスクの作成に失敗しました"})
	}

	created, err := h.Service.GetByID(task.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "作成後のタスク取得に失敗しました"})
	}

	return c.JSON(http.StatusCreated, created)
}

func (h *TaskHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
	}

	existing, err := h.Service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "タスクが見つかりません"})
	}

	if err := c.Bind(existing); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}
	existing.ID = uint(id)

	if err := h.Service.Update(existing); err != nil {
		if errors.Is(err, service.ErrWIPLimitExceeded) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスクの更新に失敗しました"})
	}

	updated, err := h.Service.GetByID(existing.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "更新後のタスク取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, updated)
}

type changeStageRequest struct {
	StageID uint `json:"stage_id"`
}

func (h *TaskHandler) ChangeStage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
	}

	var req changeStageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}

	if req.StageID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "stage_idは必須です"})
	}

	if err := h.Service.ChangeStage(uint(id), req.StageID); err != nil {
		if errors.Is(err, service.ErrWIPLimitExceeded) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "タスクが見つかりません"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ステージの変更に失敗しました"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ステージを変更しました"})
}

func (h *TaskHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスクの削除に失敗しました"})
	}

	return c.NoContent(http.StatusNoContent)
}

type addWorkLogRequest struct {
	Content string `json:"content"`
}

func (h *TaskHandler) AddWorkLog(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
	}

	var req addWorkLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "contentは必須です"})
	}

	workLog := model.WorkLog{
		TaskID:  uint(id),
		Content: req.Content,
	}

	if err := h.Service.AddWorkLog(&workLog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "作業ログの追加に失敗しました"})
	}

	return c.JSON(http.StatusCreated, workLog)
}

type addNoteRequest struct {
	Content string `json:"content"`
}

func (h *TaskHandler) AddNote(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IDが不正です"})
	}

	var req addNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "contentは必須です"})
	}

	note := model.Note{
		TaskID:  uint(id),
		Content: req.Content,
	}

	if err := h.Service.AddNote(&note); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "特記事項の追加に失敗しました"})
	}

	return c.JSON(http.StatusCreated, note)
}
