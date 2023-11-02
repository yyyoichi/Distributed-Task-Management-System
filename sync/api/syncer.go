package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
)

func NewSyncer(url string) *Syncer {
	return &Syncer{
		url: url,
	}
}

type Syncer struct {
	url string
}

func (s *Syncer) Me() string { return fmt.Sprintf("Syncer[%s]", s.url) }

// [latestVersion]よりも大きいバージョンを持つデータストア更新情報訪を[uri]から取得する
func (s *Syncer) GetDifference(latestVersion int) DiffResponse {
	reqBody := []byte(fmt.Sprintf(`{"version":%d`, latestVersion))
	resp, err := http.Post(fmt.Sprintf("%s/differences", s.url), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return DiffResponse{Err: err}
	}
	defer resp.Body.Close()

	var data []store.TodoDateset
	resBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(resBody, &data); err != nil {
		return DiffResponse{Err: err}
	}

	return DiffResponse{TodoDatasets: data, Err: nil}
}

// 同期を実行する
func (s *Syncer) Sync(currentVersion int, todos []store.TodoDateset) SyncResponse {
	reqBody, err := json.Marshal(struct {
		Version int                 `json:"version"`
		Todos   []store.TodoDateset `json:"todos"`
	}{
		Version: currentVersion,
		Todos:   todos,
	})
	if err != nil {
		return SyncResponse{err}
	}

	resp, err := http.Post(fmt.Sprintf("%s/sync", s.url), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return SyncResponse{err}
	}
	defer resp.Body.Close()
	return SyncResponse{nil}
}
