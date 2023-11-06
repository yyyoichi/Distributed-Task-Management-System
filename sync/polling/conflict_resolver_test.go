package polling

import (
	"testing"
	"time"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"
)

func TestConflictResolver(t *testing.T) {
	dataset := make(map[int]document.TodoDataset)
	dataset[1] = document.TodoDataset{Version: 1, UpdatedAt: time.Now()}
	dataset[2] = document.TodoDataset{Version: 1, UpdatedAt: time.Now().Add(time.Duration(time.Minute * 10))}

	conflictResolver := ConflictResolver{dataset}
	conflictResolver.AddToDo(document.TodoDataset{
		ID:        1,
		Task:      "ID:1 must be updated",
		Version:   100,
		UpdatedAt: time.Now().Add(time.Duration(time.Minute * 10)),
	})
	conflictResolver.AddToDo(document.TodoDataset{
		ID:        2,
		Task:      "ID:1 must not be updated",
		Version:   100,
		UpdatedAt: time.Now(),
	})

	if conflictResolver.ToDoByID[1].Version != 100 {
		t.Errorf("Expected ID:1 to be updated, but it is not")
	}
	if conflictResolver.ToDoByID[2].Version != 1 {
		t.Errorf("Expected ID:2 not to be updated, but it is")
	}
}
