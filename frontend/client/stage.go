package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *BackendClient) GetStages() ([]Stage, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/stages")
	if err != nil {
		return nil, fmt.Errorf("ステージ一覧の取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ステージ一覧の取得でエラーが発生しました（ステータスコード: %d）", resp.StatusCode)
	}

	var stages []Stage
	if err := json.NewDecoder(resp.Body).Decode(&stages); err != nil {
		return nil, fmt.Errorf("ステージ一覧のJSONデコードに失敗しました: %w", err)
	}

	return stages, nil
}
