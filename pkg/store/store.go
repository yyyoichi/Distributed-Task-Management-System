package store

import (
	"errors"
	"fmt"
	"sync"
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
	ByID        map[int]*Todo
}

type Todo struct {
	Task      string
	Completed bool
	Version   int
	Deleted   bool
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

// [version]よりも大きいバージョンを持つTODOを返す
func (s *TStore) GetLatestVersionTodo(version int) map[int]Todo {
	todos := map[int]Todo{}
	s.mu.Lock()
	for id, todo := range s.ByID {
		if version < todo.Version {
			todos[id] = *todo
		}
	}
	s.mu.Unlock()
	return todos
}

// 更新バージョンを同期する
func (s *TStore) SyncNextVersion(newVersion int) {
	s.nextVersion = newVersion
}

// データを上書きする。
func (s *TStore) SyncTodoAt(id int, newTodo Todo) error {
	s.mu.Lock()
	_, found := s.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	s.ByID[id] = &newTodo
	s.mu.Unlock()
	return nil
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
