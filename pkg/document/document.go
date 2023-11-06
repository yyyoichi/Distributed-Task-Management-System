package document

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/stream"
)

func NewTDocument() *TDocument {
	return &TDocument{
		mu:          sync.Mutex{},
		ByID:        make(map[int]*Todo),
		nextVersion: 1,
	}
}

type TDocument struct {
	mu          sync.Mutex
	nextVersion int
	ByID        map[int]*Todo
}

func (dc *TDocument) Create(task string) int {
	dc.mu.Lock()
	id := dc.nextID()
	dc.ByID[id] = &Todo{
		Task:      task,
		Completed: false,
		Deleted:   false,
		Version:   dc.nextVersion,
		UpdatedAt: time.Now(),
	}
	dc.nextVersion++
	dc.mu.Unlock()
	return id
}

func (dc *TDocument) Update(id int, completed bool) error {
	dc.mu.Lock()
	todo, found := dc.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	todo.Completed = completed
	todo.Version = dc.nextVersion
	todo.UpdatedAt = time.Now()
	dc.nextVersion++
	dc.mu.Unlock()
	return nil
}

func (dc *TDocument) Delete(id int) error {
	dc.mu.Lock()
	todo, found := dc.ByID[id]
	if !found {
		err := fmt.Sprintf("not found TODO[ID:%d]", id)
		return errors.New(err)
	}
	todo.Deleted = true
	todo.Version = dc.nextVersion
	dc.nextVersion++
	dc.mu.Unlock()
	return nil
}

// [version]**以上**のバージョンを持つTODOを返す
func (dc *TDocument) GetLatestVersionTodo(version int) map[int]Todo {
	todos := map[int]Todo{}
	dc.mu.Lock()
	for id, todo := range dc.ByID {
		if version <= todo.Version {
			todos[id] = *todo
		}
	}
	dc.mu.Unlock()
	return todos
}

// 同期を実行する
// [currentSyncVersion]今回の同期バージョン, [todos]同期するTodoDataset
func (dc *TDocument) Synchronize(cxt context.Context, currentSyncVersion int, todos []TodoDataset) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	// TASK.1 sync TODO
	todoCh := stream.Generator[TodoDataset](cxt, todos...)
	doneCh := stream.FunIO[TodoDataset, interface{}](cxt, todoCh, dc.synchronizeTodo)
	// doneChが終わるまで待機
	for range doneCh {
	}
	// TASK.2 sync nextSyncVersion
	dc.nextVersion = currentSyncVersion + 1
}

// must use in mu.Lock
func (dc *TDocument) synchronizeTodo(td TodoDataset) interface{} {
	synchronizaionToDo := ConvertTodo(td)
	todoInDocumnt, found := dc.ByID[td.ID]
	if !found {
		dc.ByID[td.ID] = &synchronizaionToDo
		return nil
	}
	// Conflict! ToDo[ID: id] exists in the document
	// Resolve conflicts at update time.
	// The latest ToDo must remain in the document.
	// If the sync todo equal to in the document, it should be synchronize. Becaouse it is todo in the same document.
	inDocumentToDoIsLatest := todoInDocumnt.UpdatedAt.After(synchronizaionToDo.UpdatedAt)
	if !inDocumentToDoIsLatest {
		*todoInDocumnt = synchronizaionToDo
	}
	return nil
}

func (dc *TDocument) nextID() int {
	max := 0
	for id := range dc.ByID {
		if max < id {
			max = id
		}
	}
	return max + 1
}
