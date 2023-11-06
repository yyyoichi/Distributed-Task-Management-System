package polling

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"
	"github.com/yyyoichi/Distributed-Task-Management-System/sync/api"
)

// テストシナリオ...複数データノードのポーリング同期をテストする。※データノードはmockで代替する。
// 3つのデータノードを使用し、それぞれデータノードA,B ..と命名する。
func TestPollingManager(t *testing.T) {
	datanodeA := api.NewSyncerMock("DatanodeA")
	datanodeB := api.NewSyncerMock("DatanodeB")
	datanodeC := api.NewSyncerMock("DatanodeC")

	// test target call
	pollingManager := PollingManager{
		SyncerStore: NewSyncerStore(
			datanodeA,
			datanodeB,
			datanodeC,
		),
		nextSyncVersion: 1,
	}
	cxt := context.Background()

	// TestCase.1 //
	// INIT:
	// key-value-store of datanodeA has ...
	// - ID:1 v2 TaskA completed
	// and there are no keys and values in the key-value-store of datanodes B and C.
	// EXP:
	// ToDo[ID:1] has the compeletion flag is created for B and C,
	// and update version to 1.
	// key-value-store of datanodeA, B, C has ...
	// - ID:1 v1 TaskA completed

	// INIT:
	id1 := datanodeA.Create("TaskA")
	datanodeA.Update(id1, true)

	// Run:
	pollingManager.Polling(cxt)

	// TEST:
	log.Println("Start TestCase1 ..")
	expCase := map[int]document.Todo{}
	expCase[1] = document.Todo{
		Task:      "TaskA",
		Completed: true,
		Deleted:   false,
		Version:   1,
		UpdatedAt: time.Now().Add(time.Duration(time.Minute * 1)),
	}
	testNode(t, datanodeA, expCase)
	testNode(t, datanodeB, expCase)
	testNode(t, datanodeC, expCase)

	// TestCase.2 //
	// INIT:
	// key-value-store of datanodeA, B and C have ...
	// - ID:1 v1 TaskA completed
	// - (equal to list of expected TestCase1 result)
	// OPS:
	// datanodeA
	// - create 'TaskB'
	// datanodeB
	// - update ID:1 no-completed
	// EXP:
	// ToDo[ID:2] is created for B and C,
	// and lower the completion flag on A and C
	// and update the version to 2.
	// key-value-store of datanodeA, B, C has ...
	// - ID:1 v2 TaskA no-complete
	// - ID:2 v2 TaskB no-complete

	// OPS:
	datanodeA.Create("TaskB")
	datanodeB.Update(1, false)

	// Run:
	pollingManager.Polling(cxt)

	// TEST:
	log.Println("Start TestCase2 ..")
	expCase[1] = document.Todo{
		Task:      "TaskA",
		Completed: false,
		Deleted:   false,
		Version:   2,
		UpdatedAt: time.Now().Add(time.Duration(time.Minute * 1)),
	}
	expCase[2] = document.Todo{
		Task:      "TaskB",
		Completed: false,
		Deleted:   false,
		Version:   2,
		UpdatedAt: time.Now().Add(time.Duration(time.Minute * 1)),
	}
	testNode(t, datanodeA, expCase)
	testNode(t, datanodeB, expCase)
	testNode(t, datanodeC, expCase)
}

func testNode(t *testing.T, node *api.SyncerMock, expByID map[int]document.Todo) {
	for id, exp := range expByID {
		todo, found := node.TDocument.ByID[id]
		if !found {
			t.Errorf("Expected ToDo[ID:%d] is found, but it was not found in %s", id, node.Me())
		}
		if todo.Completed != exp.Completed {
			t.Errorf("Expected ToDo[ID:%d].Compeleted is %v, but got='%v' in %s", id, exp.Completed, todo.Completed, node.Me())
		}
		if todo.Version != exp.Version {
			t.Errorf("Expected ToDo[ID:%d].Version is %d, but got='%d' completed in %s", id, exp.Version, todo.Version, node.Me())
		}
	}
}
