package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *BackendClient) GetLabels() ([]Label, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/labels")
	if err != nil {
		return nil, fmt.Errorf("ラベル一覧の取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ラベル一覧の取得でエラーが発生しました（ステータスコード: %d）", resp.StatusCode)
	}

	var labels []Label
	if err := json.NewDecoder(resp.Body).Decode(&labels); err != nil {
		return nil, fmt.Errorf("ラベル一覧のJSONデコードに失敗しました: %w", err)
	}

	return labels, nil
}
