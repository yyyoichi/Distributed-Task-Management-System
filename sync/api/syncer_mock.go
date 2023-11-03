package api

import (
	"context"
	"fmt"
	"log"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"
)

func NewSyncerMock(name string) *SyncerMock {
	return &SyncerMock{name, document.NewTDocument()}
}

type SyncerMock struct {
	name string
	*document.TDocument
}

func (s *SyncerMock) Me() string {
	return fmt.Sprintf("SyncerMock[%s]: Key-Value-Store has %d itmes", s.name, len(s.ByID))
}

func (s *SyncerMock) GetDifference(currentSyncVersion int) DiffResponse {
	resp := DiffResponse{}

	todos := s.TDocument.GetLatestVersionTodo(currentSyncVersion)
	log.Printf("\t[%s] There are %d Todos(/%d) defferenced from v%d", s.name, len(todos), len(s.TDocument.ByID), currentSyncVersion)
	for id, todo := range todos {
		resp.TodoDatasets = append(resp.TodoDatasets, document.ConvertTodoDataset(id, todo))
	}
	return resp
}

func (s *SyncerMock) Synchronize(currentSyncVersion int, todos []document.TodoDataset) SynchronizeResponse {
	s.TDocument.Synchronize(context.Background(), currentSyncVersion, todos)
	return SynchronizeResponse{nil}
}
