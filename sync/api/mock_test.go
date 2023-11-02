package api

import "testing"

func TestMock(t *testing.T) {
	datanodeA := NewSyncerMock("DatanodeA")
	datanodeA.Create("TaskA")
	datanodeA.Create("TaskB")

	diff := datanodeA.GetDifference(1)
	if len(diff.TodoDatasets) != 2 {
		t.Errorf("Expected length is 2, but got='%d'", len(diff.TodoDatasets))
	}
}
