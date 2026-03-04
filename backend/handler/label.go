package handler

import (
	"net/http"

	"local-kanban/backend/model"
	"local-kanban/backend/service"

	"github.com/labstack/echo/v4"
)

type LabelHandler struct {
	Service *service.LabelService
}

func NewLabelHandler(svc *service.LabelService) *LabelHandler {
	return &LabelHandler{Service: svc}
}

func (h *LabelHandler) GetAll(c echo.Context) error {
	labels, err := h.Service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラベル一覧の取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, labels)
}

func (h *LabelHandler) Create(c echo.Context) error {
	var label model.Label
	if err := c.Bind(&label); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
	}

	if label.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "nameは必須です"})
	}

	if label.Color == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "colorは必須です"})
	}

	if err := h.Service.Create(&label); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラベルの作成に失敗しました"})
	}

	return c.JSON(http.StatusCreated, label)
}
