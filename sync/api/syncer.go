package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"
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

// [currentSyncVersion]以上のバージョンを持つデータストア差異情報を取得する
func (s *Syncer) GetDifference(currentSyncVersion int) DiffResponse {
	reqBody, err := json.Marshal(struct {
		Version int `json:"version"`
	}{
		Version: currentSyncVersion,
	})
	if err != nil {
		return DiffResponse{Err: err}
	}
	resp, err := http.Post(fmt.Sprintf("%s/differences", s.url), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return DiffResponse{Err: err}
	}
	defer resp.Body.Close()

	var data []document.TodoDataset
	resBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(resBody, &data); err != nil {
		return DiffResponse{Err: err}
	}

	return DiffResponse{TodoDatasets: data, Err: nil}
}

// 同期を実行する
func (s *Syncer) Synchronize(currentSyncVersion int, todos []document.TodoDataset) SynchronizeResponse {
	reqBody, err := json.Marshal(struct {
		Version int                    `json:"version"`
		Todos   []document.TodoDataset `json:"todos"`
	}{
		Version: currentSyncVersion,
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
