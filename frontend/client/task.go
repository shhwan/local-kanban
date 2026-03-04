package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *BackendClient) GetTasks() ([]Task, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/tasks")
	if err != nil {
		return nil, fmt.Errorf("タスク一覧の取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("タスク一覧の取得でエラーが発生しました（ステータスコード: %d）", resp.StatusCode)
	}

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("タスク一覧のJSONデコードに失敗しました: %w", err)
	}

	return tasks, nil
}

func (c *BackendClient) GetTask(id uint) (*Task, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/tasks/%d", c.BaseURL, id))
	if err != nil {
		return nil, fmt.Errorf("タスク(ID=%d)の取得に失敗しました: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("タスク(ID=%d)の取得でエラーが発生しました（ステータスコード: %d）", id, resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("タスク(ID=%d)のJSONデコードに失敗しました: %w", id, err)
	}

	return &task, nil
}

type CreateTaskRequest struct {
	Title   string `json:"title"`
	LabelID uint   `json:"label_id"`
	StageID uint   `json:"stage_id"`
}

func (c *BackendClient) CreateTask(title string, labelID, stageID uint) (*Task, error) {
	reqBody := CreateTaskRequest{
		Title:   title,
		LabelID: labelID,
		StageID: stageID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("タスク作成リクエストのJSONエンコードに失敗しました: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/tasks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("タスクの作成に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("タスクの作成でエラーが発生しました（ステータスコード: %d）", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("作成されたタスクのJSONデコードに失敗しました: %w", err)
	}

	return &task, nil
}

type UpdateTaskRequest struct {
	Title   *string `json:"title,omitempty"`
	LabelID *uint   `json:"label_id,omitempty"`
	StageID *uint   `json:"stage_id,omitempty"`
}

func (c *BackendClient) UpdateTask(id uint, req UpdateTaskRequest) (*Task, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("タスク更新リクエストのJSONエンコードに失敗しました: %w", err)
	}

	// Backend は PUT /api/tasks/:id を期待するためPUTを使用
	httpReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/tasks/%d", c.BaseURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("タスク更新リクエストの作成に失敗しました: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("タスク(ID=%d)の更新に失敗しました: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("タスク(ID=%d)の更新でエラーが発生しました（ステータスコード: %d）", id, resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("更新されたタスク(ID=%d)のJSONデコードに失敗しました: %w", id, err)
	}

	return &task, nil
}

type ChangeStageRequest struct {
	StageID uint `json:"stage_id"`
}

func (c *BackendClient) ChangeStage(id, stageID uint) error {
	reqBody := ChangeStageRequest{StageID: stageID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ステージ変更リクエストのJSONエンコードに失敗しました: %w", err)
	}

	// http.Clientには直接PATCHメソッドがないためNewRequestを使用
	httpReq, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/tasks/%d/stage", c.BaseURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ステージ変更リクエストの作成に失敗しました: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("タスク(ID=%d)のステージ変更に失敗しました: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("タスク(ID=%d)のステージ変更でエラーが発生しました（ステータスコード: %d）", id, resp.StatusCode)
	}

	return nil
}

func (c *BackendClient) DeleteTask(id uint) error {
	httpReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/tasks/%d", c.BaseURL, id), nil)
	if err != nil {
		return fmt.Errorf("タスク削除リクエストの作成に失敗しました: %w", err)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("タスク(ID=%d)の削除に失敗しました: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("タスク(ID=%d)の削除でエラーが発生しました（ステータスコード: %d）", id, resp.StatusCode)
	}

	return nil
}

type AddWorkLogRequest struct {
	Content string `json:"content"`
}

func (c *BackendClient) AddWorkLog(taskID uint, content string) (*WorkLog, error) {
	reqBody := AddWorkLogRequest{Content: content}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("作業ログ追加リクエストのJSONエンコードに失敗しました: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/api/tasks/%d/logs", c.BaseURL, taskID),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("タスク(ID=%d)への作業ログ追加に失敗しました: %w", taskID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("作業ログの追加でエラーが発生しました（ステータスコード: %d）", resp.StatusCode)
	}

	var workLog WorkLog
	if err := json.NewDecoder(resp.Body).Decode(&workLog); err != nil {
		return nil, fmt.Errorf("作成された作業ログのJSONデコードに失敗しました: %w", err)
	}

	return &workLog, nil
}

type AddNoteRequest struct {
	Content string `json:"content"`
}

func (c *BackendClient) AddNote(taskID uint, content string) (*Note, error) {
	reqBody := AddNoteRequest{Content: content}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("特記事項追加リクエストのJSONエンコードに失敗しました: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/api/tasks/%d/notes", c.BaseURL, taskID),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("タスク(ID=%d)への特記事項追加に失敗しました: %w", taskID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("特記事項の追加でエラーが発生しました（ステータスコード: %d）", resp.StatusCode)
	}

	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return nil, fmt.Errorf("作成された特記事項のJSONデコードに失敗しました: %w", err)
	}

	return &note, nil
}
