package handler

import (
	"net/http"

	"local-kanban/backend/model"
	"local-kanban/backend/service"

	"github.com/labstack/echo/v4"
)

type StageHandler struct {
	Service *service.StageService
}

func NewStageHandler(svc *service.StageService) *StageHandler {
	return &StageHandler{Service: svc}
}

func (h *StageHandler) GetAll(c echo.Context) error {
	stages, err := h.Service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ステージ一覧の取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, stages)
}

func (h *StageHandler) Create(c echo.Context) error {
	var stage model.Stage
	if err := c.Bind(&stage); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}

	if stage.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "nameは必須です"})
	}

	if err := h.Service.Create(&stage); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ステージの作成に失敗しました"})
	}

	return c.JSON(http.StatusCreated, stage)
}
