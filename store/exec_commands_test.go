package main

import (
	"errors"
	"testing"
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
)

func TestExecCommands(t *testing.T) {
	tStore := store.NewStore()

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
		resp, err := Exec(test.cmds, tStore)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("For commands %v, expected error: %v, but got: %v", test.cmds, test.err, err)
		}

		if resp != test.expected {
			t.Errorf("For commands %v, expected: %s, but got: %s", test.cmds, test.expected, resp)
		}
	}
}
