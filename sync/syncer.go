package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TodoDataset struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
	Deleted   bool   `json:"deleted"`
	Version   int    `json:"version"`
}

type Syncer struct {
	url string
}

func NewSyncer(url string) Syncer {
	return Syncer{
		url: url,
	}
}

// [latestVersion]よりも大きいバージョンを持つデータストア更新情報訪を[uri]から取得する
func (s *Syncer) GetDiff(latestVersion int) ([]TodoDataset, error) {
	reqBody := []byte(fmt.Sprintf(`{"version":%d`, latestVersion))
	resp, err := http.Post(fmt.Sprintf("%s/differences", s.url), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []TodoDataset
	resBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(resBody, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// 同期を実行する
func (s *Syncer) Sync(nextVersion int, todos []TodoDataset) error {
	reqBody, err := json.Marshal(struct {
		NextVersion int           `json:"nextVersion"`
		Todos       []TodoDataset `json:"todos"`
	}{
		NextVersion: nextVersion,
		Todos:       todos,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/sync", s.url), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
