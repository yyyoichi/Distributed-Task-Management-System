package main

import (
	"sync"
	"testing"
)

func TestCreateTODO(t *testing.T) {
	test := []struct {
		tasks       []string
		expectedLen int
	}{
		{[]string{"TODO: A", "TODO: B", "TODO: C"}, 3},
		{[]string{"TODO: D", "TODO: E"}, 2},
	}
	for i, tt := range test {
		store := NewStore()
		var wg sync.WaitGroup
		wg.Add(len(tt.tasks))
		for _, task := range tt.tasks {
			go func(t string) {
				defer wg.Done()
				store.Create(t)
			}(task)
		}
		wg.Wait()
		if len(store.ByID) != tt.expectedLen {
			t.Errorf("%d:Expected length is %d, but got='%d'", i, tt.expectedLen, len(store.ByID))
		}
	}
}
