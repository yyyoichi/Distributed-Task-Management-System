package api

import (
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

func (s *SyncerMock) Sync(nextVersion int, todos []store.TodoDateset) SyncResponse {
	s.SyncNextVersion(nextVersion)
	for _, dataset := range todos {
		s.SyncTodoAt(dataset.ID, store.ConvertTodo(dataset))
	}
	return SyncResponse{nil}
}
