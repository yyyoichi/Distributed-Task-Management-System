package polling

import "github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"

type ConflictResolver struct {
	ToDoByID map[int]document.TodoDataset
}

// add tododataset
func (cr *ConflictResolver) AddToDo(td document.TodoDataset) {
	id := td.ID
	todo, found := cr.ToDoByID[id]
	// if the id is not found in resolver,
	// it is added.
	if !found {
		cr.ToDoByID[id] = td
		return
	}

	// if the id is found, the latest todo should be added.
	isInResolverIsLatest := todo.UpdatedAt.After(td.UpdatedAt)
	if isInResolverIsLatest {
		return
	}
	cr.ToDoByID[id] = td
}

func (cs *ConflictResolver) GetSlice() []document.TodoDataset {
	tds := []document.TodoDataset{}
	for _, todo := range cs.ToDoByID {
		tds = append(tds, todo)
	}
	return tds
}
