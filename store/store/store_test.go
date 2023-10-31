package store

import (
	"errors"
	"testing"
)

func TestTStore_Read(t *testing.T) {
	tStore := NewStore()

	tests := []struct {
		cmds     []string
		expected string
		err      error
	}{
		// 正しい "create" コマンド
		{cmds: []string{"create", "Task 1"}, expected: "Created TODO[ID:1]", err: nil},

		// 正しい "list" コマンド
		{cmds: []string{"list"}, expected: "TODO: \n1: Task 1\nCOMPLETED TODO:", err: nil},

		// 正しい "update" コマンド
		{cmds: []string{"update", "1", "complete"}, expected: "Updated TODO[ID:1]", err: nil},

		// 正しい "delete" コマンド
		{cmds: []string{"delete", "1"}, expected: "Deleted TODO[ID:1]", err: nil},

		// 正しい "help" コマンド
		{cmds: []string{"help"}, expected: "create <task>: Create a new todo with the specified task.\nlist: List all todos, separated into incomplete and completed todos.\nupdate <id> <status>: Update the status of the todo with the specified id. Status can be 'complete' or 'open'.\ndelete <id>: Delete the todo with the specified id.", err: nil},

		// 不正なコマンド
		{cmds: []string{"invalidCmd"}, expected: "", err: errors.New("syntax error: invalid comand: 'invalidCmd'")},
	}

	for _, test := range tests {
		resp, err := tStore.Read(test.cmds)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("For commands %v, expected error: %v, but got: %v", test.cmds, test.err, err)
		}

		if resp != test.expected {
			t.Errorf("For commands %v, expected: %s, but got: %s", test.cmds, test.expected, resp)
		}
	}
}

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
		if todo.task != "Task 1" {
			t.Errorf("Expected task: 'Task 1', but got: %s", todo.task)
		}
		if todo.completed {
			t.Errorf("Expected completed: false, but got: true")
		}
		if todo.deleted {
			t.Errorf("Expected deleted: false, but got: true")
		}
		if todo.version != 1 {
			t.Errorf("Expected version: 1, but got: %d", todo.version)
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
	if !found || !todo.completed {
		t.Errorf("Expected completed: true, but got: false")
	}
	if todo.version != 2 {
		t.Errorf("Expected version: 2, but got: %d", todo.version)
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
	if !todo.deleted {
		t.Errorf("Expected deleted: true, but got: false")
	}
	if todo.version != 2 {
		t.Errorf("Expected version: 2, but got: %d", todo.version)
	}

	// Test deleting with invalid ID
	err = tStore.Delete(2)
	expectedErr := "not found TODO[ID:2]"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error: '%s', but got: %v", expectedErr, err)
	}
}
