package store

import (
	"context"
	"testing"
)

func TestTDocument_Create(t *testing.T) {
	tDocument := NewTDocument()

	id := tDocument.Create("Task 1")
	if id != 1 {
		t.Errorf("Expected ID: 1, but got: %d", id)
	}

	// Check if the task is created correctly
	todo, found := tDocument.ByID[id]
	if !found {
		t.Errorf("Expected to find TODO[ID:1], but it was not found")
	} else {
		if todo.Task != "Task 1" {
			t.Errorf("Expected task: 'Task 1', but got: %s", todo.Task)
		}
		if todo.Completed {
			t.Errorf("Expected completed: false, but got: true")
		}
		if todo.Deleted {
			t.Errorf("Expected deleted: false, but got: true")
		}
		if todo.Version != 1 {
			t.Errorf("Expected version: 1, but got: %d", todo.Version)
		}
	}
}

func TestTDocument_Update(t *testing.T) {
	tDocument := NewTDocument()
	id := tDocument.Create("Task 1")

	// Test updating with valid ID and status
	err := tDocument.Update(id, true)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the status is updated correctly
	todo, found := tDocument.ByID[id]
	if !found || !todo.Completed {
		t.Errorf("Expected completed: true, but got: false")
	}
	if todo.Version != 2 {
		t.Errorf("Expected version: 2, but got: %d", todo.Version)
	}

	// Test updating with invalid ID
	err = tDocument.Update(2, true)
	expectedErr := "not found TODO[ID:2]"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: '%s', but got: %v", expectedErr, err)
	}
}

func TestTDocument_Delete(t *testing.T) {
	tDocument := NewTDocument()
	id := tDocument.Create("Task 1")

	// Test deleting with valid ID
	err := tDocument.Delete(id)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the TODO is deleted correctly
	todo, found := tDocument.ByID[id]
	if !found {
		t.Errorf("Expected to find TODO[ID:1], but it was not found")
	}
	if !todo.Deleted {
		t.Errorf("Expected deleted: true, but got: false")
	}
	if todo.Version != 2 {
		t.Errorf("Expected version: 2, but got: %d", todo.Version)
	}

	// Test deleting with invalid ID
	err = tDocument.Delete(2)
	expectedErr := "not found TODO[ID:2]"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: '%s', but got: %v", expectedErr, err)
	}
}

func TestTDocument_GetLatestVersionTodo(t *testing.T) {
	test := []struct {
		todos          []string
		tagetVersion   int
		expectedLength int
	}{
		{
			todos:          []string{"TaskA", "TaskB", "TaskC"},
			tagetVersion:   3,
			expectedLength: 1,
		},
		{
			todos:          []string{"TaskA"},
			tagetVersion:   2,
			expectedLength: 0,
		},
	}

	for i, tt := range test {
		tDocument := NewTDocument()
		for _, todo := range tt.todos {
			_ = tDocument.Create(todo)
		}
		actTODO := tDocument.GetLatestVersionTodo(tt.tagetVersion)
		if len(actTODO) != tt.expectedLength {
			t.Errorf("%d: Expected len(actTODO) is '%d', but got='%d'", i, tt.expectedLength, len(actTODO))
		}
	}
}

func TestTDocument_SynchronizeTodoAt(t *testing.T) {
	tDocument := NewTDocument()
	id := tDocument.Create("TaskA")  // version 1
	tDocument.Update(id, true)       // version 2
	id2 := tDocument.Create("TaskB") // version 3
	// now
	// ID:1 version:2 TaskA completed
	// ID 2 version:3 TaskB no-complete

	// sync version 1
	syncVersion := 1
	// sync todo (update no-completed)
	todo := []TodoDateset{{
		ID:        1,
		Task:      "TaskA",
		Completed: false,
		Deleted:   false,
		Version:   1,
	}}
	// exec
	tDocument.Synchronize(context.Background(), syncVersion, todo)

	if tDocument.ByID[id].Completed {
		t.Error("Expected completed: false, but got: true")
	}
	if tDocument.ByID[id].Version != 1 {
		t.Errorf("Expected Version is 1, but got='%d'", tDocument.ByID[id].Version)
	}
	if tDocument.ByID[id2].Version != 3 {
		t.Errorf("Expected Version is 3, but got='%d'", tDocument.ByID[id].Version)
	}
}
