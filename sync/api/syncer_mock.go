package api

import (
	"context"
	"fmt"
	"log"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/store"
)

func NewSyncerMock(name string) *SyncerMock {
	return &SyncerMock{name, store.NewStore()}
}

type SyncerMock struct {
	name string
	*store.TStore
}

func (s *SyncerMock) Me() string {
	return fmt.Sprintf("SyncerMock[%s]: Key-Value-Store has %d itmes", s.name, len(s.ByID))
}

func (s *SyncerMock) GetDifference(currentSyncVersion int) DiffResponse {
	resp := DiffResponse{}

	todos := s.TStore.GetLatestVersionTodo(currentSyncVersion)
	log.Printf("\t[%s] There are %d Todos(/%d) defferenced from v%d", s.name, len(todos), len(s.TStore.ByID), currentSyncVersion)
	for id, todo := range todos {
		resp.TodoDatasets = append(resp.TodoDatasets, store.ConvertTodoDataset(id, todo))
	}
	return resp
}

func (s *SyncerMock) Synchronize(currentVersion int, todos []store.TodoDateset) SynchronizeResponse {
	s.TStore.Sync(context.Background(), currentVersion, todos)
	return SynchronizeResponse{nil}
}
