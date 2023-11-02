package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/stream"
)

func NewStore() *TStore {
	return &TStore{
		mu:          sync.Mutex{},
		ByID:        make(map[int]*Todo),
		nextVersion: 1,
	}
}

type TStore struct {
	mu          sync.Mutex
	nextVersion int
	ByID        TodoKeyValueStore
}

func (s *TStore) Create(task string) int {
	s.mu.Lock()
	id := s.nextID()
	s.ByID[id] = &Todo{
		Task:      task,
		Completed: false,
		Deleted:   false,
		Version:   s.nextVersion,
	}
	s.nextVersion++
	s.mu.Unlock()
	return id
}

func (s *TStore) Update(id int, completed bool) error {
	s.mu.Lock()
	todo, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	todo.Completed = completed
	todo.Version = s.nextVersion
	s.nextVersion++
	s.mu.Unlock()
	return nil
}

func (s *TStore) Delete(id int) error {
	s.mu.Lock()
	todo, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	todo.Deleted = true
	todo.Version = s.nextVersion
	s.nextVersion++
	s.mu.Unlock()
	return nil
}

// [version]**以上**のバージョンを持つTODOを返す
func (s *TStore) GetLatestVersionTodo(version int) map[int]Todo {
	todos := map[int]Todo{}
	s.mu.Lock()
	for id, todo := range s.ByID {
		if version <= todo.Version {
			todos[id] = *todo
		}
	}
	s.mu.Unlock()
	return todos
}

// 同期を実行する
// [currentSyncVersion]今回の同期バージョン, [todos]同期するTodoDataset
func (s *TStore) Sync(cxt context.Context, currentSyncVersion int, todos []TodoDateset) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TASK.1 sync TODO
	todoCh := stream.Generator[TodoDateset](cxt, todos...)
	doneCh := stream.FunIO[TodoDateset, interface{}](cxt, todoCh, func(td TodoDateset) interface{} {
		todo := ConvertTodo(td)
		s.ByID[td.ID] = &todo
		return nil
	})
	// doneChが終わるまで待機
	for range doneCh {
	}
	// TASK.2 sync nextSyncVersion
	s.nextVersion = currentSyncVersion + 1
}

func (s *TStore) nextID() int {
	max := 0
	for id := range s.ByID {
		if max < id {
			max = id
		}
	}
	return max + 1
}
