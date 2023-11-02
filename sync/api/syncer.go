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

// [currentVersion]以上のバージョンを持つデータストア差異情報を取得する
func (s *Syncer) GetDifference(currentVersion int) DiffResponse {
	reqBody := []byte(fmt.Sprintf(`{"version":%d`, currentVersion))
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
func (s *Syncer) Synchronize(currentVersion int, todos []store.TodoDateset) SynchronizeResponse {
	reqBody, err := json.Marshal(struct {
		Version int                 `json:"version"`
		Todos   []store.TodoDateset `json:"todos"`
	}{
		Version: currentVersion,
		Todos:   todos,
	})
	if err != nil {
		return SynchronizeResponse{err}
	}

	resp, err := http.Post(fmt.Sprintf("%s/sync", s.url), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return SynchronizeResponse{err}
	}
	defer resp.Body.Close()
	return SynchronizeResponse{nil}
}
