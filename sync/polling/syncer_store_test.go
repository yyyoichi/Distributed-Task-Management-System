package polling

import (
	"context"
	"testing"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

func TestSyncerStore(t *testing.T) {
	datanodeA := api.NewSyncerMock("DatanodeA")
	datanodeB := api.NewSyncerMock("DatanodeB")
	store := NewSyncerStore(datanodeA, datanodeB)
	datanodeA.Create("TaskA")
	datanodeA.Create("TaskB")

	cxt := context.Background()
	id := 0
	num := 0
	for detector := range store.getDifferenceDetectorCh(cxt, 0) {
		id += detector.SyncerID
		num++
	}
	if id != 1 {
		t.Errorf("Expected IDs are '0' and '1', but got='%d' and '%d'", id, id)
	}
	if num != 2 {
		t.Errorf("Expected channel length is 2, but got='%d'", num)
	}

}
