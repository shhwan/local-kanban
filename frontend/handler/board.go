package handler

import (
	"net/http"
	"strconv"

	"local-kanban/frontend/client"
	"local-kanban/frontend/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type BoardHandler struct {
	Client *client.BackendClient
}

func NewBoardHandler(c *client.BackendClient) *BoardHandler {
	return &BoardHandler{Client: c}
}

func render(c echo.Context, statusCode int, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(statusCode)
	return component.Render(c.Request().Context(), c.Response())
}

func (h *BoardHandler) HandleBoard(c echo.Context) error {
	// API未準備の場合（初回起動時など）に備えて空スライスで続行する
	tasks, err := h.Client.GetTasks()
	if err != nil {
		tasks = []client.Task{}
	}

	labels, err := h.Client.GetLabels()
	if err != nil {
		labels = []client.Label{}
	}

	stages, err := h.Client.GetStages()
	if err != nil {
		stages = []client.Stage{}
	}

	data := templates.BoardData{
		Stages: stages,
		Tasks:  tasks,
		Labels: labels,
	}

	return render(c, http.StatusOK, templates.Board(data))
}

func (h *BoardHandler) HandleBoardPartial(c echo.Context) error {
	return h.HandleBoard(c)
}

func (h *BoardHandler) HandleCreateTask(c echo.Context) error {
	title := c.FormValue("title")

	labelIDStr := c.FormValue("label_id")
	labelID, err := strconv.ParseUint(labelIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なラベルIDです")
	}

	stageIDStr := c.FormValue("stage_id")
	stageID, err := strconv.ParseUint(stageIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なステージIDです")
	}

	_, err = h.Client.CreateTask(title, uint(labelID), uint(stageID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "タスクの作成に失敗しました: "+err.Error())
	}

	return h.HandleBoardPartial(c)
}

func (h *BoardHandler) HandleChangeStage(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なタスクIDです")
	}

	stageIDStr := c.QueryParam("stage_id")
	stageID, err := strconv.ParseUint(stageIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なステージIDです")
	}

	err = h.Client.ChangeStage(uint(id), uint(stageID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "ステージの変更に失敗しました: "+err.Error())
	}

	return h.HandleBoardPartial(c)
}

func (h *BoardHandler) HandleAddWorkLog(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なタスクIDです")
	}

	content := c.FormValue("content")
	if content == "" {
		return c.String(http.StatusBadRequest, "作業内容を入力してください")
	}

	_, err = h.Client.AddWorkLog(uint(id), content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "作業ログの追加に失敗しました: "+err.Error())
	}

	return h.HandleBoardPartial(c)
}

func (h *BoardHandler) HandleAddNote(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なタスクIDです")
	}

	content := c.FormValue("content")
	if content == "" {
		return c.String(http.StatusBadRequest, "特記事項を入力してください")
	}

	_, err = h.Client.AddNote(uint(id), content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "特記事項の追加に失敗しました: "+err.Error())
	}

	return h.HandleBoardPartial(c)
}

func (h *BoardHandler) HandleDeleteTask(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なタスクIDです")
	}

	err = h.Client.DeleteTask(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "タスクの削除に失敗しました: "+err.Error())
	}

	return h.HandleBoardPartial(c)
}
