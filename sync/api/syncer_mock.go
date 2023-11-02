package api

import (
	"context"
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
)

func NewSyncerMock() *SyncerMock {
	return &SyncerMock{store.NewStore()}
}

type SyncerMock struct {
	*store.TStore
}

func (s *SyncerMock) Me() string { return "MockSyncer" }

func (s *SyncerMock) GetDifference(latestVersion int) DiffResponse {
	resp := DiffResponse{}
	for id, todo := range s.GetLatestVersionTodo(latestVersion) {
		resp.TodoDatasets = append(resp.TodoDatasets, store.ConvertTodoDataset(id, todo))
	}
	return resp
}

func (s *SyncerMock) Sync(currentVersion int, todos []store.TodoDateset) SyncResponse {
	s.TStore.Sync(context.Background(), currentVersion, todos)
	return SyncResponse{nil}
}
