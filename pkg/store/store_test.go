package store

import (
	"testing"
)

func TestTStore_Create(t *testing.T) {
	tStore := NewStore()

	id := tStore.Create("Task 1")
	if id != 1 {
		t.Errorf("Expected ID: 1, but got: %d", id)
	}

	// Check if the task is created correctly
	todo, found := tStore.ByID[id]
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

func TestTStore_Update(t *testing.T) {
	tStore := NewStore()
	id := tStore.Create("Task 1")

	// Test updating with valid ID and status
	err := tStore.Update(id, true)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the status is updated correctly
	todo, found := tStore.ByID[id]
	if !found || !todo.Completed {
		t.Errorf("Expected completed: true, but got: false")
	}
	if todo.Version != 2 {
		t.Errorf("Expected version: 2, but got: %d", todo.Version)
	}

	// Test updating with invalid ID
	err = tStore.Update(2, true)
	expectedErr := "not found TODO[ID:2]"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: '%s', but got: %v", expectedErr, err)
	}
}

func TestTStore_Delete(t *testing.T) {
	tStore := NewStore()
	id := tStore.Create("Task 1")

	// Test deleting with valid ID
	err := tStore.Delete(id)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the TODO is deleted correctly
	todo, found := tStore.ByID[id]
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
	err = tStore.Delete(2)
	expectedErr := "not found TODO[ID:2]"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: '%s', but got: %v", expectedErr, err)
	}
}

func TestTStore_GetLatestVersionTodo(t *testing.T) {
	test := []struct {
		todos          []string
		tagetVersion   int
		expectedLength int
	}{
		{
			todos:          []string{"TaskA", "TaskB", "TaskC"},
			tagetVersion:   2,
			expectedLength: 1,
		},
		{
			todos:          []string{"TaskA"},
			tagetVersion:   1,
			expectedLength: 0,
		},
	}

	for i, tt := range test {
		tStore := NewStore()
		for _, todo := range tt.todos {
			_ = tStore.Create(todo)
		}
		actTODO := tStore.GetLatestVersionTodo(tt.tagetVersion)
		if len(actTODO) != tt.expectedLength {
			t.Errorf("%d: Expected len(actTODO) is '%d', but got='%d'", i, tt.expectedLength, len(actTODO))
		}
	}
}

func TestTStore_SyncTodoAt(t *testing.T) {
	tStore := NewStore()
	id := tStore.Create("TaskA")

	// sync
	todo := Todo{
		Task:      "Task-1st",
		Completed: false,
		Deleted:   false,
		Version:   100,
	}
	if err := tStore.SyncTodoAt(id, todo); err != nil {
		t.Error(err)
	}

	if tStore.ByID[id].Task != "Task-1st" {
		t.Errorf("Expected TODO.task 'Task-1st', but got='%s'", tStore.ByID[id].Task)
	}

	// error
	if err := tStore.SyncTodoAt(9999, todo); err == nil {
		t.Errorf("Expected error is not nil, but got='nil'")
	}
}

func TestTStore_SyncNextVersion(t *testing.T) {
	tStore := NewStore()
	tStore.SyncNextVersion(100)
	if tStore.nextVersion != 100 {
		t.Errorf("Expected nextVersion is 100, but got='%d'", tStore.nextVersion)
	}
}
